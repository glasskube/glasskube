# Changelog

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
