package api

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/mora/session"
)

func RegisterDocumentDatabaseResource(alias, database string, sessMng *session.SessionManager, container *restful.Container, cors bool) {
	dc := DocumentResource{sessMng}
	dc.AddWebServiceTo(container, cors)
}

func (d DocumentResource) AddDatabaseWebServiceTo(alias, database string, container *restful.Container, cors bool) {
	ws := d.GetDatabaseWebService(alias, database)

	if cors {
		corsRule := restful.CrossOriginResourceSharing{ExposeHeaders: []string{"Content-Type"}, CookiesAllowed: false, Container: container}
		ws.Filter(corsRule.Filter)
	}

	container.Add(ws)
}

func (d DocumentResource) GetDatabaseWebService(alias, database string) (ws *restful.WebService) {
	ws = new(restful.WebService)
	ws.Path("/docs")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)

	ws.Filter(func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		req.SetAttribute("alias", alias)
		req.SetAttribute("database", database)
		chain.ProcessFilter(req, resp)
	})

	collection := ws.PathParameter("collection", "Collection name from the database")
	id := ws.PathParameter("_id", "Storage identifier of the document")

	ws.Route(ws.DELETE("/{collection}").To(d.DeleteDocuments).
		Doc("Deletes documents from collection if query present, otherwise removes the entire collection.").
		Operation("DeleteDocuments").
		Param(collection).
		Param(ws.QueryParameter("query", "query in json format")))

	ws.Route(ws.GET("/{collection}/{_id}").To(d.GetDocument).
		Doc("Return a document from a collection from the database by its internal _id").
		Operation("GetDocument").
		Param(collection).
		Param(id).
		Param(ws.QueryParameter("fields", "comma separated list of field names")))

	ws.Route(ws.DELETE("/{collection}/{_id}").To(d.DeleteDocument).
		Doc("Deletes a document from a collection from the database by its internal _id").
		Operation("DeleteDocument").
		Param(collection).
		Param(id))

	ws.Route(ws.PUT("/{collection}/{_id}").To(d.PutDocument).
		Doc("Store a document to a collection from the database using its internal _id").
		Operation("PutDocument").
		Param(collection).
		Param(id).
		Reads(""))

	ws.Route(ws.POST("/{collection}").To(d.PostDocument).
		Doc("Store a document to a collection from the database").
		Operation("PostDocument").
		Param(collection).
		Reads(""))

	ws.Route(ws.GET("/{collection}").To(d.GetDocuments).
		Doc("Return documents (max 10 by default) from a collection from the database.").
		Operation("GetDocuments").
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
