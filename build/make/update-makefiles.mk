MAKEFILES_VERSION=0.0.1b2

.PHONY: update-makefiles
update-makefiles:
	@echo Updating makefiles...
	curl -L --silent https://github.com/cloudogu/makefiles/archive/v$(MAKEFILES_VERSION).tar.gz > $(TMPDIR)/makefiles-v$(MAKEFILES_VERSION).tar.gz

	tar -xzf $(TMPDIR)/makefiles-v$(MAKEFILES_VERSION).tar.gz -C $(TMPDIR)
	cp -r $(TMPDIR)/makefiles-$(MAKEFILES_VERSION)/build/make $(BUILDDIR)
