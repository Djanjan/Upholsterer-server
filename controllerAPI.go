package main

import (
	"database/sql"
	"go_service/planner"
	"os"
	"strconv"

	"github.com/ivahaev/go-logger"
)

type ControllerAPI struct {
	plan   planner.Planner
	urlImg []string
}

func (c *ControllerAPI) init() {
	c.createDir()

	c.plan.Init()
}

func (c *ControllerAPI) addTask(imgs []Image, typeTask string) {
	var dateInput []planner.Date
	for _, item := range imgs {
		img := planner.Date{
			ImgsURL:        item.Url,
			ImgsOriginPath: item.OriginPath,
			ImgsIconPath:   item.IconPath,
		}

		dateInput = append(dateInput, img)
	}

	var tas planner.Task
	tas.Init()
	if typeTask == "Dowload" {
		tas.TypeTask = planner.Dowload
	} else {
		tas.TypeTask = planner.Resize
	}
	tas.DateInput = dateInput
	tas.Name = "Task Controller Api -- " + typeTask

	c.plan.AddTask(&tas)
}

func (c *ControllerAPI) runAllTaskType(typeTask string) {
	dateC := make(chan []planner.Date)
	var img []planner.Date
	if typeTask == "Dowload" {
		go c.plan.RunAllTaskType(dateC, planner.Dowload)
	} else {
		go c.plan.RunAllTaskType(dateC, planner.Resize)
	}

	img = <-dateC
	for _, value := range img {
		logger.Info("Task: date -- " + value.ImgsOriginPath + value.ImgsIconPath)
		UpdateUrlDB(value.ImgsURL, Image{
			OriginPath: value.ImgsOriginPath,
			IconPath:   value.ImgsIconPath,
		})
	}

	c.plan.PoolsTaskComplite()
}

func (c *ControllerAPI) runAllTaskTypeA(typeTask string) {
	dateC := make(chan []planner.Date)
	//var img []planner.Date
	if typeTask == "Dowload" {
		go c.plan.RunAllTaskType(dateC, planner.Dowload)
	} else {
		go c.plan.RunAllTaskType(dateC, planner.Resize)
	}
	<-dateC
	c.getDateCompliteTask()
	/*for _, value := range img {
		fmt.Println("Task: date -- " + value.ImgsOriginPath)
		UpdateUrlDB(value.ImgsURL, Image{
			OriginPath: value.ImgsOriginPath,
			IconPath:   value.ImgsIconPath,
		})
	}*/
}

func (c *ControllerAPI) getDateCompliteTask() {
	var massTas []*planner.Task

	massTas = c.plan.PoolsTaskComplite()
	//fmt.Println(len(massTas))

	for key, item := range massTas {
		for _, value := range item.DateOut {
			logger.Info("Task: key -- " + strconv.Itoa(key))
			UpdateUrlDB(value.ImgsURL, Image{
				OriginPath: value.ImgsOriginPath,
				IconPath:   value.ImgsIconPath,
			})
		}
	}

	/*for _, item := range c.plan.Tasks {
		if item.Done {
			for _, value := range item.DateOut {
				fmt.Println("Task: date -- " + value.ImgsOriginPath)
				UpdateUrlDB(value.ImgsURL, Image{
					OriginPath: value.ImgsOriginPath,
					IconPath:   value.ImgsIconPath,
				})
			}
		}
	}*/
}

//rows, err := db.Query("select * from Images where iconPath IS NULL OR iconPath = ''")

func (c *ControllerAPI) getImagsNoneIcon() ([]Image, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from Images where iconPath IS NULL OR iconPath = ''")
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	imgs := []Image{}

	for rows.Next() {
		p := ImageNil{}
		err := rows.Scan(&p.ID, &p.Url, &p.Catalog, &p.OriginPath, &p.IconPath)
		if err != nil {
			logger.Error("Start -- PlannerAPI -- " + err.Error())
			continue
		}
		imgs = append(imgs, p.toImage())
	}

	return imgs, err
}

func (c *ControllerAPI) getImagsNoneOriginPath() ([]Image, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from Images where originPath IS NULL OR originPath = ''")
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	imgs := []Image{}

	for rows.Next() {
		p := ImageNil{}
		err := rows.Scan(&p.ID, &p.Url, &p.Catalog, &p.OriginPath, &p.IconPath)
		if err != nil {
			logger.Error("Start -- PlannerAPI -- " + err.Error())
			continue
		}
		imgs = append(imgs, p.toImage())
	}

	return imgs, err
}

func (c *ControllerAPI) startRejection() {
	c.checkRejection()

	errImgs, err := GetsErrorDB()
	if err != nil {
		logger.Error("ControllerApi -- SQLITE -- " + err.Error())
	}

	for _, item := range errImgs {
		var newPriority = (item.Count * item.Priority) / 2
		if newPriority < 1 {
			newPriority = 1
		}

		if newPriority >= 5 {
			_, err = DeleteUrlDB(item.Url)
			if err != nil {
				logger.Error("ControllerApi -- SQLITE -- " + err.Error())
			}

			logger.Info("ControllerApi -- Rejection imag -- ", item.Url)

			_, err = DeleteErrorDB(item.Url)
			if err != nil {
				logger.Error("ControllerApi -- SQLITE -- " + err.Error())
			}
		}

		_, err = UpdateErrorDB(item.Url, ErrorImage{
			Url:      item.Url,
			Count:    item.Count,
			Priority: newPriority,
		})
		if err != nil {
			logger.Error("ControllerApi -- SQLITE -- " + err.Error())
		}
	}
}

func (c *ControllerAPI) checkRejection() {
	var errImgs = c.plan.PoolErrorImages()
	for _, item := range errImgs {
		_, err := AddErrorDB(ErrorImage{
			Url:      item,
			Count:    0,
			Priority: 1,
		})
		if err != nil {
			errImg, err := GetErrorDB(item)
			if err != nil {
				logger.Error("ControllerApi -- SQLITE -- " + err.Error())
			}
			//logger.Debug("ControllerApi -- SQLITE -- ", errImg)
			_, err = UpdateErrorDB(errImg.Url, ErrorImage{
				Url:      errImg.Url,
				Count:    errImg.Count + 1,
				Priority: errImg.Priority,
			})
			if err != nil {
				logger.Error("ControllerApi -- SQLITE -- " + err.Error())
			}
		}
	}
}

func (c *ControllerAPI) createDir() {
	if err := os.MkdirAll("img", 666); err != nil {
		panic(err)
	}
	if err := os.MkdirAll("img/origin", 666); err != nil {
		panic(err)
	}
	if err := os.MkdirAll("img/icon", 666); err != nil {
		panic(err)
	}
}
