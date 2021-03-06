PROJECT_NAME := "pace-api"
PKG := "github.com/coma-toast/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

dep: ## Get the dependencies
	@go get -v -d ./...

build: dep ## Build the binary file
	@go build -i -v $(PKG)

build-mac: dep ## Build for Mac dev machine. Prod will be linux
	@GOOS="darwin" go build -i -v $(PKG)

deploy: build
	ssh jjd "sudo service pace-api stop"
	rsync pace-api jjd:/home/jason/www-data/pace-api/
	rsync pace-api.sh jjd:/home/jason/www-data/pace-api/
	ssh jjd "sudo service pace-api start"
