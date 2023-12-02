<h1 align="center">
  <img alt="logo" src="https://dagger.io/img/astronaut.png" width="224px"/><br/>
  TerraDagger 👨🏻‍🚀
</h1>
<p align="center">An easy to understand GO library for building portables CI/CD pipelines (as code) using <b>Dagger</b> for your <b> infrastructure-as-code ☁️</b>.<br/><br/>

---

[![Auto Release](https://github.com/Excoriate/vault-labs/actions/workflows/release.yml/badge.svg)](https://github.com/Excoriate/vault-labs/actions/workflows/release.yml)
[![Terraform Check](https://github.com/Excoriate/terraform-registry-aws-accounts-creator/actions/workflows/ci-check-terraform.yml/badge.svg)](https://github.com/Excoriate/terraform-registry-aws-accounts-creator/actions/workflows/ci-check-terraform.yml)
[![Run pre-commit](https://github.com/Excoriate/terraform-registry-aws-accounts-creator/actions/workflows/ci-check-precommit.yml/badge.svg)](https://github.com/Excoriate/terraform-registry-aws-accounts-creator/actions/workflows/ci-check-precommit.yml)
[![Terratest](https://github.com/Excoriate/terraform-registry-aws-accounts-creator/actions/workflows/ci-pr-terratest.yml/badge.svg)](https://github.com/Excoriate/terraform-registry-aws-accounts-creator/actions/workflows/ci-pr-terratest.yml)

---
**TerraDagger** is a **GO library** that provides a set of functions and patterns for building portable CI/CD pipelines (as code) for your infrastructure-as-code. It's based on the wonderful [Dagger](https://dagger.io) pipeline-as-code project, and heavily inspired by [Terratest](https://terratest.gruntwork.io). The problem that TerraDagger tries to solve is to provide a simple way to run your [Terraform](https://www.terraform.io/) code in a portable way, and also to provide a way to run your pipelines in a containerized way, so you can run your pipelines in any environment, and also in any CI/CD platform.

---

## Installation 🛠️

### Pre-requisites 📋

- [Go](https://golang.org/doc/install) >= 1.18
- [Docker](https://docs.docker.com/get-docker/) >= 20.10.7
- [Dagger](https://dagger.io)
- [Terraform](https://www.terraform.io/downloads.html)

>**NOTE**: For the tools used in this project, please check the [Makefile](./Makefile), and the [Taskfile](./Taskfile.yml) files. You'll also need [pre-commit](https://pre-commit.com/) installed.

---

### Using TerraDagger as a GO library 📚

```bash
go get github.com/Excoriate/terradagger
```


## Roadmap 🗓️

- [ ] Add support for [Terragrunt](https://terragrunt.gruntwork.io/).
- [ ] Add support for [Terratest](https://terratest.gruntwork.io/).
- [ ] Mature a CLI 🤖 as a wrapper (and non-programmatic) way to use TerraDagger.

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
