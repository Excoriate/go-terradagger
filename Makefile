# Default variables
ENV ?= dev
MODULE ?= default
EXAMPLE ?= basic


# Root directories for different Terraform components
MODULE_ROOT_DIR ?= modules
EXAMPLE_ROOT_DIR ?= examples
VARS_FILE_ROOT_DIR ?= config

# Default vars file
VARS ?= fixtures.tfvars

.PHONY: default clean prune check-workdir tf-init

default:
	@echo "Available make targets ..."

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

TF_WORKING_DIR_MODULES = $(MODULE_ROOT_DIR)/$(MODULE) # modules/default
TF_WORKING_DIR_EXAMPLES = $(EXAMPLE_ROOT_DIR)/$(MODULE)/$(EXAMPLE) # examples/default/basic
WORKDIR =
TF_VARS_FILE = $(VARS_FILE_ROOT_DIR)/$(VARS) # config/fixtures.tfvars

check-workdir:
	@if [ ! -d "$(WORKDIR)" ]; then \
		echo "Module directory does not exist: $(WORKDIR)"; \
		echo "Full path: $(shell pwd)/$(WORKDIR)"; \
		exit 1; \
	fi
	@if [ -z "$(shell find $(WORKDIR) -mindepth 1 -name '*.tf' -print -quit)" ]; then \
		echo "No Terraform files found in the module directory $(WORKDIR)"; \
		exit 1; \
	fi
#####################
# Common targets #
#####################
hooks-init:
	@pre-commit install --hook-type pre-commit
	@pre-commit install --hook-type pre-push
	@pre-commit install --hook-type commit-msg
	@pre-commit autoupdate

hooks-run:
	@pre-commit run --show-diff-on-failure \
		--all-files \
		--color always

#####################
# Terraform targets #
#####################
tf-init: clean check-workdir
	@cd $(WORKDIR) && terraform init

tf-validate: check-workdir clean tf-init
	@cd $(WORKDIR) && terraform validate

tf-fmt-check: check-workdir clean
	@cd $(WORKDIR) && terraform fmt -check=true -diff=true -write=false

tf-fmt: check-workdir clean
	@cd $(WORKDIR) && terraform fmt -check=false -diff=true -write=true

tf-docs: check-workdir clean
	@cd $(WORKDIR) && terraform-docs -c .terraform-docs.yml md . > README.md

tf-lint: check-workdir clean tf-init
	@cd $(WORKDIR) && tflint -v && tflint --init && tflint

tf-plan: check-workdir clean tf-init
	@cd $(WORKDIR) && terraform plan -var-file=$(TF_VARS_FILE)

tf-apply: check-workdir clean tf-plan
	@cd $(WORKDIR) && terraform apply -var-file=$(TF_VARS_FILE)

tf-destroy: check-workdir clean tf-init
	@cd $(WORKDIR) && terraform destroy -var-file=$(TF_VARS_FILE)

#####################
# Terraform module targets #
#####################
tf-module-init: check-workdir clean tf-init
tf-module-ci: check-workdir clean tf-init tf-validate tf-fmt-check tf-lint tf-docs
tf-example-ci: check-workdir clean tf-init tf-validate tf-fmt-check tf-lint tf-docs
