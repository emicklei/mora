package main

import (
	"github.com/emicklei/go-restful"
)

type StatisticsResource struct {
	sessMng *SessionManager
}

func (s StatisticsResource) AddTo(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/stats")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)

	alias := ws.PathParameter("alias", "Name of the MongoDB instance as specified in the configuration")
	database := ws.PathParameter("database", "Database name from the MongoDB instance")

	ws.Route(ws.GET("/{alias}/{database}").To(s.getDatabaseStatistics).
		Doc("Return statistics for the database").
		Operation("getDatabaseStatistics").
		Param(alias).
		Param(database))

	container.Add(ws)
}
