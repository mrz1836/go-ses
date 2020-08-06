# go-ses
> Lightweight implementation of [AWS SES](https://aws.amazon.com/ses/) - fork from [go-ses](https://github.com/publica-project/go-ses)

[![Release](https://img.shields.io/github/release-pre/mrz1836/go-ses.svg?logo=github&style=flat&v=1)](https://github.com/mrz1836/go-ses/releases)
[![Build Status](https://travis-ci.com/mrz1836/go-ses.svg?branch=master)](https://travis-ci.com/mrz1836/go-ses)
[![Report](https://goreportcard.com/badge/github.com/mrz1836/go-ses?style=flat&v=1)](https://goreportcard.com/report/github.com/mrz1836/go-ses)
[![codecov](https://codecov.io/gh/mrz1836/go-ses/branch/master/graph/badge.svg)](https://codecov.io/gh/mrz1836/go-ses)
[![Go](https://img.shields.io/github/go-mod/go-version/mrz1836/go-ses)](https://golang.org/)
[![Sponsor](https://img.shields.io/badge/sponsor-MrZ-181717.svg?logo=github&style=flat&v=3)](https://github.com/sponsors/mrz1836)
[![Donate](https://img.shields.io/badge/donate-bitcoin-ff9900.svg?logo=bitcoin&style=flat)](https://mrz1818.com/?tab=tips&af=go-ses)

<br/>

<br/>

## Table of Contents
- [Installation](#installation)
- [Documentation](#documentation)
- [Examples & Tests](#examples--tests)
- [Benchmarks](#benchmarks)
- [Code Standards](#code-standards)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

<br/>

## Installation

**go-ses** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get -u github.com/mrz1836/go-ses
```

<br/>

## Documentation
View the generated [documentation](https://pkg.go.dev/github.com/mrz1836/go-ses)

[![GoDoc](https://godoc.org/github.com/mrz1836/go-ses?status.svg&style=flat)](https://pkg.go.dev/github.com/mrz1836/go-ses)

### Features
- Send raw or html emails
- Multiple to, cc, and bcc recipients
- AWS4 signatures

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

[goreleaser](https://github.com/goreleaser/goreleaser) for easy binary or library deployment to Github and can be installed via: `brew install goreleaser`.

The [.goreleaser.yml](.goreleaser.yml) file is used to configure [goreleaser](https://github.com/goreleaser/goreleaser).

Use `make release-snap` to create a snapshot version of the release, and finally `make release` to ship to production.
</details>

<details>
<summary><strong><code>Makefile Commands</code></strong></summary>
<br/>

View all `makefile` commands
```shell script
make help
```

List of all current commands:
```text
all                    Runs multiple commands
clean                  Remove previous builds and any test cache data
clean-mods             Remove all the Go mod cache
coverage               Shows the test coverage
godocs                 Sync the latest tag with GoDocs
help                   Show this help message
install                Install the application
install-go             Install the application (Using Native Go)
lint                   Run the Go lint application
release                Full production release (creates release in Github)
release                Runs common.release then runs godocs
release-snap           Test the full release (build binaries)
release-test           Full production test release (everything except deploy)
replace-version        Replaces the version in HTML/JS (pre-deploy)
run-examples           Runs all the examples
tag                    Generate a new tag and push (tag version=0.0.0)
tag-remove             Remove a tag if found (tag-remove version=0.0.0)
tag-update             Update an existing tag to current commit (tag-update version=0.0.0)
test                   Runs vet, lint and ALL tests
test-short             Runs vet, lint and tests (excludes integration tests)
test-travis            Runs all tests via Travis (also exports coverage)
test-travis-short      Runs unit tests via Travis (also exports coverage)
uninstall              Uninstall the application (and remove files)
vet                    Run the Go vet application
```
</details>

<br/>

## Examples & Tests
All unit tests and [examples](examples/simple/simple.go) run via [Travis CI](https://travis-ci.com/mrz1836/go-ses) and uses [Go version 1.14.x](https://golang.org/doc/go1.14). View the [deployment configuration file](.travis.yml).

#### Running Tests
1. Set the environment variables `$AWS_ACCESS_KEY_ID` and `$AWS_SECRET_KEY` and `$AWS_SES_ENDPOINT`.
2. Example SES Endpoint `https://email.us-west-2.amazonaws.com`
3. Run `go test -from=user@example.com`, where `user@example.com` is a sender address that is verified
   in your Amazon SES account.

Run all tests (including integration tests)
```shell script
make test
```

Run tests (excluding integration tests)
```shell script
make test-short
```

Run the examples:
```shell script
make run-examples
```

<br/>

## Benchmarks
Run the Go benchmarks:
```shell script
make bench
```

<br/>

## Code Standards
Read more about this Go project's [code standards](CODE_STANDARDS.md).

<br/>

## Usage
View the [examples](examples/simple/simple.go)

<br/>

## Maintainers
| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) |
|:---:|
| [MrZ](https://github.com/mrz1836) |

<br/>

## Contributing

Previous Contributors:

* Forked from [publica-project](https://github.com/publica-project/go-ses)
* Quinn Slack <sqs@sourcegraph.com>
* Patrick Crosby (author of original [stathat/amzses](https://github.com/stathat/amzses))

View the [contributing guidelines](CONTRIBUTING.md) and follow the [code of conduct](CODE_OF_CONDUCT.md).

### How can I help?
All kinds of contributions are welcome :raised_hands:! 
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:. 
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/mrz1836) :clap: 
or by making a [**bitcoin donation**](https://mrz1818.com/?tab=tips&af=go-ses) to ensure this journey continues indefinitely! :rocket:

<br/>

## License

![License](https://img.shields.io/github/license/mrz1836/go-ses.svg?style=flat&v=1)

Note: the public API is experimental and subject to change until further notice.


Usage
=====

Documentation: [go-ses on Sourcegraph](https://sourcegraph.com/github.com/sourcegraph/go-ses).

Example: [example_test.go](https://github.com/sourcegraph/go-ses/blob/master/example_test.go) ([Sourcegraph](https://sourcegraph.com/github.com/sourcegraph/go-ses/tree/master/example_test.go)):