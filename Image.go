package main

import "database/sql"

type Image struct {
	ID         string `json:"id,omitempty"`
	Url        string `json:"url,omitempty"`
	Catalog    string `json:"catalog,omitempty"`
	OriginPath string `json:"originPath,omitempty"`
	IconPath   string `json:"iconPath,omitempty"`
}

type ImageNil struct {
	ID         string         `json:"id,omitempty"`
	Url        string         `json:"url,omitempty"`
	Catalog    string         `json:"catalog,omitempty"`
	OriginPath sql.NullString `json:"originPath,omitempty"`
	IconPath   sql.NullString `json:"iconPath,omitempty"`
}

func (i *ImageNil) toImage() Image {
	img := Image{
		ID:      i.ID,
		Url:     i.Url,
		Catalog: i.Catalog,
	}
	if i.OriginPath.Valid {
		img.OriginPath = i.OriginPath.String
	} else {
		img.OriginPath = ""
	}

	if i.IconPath.Valid {
		img.IconPath = i.IconPath.String
	} else {
		img.IconPath = ""
	}
	return img
}
