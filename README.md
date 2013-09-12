# ![](Letter-M-icon.png) Mora - MongoDB Rest API

#### Generic REST server for accessing MongoDB documents and meta data
	
##### Example		
	
	http://localhost:8181/docs/localhost/landskape/connections/51caec2e95c51cb63a584fde	

Returns the document from

 - alias=localhost, mongodb hosted on localhost (aliases are defined in properties file)
 - database=landskape
 - collection=connections
 - _id=51caec2e95c51cb63a584fde

#### API	
			
	GET /docs

Returns a JSON document with known aliases
	
	GET /docs/{alias}
	
In the configuration file: (e.g. mora.properties)
	
	mongod.{alias}.host=localhost
	mongod.{alias}.port=27017
	# optional
	mongod.{alias}.username=
	mongod.{alias}.password=
	mongod.{alias}.database=	

Returns a JSON document with the names of all databases	
			
	GET /docs/{alias}/{database}
	
Returns a JSON document with the names of all collections in a database	
	
	GET /docs/{alias}/{database}/{collection}/{_id}

Returns a JSON document from a collection using its _id							

	GET /docs/{alias}/{database}/{collection}
	
Returns a JSON document with the first (default 10) documents in a collection.
This method also accepts query paramters

 - query, use mongo shell syntax, e.g. {"size":42}
 - limit , maximum number of documents in the result
 - skip, offset in the result set
 - fields, comma separated list of (path-dotted) field names
 - sort, comma separated list of (path-dotted) field names

Query paramters are optional. Default values are used if left out.

	PUT /docs/{alias}/{database}/{collection}/{_id}
	(todo) POST /docs/{alias}/{database}/{collection}
	
Stores a JSON document in a colllection	

	GET /{alias}/{database}/{collection}/{_id}/{fields}

Returns selected fields of a JSON document. Currently, the fields parameter must be
a comma separated list of known fields. The document returned will always contains the internal _id.


	GET /{alias}/{database}
	
Returns statistics for the database	

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

### Run

	./mora -config mora.properties
	
&copy; 2013, http://ernestmicklei.com. MIT License
 - Icons from http://www.iconarchive.com, CC Attribution 3.0
 - Swagger from https://github.com/wordnik/swagger-core/wiki, Apache License, Version 2.0 	

## Mora API Web UI
![Mora UI](https://s3.amazonaws.com/public.philemonworks.com/mora/mora-2013-08-04.png)