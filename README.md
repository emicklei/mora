# ![](Letter-M-icon.png) Mora - MongoDB Rest API

#### Generic REST server for accessing MongoDB documents and meta data
	
##### Example		
	
	http://localhost:8181/docs/localhost/landskape/connections/51caec2e95c51cb63a584fde	

Returns the document from

	http://localhost:8181/docs/{alias}/{database}/{collection}/{_id}

 - alias: localhost (alias is a name for particular MongoDB Server defined in configuration file)
 - database: landskape
 - collection: connections
 - _id: 51caec2e95c51cb63a584fde

#### API
			
	GET /docs

returns a JSON document with known aliases
	
	GET /docs/{alias}

returns a JSON document with the names of all databases	
			
	GET /docs/{alias}/{database}
	
returns a JSON document with the names of all collections in a database	
	
	GET /docs/{alias}/{database}/{collection}/{_id}

returns a JSON document from a collection using its _id							

	GET /docs/{alias}/{database}/{collection}
	
returns a JSON document with the first (default 10) documents in a collection.
This method also accepts query paramters

 - query, use mongo shell syntax, e.g. {"size":42}
 - limit , maximum number of documents in the result
 - skip, offset in the result set
 - fields, comma separated list of (path-dotted) field names
 - sort, comma separated list of (path-dotted) field names

Query paramters are optional. Default values are used if left out.

	PUT /docs/{alias}/{database}/{collection}/{_id}
	(todo) POST /docs/{alias}/{database}/{collection}
	
stores a JSON document in a colllection	

	GET /{alias}/{database}/{collection}/{_id}/{fields}

returns selected fields of a JSON document. Currently, the fields parameter must be
a comma separated list of known fields. The document returned will always contains the internal _id.


### Install
						
	go get -u github.com/emicklei/mora
	
### Build
	
	go build 
	
### Configuration
	
	# mora server settings:
	http.server.host=localhost
	http.server.port=8181
	# enable cross site requests
	http.server.cors=true
	
	# alias is a name for particular MongoDB Server
	# you can define as many aliases as you want
	mongod.{alias}.host=localhost
	mongod.{alias}.port=27017
	# optional
	mongod.{alias}.username=
	mongod.{alias}.password=
	mongod.{alias}.database=
	
	# moro api documentation
	http.server.api.docs.enable = true
	http.server.api.docs.swagger = ./swagger-ui/dist
	http.server.api.docs.path = /apidocs.json
	http.server.api.docs.ui = /apidocs

### Run

	./mora -config mora.properties
	
&copy; 2013, http://ernestmicklei.com. MIT License
 - Icons from http://www.iconarchive.com, CC Attribution 3.0
 - Swagger from https://github.com/wordnik/swagger-core/wiki, Apache License, Version 2.0 	

## Mora API Web UI
![Mora UI](https://s3.amazonaws.com/public.philemonworks.com/mora/mora-2013-08-04.png)
