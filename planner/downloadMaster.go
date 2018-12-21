package planner

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ivahaev/go-logger"
)

type DownloadMaster struct {
	date      []Date
	col_tasks int
	start     bool
}

func (downl DownloadMaster) init() {
	downl.col_tasks = 0
	downl.start = false
}

func (downl DownloadMaster) run(dateOut chan Date) {
	logger.Notice("DownloadMaster -- Start")
	downl.start = true

	defer close(dateOut)

	dateDowload := make(chan Date, MAXDOWLOADS)

	//go allDowload(p.date, dateDowload)

	for _, item := range downl.date {
		go download(item, dateDowload)
		logger.Info("DownloadMaster -- Gorutina start: " + item.ImgsURL)
	}

	for _, i := range downl.date {
		item := <-dateDowload
		logger.Info("DownloadMaster -- Gorutina close: " + item.ImgsOriginPath + " -- " + i.ImgsURL)

		logger.Info("DownloadMaster -- patok to dataout: " + item.ImgsOriginPath + " -- " + i.ImgsURL)
		dateOut <- item
	}
	close(dateDowload)

	logger.Notice("DownloadMaster -- END")
	downl.start = false
}

/*func download(url string) (Date, error) {
	var date Date
	var err error

	fileName := "img" + "/origin/" + url[strings.LastIndex(url, "/")+1:]
	output, err := os.Create(fileName)
	if err != nil {
		err = errors.New("DownloadMaster -- Error create file")
		return date, err
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		err = errors.New("DownloadMaster -- Error while downloading" + url + " - " + err.Error())
		return date, err
	}
	defer response.Body.Close()

	//resizeImg(response.Body)
	//fmt.Println(encodeBase64(response.Body))
	//fmt.Println(encodeByte(response.Body))

	i, err := io.Copy(output, response.Body)
	if err != nil {
		err = errors.New("DownloadMaster -- Error copy new image" + url + " - " + err.Error())
		return date, err
	}

	i = i + 1

	return Date{
		imgsURL:        url,
		imgsOriginPath: fileName,
	}, nil
}*/

func download(dateInput Date, dateOut chan Date) {
	logger.Info("DownloadMaster -- item downoload: " + dateInput.ImgsURL)
	fileName := "img" + "/origin/" + dateInput.ImgsURL[strings.LastIndex(dateInput.ImgsURL, "/")+1:]
	output, err := os.Create(fileName)
	if err != nil {
		logger.Error("DownloadMaster -- Error create file")
		dateOut <- Date{
			ImgsURL:        dateInput.ImgsURL,
			ImgsOriginPath: "",
			ImgsIconPath:   dateInput.ImgsIconPath,
		}
		return
	}
	defer output.Close()

	response, err := http.Get(dateInput.ImgsURL)
	if err != nil {
		logger.Error("DownloadMaster -- Error while downloading" + dateInput.ImgsURL + " - " + err.Error())
		dateOut <- Date{
			ImgsURL:        dateInput.ImgsURL,
			ImgsOriginPath: "",
			ImgsIconPath:   dateInput.ImgsIconPath,
		}
		return
	}
	defer response.Body.Close()

	i, err := io.Copy(output, response.Body)
	if err != nil {
		logger.Error("DownloadMaster -- Error copy new image" + dateInput.ImgsURL + " - " + err.Error())
		dateOut <- Date{
			ImgsURL:        dateInput.ImgsURL,
			ImgsOriginPath: "",
			ImgsIconPath:   dateInput.ImgsIconPath,
		}
		return
	}
	//KEK
	i = i + 1

	dateOut <- Date{
		ImgsURL:        dateInput.ImgsURL,
		ImgsOriginPath: fileName,
		ImgsIconPath:   dateInput.ImgsIconPath,
	}
}
