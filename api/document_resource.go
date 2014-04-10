package api

import (
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/mora/session"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strconv"
	"strings"
)

type DocumentResource struct {
	SessMng *session.SessionManager
}

func (d *DocumentResource) GetAllAliases(req *restful.Request, resp *restful.Response) {
	resp.WriteAsJson(d.SessMng.GetAliases())
}

func (d *DocumentResource) GetAllDatabaseNames(req *restful.Request, resp *restful.Response) {
	// filter invalids
	hostport := getParam("alias", req)
	if hostport == "" || strings.Index(hostport, ".") != -1 {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	session, needsClose, err := d.SessMng.Get(getParam("alias", req))
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}
	names, err := session.DatabaseNames()
	if err != nil {
		handleError(err, resp)
		return
	}
	resp.WriteEntity(names)
}

func (d *DocumentResource) GetAllCollectionNames(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.SessMng.Get(getParam("alias", req))
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}
	dbname := getParam("database", req)
	names, err := session.DB(dbname).CollectionNames()
	if err != nil {
		handleError(err, resp)
		return
	}
	resp.WriteEntity(names)
}

func (d *DocumentResource) GetDocuments(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.SessMng.Get(getParam("alias", req))
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}
	col := d.GetMongoCollection(req, session)
	query, err := d.ComposeQuery(col, req)
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
	if c, _ := strconv.ParseBool(req.QueryParameter("count")); c {
		query.Limit(0)
		if n, err := query.Count(); err == nil {
			resp.AddHeader("X-Object-Count", strconv.Itoa(n))
		}
	}
	resp.WriteEntity(result)
}

func (d *DocumentResource) DeleteDocuments(req *restful.Request, resp *restful.Response) {
	exp, err := getQuery(req)
	if err != nil {
		handleError(err, resp)
		return
	}
	// get session
	session, needsClose, err := d.SessMng.Get(getParam("alias", req))
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}
	col := d.GetMongoCollection(req, session)
	if len(exp) == 0 {
		// Remove Entire collection
		if err = col.DropCollection(); err != nil {
			handleError(err, resp)
			return
		}
		resp.WriteHeader(http.StatusOK)
		return
	}
	// remove documents
	err = col.Remove(exp)
	if err != nil {
		handleError(err, resp)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

//Param(ws.QueryParameter("query", "query in json format")).
//Param(ws.QueryParameter("fields", "comma separated list of field names")).
//Param(ws.QueryParameter("skip", "number of documents to skip in the result set, default=0")).
//Param(ws.QueryParameter("limit", "maximum number of documents in the result set, default=10")).
//Param(ws.QueryParameter("sort", "comma separated list of field names")).
func (d *DocumentResource) ComposeQuery(col *mgo.Collection, req *restful.Request) (*mgo.Query, error) {
	expression, err := getQuery(req)
	if err != nil {
		return nil, err
	}
	query := col.Find(expression)

	if fields := getFields(req); len(fields) > 0 {
		query.Select(fields)
	}

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

func (d *DocumentResource) GetDocument(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.SessMng.Get(getParam("alias", req))
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}
	col := d.GetMongoCollection(req, session)
	doc := bson.M{}

	id := req.PathParameter("_id")
	if bson.IsObjectIdHex(id) {
		doc["_id"] = bson.ObjectIdHex(id)
	} else {
		doc["_id"] = id
	}
	query := col.Find(doc)

	if fields := getFields(req); len(fields) > 0 {
		query.Select(fields)
	}

	if err := query.One(&doc); err != nil {
		handleError(err, resp)
		return
	}
	resp.WriteEntity(doc)
}

func (d *DocumentResource) DeleteDocument(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.SessMng.Get(getParam("alias", req))
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}
	col := d.GetMongoCollection(req, session)
	id := req.PathParameter("_id")
	exp := bson.M{}
	if bson.IsObjectIdHex(id) {
		exp["_id"] = bson.ObjectIdHex(id)
	} else {
		exp["_id"] = id
	}

	if err := col.Remove(exp); err != nil {
		handleError(err, resp)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

// TODO check for conflict
// A document must have no _id set or one that matches the path parameter
func (d *DocumentResource) PutDocument(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.SessMng.Get(getParam("alias", req))
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}
	col := d.GetMongoCollection(req, session)
	doc := bson.M{}
	if err := req.ReadEntity(&doc); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, "Cannot read entity from request")
		return
	}
	// Transform document id string to a ObjectIdHex
	if docId, ok := doc["_id"].(string); ok && bson.IsObjectIdHex(docId) {
		doc["_id"] = bson.ObjectIdHex(docId)
	}
	// Create selector with id
	var id interface{}
	strId := req.PathParameter("_id")
	if bson.IsObjectIdHex(strId) {
		id = bson.ObjectIdHex(strId)
	} else {
		id = strId
	}
	sel := bson.M{"_id": id} // query selector
	_, err = col.Upsert(sel, doc)
	if err != nil {
		handleError(err, resp)
		return
	}
	d.HandleCreated(req, resp, strId)
}

// A document cannot have an _id set. Use PUT in that case
func (d *DocumentResource) PostDocument(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.SessMng.Get(getParam("alias", req))
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}
	col := d.GetMongoCollection(req, session)
	doc := bson.M{}
	if err := req.ReadEntity(&doc); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, "Cannot read entity from request")
		return
	}
	if doc["_id"] != nil {
		resp.WriteErrorString(http.StatusBadRequest, "Document cannot have _id ; use PUT instead to create one")
		return
	}
	newObjectId := bson.NewObjectId()
	doc["_id"] = newObjectId
	if err = col.Insert(doc); err != nil {
		handleError(err, resp)
		return
	}
	d.HandleCreated(req, resp, newObjectId.Hex())
}

func (d *DocumentResource) HandleCreated(req *restful.Request, resp *restful.Response, id string) {
	location := strings.TrimRight(req.Request.URL.RequestURI(), "/")
	if noid := strings.TrimRight(location, id); noid == location {
		location = noid + "/" + id
	}
	resp.AddHeader("Content-Location", location)
	resp.WriteHeader(http.StatusCreated)
}

func (d *DocumentResource) GetMongoCollection(req *restful.Request, session *mgo.Session) *mgo.Collection {
	return session.DB(getParam("database", req)).C(req.PathParameter("collection"))
}

func getFields(req *restful.Request) bson.M {
	exp := bson.M{}
	fields := req.QueryParameter("fields")
	if len(fields) > 0 {
		for _, v := range strings.Split(fields, ",") {
			exp[v] = 1
		}
	}
	return exp
}

func getQuery(req *restful.Request) (exp bson.M, err error) {
	exp = bson.M{}
	qp := req.QueryParameter("query")
	if len(qp) == 0 {
		return
	}
	err = json.Unmarshal([]byte(qp), &exp)
	return
}

func getParam(name string, req *restful.Request) string {
	if param := req.PathParameter(name); param != "" {
		return param
	}
	if attr := req.Attribute(name); attr != nil {
		param, _ := attr.(string)
		return param
	}
	return ""
}
