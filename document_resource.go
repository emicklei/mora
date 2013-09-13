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

func (d *DocumentResource) getAllAliases(req *restful.Request, resp *restful.Response) {
	resp.WriteAsJson(d.sessMng.GetAliases())
}

func (d *DocumentResource) getAllDatabaseNames(req *restful.Request, resp *restful.Response) {
	// filter invalids
	hostport := req.PathParameter("alias")
	if hostport == "" || strings.Index(hostport, ".") != -1 {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	session, needsClose, err := d.sessMng.Get(req.PathParameter("alias"))
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

func (d *DocumentResource) getAllCollectionNames(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.sessMng.Get(req.PathParameter("alias"))
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}
	dbname := req.PathParameter("database")
	names, err := session.DB(dbname).CollectionNames()
	if err != nil {
		handleError(err, resp)
		return
	}
	resp.WriteEntity(names)
}

func (d *DocumentResource) getDocuments(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.sessMng.Get(req.PathParameter("alias"))
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
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
func (d *DocumentResource) composeQuery(col *mgo.Collection, req *restful.Request) (*mgo.Query, error) {
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

func (d *DocumentResource) getDocument(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.sessMng.Get(req.PathParameter("alias"))
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}
	d.fetchDocument(d.getMongoCollection(req, session), req.PathParameter("_id"), bson.M{}, resp)
}

func (d *DocumentResource) deleteDocument(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.sessMng.Get(req.PathParameter("alias"))
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}
	col := d.getMongoCollection(req, session)
	id := req.PathParameter("_id")

	if bson.IsObjectIdHex(id) {
		err = col.RemoveId(bson.ObjectIdHex(id))
	} else {
		err = col.Remove(bson.M{"_id": id})
	}
	
	if err != nil {
		handleError(err, resp)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func (d *DocumentResource) getSubDocument(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.sessMng.Get(req.PathParameter("alias"))
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}
	fields := req.PathParameter("fields")
	selector := bson.M{}
	for _, each := range strings.Split(fields, ",") {
		selector[each] = 1
	}
	d.fetchDocument(d.getMongoCollection(req, session), req.PathParameter("_id"), selector, resp)
}

func (d *DocumentResource) fetchDocument(col *mgo.Collection, id string, selector bson.M, resp *restful.Response) {
	doc := bson.M{}
	var sel *mgo.Query
	if bson.IsObjectIdHex(id) {
		sel = col.FindId(bson.ObjectIdHex(id))
	} else {
		sel = col.Find(bson.M{"_id": id})
	}
	if err := sel.Select(selector).One(&doc); err != nil {
		handleError(err, resp)
		return
	}
	resp.WriteEntity(doc)
}

// TODO check for conflict
// A document must have no _id set or one that matches the path parameter
func (d *DocumentResource) putDocument(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.sessMng.Get(req.PathParameter("alias"))
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}
	col := d.getMongoCollection(req, session)
	doc := bson.M{}
	req.ReadEntity(&doc)
	// Apply internal _id
	newId := req.PathParameter("_id")
	if bson.IsObjectIdHex(newId) {
		doc["_id"] = bson.ObjectIdHex(newId)
	} else {
		doc["_id"] = newId
	}
	_, err = col.Upsert(bson.M{"_id": doc["_id"]}, doc)
	if err != nil {
		handleError(err, resp)
		return
	}
	d.handleCreated(req, resp, newId)
}

// A document cannot have an _id set. Use PUT in that case
func (d *DocumentResource) postDocument(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.sessMng.Get(req.PathParameter("alias"))
	if err != nil {
		handleError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}
	col := d.getMongoCollection(req, session)
	doc := bson.M{}
	req.ReadEntity(&doc)
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
	d.handleCreated(req, resp, newObjectId.Hex())
}

func (d *DocumentResource) handleCreated(req *restful.Request, resp *restful.Response, id string) {
	location := req.Request.URL.RequestURI() + "/" + id
	resp.AddHeader("Content-Location", location)
	resp.WriteHeader(http.StatusCreated)
}

func (d *DocumentResource) getMongoCollection(req *restful.Request, session *mgo.Session) *mgo.Collection {
	return session.DB(req.PathParameter("database")).C(req.PathParameter("collection"))
}
