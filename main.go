package main

import (
	"flag"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/emicklei/goproperties"
	"github.com/emicklei/mora/api"
	"github.com/emicklei/mora/session"
	"log"
	"net/http"
	"path/filepath"
)

var (
	props          properties.Properties
	propertiesFile = flag.String("config", "mora.properties", "the configuration file")

	SwaggerPath string
	MoraIcon    string
)

func main() {
	flag.Parse()

	// Load configurations from a file
	info("loading configuration from [%s]", *propertiesFile)
	var err error
	if props, err = properties.Load(*propertiesFile); err != nil {
		log.Fatalf("[mora] Unable to read properties:%v\n", err)
	}

	// Swagger configuration
	SwaggerPath = props["swagger.path"]
	MoraIcon = filepath.Join(SwaggerPath, "images/mora.ico")

	// New, shared session manager
	sessMng := session.NewSessionManager(props.SelectProperties("mongod.*"))
	defer sessMng.CloseAll()

	// Enable content encoding
	restful.EnableContentEncoding = true

	// Default Response serialize method (JSON)
	restful.DefaultResponseMimeType = restful.MIME_JSON
	restful.DefaultContainer.Router(new(restful.CurlyRouter))

	// API Cross-origin requests
	apiCors := props.GetBool("http.server.cors", false)

	// Documents API
	api.RegisterDocumentResource(sessMng, restful.DefaultContainer, apiCors)

	// Statistics API
	api.RegisterStatisticsResource(sessMng, restful.DefaultContainer)

	basePath := "http://" + props["http.server.host"] + ":" + props["http.server.port"]

	// Register Swagger UI
	swagger.InstallSwaggerService(swagger.Config{
		WebServices:     restful.RegisteredWebServices(),
		WebServicesUrl:  basePath,
		ApiPath:         "/apidocs.json",
		SwaggerPath:     SwaggerPath,
		SwaggerFilePath: props["swagger.file.path"],
	})

	// If swagger is not on `/` redirect to it
	if SwaggerPath != "/" {
		http.HandleFunc("/", index)
	}

	// Serve favicon.ico
	http.HandleFunc("/favion.ico", icon)

	info("ready to serve on %s", basePath)
	log.Fatal(http.ListenAndServe(":"+props["http.server.port"], nil))
}

// If swagger is not on `/` redirect to it
func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, SwaggerPath, http.StatusMovedPermanently)
}

func icon(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, MoraIcon, http.StatusMovedPermanently)
}

// Log wrapper
func info(template string, values ...interface{}) {
	log.Printf("[mora] "+template+"\n", values...)
}
