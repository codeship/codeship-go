# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased

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
