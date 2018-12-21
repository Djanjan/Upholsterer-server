package planner

import (
	"bufio"
	"encoding/base64"
	"image/jpeg"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ivahaev/go-logger"
	"github.com/nfnt/resize"
)

type ResizeMaster struct {
	date      []Date
	col_tasks int
	start     bool
}

func (res ResizeMaster) init() {
	res.col_tasks = 0
	res.start = false
}

func (res ResizeMaster) run(dateOut chan Date) {
	logger.Notice("ResizeMaster -- Start")
	res.start = true

	defer close(dateOut)

	dateResized := make(chan Date, MAXDOWLOADS)

	for _, item := range res.date {
		go resizeFileImg(item, dateResized)
		logger.Info("ResizeMaster -- Gorutina start: " + item.ImgsURL)
	}

	for _, i := range res.date {
		item := <-dateResized
		logger.Info("ResizeMaster -- Gorutina close: " + item.ImgsOriginPath + " -- " + i.ImgsURL)

		logger.Info("ResizeMaster -- patok to dataout: " + item.ImgsOriginPath + " -- " + i.ImgsURL)
		dateOut <- item
	}
	close(dateResized)

	logger.Notice("ResizeMaster -- END")
	res.start = false
}

func resizeFileImg(dateInput Date, dateOut chan Date) {
	logger.Info("ResizeMaster -- item resize: "+dateInput.ImgsOriginPath, nil)
	imgIn, _ := os.Open(dateInput.ImgsOriginPath)
	defer imgIn.Close()
	imgJpg, err := jpeg.Decode(imgIn)
	if err != nil {
		logger.Error("ResizeMaster -- Error img decoder:  ", err)
		dateOut <- dateInput
		return
	}

	imgJpg = resize.Resize(600, 400, imgJpg, resize.Bicubic)

	pathSplit := strings.Split(dateInput.ImgsOriginPath, "/")
	newPathImag := "img/icon/" + pathSplit[2]

	imgOut, err := os.Create(newPathImag)
	if err != nil {
		logger.Error("ResizeMaster -- Error img create:  ", err)
		dateOut <- dateInput
		return
	}

	err = jpeg.Encode(imgOut, imgJpg, nil)
	if err != nil {
		logger.Error("ResizeMaster -- Error img encode:  ", err)
		dateOut <- dateInput
		return
	}

	imgOut.Close()

	dateInput.ImgsIconPath = newPathImag
	dateOut <- dateInput
}

func resizeByteImg(res io.Reader) {
	/*imgIn, _ := os.Open("more-galka-skala-zaplyv.jpg")
	  defer imgIn.Close()
	  //imgJpg, _ := jpeg.Decode(imgIn)
	  imgJpg, err := jpeg.Decode(imgIn)
	  if err != nil {
	  	fmt.Println("Error img decoder:  ", err)
	  }*/

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(encodeBase64(res)))
	imgJpg, err := jpeg.Decode(reader)
	if err != nil {
		logger.Error("Error img decoder:  ", err)
	}

	imgJpg = resize.Resize(600, 0, imgJpg, resize.Bicubic) // <-- Собственно изменение размера картинки

	imgOut, err := os.Create("img/icon/test-out.jpg")
	if err != nil {
		logger.Error("Error img create:  ", err)
	}

	err = jpeg.Encode(imgOut, imgJpg, nil)
	if err != nil {
		logger.Error("Error img encode:  ", err)
	}

	imgOut.Close()
}

func encodeBase64(res io.Reader) string {
	// Read entire JPG into byte slice.
	reader := bufio.NewReader(res)
	content, _ := ioutil.ReadAll(reader)

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)
	if encoded == "" {
		logger.Error("Encoder error")
	}
	// Print encoded data to console.
	// ... The base64 image can be used as a data URI in a browser.
	return encoded
}

func encodeByte(res io.Reader) []byte {
	// Read entire JPG into byte slice.
	reader := bufio.NewReader(res)
	content, _ := ioutil.ReadAll(reader)

	return content
}
