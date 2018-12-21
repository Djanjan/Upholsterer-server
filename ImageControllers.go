package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	logger "github.com/ivahaev/go-logger"
	_ "github.com/mattn/go-sqlite3"
)

var images []Image

func getImages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	images, err := GetsDB(true)
	if err != nil {
		w.WriteHeader(http.StatusFound)
		logger.Error("GET Request", "Func getImages")
		return
	}
	logger.Info("GET Request", "Func getImages")
	json.NewEncoder(w).Encode(images)
	//w.WriteHeader(http.StatusOK)
}

func getImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	params := mux.Vars(r)
	image, err := GetDB(params["id"], true)
	if err != nil {
		w.WriteHeader(http.StatusFound)
		logger.Error("GET Request", "Func getImages")
		return
	}
	logger.Info("GET Request", "Func getImage")
	json.NewEncoder(w).Encode(image)
}

func getImageCatalog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	filter := r.FormValue("filter")
	var catalogs []Image
	if filter == "all" {
		catalogs, err := GetCatalogDB(true)
		if err != nil {
			w.WriteHeader(http.StatusFound)
			logger.Error("GET Request", "Func getImageCatalog")
			return
		}
		logger.Info("GET Request", "Func getImageCatalog")
		json.NewEncoder(w).Encode(catalogs)
		return
	}

	catalogs, err := GetImageCatalogDB(filter, true)
	if err != nil {
		w.WriteHeader(http.StatusFound)
		logger.Error("GET Request", "Func getImageCatalog")
		return
	}
	json.NewEncoder(w).Encode(catalogs)
	//w.WriteHeader(http.StatusOK)
}

func createImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var imags []Image
	_ = json.NewDecoder(r.Body).Decode(&imags)

	for _, item := range imags {
		_, err := AddDB(item)
		if err != nil {
			logger.Error(err)
		}
	}

	json.NewEncoder(w).Encode(http.StatusOK)
}

func updateImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	params := mux.Vars(r)
	var imag Image
	_ = json.NewDecoder(r.Body).Decode(&imag)
	_, err := UpdateDB(params["id"], imag)
	if err != nil {
		w.WriteHeader(http.StatusFound)
		return
	}
	json.NewEncoder(w).Encode(http.StatusOK)
}

func deleteImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	params := mux.Vars(r)
	_, err := DeleteDB(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusFound)
		return
	}
	json.NewEncoder(w).Encode(http.StatusOK)
}

func queriImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	queri := r.FormValue("queri")
	switch queri {

	case "deleteAllImages":
		_, err := DeleteAllDB()
		if err != nil {
			w.WriteHeader(http.StatusFound)
			return
		}
	}
	json.NewEncoder(w).Encode(http.StatusOK)
}
