package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/emicklei/mora/api/documents"
	"github.com/emicklei/mora/api/statistics"
	"github.com/emicklei/mora/session"
	"github.com/magiconair/properties"
)

var (
	props          *properties.Properties
	propertiesFile = flag.String("config", "mora.properties", "the configuration file")

	SwaggerPath string
	MoraIcon    string
)

func main() {
	flag.Parse()

	// Load configurations from a file
	info("loading configuration from [%s]", *propertiesFile)
	var err error
	if props, err = properties.LoadFile(*propertiesFile, properties.UTF8); err != nil {
		log.Fatalf("[mora][error] Unable to read properties:%v\n", err)
	}

	// Swagger configuration
	SwaggerPath = props.GetString("swagger.path", "")
	MoraIcon = filepath.Join(SwaggerPath, "images/mora.ico")

	// New, shared session manager
	sessMng := session.NewSessionManager(props.FilterPrefix("mongod."))
	defer sessMng.CloseAll()

	// accept and respond in JSON unless told otherwise
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)
	// gzip if accepted
	restful.DefaultContainer.EnableContentEncoding(true)
	// faster router
	restful.DefaultContainer.Router(restful.CurlyRouter{})
	// no need to access body more than once
	restful.SetCacheReadEntity(false)

	// API Cross-origin requests
	apiCors := props.GetBool("http.server.cors", false)

	// Documents API
	documents.Register(sessMng, restful.DefaultContainer, apiCors)

	// Statistics API
	if ok := props.GetBool("mora.statistics.enable", false); ok {
		statistics.Register(sessMng, restful.DefaultContainer)
	}

	addr := props.MustGet("http.server.host") + ":" + props.MustGet("http.server.port")
	basePath := "http://" + addr

	// Register Swagger UI
	swagger.InstallSwaggerService(swagger.Config{
		WebServices:     restful.RegisteredWebServices(),
		WebServicesUrl:  basePath,
		ApiPath:         "/apidocs.json",
		SwaggerPath:     SwaggerPath,
		SwaggerFilePath: props.GetString("swagger.file.path", ""),
	})

	// If swagger is not on `/` redirect to it
	if SwaggerPath != "/" {
		http.HandleFunc("/", index)
	}

	// Serve favicon.ico
	http.HandleFunc("/favion.ico", icon)

	info("ready to serve on %s", basePath)
	log.Fatal(http.ListenAndServe(addr, nil))
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
	log.Printf("[mora][info] "+template+"\n", values...)
}
