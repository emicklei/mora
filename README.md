# ![](static/Letter-M-icon.png) Mora - Mongo Rest API

#### REST server for accessing MongoDB documents and meta data
	
##### Documents

When querying on collections those parameters are available:

	query  - use mongo shell syntax, e.g. {"size":42}
	limit  - maximum number of documents in the result
	skip   - offset in the result set
	fields - comma separated list of (path-dotted) field names
	sort   - comma separated list of (path-dotted) field names

##### Examples

###### Listing aliases


	$ curl 'http://127.0.0.1:8181/docs/' \
	>   -D - \
	>   -H 'Accept: application/json'
	HTTP/1.1 200 OK
	Content-Type: application/json
	Date: Tue, 22 Apr 2014 07:06:30 GMT
	Content-Length: 61

	{
	  "success": true,
	  "data": [
	   "test",
	   "local"
	  ]
	}

###### Listing databases

	$ curl 'http://127.0.0.1:8181/docs/local/' \
	>   -D - \
	>   -H 'Accept: application/json'
	HTTP/1.1 200 OK
	Content-Type: application/json
	Date: Tue, 22 Apr 2014 07:07:11 GMT
	Content-Length: 61

	{
	  "success": true,
	  "data": [
	   "local",
	   "use1"
	  ]
	}

###### Listing collections

	$ curl 'http://127.0.0.1:8181/docs/local/local' \
	>   -D - \
	>   -H 'Accept: application/json'
	HTTP/1.1 200 OK
	Content-Type: application/json
	Date: Tue, 22 Apr 2014 07:24:10 GMT
	Content-Length: 98

	{
	  "success": true,
	  "data": [
	   "new-collection",
	   "startup_log",
	   "system.indexes"
	  ]
	}

###### Inserting document

	$ curl 'http://127.0.0.1:8181/docs/local/local/new-collection/document-id' \
	>   -D - \
	>   -X POST \
	>   -H 'Content-Type: application/json' \
	>   -H 'Accept: application/json' \
	>   --data '{"title": "Some title", "content": "document content"}'
	HTTP/1.1 201 Created
	Content-Location: /docs/local/local/new-collection/document-id
	Content-Type: application/json
	Date: Tue, 22 Apr 2014 07:23:33 GMT
	Content-Length: 116

	{
	  "success": true,
	  "data": {
	   "created": true,
	   "url": "/docs/local/local/new-collection/document-id"
	  }
	}

###### Finding document

	$ curl 'http://127.0.0.1:8181/docs/local/local/new-collection/document-id' \
	>   -D - \
	>   -H 'Accept: application/json'
	HTTP/1.1 200 OK
	Content-Type: application/json
	Date: Tue, 22 Apr 2014 07:32:33 GMT
	Content-Length: 123

	{
	  "success": true,
	  "data": {
	   "_id": "document-id",
	   "content": "document content",
	   "title": "Some title"
	  }
	}

###### Finding documents

	$ curl 'http://127.0.0.1:8181/docs/local/local/new-collection?limit=1&skip=1' \
	>    -D - \
	>    -H 'Accept: application/json'
	HTTP/1.1 200 OK
	Content-Type: application/json
	Date: Wed, 23 Apr 2014 23:18:39 GMT
	Content-Length: 387

	{
	  "success": true,
	  "prev_url": "/docs/local/local/new-collection?limit=1\u0026skip=0",
	  "next_url": "/docs/local/local/new-collection?limit=1\u0026skip=2",
	  "data": [
	   {
	    "_id": "535849cfb734f91cdc000002",
	    "content": "document content",
	    "title": "Some title"
	   }
	  ]
	}

###### Updating document


	$ curl 'http://127.0.0.1:8181/docs/local/database/new-collection/document-id' \
	>  -D - \
	>  -X PUT \
	>  -H 'Content-Type: application/json' \
	>  -H 'Accept: application/json' \
	>  --data '{"title": "New title"}'
	HTTP/1.1 200 OK
	Content-Location: /docs/local/database/new-collection/document-id
	Content-Type: application/json
	Date: Tue, 22 Apr 2014 06:37:02 GMT
	Content-Length: 133

	{
	  "success": true,
	  "data": {
	   "created": false,
	   "url": "/docs/local/database/new-collection/document-id"
	  }
	}

