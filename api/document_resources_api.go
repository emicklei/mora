package api

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/mora/session"
)

func RegisterDocumentResource(sessMng *session.SessionManager, container *restful.Container, cors bool) {
	dc := DocumentResource{sessMng}
	dc.AddWebServiceTo(container, cors)
}

func (d DocumentResource) AddWebServiceTo(container *restful.Container, cors bool) {
	ws := d.GetWebService(cors)
	container.Add(ws)
}

func (d DocumentResource) GetWebService(cors bool) (ws *restful.WebService) {
	ws = new(restful.WebService)
	ws.Path("/docs")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)

	if cors {
		corsRule := restful.CrossOriginResourceSharing{ExposeHeaders: []string{"Content-Type"}, CookiesAllowed: false, Container: container}
		ws.Filter(corsRule.Filter)
	}

	alias := ws.PathParameter("alias", "Name of the MongoDB instance as specified in the configuration")
	database := ws.PathParameter("database", "Database name from the MongoDB instance")
	collection := ws.PathParameter("collection", "Collection name from the database")
	id := ws.PathParameter("_id", "Storage identifier of the document")

	ws.Route(ws.GET("/").To(d.GetAllAliases).
		Doc("Return all Mongo DB aliases from the configuration").
		Operation("GetAllAliases"))

	ws.Route(ws.GET("/{alias}").To(d.GetAllDatabaseNames).
		Doc("Return all database names").
		Operation("GetAllDatabaseNames").
		Param(alias))

	ws.Route(ws.GET("/{alias}/{database}").To(d.GetAllCollectionNames).
		Doc("Return all collections for the database").
		Operation("GetAllCollectionNames").
		Param(alias).
		Param(database))

	ws.Route(ws.DELETE("/{alias}/{database}/{collection}").To(d.DeleteDocuments).
		Doc("Deletes documents from collection if query present, otherwise removes the entire collection.").
		Operation("DeleteDocuments").
		Param(alias).
		Param(database).
		Param(collection).
		Param(ws.QueryParameter("query", "query in json format")))

	ws.Route(ws.GET("/{alias}/{database}/{collection}/{_id}").To(d.GetDocument).
		Doc("Return a document from a collection from the database by its internal _id").
		Operation("GetDocument").
		Param(alias).
		Param(database).
		Param(collection).
		Param(id).
		Param(ws.QueryParameter("fields", "comma separated list of field names")))

	ws.Route(ws.DELETE("/{alias}/{database}/{collection}/{_id}").To(d.DeleteDocument).
		Doc("Deletes a document from a collection from the database by its internal _id").
		Operation("DeleteDocument").
		Param(alias).
		Param(database).
		Param(collection).
		Param(id))

	ws.Route(ws.PUT("/{alias}/{database}/{collection}/{_id}").To(d.PutDocument).
		Doc("Store a document to a collection from the database using its internal _id").
		Operation("PutDocument").
		Param(alias).
		Param(database).
		Param(collection).
		Param(id).
		Reads(""))

	ws.Route(ws.POST("/{alias}/{database}/{collection}").To(d.PostDocument).
		Doc("Store a document to a collection from the database").
		Operation("PostDocument").
		Param(alias).
		Param(database).
		Param(collection).
		Reads(""))

	ws.Route(ws.GET("/{alias}/{database}/{collection}").To(d.GetDocuments).
		Doc("Return documents (max 10 by default) from a collection from the database.").
		Operation("GetDocuments").
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

	return
}
