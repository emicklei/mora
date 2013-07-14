package main

import (
	"flag"
	"github.com/dmotylev/goproperties"
	"github.com/emicklei/go-restful"
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

	restful.EnableContentEncoding = true
	restful.DefaultResponseMimeType = restful.MIME_JSON
	DocumentResource{}.Register()
	defer func() {
		closeSessions()
	}()

	basePath := "http://" + props["http.server.host"] + ":" + props["http.server.port"]
	info("ready to serve on %s", basePath)
	log.Fatal(http.ListenAndServe(":"+props["http.server.port"], nil))
}

func info(template string, values ...interface{}) {
	log.Printf("[mora] "+template+"\n", values...)
}
