# ------------------------------------------------------------
#  rawkit
# ------------------------------------------------------------
# make release  →  clean → build → test → bump VERSION → git tag (pushes to Github!!)
# make verify   →  clean → build → test    (no version bump)
# ------------------------------------------------------------

UNAME_S       := $(shell uname -s)
VERSION_FILE  := .env
OLD_VER       := $(shell grep -E '^VERSION=' $(VERSION_FILE) | cut -d= -f2)

.PHONY: release verify nextver bump clean build libs libs-current test


release: clean build test bump
	gomarkdoc ./ >> docs.md
	@echo "✔ release pipeline succeeded for $(NEW_VER)"


verify: clean libs-current test
	gomarkdoc ./ >> docs.md
	@echo "✔ tests passed on current version $(OLD_VER)"


nextver:
	$(eval NEW_VER := $(shell echo $(OLD_VER) | \
	    awk -F. '{sub("v","",$$1); printf "v%d.%d.%d", $$1, $$2, $$3+1}'))
	@echo "• next version will be $(NEW_VER)"


clean:
	@echo "• cleaning artefacts"
	@rm -rf docs.md
	@rm -rf libs/*/*/current libs/*/*/*/*.a include/libraw/* wrapper/*.o
	@go clean -cache -testcache


build: nextver libs
libs: nextver
	@if [ "$(OS)" = "Windows_NT" ]; then \
	    bash ./scripts/build_windows.sh $(NEW_VER); \
	elif [ "$(UNAME_S)" = "Darwin" ]; then \
	    bash ./scripts/build_darwin.sh $(NEW_VER); \
	else \
	    bash ./scripts/build_linux.sh  $(NEW_VER); \
	fi

libs-current:
	@if [ "$(OS)" = "Windows_NT" ]; then \
	    bash ./scripts/build_windows.sh $(OLD_VER); \
	elif [ "$(UNAME_S)" = "Darwin" ]; then \
	    bash ./scripts/build_darwin.sh $(OLD_VER); \
	else \
	    bash ./scripts/build_linux.sh  $(OLD_VER); \
	fi


test:
	@echo "• running go tests"
	@go test -v ./...


bump: nextver
	@echo "• writing VERSION=$(NEW_VER) to $(VERSION_FILE)"
	@sed -Ei.bak 's/^VERSION=.*/VERSION=$(NEW_VER)/' $(VERSION_FILE) && rm -f $(VERSION_FILE).bak
	@git add $(VERSION_FILE)
	@git commit -m "release $(NEW_VER)"
	@git tag -a $(NEW_VER) -m "RawKit $(NEW_VER)"
	@git push origin HEAD:main --follow-tags
