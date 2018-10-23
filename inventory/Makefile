APP ?= inventory
INSTALL_PATH ?= /opt/inventory
PROJECT_HOME ?= github.com/comolago/shop/inventory
GOBIN ?= /usr/local/src

all: app 

clean:
	rm -f ${GOBIN}/${APP}

app:
	mkdir -p ${GOPATH}/src/${PROJECT_HOME}
	cp -dpR * ${GOPATH}/src/${PROJECT_HOME}
	cd ${GOPATH}/src/${PROJECT_HOME} && make build

build:
	go build -tags netgo -o ${GOBIN}/${APP} main.go

install:
	mkdir -p ${INSTALL_PATH}/bin
	cp ${GOBIN}/* ${INSTALL_PATH}/bin

define-entrypoint:
	ln -s ${INSTALL_PATH}/bin/${APP} ${INSTALL_PATH}/bin/container-entrypoint

