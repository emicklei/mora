package main

import (
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo"
)

type MetaResource struct {
	session *mgo.Session
}

func (m MetaResource) Register() {
	ws := new(restful.WebService)
	ws.Consumes("*/*")
	restful.DefaultResponseMimeType = restful.MIME_JSON
	ws.Route(ws.GET("/databases/{database}/collections").To(m.getAllCollectionNames))
	restful.Add(ws)
}

func (m MetaResource) getAllCollectionNames(req *restful.Request, resp *restful.Response) {
	dbname := req.PathParameter("database")
	names, err := m.session.DB(dbname).CollectionNames()
	if err != nil {
		resp.WriteError(500, err)
		return
	}
	resp.WriteEntity(names)
}
