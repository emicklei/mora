package main

import "net/http"

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
	<html>
	<body>
	<h1>
		Mora - REST api server for MongoDB
	</h1>
	<h3>
		<a href="/apidocs">Documentation</a>
	</h3>	
	<h5>
		&copy;2013, <a href="https://github.com/emicklei/mora">https://github.com/emicklei/mora</a> 
	</h5>
	</body>
	</html>
	`))
}
