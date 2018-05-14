.PHONY: clean
clean:
	rm -rf ${TARGET_DIR}
	rm -rf ${TMPDIR}

.PHONY: dist-clean
dist-clean: clean
	rm -rf node_modules
	rm -rf public/vendor
	rm -rf vendor
	rm -rf npm-cache
	rm -rf bower
