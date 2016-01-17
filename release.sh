#!/bin/sh

rm -rf target
mkdir target

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

go get -v -u github.com/magiconair/properties
go get -v -u github.com/emicklei/go-restful
go get -v -u labix.org/v2/mgo 

# build Mora
go build -o target/mora main.go
cp mora.properties target

# fetch or update Swagger UI
if [ ! -d ./target/swagger-ui ]; then
	git clone https://github.com/wordnik/swagger-ui.git ./target/swagger-ui
fi

# apply customizations to Swagger UI
cp $DIR/scripts/patches/index.html ./target/swagger-ui/dist
cp $DIR/scripts/patches/logo_small.png ./target/swagger-ui/dist/images
cp $DIR/scripts/patches/mora.ico ./target/swagger-ui/dist/images
sed "s/89bf04/89bfAA/" ./target/swagger-ui/dist/css/screen.css > ./target/swagger-ui/dist/css/screen-mora.css

ARCHIVE="mora-"`date +"%Y%m%d"`".zip"
cd target
	rm -f *.zip
	zip $ARCHIVE mora
	zip $ARCHIVE mora.properties
	zip -r $ARCHIVE ./swagger-ui/dist
	cd ..
