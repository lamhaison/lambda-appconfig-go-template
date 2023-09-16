build: clean
	if [ ! -d bin ]; then mkdir bin; fi
	env GOOS=linux GOARCH=amd64 go build -o bin/appconfig_demo src/appconfig_demo.go
	cd bin && zip appconfig_demo.zip appconfig_demo

test:
	cd src && go test

clean:
	go clean
	rm -fr bin
