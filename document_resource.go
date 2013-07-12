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
	ws.Path("/documents/{database}")
	ws.Consumes("*/*")
	ws.Route(ws.GET("/{collection}/{_id}").To(d.getDocument))
	restful.Add(ws)
}

func (d DocumentResource) getDocument(req *restful.Request, resp *restful.Response) {
	db := d.session.DB(req.PathParameter("database"))
	col := db.C(req.PathParameter("collection"))
	doc := bson.M{}
	err := col.FindId(bson.ObjectIdHex(req.PathParameter("_id"))).One(&doc)
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
