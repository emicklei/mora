package documents

import (
	"encoding/json"
	"errors"
	"github.com/emicklei/go-restful"
	. "github.com/emicklei/mora/api/response"
	"github.com/emicklei/mora/session"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/compose/mejson"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Resource struct {
	SessMng *session.SessionManager
}

//
// Returns all available aliases
//
func (d *Resource) AliasListHandler(req *restful.Request, resp *restful.Response) {
	// Get aliases from session manager
	aliases := d.SessMng.GetAliases()

	// Write response back to client
	WriteResponse(aliases, resp)
}

//
// Returns all databases in alias
//
func (d *Resource) AliasDatabasesHandler(req *restful.Request, resp *restful.Response) {
	// filter invalids
	alias := getParam("alias", req)
	if alias == "" || strings.Index(alias, ".") != -1 {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	// Mongo session
	session, needclose, err := d.SessMng.Get(alias)
	if err != nil {
		WriteError(err, resp)
		return
	}
	if needclose {
		defer session.Close()
	}

	// Get all databases in mongo
	names, err := session.DatabaseNames()
	if err != nil {
		WriteError(err, resp)
		return
	}

	// Write response back to client
	WriteResponse(names, resp)
}

//
// Returns all collections in database
//
func (d *Resource) DatabaseCollectionsHandler(req *restful.Request, resp *restful.Response) {
	// Mongo session
	session, needclose, err := d.SessMng.Get(getParam("alias", req))
	if err != nil {
		WriteError(err, resp)
		return
	}
	if needclose {
		defer session.Close()
	}

	// Database request parameter
	dbname := getParam("database", req)

	// Get collections from database
	collections, err := session.DB(dbname).CollectionNames()
	if err != nil {
		WriteError(err, resp)
		return
	}

	if collections == nil {
		err = errors.New("Unknown database: " + dbname)
		WriteError(err, resp)
		return
	}

	// Write collections back to client
	WriteResponse(collections, resp)
}

//
// Updates or inserts document(/s) in collection.
// Depending on request method
// 	POST - insert
// 	PUT  - update
//
func (d *Resource) CollectionUpdateHandler(req *restful.Request, resp *restful.Response) {
	// Read a document from request
	document := bson.M{}
	decoder := json.NewDecoder(req.Request.Body)
	err := decoder.Decode(&document)
	if err != nil {
		WriteStatusError(http.StatusBadRequest, err, resp)
		return
	}

	document, err = mejson.Unmarshal(document)
	if err != nil {
		WriteStatusError(http.StatusBadRequest, err, resp)
		return
	}

	// Mongo session
	session, needclose, err := d.SessMng.Get(getParam("alias", req))
	if err != nil {
		WriteError(err, resp)
		return
	}

	// Close session if it's needed
	if needclose {
		defer session.Close()
	}

	// Mongo Collection
	col := d.GetMongoCollection(req, session)

	// Compose a selector from request
	selector, one, err := getSelector(req)
	if err != nil {
		WriteError(err, resp)
		return
	}

	// Insert if request method is POST or no selector otherwise update
	if req.Request.Method == "POST" {
		d.handleInsert(col, selector, document, req, resp)
		return
	}

	d.handleUpdate(col, one, selector, document, req, resp)
}

func (d *Resource) successUpdate(id string, created bool, req *restful.Request, resp *restful.Response) {
	// Updated document API location
	docpath := d.documentLocation(req, id)

	// Content-Location header
	resp.AddHeader("Content-Location", docpath)

	// Information about updated document
	info := struct {
		Created bool   `json:"created"`
		Url     string `json:"url"`
	}{created, docpath}

	if created {
		WriteResponseStatus(http.StatusCreated, info, resp)
	} else {
		WriteResponse(info, resp)
	}
}

func (d *Resource) handleUpdate(col *mgo.Collection, one bool, selector, document bson.M, req *restful.Request, resp *restful.Response) {
	// Update document(/s)
	var (
		err  error
		info *mgo.ChangeInfo
	)

	// Trasform id to ObjectId if needed
	if id, _ := document[ParamID].(string); id != "" && bson.IsObjectIdHex(id) {
		document[ParamID] = bson.ObjectIdHex(id)
	}

	// Update document by id
	if one {
		info, err = col.UpsertId(selector[ParamID], document)
	} else {
		// Otherwise update all matching selector
		_, err = col.UpdateAll(selector, document)
	}
	if err != nil {
		WriteError(err, resp)
		return
	}

	var docid string
	// Get id from mongo
	if info != nil && info.UpsertedId != nil {
		docid, _ = info.UpsertedId.(string)
	}
	// Otherwise from selector
	if docid == "" {
		if id, ok := selector[ParamID].(string); ok {
			docid = id
		} else if id, ok := selector[ParamID].(bson.ObjectId); ok {
			docid = id.Hex()
		}
	}

	// Write info about updated document
	if one {
		d.successUpdate(docid, (info.Updated == 0), req, resp)
		return
	}

	// Write success response
	WriteSuccess(resp)
}

func (d *Resource) handleInsert(col *mgo.Collection, selector, document bson.M, req *restful.Request, resp *restful.Response) {
	var id string
	// Set document _id if not set
	if document[ParamID] == nil {
		// If id in selector use it
		if selector[ParamID] != nil {
			// Set document id from selector
			document[ParamID] = selector[ParamID]
			// Get string ID for content-location
			if hexid, ok := document[ParamID].(bson.ObjectId); ok {
				id = hexid.Hex()
			} else {
				id, _ = document[ParamID].(string)
			}
		} else {
			// Create new ObjectId
			docid := bson.NewObjectId()
			// Set new ID for document
			document[ParamID] = docid
			// Get string ID for content-location
			id = docid.Hex()
		}
	}

	// Insert document to collection
	if err := col.Insert(document); err != nil {
		log.Printf("Error inserting document: %v", err)
		WriteError(err, resp)
		return
	}

	d.successUpdate(id, true, req, resp)
}

//
// Finds document(/s) in collection
//
func (d *Resource) CollectionFindHandler(req *restful.Request, resp *restful.Response) {
	// Mongo session
	session, needclose, err := d.SessMng.Get(getParam("alias", req))
	if err != nil {
		WriteError(err, resp)
		return
	}

	// Close session if it's needed
	if needclose {
		defer session.Close()
	}

	// Mongo Collection
	col := d.GetMongoCollection(req, session)

	// Compose a query from request
	query, one, err := d.ComposeQuery(col, req)
	if err != nil {
		WriteStatusError(400, err, resp)
		return
	}

	// If _id parameter is included in path
	// 	queries only one document.
	// Get documents from database
	if one {
		// Get one document
		document := bson.M{}
		err = query.One(&document)
		if err != nil {
			WriteError(err, resp)
			return
		}
		jsonDocument, err := mejson.Marshal(document)
		if err != nil {
			WriteError(err, resp)
			return
		}
		WriteResponse(jsonDocument, resp)
		return
	}

	// Get all documents
	documents := []bson.M{}
	err = query.All(&documents)
	if err != nil {
		WriteError(err, resp)
		return
	}

	jsonDocuments, err := mejson.Marshal(documents)
	if err != nil {
		WriteError(err, resp)
		return
	}

	res := struct {
		Success bool        `json:"success"`
		Count   interface{} `json:"count,omitempty"`
		Prev    string      `json:"prev_url,omitempty"`
		Next    string      `json:"next_url,omitempty"`
		Data    interface{} `json:"data"`
	}{Success: true, Data: jsonDocuments}

	// Get limit amount
	limitnum := 10
	if limit := req.QueryParameter("limit"); len(limit) > 0 {
		limitnum, _ = strconv.Atoi(limit)
	}

	// If got full limit set next link
	if len(documents) == limitnum {
		res.Prev, res.Next = d.prevnexturl(req)
	}

	// Count documents if count parameter is included in query
	if c, _ := strconv.ParseBool(req.QueryParameter("count")); c {
		query.Skip(0)
		query.Limit(0)
		if n, err := query.Count(); err == nil {
			res.Count = n
			resp.AddHeader("X-Object-Count", strconv.Itoa(n))
		}
	}

	// Write result back to client
	resp.WriteEntity(res)
}

//
// Removes document(/s) from collection
//
func (d *Resource) CollectionRemoveHandler(req *restful.Request, resp *restful.Response) {
	// Mongo session
	session, needclose, err := d.SessMng.Get(getParam("alias", req))
	if err != nil {
		WriteError(err, resp)
		return
	}

	// Close session if it's needed
	if needclose {
		defer session.Close()
	}

	// Mongo Collection
	col := d.GetMongoCollection(req, session)

	// Compose a selector from request
	// Get selector from `_id` path parameter and `query` query parameter
	selector, one, err := getSelector(req)
	if err != nil {
		WriteError(err, resp)
		return
	}

	// If no selector at all - drop entire collection
	if len(selector) == 0 {
		err = col.DropCollection()
		if err != nil {
			WriteError(err, resp)
			return
		}
		WriteSuccess(resp)
		return
	}

	// Remove one document if no query, otherwise remove all matching query
	if one {
		err = col.Remove(selector)
	} else {
		_, err = col.RemoveAll(selector)
	}

	if err != nil {
		WriteError(err, resp)
		return
	}

	// Write success response
	WriteSuccess(resp)
}

//
// Composes a query for finding documents
//
// Param(ws.PathParameter(ParamID, "query in json format")).
// Param(ws.QueryParameter("query", "query in json format")).
// Param(ws.QueryParameter("fields", "comma separated list of field names")).
// Param(ws.QueryParameter("skip", "number of documents to skip in the result set, default=0")).
// Param(ws.QueryParameter("limit", "maximum number of documents in the result set, default=10")).
// Param(ws.QueryParameter("sort", "comma separated list of field names")).
//
func (d *Resource) ComposeQuery(col *mgo.Collection, req *restful.Request) (query *mgo.Query, one bool, err error) {
	// Get selector from `_id` path parameter and `query` query parameter
	selector, one, err := getSelector(req)
	if err != nil {
		return
	}

	// Create a Mongo Query
	query = col.Find(selector)

	// Fields of document to select
	if fields := getFields(req); len(fields) > 0 {
		query.Select(fields)
	}

	// If selects one from `_id` parameter that's all
	if one {
		return
	}

	// Number of documents to skip in result set
	query.Skip(queryIntParam(req, "skip", 0))

	// Maximum number of documents in the result set
	query.Limit(queryIntParam(req, "limit", 10))

	// Compose sort from comma separated list in request query
	sort := req.QueryParameter("sort")
	if len(sort) > 0 {
		query.Sort(strings.Split(sort, ",")...)
	}

	return query, false, nil
}

//
// Return document location URL
//
func (d *Resource) documentLocation(req *restful.Request, id string) (location string) {
	// Get current location url
	location = strings.TrimRight(req.Request.URL.RequestURI(), "/")

	// Remove id from current location url if any
	if reqId := req.PathParameter(ParamID); reqId != "" {
		idlen := len(reqId)
		// If id is in current location remove it
		if noid := len(location) - idlen; noid > 0 {
			if id := location[noid : noid+idlen]; id == reqId {
				location = location[:noid]
			}
		}
		location = strings.TrimRight(location, "/")
	}

	// Add id of the document
	return location + "/" + id
}

func (d *Resource) prevnexturl(req *restful.Request) (prev string, next string) {
	return d.collectionurl(false, req), d.collectionurl(true, req)
}

func (d *Resource) collectionurl(next bool, req *restful.Request) string {
	// Get current location url
	uri, _ := url.Parse(req.Request.URL.RequestURI())
	q := uri.Query()

	// Skip/limit values
	limitnum := queryIntParam(req, "limit", 10)
	skipnum := queryIntParam(req, "skip", 0)

	// Number of documents to skip
	if next {
		q.Set("skip", strconv.Itoa(skipnum+limitnum))
	} else {
		// prev
		prevskip := skipnum - limitnum
		if prevskip < 0 {
			prevskip = 0
		}
		if prevskip == skipnum {
			return ""
		}
		q.Set("skip", strconv.Itoa(prevskip))
	}

	// URL query
	uri.RawQuery = q.Encode()

	return uri.String()
}

func (d *Resource) GetMongoCollection(req *restful.Request, session *mgo.Session) *mgo.Collection {
	return session.DB(getParam("database", req)).C(req.PathParameter("collection"))
}

func getFields(req *restful.Request) bson.M {
	selector := bson.M{}
	fields := req.QueryParameter("fields")
	if len(fields) > 0 {
		for _, v := range strings.Split(fields, ",") {
			selector[v] = 1
		}
	}
	return selector
}

//
// Composes a mongo selector from request
// If _id in the path is present `one` is true and query parameter is not inclued.
//
// Param(ws.PathParameter(ParamID, "query in json format")).
// Param(ws.QueryParameter("query", "query in json format")).
func getSelector(req *restful.Request) (selector bson.M, one bool, err error) {
	selector = make(bson.M)
	// If id is included in path, dont include query
	// It only select's one item
	if id := req.PathParameter(ParamID); id != "" {
		selector[ParamID] = id
	} else {
		// Unmarshal json query if any
		if query := req.QueryParameter("query"); len(query) > 0 {
			query, err = url.QueryUnescape(query)
			if err != nil {
				return
			}
			err = json.Unmarshal([]byte(query), &selector)
			if err != nil {
				return
			}
			selector, err = mejson.Unmarshal(selector)
			if err != nil {
				return
			}
		}
	}

	// Transform string HexId to ObjectIdHex
	if selid, _ := selector[ParamID].(string); selid != "" {
		// Transform to ObjectIdHex if required
		if bson.IsObjectIdHex(selid) {
			selector[ParamID] = bson.ObjectIdHex(selid)
		} else {
			selector[ParamID] = selid
		}
		one = true
	}

	return
}

// Returns a string parameter from request path or req.Attributes
func getParam(name string, req *restful.Request) (param string) {
	// Get parameter from request path
	param = req.PathParameter(name)
	if param != "" {
		return param
	}

	// Get parameter from request attributes (set by intermediates)
	attr := req.Attribute(name)
	if attr != nil {
		param, _ = attr.(string)
	}
	return
}

func queryIntParam(req *restful.Request, name string, def int) int {
	num := def
	if strnum := req.QueryParameter(name); len(strnum) > 0 {
		num, _ = strconv.Atoi(strnum)
	}
	return num
}