###### Updating documents

	$ curl 'http://127.0.0.1:8181/docs/local/local/new-collection' \
	>   -D - \
	>   -X PUT \
	>   -H 'Content-Type: application/json' \
	>   -H 'Accept: application/json' \
	>   --data '{"$set": {"title": "New title"}}'
	HTTP/1.1 200 OK
	Content-Type: application/json
	Date: Wed, 23 Apr 2014 23:33:11 GMT
	Content-Length: 22

	{
	  "success": true
	}

###### Removing document

	$ curl 'http://127.0.0.1:8181/docs/local/local/new-collection/document-id'  \
	>   -D - \
	>   -X DELETE \
	>   -H 'Accept: application/json'
	HTTP/1.1 200 OK
	Content-Type: application/json
	Date: Tue, 22 Apr 2014 07:42:47 GMT
	Content-Length: 22

	{
	  "success": true
	}

###### Removing collection

	$ curl 'http://127.0.0.1:8181/docs/local/local/new-collection'  \
	>   -D - \
	>   -X DELETE \
	>   -H 'Accept: application/json'
	HTTP/1.1 200 OK
	Content-Type: application/json
	Date: Tue, 22 Apr 2014 07:43:24 GMT
	Content-Length: 22

	{
	  "success": true
	}

##### Statistics

###### Database statistics

	$ curl http://127.0.0.1:8181/stats/local/local -D -
	HTTP/1.1 200 OK
	Content-Type: application/json
	Date: Tue, 22 Apr 2014 08:17:46 GMT
	Content-Length: 341

	{
	  "success": true,
	  "data": {
	   "avgObjSize": 595.6,
	   "collections": 3,
	   "dataFileVersion": {
	    "major": 4,
	    "minor": 5
	   },
	   "dataSize": 5956,
	   "db": "local",
	   "fileSize": 67108864,
	   "indexSize": 0,
	   "indexes": 0,
	   "nsSizeMB": 16,
	   "numExtents": 3,
	   "objects": 10,
	   "ok": 1,
	   "storageSize": 10502144
	  }
	}

###### Collection statistics

	$ curl http://127.0.0.1:8181/stats/local/local/startup_log -D -
	HTTP/1.1 200 OK
	Content-Type: application/json
	Date: Tue, 22 Apr 2014 08:18:16 GMT
	Content-Length: 389

	{
	  "success": true,
	  "data": {
	   "avgObjSize": 728,
	   "capped": true,
	   "count": 8,
	   "indexSizes": {},
	   "lastExtentSize": 10485760,
	   "max": 9223372036854775807,
	   "nindexes": 0,
	   "ns": "local.startup_log",
	   "numExtents": 1,
	   "ok": 1,
	   "paddingFactor": 1,
	   "size": 5824,
	   "storageSize": 10485760,
	   "systemFlags": 0,
	   "totalIndexSize": 0,
	   "userFlags": 0
	  }
	}

### Install from source
						
	go get -u github.com/emicklei/mora
	
### Create a release
	
	sh release.sh

### Configuration

Mora uses a simple properties file to specify host,port,aliases and other options

	# listener info is required
	http.server.host=localhost
	http.server.port=8181
	
	# enable cross site requests
	http.server.cors=true

	# for swagger support (optional)
	swagger.path=/apidocs/
	swagger.file.path=./swagger-ui/dist

	# mongo instances are listed here; specify an alias for each
	mongod.{alias}.host=localhost
	mongod.{alias}.port=27017
	# initial and operational timeout in seconds
	mongod.{alias}.timeout=5
	# optional authentication
	mongod.{alias}.username=
	mongod.{alias}.password=
	mongod.{alias}.database=
	# alternatively, a mongodb connection string uri can be used instead
	# supported options: http://godoc.org/labix.org/v2/mgo#Dial
	mongod.{alias}.uri=mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb

	# enable /stats/ endpoint
	mora.statistics.enable=true


### Run

	$ mora -config mora.properties

### Swagger

Swagger UI is displaying automatically generated API documentation and playground.

![Swagger UI](https://s3.amazonaws.com/public.philemonworks.com/mora/mora-2013-08-04.png)

	
&copy; 2013, http://ernestmicklei.com. MIT License
 - Icons from http://www.iconarchive.com, CC Attribution 3.0
 - Swagger from https://github.com/wordnik/swagger-core/wiki, Apache License, Version 2.0
