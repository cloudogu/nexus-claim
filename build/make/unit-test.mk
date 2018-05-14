XUNIT_XML=target/unit-tests.xml

.PHONY: unit-test
unit-test: ${XUNIT_XML}

${XUNIT_XML}: ${GOPATH}/bin/go2xunit
	mkdir -p $(TARGET_DIR)
	go test -v $(PACKAGES) | tee target/unit-tests.log
	@if grep '^FAIL' target/unit-tests.log; then \
		exit 1; \
	fi
	@if grep '^=== RUN' target/unit-tests.log; then \
	  cat target/unit-tests.log | go2xunit -output $@; \
	fi

${GOPATH}/bin/go2xunit:
	go get github.com/tebeka/go2xunit
