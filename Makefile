## Copyright (c) 2019 Kevin L. Mitchell
##
## Licensed under the Apache License, Version 2.0 (the "License"); you
## may not use this file except in compliance with the License.  You
## may obtain a copy of the License at
##
##      http://www.apache.org/licenses/LICENSE-2.0
##
## Unless required by applicable law or agreed to in writing, software
## distributed under the License is distributed on an "AS IS" BASIS,
## WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
## implied.  See the License for the specific language governing
## permissions and limitations under the License.

GO        = go
GOFMT     = gofmt
GOIMPORTS = goimports
GOLINT    = golint

SOURCES   = $(shell find . -name \*.go -print)

_mainRE   = ^\s*package\s\s*main\s*\(\#.*\)*$$
BINSRC    = $(shell echo "$(SOURCES)" | xargs grep '$(_mainRE)' | awk -F: '{print $$1}' | sort -u )
BINS      = $(patsubst %.go,%,$(BINSRC))

CLEAN     = $(BINS) cover.out coverage.html

all: test build

build: $(BINS)

format:
ifeq ($(CI_TEST),true)
	@imports=`$(GOIMPORTS) -l $(SOURCES)`; \
	fmts=`$(GOFMT) -l -s $(SOURCES)`; \
	all=`(echo $$imports; echo $$fmts) | sort -u`; \
	if [ "$$all" != "" ]; then \
		echo "Following files need updates:"; \
		echo; \
		echo $$all; \
		exit 1;\
	fi
else
	$(GOIMPORTS) -l -w $(SOURCES)
	$(GOFMT) -l -s -w $(SOURCES)
endif

lint:
	$(GOLINT) -set_exit_status ./...

vet:
	$(GO) vet ./...

test: format lint vet
	$(GO) test -race -coverprofile=cover.out ./...

cover: test
	$(GO) tool cover -html=cover.out -o coverage.html

clean:
	rm -f $(CLEAN)

%: %.go
	$(GO) build -o $@ $<
