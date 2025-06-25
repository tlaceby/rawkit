# ------------------------------------------------------------
#  rawkit
# ------------------------------------------------------------
# make release  →  clean → build → test → bump VERSION → git tag (pushes to Github!!)
# make verify   →  clean → build → test    (no version bump)
# ------------------------------------------------------------

UNAME_S       := $(shell uname -s)
VERSION_FILE  := .env
OLD_VER       := $(shell grep -E '^VERSION=' $(VERSION_FILE) | cut -d= -f2)
LINK_FILES    = $(wildcard link_*.go)

.PHONY: release verify nextver bump clean build libs libs-current test generate-link-files clean-link-files describe

describe:
	@echo "• rawkit version: $(OLD_VER)"

release: nextver clean-link-files verify
	@$(MAKE) clean-link-files
	@$(MAKE) VER=$(NEW_VER) generate-link-files
	@$(MAKE) NEW_VER=$(NEW_VER) bump
	@echo "✔ release pipeline succeeded for $(NEW_VER)"


verify: clean libs-current generate-link-files test
	@go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
	@gomarkdoc ./ >> docs.md
	@echo "✔ tests passed on current version $(OLD_VER)"

nextver:
	$(eval NEW_VER := $(shell echo $(OLD_VER) | awk -F. '{sub("v","", $$1); printf "v%d.%d.%d", $$1, $$2, $$3+1}'))
	@echo "• next version will be $(NEW_VER)"

clean:
	@echo "• cleaning artefacts"
	@rm -rf docs.md
	@rm -rf libs/*/*/current libs/*/*/*/*.a include/libraw/* wrapper/*.o
	@rm -f $(LINK_FILES)
	@go clean -cache -testcache

clean-link-files:
	@echo "• removing old link_*.go files"
	@rm -f $(LINK_FILES)

generate-link-files:
	@echo "• generating version-specific link files for $${VER:-$(OLD_VER)}"
	@bash ./scripts/gen_link_files.sh $${VER:-$(OLD_VER)}

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

bump:
	@test -n "$(NEW_VER)" || (echo "NEW_VER is not set" && exit 1)
	@echo "• writing VERSION=$(NEW_VER) to $(VERSION_FILE)"
	@sed -Ei.bak 's/^VERSION=.*/VERSION=$(NEW_VER)/' $(VERSION_FILE) && rm -f $(VERSION_FILE).bak
	@git add -A
	@git commit -m "release $(NEW_VER)"
	@git tag -a $(NEW_VER) -m "RawKit $(NEW_VER)"
	@git push origin HEAD:main
	@git push origin $(NEW_VER)