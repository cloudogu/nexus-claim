.PHONY: prepare-package
prepare-package:
	@echo "Default prepare-package, to write your own, simply define a prepare-package goal in the base Makefile (AFTER importing package-tar.mk)"

.PHONY: package
package: targz-package sign-package

targz-package: $(TARGET_DIR)/$(ARTIFACT_ID) prepare-package
	@cd $(TARGET_DIR) && tar cvzf $(ARTIFACT_ID)-$(VERSION).tar.gz $(ARTIFACT_ID)

sign-package: targz-package
	@echo "Signing tar.gz package"
	@cd $(TARGET_DIR) ; shasum -a 256 $(ARTIFACT_ID)-$(VERSION).tar.gz > $(ARTIFACT_ID)-$(VERSION).tar.gz.sha256sum
