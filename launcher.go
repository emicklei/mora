package main

import (
	"flag"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/emicklei/goproperties"
	"log"
	"net/http"
)

var propertiesFile = flag.String("config", "mora.properties", "the configuration file")

func main() {
	flag.Parse()
	info("loading configuration from [%s]", *propertiesFile)
	props, err := properties.Load(*propertiesFile)
	if err != nil {
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
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: props["swagger.file.path"]}
	swagger.InstallSwaggerService(config)

	http.HandleFunc("/", index)
	http.HandleFunc("/favion.ico", icon)

	info("ready to serve on %s", basePath)
	log.Fatal(http.ListenAndServe(":"+props["http.server.port"], nil))
}

func info(template string, values ...interface{}) {
	log.Printf("[mora] "+template+"\n", values...)
}
