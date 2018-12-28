package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/ivahaev/go-logger"
	_ "github.com/mattn/go-sqlite3"
)

var Planner = ControllerAPI{}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index/index.html")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := "Ухади"
	tmpl, _ := template.ParseFiles("index/about.html")
	tmpl.Execute(w, data)
}

func callAt(hour, min, sec int, f func()) error {
	logger.Notice("TIMER -- START")

	go func() {
		time.Sleep(time.Second * 5)
		for {
			f()
			logger.Info("TIMER -- SLEEP")
			time.Sleep(time.Hour * 10)
		}
	}()

	return nil
}

func planner_check() {
	logger.Notice("TIMER -- start task")
	logger.Info("TIMER -- Planner task count = " + strconv.Itoa(Planner.plan.GetLength()))
	imgsO, _ := Planner.getImagsNoneOriginPath()
	Planner.addTask(imgsO, "Dowload")
	Planner.runAllTaskType("Dowload")

	imgsI, _ := Planner.getImagsNoneIcon()
	Planner.addTask(imgsI, "Resize")
	Planner.runAllTaskType("Resize")

	//Rejection image
	Planner.startRejection()
}

func main() {
	logger.SetLevel("DEBUG")
	Planner.init()

	router := mux.NewRouter()
	/*router.HandleFunc("/api/v1/images", getImages).Methods("GET")
	router.HandleFunc("/api/v1/image/id/{id}", getImage).Methods("GET")
	router.HandleFunc("/api/v1/image/catalog", getImageCatalog).Queries("filter", "{filter}").Methods("GET")
	router.HandleFunc("/api/v1/image", queriImage).Queries("queri", "{queri}").Methods("GET")
	*/router.HandleFunc("/api/v1/image", createImage).Methods("POST")
	/*router.HandleFunc("/api/v1/image/{id}", updateImage).Methods("PUT")
	router.HandleFunc("/api/v1/image/{id}", deleteImage).Methods("DELETE")*/

	router.HandleFunc("/api/v1/images", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	router.HandleFunc("/about", aboutHandler)
	router.HandleFunc("/", homeHandler)

	//Absolute path Images
	http.Handle("/img/origin/", http.StripPrefix("/img/origin/", http.FileServer(http.Dir("./img/origin/"))))
	http.Handle("/img/icon/", http.StripPrefix("/img/icon/", http.FileServer(http.Dir("./img/icon/"))))

	http.Handle("/", router)

	err := callAt(0, 0, 0, planner_check)
	if err != nil {
		logger.Error("error: %v\n", err)
	}

	logger.Notice("Server is listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
