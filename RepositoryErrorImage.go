package main

import (
	"database/sql"
	"errors"

	"github.com/ivahaev/go-logger"
	_ "github.com/mattn/go-sqlite3"
)

func GetsErrorDB() ([]ErrorImage, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var rows *sql.Rows
	rows, err = db.Query("select * from ErrorImages")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	imgs := []ErrorImage{}

	for rows.Next() {
		p := ErrorImage{}
		err := rows.Scan(&p.Url, &p.Count, &p.Priority)
		if err != nil {
			logger.Error(err)
			continue
		}
		imgs = append(imgs, p)
	}

	return imgs, err
}

func GetErrorDB(url string) (ErrorImage, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var rows *sql.Rows
	rows, err = db.Query("select * from ErrorImages where Url = $1", url)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	imgs := ErrorImage{}

	for rows.Next() {
		p := ErrorImage{}
		err := rows.Scan(&p.Url, &p.Count, &p.Priority)
		if err != nil {
			logger.Error(err)
			continue
		}
		imgs = p
	}

	return imgs, err
}

func AddErrorDB(img ErrorImage) (int64, error) {
	if img.Url == "" {
		return 0, errors.New("img not found")
	}

	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("insert into ErrorImages (Url, Count, Priority) values ($1, $2, $3)",
		img.Url, img.Count, img.Priority)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func UpdateErrorDB(url string, imag ErrorImage) (int64, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("update ErrorImages set Count = $1, Priority = $2 where url = $3", imag.Count, imag.Priority, url)
	if err != nil {
		panic(err)
	}

	return result.LastInsertId()
}

func DeleteErrorDB(url string) (int64, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("delete from ErrorImages where url = $1", url)
	if err != nil {
		panic(err)
	}

	return result.LastInsertId()
}

func DeleteErrorAllDB() (int64, error) {
	db, err := sql.Open("sqlite3", "imagedb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("delete from ErrorImages")
	if err != nil {
		panic(err)
	}

	return result.LastInsertId()
}
