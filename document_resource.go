package main

import (
	"encoding/json"
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type DocumentResource struct{}

func (d DocumentResource) Register() {
	ws := new(restful.WebService)
	ws.Path("/docs")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)
	alias := ws.PathParameter("alias", "Name of the MongoDB instance as specified in the configuration")

	ws.Route(ws.GET("/").To(d.getAllAliases).
		Doc("Return all Mongo DB aliases from the configuration").
		Operation("getAllAliases").
		Writes(""))

	ws.Route(ws.GET("/{alias}").To(d.getAllDatabaseNames).
		Doc("Return all database names").
		Operation("getAllDatabaseNames").
		Param(alias).
		Writes(""))

	database := ws.PathParameter("database", "Database name from the MongoDB instance")

	ws.Route(ws.GET("/{alias}/{database}").To(d.getAllCollectionNames).
		Doc("Return all collections for the database").
		Operation("getAllCollectionNames").
		Param(alias).
		Param(database).
		Writes(""))

	collection := ws.PathParameter("collection", "Collection name from the database")
	id := ws.PathParameter("_id", "Storage identifier of the document")

	ws.Route(ws.GET("/{alias}/{database}/{collection}/{_id}").To(d.getDocument).
		Doc("Return a document from a collection from the database by its internal _id").
		Operation("getDocument").
		Param(alias).
		Param(database).
		Param(collection).
		Param(id).
		Writes(""))

	ws.Route(ws.PUT("/{alias}/{database}/{collection}/{_id}").To(d.putDocument).
		Doc("Store a document to a collection from the database using its internal _id").
		Operation("putDocument").
		Param(alias).
		Param(database).
		Param(collection).
		Param(id).
		Reads("").
		Writes(""))

	ws.Route(ws.POST("/{alias}/{database}/{collection}").To(d.postDocument).
		Doc("Store a document to a collection from the database").
		Operation("postDocument").
		Param(alias).
		Param(database).
		Param(collection).
		Reads("").
		Writes(""))

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
		Param(ws.QueryParameter("sort", "comma separated list of field names")).
		Writes(""))

	ws.Route(ws.GET("/{alias}/{database}/{collection}/{_id}/{fields}").To(d.getSubDocument).
		Doc("Get a partial document using the internal _id and fields (comma separated field names)").
		Operation("getSubDocument").
		Param(alias).
		Param(database).
		Param(collection).
		Param(id).
		Reads("").
		Writes(""))

	restful.Add(ws)
}

func (d DocumentResource) getAllAliases(req *restful.Request, resp *restful.Response) {
	aliases := []string{}
	for k, _ := range configurationMap {
		aliases = append(aliases, k)
	}
	resp.WriteAsJson(aliases)
}

func (d DocumentResource) getAllDatabaseNames(req *restful.Request, resp *restful.Response) {
	// filter invalids
	hostport := req.PathParameter("alias")
	if hostport == "" || strings.Index(hostport, ".") != -1 {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	session, needsClose, err := d.getMongoSession(req)
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer func() { session.Close() }()
	}
	names, err := session.DatabaseNames()
	if err != nil {
		handleError(err, resp)
		return
	}
	resp.WriteEntity(names)
}

func (d DocumentResource) getAllCollectionNames(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.getMongoSession(req)
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer func() { session.Close() }()
	}
	dbname := req.PathParameter("database")
	names, err := session.DB(dbname).CollectionNames()
	if err != nil {
		handleError(err, resp)
		return
	}
	resp.WriteEntity(names)
}

func (d DocumentResource) getDocuments(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.getMongoSession(req)
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer func() { session.Close() }()
	}
	col := d.getMongoCollection(req, session)
	query, err := d.composeQuery(col, req)
	if err != nil {
		resp.WriteError(400, err) // TODO handleError(err, resp)
		return
	}
	result := []bson.M{}
	err = query.All(&result)
	if err != nil {
		handleError(err, resp)
		return
	}
	resp.WriteEntity(result)
}

