# Changelog

## [1.3.0](https://github.com/Excoriate/go-terradagger/compare/v1.2.1...v1.3.0) (2024-02-16)


### Features

* add args validation through all core apis ([3214e0e](https://github.com/Excoriate/go-terradagger/commit/3214e0e27c80a856585bf42c63d6097c38f68095))


### Refactoring

* add more test in config and env packages ([55fde49](https://github.com/Excoriate/go-terradagger/commit/55fde49ff68666b4d35edc29b8137f54b8118a9d))
* enhance the error logic ([e518b4b](https://github.com/Excoriate/go-terradagger/commit/e518b4b32cc9bc14ba153bc59964905cae083aba))
* minor fixes, add validation for arguments, and tfvars ([9e1cbae](https://github.com/Excoriate/go-terradagger/commit/9e1cbaed92e3bbb136eb1e441a5c2e57f0b259c3))


### Other

* add CHANGELOG to ignored files for markdown linter ([081ad83](https://github.com/Excoriate/go-terradagger/commit/081ad83a3f5908d8dc2bd43a813e839526e38ffb))
* add missing tflint configurfation file in new terraform module's example ([774c35d](https://github.com/Excoriate/go-terradagger/commit/774c35dfb168b68e64777bfff8ade3aea7b73a82))
* fix terraform-docs configuration for example modules ([45ea705](https://github.com/Excoriate/go-terradagger/commit/45ea705002529cdf199eb5b66408bb0f697f5464))
* generated missing docs ([5715782](https://github.com/Excoriate/go-terradagger/commit/5715782040cd393fd32d77586078a4135424bdde))
* remove go-vet from hooks ([c83d2f3](https://github.com/Excoriate/go-terradagger/commit/c83d2f3ff55f8932a725988acf430bd60c7c9a5e))

## [1.2.1](https://github.com/Excoriate/go-terradagger/compare/v1.2.0...v1.2.1) (2024-02-11)


### Bug Fixes

* add auto-inject for terraform commands ([8930c49](https://github.com/Excoriate/go-terradagger/commit/8930c49f98489bea292ceb8bccafe5e97696db28))


### Refactoring

* add centralised command string generator ([90894bb](https://github.com/Excoriate/go-terradagger/commit/90894bbf671ba9fa76ad4785ad6ccdc5c1774e17))
* add high-level runner api ([4d04fd9](https://github.com/Excoriate/go-terradagger/commit/4d04fd99a497f499496a0594cc8d1d8428c0b163))
* add terraformcore container configurator ([b7cd389](https://github.com/Excoriate/go-terradagger/commit/b7cd38965a66a9ef3150c850e16059b831475782))
* centralise command resolver for terraform lifecycle commands ([49f9716](https://github.com/Excoriate/go-terradagger/commit/49f9716889ae05fd580cbb3c13b59d80eaf09ac4))
* centralise the image configuration resolution ([8e72843](https://github.com/Excoriate/go-terradagger/commit/8e728435281f66da60ce6ec7c1318460a0aad896))
* fix naming ([7cd7373](https://github.com/Excoriate/go-terradagger/commit/7cd737364858b9b2a2a9f427f8c32a0ddf8bac16))
* move args to separated go files, enhance naming ([460b514](https://github.com/Excoriate/go-terradagger/commit/460b51490aeb38d991abe8c5ba860746429abcc7))
* move tests. Add rest of pending apis for tg and tf ([720be28](https://github.com/Excoriate/go-terradagger/commit/720be2840cc02105dfbba9e4e9bd1522e661111f))
* remove duplicated apis, leverage runner for terraform and terragrunt ([42d124e](https://github.com/Excoriate/go-terradagger/commit/42d124e01fb12342c0a9637323c14d85262a161d))


### Other

* Add params through command flags ([bd2142c](https://github.com/Excoriate/go-terradagger/commit/bd2142ca8bf205c02d502d11eab13cb29ed87580))
* update hooks ([e116f4d](https://github.com/Excoriate/go-terradagger/commit/e116f4d1e2f63c3a5d8589fba29117f6c11f028e))


### Docs

* update documentation ([3fbfc0e](https://github.com/Excoriate/go-terradagger/commit/3fbfc0e94593c8f9156d3f9d354f6a873d4476dd))
* update getting started ([688a7eb](https://github.com/Excoriate/go-terradagger/commit/688a7ebc36816e949f88a8cd5c7ace66f7662b1c))

## [1.2.0](https://github.com/Excoriate/go-terradagger/compare/v1.1.0...v1.2.0) (2024-02-08)


### Features

* Add advanced export APIs ([3729d62](https://github.com/Excoriate/go-terradagger/commit/3729d621bf4060281dee8a00f70fd6937ba5c4cd))
* Add backup option for the exporter API ([e99eed0](https://github.com/Excoriate/go-terradagger/commit/e99eed0706a18e405c369c60df712670128b9fef))
* Add security group module ([814faf6](https://github.com/Excoriate/go-terradagger/commit/814faf62b6fbf65c6eccb34d4fbd2b7b10993086))
* add support for forwarding unix sockets, and ssh authentication ([e4440ff](https://github.com/Excoriate/go-terradagger/commit/e4440ff5e6b57057eaff494393369053691d5bb8))
* add support for terragrunt ([6270a94](https://github.com/Excoriate/go-terradagger/commit/6270a94ee65142d820af6efc031763f030235bc4))
* Add terraform core apis to be reused by different iaac tools ([9f2f0a1](https://github.com/Excoriate/go-terradagger/commit/9f2f0a150f491fe65910918c477fc333adb1d1ad))
* add terraformcore api for plan ([9a00148](https://github.com/Excoriate/go-terradagger/commit/9a00148fa8a228d6cd258ea534aad620bddb30fc))
* Add tests for dirs.go in TerraDagger ([9854117](https://github.com/Excoriate/go-terradagger/commit/9854117a981626255febd61222546c118f39cb55))


### Refactoring

* Add client configuration API ([d821f58](https://github.com/Excoriate/go-terradagger/commit/d821f58cc538a2a6ce7fbd34f685d57197cc626c))
* Add export and import functionality ([8f149d8](https://github.com/Excoriate/go-terradagger/commit/8f149d8f36a585bb8d4cff25cdc88e897eeb0ca4))
* Add import and advanced export ([c57772b](https://github.com/Excoriate/go-terradagger/commit/c57772bb19e5b4874556953d29e9c493552c418c))
* Add instance validator interface ([38d4045](https://github.com/Excoriate/go-terradagger/commit/38d4045671c9c3002df4e77a3036b34cdc43ac74))
* add simplified version of the tfinit command ([c8e1f9d](https://github.com/Excoriate/go-terradagger/commit/c8e1f9dba850becc6a32b4a898a1c1e03c9340e7))
* Add unit tests for commands, args and command's package ([79ae9e5](https://github.com/Excoriate/go-terradagger/commit/79ae9e5734bd00d833d2d73b465e866d583b81c2))
* full rewrite of terradagger ([ea40e68](https://github.com/Excoriate/go-terradagger/commit/ea40e68168b2de9696ec21e2cf838073d4c9a674))
* New logic ([134f531](https://github.com/Excoriate/go-terradagger/commit/134f5315282164b68ad499444d904f5ffdfaa7d1))


### Other

* adjust dead link in docs ([db3b788](https://github.com/Excoriate/go-terradagger/commit/db3b788fe9aacaa41054256eac9aed737493ba52))

## [1.1.0](https://github.com/Excoriate/go-terradagger/compare/v1.0.0...v1.1.0) (2023-12-24)


### Features

* Add logic for exporting directories and files ([dd2fda7](https://github.com/Excoriate/go-terradagger/commit/dd2fda78a8492d0e382ceac07271076321ee4e19))


### Refactoring

* Enhance the code a bit ([40183d1](https://github.com/Excoriate/go-terradagger/commit/40183d1318beabaf8df429bafac36eca4440fa9c))
* remove ununsed task, make golangci-linter happy ([22974d4](https://github.com/Excoriate/go-terradagger/commit/22974d4244a537f7e089bc00321430719e3deffa))

## 1.0.0 (2023-12-03)


### Features

* add commit with basic structure ([8b08bfa](https://github.com/Excoriate/go-terradagger/commit/8b08bfa7c712803e24ef6a2dba9090a3ce63ddae))
* Add fixed configuration for Golangci linter ([c9ee52c](https://github.com/Excoriate/go-terradagger/commit/c9ee52c7420acc568f508d86d38d26f5c6363b5a))
* Add makeFile ([89b9e39](https://github.com/Excoriate/go-terradagger/commit/89b9e39b08d84573dfd9958d0095c8862ebd0154))
* Add mvp for terraform commands, on Dagger ([6e5529d](https://github.com/Excoriate/go-terradagger/commit/6e5529d54774bd5f0a1ad8832bf8dc123be38a1e))
* Add TaskFile with necessary tasks ([e005292](https://github.com/Excoriate/go-terradagger/commit/e005292a1e9ba4a1233e6cd5c47f756cad4488b1))
* first commit ([e507f17](https://github.com/Excoriate/go-terradagger/commit/e507f17c775939cf5b411ad2476a91577cef9f15))


### Other

* Amend README.md ([3b0fe16](https://github.com/Excoriate/go-terradagger/commit/3b0fe160ffe743355a76ed13bb81ed3e795adadc))


### Refactoring

* Add proper names, and issue-templates ([731f627](https://github.com/Excoriate/go-terradagger/commit/731f627074a016b0a4972c727e3fb5dc5e1022bc))
* Adjust markdown-link-check configuration ([b64a328](https://github.com/Excoriate/go-terradagger/commit/b64a3286f12b1c9bc9d3caffa070aa0d75457d3f))
* Adjust names, fix dead link in contribution guideline ([128c462](https://github.com/Excoriate/go-terradagger/commit/128c4627a494b5c57a7be01f9350a74386baf58f))
* Amend incorrect name change after removed boilerplate naming convetion ([4411953](https://github.com/Excoriate/go-terradagger/commit/44119532d0ef58195000bb94cb773b72994db23f))
* Amend pre-commit configuration ([7c36959](https://github.com/Excoriate/go-terradagger/commit/7c36959ae0c528d435324d97403976fd942ac38a))
* Amend pre-commit configuration ([2618b05](https://github.com/Excoriate/go-terradagger/commit/2618b055be183bab0f0800c07d1a5136ab70688f))
* changed name of workflows ([80a6bff](https://github.com/Excoriate/go-terradagger/commit/80a6bff71edebc4c64a88259660763910d4a83a2))
* Changed names of certain workflows ([7d667d9](https://github.com/Excoriate/go-terradagger/commit/7d667d98f1afb65bc4db0e6f615cf6d33edc15bf))
* Unified pre-commit hooks ([5cd85ff](https://github.com/Excoriate/go-terradagger/commit/5cd85ff0af815f6dc597665221a1dcb20d4738f0))


### Docs

* Add example implementation in the README.md ([ee9ec77](https://github.com/Excoriate/go-terradagger/commit/ee9ec77dedd0c4947c7f963cff47ce93307eada8))
* Add updated badges to README.md ([7d76d94](https://github.com/Excoriate/go-terradagger/commit/7d76d9484c2d99cadf75fe074494120107df3f7c))
* basic structure of the readme.md ([7676e51](https://github.com/Excoriate/go-terradagger/commit/7676e5177d7a5723aa9e224f654fc8493c72d4ac))
