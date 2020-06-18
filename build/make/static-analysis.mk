STATIC_ANALYSIS_DIR=$(TARGET_DIR)/static-analysis

LINT=$(GOPATH)/bin/golangci-lint
# ignore tests and mocks
LINTFLAGS=--tests=false --skip-files="^.*_mock.go$$" --skip-files="^.*/mock.*.go$$"

.PHONY: static-analysis
static-analysis: $(GOPATH)/bin/reviewdog static-analysis-$(ENVIRONMENT)

.PHONY: static-analysis-ci
static-analysis-ci: $(STATIC_ANALYSIS_DIR)/static-analysis-cs.log $(STATIC_ANALYSIS_DIR)/static-analysis.log
	@if [ X"$(CI_PULL_REQUEST)" != X"" -a X"$(CI_PULL_REQUEST)" != X"null" ] ; then cat $< | CI_COMMIT=$(COMMIT_ID) reviewdog -f=checkstyle -reporter="github-pr-review"; fi

.PHONY: static-analysis-local
static-analysis-local: $(STATIC_ANALYSIS_DIR)/static-analysis-cs.log $(STATIC_ANALYSIS_DIR)/static-analysis.log
	@echo ""
	@echo "differences to develop branch:"
	@echo ""
	@cat $< | $(GOPATH)/bin/reviewdog -f checkstyle -diff "git diff develop"

$(LINT): 
	@${GO_CALL} get -u github.com/golangci/golangci-lint/cmd/golangci-lint

$(STATIC_ANALYSIS_DIR)/static-analysis.log: $(LINT)
	@mkdir -p $(STATIC_ANALYSIS_DIR)
	@echo ""
	@echo "complete static analysis:"
	@echo ""
	@$(LINT) $(LINTFLAGS) run ./... | tee $@

$(STATIC_ANALYSIS_DIR)/static-analysis-cs.log: $(LINT)
	@mkdir -p $(STATIC_ANALYSIS_DIR)
	@$(LINT) $(LINTFLAGS) run --out-format=checkstyle ./... > $@ | true

$(GOPATH)/bin/reviewdog:
	@${GO_CALL} get -u github.com/haya14busa/reviewdog/cmd/reviewdog
