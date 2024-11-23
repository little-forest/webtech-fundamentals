# Common make file for Go applications
.PHONY: clean build test run

FINAL_NAME ?= example-app
ASSETS     ?= static templates
DIST_DIR   ?= .
PORT       ?= 8080

ifeq ($(GOOS),windows)
	FINAL_NAME := $(FINAL_NAME).exe
endif

.DEFAULT_GOAL := build

test: build
	@echo "\n==== Test ===="
	$(DIST_DIR)/$(FINAL_NAME) &
	sleep 2; nc -vz localhost $(PORT)
	PID=`lsof -i:$(PORT) | sed -e 1d | awk '{print $$2}'`; [ -n "$${PID}" ] && kill $$PID

build: clean
	@echo "\n==== Build ===="
	go build -ldflags="-s -w -buildid=" -trimpath -o $(DIST_DIR)/$(FINAL_NAME) *.go
	[ "$(DIST_DIR)" != "." ] && for D in $(ASSETS); do cp -rpv $$D $(DIST_DIR)/; done || true

ifeq ($(DIST_DIR),.)
clean:
		-rm -f $(DIST_DIR)/$(FINAL_NAME)
else
clean:
		rm -rf $(DIST_DIR)
endif

run: build
	$(DIST_DIR)/$(FINAL_NAME)

