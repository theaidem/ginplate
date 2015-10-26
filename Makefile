BIN = $(GOPATH)/bin
APPNAME = app
BACKEND = ./backend
FRONTEND = ./frontend
BINDATA = $(BACKEND)/bindata.go
BINDATA_FLAGS = -pkg=main -o=$(BINDATA) -prefix=frontend/dist

setup_go: 
	@go get github.com/jteeuwen/go-bindata/...
	@go get $(BACKEND)/...
	@echo setup: OK!

setup_node:
	@cd $(FRONTEND) && npm i

build: bindata
	@echo build: OK!

clean:
	@rm -rf  $(FRONTEND)/dist
	@rm -rf $(BINDATA)
	@echo clean: OK!

serve:
	@go build -o=./$(APPNAME) ./backend && ./$(APPNAME)

bindata:
	$(BIN)/go-bindata $(BINDATA_FLAGS) $(FRONTEND)/dist/...

