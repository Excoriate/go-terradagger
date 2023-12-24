# Default variables
ENV ?= dev

# Root directories for different Terraform components
MODULE ?= root-module-1
VARS ?= ""
TERRAFORM_TESTS_MODULES_DIR ?= tests/terraform

# Tools, and scripts.
SCRIPTS_FOLDER ?= scripts
GO=go
GO_BUILD_SCRIPT := $(shell pwd)/scripts/golang/go_build.sh
TERRADAGGER_CLI_SRC_DIR := $(shell pwd)/cli
TERRADAGGER_CLI_NAME := terradagger-cli

.PHONY: default clean prune check-workdir tf-init

clean:
	@echo "Cleaning directories..."
	@if [ -d "$(MODULE_ROOT_DIR)" ]; then \
		find . -type d -name ".terraform" -exec echo "Removing {}" \; -exec rm -rf '{}' \;; \
		echo "Cleaned .terraform directories."; \
	else \
		echo "$(MODULE_ROOT_DIR) directory not found, skipping cleanup of .terraform directories."; \
	fi
	@if [ -d "$(MODULE_ROOT_DIR)" ]; then \
		find . -type d -name ".terragrunt-cache" -exec echo "Removing {}" \; -exec rm -rf '{}' \;; \
		echo "Cleaned .terragrunt-cache directories."; \
	else \
		echo "$(MODULE_ROOT_DIR) directory not found, skipping cleanup of .terragrunt-cache directories."; \
	fi
	@if [ -d "$(MODULE_ROOT_DIR)" ]; then \
		find . -type f -name "terraform.tfstate" -exec echo "Removing {}" \; -exec rm -rf '{}' \;; \
		echo "Removed terraform.tfstate files."; \
	else \
		echo "$(MODULE_ROOT_DIR) directory not found, skipping removal of terraform.tfstate files."; \
	fi
	@if [ -d "$(MODULE_ROOT_DIR)" ]; then \
		find . -type f -name "terraform.tfstate.backup" -exec echo "Removing {}" \; -exec rm -rf '{}' \;; \
		echo "Removed terraform.tfstate.backup files."; \
	else \
		echo "$(MODULE_ROOT_DIR) directory not found, skipping removal of terraform.tfstate.backup files."; \
	fi
	@if [ -d "$(MODULE_ROOT_DIR)" ]; then \
		find . -type f -name "terraform.tfplan" -exec echo "Removing {}" \; -exec rm -rf '{}' \;; \
		echo "Removed terraform.tfplan files."; \
	else \
		echo "$(MODULE_ROOT_DIR) directory not found, skipping removal of terraform.tfplan files."; \
	fi

prune: clean
	@git clean -f -xd --exclude-list

#####################
# Common targets #
#####################
pc-init:
	@pre-commit install --hook-type pre-commit
	@pre-commit install --hook-type pre-push
	@pre-commit install --hook-type commit-msg
	@pre-commit autoupdate

pc-run:
	@pre-commit run --show-diff-on-failure \
		--all-files \
		--color always

#####################
# Terraform targets #
#####################
tf-init: clean
	@cd $(TERRAFORM_TESTS_MODULES_DIR)/$(MODULE) && terraform init

tf-validate: clean tf-init
	@cd $(TERRAFORM_TESTS_MODULES_DIR)/$(MODULE) && terraform validate

tf-fmt-check: clean
	@cd $(TERRAFORM_TESTS_MODULES_DIR)/$(MODULE) && terraform fmt -check=true -diff=true -write=false

tf-fmt: clean
	@cd $(TERRAFORM_TESTS_MODULES_DIR)/$(MODULE) && terraform fmt -check=false -diff=true -write=true

tf-docs: clean
	@cd $(TERRAFORM_TESTS_MODULES_DIR)/$(MODULE) && terraform-docs -c .terraform-docs.yml md . > README.md

tf-lint: clean tf-init
	@cd $(TERRAFORM_TESTS_MODULES_DIR)/$(MODULE) && tflint -v && tflint --init && tflint

tf-plan: clean tf-init
	@if [ -z "$(VARS)" ]; then \
		echo "No vars file provided, skipping terraform plan with vars."; \
		cd $(TERRAFORM_TESTS_MODULES_DIR)/$(MODULE) && terraform plan; \
	else \
		cd $(TERRAFORM_TESTS_MODULES_DIR)/$(MODULE) && terraform plan -var-file=$(VARS); \
	fi

tf-apply: clean tf-plan
	@if [ -z "$(VARS)" ]; then \
		echo "No vars file provided, skipping terraform apply with vars."; \
		cd $(TERRAFORM_TESTS_MODULES_DIR)/$(MODULE) && terraform apply -auto-approve; \
	else \
		cd $(TERRAFORM_TESTS_MODULES_DIR)/$(MODULE) && terraform apply -auto-approve -var-file=$(VARS); \
	fi

tf-destroy: clean tf-init
	@if [ -z "$(VARS)" ]; then \
		echo "No vars file provided, skipping terraform destroy with vars."; \
		cd $(TERRAFORM_TESTS_MODULES_DIR)/$(MODULE) && terraform destroy -auto-approve; \
	else \
		cd $(TERRAFORM_TESTS_MODULES_DIR)/$(MODULE) && terraform destroy -auto-approve -var-file=$(VARS); \
	fi

#####################
# Go targets #
#####################
## tidy: tidy go.mod
.PHONY: go-tidy
go-tidy:
	@$(GO) mod tidy

## fmt: Run go fmt against code.
.PHONY: go-fmt
go-fmt:
	@$(GO) fmt -x ./...

## vet: Run go vet against code.
.PHONY: go-vet
go-vet:
	@$(GO) vet ./...

## lint: Run go lint against code.
.PHONY: go-lint
go-lint:
	@golangci-lint run -v --config .golangci.yml

## style: Code style -> fmt,vet,lint
.PHONY: go-style
go-style: go-fmt go-vet go-lint

## test: Run unit test
.PHONY: go-test
go-test:
	@echo "===========> Run unit test"
	@$(GO) test -race -v ./...

.PHONY: go-test-nocache
go-test-nocache:
	@echo "===========> Run unit test without cache"
	@go test ./... -count=1

## Build Go Binary
.PHONY: go-build
cli-build:
	@echo "===========> Building binary"
	@cd $(TERRADAGGER_CLI_SRC_DIR) &&  $(GO_BUILD_SCRIPT) --binary $(TERRADAGGER_CLI_NAME) --path ./main.go

## Run Go source code
.PHONY: go-run
cli-run:
	@echo "===========> Running binary of the terradagger-cli"
	@cd $(TERRADAGGER_CLI_SRC_DIR) && ./$(TERRADAGGER_CLI_NAME) $(ARGS)

.PHONY: go-run
cli-run-src:
	@echo "===========> Running source code"
	@$(GO) run $(TERRADAGGER_CLI_SRC_DIR)/main.go $(ARGS)
