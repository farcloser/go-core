#   Copyright Farcloser.

#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at

#       http://www.apache.org/licenses/LICENSE-2.0

#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

ORG_PREFIXES := "go.farcloser.world"
ICON := "🧿"

MAKEFILE_DIR := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
VERSION ?= $(shell git -C $(MAKEFILE_DIR) describe --match 'v[0-9]*' --dirty='.m' --always --tags 2>/dev/null \
	|| echo "no_git_information")
VERSION_TRIMMED := $(VERSION:v%=%)
REVISION ?= $(shell git -C $(MAKEFILE_DIR) rev-parse HEAD 2>/dev/null || echo "no_git_information")$(shell \
	if ! git -C $(MAKEFILE_DIR) diff --no-ext-diff --quiet --exit-code 2>/dev/null; then echo .m; fi)
LINT_COMMIT_RANGE ?= main..HEAD

ifdef VERBOSE
	VERBOSE_FLAG := -v
	VERBOSE_FLAG_LONG := --verbose
endif

ifndef NO_COLOR
    NC := \033[0m
    GREEN := \033[1;32m
    ORANGE := \033[1;33m
    BLUE := \033[1;34m
    RED := \033[1;31m
endif

# Helpers
recursive_wildcard=$(wildcard $1$2) $(foreach e,$(wildcard $1*),$(call recursive_wildcard,$e/,$2))

define title
	@printf "$(GREEN)____________________________________________________________________________________________________\n" 1>&2
	@printf "$(GREEN)%*s\n" $$(( ( $(shell echo "$(ICON)$(1) $(ICON)" | wc -c ) + 100 ) / 2 )) "$(ICON)$(1) $(ICON)" 1>&2
	@printf "$(GREEN)____________________________________________________________________________________________________\n$(ORANGE)" 1>&2
endef

define footer
	@printf "$(GREEN)> %s: done!\n" "$(1)" 1>&2
	@printf "$(GREEN)____________________________________________________________________________________________________\n$(NC)" 1>&2
endef

# Tasks
lint: lint-go-all lint-commits lint-mod lint-licenses-all lint-headers lint-yaml lint-shell

fix: fix-mod fix-go-all

test: unit

unit: test-unit test-unit-race test-unit-bench

##########################
# Linting tasks
##########################
lint-go:
	$(call title, $@)
	@cd $(MAKEFILE_DIR) \
		&& golangci-lint run $(VERBOSE_FLAG_LONG) ./...
	$(call footer, $@)

lint-go-all:
	$(call title, $@)
	@cd $(MAKEFILE_DIR) \
		&& GOOS=darwin make lint-go \
		&& GOOS=linux make lint-go \
		&& GOOS=freebsd make lint-go \
		&& GOOS=windows make lint-go
	$(call footer, $@)

lint-yaml:
	$(call title, $@)
	@cd $(MAKEFILE_DIR) \
		&& yamllint .
	$(call footer, $@)

lint-shell: $(call recursive_wildcard,$(MAKEFILE_DIR)/,*.sh)
	$(call title, $@)
	@shellcheck -a -x $^
	$(call footer, $@)

# See https://github.com/andyfeller/gh-ssh-allowed-signers for automation to retrieve contributors keys
lint-commits:
	$(call title, $@)
	@cd $(MAKEFILE_DIR) \
		&& git config --add gpg.ssh.allowedSignersFile hack/allowed_signers \
		&& git-validation $(VERBOSE_FLAG) -run DCO,short-subject,dangling-whitespace -range "$(LINT_COMMIT_RANGE)"
	$(call footer, $@)

lint-headers:
	$(call title, $@)
	@cd $(MAKEFILE_DIR) \
		&& ltag -t "./hack/headers" --check -v
	$(call footer, $@)

lint-mod:
	$(call title, $@)
	@cd $(MAKEFILE_DIR) \
		&& go mod tidy --diff
	$(call footer, $@)

# FIXME: go-licenses cannot find LICENSE from root of repo when submodule is imported:
# https://github.com/google/go-licenses/issues/186
# This is impacting gotest.tools
lint-licenses:
	$(call title, $@: $(GOOS))
	@cd $(MAKEFILE_DIR) \
		&& go-licenses check --include_tests --allowed_licenses=Apache-2.0,BSD-2-Clause,BSD-3-Clause,MIT,MPL-2.0 \
		  --ignore gotest.tools \
		  ./...
	$(call footer, $@)

lint-licenses-all:
	$(call title, $@)
	@cd $(MAKEFILE_DIR) \
		&& GOOS=darwin make lint-licenses \
		&& GOOS=linux make lint-licenses \
		&& GOOS=freebsd make lint-licenses \
		&& GOOS=windows make lint-licenses
	$(call footer, $@)

##########################
# Automated fixing tasks
##########################
fix-go:
	$(call title, $@)
	@cd $(MAKEFILE_DIR) \
		&& golangci-lint run --fix
	$(call footer, $@)

fix-go-all:
	$(call title, $@)
	@cd $(MAKEFILE_DIR) \
		&& GOOS=darwin make fix-go \
		&& GOOS=linux make fix-go \
		&& GOOS=freebsd make fix-go \
		&& GOOS=windows make fix-go
	$(call footer, $@)

fix-mod:
	$(call title, $@)
	@cd $(MAKEFILE_DIR) \
		&& go mod tidy
	$(call footer, $@)

up:
	$(call title, $@)
	@cd $(MAKEFILE_DIR) \
		&& go get -u ./...
	$(call footer, $@)

##########################
# Development tools installation
##########################
install-dev-gotestsum:
	# gotestsum: 1.12.1 (2025-03-15)
	$(call title, $@)
	@cd $(MAKEFILE_DIR) \
		&& go install gotest.tools/gotestsum@3f7ff0ec4aeb6f95f5d67c998b71f272aa8a8b41
	$(call footer, $@)

install-dev-tools: install-dev-gotestsum
	$(call title, $@)
	# golangci: v2.0.2 (2024-03-26)
	# git-validation: main (2025-02-25)
	# ltag: main (2025-03-04)
	# go-licenses: v2.0.0-alpha.1 (2024-06-27)
	@cd $(MAKEFILE_DIR) \
		&& go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@2b224c2cf4c9f261c22a16af7f8ca6408467f338 \
		&& go install github.com/vbatts/git-validation@7b60e35b055dd2eab5844202ffffad51d9c93922 \
		&& go install github.com/containerd/ltag@66e6a514664ee2d11a470735519fa22b1a9eaabd \
		&& go install github.com/google/go-licenses/v2@d01822334fba5896920a060f762ea7ecdbd086e8
	@echo "Remember to add \$$HOME/go/bin to your path"
	$(call footer, $@)

test-unit:
	$(call title, $@)
	@go test $(VERBOSE_FLAG) -count 1 $(MAKEFILE_DIR)/...
	$(call footer, $@)

test-unit-bench:
	$(call title, $@)
	@go test $(VERBOSE_FLAG) -count 1 $(MAKEFILE_DIR)/... -bench=.
	$(call footer, $@)

test-unit-race:
	$(call title, $@)
	@CGO_ENABLED=1 go test $(VERBOSE_FLAG) $(MAKEFILE_DIR)/... -race
	$(call footer, $@)

.PHONY: \
	lint \
	fix \
	test \
	up \
	unit \
	install-dev-tools \
	lint-commits lint-go lint-go-all lint-headers lint-licenses lint-licenses-all lint-mod lint-shell lint-yaml \
	fix-go fix-go-all fix-mod \
	test-unit test-unit-race test-unit-bench