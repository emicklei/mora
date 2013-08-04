# ![](Letter-M-icon.png) Mora - MongoDB Rest API

#### generic REST server in Go for accessing MongoDB documents and meta data
	
##### Example		
	
	http://localhost:8181/docs/localhost/landskape/connections/51caec2e95c51cb63a584fde	

returns the document from

 - alias=localhost, mongodb hosted on localhost (aliases are defined in properties file)
 - database=landskape
 - collection=connections
 - _id=51caec2e95c51cb63a584fde

##### API	
			
	GET /docs/{alias}
	
	In the configuration file: (e.g. mora.properties)
	
	mongod.{alias}.host=localhost
	mongod.{alias}.port=27017
	# optional
	mongod.{alias}.username=
	mongod.{alias}.password=
	mongod.{alias}.database=

returns a JSON document with the names of all databases	
			
	GET /docs/{alias}/{database}
	
returns a JSON document with the names of all collections in a database	
	
	GET /docs/{alias}/{database}/{collection}/{_id}

returns a JSON document from a collection using its _id							

	GET /docs/{alias}/{database}/{collection}
	
returns a JSON document with the first (max 10) documents in a collection		

	PUT /docs/{alias}/{database}/{collection}/{_id}
	(todo) POST /docs/{alias}/{database}/{collection}
	
stores a JSON document in a colllection	

	GET /{alias}/{database}/{collection}/{_id}/{fields}

returns selected fields of a JSON document. Currently, the fields parameter must be
a comma separated list of known fields. The document returned will always contains the internal _id.


Install
						
	go get -u github.com/emicklei/mora
	
Build (inside mora folder)
	
	go build 

Run

	./mora -config mora.properties
	
&copy; 2013, http://ernestmicklei.com. MIT License
 - Icons from http://www.iconarchive.com, CC Attribution 3.0
 - Swagger from https://github.com/wordnik/swagger-core/wiki, Apache License, Version 2.0 	

## Mora API Web UI
![Mora UI](https://s3.amazonaws.com/public.philemonworks.com/mora/mora-2013-08-04.png)