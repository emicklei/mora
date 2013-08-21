package main

import (
	"flag"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/emicklei/goproperties"
	"log"
	"net/http"
)

var (
	props properties.Properties
	propertiesFile = flag.String("config", "mora.properties", "the configuration file")
)

func main() {
	flag.Parse()
	info("loading configuration from [%s]", *propertiesFile)
	var err error
	if props, err = properties.Load(*propertiesFile); err != nil {
		log.Fatalf("[mora] Unable to read properties:%v\n", err)
	}
	initConfiguration(props)

	restful.EnableContentEncoding = true
	restful.DefaultResponseMimeType = restful.MIME_JSON
	DocumentResource{}.Register()
	defer func() {
		closeSessions()
	}()

	basePath := "http://" + props["http.server.host"] + ":" + props["http.server.port"]

	config := swagger.Config{
		WebServices:     restful.RegisteredWebServices(),
		WebServicesUrl:  basePath,
		ApiPath:         "/apidocs.json",
		SwaggerPath:     props["swagger.path"],
		SwaggerFilePath: props["swagger.file.path"],
	}
	swagger.InstallSwaggerService(config)

	log.Println("swagger.path is "+ props["swagger.path"])
	if props["swagger.path"] != "/" {
		http.HandleFunc("/", index)
	}
	http.HandleFunc("/favion.ico", icon)

	info("ready to serve on %s", basePath)
	log.Fatal(http.ListenAndServe(":"+props["http.server.port"], nil))
}

func info(template string, values ...interface{}) {
	log.Printf("[mora] "+template+"\n", values...)
}
