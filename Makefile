# Set these to the desired values
ARTIFACT_ID=nexus-claim
VERSION=0.1.0

MAKEFILES_VERSION= # Set this once we have a stable release

.DEFAULT_GOAL:=compile

include build/make/variables.mk

include build/make/info.mk

include build/make/dependencies_glide.mk

include build/make/build.mk

include build/make/unit-test.mk

include build/make/static-analysis.mk

include build/make/clean.mk

include build/make/package-tar.mk

.PHONY: update-makefiles
update-makefiles:
	@echo Updating makefiles...
	@curl -L --silent https://github.com/cloudogu/makefiles/archive/v$(MAKEFILES_VERSION).tar.gz > $(TMPDIR)/makefiles-v$(MAKEFILES_VERSION).tar.gz

	@tar -xzf $(TMPDIR)/makefiles-v$(MAKEFILES_VERSION).tar.gz -C $(TMPDIR)
	@cp -r $(TMPDIR)/makefiles-$(MAKEFILES_VERSION)/build/make $(BUILDDIR)
