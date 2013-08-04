package main

import "net/http"

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/apidocs", http.StatusMovedPermanently)
}

func icon(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/apidocs/images/mora.ico", http.StatusMovedPermanently)
}
