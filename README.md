# Mora - MongoDB Rest API

#### generic REST server in Go for accessing MongoDB documents and meta data
	
##### Example		
	
	http://localhost:8181/docs/localhost/landskape/connections/51caec2e95c51cb63a584fde	

returns the document from
 - mongodb hosted on localhost (no port so using default :27017)
 - database=landskape
 - collection=connections
 - _id=51caec2e95c51cb63a584fde

##### API	
	
Return a JSON document with the names of all databases
		
	GET /docs/{host}
	
	host 	::= <address>[:<port>]
	address ::= <hostname>|<ip>
	
	e.g. /localhost:27017/docs
	
Return a JSON document with the names of all collections in a database
		
	GET /docs/{host}/{database}

Return a JSON document from a collection using its _id
	
	GET /docs/{host}/{database}/{collection}/{_id}
		
Return a JSON document with the first (max 10) documents in a collection			

	GET /docs/{host}/{database}/{collection}
					
Store a JSON document in a colllection

	PUT /docs/{host}/{database}/{collection}/{_id}
	(todo) POST /docs/{host}/{database}/{collection}

Install
						
	go get -u github.com/emicklei/mora
	
Build (inside mora folder)
	
	go build 

Run

	./mora -config mora.properties
	
&copy; 2013, http://ernestmicklei.com. MIT License	