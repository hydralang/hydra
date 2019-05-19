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
