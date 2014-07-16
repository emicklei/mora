package documents

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/mora/session"
)

const ParamID = "_id" // mongo id parameter

// Creates and returns documents webservice
func WebService(sessMng *session.SessionManager) *restful.WebService {
	dc := Resource{sessMng}
	return dc.WebService()
}

// Creates and adds documents resource to container
func Register(sessMng *session.SessionManager, container *restful.Container, cors bool) {
	dc := Resource{sessMng}
	dc.Register(container, cors)
}

// Adds documents resource to container
func (d Resource) Register(container *restful.Container, cors bool) {
	ws := d.WebService()

	// Cross Origin Resource Sharing filter
	if cors {
		corsRule := restful.CrossOriginResourceSharing{ExposeHeaders: []string{"Content-Type"}, CookiesAllowed: false, Container: container}
		ws.Filter(corsRule.Filter)
	}

	// Add webservice to container
	container.Add(ws)
}

func (d Resource) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/docs")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)

	alias := ws.PathParameter("alias", "Name of the MongoDB instance as specified in the configuration")
	database := ws.PathParameter("database", "Database name from the MongoDB instance")
	collection := ws.PathParameter("collection", "Collection name from the database")
	id := ws.PathParameter(ParamID, "Storage identifier of the document")

	paramID := "{" + ParamID + "}"

	ws.Route(ws.GET("/").To(d.AliasListHandler).
		Doc("Return all Mongo DB aliases from the configuration").
		Operation("AliasListHandler"))

	//
	// Returns all available aliases
	//
	// Curl example
	//
	// 	curl http://127.0.0.1:8181/docs/
	//
	ws.Route(ws.GET("/{alias}").To(d.AliasDatabasesHandler).
		Doc("Return all database names").
		Operation("AliasDatabasesHandler").
		Param(alias))

	//
	// Returns all available databases in alias
	//
	// Curl example
	//
	// 	curl http://127.0.0.1:8181/docs/local
	//
	ws.Route(ws.GET("/{alias}/{database}").To(d.DatabaseCollectionsHandler).
		Doc("Return all collections for the database").
		Operation("DatabaseCollectionsHandler").
		Param(alias).
		Param(database))

	//
	// Inserts a document into collection
	//
	// Curl example
	//
	// curl 'http://127.0.0.1:8181/docs/local/database/new-collection'  \
	//   -D - \
	//   -X POST \
	//   -H 'Content-Type: application/json' \
	//   -H 'Accept: application/json' \
	//   --data '{"title": "test", "content": "high value content"}'
	//
	ws.Route(ws.POST("/{alias}/{database}/{collection}").To(d.CollectionUpdateHandler).
		Doc("Store a document to a collection from the database").
		Operation("CollectionUpdateHandler").
		Param(alias).
		Param(database).
		Param(collection).
		Reads(""))

	//
	// Inserts a document into collection under specified ID
	//
	// Curl example
	//
	// curl 'http://127.0.0.1:8181/docs/local/database/new-collection/document-id'  \
	//   -D - \
	//   -X POST \
	//   -H 'Content-Type: application/json' \
	//   -H 'Accept: application/json' \
	//   --data '{"title": "test", "content": "high value content"}'
	//
	ws.Route(ws.POST("/{alias}/{database}/{collection}/" + paramID).To(d.CollectionUpdateHandler).
		Doc("Store a document to a collection from the database").
		Operation("CollectionUpdateHandler").
		Param(alias).
		Param(database).
		Param(collection).
		Param(id).
		Reads(""))

	//
	// Finds document in collection by it's id
	//
	// Curl example
	//
	// curl 'http://127.0.0.1:8181/docs/local/database/new-collection/document-id'  \
	//   -D - \
	//   -H 'Accept: application/json'
	//
	ws.Route(ws.GET("/{alias}/{database}/{collection}/" + paramID).To(d.CollectionFindHandler).
		Doc("Return a document from a collection from the database by its internal " + ParamID).
		Operation("GetDocument").
		Param(alias).
		Param(database).
		Param(collection).
		Param(id).
		Param(ws.QueryParameter("fields", "comma separated list of field names")))

	//
	// Finds documents in collection
	//
	// Curl example
	//
	// curl 'http://127.0.0.1:8181/docs/local/database/new-collection?query={"new": true}'  \
	//   -D - \
	//   -H 'Accept: application/json'
	//
	ws.Route(ws.GET("/{alias}/{database}/{collection}").To(d.CollectionFindHandler).
		Doc("Return documents (max 10 by default) from a collection from the database.").
		Operation("CollectionFindHandler").
		Param(alias).
		Param(database).
		Param(collection).
		Param(ws.QueryParameter("count", "counts documents in collection and returns result in X-Object-Count header").DataType("boolean")).
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
	// curl 'http://127.0.0.1:8181/docs/local/database/new-collection/document-id'  \
	//   -D - \
	//   -X PUT \
	//   -H 'Content-Type: application/json' \
	//   -H 'Accept: application/json' \
	//   --data '{"title": "New title"}'
	//
	ws.Route(ws.PUT("/{alias}/{database}/{collection}/" + paramID).To(d.CollectionUpdateHandler).
		Doc("Updates documents in collection selected by " + ParamID + " parameter").
		Operation("CollectionUpdateHandler").
		Param(alias).
		Param(database).
		Param(collection).
		Param(id).
		Reads(""))

	//
	// Updates documents in collection
	//
	// Curl example
	//
	// curl 'http://127.0.0.1:8181/docs/local/database/new-collection?query={"new": true}'  \
	//   -D - \
	//   -X PUT \
	//   -H 'Content-Type: application/json' \
	//   -H 'Accept: application/json' \
	//   --data '{"title": "New title"}'
	//
	ws.Route(ws.PUT("/{alias}/{database}/{collection}").To(d.CollectionUpdateHandler).
		Doc("Updates documents in collection selected by query parameter").
		Operation("CollectionUpdateHandler").
		Param(alias).
		Param(database).
		Param(collection).
		Param(ws.QueryParameter("query", "query in json format")).
		Reads(""))

	//
	// Removes a document from collection
	//
	// Curl example
	//
	// 	curl -X DELETE 'http://127.0.0.1:8181/docs/local/database/new-collection/document-id'
	//
	ws.Route(ws.DELETE("/{alias}/{database}/{collection}/" + paramID).To(d.CollectionRemoveHandler).
		Doc("Deletes a document from a collection from the database by its internal " + ParamID).
		Operation("CollectionRemoveHandler").
		Param(alias).
		Param(database).
		Param(collection).
		Param(id))

	//
	// Removes documents from collection
	//
	// Curl example
	//
	// 	curl -X DELETE 'http://127.0.0.1:8181/docs/local/database/new-collection?query={"new": false}'
	//
	ws.Route(ws.DELETE("/{alias}/{database}/{collection}").To(d.CollectionRemoveHandler).
		Doc("Deletes documents from collection if query present, otherwise removes the entire collection.").
		Operation("CollectionRemoveHandler").
		Param(alias).
		Param(database).
		Param(collection).
		Param(ws.QueryParameter("query", "query in json format")))

	return ws
}
