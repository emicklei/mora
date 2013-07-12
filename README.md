# Mora - MongoDB Rest API

#### generic REST server in Go for accessing MongoDB documents and meta data
	
##### API	
	
Return a JSON document with the names of all databases
		
	GET /databases
	
Return a JSON document with the names of all collections in a databases
		
	GET /databases/{database}/collections

Return a JSON document from a collection using its _id
	
	GET /documents/{database}/{collection}/{_id}
		
Return a JSON document with the first (max 10) documents in a collection			

	GET /documents/{database}/{collection}
					
Store a JSON document in a colllection

	PUT /documents/{database}/{collection}/{_id}	
	(todo) POST /documents/{database}/{collection}

Install
						
	go get -u github.com/emicklei/mora
	
Build (inside mora folder)
	
	go build 

Run

	./mora -config mora.properties
	
&copy; 2013, http://ernestmicklei.com. MIT License	