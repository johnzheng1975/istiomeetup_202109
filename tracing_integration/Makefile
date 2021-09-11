include Makefile.defs

.PHONY: help
help:
	@echo 'Management commands for AuthzService:'
	@echo
	@echo 'Usage:'
	@echo '    make init                Init the project'
	@echo '    make build               Compile the project'
	@echo '    make get-deps            Runs dep ensure'
	@echo '    make gen-openapi         Generate openapi code'
	@echo '    make gen-openapi-server  Generate openapi code for server'
	@echo '    make run-openapi         Run openapi server'
	@echo '    make run                 Run service'
	@echo '    make docker_build        Docker build'
	@echo '    make docker_run          Docker run'
	@echo '    make mocks               Generate mocks'
	@echo '    make test                Run tests on a compiled project.'
	@echo '    make lint                Run linters.'
	@echo '    make clean               Clean the directory tree'
	@echo '    make full_clean          Clean the directory tree (and vendor directory)'
	@echo

bin/dep: bin/dep-${DEP_VERSION}
	@ln -sf dep-${DEP_VERSION} bin/dep
bin/dep-${DEP_VERSION}:
	@mkdir -p bin
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | INSTALL_DIRECTORY=bin DEP_RELEASE_TAG=${DEP_VERSION} sh
	@mv bin/dep $@

.PHONY: init
init:
	@make vendor

.PHONY: vendor
vendor: bin/dep
	#bin/dep ensure -v -vendor-only
	@echo 'vendor ...'

.PHONY: build
build: init 
	CGO_ENABLED=0 GOOS=${GOOS} go build -a -installsuffix cgo $(GOBUILD) ./cmd/main.go

.PHONY: get-deps
get-deps: bin/dep
#	bin/dep ensure -v
	@echo 'dep ...'

.PHONY: run
run:
	VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=root MIGRATIONS_DIRECTORY=$(MIGRATIONS_DIR) AUTHZ_DISABLED=true LOG_USE_JSON=false go run cmd/main.go
