<h1 align="center">
  <img alt="logo" src="https://dagger.io/img/astronaut.png" width="224px"/><br/>
  TerraDagger 👨🏻‍🚀
</h1>
<p align="center">An easy to understand GO library for building portables CI/CD pipelines (as code) using <b>Dagger</b> for your <b> infrastructure-as-code ☁️</b>.<br/><br/>

---
[![Release](https://github.com/Excoriate/go-terradagger/actions/workflows/release.yaml/badge.svg)](https://github.com/Excoriate/go-terradagger/actions/workflows/release.yaml)
[![Go Build TerraDagger CLI](https://github.com/Excoriate/go-terradagger/actions/workflows/golang-build-cli.yml/badge.svg)](https://github.com/Excoriate/go-terradagger/actions/workflows/golang-build-cli.yml)
[![Go Tests Library](https://github.com/Excoriate/go-terradagger/actions/workflows/golang-tests-library.yml/badge.svg)](https://github.com/Excoriate/go-terradagger/actions/workflows/golang-tests-library.yml)
[![Go Linter Library](https://github.com/Excoriate/go-terradagger/actions/workflows/golang-linter-library.yaml/badge.svg)](https://github.com/Excoriate/go-terradagger/actions/workflows/golang-linter-library.yaml)

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
- [Terraform](https://www.terraform.io/downloads.html)

>**NOTE**: For the tools used in this project, please check the [Makefile](./Makefile), and the [Taskfile](./Taskfile.yml) files. You'll also need [pre-commit](https://pre-commit.com/) installed.

---

### Getting Started 🚀

```bash
```


## Roadmap 🗓️

- [ ] Add support for [Terragrunt](https://terragrunt.gruntwork.io/).
- [ ] Add support for [Terratest](https://terratest.gruntwork.io/).
- [ ] Mature a CLI 🤖 as a wrapper (and non-programmatic) way to use TerraDagger.
- [ ] Add official Docker images for TerraDagger.

>**Note**: This is still work in progress, however, I'll be happy to receive any feedback or contribution. Ensure you've read the [contributing guide](./CONTRIBUTING.md) before doing so.


## Contributing

Please read our [contributing guide](./CONTRIBUTING.md).

## Community

Feel free to contact me through:

- 📧 [Email](mailto:alex@ideaup.cl)
- 🧳 [Linkedin](https://www.linkedin.com/in/alextorresruiz/)


<a href="https://github.com/Excoriate/stilettov2/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=Excoriate/stiletto" />
</a>
