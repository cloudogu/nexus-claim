DEBIAN_TARGET_DIR=$(TARGET_DIR)/deb/content
DEBIAN_PACKAGE=$(TARGET_DIR)/deb/$(ARTIFACT_ID)_$(VERSION).deb

package: $(DEBIAN_PACKAGE) sign-package

.PHONY: prepare-package
prepare-package:
	@echo "Default prepare-package, to write your own, simply define a prepare-package goal in the base Makefile (AFTER importing package-debian.mk)"

$(TARGET_DIR)/debian-binary:
	@echo "2.0" > $@

$(DEBIAN_PACKAGE): compile $(TARGET_DIR)/debian-binary prepare-package
	@echo "Creating .deb package..."

	@install -m 0755 -d $(DEBIAN_TARGET_DIR)/control
	@sed -e "s/^Version:.*/Version: $(VERSION)/g" deb/DEBIAN/control > $(DEBIAN_TARGET_DIR)/_control
	@install -m 0644 $(DEBIAN_TARGET_DIR)/_control $(DEBIAN_TARGET_DIR)/control/control

# creating control.tar.gz
	@tar cvzf $(DEBIAN_TARGET_DIR)/control.tar.gz -C $(DEBIAN_TARGET_DIR)/control  --owner=0 --group=0 .
# populating data directory
	@install -m 0755 -d $(DEBIAN_TARGET_DIR)/data/usr/sbin

	@for dir in $$(find deb -mindepth 1 -not -name "DEBIAN" -a -type d |sed s@"^deb/"@"$(DEBIAN_TARGET_DIR)/data/"@) ; do install -m 0755 -d $${dir} ; done
	@for file in $$(find deb -mindepth 1 -type f | grep -v "DEBIAN") ; do install -m 0644 $${file} $(DEBIAN_TARGET_DIR)/data/$${file#deb/}; done
	@install -m 0755 $(TARGET_DIR)/$(ARTIFACT_ID) $(DEBIAN_TARGET_DIR)/data/usr/sbin/

# creating data.tar.gz
	@tar cvzf $(DEBIAN_TARGET_DIR)/data.tar.gz -C $(DEBIAN_TARGET_DIR)/data --owner=0 --group=0 .
# creating package
	@ar rc $@ $(TARGET_DIR)/debian-binary $(DEBIAN_TARGET_DIR)/control.tar.gz $(DEBIAN_TARGET_DIR)/data.tar.gz
	@echo "... deb package can be found at $@"

sign-package: $(DEBIAN_PACKAGE)
	@echo "Signing .deb package"
	@cd $(TARGET_DIR)/deb/; shasum -a 256 $(ARTIFACT_ID)_$(VERSION).deb > $(ARTIFACT_ID)_$(VERSION).deb.sha256sum

# deployment
#
deploy:
	@case X"${VERSION}" in *-SNAPSHOT) echo "i will not upload a snaphot version for you" ; exit 1; esac;
	@if [ X"${APT_API_USERNAME}" = X"" ] ; then echo "supply an APT_API_USERNAME environment variable"; exit 1; fi;
	@if [ X"${APT_API_PASSWORD}" = X"" ] ; then echo "supply an APT_API_PASSWORD environment variable"; exit 1; fi;
	@if [ X"${APT_API_SIGNPHRASE}" = X"" ] ; then echo "supply an APT_API_SIGNPHRASE environment variable"; exit 1; fi;
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" -F file=@"${DEBIAN_PACKAGE}" "${APT_API_BASE_URL}/files/trusty" |jq
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" -F file=@"${DEBIAN__PACKAGE}" "${APT_API_BASE_URL}/files/xenial" |jq
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" -X POST "${APT_API_BASE_URL}/repos/trusty/file/trusty/${ARTIFACT_ID}_${VERSION}.deb" |jq
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" -X POST "${APT_API_BASE_URL}/repos/xenial/file/xenial/${ARTIFACT_ID}_${VERSION}.deb" |jq
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" -X PUT -H "Content-Type: application/json" --data '{"Signing": { "Batch": true, "Passphrase": "${APT_API_SIGNPHRASE}"}}' ${APT_API_BASE_URL}/publish/trusty/trusty
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" -X PUT -H "Content-Type: application/json" --data '{"Signing": { "Batch": true, "Passphrase": "${APT_API_SIGNPHRASE}"}}' ${APT_API_BASE_URL}/publish/xenial/xenial

undeploy:
	@case X"${VERSION}" in *-SNAPSHOT) echo "i will not upload a snaphot version for you" ; exit 1; esac;
	@if [ X"${APT_API_USERNAME}" = X"" ] ; then echo "supply an APT_API_USERNAME environment variable"; exit 1; fi;
	@if [ X"${APT_API_PASSWORD}" = X"" ] ; then echo "supply an APT_API_PASSWORD environment variable"; exit 1; fi;
	@if [ X"${APT_API_SIGNPHRASE}" = X"" ] ; then echo "supply an APT_API_SIGNPHRASE environment variable"; exit 1; fi;
	PREF=$$(curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" "${APT_API_BASE_URL}/repos/xenial/packages?q=${ARTIFACT_ID}%20(${VERSION})"); \
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" -X DELETE -H 'Content-Type: application/json' --data "{\"PackageRefs\": $${PREF}}" ${APT_API_BASE_URL}/repos/xenial/packages
	PREF=$$(curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" "${APT_API_BASE_URL}/repos/trusty/packages?q=${ARTIFACT_ID}%20(${VERSION})"); \
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" -X DELETE -H 'Content-Type: application/json' --data "{\"PackageRefs\": $${PREF}}" ${APT_API_BASE_URL}/repos/trusty/packages
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" -X PUT -H "Content-Type: application/json" --data '{"Signing": { "Batch": true, "Passphrase": "${APT_API_SIGNPHRASE}"}}' ${APT_API_BASE_URL}/publish/trusty/trusty
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" -X PUT -H "Content-Type: application/json" --data '{"Signing": { "Batch": true, "Passphrase": "${APT_API_SIGNPHRASE}"}}' ${APT_API_BASE_URL}/publish/xenial/xenial

upload-info:
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" "${APT_API_BASE_URL}/files" |jq
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" "${APT_API_BASE_URL}/files/trusty" |jq
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" "${APT_API_BASE_URL}/files/xenial" |jq

repo-info:
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" "${APT_API_BASE_URL}/repos" |jq
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" "${APT_API_BASE_URL}/repos/xenial" |jq
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" "${APT_API_BASE_URL}/repos/xenial/packages" |jq
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" "${APT_API_BASE_URL}/repos/trusty" |jq
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" "${APT_API_BASE_URL}/repos/trusty/packages" |jq

pub-info:
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" "${APT_API_BASE_URL}/publish" |jq

create-repos:
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" -X POST -H 'Content-Type: application/json' --data '{"Name": "trusty", "DefaultDistribution": "trusty", "DefaultComponent": "main"}' "${APT_API_BASE_URL}/repos" |jq
	curl --silent -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}" -X POST -H 'Content-Type: application/json' --data '{"Name": "xenial", "DefaultDistribution": "xenial", "DefaultComponent": "main"}' "${APT_API_BASE_URL}/repos" |jq
