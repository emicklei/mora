package main

import (
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo"
	"log"
	"net/http"
)

func getMongoSession(req *restful.Request) (*mgo.Session, bool, error) {
	alias := req.PathParameter("alias")
	config, err := configuration(alias)
	if err != nil {
		return nil, false, err
	}
	session, needsClose, err := openSession(config)
	if err != nil {
		return nil, false, err
	}
	return session, needsClose, nil
}

func handleError(err error, resp *restful.Response) {
	if err.Error() == "not found" {
		resp.WriteError(http.StatusNotFound, err)
		return
	}
	if err.Error() == "unauthorized" {
		resp.WriteError(http.StatusUnauthorized, err)
		return
	}
	log.Printf("[mora] error:%v", err)
	resp.AddHeader("Content-Type", "text/plain") // consider making ServiceError and write JSON
	resp.WriteErrorString(500, err.Error())
}

func optionsOK(req *restful.Request, resp *restful.Response) {
	resp.WriteHeader(http.StatusOK)
}

func enableCORSFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if origin := req.Request.Header.Get("Origin"); origin != "" {
		resp.AddHeader("Access-Control-Allow-Origin", origin)
	} else {
		resp.AddHeader("Access-Control-Allow-Origin", "*")
	}

	resp.AddHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	resp.AddHeader("Access-Control-Allow-Headers", "Content-Type")
	chain.ProcessFilter(req, resp)
}
