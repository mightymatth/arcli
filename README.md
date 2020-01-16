<p align="center">
  <img alt="arcli" src="img/arcli.png" height="200" width="200" />
  <h1 align="center">arcli</h1>
  <h3 align="center">Awesome Redmine CLI</h3>
</p>

[![Go Report Card](https://goreportcard.com/badge/github.com/mightymatth/arcli)](https://goreportcard.com/report/github.com/mightymatth/arcli)
[![arcli](https://snapcraft.io//arcli/badge.svg)](https://snapcraft.io/arcli)
## About

arcli is CLI for [Redmine](https://www.redmine.org/) that simplifies some actions such as checking for issue details and tracking time. It uses Redmine REST API and should work with any Redmine server which arcli currently supports. 

## Installation

### macOS

```sh
$ brew tap mightymatth/arcli https://github.com/mightymatth/arcli
$ brew install arcli
```

### Linux
[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/arcli)

```sh
snap install arcli
```

## Usage
```sh
➜  ~ arcli -h
Client for Redmine. Wrapper around Redmine API

Usage:
  arcli [command]

Available Commands:
  aliases     Words that can be used instead of issue or project ids.
  defaults    User session defaults.
  help        Help about any command
  issues      Shows issue details.
  log         Time entries on projects and issues.
  login       Authenticate to Redmine server.
  logout      Logout current user.
  projects    Shows project details.
  search      Search Redmine
  status      Overall account info

Flags:
  -h, --help      help for arcli
      --version   version for arcli

Use "arcli [command] --help" for more information about a command.
```
