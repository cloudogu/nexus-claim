#
# useful targets:
#
# update-dependencies
#	calls glide to recreate glide.lock and update dependencies
#
# info
#	prints build time information
#
# build
#	builds the executable and ubuntu packages for trusty and xenial
#
# generate
#	generates required go file
#
# unit-test
#	performs unit-testing
#
# integration-test
#	not implemented yet
#
# static-analysis
#	performs static source code analysis
#
# deploy
#	deploys ubuntu packages for trusty and xenial to the according repositories
#
# undeploy
#	undeploys ubuntu packages for trusty and xenial from the according repositories
#
# clean
#	remove target folder
#
# dist-clean
#	also removes any downloaded dependencies
#

# collect packages and dependencies for later usage
PACKAGES=$(shell go list ./... | grep -v /vendor/)

ARTIFACT_ID=nexus-claim
VERSION=0.1.0
COMMIT_ID:=$(shell git rev-parse HEAD)


# directory settings
TARGET_DIR=target

# make target files
EXECUTABLE=${TARGET_DIR}/${ARTIFACT_ID}
PACKAGE=${TARGET_DIR}/${ARTIFACT_ID}-${VERSION}.tar.gz
XUNIT_XML=${TARGET_DIR}/unit-tests.xml

# deployment
APT_API_BASE_URL=https://apt-api.cloudogu.com/api


# tools
LINT=gometalinter
GLIDE=glide
GO2XUNIT=go2xunit


# flags
LINTFLAGS=--vendor --exclude="vendor" --exclude="_test.go"
LINTFLAGS+=--disable-all --enable=errcheck --enable=vet --enable=golint
LINTFLAGS+=--deadline=2m
LDFLAGS=-ldflags "-linkmode external -extldflags -static -X main.Version=${VERSION} -X main.CommitID=${COMMIT_ID}"
GLIDEFLAGS=



# choose the environment, if BUILD_URL environment variable is available then we are on ci (jenkins)
ifdef BUILD_URL
ENVIRONMENT=ci
GLIDEFLAGS+=--no-color
else
ENVIRONMENT=local
endif


# default goal is "build"
#
.DEFAULT_GOAL:=build

# updating dependencies
#
update-dependencies: glide.lock

glide.lock: glide.yaml
	${GLIDE} ${GLIDEFLAGS} up

# build steps: dependencies, compile, package
#
# XXX dependencies- target can not be associated to a file.
# As a consequence make build will always trigger a full build, even if targets already exist.
#
info:
	@echo "dumping build information ..."
	@echo "Version    : $(VERSION)"
	@echo "Snapshot   : $(SNAPSHOT)"
	@echo "Build-Time : $(BUILD_TIME)"
	@echo "Commit-ID  : $(COMMIT_ID)"
	@echo "Environment: $(ENVIRONMENT)"
	@echo "Branch     : $(BRANCH)"
	@echo "Branch-Type: $(BRANCH_TYPE)"
	@echo "Packages   : $(PACKAGES)"

dependencies: info
	@echo "installing dependencies ..."
	${GLIDE} ${GLIDEFLAGS} install

#generate
generate:
	@echo "generating go files"
	go generate

${EXECUTABLE}: dependencies generate
	@echo "compiling ..."
	mkdir -p $(TARGET_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo ${LDFLAGS} -o $@
	@echo "... executable can be found at $@"

${PACKAGE}: ${EXECUTABLE}
	cd ${TARGET_DIR} && tar cvzf ${ARTIFACT_ID}-${VERSION}.tar.gz ${ARTIFACT_ID}

build: ${PACKAGE}


# unit tests
#
unit-test: ${XUNIT_XML}

${XUNIT_XML}:
	mkdir -p $(TARGET_DIR)
	go test -v $(PACKAGES) | tee ${TARGET_DIR}/unit-tests.log
	@if grep '^FAIL' ${TARGET_DIR}/unit-tests.log; then \
		exit 1; \
	fi
	cat ${TARGET_DIR}/unit-tests.log | go2xunit -output $@


# integration tests, not yet
#
integration-test:
	@echo "not yet implemented"


# static analysis
#
static-analysis: static-analysis-${ENVIRONMENT}

static-analysis-ci: ${TARGET_DIR}/static-analysis-cs.log
	@if [ X"$${CI_PULL_REQUEST}" != X"" -a X"$${CI_PULL_REQUEST}" != X"null" ] ; then cat $< | CI_COMMIT=$(COMMIT_ID) reviewdog -f=checkstyle -ci="common" ; fi

static-analysis-local: ${TARGET_DIR}/static-analysis-cs.log ${TARGET_DIR}/static-analysis.log
	@echo ""
	@echo "differences to develop branch:"
	@echo ""
	@cat $< | reviewdog -f checkstyle -diff "git diff develop"

${TARGET_DIR}/static-analysis.log:
	@mkdir -p ${TARGET_DIR}
	@echo ""
	@echo "complete static analysis:"
	@echo ""
	@$(LINT) ${LINTFLAGS} ./... | tee $@

${TARGET_DIR}/static-analysis-cs.log:
	@mkdir -p ${TARGET_DIR}
	@$(LINT) ${LINTFLAGS} --checkstyle ./... > $@ | true


# clean lifecycle
#
clean:
	rm -rf ${TARGET_DIR}

dist-clean: clean
	rm -rf node_modules
	rm -rf public/vendor
	rm -rf vendor
	rm -rf npm-cache
	rm -rf bower

.PHONY: update-dependencies
.PHONY: build dependencies info
.PHONY: static-analysis static-analysis-ci static-analysis-local
.PHONY: integration-test
.PHONY: unit-test
.PHONY: clean dist-clean
