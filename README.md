<p align="center">
  <img alt="arcli" src="img/arcli.png" height="200" width="200" />
  <h1 align="center">arcli</h1>
  <h3 align="center">Awesome Redmine CLI</h3>
</p>


## About

arcli is CLI for [Redmine](https://www.redmine.org/) that simplifies some actions such as checking for issue details and tracking time. It uses Redmine REST API and should work with any Redmine server which arcli currently supports. 


## Usage
```bash
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
