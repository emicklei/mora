package api

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/mora/session"
)

func RegisterStatisticsResource(sessMng *session.SessionManager, container *restful.Container) {
	dc := StatisticsResource{sessMng}
	dc.AddWebServiceTo(container)
}

func (s StatisticsResource) AddWebServiceTo(container *restful.Container) {
	ws := d.GetWebService(cors)
	container.Add(ws)
}

func (s StatisticsResource) GetWebService() (ws *restful.WebService) {
	ws = new(restful.WebService)
	ws.Path("/stats")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)

	alias := ws.PathParameter("alias", "Name of the MongoDB instance as specified in the configuration")
	database := ws.PathParameter("database", "Database name from the MongoDB instance")
	collection := ws.PathParameter("collection", "Collection name from the database")

	ws.Route(ws.GET("/{alias}/{database}").To(s.GetDatabaseStatistics).
		Doc("Return statistics for the database").
		Operation("GetDatabaseStatistics").
		Param(alias).
		Param(database))

	ws.Route(ws.GET("/{alias}/{database}/{collection}").To(s.GetCollectionStatistics).
		Doc("Return statistics for the collection of a database").
		Operation("GetCollectionStatistics").
		Param(alias).
		Param(database).
		Param(collection))

	return
}
