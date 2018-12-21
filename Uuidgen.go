package main

import uuid "github.com/satori/go.uuid"

func uuID() string {
	u1 := uuid.Must(uuid.NewV4())
	return u1.String()
}
