package main

import (
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

type DocumentResource struct {
	session *mgo.Session
}

func (d DocumentResource) Register() {
	ws := new(restful.WebService)
	ws.Consumes("*/*")
	restful.DefaultResponseMimeType = restful.MIME_JSON
	ws.Route(ws.GET("/databases/{database}/{collection}/{id}").To(d.getDocument))
	ws.Route(ws.GET("/databases/{database}/collections").To(d.getAllCollectionNames))
	restful.Add(ws)
}

func (d DocumentResource) getDocument(req *restful.Request, resp *restful.Response) {
	db := d.session.DB(req.PathParameter("database"))
	col := db.C(req.PathParameter("collection"))
	doc := bson.M{}
	err := col.FindId(bson.ObjectIdHex(req.PathParameter("id"))).One(&doc)
	if err != nil {
		if "not found" == err.Error() {
			resp.WriteError(404, err)
		} else {
			log.Printf("[mora] error:%v", err)
			resp.WriteError(500, err)
		}
		return
	}
	resp.WriteEntity(doc)
}

func (d DocumentResource) getAllCollectionNames(req *restful.Request, resp *restful.Response) {
	dbname := req.PathParameter("database")
	names, err := d.session.DB(dbname).CollectionNames()
	if err != nil {
		resp.WriteError(500, err)
		return
	}
	resp.WriteEntity(names)
}
