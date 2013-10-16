package main

import (
	"github.com/emicklei/go-restful"
)

func (s StatisticsResource) AddWebServiceTo(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/stats")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)

	alias := ws.PathParameter("alias", "Name of the MongoDB instance as specified in the configuration")
	database := ws.PathParameter("database", "Database name from the MongoDB instance")
	collection := ws.PathParameter("collection", "Collection name from the database")

	ws.Route(ws.GET("/{alias}/{database}").To(s.getDatabaseStatistics).
		Doc("Return statistics for the database").
		Operation("getDatabaseStatistics").
		Param(alias).
		Param(database))

	ws.Route(ws.GET("/{alias}/{database}/{collection}").To(s.getCollectionStatistics).
		Doc("Return statistics for the collection of a database").
		Operation("getCollectionStatistics").
		Param(alias).
		Param(database).
		Param(collection))

	container.Add(ws)
}
