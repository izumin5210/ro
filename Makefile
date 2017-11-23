GO_PKGS := $(shell go list ./... | grep -v vendor)
GO_SRCS := $(shell git ls-files | grep -E "\.go$$" | grep -v -E "\.pb(:?\.gw)?\.go$$")
GO_TEST_FLAGS  := -v -race
COVER_FILE := coverage.txt
COVER_MODE := atomic

#  Commands
#-----------------------------------------------
.PHONY: dep
dep: Gopkg.toml Gopkg.lock
	@dep ensure -v -vendor-only

.PHONY: lint
lint:
	@gofmt -e -d -s $(GO_SRCS) | awk '{ e = 1; print $0 } END { if (e) exit(1) }'
	@echo $(GO_SRCS) | xargs -n1 golint -set_exit_status
	@go vet ./...

.PHONY: test
test: lint
	@go test $(GO_TEST_FLAGS) ./...

.PHONY: cover
cover: lint
	@echo "mode: $(COVER_MODE)" > $(COVER_FILE)
	@set -e; \
	for pkg in $(GO_PKGS); do \
		tmp=/tmp/ro-coverage.out; \
		go test $(GO_TEST_FLAGS) -coverprofile=$$tmp -covermode=$(COVER_MODE) $$pkg; \
		if [ -f $$tmp ]; then \
			cat $$tmp | tail -n +2 >> $(COVER_FILE); \
			rm $$tmp; \
		fi \
	done
