package documents

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/mora/session"
)

// Creates and returns documents webservice
func WebServiceDefaults(alias, database string, sessMng *session.SessionManager) *restful.WebService {
	dc := Resource{sessMng}
	return dc.WebServiceDefaults(alias, database)
}

// Creates and adds documents resource to container with default alias and database
func RegisterDefaults(alias, database string, sessMng *session.SessionManager, container *restful.Container, cors bool) {
	dc := Resource{sessMng}
	dc.RegisterDefaults(alias, database, container, cors)
}

// Adds documents resource to container with default alias and database
func (d Resource) RegisterDefaults(alias, database string, container *restful.Container, cors bool) {
	ws := d.WebServiceDefaults(alias, database)

	// Cross Origin Resource Sharing filter
	if cors {
		corsRule := restful.CrossOriginResourceSharing{ExposeHeaders: []string{"Content-Type"}, CookiesAllowed: false, Container: container}
		ws.Filter(corsRule.Filter)
	}

	// Add webservice to container
	container.Add(ws)
}

func (d Resource) WebServiceDefaults(alias, database string) (ws *restful.WebService) {
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
	id := ws.PathParameter(ParamID, "Storage identifier of the document")

	paramID := "{" + ParamID + "}"

	//
	// Inserts a document into collection under specified ID
	//
	// Curl example
	//
	// 	curl -H 'Accept: application/json' -X POST '{"title": "test", "content": "high value content"}' \
	//		http://127.0.0.1:8181/local/database/collection/document-id
	//
	ws.Route(ws.POST("/{collection}/" + paramID).To(d.CollectionUpdateHandler).
		Doc("Store a document to a collection from the database").
		Operation("CollectionUpdateHandler").
		Param(collection).
		Reads(""))

	//
	// Inserts a document into collection
	//
	// Curl example
	//
	// 	curl -H 'Accept: application/json' -X POST '{"title": "test", "content": "high value content"}' \
	//		http://127.0.0.1:8181/local/database/collection
	//
	ws.Route(ws.POST("/{collection}").To(d.CollectionUpdateHandler).
		Doc("Store a document to a collection from the database").
		Operation("CollectionUpdateHandler").
		Param(collection).
		Reads(""))

	//
	// Finds document in collection by it's id
	//
	// Curl example
	//
	// 	curl -H 'Accept: application/json' \
	// 		http://127.0.0.1:8181/local/database/collection/document-id
	//
	ws.Route(ws.GET("/{collection}/" + paramID).To(d.CollectionFindHandler).
		Doc("Return a document from a collection from the database by its internal " + ParamID).
		Operation("GetDocument").
		Param(collection).
		Param(id).
		Param(ws.QueryParameter("fields", "comma separated list of field names")))

	//
	// Finds documents in collection
	//
	// Curl example
	//
	// 	curl -H 'Accept: application/json' \
	// 		http://127.0.0.1:8181/local/database/collection?query={"new": true}
	//
	ws.Route(ws.GET("/{collection}").To(d.CollectionFindHandler).
		Doc("Return documents (max 10 by default) from a collection from the database.").
		Operation("CollectionFindHandler").
		Param(collection).
		Param(ws.QueryParameter("count", "counts documents in collection and returns result in X-Object-Count header"+
		"(value should be `true` to activate)")).
		Param(ws.QueryParameter("query", "query in json format")).
		Param(ws.QueryParameter("fields", "comma separated list of field names")).
		Param(ws.QueryParameter("skip", "number of documents to skip in the result set, default=0")).
		Param(ws.QueryParameter("limit", "maximum number of documents in the result set, default=10")).
		Param(ws.QueryParameter("sort", "comma separated list of field names")))

	//
	// Updates a document in collection
	//
	// Curl example
	//
	// 	curl -H 'Accept: application/json' -X PUT '{"title": "New title"}' \
	//		http://127.0.0.1:8181/local/database/collection/document-id
	//
	ws.Route(ws.PUT("/{collection}/" + paramID).To(d.CollectionUpdateHandler).
		Doc("Updates documents in collection selected by " + ParamID + " parameter").
		Operation("CollectionUpdateHandler").
		Param(collection).
		Param(id).
		Reads(""))

	//
	// Updates documents in collection
	//
	// Curl example
	//
	// 	curl -H 'Accept: application/json' -X PUT '{"new": false}' \
	//		http://127.0.0.1:8181/local/database/collection?query={"new": true}
	//
	ws.Route(ws.PUT("/{collection}").To(d.CollectionUpdateHandler).
		Doc("Updates documents in collection selected by query parameter").
		Operation("CollectionUpdateHandler").
		Param(collection).
		Param(ws.QueryParameter("query", "query in json format")).
		Reads(""))

	//
	// Removes a document from collection
	//
	// Curl example
	//
	// 	curl -X DELETE http://127.0.0.1:8181/local/database/collection/document-id
	//
	ws.Route(ws.DELETE("/{collection}/" + paramID).To(d.CollectionRemoveHandler).
		Doc("Deletes a document from a collection from the database by its internal " + ParamID).
		Operation("CollectionRemoveHandler").
		Param(collection).
		Param(id))

	//
	// Removes documents from collection
	//
	// Curl example
	//
	// 	curl -X DELETE \
	//		http://127.0.0.1:8181/local/database/collection?query={"new": false}
	//
	ws.Route(ws.DELETE("/{collection}").To(d.CollectionRemoveHandler).
		Doc("Deletes documents from collection if query present, otherwise removes the entire collection.").
		Operation("CollectionRemoveHandler").
		Param(collection).
		Param(ws.QueryParameter("query", "query in json format")))

	return
}
