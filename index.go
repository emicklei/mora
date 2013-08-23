package main

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

func enableCORS(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	resp.AddHeader("Access-Control-Allow-Origin", "*")
	chain.ProcessFilter(req, resp)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, props["swagger.path"], http.StatusMovedPermanently)
}

func icon(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, props["swagger.path"] + "images/mora.ico", http.StatusMovedPermanently)
}
