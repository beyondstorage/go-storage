# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [v5.0.0] - 2021-10-15

Since this release, we will use `go.beyondstorage.io/v5` as new module name.

### Added

- GSP-840: Convert to monorepo (#840)
- GSP-839: Add Support for Content-Disposition (#839)

### Changed

- *: Rename to go.beyondstorage.io (#842)
- credential,endpoint: Split from pkg (#843)
- refactor(tests): Move go-integration-tests here (#847)

## [v4.8.0] - 2021-09-30

### Added

- feat: Generate code via Xuanwo/gg (#728, #737, #738, #739, #740, #742)
- docs/rfcs: GSP-749 Unify Path Behavior (#749)
- docs: Add README.zh-CN.md (#750)
- docs/rfcs: GSP-751 Write Empty File Behavior (#751)
- docs: Add description of GSP-751 (#759)
- feat: Implement GSP-749 unify path behavior (#822)
- feat: Replace separator with / in path for service (#823)
- docs/rfcs: GSP-826 Add Multipart HTTP Signer Support (#826)
- feat: Add Multipart HTTP Signer support (#829)

### Changed

- docs: Add the services bos and obs to the README (#748)
- refactor: Cleanup definition generate logic (#754)
- feat: Migrate internal iterator generator to definitions (#755)
- ci: Sync github actions (#757)

### Fixed

- fix: Compilation fails if Default*Pairs unimplemented in service (#821)
- fix: Not enough arguments in call to Service.formatError (#824)
- fix: Replace each backslash in path with a slash (#825)
- fix: Fixed the error when no function is declared under the namespace (#831)

### Upgraded

- build(deps): bump github.com/pelletier/go-toml from 1.9.3 to 1.9.4 (#746)

## [v4.7.0] - 2021-09-02

### Added

- tests: Add test for NewStorageFromString (#724)
- docs/rfcs: GSP-725 Add Defaultable Property for Pair (#725)
- cmd/definitions: Implement GSP-725 Add Defaultable Property for Pair (#730)
- docs/rfcs: GSP-729 Redesign HTTP Signer (#729)
- definitions: Implement GSP-729 Redesign HTTP Signer (#732)

### Changed

- docs/rfc: Remove pair declare description for GSP-729 (#735)

### Deprecated

- definitions: Deprecate Reacher and expire (#734)

## [v4.6.0] - 2021-08-25

### Added

- feat: Support features and default paris in connection string (#713)
- feat: Set Features with feature pairs (#716)
- feat: Allow append pairs for New{Servicer,Storager}FromString (#718)
- feat: Set DefaultPairs with defaultable pairs (#719)
- ci: Add cross build (#721)

### Changed

- cmd/definitions: Remove conflict from pair (#720)

### Fixed

- definitions: Fix naming format for http signer (#714)

## [v4.5.0] - 2021-08-19

### Added

- tests: Implement bench for List (#687)
- ci: Add golangci-lint (#690)
- docs/rfcs: GSP-700 Config Features and DefaultPairs via Connection String (#700)
- docs/rfcs: GSP-706 Support HTTP Signer (#706)
- feat: Generate feature pairs (#708)
- definitions: Implement GSP-706 Support HTTP Signer (#709)
- feat: Generate default pairs (#710)

### Changed

- docs: Update go-storage docs link (#696)
- docs: Update links to matrix (#701)
- docs: Update readme about services support (#705)

### Fixed

- *: Fix ineffectual assignments (#691)
- definitions: Fixed a description error in CreateLink (#698)

## [v4.4.0] - 2021-07-22

### Added

- docs: GSP-669 Feature Lifecycle (#669)
- definitions: Implement GSP-86 Add CreateLink (#674)

### Changed

- docs: Update rfc README, fix GSP-669 typo (#673)
- cmd: Remove not used package in service.tmpl (#679)
- definitions: Implement GSP-41 Turn Pair Expire into Duration (#681)
- tests: Split from cmd definitions (#685)

### Upgraded

- build(deps): Bump github.com/google/uuid from 1.2.0 to 1.3.0 (675)

## [v4.3.2] - 2021-07-15

### Fixed

- cmd: Add warning for type convert instead (#670)

## [v4.3.1] - 2021-07-15

### Fixed

- pkg: Add pkg/endpoint back to avoid breaking other projects (#667)

### Upgraded

- build(deps): bump github.com/golang/mock from 1.5.0 to 1.6.0 (#597)

## [v4.3.0] - 2021-07-14

### Added

- cmd/definitions: Add support for type uint64 (#615)
- services: Implement GSP-127 Add ServiceInternal and RequestThrottled Errors (#616)
- ci: Add codecov in unittests (#622)
- docs: GSP-654 Unify List Behavior (#654)
- ci: Add gofmt to action (#660)

### Changed

- docs: Reorganize README for better understanding (#613)
- types: Add comments for Object and ObjectMode (#614)
- cmd/definitions: Remove deprecated fields introduced in v4.2.0 (#619)
- pkg/endpoint: Remove deprecated package pkg/endpoint (#621)
- docs: Move rfcs and spec here as required in GSP-139 (#628)
- cmd: Move definitions here as described in GSP-139 (#652)
- docs: Update RFC template (#655)
- cmd/definitions: Implement GSP-134 Write Behavior Consistency (#656)

### Fixed

- makefile: Fix code not formatted correctly (#661)

## [v4.2.0] - 2021-06-24

### Added

- cmd/definitions: Implement GSP-97 Add Restrictions In Storage Metadata (#601)
- cmd/definitions: Implement GSP-111 Add System Metadata in Storage Metadata (#605)

### Changed

- cmd/definitions: Trim Service Pair Prefix (#598)
- services: Implement GSP-90 Re-support Initialization Via Connection String (#589)
- cmd/definitions: Implement GSP-109 Redesign Features (#606)
- cmd/definitions: Implement GSP-117 Rename Service to System as the Opposite to Global (#608)

### Deprecated

- pkg/endpoint: Deprecated by github.com/beyondstorage/go-endpoint (#595)
- services: Related fields `ServiceFeatures` and `StorageFeatures` have been deprecated (#606)
- services: `GetObjectMetadata` has been deprecated by `GetObjectSystemMetadata` (#608)

## [v4.1.0] - 2021-06-08

### Added

- cmd/definitions: Implement GSP-87 Feature Gates (#586)
- cmd/definitions: Implement GSP-93 Add ObjectMode Pair (#592)

### Fixed

- cmd/definitions: Fix storage features not generated (#588)

## [v4.0.0] - 2021-05-24

### Added

- storage: Implement GSP-44 (#547)
- services: Implement GSP-48 Service Registry (#556)
- services: Implement GSP-47 Additional Error Specification (#555)
- storage: Implement GSP-49 Add CreateDir Operation (#561)
- services: Implement GSP-51 Distinguish Errors by IsInternalError (#564)
- types: Implement GSP-61 Add object mode check for operations (#573)
- cmd/definitions: Add functions Add and Del for ObjectMode (#574)
- cmd/definitions: Declare supports interfaces in generator (#584)

### Changed

- cmd/definitions: Remove redundant // in object.tmpl (#565)
- docs: Move spec to specs (#568)
- docs: Add comment to illustrate ignoring generated.go in tests package (#567)
- cmd/definitions: Update obsolete file paths in comments (#570)
- object: Implement GSP-62 WriteMultipart returns Part (#572)
- iowrap: Explicit declare implemented interface (#576)
- *: Implement GSP-73 Organization Rename (#578)
- storage: Implement GSP-76 Local Function Metadata (#77) (#580)

### Fixed

- types: Fix IterateDone doesn't implement InternalError (#582)

## [v3.6.0] - 2021-04-24

### Added

- object: Add AppendOffset for appendable object (#532)
- cmd/definitions: Add a new type byte_array (#534)
- pkg: Implement wrapper for fs.FS and http.FileSystem (#538)
- cmd/definitions: Implement GSP-40 (#541)

### Changed

- docs: Fix broken links in README (#533)
- docs/specs: Replace spaces by tab in the code in 1-error-handling (#535)

## [v3.5.0] - 2021-04-20

### Added

- cmd/definitions: Add tests service (#518)
- cmd/definitions: Increase stability of generated code (#526)
- types: Add UnimplementedStub to have forward compatible implementations (#524)
- pkg/iowrap: Add cancelable io wrappers (#528)

### Changed

- docs: Update README
- docs: Migrate zulip to matrix

### Fixed

- docs: Fix quick start is not runnable

## [v3.4.2] - 2021-03-16

### Added

- object: Add Multipart related fields into object (#516)

## [v3.4.1] - 2021-03-04

### Fixed

- cmd: Fix support for local function generation (#513)

## [v3.4.0] - 2021-03-04

### Fixed

- storage: New is conflict with Storage init logic, rename to Create instead (#511)

## [v3.3.0] - 2021-03-04 (deprecated)

### Added

- pkg/iowrap: Implement Pipe (#508)
- types: Add "New" operation to create an object locally (#509)

## [v3.2.0] - 2021-02-22

### Added

- iowrap: Implement CallbackWriter (#502)
- types: Implement Stringer for ObjectMode (#503)
- service: Add template for generating default pair for each service (#504)

## [v3.1.0] - 2021-02-18

### Added

- *: Implement proposal unify object metadata (#498)

### Changed

- cmd/definitions: Remove parsed pairs pointer (#500)

### Removed

- pkg: Remove not used storageclass package

## [v3.0.0] - 2021-01-15

### Added

- pairs: Add support for user-agent (#477)
- operation: Add fetcher (#480)
- Proposal: Add default pair for operations (#484)
- types: Implement proposal Unify List Operation (#489)
- types: Implement proposal segment api redesign (#490)
- cmd: Implement code generate and format (#491)
- types: Implement proposal Object Mode (#493)

### Changed

- cmd/definitions: Don't need to store definitions to bindata (#476)
- cmd: Introduce beyondstorage/specs to maintain specs (#481)
- docs: Migrate design to beyondstorage/specs (#488)
- cmd: Migrate from hcl to toml (#496)

## [v2.0.0] - 2020-11-12

### Changed

- cmd/definitions: Merge into main modules (#465)
- cmd: Add tools tag into build (#468)

### Fixed

- cmd/definitions: Fix service not generated correctly (#466)
- cmd/definitions: Fix server pair not handled correctly (#472)

## [v2.0.0-beta] - 2020-11-09

### Added

- types: Implement pair policy (#453)
- pkg/storageclass: Add sotrageclass support (#456)

### Changed

- build: Use beyondstorage/go-dev-tools to tidy go mod files (#454)
- cmd/install: Move to beyondstorage/go-dev-tools/setup
- pairs: Use dot import to avoid type conflicts (#459)
- build(deps): bump github.com/google/uuid from 1.1.1 to 1.1.2 (#461)

## [v2.0.0-alpha.1] - 2020-11-02

### Added

- Support Iterator based list operation
- types/iterator: Allow store current status in iterator
- types/object: Support linked set
- types/iterator: Make Page.Status continuable (#433)
- types/object: Add link support (#438)
- types: Add interceptor support (#449)

### Changed

- *: Moving to beyondstorage/go-storage (#414)
- cmd: Move definitions to cmd to support service split (#416)
- types/object: Move all meta into ObjectMeta
- Return count in storager when read and write (#427) 
- types: Refactor into struct for object stat support
- types/pair: Use struct instead of pointer (#435)
- *: Improve minimum supported version to go 1.14 (#444)
- makefile: Manage build tools via go modules (#447)

### Fixed

- types: Fix object stat not updated correctly
- types: Fix bit operations not correctly (#434)

### Removed

- coreutils: Split into beyondstorage/go-coreutils (#417)
- tests: Move to beyondstorage/go-storage-integration-test (#418)
- services: Split all services into separate repos (#419)
- types: Remove no used object meta

## [v1.2.1] - 2020-06-30

### Changed

- internal: Generate all exported APIs (#361)
- services/fs: Convert system specific separator to slash (#408)

### Fixed

- services/qingstor: Fix WorkDir listed in keys while ListDir (#366)
- definitions: Fix statistical's result is incorrect (#367)
- services/qingstor: Fix unit test for ListDir (#404)

## [v1.2.0] - 2020-05-20

### Added

- tests: Add integration test for qingstor (#325)
- pkg/iowrap: Add SizedReadSeekCloser support (#329)
- tests: Add integration test for azblob (#331)
- tests: Add integration test for s3 (#338)
- tests: Add integration test for gcs (#341)
- tests: Add integration test for cos (#342)
- tests: Add integration test for oss (#343)
- tests: Add integration test for kodo (#347)

### Changed

- services/fs: Auto create work dir (#324)
- pkg/endpoint: Allow omit protocol default port (#346)
- tests: Compare content sha256 instead of full content

### Fixed

- services/*: Fix context not initiated (#328)
- services/azblob: Fix content length not set correctly (#330)
- services/s3: Fix NotFound error not handled correctly (#332)
- services/s3: Fix ListDir && ListPrefix not ended correctly (#333)
- services/s3: Fix Read's pair not parsed correctly (#334)
- services/s3: Fix bucket in request input not assigned (#335)
- services/s3: Location should be required in Storager Init (#336)
- services/s3: Don't calculate content-sha256 as default (#337)
- services/gcs: Fix ListDir and ListPrefix not ended correctly (#339)
- services/gcs: Fix oauth2 token source not configured correctly (#340)
- services/oss: Fix ListPrefix && ListDir not ended correctly
- services/oss: Fix Read's pairs not parsed correctly
- services/oss: Fix LastModified not parsed correctly
- services/uss: Fix content-length header not filled
- services/uss: Fix ListDir && ListPrefix data race while returning error
- services/uss: Fix object channel double close
- services/uss: Fix object not iterated fully
- services/uss: Fix dead loop on reading from io.Pipe
- services/uss: Use async delete to avoid 'concurrent' delete
- services/uss: Disable async delete for unexpected behavior (#345)
- services/kodo: Fix kodo's domain not setup correctly
- services/kodo: Fix object not found error not formatted correctly
- services/kodo: Fix request deadline not set correctly
- services/kodo: Fix service level error not parsed

## [v1.1.0] - 2020-05-14

### Proposal

- services, design: Propose and implement 21-split-segmenter (#270)

### Added

- types/pairs: Add more comment for work dir
- pkg/httpclient: Add stream-oriented http client support (#274)
- services/qingstor: Detect bucket location automatically (#278)
- ci: Setup drone for integration test (#282)
- services/qingstor: Allow read with offset and size (#283)
- services: Add service level pair support (#311)
- services/*: Add service level storage class support (#313)
- services/qingstor: Add disable uri cleaning support (#314)

### Changed

- types/pairs: Allow parse from plain string value (#281)
- docs: Use vuepress instead (#286)
- docs: Import user experience (#287)
- docs: Get ready for i18n docs (#289)
- docs: Add crowdin support for i18n (#290)
- *: Refactor definitions generator (#303)
- definitions: Auto inject http_client_options for all services (#321)

### Fixed

- internal/cmd: Fix -ignore flag in go generate matches unexpected files (#268)
- services/azblob: Fix context deadline exceeded while reading (#275)
- pkg/httpclient: Fix connection closed while writing or reading (#305)

## [v1.0.0] - 2020-04-23

### Proposal

- storager, services: Implement proposal 19-split-storage-list (#249) 
- types, services: Implement proposal 20-remove-loose-mode (#250) 

### Added

- types/pairs: Add description for paris (#238)
- services/azblob: Add offset and size support for Read (#252) 
- types/pairs: Add Parse support (#260) 

### Changed

- services/qingstor: Dir support is not needed in stat (#223)
- services/*: Set content-type as best errort (#234)
- services/*: Check metadata value before set (#235) 
- pkg/headers: Move all explicit header keys into const (#236)
- tests: Refactor integration test (#261) 
- coreutils, services: Refactor Open related logic (#264)

### Fixed

- services/*: Fix error not handled correctly with empty body (#227)
- services/qingstor: Fix ListBuckets incomplete (#228)
- services/*: Fix getAbsPath and getRelPath not handled correctly (#253) 
- internal/cmd/service: Fix context pair not handled correctly (#254)


## [v0.9.0] - 2020-03-23

### Proposal

- docs/design: Add 14-normalize-content-hash-check (#186)
- docs/design: Add proposal release policy (#192)
- docs/design: Add proposal loose mode (#199)
- docs: Add proposal 17-proposal-process (#210)
- docs/design: Add proposal return-segment-interface-instead (#216)

### Added

- services/*: Implement 14-normalize-content-hash-check (#189)
- services/*, types/pairs: Implement proposal loose-mode (#200)
- pkg/segment, docs/design: Add and implement proposal return-segment-interface-instead (#216)
- services/s3: Add multipart support (#220)

## [v0.8.0] - 2020-03-09

### Added

- services/*: Add ReadCallbackFunc for WriteSegment (#169)
- docs/design, coreutils: Propose and implement proposal remove config string (#172)

### Changed

- services/cos: Refactor service newStorage (#176)

### Fixed

- services/fs: Fix size and offset pair not handled correctly (#175)
- services/kodo: Fix ID not set (#178)
- services/*: Handle errors returned by New (#179) 

## [v0.7.2] - 2020-03-05

### Added

- services/*: Implement proposal 11-error-handling (#143)
- docs/design: Support both directory and prefix based list (#157)
- services/*: Implement proposal 12-support-both-directory-and-prefix-based-list (#158)

### Changed

- services/*: Refactor format error (#166)

### Fixed

- services/{gcs,kodo,oss,s3}: Fix Object type in List incorrect (#162)

## [v0.7.1] - 2020-02-29

### Added

- docs/{design,spec}: Add proposal for error handling (#106)
- pkg/*: Implement proposal 11-error-handling (#109)
- services/qingstor: Implment proposal 11-error-handling (#117)
- services/fs: Implement proposal 11-error-handling (#141)

### Changed

- services/qingstor: Refactor work dir handler with unit tests (#139)

### Fixed

- services: Fix WorkDir support missing in some services (#131)
- services/qingstor: Fix error not handled as intended (#135)
- services/qingstor: Fix service qingstor error not handled correctly


## [v0.7.0] - 2020-02-10

### Added

- tests: Add bdd test for integration test (#81)
- docs/design, pkg/iowrap: Add and implement proposal 10-callback-reader (#88)

### Removed

- docs/design, services: Add and implement proposal 9-remove-storager-init (#79)

## [v0.6.0] - 2020-01-13

### Added

- services,types: Implement proposal add-id-in-object (#56)
- services, types/metadata: Implement proposal 6-normalize-metadata (#59)
- services: Add basic kodo support (#49)
- services: Add basic cos support (#65)
- services: Add dropbox basic support (#53)
- services: Add uss basic support (#67)
- *: Implement proposal 7-support-context (#68)
- services, pkg/storageclass: Add and implement proposal 8-normalize-metadata-storage-class #71

### Changed

- storager: Rename ListDir to List (#52)

## [v0.5.0] - 2019-12-30

### Added

- services: Add support for s3 (#41)
- services: Add basic oss support (#42)
- services: Add basic gcs support (#48)
- services: Add basic support for azblob (#50)

### Changed

- pkg/config: Allow emit host instead of credential
- pkg/credential: Implement proposal 4-credential-refactor

### Fixed

- sercices/s3: Fix error message for servicer (#44)

## [v0.4.0] - 2019-12-23

### Added

- servicer: Add String() for debug (#23)
- Implement proposal support service init via config string (#38)

### Changed

- internal: Refactor generator (#24)
- internal: Don't preserve files' metadata
- storager: Implement proposal 1-unify-storager-behavior (#30)
- *: Implement proposal 2-use-callback-in-list-operations (#31)
- services: Promote values into struct instead of metadata (#33)
- services: Split endpoint and credential into different pairs (#34)
- storager: Split Metadata to Metadata and Statistical (#39)

### Fixed

- services/posixfs: Fix std{in/out} support for Stat (#35)

## [v0.3.0] - 2019-11-11

### Added

- services: Add POSIX fs support (#1)
- services/posixfs: Add support for write stdout
- services/posixfs: Implement size and offset in Read (#8)
- services/qingstor: Add bucket name validate (#9)
- storager: Add String interface for debug print (#16)
- services/posixfs: Set updatedAt for regular file (#19)

### Changed

- storager: Merge Capablity and IsPairAvailable into Capable
- services: Unify behavior for ListDir recursively
- storager: All API now use relative path instead
- services/qingstor: Set default base value (#13)
- types: Rename Base to WorkDir (#17)

### Fixed

- services/qingstor: Fix segment ID used incorrectly
- services/posixfs: Fix ListDir not returned ErrDone
- pkg/segment: Fix data while update segments
- pkg/segment: Fix GetPartIndex bug in concurrent write segment
- pkg/segment: Return sorted parts
- services/posixfs: Fix ListDir recursive not work
- pkg/iterator: Fix next panic while under lying func return empty buf
- services/qingstor: Fix recursive not handled correctly
- services/posixfs: Fix relative path not returned correctly
- services/qingstor: Part number must be in ascending order (#14)
- services/qingstor: Fix get abs and rel path (#15)
- services/qingstor: Handle not found returned via head (#20)
- services/posixfs: Handle not found returned via os not exist error. (#21)

### Removed

- storager: Remove CreateDir (#18)

## [v0.2.0] - 2019-10-17

### Added

- services/qingstor: Implement ListSegments (#7)

### Changed

- segment: Replace identifier and add ListSegments support (#6)

### Fixed

- services/qingstor: Fix ListDir not handled correctly (#2)
- services/qingstor: Fix object size and last modified not filled (#4)
- services/qingstor: Add stat updated at support (#5)

## v0.1.0 - 2019-10-12

### Added

- Add storager and servicer interface.
- Add segment and iterator support.
- Add pair based option and metadata support.
- Add qingstor services.

[Unreleased]: https://github.com/beyondstorage/go-storage/compare/v5.0.0...HEAD
[v5.0.0]: https://github.com/beyondstorage/go-storage/compare/v4.8.0...v5.0.0
[v4.8.0]: https://github.com/beyondstorage/go-storage/compare/v4.7.0...v4.8.0
[v4.7.0]: https://github.com/beyondstorage/go-storage/compare/v4.6.0...v4.7.0
[v4.6.0]: https://github.com/beyondstorage/go-storage/compare/v4.5.0...v4.6.0
[v4.5.0]: https://github.com/beyondstorage/go-storage/compare/v4.4.0...v4.5.0
[v4.4.0]: https://github.com/beyondstorage/go-storage/compare/v4.3.2...v4.4.0
[v4.3.2]: https://github.com/beyondstorage/go-storage/compare/v4.3.1...v4.3.2
[v4.3.1]: https://github.com/beyondstorage/go-storage/compare/v4.3.0...v4.3.1
[v4.3.0]: https://github.com/beyondstorage/go-storage/compare/v4.2.0...v4.3.0
[v4.2.0]: https://github.com/beyondstorage/go-storage/compare/v4.1.0...v4.2.0
[v4.1.0]: https://github.com/beyondstorage/go-storage/compare/v4.0.0...v4.1.0
[v4.0.0]: https://github.com/beyondstorage/go-storage/compare/v3.6.0...v4.0.0
[v3.6.0]: https://github.com/beyondstorage/go-storage/compare/v3.5.0...v3.6.0
[v3.5.0]: https://github.com/beyondstorage/go-storage/compare/v3.4.2...v3.5.0
[v3.4.2]: https://github.com/beyondstorage/go-storage/compare/v3.4.1...v3.4.2
[v3.4.1]: https://github.com/beyondstorage/go-storage/compare/v3.4.0...v3.4.1
[v3.4.0]: https://github.com/beyondstorage/go-storage/compare/v3.3.0...v3.4.0
[v3.3.0]: https://github.com/beyondstorage/go-storage/compare/v3.2.0...v3.3.0
[v3.2.0]: https://github.com/beyondstorage/go-storage/compare/v3.1.0...v3.2.0
[v3.1.0]: https://github.com/beyondstorage/go-storage/compare/v3.0.0...v3.1.0
[v3.0.0]: https://github.com/beyondstorage/go-storage/compare/v2.0.0...v3.0.0
[v2.0.0]: https://github.com/beyondstorage/go-storage/compare/v2.0.0-beta...v2.0.0
[v2.0.0-beta]: https://github.com/beyondstorage/go-storage/compare/v2.0.0-alpha.1...v2.0.0-beta
[v2.0.0-alpha.1]: https://github.com/beyondstorage/go-storage/compare/v1.2.1...v2.0.0-alpha.1
[v1.2.1]: https://github.com/beyondstorage/go-storage/compare/v1.2.0...v1.2.1
[v1.2.0]: https://github.com/beyondstorage/go-storage/compare/v1.1.0...v1.2.0
[v1.1.0]: https://github.com/beyondstorage/go-storage/compare/v1.0.0...v1.1.0
[v1.0.0]: https://github.com/beyondstorage/go-storage/compare/v0.9.0...v1.0.0
[v0.9.0]: https://github.com/beyondstorage/go-storage/compare/v0.8.0...v0.9.0
[v0.8.0]: https://github.com/beyondstorage/go-storage/compare/v0.7.2...v0.8.0
[v0.7.2]: https://github.com/beyondstorage/go-storage/compare/v0.7.1...v0.7.2
[v0.7.1]: https://github.com/beyondstorage/go-storage/compare/v0.7.0...v0.7.1
[v0.7.0]: https://github.com/beyondstorage/go-storage/compare/v0.6.0...v0.7.0
[v0.6.0]: https://github.com/beyondstorage/go-storage/compare/v0.5.0...v0.6.0
[v0.5.0]: https://github.com/beyondstorage/go-storage/compare/v0.4.0...v0.5.0
[v0.4.0]: https://github.com/beyondstorage/go-storage/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/beyondstorage/go-storage/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/beyondstorage/go-storage/compare/v0.1.0...v0.2.0
