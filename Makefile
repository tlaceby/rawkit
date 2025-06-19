# ------------------------------------------------------------
#  rawkit — simplified
# ------------------------------------------------------------

VERSION_FILE  := .env
OLD_VER       := $(shell grep -E '^VERSION=' $(VERSION_FILE) | cut -d= -f2)

.PHONY: release verify nextver clean ensure_current test bump

# ----------------------------------------------------------------------------
# make verify → clean → populate current/ → go test
# ----------------------------------------------------------------------------
verify: clean ensure_current test
	@echo "✔ tests passed on current version $(OLD_VER)"

# ----------------------------------------------------------------------------
# make release → clean → populate current/ → go test → bump VERSION → git tag
# ----------------------------------------------------------------------------
release:
	@go generate ./...
	@make clean libs-current test
	@echo "✔ release pipeline succeeded for $(NEW_VER)"

# ----------------------------------------------------------------------------
# bump next version
# ----------------------------------------------------------------------------
nextver:
	$(eval NEW_VER := $(shell echo $(OLD_VER) | \
	  awk -F. '{sub("v","",$$1); printf "v%d.%d.%d", $$1, $$2, $$3+1}'))
	@echo "• next version will be $(NEW_VER)"

# ----------------------------------------------------------------------------
# drop out-of-date artefacts (but leave current/ scripts happy)
# ----------------------------------------------------------------------------
clean:
	@echo "• cleaning artefacts"
	@rm -rf libs/*/*/current include/libraw/* wrapper/*.o
	@go clean -cache -testcache

# ----------------------------------------------------------------------------
# ensure libs/*/*/current is populated (build or fetch)
# ----------------------------------------------------------------------------
ensure_current:
	@echo "• ensuring current/ libs are in place"
	@bash scripts/ensure_current.sh

# ----------------------------------------------------------------------------
# run Go tests
# ----------------------------------------------------------------------------
test:
	@echo "• running go tests"
	@go test -v ./...

# ----------------------------------------------------------------------------
# bump VERSION in .env, commit & tag
# ----------------------------------------------------------------------------
bump: nextver
	@echo "• writing VERSION=$(NEW_VER) to $(VERSION_FILE)"
	@sed -Ei.bak 's/^VERSION=.*/VERSION=$(NEW_VER)/' $(VERSION_FILE) && rm -f $(VERSION_FILE).bak
	@git add $(VERSION_FILE)
	@git commit -m "release $(NEW_VER)"
	@git tag -a $(NEW_VER) -m "RawKit $(NEW_VER)"
	@git push origin HEAD:main --follow-tags
