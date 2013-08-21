package main

import (
	"net/http"
	"log"
)

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("Redirecting to "+ props["swagger.path"])
	http.Redirect(w, r, props["swagger.path"], http.StatusMovedPermanently)
}

func icon(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, props["swagger.path"] + "images/mora.ico", http.StatusMovedPermanently)
}
