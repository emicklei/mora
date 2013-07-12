package main

import (
	"flag"
	"github.com/dmotylev/goproperties"
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo"
	"log"
	"net/http"
)

var propertiesFile = flag.String("config", "mora.properties", "the configuration file")

func main() {
	log.Println("[mora] loading configuration ...")
	flag.Parse()
	props, err := properties.Load(*propertiesFile)
	if err != nil {
		log.Fatalf("[mora] Unable to read properties:%v\n", err)
	}
	mongo := props["mongo.connection"]
	log.Printf("[mora] connecting to %s ...\n", mongo)
	session, err := mgo.Dial(mongo)
	if err != nil {
		log.Fatalf("[mora] Unable to dial mongo [%s]:%v\n", mongo, err)
	}
	defer session.Close()

	restful.EnableContentEncoding = true
	restful.DefaultResponseMimeType = restful.MIME_JSON
	DocumentResource{session}.Register()
	MetaResource{session}.Register()

	basePath := "http://" + props["http.server.host"] + ":" + props["http.server.port"]
	log.Printf("[mora] ready to serve on %s\n", basePath)
	log.Fatal(http.ListenAndServe(":"+props["http.server.port"], nil))
}
