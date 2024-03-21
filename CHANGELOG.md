# Changelog

## [0.1.0](https://github.com/glasskube/glasskube/compare/v0.0.4...v0.1.0) (2024-03-21)


### Features

* add cached package client ([4c2a18f](https://github.com/glasskube/glasskube/commit/4c2a18fc9785ff87e8e2fa214145f8edf169b154))
* add generic caching to repo client ([c69eff9](https://github.com/glasskube/glasskube/commit/c69eff901e03c15ad4451d6059ac4d1d7850f1d0))
* add version mismatch detection for package-operator ([#373](https://github.com/glasskube/glasskube/issues/373)) ([af2e2ed](https://github.com/glasskube/glasskube/commit/af2e2ed364079b44f3ecab89a18761c310db1345))
* **cli, ui:** check dependencies and show newly installed packages at update ([#113](https://github.com/glasskube/glasskube/issues/113), [#114](https://github.com/glasskube/glasskube/issues/114)) ([1b69fde](https://github.com/glasskube/glasskube/commit/1b69fde8cfb7c68906237142c257e581e17393ae))
* **cli:** add `--latest` flag to bootstrap command ([#361](https://github.com/glasskube/glasskube/issues/361)) ([11cbca8](https://github.com/glasskube/glasskube/commit/11cbca899b56d0eb0df97a647cbe34047d331e62))
* **cli:** add showing dependencies in `describe` command ([a495ebc](https://github.com/glasskube/glasskube/commit/a495ebcae88da6abf6cba0fd817d2fe987ce0843))
* **cli:** add validating dependencies in install command ([f8f72c3](https://github.com/glasskube/glasskube/commit/f8f72c3368c28213688fee06806785cdb619324d))
* **cli:** change `describe` command to be more clear ([9e3faf3](https://github.com/glasskube/glasskube/commit/9e3faf384955e49769265d36065f27cba3064195))
* **cli:** show entrypoints in glasskube describe command ([#346](https://github.com/glasskube/glasskube/issues/346)) ([fb0a824](https://github.com/glasskube/glasskube/commit/fb0a8241896a5383ce2766428021ea0c2b351d3c))
* mandatory package version ([#341](https://github.com/glasskube/glasskube/issues/341)) ([7cc7ba8](https://github.com/glasskube/glasskube/commit/7cc7ba8550d8581fd0bdb570a30d946a3b1aee16))
* **package-operator:** add validating webhook and cert generation ([328dd58](https://github.com/glasskube/glasskube/commit/328dd58a1aaa9faff37e0c19b16dca2c7526ab93))
* **package-operator:** prevent invalid updates and deletions ([#364](https://github.com/glasskube/glasskube/issues/364)) ([26e3ddb](https://github.com/glasskube/glasskube/commit/26e3ddb9a320d658834386596db7e7c1b7877eab))
* prune package dependencies ([#318](https://github.com/glasskube/glasskube/issues/318)) ([4c2af36](https://github.com/glasskube/glasskube/commit/4c2af3611d41b63895a11b70a27cf7c1eb7abada))
* **ui:** add syncing update notification via websocket ([c249fa1](https://github.com/glasskube/glasskube/commit/c249fa12ed660c7e419e489c51fe0b4a79a85c35))
* **ui:** handle dependencies at installation ([#114](https://github.com/glasskube/glasskube/issues/114)) ([3f6f510](https://github.com/glasskube/glasskube/commit/3f6f510cae031fe8581f08479b05fdfc91fe2d27))
* **ui:** show a notification when "open" fails ([#393](https://github.com/glasskube/glasskube/issues/393)) ([bbe36dc](https://github.com/glasskube/glasskube/commit/bbe36dc1e0c0c0e6381c0acf075490586373a8fd))
* **ui:** use updater for single package update check ([90b7661](https://github.com/glasskube/glasskube/commit/90b7661bf39b51e0bdc3d46953204dc80ab7ca23))


### Bug Fixes

* **cli, ui:** add support for bootstrap latest via ui, implicit request ([77f72f2](https://github.com/glasskube/glasskube/commit/77f72f2d49367ffa241c06419da6d558e0d115fc))
* **cli:** update check fails in dev environment ([e2f2d2d](https://github.com/glasskube/glasskube/commit/e2f2d2d657189fe67f83b8cbc2ffe2d6b1cf0936))
* consider semver metadata when comparing installed vs latest version ([#397](https://github.com/glasskube/glasskube/issues/397)) ([5121a95](https://github.com/glasskube/glasskube/commit/5121a95f461ad512e388034b655a7230d601a85e))
* **deps:** update dependency htmx.org to v1.9.11 ([86eec99](https://github.com/glasskube/glasskube/commit/86eec99c029e3c46957e22abf89d0ff8af9d7b71))
* **deps:** update kubernetes packages to v0.29.3 ([510907e](https://github.com/glasskube/glasskube/commit/510907e4ebb1077b7dd0d4efec6e0f73683c5c07))
* **deps:** update module github.com/google/go-containerregistry to v0.19.1 ([fc868b1](https://github.com/glasskube/glasskube/commit/fc868b169321dccafac6aa1c3d6c17d1e189ea29))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.17.0 ([04fa04f](https://github.com/glasskube/glasskube/commit/04fa04fc55cf3477b87cbd267ae514397ddd999f))
* **deps:** update module github.com/onsi/gomega to v1.32.0 ([48c2669](https://github.com/glasskube/glasskube/commit/48c2669ba42b417e1a7bf1d5bd889b25c1f8a8f6))
* **package-operator:** add missing permissions for webhook-cert ([82c5cd0](https://github.com/glasskube/glasskube/commit/82c5cd054c6d1e7388a98a69fce7b80898bf8ed0))
* **ui:** change server informer to also broadcast after errors ([497c2b8](https://github.com/glasskube/glasskube/commit/497c2b87dd6313402d83eb800a23ce23bd69505a))
* **ui:** make link color readable on dark mode ([28991c9](https://github.com/glasskube/glasskube/commit/28991c94a1e8c894019fcebdbf7856076eb1e64d))


### Other

* change config manifest structure to use default ([0fb1a4d](https://github.com/glasskube/glasskube/commit/0fb1a4dc1f8e4c988e6507716dfbe0a3862461c4))
* change next release to 0.1.0 ([e69df6b](https://github.com/glasskube/glasskube/commit/e69df6b52103e81794e3df6343a1f1d7402dba23))
* **deps:** update actions/cache action to v4.0.2 ([0dfcbaa](https://github.com/glasskube/glasskube/commit/0dfcbaad09ab05b068b5285fdcb612646026c39b))
* **deps:** update actions/checkout digest to b4ffde6 ([e11aeec](https://github.com/glasskube/glasskube/commit/e11aeec174323ca6f4f62457f63ed219aaf8af9f))
* **deps:** update dependency @commitlint/cli to v19.2.0 ([f7c7ac1](https://github.com/glasskube/glasskube/commit/f7c7ac1936dbf02413cea4c1f66caab3929c5c88))
* **deps:** update dependency @commitlint/cli to v19.2.1 ([5e2d7f6](https://github.com/glasskube/glasskube/commit/5e2d7f6f14590e20cc3e6887cc6e2a2e07cddc4a))
* **deps:** update dependency esbuild to v0.20.2 ([d079bf2](https://github.com/glasskube/glasskube/commit/d079bf2383ffa41db37dab9fcf89c20fecfbe909))
* **deps:** update dependency typescript to v5.4.3 ([1f217c9](https://github.com/glasskube/glasskube/commit/1f217c9702b610c8a83626b6bf2d71cffe7a9bbb))
* **deps:** update docker/login-action digest to e92390c ([fd0b915](https://github.com/glasskube/glasskube/commit/fd0b9157b32a4b48d6c58f8618ea89177cf214cd))
* **deps:** update ghcr.io/glasskube/package-operator docker tag to v0.0.4 ([30fec9c](https://github.com/glasskube/glasskube/commit/30fec9c5bfe73e40bb1886ac59cea9951979330b))
* remove redundant image in manager kustomization ([9b9ea29](https://github.com/glasskube/glasskube/commit/9b9ea29db82e6fe8401cea8f3c1c4a0ded2b4803))
* update release-please extra files ([48f6d43](https://github.com/glasskube/glasskube/commit/48f6d431c85186e274f8c74c7d5285c2649dfee1))


### Docs

* add web development section in contributing guide ([8876f98](https://github.com/glasskube/glasskube/commit/8876f982ad61f8e82befdbe1c450083215ebef48))
* describe dependencies and changes to package repo ([#316](https://github.com/glasskube/glasskube/issues/316)) ([f34dfb4](https://github.com/glasskube/glasskube/commit/f34dfb4ae7d423badeeceb8feb7c79a0cca2858f))
* update contributing guide on how to run operator ([5f03733](https://github.com/glasskube/glasskube/commit/5f037338c1f4a2813c9f9d5e4186b3d9bbef1407))
* update guides sidebar ([2fced23](https://github.com/glasskube/glasskube/commit/2fced239a6be69515afb77840e637099f1a3fd5e))
* updates guides sidebar title ([761faec](https://github.com/glasskube/glasskube/commit/761faecb8ea715261b19257b037c79609cadacbc))
* versioning ([b526cf1](https://github.com/glasskube/glasskube/commit/b526cf14f51ca86b348d6803f487d6350fbe1305))
* **website:** add hs-script to newsletter ([6b6a261](https://github.com/glasskube/glasskube/commit/6b6a2618206f8e28463b6b6cbaa7c011050ed0c0))


### Refactoring

* add repo client interface and fake client for testing ([5a9fc7a](https://github.com/glasskube/glasskube/commit/5a9fc7af3925fa1d5d9fa5bc76269a0c1060c037))

## [0.0.4](https://github.com/glasskube/glasskube/compare/v0.0.3...v0.0.4) (2024-03-12)


### Features

* **cli:** add --no-await in install command ([51bc1be](https://github.com/glasskube/glasskube/commit/51bc1be5444deb722fa96d26804e3cf8b2560aec))
* **cli:** add --no-await in uninstall command ([a166b36](https://github.com/glasskube/glasskube/commit/a166b36460ff7ce6c903f244c8baa1c69af50f12))
* **client:** parallelize IsBootstrap function ([8d8559e](https://github.com/glasskube/glasskube/commit/8d8559e0fc938ff54ea348efc02095579dc3d24b))
* **package-controller:** install package dependencies ([#111](https://github.com/glasskube/glasskube/issues/111)) ([3c83668](https://github.com/glasskube/glasskube/commit/3c83668ea65d240b678864e522d98a72f87d4cee))
* **package-controller:** support versions for dependent packages ([#311](https://github.com/glasskube/glasskube/issues/311)) ([e22a401](https://github.com/glasskube/glasskube/commit/e22a4019405e2beb35595484b9207d0ad7b00350))
* **package-operator:** add default namespace handling for manifests ([97e17a5](https://github.com/glasskube/glasskube/commit/97e17a5ff8bbb57f73fbcd9d39d1700797b0be6d))
* **package-operator:** add handling of packages with helm and manifests ([149a8c7](https://github.com/glasskube/glasskube/commit/149a8c7e5d7deb7f192cb76e189e4f2b42a6bf3c))
* **ui:** add dark mode ([60ad43e](https://github.com/glasskube/glasskube/commit/60ad43ebda9cb9d34fdebb80635ec2667d12db15))


### Bug Fixes

* change opener to choose first ready pod ([b018081](https://github.com/glasskube/glasskube/commit/b018081693ce689c26c2fe11a7dd1de4293d95a9))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.16.0 ([b268277](https://github.com/glasskube/glasskube/commit/b26827731447b2f7e04152057b1fede9bd753f5d))
* **ui, cli:** add fallback to manifest from repo ([0e5d024](https://github.com/glasskube/glasskube/commit/0e5d0243244be949db6ac502db4eb33f3a835332))
* **ui:** hide global update notification ([#355](https://github.com/glasskube/glasskube/issues/355)) ([d2e60ca](https://github.com/glasskube/glasskube/commit/d2e60cadc657942ce23ac23da7210381faae2266))


### Other

* add commitlint to root project ([54fd3a9](https://github.com/glasskube/glasskube/commit/54fd3a90348f7b38dd3da8ea2a33b048d5517f6c))
* add prettier for web formatting ([1493268](https://github.com/glasskube/glasskube/commit/1493268e46bb72c81dba19d430e1d0ccf501045f))
* **deps:** update actions/checkout digest to 9bb5618 ([12303a2](https://github.com/glasskube/glasskube/commit/12303a22b2f1e526a2379295958936295d8f7f94))
* **deps:** update azure/setup-kubectl action to v4 ([cbb4712](https://github.com/glasskube/glasskube/commit/cbb4712254115dbfba3f12fa91ce7b85b54262eb))
* **deps:** update commitlint monorepo to v19.1.0 ([170e905](https://github.com/glasskube/glasskube/commit/170e905410d6e464355b09e9b65ff5e7e2511cc6))
* **deps:** update dependency typescript to ~5.4.0 ([004dcbe](https://github.com/glasskube/glasskube/commit/004dcbeb3cf776c57badff153eff3f00ca44d5af))
* **deps:** update google-github-actions/release-please-action action to v4.1.0 ([30871fa](https://github.com/glasskube/glasskube/commit/30871faf683ce392c78e89b955520fab971ad3fe))


### Docs

* remove disclaimer from package update section ([7cf43eb](https://github.com/glasskube/glasskube/commit/7cf43ebde2ec47316ce2a111572488079f15d368))
* **website:** adapt roadmap to match our new release convention ([ff41302](https://github.com/glasskube/glasskube/commit/ff413026b39b57e1b95bf4fbcad787aded2ca739))
* **website:** add latest release to roadmap ([319d33e](https://github.com/glasskube/glasskube/commit/319d33e9aa4a74b113e959f52f15a709a82c6ba1))
* **website:** add latest release video ([b58eb3b](https://github.com/glasskube/glasskube/commit/b58eb3bd651bb6eb60dcdcfec309a4a5caa98229))
* **website:** add shadow to screenshot ([d1b1b9d](https://github.com/glasskube/glasskube/commit/d1b1b9d1b51b0e4c27269c6930b16960436ad701))
* **website:** add v0.0.3 release blogpost ([e080966](https://github.com/glasskube/glasskube/commit/e08096676f1e49ede03a07711d00bd3e7d2c4dda))
* **website:** fix release video url ([3b8a07b](https://github.com/glasskube/glasskube/commit/3b8a07b245a4918cd916a5902aa72c1e4d3f91ae))

## [0.0.3](https://github.com/glasskube/glasskube/compare/v0.0.2...v0.0.3) (2024-02-27)


### Features

* add current context in uninstall ([#273](https://github.com/glasskube/glasskube/issues/273)) ([470e01e](https://github.com/glasskube/glasskube/commit/470e01e8ecccec8d0b16cb9e33e898de6b892f64))
* add foreground propagation for uninstall cmd ([54be740](https://github.com/glasskube/glasskube/commit/54be740aae7d435b96adb37df9624e947024019d))
* add shortName and extra cols to crds ([ab439c9](https://github.com/glasskube/glasskube/commit/ab439c947e629b2b5cb172ad395195a1e7c09079))
* add version command ([#265](https://github.com/glasskube/glasskube/issues/265)) ([c4d1bc5](https://github.com/glasskube/glasskube/commit/c4d1bc5ca6df1d5e01ba2a1d4c73543e6c501042))
* added a feature to detect outdated client for cli ([#210](https://github.com/glasskube/glasskube/issues/210)) ([5e6f1a7](https://github.com/glasskube/glasskube/commit/5e6f1a79f6f38ab8742f484560b5688ff10dab68))
* added Uninstall Blocking, Progress Spinner and enhanced CLI UI for Uninstall Command hash168 ([f044a32](https://github.com/glasskube/glasskube/commit/f044a32423cc9f0c6233dad6fa76cab8d6e4d294))
* **api:** add LocalPort, Scheme to PackageEntrypoint ([73fd8fa](https://github.com/glasskube/glasskube/commit/73fd8fa1c46ba52d9de5b7b6fae8fa0a9122e004))
* **cli:** add confim dialog in install cmd ([47af3f1](https://github.com/glasskube/glasskube/commit/47af3f1ed4757f8536839d33a9e387ca8c8a4747))
* **cli:** add handling of LocalPort, Scheme in Open cmd ([8557f12](https://github.com/glasskube/glasskube/commit/8557f127e8689193de76d7109e3904871b88b959))
* **cli:** add installing specific package version ([#203](https://github.com/glasskube/glasskube/issues/203)) ([23b2943](https://github.com/glasskube/glasskube/commit/23b2943ea4d8d3a9afd85b0faa9a5d04aa357643))
* **cli:** add outdated flag for list cmd ([#201](https://github.com/glasskube/glasskube/issues/201)) ([8d93c61](https://github.com/glasskube/glasskube/commit/8d93c619a4404f356f0dabaf4b9e136c2efd2458))
* **cli:** add showing package version in list cmd ([#200](https://github.com/glasskube/glasskube/issues/200)) ([fe856bc](https://github.com/glasskube/glasskube/commit/fe856bc30aaaee861dba7b78e5cac6647de1db08))
* **cli:** add update cmd ([#202](https://github.com/glasskube/glasskube/issues/202)) ([ddfe2cf](https://github.com/glasskube/glasskube/commit/ddfe2cf0c7ebc0e98a693f60ca1bd552f28ab47e))
* **cli:** automatic bootstrap in CLI commands ([#196](https://github.com/glasskube/glasskube/issues/196)) ([5d86eb1](https://github.com/glasskube/glasskube/commit/5d86eb1f57ae236d42fd94ac56de75623397d1cf))
* **cli:** change update message to be less obtrusive ([0b2dfd2](https://github.com/glasskube/glasskube/commit/0b2dfd22799cccee45854485a52427c37a0c04c1))
* **cli:** glasskube describe ([#241](https://github.com/glasskube/glasskube/issues/241)) ([55bebd6](https://github.com/glasskube/glasskube/commit/55bebd69c8c55ccaaa293d1b2bd74c8292759674))
* include current context in install ([5d3d390](https://github.com/glasskube/glasskube/commit/5d3d390b0c86780cfe34dbb3ab02f3ea43645cc8))
* **package-operator:** add blockOwnerDeletion on OwnerReferences ([5ec721b](https://github.com/glasskube/glasskube/commit/5ec721b17af9c150aad1c543c34cb901751d2558))
* **package-operator:** add version aware package updates ([8c56780](https://github.com/glasskube/glasskube/commit/8c56780f93618e038a96d5bd977184f5b6ffa02e))
* **package-operator:** add version aware packageinfo updates ([b6ebf5b](https://github.com/glasskube/glasskube/commit/b6ebf5b9077e8a612cf4e4fef9fd183499b2e3cc))
* **ui:** add current context in navbar ([#263](https://github.com/glasskube/glasskube/issues/263)) ([067793c](https://github.com/glasskube/glasskube/commit/067793cf07e8feb6ee50bd86deec75aa4db218c1))
* **ui:** add infinite progress bar ([8d9b291](https://github.com/glasskube/glasskube/commit/8d9b2919413fcd083ee017f0a8ab91522b89e115))
* **ui:** add package detail page [#172](https://github.com/glasskube/glasskube/issues/172) ([1298292](https://github.com/glasskube/glasskube/commit/1298292aea123276c5a3983cf4062d38f122568e))
* **ui:** add selecting kubeconfig ([#140](https://github.com/glasskube/glasskube/issues/140)) ([620f5d5](https://github.com/glasskube/glasskube/commit/620f5d59f4d2a4f268e7162fde1b8625ef929104))
* **ui:** bootstrap via UI [#123](https://github.com/glasskube/glasskube/issues/123) ([c8c0576](https://github.com/glasskube/glasskube/commit/c8c05760680ba3e71f1ee72aa2c6b444d495aef6))
* **ui:** install package in specific version ([#269](https://github.com/glasskube/glasskube/issues/269)) ([5d6067b](https://github.com/glasskube/glasskube/commit/5d6067bfe3cb80cf9d80c075ab974223520711c8))
* **ui:** update all packages [#289](https://github.com/glasskube/glasskube/issues/289) ([cdce20f](https://github.com/glasskube/glasskube/commit/cdce20f50ef2124e4dd5b95185b39aea07d1f3c9))
* **ui:** update packages and push cluster events to UI ([#269](https://github.com/glasskube/glasskube/issues/269)) ([d81d1bb](https://github.com/glasskube/glasskube/commit/d81d1bb9b56039cd379c5d85be7d7bdfab52e89c))
* **website:** add blogpost about kubernetes frontends ([f5c16e3](https://github.com/glasskube/glasskube/commit/f5c16e34bcdd3eda217a0b7496cbfcd33adf318a))


### Bug Fixes

* **action:** build package-operator image for amd64, arm64 ([38e1cc5](https://github.com/glasskube/glasskube/commit/38e1cc5a810dc2767d937cc5491f5377a6b2ea29))
* add path escaping repository urls ([8bffd3b](https://github.com/glasskube/glasskube/commit/8bffd3b79d9d9aabbdb8cd3ca1f885b01bac9190))
* change opener to use real PackageInfo name ([879e647](https://github.com/glasskube/glasskube/commit/879e647f2dda51d0a679cff30345449717c40365))
* **deps:** update dependency @mdx-js/react to v3.0.1 ([a74e4b9](https://github.com/glasskube/glasskube/commit/a74e4b9b9db0bc5c8c6a2340ba2b999b7e5ec966))
* **deps:** update dependency asciinema-player to v3.6.4 ([a1e8e69](https://github.com/glasskube/glasskube/commit/a1e8e69672557b95c0387d2f20ea8d1199fdbc9f))
* **deps:** update dependency asciinema-player to v3.7.0 ([7093e80](https://github.com/glasskube/glasskube/commit/7093e806a4ffcf0e1c7375b560d850b5639d0f1e))
* **deps:** update kubernetes packages to v0.29.2 ([5db9fcb](https://github.com/glasskube/glasskube/commit/5db9fcb2c2599bcb7fcaa8bbb47ba4551665583a))
* **deps:** update module github.com/schollz/progressbar/v3 to v3.14.2 ([1d20832](https://github.com/glasskube/glasskube/commit/1d20832fda5134dfd6da95488f1f5e638c8522e3))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.17.2 ([8a6dd83](https://github.com/glasskube/glasskube/commit/8a6dd83a3173c006330dd577891076172b0f394e))
* made the necessary changes and code cleanup ([d184784](https://github.com/glasskube/glasskube/commit/d184784e2d2d469ac9c868f2e8ea50dc4e8cdbe6))
* **package manifest:** supported svg files and change the schema accordingly ([070a7e2](https://github.com/glasskube/glasskube/commit/070a7e21424db5a7bde162534d575e0302f43b37))
* **package-operator:** set pending status when package is being deleted ([718aadc](https://github.com/glasskube/glasskube/commit/718aadc40b33495fd87ed93c22e0251e8c1595b3))
* **packagemanifest:** updated img element and schema ([b0925db](https://github.com/glasskube/glasskube/commit/b0925db2a654390fe84570cb2f3a791f91191951))
* **ui:** reload content when pages are restored from history ([9c282ce](https://github.com/glasskube/glasskube/commit/9c282ce29c5925c45dc3c14ceb8fdd0be216a2f9))
* **ui:** support page error when no context given ([62faa9b](https://github.com/glasskube/glasskube/commit/62faa9be5d5658f48420b4a1f4fe55f23f053b1e))
* **ui:** use real PackageInfo name when getting installed manifest ([13f9e60](https://github.com/glasskube/glasskube/commit/13f9e602793b9dbf54b143be535bb5b9a853b991))


### Other

* add generating schema files for repo types ([c05e5ce](https://github.com/glasskube/glasskube/commit/c05e5ced509d6332434170f615c16fc417fa12e5))
* change schema to not use $ref and move to website ([d42fcb1](https://github.com/glasskube/glasskube/commit/d42fcb11c51631a558f32dfd58d251f7b5dcb2ee))
* cleanup uninstaller creation ([cca98c6](https://github.com/glasskube/glasskube/commit/cca98c647333a723a82d9110100dbec5a2780bc6))
* **deps:** update go.sum file ([be90842](https://github.com/glasskube/glasskube/commit/be908428ef2b5939327b42dfe5bbe7a080f9f434))
* make renovate run "go mod tidy" ([88e93aa](https://github.com/glasskube/glasskube/commit/88e93aa89f87935aeffccfed941fcf92ad87e050))
* remove path escape workaround ([424086f](https://github.com/glasskube/glasskube/commit/424086f736c0f42d78894e4f252a7b2d145f0617))
* solved errors for install ([b947042](https://github.com/glasskube/glasskube/commit/b947042382c01f72085f7892f728a60cea1895fb))
* update controller-tools to v0.14.0 ([e9573e6](https://github.com/glasskube/glasskube/commit/e9573e653e24843c1014288655ad0386795cab89))


### Docs

* add docs for outdated flag for list cmd ([#201](https://github.com/glasskube/glasskube/issues/201)) ([e96c769](https://github.com/glasskube/glasskube/commit/e96c769d3268f81307840bf89856e9f4d79fe19c))
* add flags to install cmd docs ([ad294e4](https://github.com/glasskube/glasskube/commit/ad294e4dcac4ec08debe9e07524c5ebf2abc3d7e))
* add package reconciliation diagram ([28e9266](https://github.com/glasskube/glasskube/commit/28e92660b22c609bb94f2c244e33a2db229b156b))
* add update cmd docs ([c5ccf94](https://github.com/glasskube/glasskube/commit/c5ccf94b99cf8ef34fd016f49dd3ffdedaf5b45c))
* **webiste:** add release video to v0.0.2 blog post ([2a2a7a9](https://github.com/glasskube/glasskube/commit/2a2a7a9346e0740b89d5214655bac1c15099a5e6))
* **website:** adapt roadmap and README ([346a18f](https://github.com/glasskube/glasskube/commit/346a18f335e600a29875ae6e73d23b17913dedf8))
* **website:** add missing descriptions to blog posts ([53a4efc](https://github.com/glasskube/glasskube/commit/53a4efca1f8b0bcac29c6b0bea4fcca780937aae))
* **website:** add v0.0.2 release blogpost ([471e26a](https://github.com/glasskube/glasskube/commit/471e26a0bcb2f9eed41c64c548807fa6cc24ca6a))
* **website:** fix typo in footer ([aaafe03](https://github.com/glasskube/glasskube/commit/aaafe03086e462c0f88249222d38d668742d29e0))
* **website:** fix typos ([3e7236a](https://github.com/glasskube/glasskube/commit/3e7236a5039f49ca4e336d9639966a4ebc8235ca))
* **website:** improve comparision documentation ([0bbc986](https://github.com/glasskube/glasskube/commit/0bbc9864d31d3c9a8e5d92997ce9ab55828f1163))
* **website:** make newsletter subscription email required ([0b80421](https://github.com/glasskube/glasskube/commit/0b80421991fa1bbdd80488c7bf978e197ced5d5a))


### Refactoring

* **ui:** add deferred config loading, use request context ([#140](https://github.com/glasskube/glasskube/issues/140)) ([1fe219c](https://github.com/glasskube/glasskube/commit/1fe219c6ec05a754578193c9e3bcc93eb41725e4))

## [0.0.2](https://github.com/glasskube/glasskube/compare/v0.0.1...v0.0.2) (2024-02-09)


### Features

* add 'port' flag to 'glasskube serve' ([2ac1485](https://github.com/glasskube/glasskube/commit/2ac14851f2451e2d7e43614b8543d72b77481435))
* add event recording ([#139](https://github.com/glasskube/glasskube/issues/139)) ([3ccc35c](https://github.com/glasskube/glasskube/commit/3ccc35cedcd9712b02118c949c26d9253e5a10e3))
* **cli, ui:** glasskube open ([8af0072](https://github.com/glasskube/glasskube/commit/8af00727e59cf8458fbbb1df62f392cd8a5602de))
* **cli:** add open for packages ([d991c5e](https://github.com/glasskube/glasskube/commit/d991c5ebb0bd9e21acb5403787fb88319beeedf0))
* **cli:** add progress spinner for install cmd ([6a99408](https://github.com/glasskube/glasskube/commit/6a994087d673fc935278b0f0ea9c395a76261c42))
* **ui:** add open button on list page ([3e62c4b](https://github.com/glasskube/glasskube/commit/3e62c4bbb3ba9ded46d2c22b79c1ae95abf1f7d1))
* **ui:** add real-time updates [#164](https://github.com/glasskube/glasskube/issues/164) and refactor [#126](https://github.com/glasskube/glasskube/issues/126) ([11d13d2](https://github.com/glasskube/glasskube/commit/11d13d274b134c566c17ced926c146e95db72717))
* **website:** add a latest release version json file for outdated dedection ([d9b4c8a](https://github.com/glasskube/glasskube/commit/d9b4c8a5273fb63221e14ed89fa5a65e6c25a14e))
* **website:** allow search engies to index glasskube.dev ([04f6dc2](https://github.com/glasskube/glasskube/commit/04f6dc2fa0ef9671db8855e65404569b210dcfb0))
* **website:** improve glasskube vs helm meta title ([7a79ba1](https://github.com/glasskube/glasskube/commit/7a79ba16f27506f68f69dfc8ed9d4e933eb044a0))
* **website:** mark keptn as a supported package ([#137](https://github.com/glasskube/glasskube/issues/137)) ([57137b4](https://github.com/glasskube/glasskube/commit/57137b43bd33972532cca75de4ce09d92d4f76bf))


### Bug Fixes

* **cli:** add port check before forwarding ([68b273d](https://github.com/glasskube/glasskube/commit/68b273d82ae95cd17ad57ecbbf8be87ac6533073))
* **deps:** update module github.com/fluxcd/helm-controller/api to v0.37.3 ([97950df](https://github.com/glasskube/glasskube/commit/97950df2d0f455342ca7d564f51f526c48bb9be0))
* **deps:** update module github.com/fluxcd/helm-controller/api to v0.37.4 ([f734c3c](https://github.com/glasskube/glasskube/commit/f734c3c0b0139f768ba583fb5551a84e1703258e))
* **deps:** update module github.com/fluxcd/source-controller/api to v1.2.4 ([bc81723](https://github.com/glasskube/glasskube/commit/bc8172342239b9a2d71a895a7cd4f659e98f2579))
* **deps:** update module github.com/gorilla/websocket to v1.5.1 ([2ab5393](https://github.com/glasskube/glasskube/commit/2ab539323bea9fef3b269b93ba90075a9fee1a4f))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.17.1 ([6120f45](https://github.com/glasskube/glasskube/commit/6120f45d71b89b8058b60ff251475eed891d88f4))
* fix broken links ([#103](https://github.com/glasskube/glasskube/issues/103)) ([3dd8ac6](https://github.com/glasskube/glasskube/commit/3dd8ac6fd7032007d6bd634d31fffd6d60574778))
* fixed typo in mariadb-operator reference ([ffc8edb](https://github.com/glasskube/glasskube/commit/ffc8edb3bd04ad01fa5614476c1d701da8a794e9))
* folder name of technical preview  blogpost ([04ec42b](https://github.com/glasskube/glasskube/commit/04ec42bd76ab7cec322bcfd90ec651e032af2788))
* **package-operator:** add version in manifest ([ed4a657](https://github.com/glasskube/glasskube/commit/ed4a657770b3a6196401e62d25095d7ad5b2510d))
* port conflict ([0f7fbcf](https://github.com/glasskube/glasskube/commit/0f7fbcf3f5cb7a2c6686a3835ee68cb7789c1f31))
* removed logger dependency issue from the repo client ([8fd9bff](https://github.com/glasskube/glasskube/commit/8fd9bff13759136e1a94dd7fb4fd2734ea8487cb))
* removed unused ctx and removed null logger ([233ab22](https://github.com/glasskube/glasskube/commit/233ab221b9ba090a44ff19f24b4289551a502700))
* **website:** add trailing slash, don't index tag pages on google ([76db2d1](https://github.com/glasskube/glasskube/commit/76db2d1f545918a1705db698ea48e308ff7d4b44))


### Other

* add codeowners ([16e2b84](https://github.com/glasskube/glasskube/commit/16e2b8451d0acd1547360cba086a062feb8476bf))
* add issue templates ([fe52235](https://github.com/glasskube/glasskube/commit/fe5223577f0517d430ce3ceb60312c4bbcbde491))
* **deps:** update actions/setup-node action to v4.0.2 ([8f47904](https://github.com/glasskube/glasskube/commit/8f479044b80c532c59f8219ec66e1e8c3d1f2b40))
* **deps:** update actions/setup-node digest to 60edb5d ([31f3abd](https://github.com/glasskube/glasskube/commit/31f3abdafe48feab75fe8d4970f049b9b8258520))
* **deps:** update aws-actions/configure-aws-credentials digest to e3dd6a4 ([56b8a84](https://github.com/glasskube/glasskube/commit/56b8a8471e15dcb1ce0c2cc41be6a953a1c29596))
* fix wrong path in release please config ([e59e4b0](https://github.com/glasskube/glasskube/commit/e59e4b004fbc5a28a48db3de008ac516246a49fa))


### Docs

* add CNCF landscape badge ([f3f46be](https://github.com/glasskube/glasskube/commit/f3f46beda7a8f81d35cd0b2b699ef3106fd5414c))
* add shell completion to install guide ([ef3bf90](https://github.com/glasskube/glasskube/commit/ef3bf909f4407ca141c5171b2747758a7a7b4f2b))
* add technical preview blogpost ([2af5351](https://github.com/glasskube/glasskube/commit/2af5351c266ff30b88e6a45cc1b0f4a3dd3e159a))
* change release-please annotations to mdx comments ([aed9f08](https://github.com/glasskube/glasskube/commit/aed9f084e2e420f418f930d56247d5136df69796))
* fix sidebar ([a223e89](https://github.com/glasskube/glasskube/commit/a223e8904d7dbf0255c9e8ae93e91e2c4daff1c2))
* fix some typos ([#178](https://github.com/glasskube/glasskube/issues/178)) ([f9ee2e2](https://github.com/glasskube/glasskube/commit/f9ee2e2af21a1c2af10946e0c893ac3bcfb02b4e))
* fix typo ([ea6bed8](https://github.com/glasskube/glasskube/commit/ea6bed838f8fb207325a3a2f83b546920c42b9f5))
* fix typo in Readme ([8bd8aa7](https://github.com/glasskube/glasskube/commit/8bd8aa7d6999ba0477d5de76fc8d6ba3b0a6f028))
* fix wrong package version ([bdf8074](https://github.com/glasskube/glasskube/commit/bdf807404d49ef9a86e70353baa4d305158b45ed))
* improve technical preview blog post ([7c409da](https://github.com/glasskube/glasskube/commit/7c409da9a2779f637f799434593d6eb45647ca9d))
* replace gif in readme with svg ([b197138](https://github.com/glasskube/glasskube/commit/b197138af062aee566382dfa67e80e6a13666510))
* update client commands ([e5b4df1](https://github.com/glasskube/glasskube/commit/e5b4df12dbff1c6a8a6442deafa489ca3b3f3218))
* update helm comparison ([150d319](https://github.com/glasskube/glasskube/commit/150d3194f1d9469cb5fa1bf3277fdd28730a8e6e))
* update roadmap ([730cb73](https://github.com/glasskube/glasskube/commit/730cb73495ed8ebcce9d2f6a32d1209f63d62ed0))
* use a custom domain for binary downloads ([cbfc1ef](https://github.com/glasskube/glasskube/commit/cbfc1efbbd6f9bee0421edf5844401021fef4120))
* **website:** add [@kubesimplify](https://github.com/kubesimplify) video to release blog post ([18a9fa4](https://github.com/glasskube/glasskube/commit/18a9fa4365af43456cd7361f0a0c1f30ba250397))
* **website:** add CTA at the bottom ([dc90ece](https://github.com/glasskube/glasskube/commit/dc90ece6c177be4322522383141cf9298afe9ada))
* **website:** add glasskube is part of the cncf landscape blog post ([607c158](https://github.com/glasskube/glasskube/commit/607c1589faa6a148480c36791abc404673d4ab7f))
* **website:** update roadmap ([297a172](https://github.com/glasskube/glasskube/commit/297a17266ed74c68d1a236829f6575aaed44f958))


### Refactoring

* **website:** move guides to own folder for a clean url ([b909508](https://github.com/glasskube/glasskube/commit/b90950823dc4d30d20c5c0502cb61aff3fc7dae7))

## 0.0.1 (2024-01-31)


### Features

* add initial project structure ([bdf7434](https://github.com/glasskube/glasskube/commit/bdf7434387eb701d279b887a458f6b857de1b974))
* **cli, client:** install package [#23](https://github.com/glasskube/glasskube/issues/23) [#28](https://github.com/glasskube/glasskube/issues/28) [#35](https://github.com/glasskube/glasskube/issues/35) ([54a5237](https://github.com/glasskube/glasskube/commit/54a5237ffb2088421c6194b1cbac9890b3db8097))
* **cli, client:** list packages [#22](https://github.com/glasskube/glasskube/issues/22) [#26](https://github.com/glasskube/glasskube/issues/26) [#34](https://github.com/glasskube/glasskube/issues/34) [#41](https://github.com/glasskube/glasskube/issues/41) ([ea4cedc](https://github.com/glasskube/glasskube/commit/ea4cedc9b65959a5746b8333d8f9eb32565a008d))
* **cli, client:** uninstall package [#30](https://github.com/glasskube/glasskube/issues/30) [#32](https://github.com/glasskube/glasskube/issues/32) ([13ccf96](https://github.com/glasskube/glasskube/commit/13ccf96a28c920e31e7182bff1afc127e24904ee))
* **cli, ui:** add serve command [#25](https://github.com/glasskube/glasskube/issues/25) ([50a384c](https://github.com/glasskube/glasskube/commit/50a384c24663c55b5c7f975b76d5f3a2026b670c))
* **cli, ui:** validate bootstrap [#149](https://github.com/glasskube/glasskube/issues/149) ([46ee71f](https://github.com/glasskube/glasskube/commit/46ee71fec8e791e5b558c214906f7796e2f0ee0e))
* **cli:** add completion for install cmd ([fe7ecec](https://github.com/glasskube/glasskube/commit/fe7ecec9fc0c9b25bbbc87ee3d057e41a1abe77b))
* **cli:** add helpful message if kubeconfig is empty ([#31](https://github.com/glasskube/glasskube/issues/31)) ([d221f89](https://github.com/glasskube/glasskube/commit/d221f89200917205e2647642ccd0c9ccb5537789))
* **cli:** initial bootstrap command ([6073f86](https://github.com/glasskube/glasskube/commit/6073f8653d2c0e381231df7981253f05e9ebaa47))
* **package-operator:** add aio config that includes flux dependencies ([9fdd43f](https://github.com/glasskube/glasskube/commit/9fdd43f3045fa3ac44f05cb733a22241d2ea58d2))
* **package-operator:** add basic fields to package crd ([#10](https://github.com/glasskube/glasskube/issues/10)) ([f192ef2](https://github.com/glasskube/glasskube/commit/f192ef2984bff96e84ebaa1ead3e5c072703f797))
* **package-operator:** add error handling in manifest adapter initialization ([#14](https://github.com/glasskube/glasskube/issues/14)) ([5089bda](https://github.com/glasskube/glasskube/commit/5089bda57d947f6f7c3da5960e68add31f31a4fd))
* **package-operator:** add error message to packageinfo condition ([#13](https://github.com/glasskube/glasskube/issues/13)) ([aa2971a](https://github.com/glasskube/glasskube/commit/aa2971a19149c63f6f9576025fb61dd94a9b20eb))
* **package-operator:** add fields to package info crd and manifest schema ([5b78e1f](https://github.com/glasskube/glasskube/commit/5b78e1f90c6c8de6d6b144864c6088d9f501340e))
* **package-operator:** add handling of package manifest and helm adapter ([#14](https://github.com/glasskube/glasskube/issues/14)) ([7b75591](https://github.com/glasskube/glasskube/commit/7b755918f81ff5b0eb833f0222ac91f69d37a45d))
* **package-operator:** add HelmRelease creation ([#14](https://github.com/glasskube/glasskube/issues/14)) ([f0ef0e2](https://github.com/glasskube/glasskube/commit/f0ef0e261a53445cf9442eb545d7521ff34c5795))
* **package-operator:** add initial packageinfo crd ([#9](https://github.com/glasskube/glasskube/issues/9)) ([1798d9d](https://github.com/glasskube/glasskube/commit/1798d9d5d4fe35f442d0ebfbd8be7f0ee116bdfb))
* **package-operator:** add missing namespace in aio config ([e3276db](https://github.com/glasskube/glasskube/commit/e3276db80672435fd5ba52890a931a7ceb98e6e7))
* **package-operator:** add multi owner references ([#12](https://github.com/glasskube/glasskube/issues/12)) ([5d8862d](https://github.com/glasskube/glasskube/commit/5d8862de8c0bfa8164feae3685c0f5d2ad15535a))
* **package-operator:** add package controller creates dependent package info ([#12](https://github.com/glasskube/glasskube/issues/12)) ([23c87be](https://github.com/glasskube/glasskube/commit/23c87be86951aa0dd0459b1bf3113b7ab3176556))
* **package-operator:** add package info controller fetches manifest ([#13](https://github.com/glasskube/glasskube/issues/13)) ([1022583](https://github.com/glasskube/glasskube/commit/102258334fb4a9dae390161ea61e7841561f0f27))
* **package-operator:** add support for plain manifests ([#88](https://github.com/glasskube/glasskube/issues/88)) ([334ce23](https://github.com/glasskube/glasskube/commit/334ce2366e8c736199a68c57105c67b2849be038))
* **package-operator:** add type flag for bootstrap cmd ([197b62b](https://github.com/glasskube/glasskube/commit/197b62b4875642d222ad9ff0c2a34640d49cf9ed))
* **package-operator:** change package to be cluster scoped ([a7c9797](https://github.com/glasskube/glasskube/commit/a7c97971d425f8ba6dfbc2f1d4ddf246cfbe3184))
* **package-operator:** mark some crd fields as required ([3b2785c](https://github.com/glasskube/glasskube/commit/3b2785cb43a7277e86922595e9d51b30df9fd869))
* **ui:** install and uninstall trivial packages [#29](https://github.com/glasskube/glasskube/issues/29) [#33](https://github.com/glasskube/glasskube/issues/33) ([975f2fb](https://github.com/glasskube/glasskube/commit/975f2fbfc46f78b840ab66a5a7f941a119c5253c))
* **ui:** show packages [#27](https://github.com/glasskube/glasskube/issues/27) [#42](https://github.com/glasskube/glasskube/issues/42) ([90bb333](https://github.com/glasskube/glasskube/commit/90bb333583014381e10869f05a497e5b8ee3f9ed))
* **ui:** support user in setting up kubeconfig [#31](https://github.com/glasskube/glasskube/issues/31) ([d8e8382](https://github.com/glasskube/glasskube/commit/d8e8382e19c7fcfd9f8fe9bb0284d6dc301938bf))
* **webiste:** add cert-manager guide ([da6d987](https://github.com/glasskube/glasskube/commit/da6d9870b2a8be5a40ed7dad794c0aff6f7192aa))
* **webiste:** improve cert-manager guide ([49a98bc](https://github.com/glasskube/glasskube/commit/49a98bcbcbd9220b5be2c4384f36cda140abf0e0))
* **website:** add announcement bar ([#56](https://github.com/glasskube/glasskube/issues/56)) ([5c6ee1c](https://github.com/glasskube/glasskube/commit/5c6ee1ca41cbb6539232029aace3e10b951b6cef))
* **website:** add imprint & data privacy page ([fef7d7c](https://github.com/glasskube/glasskube/commit/fef7d7ccfb0ffd0aeabc73224b1b7fad5d8ffbc4))
* **website:** add initial website version ([d03b807](https://github.com/glasskube/glasskube/commit/d03b80758a074becdb9ae1842cbab86c77436ecf))
* **website:** add missing features ([#56](https://github.com/glasskube/glasskube/issues/56)) ([6f9faae](https://github.com/glasskube/glasskube/commit/6f9faaee166a1b55a43fb244e0c5e1899f532259))
* **website:** add newsletter signup component ([#56](https://github.com/glasskube/glasskube/issues/56)) ([f8b1f68](https://github.com/glasskube/glasskube/commit/f8b1f688fde8d0a9d67dc16822e22990dd6ccc9a))
* **website:** add package overview ([44854c1](https://github.com/glasskube/glasskube/commit/44854c188b3d8d326c5156e410f755f2e14c0ddf))
* **website:** add roadmap ([f89be7d](https://github.com/glasskube/glasskube/commit/f89be7d7e0a8f92ddde0ed979ff16c44fe6b5762))
* **website:** improve hero section ([#56](https://github.com/glasskube/glasskube/issues/56)) ([4ff4b5f](https://github.com/glasskube/glasskube/commit/4ff4b5f31f39d8f73329875d5d4c807edbc71f89))
* **website:** improve seo preview for homepage ([205fa1d](https://github.com/glasskube/glasskube/commit/205fa1dfe81a7519b3d647a11c3b410a1692c9ac))
* **website:** migrate asciiflow chart to mermaid ([618737c](https://github.com/glasskube/glasskube/commit/618737c25f9a7f788f42680cbe9c12e6072aab1d))
* **website:** use GitHub meta image for website ([600b31e](https://github.com/glasskube/glasskube/commit/600b31e081524d3d5c91a1ec0b7ad1992057c02b))


### Bug Fixes

* **deps:** update docusaurus monorepo to v3.1.1 ([54d0cd9](https://github.com/glasskube/glasskube/commit/54d0cd903b2284e2b01468a298912f2a6766fa1f))
* **deps:** update kubernetes packages to v0.29.1 ([73b9412](https://github.com/glasskube/glasskube/commit/73b94127a6e0a2ecd6d888a3b597bca6ba97de03))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.15.0 ([fcf5deb](https://github.com/glasskube/glasskube/commit/fcf5deb897189faf96cb6f3fe20750ea8d632076))
* **deps:** update module github.com/onsi/gomega to v1.31.0 ([221a05a](https://github.com/glasskube/glasskube/commit/221a05aabf93f30b7e328882bd02e1d204ae31d3))
* **deps:** update module github.com/onsi/gomega to v1.31.1 ([5c598ae](https://github.com/glasskube/glasskube/commit/5c598aed5399940e7ca0cd9ac3f640f440e0f046))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.17.0 ([26e0daf](https://github.com/glasskube/glasskube/commit/26e0dafe0d8796d0b136ed6895b1637c8ae3d9e0))
* **helm-flux-adapter:** avoid panic when helm values are empty ([d9af61b](https://github.com/glasskube/glasskube/commit/d9af61bda7379ee847459a2c0d66eed89bb20861))
* **package-operator:** add json un-/marshalling of helm values ([#14](https://github.com/glasskube/glasskube/issues/14)) ([ed4a4a5](https://github.com/glasskube/glasskube/commit/ed4a4a5b80376747978a188549ea67c2e8e5ef29))
* **package-operator:** add missing rbac permissions ([d87510e](https://github.com/glasskube/glasskube/commit/d87510e6bade856f0cc3d455172d2a2b050724af))
* **package-operator:** change command controller container ([7c575e8](https://github.com/glasskube/glasskube/commit/7c575e819b9a2691bfb412c928c4c3725b9f0efa))
* **package-operator:** change image name in controller manifest ([8083831](https://github.com/glasskube/glasskube/commit/808383191d547fdf8d0326d9a033a63504cd65d0))
* **package-operator:** dont set package status failed for specific flux failure ([6a5a8d6](https://github.com/glasskube/glasskube/commit/6a5a8d6d1f33c5ab9f38388799f5ad8e735fdfa1))
* **package-operator:** make requireBootstrap optional during init ([04364cc](https://github.com/glasskube/glasskube/commit/04364cc6d81ec99e494f834730b6c1bb7e98798d))
* **package-operator:** reconcile package after helmrelease changes ([3f2949a](https://github.com/glasskube/glasskube/commit/3f2949a002ddafe180f96ec2a7daf2123bebb4f8))
* **package-operator:** remove cluster.local dns suffix from flux manifests ([0a4e748](https://github.com/glasskube/glasskube/commit/0a4e748653613d7d9ce442062d3eab619d4b60e7))
* **package-operator:** remove unused arguments ([71736ab](https://github.com/glasskube/glasskube/commit/71736abed7403ba23a0fefd976f74f8b2bea4ac0))
* **repo-client:** unmarshall helm values ([d8e6063](https://github.com/glasskube/glasskube/commit/d8e60635e1792aba642c6e2d21eb93db34d83d4f))
* **website:** editUrl ([7c81ce7](https://github.com/glasskube/glasskube/commit/7c81ce7d3efb362088d4d603ad9097194853935d))
* **website:** fix header in dark mode, enable matomo tracking ([cbf18a8](https://github.com/glasskube/glasskube/commit/cbf18a89567c162599915d43f0f0e57fa8bd0bb2))
* **website:** improve general layout responsiveness ([#94](https://github.com/glasskube/glasskube/issues/94)) ([699741c](https://github.com/glasskube/glasskube/commit/699741cabbeb7a2e203ddbda39c71cd50fb981be))


### Other

* add docs files to release please extra files ([ba17806](https://github.com/glasskube/glasskube/commit/ba17806e37ed1d72ffba86c25a0290851d83552c))
* add initial release-please version ([980eb93](https://github.com/glasskube/glasskube/commit/980eb9377fa64f40e27c5522a81804f64e1357ca))
* add intellij config file to .gitignore ([29a0776](https://github.com/glasskube/glasskube/commit/29a0776c01a339b9df852eec0fc9a83bef7d4e26))
* add rc to pre-release detection ([b9bca48](https://github.com/glasskube/glasskube/commit/b9bca48536d0b1f65fe1e7583170f77d93e5146f))
* add release please configuration ([953cf47](https://github.com/glasskube/glasskube/commit/953cf47731de8047b2619a63a53bc6f780c85349))
* Configure renovate ([79a7cb2](https://github.com/glasskube/glasskube/commit/79a7cb2a3d13fe80339dadb3c5ab1d2c9f2c4ca9))
* **deps:** update actions/setup-go digest to 0c52d54 ([a14f614](https://github.com/glasskube/glasskube/commit/a14f6146d8b439f37b804666ca844f4663f71090))
* **deps:** update actions/upload-artifact digest to 26f96df ([20d1aac](https://github.com/glasskube/glasskube/commit/20d1aac99dca945f5c494f1533221b2254437611))
* **deps:** update actions/upload-artifact digest to 694cdab ([72e8a84](https://github.com/glasskube/glasskube/commit/72e8a84e9909b4970f65e673028905c6bc37503e))
* **deps:** update dependency typescript to ~5.3.0 ([acf22bd](https://github.com/glasskube/glasskube/commit/acf22bd3b472c3dbb3c837d1fc936f1eca2d381a))
* **deps:** update golang docker tag to v1.21 ([cfe0030](https://github.com/glasskube/glasskube/commit/cfe00301a3c8953e285f225415a6863fbe8080a9))
* fix current release number ([23a2241](https://github.com/glasskube/glasskube/commit/23a2241dcdfe5d38256a9e78153e3346933d47e4))
* fix typos, remove default settings ([0690063](https://github.com/glasskube/glasskube/commit/0690063ab5256bcb4e35ca0a153322134a619844))
* Initial Glasskube commit ([dd7f9ee](https://github.com/glasskube/glasskube/commit/dd7f9eec89d98f80f1d76e86bf3a6194991cf051))
* **package-operator:** add missing rbac ([5a6b316](https://github.com/glasskube/glasskube/commit/5a6b316c6c322e073bf99dd79152516a57b1244e))
* remove redundant newline ([860ddb6](https://github.com/glasskube/glasskube/commit/860ddb6b060a28728e7b245dc4c71dfa7735f3f5))
* **website:** fix typos, improve typwriter speeed ([#56](https://github.com/glasskube/glasskube/issues/56)) ([9c28b4c](https://github.com/glasskube/glasskube/commit/9c28b4c377e7786da9da5470ac91fafb9e6e4d8a))


### Docs

* add architecture diagram to readme ([5e23d28](https://github.com/glasskube/glasskube/commit/5e23d28e8d7caabf6ff7713c3d139a9ffcb9e9df))
* add asciinema cast to website ([620d7ea](https://github.com/glasskube/glasskube/commit/620d7eaa360e56224330cfbd23f8174e46c7f2c1))
* add bootstrap guide ([f4965eb](https://github.com/glasskube/glasskube/commit/f4965ebfe9c7e3f619f52848e9e63ffd99e82abe))
* add bootstrap info in install segment ([622d862](https://github.com/glasskube/glasskube/commit/622d8627d6af0bd0f903c7d1cd0e54ed1b71c915))
* add brew installation instructions ([d1b575c](https://github.com/glasskube/glasskube/commit/d1b575c972b94cd569a61c83e5cfc10873ed4521))
* add gui mockup image ([cb48b08](https://github.com/glasskube/glasskube/commit/cb48b08ac564e86ea5b51f5f215c80ad2240cb37))
* add helm comparison ([bd6579a](https://github.com/glasskube/glasskube/commit/bd6579a0a4c9fef9148fd86576d74c8a0799353a))
* add ingress-nginx guide ([05ab4cf](https://github.com/glasskube/glasskube/commit/05ab4cf23aa0cafe6881deb0164207e702a8b8ff))
* add installation guide ([c76d5e0](https://github.com/glasskube/glasskube/commit/c76d5e0fcc043e05861e37d39f71928e5bd734eb))
* add note about future version bootstrap ([8d27be4](https://github.com/glasskube/glasskube/commit/8d27be4d40017e226cb3230be860b8a7af4a486a))
* change asciinema link to gif ([1752b15](https://github.com/glasskube/glasskube/commit/1752b1502c89c30d29b90301db1ecd848702733e))
* create beautiful readme ([#7](https://github.com/glasskube/glasskube/issues/7)) ([cf099f0](https://github.com/glasskube/glasskube/commit/cf099f08ad2415e8f709aa217be1b3bb3e2a10b5))
* create community standards ([#7](https://github.com/glasskube/glasskube/issues/7)) ([aa75153](https://github.com/glasskube/glasskube/commit/aa75153ebcabf71710fdc897a59fd6fb4e3b678a))
* extract guides as separate page ([92bfb28](https://github.com/glasskube/glasskube/commit/92bfb28d0b52a6feca4125563b0cf77c0d43f799))
* fix DevOps / GitOps confusion ([2c7fe8d](https://github.com/glasskube/glasskube/commit/2c7fe8d76c9f274ba21e439e7716dc60031ad4bd))
* fix downloads shield, make sure all shields use flat styling ([fa82ccc](https://github.com/glasskube/glasskube/commit/fa82ccc4553c8a61e173f626c1aa1ad1db9c5205))
* fix typos, reformat code ([1c5982e](https://github.com/glasskube/glasskube/commit/1c5982ea6d861e3170b96d645060398c5597100c))
* remove release please comments from readme ([c46df8a](https://github.com/glasskube/glasskube/commit/c46df8a93d7b254d0cd681ca12d9bb88bb6105a1))
* replace link on homepage with asciinema-player ([a60670e](https://github.com/glasskube/glasskube/commit/a60670ef01c8e583774998de480c21e338ec0e65))
* update component descriptions ([af24d3e](https://github.com/glasskube/glasskube/commit/af24d3eee6e6cec1753720c9df5dd4652781e403))
* update quick start section in readme ([6072cf2](https://github.com/glasskube/glasskube/commit/6072cf2c2a016bab551545baae98dad149df855d))
* update readme and roadmap ([176d38a](https://github.com/glasskube/glasskube/commit/176d38aca90df693b1263ce6608bdd033e84ff40))
* update supported packages ([f061395](https://github.com/glasskube/glasskube/commit/f061395763ea0ba7748c1a039effcbef5c18972c))


### Refactoring

* **cli:** reuse existing code for bootstrap ([dac2f45](https://github.com/glasskube/glasskube/commit/dac2f453f6a5b968e8b9f40c7cffb823837dc38e))
* **package-operator:** move condition utils to separate package ([577f338](https://github.com/glasskube/glasskube/commit/577f338358b8b994509be4780b8f312ca494d5ee))
* **package-operator:** move repo client to own package ([c43d77e](https://github.com/glasskube/glasskube/commit/c43d77ed7a8d2df2b6642bef0db6191c8cb7b838))
* use css variables for theme colors ([e0939d1](https://github.com/glasskube/glasskube/commit/e0939d1b72d3ceae7f5787efd8236acfdb1971d9))
