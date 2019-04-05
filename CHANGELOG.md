# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## 0.5.0 - 2019-04-05

### Changed

- Allow `Logger()` option to accept interface `StdLogger` instead of enforcing a specific logger implementation

## 0.4.0 - 2019-02-01

### Added

- Added `id` to `TestPipeline` and `DeploymentPipeline` returned from project calls

## 0.3.1 - 2018-07-17

### Fixed

- Remove `omitempty` struct tag from Project type field to fix creating Basic projects

## 0.3.0 - 2018-05-14

### Added

- Added `branch` to `Build`

## 0.2.2 - 2018-04-11

### Fixed

- Changed ProjectType `MarshalJSON` and `String` methods to be value receivers
so that they perform the necessary type name lookups when marshalling to JSON.

## 0.2.1 - 2018-04-08

### Added

- Return error when Authentication fails because 2FA is enabled

## 0.2.0 - 2018-01-26

### Added

- Added `id` to `Project`
- Added `project_id` to `Build`

## 0.1.0 - 2018-01-05

- Initial Release
