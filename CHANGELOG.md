# Changelog

## [0.12.0](https://github.com/glasskube/glasskube/compare/v0.11.0...v0.12.0) (2024-07-01)


### Features

* add namespace support for enable/disable auto-update ([#834](https://github.com/glasskube/glasskube/issues/834)) ([2d0c02b](https://github.com/glasskube/glasskube/commit/2d0c02b9183ef7b97f12e7afb47b2fd8bd516dfb))
* **cli:** add --dry-run to bootstrap command ([#819](https://github.com/glasskube/glasskube/issues/819)) ([3a5b2d2](https://github.com/glasskube/glasskube/commit/3a5b2d2c147ec7b9657ed38b0d52b8926bb8340a))
* **cli:** add installing namespace-scoped packages ([#851](https://github.com/glasskube/glasskube/issues/851)) ([c7e139f](https://github.com/glasskube/glasskube/commit/c7e139ffc89cb93a76066c475b92dce0c86fcf22))
* **cli:** add support for auto-updating namespace-scoped packages ([#855](https://github.com/glasskube/glasskube/issues/855)) ([93926ad](https://github.com/glasskube/glasskube/commit/93926adfe3ca36653a00b65a88644fa8c5b73657))
* **cli:** add support for describing namespaced packages ([#877](https://github.com/glasskube/glasskube/issues/877)) ([aba9651](https://github.com/glasskube/glasskube/commit/aba965152e70a6f60544cc34d7229a35fe1698e1))
* **cli:** add support for updating namespace-scoped packages ([#859](https://github.com/glasskube/glasskube/issues/859)) ([25cc84c](https://github.com/glasskube/glasskube/commit/25cc84c76454bc9eb24b9d84a4c295f7190cc603))
* **cli:** add uninstalling namespaced packages ([#857](https://github.com/glasskube/glasskube/issues/857)) ([41b3773](https://github.com/glasskube/glasskube/commit/41b3773e88e064ce68958ee45c5d89885777204d))
* **cli:** support repo deletion with `repo delete [repoName]` ([#909](https://github.com/glasskube/glasskube/issues/909)) ([3152412](https://github.com/glasskube/glasskube/commit/31524127caf735ef5cde1d492802466a5a53e811))
* **deploy:** add autoupdate manifest ([#902](https://github.com/glasskube/glasskube/issues/902)) ([6c9e24d](https://github.com/glasskube/glasskube/commit/6c9e24d9a963a0059aad58cd92ddb0d3619683c9))
* **ui:** add support for namespace-scoped packages ([#817](https://github.com/glasskube/glasskube/issues/817)) ([ea9bbfb](https://github.com/glasskube/glasskube/commit/ea9bbfb161c4d09337d8c36a61fad6991e29798b))
* **ui:** added asterisk for mandatory form inputs ([#904](https://github.com/glasskube/glasskube/issues/904)) ([d4c959c](https://github.com/glasskube/glasskube/commit/d4c959c69d2e437eb4d15f6354a90ce329adee05))


### Bug Fixes

* **cli:** remove bootstrap requirement from `auto-update` command ([#901](https://github.com/glasskube/glasskube/issues/901)) ([6547fbb](https://github.com/glasskube/glasskube/commit/6547fbb5119cb3f328fb106e538c1eb04cad142f))
* **ui:** autoupdate label for discussion page ([#907](https://github.com/glasskube/glasskube/issues/907)) ([a4d3543](https://github.com/glasskube/glasskube/commit/a4d35432fb5cecf49ec7c1c9aa3215a75d53a23c))
* **ui:** improve navbar for smaller screens ([#876](https://github.com/glasskube/glasskube/issues/876)) ([e92cecd](https://github.com/glasskube/glasskube/commit/e92cecdb38fba125129d10109ce397a9109bf229))
* **ui:** show an error if the client cache is out of sync ([#911](https://github.com/glasskube/glasskube/issues/911)) ([9b6f854](https://github.com/glasskube/glasskube/commit/9b6f854396a229d76adb0512aa259b8a2bb8a8a2))
* **ui:** sticky navbar ([#903](https://github.com/glasskube/glasskube/issues/903)) ([4870c3a](https://github.com/glasskube/glasskube/commit/4870c3a257f43f87e4215aea6349f0e660a48d82))


### Other

* **cli:** enable `--kind`, `--namespace` flags ([#856](https://github.com/glasskube/glasskube/issues/856)) ([a7b5627](https://github.com/glasskube/glasskube/commit/a7b5627899c0ea6eea4bf2fd283f5a9bdbc33ee7))
* **deps:** update amannn/action-semantic-pull-request digest to 0723387 ([#906](https://github.com/glasskube/glasskube/issues/906)) ([09a301c](https://github.com/glasskube/glasskube/commit/09a301c7d6f3b27bb850e41261e169639cae6f86))
* **deps:** update dependency @eslint/js to v9.6.0 ([#913](https://github.com/glasskube/glasskube/issues/913)) ([032568e](https://github.com/glasskube/glasskube/commit/032568e552f3c6364f093a7930e70c712a2e95b1))
* **deps:** update dependency esbuild to v0.22.0 ([#916](https://github.com/glasskube/glasskube/issues/916)) ([ddd0e07](https://github.com/glasskube/glasskube/commit/ddd0e0715f749e8a233c6022cc3d2539cb610fdf))
* **deps:** update dependency globals to v15.7.0 ([#915](https://github.com/glasskube/glasskube/issues/915)) ([39dda0c](https://github.com/glasskube/glasskube/commit/39dda0ca7351477d2bb1a2319444afb1bc866821))
* **ui:** telemetry should exclude certain paths ([#921](https://github.com/glasskube/glasskube/issues/921)) ([05dec1a](https://github.com/glasskube/glasskube/commit/05dec1a0c53c973d5c1eb2cb322de37f77b3c6a5))


### Docs

* document purge command ([#908](https://github.com/glasskube/glasskube/issues/908)) ([51b4e2f](https://github.com/glasskube/glasskube/commit/51b4e2f45a0954c77582d6192794da9f8042f0a2))
* **website:** change YouTube embeds to use youtube-nocookie ([#888](https://github.com/glasskube/glasskube/issues/888)) ([d25231a](https://github.com/glasskube/glasskube/commit/d25231a88175f0f55378ae3714008f47d3982b6a))
* **website:** exchange AsciinemaPlayer with youtube demo video embed ([#883](https://github.com/glasskube/glasskube/issues/883)) ([52e0925](https://github.com/glasskube/glasskube/commit/52e0925acea2b14ad9e0292bf6aca307bb600bb9))

## [0.11.0](https://github.com/glasskube/glasskube/compare/v0.10.1...v0.11.0) (2024-06-27)


### Features

* **cli:** add ascii art on glasskube version ([#879](https://github.com/glasskube/glasskube/issues/879)) ([3040ab1](https://github.com/glasskube/glasskube/commit/3040ab10f156a551c1bcbbaa06c79e819460dad3))


### Bug Fixes

* **cli:** standardize usage texts ([#848](https://github.com/glasskube/glasskube/issues/848)) ([7d23c1e](https://github.com/glasskube/glasskube/commit/7d23c1e638c406e827a25388a56efd0707eeecac))
* **deps:** update module github.com/yuin/goldmark to v1.7.4 ([#868](https://github.com/glasskube/glasskube/issues/868)) ([d7ce5fa](https://github.com/glasskube/glasskube/commit/d7ce5fa72434e93720a9681279e8ad3e5e058cfe))
* **open:** fix typo in service name candidate ([#885](https://github.com/glasskube/glasskube/issues/885)) ([921d049](https://github.com/glasskube/glasskube/commit/921d049ff4f3575ee863a1f2ed3f5b78ea94bf47))


### Other

* **website:** configure eslint with docusaurus, react-ts and prettier plugins ([#858](https://github.com/glasskube/glasskube/issues/858)) ([613cbb7](https://github.com/glasskube/glasskube/commit/613cbb728da7cd1329b75b3148b17c2cb01ea50b))


### Docs

* exchange static image with gif ([#862](https://github.com/glasskube/glasskube/issues/862)) ([946baf4](https://github.com/glasskube/glasskube/commit/946baf46f4872ed2b45188dfb378ed0f2df6cb24))
* **website:** exchange repo mockup with actual screenshots ([#852](https://github.com/glasskube/glasskube/issues/852)) ([8adf8fb](https://github.com/glasskube/glasskube/commit/8adf8fb8e20f29e635eb9ce812338dd068f297bb))
* **website:** fix broken link ([#886](https://github.com/glasskube/glasskube/issues/886)) ([146dc25](https://github.com/glasskube/glasskube/commit/146dc25b11771cb81aa782fa9ec4895bccdd4a07))
* **website:** fix typo ([#878](https://github.com/glasskube/glasskube/issues/878)) ([e6ebb8c](https://github.com/glasskube/glasskube/commit/e6ebb8c16b3f41976e41b86ed5d4d130ed80fa32))
* **website:** glasskube is backed by Y Combinator ([#853](https://github.com/glasskube/glasskube/issues/853)) ([05e2ef7](https://github.com/glasskube/glasskube/commit/05e2ef7ce37af1ee31618dbb49258ab45d3a8a37))

## [0.10.1](https://github.com/glasskube/glasskube/compare/v0.10.0...v0.10.1) (2024-06-24)


### Bug Fixes

* **client:** propagate list options and apply timeout ([#843](https://github.com/glasskube/glasskube/issues/843)) ([7829778](https://github.com/glasskube/glasskube/commit/7829778a71eaef3fb32088fd0a8fbb0ac3a417a6))
* **deps:** update module github.com/yuin/goldmark to v1.7.3 ([#840](https://github.com/glasskube/glasskube/issues/840)) ([cc74656](https://github.com/glasskube/glasskube/commit/cc746566236103de0d5feb529992835f9d0fa251))
* **open:** try different service names ([#847](https://github.com/glasskube/glasskube/issues/847)) ([9a1ef27](https://github.com/glasskube/glasskube/commit/9a1ef27ce72b3ad6d6f02bcbbbc862550c63451e))

## [0.10.0](https://github.com/glasskube/glasskube/compare/v0.9.0...v0.10.0) (2024-06-21)


### ⚠ BREAKING CHANGES

* add `ClusterPackage` CRD and change `Package` CRD scope to Namespaced ([#792](https://github.com/glasskube/glasskube/issues/792))

### Features

* add `ClusterPackage` CRD and change `Package` CRD scope to Namespaced ([#792](https://github.com/glasskube/glasskube/issues/792)) ([9dd481f](https://github.com/glasskube/glasskube/commit/9dd481f5560ed725c1940a9c79ba6a30b22e6be3))
* add verifying breaking changes during bootstrap ([#824](https://github.com/glasskube/glasskube/issues/824)) ([9b53303](https://github.com/glasskube/glasskube/commit/9b53303a6d9ec135ec979c6ce1198c9416937a9c))
* **cli:** add `purge` command to remove installation from a cluster ([#783](https://github.com/glasskube/glasskube/issues/783)) ([4ebe30d](https://github.com/glasskube/glasskube/commit/4ebe30d4d1896f9fd9864ce823408b87c374d81a))
* **cli:** add `repo update` command ([#808](https://github.com/glasskube/glasskube/issues/808)) ([38719a8](https://github.com/glasskube/glasskube/commit/38719a8d2ccef94f1d0548e33a661127c45b276b))
* **cli:** bootstrap shows different prompt for bootstrapped clusters ([#822](https://github.com/glasskube/glasskube/issues/822)) ([df63fa4](https://github.com/glasskube/glasskube/commit/df63fa447775f78bb5403463b7f3de025c3c9699))


### Bug Fixes

* **cli:** `repo add --default` removes annotation for current default repo ([#827](https://github.com/glasskube/glasskube/issues/827)) ([ac27553](https://github.com/glasskube/glasskube/commit/ac2755315ff958bf736fdfb1cc46036adce65f88))
* **deps:** update dependency @easyops-cn/docusaurus-search-local to v0.44.1 ([#826](https://github.com/glasskube/glasskube/issues/826)) ([10c2797](https://github.com/glasskube/glasskube/commit/10c2797c803301a83b3c5cae19481c7ccd4f6883))
* **deps:** update dependency @easyops-cn/docusaurus-search-local to v0.44.2 ([#830](https://github.com/glasskube/glasskube/issues/830)) ([b9f56ba](https://github.com/glasskube/glasskube/commit/b9f56ba905c039621cf2267edd4c579f1e28987b))
* **deps:** update dependency asciinema-player to v3.8.0 ([#815](https://github.com/glasskube/glasskube/issues/815)) ([2b04852](https://github.com/glasskube/glasskube/commit/2b04852c09eb7da6e6dc4971eee668fbaeca27d9))
* **deps:** update module github.com/fluxcd/helm-controller/api to v1 ([#622](https://github.com/glasskube/glasskube/issues/622)) ([01dca18](https://github.com/glasskube/glasskube/commit/01dca1844a2abf03fb5eac273f0afa52fd20624d))
* **deps:** update module github.com/fluxcd/source-controller/api to v1.3.0 ([#472](https://github.com/glasskube/glasskube/issues/472)) ([4ad5b84](https://github.com/glasskube/glasskube/commit/4ad5b8424f5a0193bc0f4ddbaa31bc5ff313ef67))
* **deps:** update module github.com/google/go-containerregistry to v0.19.2 ([#814](https://github.com/glasskube/glasskube/issues/814)) ([3a11a56](https://github.com/glasskube/glasskube/commit/3a11a56c1333aef406c9258e75864ee7fe668263))
* **deps:** update module github.com/spf13/cobra to v1.8.1 ([#812](https://github.com/glasskube/glasskube/issues/812)) ([38392e2](https://github.com/glasskube/glasskube/commit/38392e2c1d1afc93112d68573ca15a7795494214))
* **deps:** update module github.com/yuin/goldmark to v1.7.2 ([#811](https://github.com/glasskube/glasskube/issues/811)) ([b74546d](https://github.com/glasskube/glasskube/commit/b74546d7637b6a257a34945f9aff51d5a6be8602))
* **deps:** update module k8s.io/klog/v2 to v2.130.0 ([#816](https://github.com/glasskube/glasskube/issues/816)) ([d9095d5](https://github.com/glasskube/glasskube/commit/d9095d57058120fccc414235f219ab7c5e13a345))
* **deps:** update module k8s.io/klog/v2 to v2.130.1 ([#831](https://github.com/glasskube/glasskube/issues/831)) ([93eea22](https://github.com/glasskube/glasskube/commit/93eea223ce88d0a71f690043c72f94b8b7f3969f))
* temporarily disable considering packages in dependency manager ([#839](https://github.com/glasskube/glasskube/issues/839)) ([a03a08e](https://github.com/glasskube/glasskube/commit/a03a08e65a7e22106cb7aadaff1e425ee439a062))
* **ui:** open package description links in new tab ([#837](https://github.com/glasskube/glasskube/issues/837)) ([4689c1d](https://github.com/glasskube/glasskube/commit/4689c1d9a18b801f441d8fa2907d40c6b0d15847))


### Other

* **deps:** update actions/checkout digest to 692973e ([#809](https://github.com/glasskube/glasskube/issues/809)) ([db43364](https://github.com/glasskube/glasskube/commit/db43364c62fd876130e6ecedef84f5699f3bca7c))
* **deps:** update dependency typescript to ~5.5.0 ([#835](https://github.com/glasskube/glasskube/issues/835)) ([b716b35](https://github.com/glasskube/glasskube/commit/b716b35e7bfa1b967064f9822152b7fe2268f5f3))
* **deps:** update website dependency ws to v8.17.1 ([#828](https://github.com/glasskube/glasskube/issues/828)) ([a7c4f19](https://github.com/glasskube/glasskube/commit/a7c4f19c41d4cf5cb82d4287c25fb989b810aa7e))


### Docs

* add upgrading guide ([#825](https://github.com/glasskube/glasskube/issues/825)) ([bc11f0b](https://github.com/glasskube/glasskube/commit/bc11f0b8cf6d23fd04bb2dcb6fced8a4252ef138))
* update README.md ([#807](https://github.com/glasskube/glasskube/issues/807)) ([307a330](https://github.com/glasskube/glasskube/commit/307a33047b365e2093f29b1a61138ce465e87ea4))

## [0.9.0](https://github.com/glasskube/glasskube/compare/v0.8.0...v0.9.0) (2024-06-13)


### Features

* **cli:** add `--dry-run` support for `glasskube install` to simulate package installation ([#727](https://github.com/glasskube/glasskube/issues/727)) ([05d6b02](https://github.com/glasskube/glasskube/commit/05d6b028275e0241e040694f4158dd3928b825fb))
* **cli:** add `--output` flag for `glasskube bootstrap` ([#779](https://github.com/glasskube/glasskube/issues/779)) ([b427e0a](https://github.com/glasskube/glasskube/commit/b427e0ad7c424f3c2c74a054cd9382f3b45fab34))
* **cli:** add `--output` flag for `glasskube update` ([#669](https://github.com/glasskube/glasskube/issues/669)) ([7bd44bf](https://github.com/glasskube/glasskube/commit/7bd44bf47eccdb231670621435711edf86adc0a4))
* **cli:** add `--output` flag to `glasskube describe` ([#717](https://github.com/glasskube/glasskube/issues/717)) ([d3562df](https://github.com/glasskube/glasskube/commit/d3562df13ee9ac3e8ae1f454fda329bf6b71e399))
* **cli:** add `auto-update` and related commands ([#772](https://github.com/glasskube/glasskube/issues/772)) ([5f441aa](https://github.com/glasskube/glasskube/commit/5f441aa2c098aeeb1e995230e8a5cb62ecbcc7ca))
* **cli:** change flag name of --force to --yes for glasskube uninstall ([#760](https://github.com/glasskube/glasskube/issues/760)) ([e1adc7d](https://github.com/glasskube/glasskube/commit/e1adc7d561901b6abcd15481bedf2ae4b4937c65))
* **cli:** introduce `--skip-open` support for `glasskube serve` ([#776](https://github.com/glasskube/glasskube/issues/776)) ([d559cbd](https://github.com/glasskube/glasskube/commit/d559cbddcb3776e42b424ca7847ca9e90116c474))
* **ui:** add discord link in glasskube footer ui ([#801](https://github.com/glasskube/glasskube/issues/801)) ([f94ddb7](https://github.com/glasskube/glasskube/commit/f94ddb71ea75113ec11cc2705ebef1768c461261))
* **ui:** cache package repositories ([#763](https://github.com/glasskube/glasskube/issues/763)) ([#791](https://github.com/glasskube/glasskube/issues/791)) ([4f3bc4f](https://github.com/glasskube/glasskube/commit/4f3bc4f5464f197ee70cc8611c7e656c6f8f06b8))
* **ui:** introduce additional logging ([#770](https://github.com/glasskube/glasskube/issues/770)) ([d73f7a9](https://github.com/glasskube/glasskube/commit/d73f7a9ec4679fbe9caeb57418f2167799b95660))
* **ui:** show repository status on settings page ([f1abe91](https://github.com/glasskube/glasskube/commit/f1abe917e828abdfebbd526c5bcaad764e2fc509)), closes [#751](https://github.com/glasskube/glasskube/issues/751)


### Bug Fixes

* **cli:** set autoUpdate a boolean for gk describe yaml/json output ([#780](https://github.com/glasskube/glasskube/issues/780)) ([87ad42b](https://github.com/glasskube/glasskube/commit/87ad42b06097ebf80dbe67c2c19947ddec9888b2))
* **deps:** update dependency @easyops-cn/docusaurus-search-local to ^0.43.0 ([fe0df91](https://github.com/glasskube/glasskube/commit/fe0df91d283da035adf5ceacdd5bffc101c031e7))
* **deps:** update dependency @easyops-cn/docusaurus-search-local to ^0.44.0 ([#761](https://github.com/glasskube/glasskube/issues/761)) ([1007409](https://github.com/glasskube/glasskube/commit/1007409046c3fa526e7e56754bef0991dd14c2f6))
* **deps:** update kubernetes packages to v0.30.2 ([#805](https://github.com/glasskube/glasskube/issues/805)) ([dcc2784](https://github.com/glasskube/glasskube/commit/dcc2784bd432723fa6193089d02057146d6a582a))
* **deps:** update module github.com/schollz/progressbar/v3 to v3.14.4 ([#785](https://github.com/glasskube/glasskube/issues/785)) ([647e286](https://github.com/glasskube/glasskube/commit/647e286e5bf7f159b32a50b1b85ec103c2473bf5))
* **deps:** update module golang.org/x/term to v0.21.0 ([861d695](https://github.com/glasskube/glasskube/commit/861d6953ef10544f785af033f39f5115207b9d50))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.18.4 ([#759](https://github.com/glasskube/glasskube/issues/759)) ([429f9df](https://github.com/glasskube/glasskube/commit/429f9dfd5ec9e7343452408040dcbf3ffa3f833c))
* suppress 404 error for packages w/o a giscus discussion ([cb927d2](https://github.com/glasskube/glasskube/commit/cb927d28b2942d8f8d2514d24794981d7d15c3d4))
* **website:** avoid full page crash ([#774](https://github.com/glasskube/glasskube/issues/774)) ([9580559](https://github.com/glasskube/glasskube/commit/95805597c7cdb783b0b822b339cce140bb31ffe9))


### Other

* **deps:** update amannn/action-semantic-pull-request digest to e32d7e6 ([#798](https://github.com/glasskube/glasskube/issues/798)) ([5d21649](https://github.com/glasskube/glasskube/commit/5d216495755a438664e746ab78911c6e5d77d52e))
* **deps:** update dependency esbuild to v0.21.5 ([#786](https://github.com/glasskube/glasskube/issues/786)) ([998023e](https://github.com/glasskube/glasskube/commit/998023e3d74566c47df1757cdab5c099e5c89138))
* **deps:** update dependency go to v1.22.4 ([e775b0b](https://github.com/glasskube/glasskube/commit/e775b0be923c167b49c8555982db3595622bf006))
* **deps:** update dependency prettier to v3.3.1 ([b2043ef](https://github.com/glasskube/glasskube/commit/b2043ef4746a68f73bfbff1b54be25ac2bc8d369))
* **deps:** update dependency prettier to v3.3.2 ([739e251](https://github.com/glasskube/glasskube/commit/739e251b26afec7ea0cd2df993965ef1712d0b32))
* **deps:** update googleapis/release-please-action action to v4.1.3 ([5c4f3c9](https://github.com/glasskube/glasskube/commit/5c4f3c9ed31ec052f47ff9930f4cb8cb2ccafff9))
* **deps:** update goreleaser/goreleaser-action action to v6 ([#745](https://github.com/glasskube/glasskube/issues/745)) ([fe1d58e](https://github.com/glasskube/glasskube/commit/fe1d58e8da8cba4f4b800733bbbb136006b83364))
* **package-operator:** use WithBlockOwnerDeletion from controllerutil ([#762](https://github.com/glasskube/glasskube/issues/762)) ([ecd72ff](https://github.com/glasskube/glasskube/commit/ecd72ff82b2eda93d2c9e78139182e7fc7a261e7))


### Docs

* clarify differences between `good first issue` and `help wanted` label ([23a7d70](https://github.com/glasskube/glasskube/commit/23a7d70d62eb0df60116e0e105c08f469d98a682))
* correction in contributing.md ([#802](https://github.com/glasskube/glasskube/issues/802)) ([2fe6589](https://github.com/glasskube/glasskube/commit/2fe65894e53798334cda23c2537fc9d9f9c61594))
* fix broken pull request template url ([#784](https://github.com/glasskube/glasskube/issues/784)) ([6e3061a](https://github.com/glasskube/glasskube/commit/6e3061a1c70492b5d2968509de56cd6526dd309b))
* fix typo in README.md ([#747](https://github.com/glasskube/glasskube/issues/747)) ([679c0b3](https://github.com/glasskube/glasskube/commit/679c0b36c7f30a32f98d036029503ff1630c97c8))
* fix typos, update supported packages ([f0f9d19](https://github.com/glasskube/glasskube/commit/f0f9d19757f06c1829725a2e481cd1616c8d8eea))
* update contributing guide with updated PR workflow ([#799](https://github.com/glasskube/glasskube/issues/799)) ([7efc686](https://github.com/glasskube/glasskube/commit/7efc686cf2151caef90c5af1a72fce562908a6c1))
* update local repo section ([dae0e82](https://github.com/glasskube/glasskube/commit/dae0e82f4dba19c29c0bbb8e052f841167785a4d))
* **website:** add Hatchet and Headlamp as planned ([82784ec](https://github.com/glasskube/glasskube/commit/82784ecbe44c48ca2a27944896cb0805b4c3b247))
* **website:** prepare website for launch ([#795](https://github.com/glasskube/glasskube/issues/795)) ([7486d31](https://github.com/glasskube/glasskube/commit/7486d31b3f8b7f1dc93c6c304d7fcfda0b832482))

## [0.8.0](https://github.com/glasskube/glasskube/compare/v0.7.0...v0.8.0) (2024-06-04)


### Features

* --output support for glasskube install ([#696](https://github.com/glasskube/glasskube/issues/696)) ([f91ac9c](https://github.com/glasskube/glasskube/commit/f91ac9ce0f4577b9f86e3dec6b6b83f19197b6c3))
* add --no-progress cli flag (glasskube[#709](https://github.com/glasskube/glasskube/issues/709)) ([7592f39](https://github.com/glasskube/glasskube/commit/7592f39162dffe461ee56cb8cca5c6530d91fdb2))
* **cli:** add a "default" column to the `glasskube repo list` command ([#738](https://github.com/glasskube/glasskube/issues/738)) ([1046690](https://github.com/glasskube/glasskube/commit/1046690eb979953177c2ea5292dc1997399eb546))
* **cli:** bootstrap command will ask for user confirmation ([#719](https://github.com/glasskube/glasskube/issues/719)) ([23c988b](https://github.com/glasskube/glasskube/commit/23c988b5e1755a98733865b1b90599a7c20c6dd0))
* **ui:** add default repository indicator on settings page ([#733](https://github.com/glasskube/glasskube/issues/733)) ([#740](https://github.com/glasskube/glasskube/issues/740)) ([b97f427](https://github.com/glasskube/glasskube/commit/b97f427124be5e405b0f5d045526db9470a9664d))
* **ui:** add support for advanced options ([#716](https://github.com/glasskube/glasskube/issues/716)) ([#726](https://github.com/glasskube/glasskube/issues/726)) ([ac0ee1b](https://github.com/glasskube/glasskube/commit/ac0ee1b2d65d75898284d9ffd0add78ac94db893))
* **ui:** show reaction count on package detail page ([#207](https://github.com/glasskube/glasskube/issues/207)) ([5d02ac2](https://github.com/glasskube/glasskube/commit/5d02ac22fcc08af3d0988c88c4d220c6447957f2))


### Bug Fixes

* add checking if error is new in dependency validation ([#737](https://github.com/glasskube/glasskube/issues/737)) ([fcf21ca](https://github.com/glasskube/glasskube/commit/fcf21cae05a325f28c67cc04439214520222ba7f))
* **deps:** update dependency @easyops-cn/docusaurus-search-local to ^0.42.0 ([e53ed41](https://github.com/glasskube/glasskube/commit/e53ed41b1a7b78f38070c21d4607f92d807a120a))
* **deps:** update dependency @easyops-cn/docusaurus-search-local to v0.41.1 ([13e72b9](https://github.com/glasskube/glasskube/commit/13e72b9c487850279cffcf0d1bc6f976bb47725e))
* **deps:** update docusaurus monorepo to v3.4.0 ([bac7028](https://github.com/glasskube/glasskube/commit/bac7028b0ea5d47b89a25e3d6912e6f96040dbb9))
* **package-operator:** mark dependency as "waitingFor" if not found ([#739](https://github.com/glasskube/glasskube/issues/739)) ([c38aacb](https://github.com/glasskube/glasskube/commit/c38aacb2919d8d812d425b02364f4608ff074a3e))
* remove optimistic cache check to prevent data race ([298e8f5](https://github.com/glasskube/glasskube/commit/298e8f5832887a8154f09ef3bef0fb6ba40c6d2d))


### Other

* **deps:** update dependency prettier to v3.3.0 ([d6afad0](https://github.com/glasskube/glasskube/commit/d6afad092b08b22c8ad118a772786d3b54a745ba))


### Docs

* add Go Reference and Go Report card badges ([d7ff2d8](https://github.com/glasskube/glasskube/commit/d7ff2d8833b6aefc1817e49a13df90462dfb80f3))

## [0.7.0](https://github.com/glasskube/glasskube/compare/v0.6.0...v0.7.0) (2024-05-28)


### Features

* **cli:** add `--output` support for `glasskube configure` ([#670](https://github.com/glasskube/glasskube/issues/670)) ([9b1a82d](https://github.com/glasskube/glasskube/commit/9b1a82d46fb22bcc7e4138c6136c31237afa838b))
* **cli:** add flag to trigger non-interactive mode ([f0f54da](https://github.com/glasskube/glasskube/commit/f0f54dad95b4011c96938b5935f9b899c57be9ae))
* **cli:** show uninstall status during package removal in list and describe commands ([#654](https://github.com/glasskube/glasskube/issues/654)) ([4347003](https://github.com/glasskube/glasskube/commit/4347003176fced705875b1e064a56801e9cdfb4d))
* **ui:** add cloud signup ([b0987d1](https://github.com/glasskube/glasskube/commit/b0987d14a88213fae3384a69ed19e18253c7304a))
* **ui:** autocomplete for reference inputs ([#495](https://github.com/glasskube/glasskube/issues/495)) ([fd72d1c](https://github.com/glasskube/glasskube/commit/fd72d1c78e4f2a8a8397f3825df2b91bd902014e))
* **ui:** integrate cloud links into the layout ([1d24aaa](https://github.com/glasskube/glasskube/commit/1d24aaa32d783d2241067f68877f60054e9e28d3))


### Bug Fixes

* **deps:** update dependency @lottiefiles/react-lottie-player to v3.5.4 ([b6602e5](https://github.com/glasskube/glasskube/commit/b6602e57abbba2de5fae8a641c2ec26737fdb2f8))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.19.0 ([3c15f65](https://github.com/glasskube/glasskube/commit/3c15f6504265ba811c27a2868929bbe0c421f77b))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.18.3 ([5e28e65](https://github.com/glasskube/glasskube/commit/5e28e652307a6384b52bdddd52f8109ae90d7c1c))
* **package-operator:** add missing repoclient ([84f1bbf](https://github.com/glasskube/glasskube/commit/84f1bbf7304e3f339baa03064bff7249aaf58aa0))
* **ui:** support cloud links in bootstrap/kubeconfig pages ([688d4fe](https://github.com/glasskube/glasskube/commit/688d4febf143c03e04d64399a4a7afbc4e02b0c6))


### Other

* **deps:** update dependency esbuild to v0.21.4 ([192dc89](https://github.com/glasskube/glasskube/commit/192dc89e9b8a234230ef2594947294704b9b2c54))
* **deps:** update docker/login-action digest to 0d4c9c5 ([27acfd9](https://github.com/glasskube/glasskube/commit/27acfd97c75022a69621ea2a037cb6553e4c1230))
* fix spacing in website ([29f6e91](https://github.com/glasskube/glasskube/commit/29f6e91e2c5502ef5da6485e064d589fc7af93ee))


### Docs

* include missing changes for multi-repo support ([fe80ada](https://github.com/glasskube/glasskube/commit/fe80adad370cbb1e93d70dc4c674aa69bd538e3b))


### Refactoring

* auto-update label to annotation ([02ccfb7](https://github.com/glasskube/glasskube/commit/02ccfb7cd3a706090974a81ac0d2d57e2e7cffab))

## [0.6.0](https://github.com/glasskube/glasskube/compare/v0.5.1...v0.6.0) (2024-05-23)


### Features

* **ui:** show broken config references ([#496](https://github.com/glasskube/glasskube/issues/496)) ([7f77ff9](https://github.com/glasskube/glasskube/commit/7f77ff9b350b96ee8d02b0cfcbc0c46a2d9415d4))


### Bug Fixes

* **deps:** update dependency @easyops-cn/docusaurus-search-local to ^0.41.0 ([edca7ca](https://github.com/glasskube/glasskube/commit/edca7ca2ea75973d51a7dfe538df27c06ceb88ab))
* **deps:** update dependency @fortawesome/react-fontawesome to v0.2.2 ([195a210](https://github.com/glasskube/glasskube/commit/195a2107db6894e35e8894620e2313e96c2bdda9))
* **ui:** use correct repository on details page ([#684](https://github.com/glasskube/glasskube/issues/684)) ([86b4e8b](https://github.com/glasskube/glasskube/commit/86b4e8b6a36614a7bbc6562f02293da10f5e47bb))

## [0.5.1](https://github.com/glasskube/glasskube/compare/v0.5.0...v0.5.1) (2024-05-22)


### Bug Fixes

* create a new restmapper after applying a crd ([7c1bbf0](https://github.com/glasskube/glasskube/commit/7c1bbf01554a26f0b54a4047cb9b1b09f21f51b2))
* set correct API version to make bootstrap work ([16dce6d](https://github.com/glasskube/glasskube/commit/16dce6d45b6545843b4b7e5aa77ea9931dcfb7bf))

## [0.5.0](https://github.com/glasskube/glasskube/compare/v0.4.1...v0.5.0) (2024-05-22)


### Features

* add support for custom package repositories ([#618](https://github.com/glasskube/glasskube/issues/618)) ([cd2931d](https://github.com/glasskube/glasskube/commit/cd2931d71943eca41b39b959b8b50ef48d2eb380))
* **cli:** add `--output` option for `glasskube list` ([#638](https://github.com/glasskube/glasskube/issues/638)) ([9758cf6](https://github.com/glasskube/glasskube/commit/9758cf6dfa678e24963b432be6d152bb843cb94e))
* **ui:** show uninstalling button if a package is currently being uninstalled ([#456](https://github.com/glasskube/glasskube/issues/456)) ([af42b03](https://github.com/glasskube/glasskube/commit/af42b036ce6f646dad5bf36e2cc22097ef7ea25d))


### Bug Fixes

* **deps:** update dependency @fortawesome/react-fontawesome to v0.2.1 ([dd8a41d](https://github.com/glasskube/glasskube/commit/dd8a41dd252054ada2a9e5d8bbcce8a7fa278dcd))
* **deps:** update kubernetes packages to v0.30.1 ([18a84cb](https://github.com/glasskube/glasskube/commit/18a84cb3150249e6924eea0f09ade59bb092e56a))
* **deps:** update module github.com/go-logr/logr to v1.4.2 ([ea59508](https://github.com/glasskube/glasskube/commit/ea59508e8e1900b76466f02452337a302fd2b64c))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.18.0 ([b49f58c](https://github.com/glasskube/glasskube/commit/b49f58ce728c6e6c4694146220d69db7f79955de))
* **deps:** update module github.com/schollz/progressbar/v3 to v3.14.3 ([cc1ea16](https://github.com/glasskube/glasskube/commit/cc1ea16da9c3b8e4dd2f53466a3b60b0989039fa))


### Other

* **deps:** update actions/checkout digest to a5ac7e5 ([151f31d](https://github.com/glasskube/glasskube/commit/151f31dc3c0ad6f023c1b20eff09c128e3cce7ff))
* **deps:** update dependency esbuild to v0.21.3 ([65f4fbf](https://github.com/glasskube/glasskube/commit/65f4fbf300aba1ef7ad7396b63ea75ff4f4eecf2))


### Docs

* added dependencies installation to CONTRIBUTING.md ([56d103f](https://github.com/glasskube/glasskube/commit/56d103fa9b3ea5b36deb642ff829237688d3188a))
* update architecture diagram to conform with configurable repositories ([880cd29](https://github.com/glasskube/glasskube/commit/880cd2930c69137aa30f9097c74c013d4fae8fe9))
* **website:** add multi repo design proposal ([be653c4](https://github.com/glasskube/glasskube/commit/be653c4c76634901fb8c7e41a1e1630b52e19f58))
* **website:** added devops blogpost ([ab67565](https://github.com/glasskube/glasskube/commit/ab675655cc3173e86bb31837591820cb9c0470f6))
* **website:** fix giscus integration for guides ([da30108](https://github.com/glasskube/glasskube/commit/da30108d42431ed96845d69483b3f026f4342e66))
* **website:** update --dry-run flag section ([0229804](https://github.com/glasskube/glasskube/commit/0229804de5e70ca5641725db292c4c677f0fea7a))

## [0.4.1](https://github.com/glasskube/glasskube/compare/v0.4.0...v0.4.1) (2024-05-15)


### Bug Fixes

* add re-creating `Job` resources when bootstrapping ([#619](https://github.com/glasskube/glasskube/issues/619)) ([ce3037b](https://github.com/glasskube/glasskube/commit/ce3037bac4b44761f344076def0c5e853bed9eeb))
* **cli:** add "v" prefix to version option for `install` and `update` commands ([#609](https://github.com/glasskube/glasskube/issues/609)) ([3b296ad](https://github.com/glasskube/glasskube/commit/3b296adb474898e0f1963e6b90241f8405763a3c))
* **cli:** add proper error handling for network related errors ([#597](https://github.com/glasskube/glasskube/issues/597)) ([c381ae0](https://github.com/glasskube/glasskube/commit/c381ae095dd073260cfdd52f380b5def880a1639))
* **deps:** update module github.com/fatih/color to v1.17.0 ([c784ce3](https://github.com/glasskube/glasskube/commit/c784ce354dc7f5a24c885586a7ec05efdc8bce3f))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.17.3 ([516118f](https://github.com/glasskube/glasskube/commit/516118fee740595413fb4275eb68a3d92edb33ee))


### Other

* **deps:** update dependency esbuild to v0.21.1 ([03f3926](https://github.com/glasskube/glasskube/commit/03f39268fd51323209f4693d591c837ac4658eb6))
* **deps:** update dependency esbuild to v0.21.2 ([32af63d](https://github.com/glasskube/glasskube/commit/32af63dc45ff4a925dafc5db1561f5fe254c0b6c))
* **deps:** update dependency go to v1.22.3 ([7704e5b](https://github.com/glasskube/glasskube/commit/7704e5b8293ae52d650a315ca54abf8e95004e0a))
* **deps:** update google-github-actions/release-please-action action to v4.1.1 ([165782c](https://github.com/glasskube/glasskube/commit/165782cf860c10094c258fe970ae68b9ef43d63b))
* **deps:** update googleapis/release-please-action action to v4.1.1 ([c2d0fcf](https://github.com/glasskube/glasskube/commit/c2d0fcf4306aa9c9f5c716143d644004e6a5f646))
* **deps:** update goreleaser/goreleaser-action digest to 5742e2a ([fd73b5c](https://github.com/glasskube/glasskube/commit/fd73b5c56971e450ee3abd79bddc21c1eaa005ea))
* **ui:** remove unnecessary console.log ([159a880](https://github.com/glasskube/glasskube/commit/159a880d96fa48906c19aed1b6e4492346c7ec96))


### Docs

* fix spelling mistake in package-config.md ([8808078](https://github.com/glasskube/glasskube/commit/8808078b4db6ef12067d28778f1a2183357d4c0a))
* **website:** added rabbitmq guide + addressed review ([5982e0e](https://github.com/glasskube/glasskube/commit/5982e0e0aee95d45060dd50665d1f0515732924a))
* **website:** update packages ([4d8afc5](https://github.com/glasskube/glasskube/commit/4d8afc5cd2137491086a7e5ec074de50a9a301d7))
* **website:** update the watch vs -w section ([ed675f9](https://github.com/glasskube/glasskube/commit/ed675f9bd5e5092992ba9845d3c827a8d01d8076))

## [0.4.0](https://github.com/glasskube/glasskube/compare/v0.3.0...v0.4.0) (2024-05-07)


### Features

* added ui footer with glasskube version ([#232](https://github.com/glasskube/glasskube/issues/232)) ([c8a1836](https://github.com/glasskube/glasskube/commit/c8a18368d914a7e7879165ca19e79806bbf43384))
* **cli:** add `--version` autocomplete for `glasskube update` ([#565](https://github.com/glasskube/glasskube/issues/565)) ([3ca9bc9](https://github.com/glasskube/glasskube/commit/3ca9bc9f597726054afa168f0036bfa1747c5101))
* **ui:** integrate giscus for package feedback ([#207](https://github.com/glasskube/glasskube/issues/207)) ([c50d3a8](https://github.com/glasskube/glasskube/commit/c50d3a86321829fb641af6da144e44d748dfb2a5))


### Bug Fixes

* add validating existing install before bootstrap ([08ee8bb](https://github.com/glasskube/glasskube/commit/08ee8bb2eb7be6ad5a7075fafc58ca947f8bf96d))
* **cli:** prevent accidental bootstrapping of an older package-operator version ([abefab1](https://github.com/glasskube/glasskube/commit/abefab1d564476420e7ee93f81012412a86a578a))
* **cli:** skip version check for `bootstrap` ([419b7b4](https://github.com/glasskube/glasskube/commit/419b7b4e2ada3af469f4f8f8a5cbc4cff125aeeb))
* **deps:** update docusaurus monorepo to v3.3.2 ([f483034](https://github.com/glasskube/glasskube/commit/f483034219653f6d82fdd049afd607e37765d7a3))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.17.2 ([0b4a5f2](https://github.com/glasskube/glasskube/commit/0b4a5f2f246b7cd148cc24b7698ec1c780f2587c))
* **deps:** update module github.com/onsi/gomega to v1.33.1 ([b5cb34c](https://github.com/glasskube/glasskube/commit/b5cb34c0bd47dc29cc5105ae2bd306191e3db4aa))
* **deps:** update module golang.org/x/term to v0.20.0 ([e474086](https://github.com/glasskube/glasskube/commit/e474086d7275ed3028b1fb74adeba481e1ca5499))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.18.1 ([7946560](https://github.com/glasskube/glasskube/commit/7946560c80fc535ff575bb7fedb2246bf80dab4d))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.18.2 ([a62f92d](https://github.com/glasskube/glasskube/commit/a62f92db68fd300938b4f56124e28028ef3959a7))
* **deps:** update react monorepo to v18.3.0 ([94fa2a8](https://github.com/glasskube/glasskube/commit/94fa2a805a507e701ebb702cf2137429e4c0a01e))
* **deps:** update react monorepo to v18.3.1 ([34ddcb6](https://github.com/glasskube/glasskube/commit/34ddcb660fc3ba2deef87653471a5aa30275082a))
* **package-operator:** make order of patch application deterministic ([bfe5e59](https://github.com/glasskube/glasskube/commit/bfe5e59205d2f64d0f64ab1421d0d7065e681a07))
* **ui:** introduce cache busting ([b91acee](https://github.com/glasskube/glasskube/commit/b91acee005f0361f79f1530bdac27b06b340170c))


### Other

* **deps:** update actions/checkout digest to 0ad4b8f ([15b8f13](https://github.com/glasskube/glasskube/commit/15b8f133270c0c49a776981e74798afc9795e0a0))
* **deps:** update actions/setup-go digest to cdcb360 ([eed4598](https://github.com/glasskube/glasskube/commit/eed4598b45d8af90a2d03b0d237a68e3d42b8f3e))
* **deps:** update dependency esbuild to v0.21.0 ([d6980f9](https://github.com/glasskube/glasskube/commit/d6980f91ed087d64b87355f5dc8e62acb58c0859))
* update copyright information ([2c8b305](https://github.com/glasskube/glasskube/commit/2c8b305f7e0a481e9dbea6cf94e9a9e4fd9ec8a7))
* **website:** update package list ([3ca994b](https://github.com/glasskube/glasskube/commit/3ca994b6e55e0f8eaea00a38c2c3b701a84ebf56))


### Docs

* **website:** add giscus for blog post comments ([f11851e](https://github.com/glasskube/glasskube/commit/f11851ec6da4e442662ef4354dde21ec336734bc))
* **website:** added beta launch blogpost + cta ([18afd3a](https://github.com/glasskube/glasskube/commit/18afd3a01f7446c1d99d65afb14c36165c0de9c1))
* **website:** added kubectl blog post ([2e7ac3f](https://github.com/glasskube/glasskube/commit/2e7ac3fc6d0e8511b5c838f0f647dd40f2acdaf4))
* **website:** enable giscus on guide section ([867873e](https://github.com/glasskube/glasskube/commit/867873eb988a525dd411e965162ddbd2b6c61416))
* **website:** update guides giscus categoryId ([12fa425](https://github.com/glasskube/glasskube/commit/12fa425033edee5ef0642961a32569a5025c2880))
* **website:** updated roadmap to a more general approach ([56674e5](https://github.com/glasskube/glasskube/commit/56674e5e9271667f4af0133374b01e30b278f114))
* **website:** updated the telemetry page ([d559d08](https://github.com/glasskube/glasskube/commit/d559d08f2af196e4773ed6b95c24c020e11b16f0))


### Refactoring

* **ui:** move to SSE and fix race conditions ([493c5e7](https://github.com/glasskube/glasskube/commit/493c5e7035c5d743df09e1330e22667629c12f31))
* **website:** move kubectl blog post to guides ([0f269f6](https://github.com/glasskube/glasskube/commit/0f269f683a6520a3f24973d855cb62798f39a8cd))

## [0.3.0](https://github.com/glasskube/glasskube/compare/v0.2.1...v0.3.0) (2024-04-25)


### Features

* **cli, ui:** add markdown support in long description ([0f4891b](https://github.com/glasskube/glasskube/commit/0f4891b2cd4601fb03303f2974aa84b79e53377e))
* **cli:** add support for custom local port in `glasskube open` ([#543](https://github.com/glasskube/glasskube/issues/543)) ([b6d98ca](https://github.com/glasskube/glasskube/commit/b6d98ca72ed87e6331b4be056e80912165aab92a))
* **ui:** sort package versions descending ([#308](https://github.com/glasskube/glasskube/issues/308)) ([98cd78c](https://github.com/glasskube/glasskube/commit/98cd78cd6e6ee1bb3870036c9fd3b11dc7ab14ec))


### Bug Fixes

* **deps:** update dependency clsx to v2.1.1 ([5cc6ae9](https://github.com/glasskube/glasskube/commit/5cc6ae9d70d3f00a75fbb3cef0169320006db14f))
* **deps:** update module golang.org/x/term to v0.19.0 ([d6057c3](https://github.com/glasskube/glasskube/commit/d6057c32c1b4c4cb1482832562df3789d9dcafb5))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.18.0 ([ae98bf0](https://github.com/glasskube/glasskube/commit/ae98bf0d4e0a681eca959c97a3cbf191967a76a7))
* make `glasskube version` more helpful when not bootstrapped ([#547](https://github.com/glasskube/glasskube/issues/547)) ([03d4fd9](https://github.com/glasskube/glasskube/commit/03d4fd94f8e237280cee82858e9db5bbabd0981e))
* **ui:** avoid boolean parsing error ([#558](https://github.com/glasskube/glasskube/issues/558)) ([3e8ea6e](https://github.com/glasskube/glasskube/commit/3e8ea6ec9e51f58447b8c8cce51a7c6c269564e0))


### Other

* **deps:** update actions/checkout digest to 1d96c77 ([d097c76](https://github.com/glasskube/glasskube/commit/d097c76dea0e934b1d8eeed5b5a1bdb68e7991b8))
* **deps:** update actions/upload-artifact action to v4.3.3 ([2ce6f25](https://github.com/glasskube/glasskube/commit/2ce6f25db861efbfe00c90f056de6786cb8bdc57))
* **deps:** update dependency @commitlint/cli to v19.3.0 ([1319b2a](https://github.com/glasskube/glasskube/commit/1319b2adf2970c1338a1db7f9585c5c0b5d4f843))
* **main:** update readme to reflect beta launch ([acbd3d6](https://github.com/glasskube/glasskube/commit/acbd3d6eb14e98ce5ad4bc8e2a111959d3ef1c50))
* update release-please configuration ([ff6c925](https://github.com/glasskube/glasskube/commit/ff6c925534b621d7ac05109e33ca6b54e96a6800))
* update telemetry ([8010bd7](https://github.com/glasskube/glasskube/commit/8010bd7b5dc23ec97dc8c9633b75760f5e5b8374))


### Docs

* add product hunt banner ([16acffc](https://github.com/glasskube/glasskube/commit/16acffc8d89c36f78b6c34fdb72073cbede7919b))
* update Product Hunt launch rul ([cf913bf](https://github.com/glasskube/glasskube/commit/cf913bf2fa80b593981aaf3cdbba4fb74a017c53))
* **website:** add producthunt launch banner ([20ad520](https://github.com/glasskube/glasskube/commit/20ad52020e9f63cba68f1e9d3a4eb4ace597f9a7))
* **website:** include GitHub star button in header ([#459](https://github.com/glasskube/glasskube/issues/459)) ([2e0e3a5](https://github.com/glasskube/glasskube/commit/2e0e3a56eb6f16637b05210fba69f056e3ff6e45))
* **website:** remove version disclaimer ([f72d0d1](https://github.com/glasskube/glasskube/commit/f72d0d15b960570425fc66e46edcd0e461e7abe7))


### Refactoring

* change Glasskube tagline ([9bdf826](https://github.com/glasskube/glasskube/commit/9bdf8267ba59348ada443a31be688b7c7dc56a32))

## [0.2.1](https://github.com/glasskube/glasskube/compare/v0.2.0...v0.2.1) (2024-04-22)


### Bug Fixes

* **client:** fallback to regular client if item not yet in cache ([bdd566e](https://github.com/glasskube/glasskube/commit/bdd566e786d50443827864599c29cbbf9ff0ef92))
* **deps:** update module github.com/onsi/gomega to v1.33.0 ([ad8da80](https://github.com/glasskube/glasskube/commit/ad8da80a64a15862f7581a900d7b47e3370cc2ea))
* **package-operator:** do not set owner reference on existing resources ([e447370](https://github.com/glasskube/glasskube/commit/e447370f1c1abe36d62d3a1f49ab3cf379e5732e))
* **ui:** restructure server initializations ([db1a585](https://github.com/glasskube/glasskube/commit/db1a585e4d7ebb6a2a6c82fdea3d56ccca76b795))


### Other

* **deps:** update actions/upload-artifact action to v4.3.2 ([e89e27b](https://github.com/glasskube/glasskube/commit/e89e27be833dfdb29a4fa09dd562aaefc248a57a))
* **deps:** update dependency go to v1.22.2 ([e24bd86](https://github.com/glasskube/glasskube/commit/e24bd86456000e9a2dc6fcaaf2fd3280261163cd))


### Docs

* **website:** add newlines ([c494978](https://github.com/glasskube/glasskube/commit/c494978331d7f36206bf286de382b1f5cd083a53))
* **website:** added guide plus updates description ([a03854e](https://github.com/glasskube/glasskube/commit/a03854e67f56f041fa2e52b3fdb55b1b25ff1656))
* **website:** unify package logos from GitHub discussion ([77d5601](https://github.com/glasskube/glasskube/commit/77d560146a5cf48020e1a83e7941900dd71e0275))


### Refactoring

* http errors handled when calling glasskube bootstrap ([b6fdf90](https://github.com/glasskube/glasskube/commit/b6fdf9024fe9d94d786311e3aa0887a251e99a23))

## [0.2.0](https://github.com/glasskube/glasskube/compare/v0.1.0...v0.2.0) (2024-04-18)


### Features

* add graph-based dependency validation ([c9957e6](https://github.com/glasskube/glasskube/commit/c9957e6516f3e056dc51bee2af3469a718cf795e))
* add kubernetes client adapter ([9623eff](https://github.com/glasskube/glasskube/commit/9623eff8c34b472b3371941db32d7bf2f9580130))
* add partial validation for package values ([5cfcb7c](https://github.com/glasskube/glasskube/commit/5cfcb7c2cf27e1266b67f359b40abf44c915fd1f))
* add resolving, validating value configurations ([54d0f0f](https://github.com/glasskube/glasskube/commit/54d0f0f556a7ddefb5e15e5d252abbad8aa24948))
* add support for transitive dependencies ([c9957e6](https://github.com/glasskube/glasskube/commit/c9957e6516f3e056dc51bee2af3469a718cf795e))
* **api:** add pattern constraint for value definitions ([de5f040](https://github.com/glasskube/glasskube/commit/de5f04047231e623b1b413d40fc1ae894c6add4f))
* **api:** add types for value configurations ([8e1c663](https://github.com/glasskube/glasskube/commit/8e1c66365ac894e3de7c406d698021b33fcb632b))
* **cli:** add package configuration ([ba0cd32](https://github.com/glasskube/glasskube/commit/ba0cd324794f60d9a5ee7b7dbd3e12085607f94e))
* **cli:** add value config flags ([9be6b1f](https://github.com/glasskube/glasskube/commit/9be6b1ffd23130574ae138b7ac9d5e225891e17d))
* **cli:** added --yes flag for non-interactive modes for install command ([#468](https://github.com/glasskube/glasskube/issues/468)) ([3ee7308](https://github.com/glasskube/glasskube/commit/3ee730825c5a74856562ec6cace3b847886f9208))
* **cli:** added --yes flag for non-interactive modes for update command ([#468](https://github.com/glasskube/glasskube/issues/468)) ([4e5d742](https://github.com/glasskube/glasskube/commit/4e5d7424effe2dc96fafcf448f085646944179c4))
* **client, ui:** introduce package info cache ([#444](https://github.com/glasskube/glasskube/issues/444)) ([24d6466](https://github.com/glasskube/glasskube/commit/24d64668df90fdc2f7ee36f332c06a5995ebab94))
* **cli:** made --enable-auto-updates flag value default even if --yes flag is used ([5d864a4](https://github.com/glasskube/glasskube/commit/5d864a4d585be6d61f307d4d267a1014f5a92c20))
* **cli:** show removed dependencies before uninstall ([3e84e09](https://github.com/glasskube/glasskube/commit/3e84e09e46a92f0dfa920f0a183bb996a2d23c11))
* **cli:** show whether a package is auto-updated ([#296](https://github.com/glasskube/glasskube/issues/296)) ([49a203a](https://github.com/glasskube/glasskube/commit/49a203a22f5b9ca355e351a5c2d22ef8cb0bcffb))
* **cli:** update to specific package version cli ([aa7649a](https://github.com/glasskube/glasskube/commit/aa7649accfcd4f385b07f386bed455343cc09774))
* **package-operator:** add base64 func to value templates ([3940dd2](https://github.com/glasskube/glasskube/commit/3940dd2a8f3aaad34cd8e0282c75a644a1b8cced))
* **package-operator:** add generating patches from value definitions ([90b264d](https://github.com/glasskube/glasskube/commit/90b264d85a2fbd6c2f2c3085f1eef1e7023dd583))
* **package-operator:** add handling package values ([8eb8e2b](https://github.com/glasskube/glasskube/commit/8eb8e2b7ce8890a6dd875f5df80e2ca74fce8346))
* **package-operator:** add validating package values in webhook ([1bfb8bf](https://github.com/glasskube/glasskube/commit/1bfb8bf2d78acc488c15a94369d5b1aece05010d))
* **ui:** alert when websocket has been closed ([#222](https://github.com/glasskube/glasskube/issues/222)) ([a86fc04](https://github.com/glasskube/glasskube/commit/a86fc04017ef0066ad289abac110cfcfcfedd827))
* **ui:** display latest & installed version in package overview ([#452](https://github.com/glasskube/glasskube/issues/452)) ([374253f](https://github.com/glasskube/glasskube/commit/374253f331be78a3c72ac4387da15100ce8f6688))
* **ui:** introduce global error handling ([d08abc0](https://github.com/glasskube/glasskube/commit/d08abc0d4dacda75a498aaf00542e64b0ce16a23))
* **ui:** package configuration ([#121](https://github.com/glasskube/glasskube/issues/121)) ([7c9e3d7](https://github.com/glasskube/glasskube/commit/7c9e3d7447d5a1fa088bd7d035e79768b5f8e2bd))
* **ui:** reuse transaction when applying updates ([#295](https://github.com/glasskube/glasskube/issues/295)) ([c26c471](https://github.com/glasskube/glasskube/commit/c26c47143123b6400586b9ee8335cf764082bfa7))
* **ui:** show removed dependencies before uninstall ([200ea5d](https://github.com/glasskube/glasskube/commit/200ea5d9952106d644f5f0f76967cb398ed19377))
* **ui:** show warning if operator and client versions differ ([#352](https://github.com/glasskube/glasskube/issues/352)) ([d50b164](https://github.com/glasskube/glasskube/commit/d50b16429e91a5aaaa5d2d89c56301bf56d2c1cc))
* **ui:** show whether a package is auto-updated ([#296](https://github.com/glasskube/glasskube/issues/296)) ([a79a636](https://github.com/glasskube/glasskube/commit/a79a636f15a5f56c1f573a8ae695b06a18ad05a8))
* **ui:** show whether a package is auto-updated ([#296](https://github.com/glasskube/glasskube/issues/296)) ([a79a636](https://github.com/glasskube/glasskube/commit/a79a636f15a5f56c1f573a8ae695b06a18ad05a8))
* **ui:** show whether a package is auto-updated ([#296](https://github.com/glasskube/glasskube/issues/296)) ([a79a636](https://github.com/glasskube/glasskube/commit/a79a636f15a5f56c1f573a8ae695b06a18ad05a8))
* **ui:** trigger refreshing package detail page ([#382](https://github.com/glasskube/glasskube/issues/382)) ([48f02bb](https://github.com/glasskube/glasskube/commit/48f02bbd593324dcda3dfeb14615e4c3e92644cb))


### Bug Fixes

* always try to resolve all values to improve error message ([b305bf9](https://github.com/glasskube/glasskube/commit/b305bf94d17654330c5c3a1b7c5380ddfd123572))
* **api:** change value constraints to pointers ([235a456](https://github.com/glasskube/glasskube/commit/235a456c003edb8fdcef60d6cab7822f9825aa64))
* avoid client rate limiting ([ab42a81](https://github.com/glasskube/glasskube/commit/ab42a81a931b389f9e4929584fcb8bc5c784d3e0))
* avoid updates to smaller versions ([4a44e4d](https://github.com/glasskube/glasskube/commit/4a44e4d9d99a10acba2267339a52dabefe5c734a))
* **cli, ui:** bootstrap handles nil values in `checkWorkloadReady` and major speedup ([d0b6683](https://github.com/glasskube/glasskube/commit/d0b66836306133e0dd507873da33b7ac2731f58b))
* **cli:** `glasskube ls` no longer shows the latest version if the installed version is newer ([#483](https://github.com/glasskube/glasskube/issues/483)) ([6725d40](https://github.com/glasskube/glasskube/commit/6725d40b97cb6d7986d14673f87a04f53dacd7be))
* **deps:** update dependency asciinema-player to v3.7.1 ([8a44cfa](https://github.com/glasskube/glasskube/commit/8a44cfa35760a5b6dd1848df1bfb9ce94386ce3b))
* **deps:** update dependency htmx.org to v1.9.12 ([20e3491](https://github.com/glasskube/glasskube/commit/20e34916da858f4914d0738b2a54a0cc4c7357cb))
* **deps:** update docusaurus monorepo to v3.2.0 ([1e65bc0](https://github.com/glasskube/glasskube/commit/1e65bc06351df0ed4b3a310f23a42ce2af3b21d5))
* **deps:** update docusaurus monorepo to v3.2.1 ([9923205](https://github.com/glasskube/glasskube/commit/9923205b1dc73832c3e026ee8b58aff9845969cf))
* **deps:** update font awesome to v6.5.2 ([995b95d](https://github.com/glasskube/glasskube/commit/995b95db7167ca1de16d21606955d1d913b03637))
* **deps:** update kubernetes packages to v0.29.4 ([40ad3c4](https://github.com/glasskube/glasskube/commit/40ad3c402b6f61773003b98c8dfd6690f900c542))
* **deps:** update module github.com/evanphx/json-patch/v5 to v5.9.0 ([62c5a55](https://github.com/glasskube/glasskube/commit/62c5a559e8d567053d56341d16330494eaa6d9fb))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.17.1 ([a1faa39](https://github.com/glasskube/glasskube/commit/a1faa393eb7bc3c2f3a7728330c306bd457f5e33))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.17.3 ([82c53f9](https://github.com/glasskube/glasskube/commit/82c53f99600ddc02f76491ce3a0882e7a4930199))
* made changes ([49a203a](https://github.com/glasskube/glasskube/commit/49a203a22f5b9ca355e351a5c2d22ef8cb0bcffb))
* **package-operator:** correct incorrect constraint error in values validation ([bea16a1](https://github.com/glasskube/glasskube/commit/bea16a1b3cbdd37833db076b6cc02dc779585fa2))
* **package-operator:** don't skip update validation if values changed ([536a7aa](https://github.com/glasskube/glasskube/commit/536a7aa3711170d6cd94e3a5c552642ad1c9d2b7))
* **package-operator:** prevent reconcile loop on packageinfo error ([aced43f](https://github.com/glasskube/glasskube/commit/aced43f742b9ea878b4de9bf816f3fddcb1d2a96))
* **ui, cli:** prevent a panic if a package has no OwnedPackageInfo ([#419](https://github.com/glasskube/glasskube/issues/419)) ([300490a](https://github.com/glasskube/glasskube/commit/300490acfdb6a0bba340e1fd542df1b7326e2918))
* **ui:** add hx-boost to fix inconsistent browser back/forth ([dc89c02](https://github.com/glasskube/glasskube/commit/dc89c029dd3146ee0ffdf298f8c4597ca4d2b73f))
* **ui:** make the whole card on package list clickable ([b4dcdd2](https://github.com/glasskube/glasskube/commit/b4dcdd215ddfecd2a541a174742bc582dfce83db))
* **ui:** parse templates before starting informer to avoid panic ([45fea63](https://github.com/glasskube/glasskube/commit/45fea63d6a0de8109d72f549b7c2409a59cb42cd))
* **ui:** take configured values into account at installation ([cd08eb9](https://github.com/glasskube/glasskube/commit/cd08eb9278ee20f841a5450ae6a8b1a2fea197e3))
* update install.go ([49a203a](https://github.com/glasskube/glasskube/commit/49a203a22f5b9ca355e351a5c2d22ef8cb0bcffb))


### Other

* add renovate automerge config ([6159c2b](https://github.com/glasskube/glasskube/commit/6159c2b00983e43124a3ddf3f3808a11611812ee))
* add telemetry ([#506](https://github.com/glasskube/glasskube/issues/506)) ([2221c64](https://github.com/glasskube/glasskube/commit/2221c64937f53521740365133d360e7ececdcccc))
* change next release to 0.2.0 ([3600c09](https://github.com/glasskube/glasskube/commit/3600c093f7f2595a5a692254caf02b9ec31764bb))
* **deps:** update commitlint monorepo to v19.2.2 ([cebb9ad](https://github.com/glasskube/glasskube/commit/cebb9adf5dc65895e338fb02cb01c8006be19036))
* **deps:** update dependency typescript to v5.4.4 ([e71dbc1](https://github.com/glasskube/glasskube/commit/e71dbc152d04758adba1c618476f5daa54751025))
* **deps:** update dependency typescript to v5.4.5 ([1bd3dfb](https://github.com/glasskube/glasskube/commit/1bd3dfb007c41481cc8621cb7ce66a4f2fa4cc09))
* **deps:** update flux manifests to version v2.2.3 ([aa1189e](https://github.com/glasskube/glasskube/commit/aa1189e1814ad9d396e736ef2fdc7a5937244024))
* remove unneeded OwnerManager in DependencyManager ([036db63](https://github.com/glasskube/glasskube/commit/036db633567b42d18ed80ce414dd760bf3920454))
* **repo:** address contributor guideline changes ([2ee73e6](https://github.com/glasskube/glasskube/commit/2ee73e68588249d16f1c672c32687c0b54ba24a5))
* **repo:** rebased and updated branch naming convention example ([69fa659](https://github.com/glasskube/glasskube/commit/69fa65955377812583aedac3f3b2cc831b4b4aef))
* **repo:** try to remove merge commit ([d7c8981](https://github.com/glasskube/glasskube/commit/d7c8981709e22c98f0885ceea400aa00ac1c1839))
* **repo:** update contributor guidelines ([bb9240c](https://github.com/glasskube/glasskube/commit/bb9240cdf3be34b34d0b71221291ac1ceb6bdaa6))
* **ui:** add reloading templates after changes in dev mode ([#170](https://github.com/glasskube/glasskube/issues/170)) ([76823e9](https://github.com/glasskube/glasskube/commit/76823e9f225a8b89d1e993280fabd3ad733e532b))
* **ui:** remove unnecessary htmx swap attributes ([08a0edb](https://github.com/glasskube/glasskube/commit/08a0edbdb353a5f52c75318ba89fd737469e8c2f))
* update golang version to v1.22 ([929bfd9](https://github.com/glasskube/glasskube/commit/929bfd9d9c0718c3b0d374bcc66bbfaa6ea38b61))
* **website:** update author avatar urls ([fe6c7ae](https://github.com/glasskube/glasskube/commit/fe6c7aed0083fe052eee65fea61a73d24c8dca4e))
* **website:** update link in blog + added video to cont guidelines ([c6a6c94](https://github.com/glasskube/glasskube/commit/c6a6c94b8fa8e8f6cd6aec06b9632a61a5082a68))
* **website:** updated video thumbnail ([e7f318e](https://github.com/glasskube/glasskube/commit/e7f318e1e92bc7bca2b33f3f540cde47e47a077b))


### Docs

* add glasskube activity chart ([1b63bb9](https://github.com/glasskube/glasskube/commit/1b63bb95177de19ea3face7a03ea25b361ed72b1))
* add instructions for installation via nix ([478d5ec](https://github.com/glasskube/glasskube/commit/478d5ec534feaa5c8179cfc08154b786919657fa))
* add proposal for package configuration ([#446](https://github.com/glasskube/glasskube/issues/446)) ([a7f3f92](https://github.com/glasskube/glasskube/commit/a7f3f9227dd05428abd0a5dcd9d11e1ce9a44e11))
* as Timoni is still being created, rephrase the word was to is ([83a7666](https://github.com/glasskube/glasskube/commit/83a76669c64a969d77852901222dc1ff6bef34b5))
* build numbers not allowed in dependency version ranges ([#405](https://github.com/glasskube/glasskube/issues/405)) ([594de21](https://github.com/glasskube/glasskube/commit/594de216c110145573be99b473641afcf8e62b8c))
* fix Timoni link to homepage ([402549d](https://github.com/glasskube/glasskube/commit/402549d2d1b4c22b90e152b5222dc78e9d1609c1))
* update CLI reference ([017c90d](https://github.com/glasskube/glasskube/commit/017c90d3eb96f10251e6167889d8624c3ba0e780))
* update contributor guidelines ([f037c2c](https://github.com/glasskube/glasskube/commit/f037c2c4d9e8a24e7458a6a5a404280f26fa5e3a))
* update go version to 1.22 ([b375170](https://github.com/glasskube/glasskube/commit/b375170c69ed9415807734e07e868f2bf04389fb))
* update readme ([0d859cf](https://github.com/glasskube/glasskube/commit/0d859cfa13cc13ff8c699a951337ea385c449403))
* **website:** address PR change requests ([25fda94](https://github.com/glasskube/glasskube/commit/25fda94d8863d8ba49f29f194ea2dc9585f86948))
* **website:** new blog post - contributor guidelines ([e9b6ad3](https://github.com/glasskube/glasskube/commit/e9b6ad312424710649a7e5a9201e5f8dd2ac797d))
* **website:** new Discord blog post + updated pic for v0.1.0 blog ([6e20133](https://github.com/glasskube/glasskube/commit/6e201333a0513e8f7735b77a55713b1bc4792d28))
* **website:** release blog v0.1.0-ammend2 ([cbf9541](https://github.com/glasskube/glasskube/commit/cbf95414628a5bb1e6dd12aa14814bdb025bbb62))
* **website:** update roadmap ([b9cfd73](https://github.com/glasskube/glasskube/commit/b9cfd73054d27ceafb6ce3666be244588722348a))
* **website:** update velero icon url to improve browser compatibility ([#488](https://github.com/glasskube/glasskube/issues/488)) ([17260ab](https://github.com/glasskube/glasskube/commit/17260abe4a75fd697381f345ec390cca0b9d49ab))


### Refactoring

* move adapters out of dependency package ([3ebba24](https://github.com/glasskube/glasskube/commit/3ebba24ea49bf0a93fd35991f4b913f46117b493))

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
