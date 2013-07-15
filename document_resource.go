package main

import (
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strings"
)

type DocumentResource struct{}

func (d DocumentResource) Register() {
	ws := new(restful.WebService)
	ws.Path("/docs")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)
	hostport := ws.PathParameter("hostport", "Address of the MongoDB instance, e.g. localhost:27017")

	ws.Route(ws.GET("/{hostport}").To(d.getAllDatabaseNames).
		Doc("Return all database names").
		Param(hostport).
		Writes(""))

	database := ws.PathParameter("database", "Database name from the MongoDB instance")

	ws.Route(ws.GET("/{hostport}/{database}").To(d.getAllCollectionNames).
		Doc("Return all collections for the database").
		Param(hostport).
		Param(database).
		Writes(""))

	collection := ws.PathParameter("collection", "Collection name from the database")
	id := ws.PathParameter("_id", "Storage identifier of the document")

	ws.Route(ws.GET("/{hostport}/{database}/{collection}/{_id}").To(d.getDocument).
		Doc("Return a document from a collection from the database by its internal _id").
		Param(hostport).
		Param(database).
		Param(collection).
		Param(id).
		Writes(""))

	ws.Route(ws.PUT("/{hostport}/{database}/{collection}/{_id}").To(d.putDocument).
		Doc("Store a document to a collection from the database using its internal _id").
		Param(hostport).
		Param(database).
		Param(collection).
		Param(id).
		Reads("").
		Writes(""))

	ws.Route(ws.POST("/{hostport}/{database}/{collection}").To(d.postDocument).
		Doc("Store a document to a collection from the database").
		Param(hostport).
		Param(database).
		Param(collection).
		Reads("").
		Writes(""))

	ws.Route(ws.GET("/{hostport}/{database}/{collection}").To(d.getDocuments).
		Doc("Return documents (max 10) from a collection from the database.").
		Param(hostport).
		Param(database).
		Param(collection).
		Writes(""))

	restful.Add(ws)
}

func (d DocumentResource) getAllDatabaseNames(req *restful.Request, resp *restful.Response) {
	// filter invalids
	hostport := req.PathParameter("hostport")
	if hostport == "" || strings.Index(hostport, ".") != -1 {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	session, needsClose, err := d.getMongoSession(req)
	if err != nil {
		resp.WriteError(500, err)
		return
	}
	if needsClose {
		defer func() { session.Close() }()
	}
	names, err := session.DatabaseNames()
	if err != nil {
		log.Printf("[mora] error:%v", err)
		resp.WriteError(500, err)
		return
	}
	resp.WriteEntity(names)
}

func (d DocumentResource) getAllCollectionNames(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.getMongoSession(req)
	if err != nil {
		resp.WriteError(500, err)
		return
	}
	if needsClose {
		defer func() { session.Close() }()
	}
	dbname := req.PathParameter("database")
	names, err := session.DB(dbname).CollectionNames()
	if err != nil {
		log.Printf("[mora] error:%v", err)
		resp.WriteError(500, err)
		return
	}
	resp.WriteEntity(names)
}

func (d DocumentResource) getDocuments(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.getMongoSession(req)
	if err != nil {
		resp.WriteError(500, err)
		return
	}
	if needsClose {
		defer func() { session.Close() }()
	}
	col := d.getMongoCollection(req, session)
	query := col.Find(bson.M{}) // all
	query.Limit(10)
	result := []bson.M{}
	err = query.All(&result)
	if err != nil {
		resp.WriteError(500, err)
	}
	resp.WriteEntity(result)
}

func (d DocumentResource) getDocument(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.getMongoSession(req)
	if err != nil {
		resp.WriteError(500, err)
		return
	}
	if needsClose {
		defer func() { session.Close() }()
	}
	col := d.getMongoCollection(req, session)
	doc := bson.M{}
	id := req.PathParameter("_id")
	var finderr error
	if bson.IsObjectIdHex(id) {
		finderr = col.FindId(bson.ObjectIdHex(id)).One(&doc)
	} else {
		finderr = col.Find(bson.M{"_id": id}).One(&doc)
	}
	if finderr != nil {
		if "not found" == finderr.Error() {
			resp.WriteError(404, finderr)
			return
		} else {
			log.Printf("[mora] error:%v", finderr)
			resp.WriteError(500, finderr)
			return
		}
	}
	resp.WriteEntity(doc)
}

// TODO check for conflict
func (d DocumentResource) putDocument(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.getMongoSession(req)
	if err != nil {
		resp.WriteError(500, err)
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
		resp.WriteError(500, err)
		return
	}
	resp.WriteHeader(http.StatusCreated)
}

func (d DocumentResource) postDocument(req *restful.Request, resp *restful.Response) {
	session, needsClose, err := d.getMongoSession(req)
	if err != nil {
		resp.WriteError(500, err)
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
		resp.WriteError(500, err)
		return
	}
	resp.WriteHeader(http.StatusCreated)
}

func (d DocumentResource) getMongoCollection(req *restful.Request, session *mgo.Session) *mgo.Collection {
	return session.DB(req.PathParameter("database")).C(req.PathParameter("collection"))
}

func (d DocumentResource) getMongoSession(req *restful.Request) (*mgo.Session, bool, error) {
	hostport := req.PathParameter("hostport")
	if strings.Index(hostport, ":") == -1 {
		// append default port
		hostport += ":27017"
	}
	session, needsClose, err := openSession(hostport)
	if err != nil {
		return nil, false, err
	}
	return session, needsClose, nil
}