//Param(ws.QueryParameter("query", "query in json format")).
//Param(ws.QueryParameter("fields", "comma separated list of field names")).
//Param(ws.QueryParameter("skip", "number of documents to skip in the result set, default=0")).
//Param(ws.QueryParameter("limit", "maximum number of documents in the result set, default=10")).
//Param(ws.QueryParameter("sort", "comma separated list of field names")).
func (d DocumentResource) composeQuery(col *mgo.Collection, req *restful.Request) (*mgo.Query, error) {
	expression := bson.M{}
	qp := req.QueryParameter("query")
	if len(qp) > 0 {
		log.Println("query=" + qp)
		if err := json.Unmarshal([]byte(qp), &expression); err != nil {
			return nil, err
		}
		log.Printf("expression=%v\n", expression)
	}
	query := col.Find(expression)

	selection := bson.M{}
	fields := req.QueryParameter("fields")
	if len(fields) > 0 {
		for _, v := range strings.Split(fields, ",") {
			selection[v] = 1
		}
	}
	query.Select(selection)

	skip := req.QueryParameter("skip")
	if len(skip) > 0 {
		v, err := strconv.Atoi(skip)
		if err != nil {
			return nil, err
		}
		query.Skip(v)
	} else {
		query.Skip(0)
	}
	limit := req.QueryParameter("limit")
	if len(limit) > 0 {
		v, err := strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}
		query.Limit(v)
	} else {
		query.Limit(10)
	}
	sort := req.QueryParameter("sort")
	if len(sort) > 0 {
		query.Sort(strings.Split(sort, ",")...)
	}

	return query, nil
}

func (d DocumentResource) getDocument(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.getMongoSession(req)
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer func() { session.Close() }()
	}
	d.fetchDocument(d.getMongoCollection(req, session), req.PathParameter("_id"), bson.M{}, resp)
}

func (d DocumentResource) getSubDocument(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.getMongoSession(req)
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer func() { session.Close() }()
	}
	fields := req.PathParameter("fields")
	selector := bson.M{}
	for _, each := range strings.Split(fields, ",") {
		selector[each] = 1
	}
	d.fetchDocument(d.getMongoCollection(req, session), req.PathParameter("_id"), selector, resp)
}

func (d DocumentResource) fetchDocument(col *mgo.Collection, id string, selector bson.M, resp *restful.Response) {
	doc := bson.M{}
	var finderr error
	if bson.IsObjectIdHex(id) {
		finderr = col.FindId(bson.ObjectIdHex(id)).Select(selector).One(&doc)
	} else {
		finderr = col.Find(bson.M{"_id": id}).Select(selector).One(&doc)
	}
	if finderr != nil {
		handleError(finderr, resp)
	}
	resp.WriteEntity(doc)
}

// TODO check for conflict
func (d DocumentResource) putDocument(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.getMongoSession(req)
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer func() { session.Close() }()
	}
	col := d.getMongoCollection(req, session)
	doc := bson.M{"_id": req.PathParameter("_id")}
	req.ReadEntity(&doc)
	err = col.Insert(doc)
	if err != nil {
		handleError(err, resp)
		return
	}
	resp.WriteHeader(http.StatusCreated)
}

func (d DocumentResource) postDocument(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.getMongoSession(req)
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer func() { session.Close() }()
	}
	col := d.getMongoCollection(req, session)
	doc := bson.M{}
	req.ReadEntity(&doc)
	err = col.Insert(doc)
	if err != nil {
		handleError(err, resp)
		return
	}
	resp.WriteHeader(http.StatusCreated)
}

func (d DocumentResource) getMongoCollection(req *restful.Request, session *mgo.Session) *mgo.Collection {
	return session.DB(req.PathParameter("database")).C(req.PathParameter("collection"))
}

func (d DocumentResource) getMongoSession(req *restful.Request) (*mgo.Session, bool, error) {
	alias := req.PathParameter("alias")
	config, err := configuration(alias)
	if err != nil {
		return nil, false, err
	}
	session, needsClose, err := openSession(config)
	if err != nil {
		return nil, false, err
	}
	return session, needsClose, nil
}

func handleError(err error, resp *restful.Response) {
	if err.Error() == "not found" {
		resp.WriteError(http.StatusNotFound, err)
		return
	}
	if err.Error() == "unauthorized" {
		resp.WriteError(http.StatusUnauthorized, err)
		return
	}
	log.Printf("[mora] error:%v", err)
	resp.WriteError(500, err)
}
