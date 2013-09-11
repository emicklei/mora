package main

import (
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo/bson"
)

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
	err = session.DB(dbname).Run(bson.M{"dbstats": 1}, &result)
	if err != nil {
		handleError(err, resp)
		return
	}
	resp.WriteEntity(result)
}
