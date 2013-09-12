#!/bin/sh

# fetch or update all dependent Go packages
export GOPATH=`pwd`/target
go get -v -u github.com/emicklei/goproperties
go get -v -u github.com/emicklei/go-restful
go get -v -u labix.org/v2/mgo 

# build Mora
go build *.go

# copy binary to target
mv common_request ./target/mora
cp mora.properties ./target

# fetch or update Swagger UI
if [ ! -d ./target/swagger-ui ]; then
	git clone https://github.com/wordnik/swagger-ui.git ./target/swagger-ui
fi
#cp -r ./swagger-ui/dist ./target/swagger-ui/dist

# apply customizations to Swagger UI
cp patches/index.html ./target/swagger-ui/dist
cp patches/logo_small.png ./target/swagger-ui/dist/images
cp patches/mora.ico ./target/swagger-ui/dist/images
sed "s/89bf04/89bfAA/" ./target/swagger-ui/dist/css/screen.css > ./target/swagger-ui/dist/css/screen-mora.css

ARCHIVE="mora-"`date +"%Y%m%d"`".zip"
cd target
	rm -f *.zip
	zip $ARCHIVE mora
	zip $ARCHIVE mora.properties
	zip -r $ARCHIVE ./swagger-ui/dist
	cd ..