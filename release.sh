#go get -v -u github.com/emicklei/goproperties
#go get -v -u github.com/emicklei/go-restful
go build *.go
rm -rf target
mkdir -p target/swagger-ui
mv configuration ./target/mora
cp mora.properties ./target
cp -r /Users/ernest/Projects/swagger-ui/dist ./target/swagger-ui/dist

# apply patches
cp patches/index.html ./target/swagger-ui/dist
cp patches/logo_small.png ./target/swagger-ui/dist/images
cp patches/mora.ico ./target/swagger-ui/dist/images
sed "s/89bf04/89bfAA/" ./target/swagger-ui/dist/css/screen.css > ./target/swagger-ui/dist/css/screen-mora.css

cd target
	zip -r mora.zip .
	cd ..
ls -l target