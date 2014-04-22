package statistics

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/mora/session"
)

// Creates and adds Statistics Resource to container
func Register(sessMng *session.SessionManager, container *restful.Container) {
	dc := Resource{sessMng}
	dc.Register(container)
}

// Adds Statistics Resource to container
func (r Resource) Register(container *restful.Container) {
	container.Add(r.WebService())
}

func (r Resource) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/stats")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)

	alias := ws.PathParameter("alias", "Name of the MongoDB instance as specified in the configuration")
	database := ws.PathParameter("database", "Database name from the MongoDB instance")
	collection := ws.PathParameter("collection", "Collection name from the database")

	//
	// Returns statistics for databases in alias
	//
	// Curl example
	//
	// 	curl http://127.0.0.1:8181/stats/local/database
	//
	ws.Route(ws.GET("/{alias}/{database}").To(r.DatabaseStatisticsHandler).
		Doc("Return statistics for the database").
		Operation("DatabaseStatisticsHandler").
		Param(alias).
		Param(database))

	//
	// Returns statistics for collection in database
	//
	// Curl example
	//
	// 	curl http://127.0.0.1:8181/stats/local/database
	//
	ws.Route(ws.GET("/{alias}/{database}/{collection}").To(r.CollectionStatisticsHandler).
		Doc("Return statistics for the collection of a database").
		Operation("CollectionStatisticsHandler").
		Param(alias).
		Param(database).
		Param(collection))

	return ws
}
