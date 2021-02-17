GO ?= go
GOFMT ?= gofmt "-s"
PACKAGES ?= $(shell $(GO) list ./...)
VETPACKAGES ?= $(shell $(GO) list ./... | grep -v /examples/)
GOFILES := $(shell find . -name "*.go")

TESTFOLDER := $(shell $(GO) list ./... | grep -E 'dbal/schema$$|dbal/query$$' | grep -v examples)
TESTTAGS ?= ""

XUN_MODE ?= "test"
XUN_UNIT_LOG ?= "/logs/mysql.log"
XUN_UNIT_DSN ?= "mysql"
XUN_UNIT_MYSQL_DSN ?= "root:123456@tcp(mysql:3306)/yao?charset=utf8mb4&parseTime=True&loc=Local"

.PHONY: test
test:
	echo "mode: count" > coverage.out
	echo "driver: ${XUN_UNIT_DSN}" > coverage.report.out
	echo "DSN: ${XUN_UNIT_MYSQL_DSN}" >> coverage.report.out
	for d in $(TESTFOLDER); do \
		$(GO) test -tags $(TESTTAGS) -v -covermode=count -coverprofile=profile.out $$d > tmp.out; \
		cat tmp.out; \
		if grep -q "^--- FAIL" tmp.out; then \
			rm tmp.out; \
			exit 1; \
		elif grep -q "build failed" tmp.out; then \
			rm tmp.out; \
			exit 1; \
		elif grep -q "setup failed" tmp.out; then \
			rm tmp.out; \
			exit 1; \
		fi; \
		if [ -f profile.out ]; then \
			cat tmp.out | grep -E "ok " >> coverage.report.out; \
			cat profile.out | grep -v "mode:" >> coverage.out; \
			rm profile.out; \
			rm tmp.out; \
		fi; \
	done

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

vet:
	$(GO) vet $(VETPACKAGES)

.PHONY: lint
lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u golang.org/x/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

.PHONY: misspell-check
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error $(GOFILES)

.PHONY: misspell
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -w $(GOFILES)

.PHONY: tools
tools:
	go install golang.org/x/lint/golint; \
	go install github.com/client9/misspell/cmd/misspell;