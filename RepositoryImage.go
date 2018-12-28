package main

import (
	"database/sql"
	"errors"
	"os"

	"github.com/ivahaev/go-logger"
	_ "github.com/mattn/go-sqlite3"
)

func GetsDB(isNULL bool) ([]Image, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var rows *sql.Rows
	if isNULL {
		rows, err = db.Query("select * from Images")
	} else {
		rows, err = db.Query("select * from Images where iconPath <> '' and iconPath is not null and originPath <> '' and originPath is not null ")
	}
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	imgs := []Image{}

	for rows.Next() {
		p := ImageNil{}
		err := rows.Scan(&p.ID, &p.Url, &p.Catalog, &p.OriginPath, &p.IconPath)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		imgs = append(imgs, p.toImage())

	}

	return imgs, err
}

func GetsCountDB(isNULL bool, count int) ([]Image, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var rows *sql.Rows
	if isNULL {
		rows, err = db.Query("select * from Images")
	} else {
		rows, err = db.Query("select * from Images where iconPath <> '' and iconPath is not null and originPath <> '' and originPath is not null ")
	}
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	imgs := []Image{}

	thisCurect := count
	for rows.Next() {
		p := ImageNil{}
		err := rows.Scan(&p.ID, &p.Url, &p.Catalog, &p.OriginPath, &p.IconPath)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		imgs = append(imgs, p.toImage())

		thisCurect--
		if thisCurect <= 0 {
			break
		}
	}

	return imgs, err
}

func GetDB(id string, isNULL bool) (Image, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var rows *sql.Rows
	if isNULL {
		rows, err = db.Query("select * from Images where id = $1", id)
	} else {
		rows, err = db.Query("select * from Images where id = $1 and iconPath <> '' and iconPath is not null and originPath <> '' and originPath is not null ", id)
	}
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	imgs := Image{}

	for rows.Next() {
		p := ImageNil{}
		err := rows.Scan(&p.ID, &p.Url, &p.Catalog, &p.OriginPath, &p.IconPath)
		if err != nil {
			logger.Error(err)
			continue
		}
		imgs = p.toImage()

	}

	return imgs, err
}

func GetUrlDB(url string) (Image, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var rows *sql.Rows
	rows, err = db.Query("select * from Images where url = $1", url)

	if err != nil {
		panic(err)
	}
	defer rows.Close()
	imgs := Image{}

	for rows.Next() {
		p := ImageNil{}
		err := rows.Scan(&p.ID, &p.Url, &p.Catalog, &p.OriginPath, &p.IconPath)
		if err != nil {
			logger.Error(err)
			continue
		}
		imgs = p.toImage()

	}

	return imgs, err
}

func GetImageCatalogDB(catalog string, isNULL bool) ([]Image, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var rows *sql.Rows
	if isNULL {
		rows, err = db.Query("select * from Images where catalog=$1", catalog)
	} else {
		rows, err = db.Query("select * from Images where catalog=$1 and iconPath <> '' and iconPath is not null and originPath <> '' and originPath is not null ", catalog)
	}

	if err != nil {
		panic(err)
	}
	defer rows.Close()
	imgs := []Image{}

	for rows.Next() {
		p := ImageNil{}
		err := rows.Scan(&p.ID, &p.Url, &p.Catalog, &p.OriginPath, &p.IconPath)
		if err != nil {
			logger.Error(err)
			continue
		}
		imgs = append(imgs, p.toImage())
	}

	return imgs, err
}

func GetCatalogDB(isNULL bool) ([]Image, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var rows *sql.Rows
	if isNULL {
		rows, err = db.Query("select * from Images group by catalog")
	} else {
		rows, err = db.Query("select * from Images where iconPath <> '' and iconPath is not null and originPath <> '' and originPath is not null  group by catalog")
	}
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	imgs := []Image{}

	for rows.Next() {
		p := ImageNil{}
		err := rows.Scan(&p.ID, &p.Url, &p.Catalog, &p.OriginPath, &p.IconPath)
		if err != nil {
			logger.Error(err)
			continue
		}
		imgs = append(imgs, p.toImage())
	}

	return imgs, err
}

func AddDB(img Image) (int64, error) {
	if img.Url == "" || img.Catalog == "" {
		return 0, errors.New("img not found")
		//panic("img not found!")
	}

	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("insert into Images (ID, Url, Catalog) values ($1, $2, $3)",
		uuID(), img.Url, img.Catalog)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func UpdateDB(id string, imag Image) (int64, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("update Images set url = $1, catalog = $2, originPath = $3, IconPath = $4 where id = $5", imag.Url, imag.Catalog, imag.OriginPath, imag.IconPath, id)
	if err != nil {
		panic(err)
	}

	return result.LastInsertId()
}

func UpdateUrlDB(url string, imag Image) (int64, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("update Images set originPath = $1, iconPath = $2 where url = $3", imag.OriginPath, imag.IconPath, url)
	if err != nil {
		panic(err)
	}

	return result.LastInsertId()
}

func DeleteDB(id string) (int64, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("delete from Images where id = $1", id)
	if err != nil {
		panic(err)
	}

	return result.LastInsertId()
}

func DeleteUrlDB(url string) (int64, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	img, err := GetUrlDB(url)
	if err != nil {
		panic(err)
	}

	result, err := db.Exec("delete from Images where url = $1", url)
	if err != nil {
		panic(err)
	}

	if img.OriginPath != "" {
		err = os.Remove(img.OriginPath)
		if err != nil {
			panic(err)
		}
	}
	if img.IconPath != "" {
		err = os.Remove(img.IconPath)
		if err != nil {
			panic(err)
		}
	}

	return result.LastInsertId()
}

func DeleteCatalogDB(catalog string) (int64, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("delete from Images where catalog=$1", catalog)
	if err != nil {
		panic(err)
	}

	return result.LastInsertId()
}

func DeleteAllDB() (int64, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("delete from Images")
	if err != nil {
		panic(err)
	}

	return result.LastInsertId()
}
