# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
- [#16] Update groovy scripts for compatibility with nexus 3.82

## [1.1.2] - 2025-07-15
### Fixed
- Fixed the repository modify/create/delete scripts. The field "repositoryManager" is not accessible by getter anymore.
  - It is now private and accessed by @repositoryManager
- Fixed the createRepository script
  - In previous versions, when null was returned, it was converted to a string "null"
  - Now the api returns an empty string on null
  - The go-code was previously checking for the "null" string

## [1.1.1] - 2024-09-24
### Fixed
- Fixed the repository modify script. It did not work before and should now update repositories

## [1.1.0] - 2024-09-18
### Changed
- Relicense to AGPL-3.0-only

## [1.0.0] - 2020-06-18
### Changed
- migrates from glide to go modules to be able to use recent go versions for compile 
- integrates makefiles [see #2](https://github.com/cloudogu/nexus-claim/issues/2)
- updates to nexus api 3.23
- fixes issue related to api changes [see #7](https://github.com/cloudogu/nexus-claim/issues/7)
