package main

import (
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

type DocumentResource struct {
	session *mgo.Session
}

func (d DocumentResource) Register() {
	ws := new(restful.WebService)
	ws.Path("/documents/{database}")
	ws.Consumes("*/*")
	ws.Route(ws.GET("/{collection}/{_id}").To(d.getDocument))
	ws.Route(ws.PUT("/{collection}/{_id}").To(d.putDocument))
	ws.Route(ws.POST("/{collection}").To(d.postDocument))
	ws.Route(ws.GET("/{collection}").To(d.getDocuments))
	restful.Add(ws)
}

func (d DocumentResource) getDocuments(req *restful.Request, resp *restful.Response) {
	col := d.getMongoCollection(req)
	query := col.Find(bson.M{}) // all
	query.Limit(10)
	result := []bson.M{}
	err := query.All(&result)
	if err != nil {
		log.Printf("[mora] error:%v", err)
		resp.WriteError(500, err)
	}
	resp.WriteEntity(result)
}

func (d DocumentResource) getDocument(req *restful.Request, resp *restful.Response) {
	col := d.getMongoCollection(req)
	doc := bson.M{}
	err := col.Find(bson.M{"_id": req.PathParameter("_id")}).One(&doc)
	if err != nil {
		// retry using hex
		err2 := col.FindId(bson.ObjectIdHex(req.PathParameter("_id"))).One(&doc)
		if err2 != nil {
			if "not found" == err2.Error() {
				resp.WriteError(404, err2)
				return
			} else {
				log.Printf("[mora] error:%v", err)
				resp.WriteError(500, err)
				return
			}
		}
	}
	resp.WriteEntity(doc)
}

func (d DocumentResource) putDocument(req *restful.Request, resp *restful.Response) {
	col := d.getMongoCollection(req)
	doc := bson.M{}
	req.ReadEntity(&doc)
	err := col.Insert(doc)
	if err != nil {
		log.Printf("[mora] error:%v", err)
		resp.WriteError(500, err)
	}
	resp.WriteHeader(http.StatusCreated)
	//resp.Write([]byte("201: Created")) json version?
}

func (d DocumentResource) postDocument(req *restful.Request, resp *restful.Response) {
	//col := d.getMongoCollection(req)
	//doc := bson.M{}
}

func (d DocumentResource) getMongoCollection(req *restful.Request) *mgo.Collection {
	db := d.session.DB(req.PathParameter("database"))
	col := db.C(req.PathParameter("collection"))
	return col
}
