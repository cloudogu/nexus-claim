TARGETDIR=target
LINT=$(GOPATH)/bin/gometalinter
LINTFLAGS=--vendor --exclude="vendor" --exclude="_test.go"
LINTFLAGS+=--disable-all --enable=errcheck --enable=vet --enable=golint
LINTFLAGS+=--deadline=2m


.PHONY: static-analysis
static-analysis: $(GOPATH)/bin/reviewdog static-analysis-$(ENVIRONMENT)

.PHONY: static-analysis-ci
static-analysis-ci: target/static-analysis-cs.log
	@if [ X"$$(CI_PULL_REQUEST)" != X"" -a X"$$(CI_PULL_REQUEST)" != X"null" ] ; then cat $< | CI_COMMIT=$(COMMIT_ID) reviewdog -f=checkstyle -ci="common" ; fi

.PHONY: static-analysis-local
static-analysis-local: target/static-analysis-cs.log target/static-analysis.log
	@echo ""
	@echo "differences to develop branch:"
	@echo ""
	@cat $< | reviewdog -f checkstyle -diff "git diff develop"

$(LINT): 
	go get -u gopkg.in/alecthomas/gometalinter.v2

target/static-analysis.log: 
	@mkdir -p $(TARGETDIR)
	@echo ""
	@echo "complete static analysis:"
	@echo ""
	@$(LINT) $(LINTFLAGS) ./... | tee $@

target/static-analysis-cs.log:
	@mkdir -p $(TARGETDIR)
	@$(LINT) $(LINTFLAGS) --checkstyle ./... > $@ | true

$(GOPATH)/bin/reviewdog:
	@go get -u github.com/haya14busa/reviewdog/cmd/reviewdog
