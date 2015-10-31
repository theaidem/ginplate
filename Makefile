BIN = $(GOPATH)/bin
APPNAME = app
BACKEND = ./backend
FRONTEND = ./frontend
BINDATA = $(BACKEND)/bindata.go
BINDATA_FLAGS = -pkg=main -o=$(BINDATA)

setup: 
	@cd $(FRONTEND) && npm i && gulp build
	@go get github.com/jteeuwen/go-bindata/...
	@make bindata
	@go get $(BACKEND)/...
	@echo setup: OK!

build: clean
	@cd $(FRONTEND) && gulp build
	@make bindata
	@go build -o=$(BIN)/$(APPNAME) -ldflags "-w -X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.mode=`echo 'release'`" ./backend
	@echo build: OK!

clean:
	@rm -rf  $(FRONTEND)/dist
	@rm -rf $(BINDATA)
	@echo clean: OK!

serve:
	@cd $(FRONTEND) && gulp &
	@go build -o=./$(APPNAME) -ldflags "-w -X main.mode=`echo 'debug'`" ./backend && ./$(APPNAME)

bindata:
	$(BIN)/go-bindata $(BINDATA_FLAGS) $(FRONTEND)/dist/...
	@echo bindata: OK!

