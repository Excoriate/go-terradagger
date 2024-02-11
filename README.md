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

- **Portable**: TerraDagger is built to be used in any CI/CD platform, and also in any environment (including your local machine).
- **Simple**: TerraDagger is built to be simple to use, if you're familiar with [Terratest](https://terratest.gruntwork.io), then you'll find this library very similar.
- **IAC Support**: Supports [Terraform](https://www.terraform.io/) and [Terragrunt](https://terragrunt.gruntwork.io/).

---

## Getting Started 🚀

Configure the terradagger client, which under the hood, will configure the Dagger client:

```go
td := terradagger.New(ctx, &terradagger.Options{
  Workspace: viper.GetString("workspace"),
})

```

Now,
it's time
to start the engine
([Dagger](https://dagger.io)).
It's important to start the engine
before running any command, so ensure that your [Docker](https://docs.docker.com/get-docker/) daemon or any compatible [OCI](https://opencontainers.org/) runtime is running.

```go
if err := td.StartEngine(); err != nil {
    return err // Handle the error properly in your code.
}

defer td.Engine.GetEngine().Close()
```

Terradagger has global options,
and also specific [terraform](https://www.terraform.io/) and [terragrunt](https://terragrunt.gruntwork.io/) options.
you can set the global options like this:

```go
tfOptions :=
   terraformcore.WithOptions(td, &terraformcore.TfOptions{
     ModulePath:                   viper.GetString("module"),
     EnableSSHPrivateGit:          true,
     TerraformVersion:             viper.GetString("terraform-version"),
     EnvVarsToInjectByKeyFromHost: []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_SESSION_TOKEN"},
   })

```

There are many options supported, options that are meant to facilitate the use of containerized pipelines based on common use-cases, such as:

- Injecting environment variables from the host to the container.
- Auto-injecting the AWS credentials from the host to the container.
- Forward your SSH agent to the container, so you can use your SSH keys in the container.

And then, you're good to go and run your desired [Terraform](https://www.terraform.io/) commands, and chain them as you wish:

```go
_, tfInitErr := terraform.InitE(td, tfOptions, terraform.InitOptions{})
if tfInitErr != nil {
    return tfInitErr
}

```

>NOTE: The `E` suffix in the function name means that the specific [terraform](https://www.terraform.io/) command will return the `stdout` and an [error object](https://golang.org/pkg/errors/). The variant without the `E` suffix will return the actual [Dagger Container](https://pkg.go.dev/github.com/excoriate/dagger/pkg/container) object, and an [error object](https://golang.org/pkg/errors/).


To see a full working example, please check the [**terradagger-cli**](cli/) that's built in this repository

---

## Roadmap 🗓️

- [x] Add basic support for Terraform commands (init, validate, plan, apply, destroy, etc).
- [ ] Add out-of-the-box support for TfLint.
- [ ] Add extra commands: Validate, Format, and Import.
- [ ] Add plenty of missing tests 🧪
- [x] Add support for [Terragrunt](https://terragrunt.gruntwork.io/).
- [ ] Enrich the [terragrunt](https://terragrunt.gruntwork.io/) API to cover all the commands supported.
- [ ] Add support for [Terratest](https://terratest.gruntwork.io/).
- [ ] Add official Docker images for TerraDagger.

>**Note**: This is still work in progress, however, I'll be happy to receive any feedback or contribution. Ensure you've read the [contributing guide](./CONTRIBUTING.md) before doing so.


## Contributing

Please read our [contributing guide](./CONTRIBUTING.md).
