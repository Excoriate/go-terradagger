<h1 align="center">
  <img alt="logo" src="docs/logo/terradagger-logo-2.png" width="450px"/><br/>
  TerraDagger 🗡️
</h1>
<p align="center">An easy to understand GO library for building portables CI/CD pipelines (as code) using <b>Dagger</b> for your <b> infrastructure-as-code ☁️</b>.<br/><br/>

---
[![Release](https://github.com/Excoriate/go-terradagger/actions/workflows/release.yaml/badge.svg)](https://github.com/Excoriate/go-terradagger/actions/workflows/release.yaml)
[![Go Build](https://github.com/Excoriate/go-terradagger/actions/workflows/go-build.yml/badge.svg)](https://github.com/Excoriate/go-terradagger/actions/workflows/go-build.yml)
[![Go Linter](https://github.com/Excoriate/go-terradagger/actions/workflows/go-ci-lint.yaml/badge.svg)](https://github.com/Excoriate/go-terradagger/actions/workflows/go-ci-lint.yaml)
[![Go Tests](https://github.com/Excoriate/go-terradagger/actions/workflows/go-ci-tests.yml/badge.svg)](https://github.com/Excoriate/go-terradagger/actions/workflows/go-ci-tests.yml)

---
**TerraDagger** is a **GO library** that provides a set of functions and patterns for building portable CI/CD pipelines (as code) for your infrastructure-as-code. It's based on the wonderful [Dagger](https://dagger.io) pipeline-as-code project, and heavily inspired by [Terratest](https://terratest.gruntwork.io). The problem that TerraDagger tries to solve is to provide a simple way to run your [Terraform](https://www.terraform.io/) code in a portable way, and also to provide a way to run your pipelines in a containerized way, so you can run your pipelines in any environment, and also in any CI/CD platform.

---

## Installation 🛠️

Install it using [Go get](https://golang.org/cmd/go/#hdr-Add_dependencies_to_current_module_and_install_them):

```bash
go get github.com/Excoriate/go-terradagger
```

### Pre-requisites 📋

- [Go](https://golang.org/doc/install) >= 1.18
- [Docker](https://docs.docker.com/get-docker/) >= 20.10.7
- [Dagger](https://dagger.io)

>**NOTE**: For the tools used in this project, please check the [Makefile](./Makefile), and the [Taskfile](./Taskfile.yml) files. You'll also need [pre-commit](https://pre-commit.com/) installed.

---

## Features 🎉
* **Portable**: TerraDagger is built to be used in any CI/CD platform, and also in any environment (including your local machine).
* **Simple**: TerraDagger is built to be simple to use, if you're familiar with [Terratest](https://terratest.gruntwork.io), then you'll find this library very similar.
* **IAC Support**: Supports [Terraform](https://www.terraform.io/) and [Terragrunt](https://terragrunt.gruntwork.io/).

---

## Getting Started 🚀


```go
td, err := terradagger.New(ctx, &terradagger.ClientOptions{
	RootDir: "../",
})

defer td.DaggerClient.Close()

if err != nil {
	return err
}

terraformOptions := &terraform.Options{
	TerraformDir: "test-data/terraform/root-module-1",
}

_ = terraform.Init(td, terraformOptions, nil)
_ = terraform.Plan(td, terraformOptions, &terraform.PlanOptions{
	Vars: map[string]interface{}{
		"is_enabled": true,
	},
})
_ = terraform.Apply(td, terraformOptions, &terraform.ApplyOptions{
	Vars: map[string]interface{}{
		"is_enabled": true,
	},
})
_ = terraform.Destroy(td, terraformOptions, &terraform.DestroyOptions{
	Vars: map[string]interface{}{
		"is_enabled": true,
	},
})

```

>**NOTE**: To see a full working example, please check the [**terradagger-cli**](cli/) that's built in this repository



## Roadmap 🗓️

- [x] Add basic support for Terraform commands (init, validate, plan, apply, destroy, etc).
- [ ] Add plenty of missing tests 🧪
- [ ] Add support for [Terragrunt](https://terragrunt.gruntwork.io/).
- [ ] Add support for [Terratest](https://terratest.gruntwork.io/).
- [ ] Mature a CLI 🤖 as a wrapper (and non-programmatic) way to use TerraDagger.
- [ ] Add official Docker images for TerraDagger.

>**Note**: This is still work in progress, however, I'll be happy to receive any feedback or contribution. Ensure you've read the [contributing guide](./CONTRIBUTING.md) before doing so.


## Contributing

Please read our [contributing guide](./CONTRIBUTING.md).
