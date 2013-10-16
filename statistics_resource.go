package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo/bson"
)

type StatisticsResource struct {
	sessMng *SessionManager
}

// GET http://localhost:8181/stats/local/landskape
func (s *StatisticsResource) getDatabaseStatistics(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := s.sessMng.Get(req.PathParameter("alias"))
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}
	dbname := req.PathParameter("database")
	result := bson.M{}
	err = session.DB(dbname).Run(bson.M{"dbStats": 1, "scale": 1}, &result)
	if err != nil {
		handleError(err, resp)
		return
	}
	fmt.Printf("result:%#v", result)
	resp.WriteEntity(result)
}

// GET http://localhost:8181/stats/local/landskape/systems
func (s *StatisticsResource) getCollectionStatistics(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := s.sessMng.Get(req.PathParameter("alias"))
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}
	dbname := req.PathParameter("database")
	colname := req.PathParameter("collection")
	result := bson.M{}
	err = session.DB(dbname).Run(bson.M{"collStats": colname, "scale": 1}, &result)
	if err != nil {
		handleError(err, resp)
		return
	}
	fmt.Printf("result:%#v", result)
	resp.WriteEntity(result)
}
