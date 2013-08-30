package main

import (
	"github.com/emicklei/go-restful"
)

// These are the route path for which CORS is allowed
// http://en.wikipedia.org/wiki/Cross-origin_resource_sharing
var corsRoutes = []string{
	"/{alias}/{database}/{collection}/{_id}",
	"/{alias}/{database}/{collection}",
	"/{alias}/{database}",
}

type DocumentResource struct{}

func (d DocumentResource) RegisterTo(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/docs")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)
	alias := ws.PathParameter("alias", "Name of the MongoDB instance as specified in the configuration")

	if props.GetBool("http.server.cors", false) {
		ws.Filter(enableCORSFilter)
		for i := 0; i < len(corsRoutes); i++ {
			ws.Route(ws.Method("OPTIONS").Path(corsRoutes[i]).To(optionsOK))
		}
	}

	ws.Route(ws.GET("/").To(d.getAllAliases).
		Doc("Return all Mongo DB aliases from the configuration").
		Operation("getAllAliases"))

	ws.Route(ws.GET("/{alias}").To(d.getAllDatabaseNames).
		Doc("Return all database names").
		Operation("getAllDatabaseNames").
		Param(alias))

	database := ws.PathParameter("database", "Database name from the MongoDB instance")

	ws.Route(ws.GET("/{alias}/{database}").To(d.getAllCollectionNames).
		Doc("Return all collections for the database").
		Operation("getAllCollectionNames").
		Param(alias).
		Param(database))

	collection := ws.PathParameter("collection", "Collection name from the database")
	id := ws.PathParameter("_id", "Storage identifier of the document")

	ws.Route(ws.GET("/{alias}/{database}/{collection}/{_id}").To(d.getDocument).
		Doc("Return a document from a collection from the database by its internal _id").
		Operation("getDocument").
		Param(alias).
		Param(database).
		Param(collection).
		Param(id))

	ws.Route(ws.DELETE("/{alias}/{database}/{collection}/{_id}").To(d.deleteDocument).
		Doc("Deletes a document from a collection from the database by its internal _id").
		Operation("deleteDocument").
		Param(alias).
		Param(database).
		Param(collection).
		Param(id))

	ws.Route(ws.PUT("/{alias}/{database}/{collection}/{_id}").To(d.putDocument).
		Doc("Store a document to a collection from the database using its internal _id").
		Operation("putDocument").
		Param(alias).
		Param(database).
		Param(collection).
		Param(id).
		Reads(""))

	ws.Route(ws.POST("/{alias}/{database}/{collection}").To(d.postDocument).
		Doc("Store a document to a collection from the database").
		Operation("postDocument").
		Param(alias).
		Param(database).
		Param(collection).
		Reads(""))

	ws.Route(ws.GET("/{alias}/{database}/{collection}").To(d.getDocuments).
		Doc("Return documents (max 10 by default) from a collection from the database.").
		Operation("getDocuments").
		Param(alias).
		Param(database).
		Param(collection).
		Param(ws.QueryParameter("query", "query in json format")).
		Param(ws.QueryParameter("fields", "comma separated list of field names")).
		Param(ws.QueryParameter("skip", "number of documents to skip in the result set, default=0")).
		Param(ws.QueryParameter("limit", "maximum number of documents in the result set, default=10")).
		Param(ws.QueryParameter("sort", "comma separated list of field names")))

	ws.Route(ws.GET("/{alias}/{database}/{collection}/{_id}/{fields}").To(d.getSubDocument).
		Doc("Get a partial document using the internal _id and fields (comma separated field names)").
		Operation("getSubDocument").
		Param(alias).
		Param(database).
		Param(collection).
		Param(id))

	container.Add(ws)
}
