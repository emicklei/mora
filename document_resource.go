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
	ws.Path("/docs/{hostport}")
	ws.Consumes("*/*")
	ws.Route(ws.GET("/").To(d.getAllDatabaseNames))
	ws.Route(ws.GET("/{database}").To(d.getAllCollectionNames))
	ws.Route(ws.GET("/{database}/{collection}/{_id}").To(d.getDocument))
	ws.Route(ws.PUT("/{database}/{collection}/{_id}").To(d.putDocument))
	ws.Route(ws.POST("/{database}/{collection}").To(d.postDocument))
	ws.Route(ws.GET("/{database}/{collection}").To(d.getDocuments))
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
	doc := bson.M{}
	req.ReadEntity(&doc)
	err = col.Insert(doc)
	if err != nil {
		resp.WriteError(500, err)
	}
	resp.WriteHeader(http.StatusCreated)
	//resp.Write([]byte("201: Created")) json version?
}

func (d DocumentResource) postDocument(req *restful.Request, resp *restful.Response) {
	//col := d.getMongoCollection(req)
	//doc := bson.M{}
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
