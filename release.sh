#!/bin/sh

rm -rf target
mkdir -p target/swagger-ui

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# build Mora
echo "building mora..."
go build -o target/mora main.go

echo "preparing archive..."
cp mora.properties target
cp -r swagger-ui/dist target/swagger-ui

echo "creating archive..."
ARCHIVE="mora-"`date +"%Y%m%d"`".zip"
cd target
	rm -f *.zip
	zip $ARCHIVE mora
	zip $ARCHIVE mora.properties
	zip -r $ARCHIVE ./swagger-ui/dist
	cd ..
echo "done"	
