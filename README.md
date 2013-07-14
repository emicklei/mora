# Mora - MongoDB Rest API

#### generic REST server in Go for accessing MongoDB documents and meta data
	
##### Example		
	
	http://localhost:8181/docs/localhost/landskape/connections/51caec2e95c51cb63a584fde	

returns the document from

 - hostport=localhost, mongodb hosted on localhost (no port so using default :27017)
 - database=landskape
 - collection=connections
 - _id=51caec2e95c51cb63a584fde

##### API	
			
	GET /docs/{hostport}
	
	hostport ::= <address>[:<port>]
	address  ::= <hostname>|<ip>
	
	e.g. /localhost:27017/docs , /localhost/docs

returns a JSON document with the names of all databases	
			
	GET /docs/{hostport}/{database}
	
returns a JSON document with the names of all collections in a database	
	
	GET /docs/{hostport}/{database}/{collection}/{_id}

returns a JSON document from a collection using its _id							

	GET /docs/{hostport}/{database}/{collection}
	
returns a JSON document with the first (max 10) documents in a collection								

	PUT /docs/{hostport}/{database}/{collection}/{_id}
	(todo) POST /docs/{hostport}/{database}/{collection}
	
stores a JSON document in a colllection	


Install
						
	go get -u github.com/emicklei/mora
	
Build (inside mora folder)
	
	go build 

Run

	./mora -config mora.properties
	
&copy; 2013, http://ernestmicklei.com. MIT License	

## Mora API Web UI
![Mora UI](https://s3.amazonaws.com/public.philemonworks.com/mora/mora-2013-07-14-swagger.png)