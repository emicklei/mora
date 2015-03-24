package statistics

import (
	"fmt"
	"github.com/emicklei/go-restful"
	. "github.com/emicklei/mora/api/response"
	"github.com/emicklei/mora/session"
	"gopkg.in/mgo.v2/bson"
)

// Statistics Resource
type Resource struct {
	SessMng *session.SessionManager
}

// GET http://localhost:8181/stats/local/landskape
func (s *Resource) DatabaseStatisticsHandler(req *restful.Request, resp *restful.Response) {
	// Mongo session
	session, needsClose, err := s.SessMng.Get(req.PathParameter("alias"))
	if err != nil {
		WriteError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}

	var (
		// Request parameters
		database = req.PathParameter("database")
		// Statistics result
		result = bson.M{}
	)

	// Get statistics for database
	err = session.DB(database).Run(bson.M{"dbStats": 1, "scale": 1}, &result)
	if err != nil {
		WriteError(err, resp)
		return
	}

	// Write result to console
	fmt.Printf("stats result:%#v", result)

	// Write result back to client
	WriteResponse(result, resp)
}

// GET http://localhost:8181/stats/local/landskape/systems
func (s *Resource) CollectionStatisticsHandler(req *restful.Request, resp *restful.Response) {
	// Mongo session
	session, needsClose, err := s.SessMng.Get(req.PathParameter("alias"))
	if err != nil {
		WriteError(err, resp)
		return
	}
	if needsClose {
		defer session.Close()
	}

	var (
		// Request parameters
		collection = req.PathParameter("collection")
		database   = req.PathParameter("database")
		// Statistics result
		result = bson.M{}
	)

	// Get statistics for collection
	err = session.DB(database).Run(bson.M{"collStats": collection, "scale": 1}, &result)
	if err != nil {
		WriteError(err, resp)
		return
	}

	// Write result to console
	fmt.Printf("stats result:%#v", result)

	// Write result back to client
	WriteResponse(result, resp)
}
