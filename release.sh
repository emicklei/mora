#go get -v -u github.com/emicklei/goproperties
#go get -v -u github.com/emicklei/go-restful
go build *.go
rm -rf target
mkdir -p target/swagger-ui
mv configuration ./target/mora
cp mora.properties ./target
cp -r /Users/ernest/Projects/swagger-ui/dist ./target/swagger-ui/dist
cd target
	zip -r mora.zip .
	cd ..
ls -l target