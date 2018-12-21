package main

type ErrorImage struct {
	Url      string `json:"url,omitempty"`
	Count    int    `json:"count,omitempty"`
	Priority int    `json:"priority,omitempty"`
}
