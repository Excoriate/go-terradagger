<h1 align="center">
  <img alt="logo" src="https://forum.huawei.com/enterprise/en/data/attachment/forum/202204/21/120858nak5g1epkzwq5gcs.png" width="224px"/><br/>
  Terraform AWS ☁️ Events ☄️
</h1>
<p align="center">An easy to understand, opinionated terraform <b>composable</b> set of modules for managing Events, Signals and Related architectures  in <b> for AWS ☁️</b>.<br/><br/>

---

[![Auto Release](https://github.com/Excoriate/vault-labs/actions/workflows/release.yml/badge.svg)](https://github.com/Excoriate/vault-labs/actions/workflows/release.yml)
[![Terraform Check](https://github.com/Excoriate/terraform-registry-aws-accounts-creator/actions/workflows/ci-check-terraform.yml/badge.svg)](https://github.com/Excoriate/terraform-registry-aws-accounts-creator/actions/workflows/ci-check-terraform.yml)
[![Run pre-commit](https://github.com/Excoriate/terraform-registry-aws-accounts-creator/actions/workflows/ci-check-precommit.yml/badge.svg)](https://github.com/Excoriate/terraform-registry-aws-accounts-creator/actions/workflows/ci-check-precommit.yml)
[![Terratest](https://github.com/Excoriate/terraform-registry-aws-accounts-creator/actions/workflows/ci-pr-terratest.yml/badge.svg)](https://github.com/Excoriate/terraform-registry-aws-accounts-creator/actions/workflows/ci-pr-terratest.yml)


## Table of Contents

1. [About The Module](#about-the-module)
2. [Module documentation](#module-documentation)
   1. [Capabilities](#capabilities)
   2. [Getting Started](#getting-started)
   3. [Roadmap](#roadmap)
3. [Contributions](#contributing)
4. [License](#license)
5. [Contact](#contact)



<!-- ABOUT THE PROJECT -->
## About The Module

This module encapsulate a set of modules that configure, and provision accounts-related resources on AWS.

---


## Module documentation

The documentation is **automatically generated** by [terraform-docs](https://terraform-docs.io), and it's available in the module's [README.md](modules/default/README.md) file.

### Capabilities

| Module                          | Status   | Description                                                                                                           |
|---------------------------------|----------|-----------------------------------------------------------------------------------------------------------------------|
| `aws-eventbridge-rule`          | Stable ✅ | Create an event rule, that can work with several destinations..                                                       |
| `aws-eventbridge-permissions`   | Stable ✅ | Opinionated module to create custom (and usually required) permissions for an eventbridge resource, normally `rules`. |
| `aws-lambda-function`           | Stable ✅ | Create a lambda function, for different use-case.                                                                     |
| `aws-cognito-user-pool`         | Stable ✅ | Create a cognito user pool, with a set of opinionated defaults.                                                       |
| `aws-cognito-user-pool-client`  | Stable ✅ | Create a cognito user pool client, with a set of opinionated defaults.                                                |
| `aws-cognito-user-pool-domain`  | Stable ✅ | Create a cognito user pool domain, with a set of opinionated defaults.                                                |
| `aws-cognito-identity-provider` | Stable ✅ | Create a cognito identity provider, with a set of opinionated defaults.                                               |
| `aws-ses`                       | Stable ✅ | Complete SES module, that supports all current available capabilities.                                                |



### Getting Started

Check the example recipes 🥗 [here](examples)

### Roadmap

- [ ] Add more modules



## Contributing

Contributions are always encouraged and welcome! ❤️. For the process of accepting changes, please refer to the [CONTRIBUTING.md](./CONTRIBUTING.md) file, and for a more detailed explanation, please refer to this guideline [here](docs/contribution_guidelines.md).

## License

![license][badge-license]

This module is licensed under the Apache License Version 2.0, January 2004.
Please see [LICENSE] for full details.

## Contact

- 📧 **Email**: [Alex T.](mailto:alex@ideaup.cl)
- 🧳 **Linkedin**: [Alex T.](https://www.linkedin.com/in/alextorresruiz/)

_made/with_ ❤️  🤟


<!-- References -->
[LICENSE]: ./LICENSE
[badge-license]: https://img.shields.io/badge/license-Apache%202.0-brightgreen.svg

<!-- BEGIN_TF_DOCS -->
## Requirements

No requirements.

## Providers

No providers.

## Modules

No modules.

## Resources

No resources.

## Inputs

No inputs.

## Outputs

No outputs.
<!-- END_TF_DOCS -->
