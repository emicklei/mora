package main

import (
	"github.com/emicklei/go-restful"
)

func (d DocumentResource) AddWebServiceTo(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/docs")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)

	if props.GetBool("http.server.cors", false) {
		cors := restful.CrossOriginResourceSharing{ExposeHeaders: []string{"Content-Type"}, CookiesAllowed: false, Container: container}
		ws.Filter(cors.Filter)
	}

	alias := ws.PathParameter("alias", "Name of the MongoDB instance as specified in the configuration")
	database := ws.PathParameter("database", "Database name from the MongoDB instance")
	collection := ws.PathParameter("collection", "Collection name from the database")
	id := ws.PathParameter("_id", "Storage identifier of the document")

	ws.Route(ws.GET("/").To(d.getAllAliases).
		Doc("Return all Mongo DB aliases from the configuration").
		Operation("getAllAliases"))

	ws.Route(ws.GET("/{alias}").To(d.getAllDatabaseNames).
		Doc("Return all database names").
		Operation("getAllDatabaseNames").
		Param(alias))

	ws.Route(ws.GET("/{alias}/{database}").To(d.getAllCollectionNames).
		Doc("Return all collections for the database").
		Operation("getAllCollectionNames").
		Param(alias).
		Param(database))

	ws.Route(ws.DELETE("/{alias}/{database}/{collection}").To(d.deleteDocuments).
		Doc("Deletes documents from collection if query present, otherwise removes the entire collection.").
		Operation("deleteDocuments").
		Param(alias).
		Param(database).
		Param(collection).
		Param(ws.QueryParameter("query", "query in json format")))

	ws.Route(ws.GET("/{alias}/{database}/{collection}/{_id}").To(d.getDocument).
		Doc("Return a document from a collection from the database by its internal _id").
		Operation("getDocument").
		Param(alias).
		Param(database).
		Param(collection).
		Param(id).
		Param(ws.QueryParameter("fields", "comma separated list of field names")))

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
		Param(ws.QueryParameter("count", "counts documents in collection and returns result in X-Object-Count header"+
		"(value should be `true` to activate)")).
		Param(ws.QueryParameter("query", "query in json format")).
		Param(ws.QueryParameter("fields", "comma separated list of field names")).
		Param(ws.QueryParameter("skip", "number of documents to skip in the result set, default=0")).
		Param(ws.QueryParameter("limit", "maximum number of documents in the result set, default=10")).
		Param(ws.QueryParameter("sort", "comma separated list of field names")))

	container.Add(ws)
}
