GODEP=$(GOPATH)/bin/dep

$(GODEP):
	@go get -u github.com/golang/dep/cmd/dep

vendor:  $(GODEP)
	@echo "Installing dependencies using go dep..."
	@dep ensure

dependencies: vendor