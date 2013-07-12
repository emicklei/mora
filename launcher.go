package main

import (
	"flag"
	"github.com/dmotylev/goproperties"
	"labix.org/v2/mgo"
	"log"
	"net/http"
)

var propertiesFile = flag.String("config", "mora.properties", "the configuration file")

func main() {
	log.Print("[mora] service startup...")
	flag.Parse()
	props, err := readProperties(*propertiesFile)
	if err != nil {
		log.Fatalf("[mora] Unable to read properties:%v\n", err)
	}
	session, err := mgo.Dial(props["mongo.connection"])
	if err != nil {
		log.Fatalf("Unable to dial mongo [%s]:%v\n", props["mongo.connection"], err)
	}
	defer session.Close()

	DocumentResource{session}.Register()
	MetaResource{session}.Register()

	basePath := "http://" + props["http.server.host"] + ":" + props["http.server.port"]
	log.Printf("[mora] ready to serve on %v\n", basePath)
	log.Fatal(http.ListenAndServe(":"+props["http.server.port"], nil))
}

func readProperties(filename string) (map[string]string, error) {
	return properties.Load(filename)
}
