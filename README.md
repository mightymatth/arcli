<img alt="arcli" src="img/arcli.png" height="200" width="200" />

[![Go Report Card](https://goreportcard.com/badge/github.com/mightymatth/arcli)](https://goreportcard.com/report/github.com/mightymatth/arcli)
[![arcli](https://snapcraft.io//arcli/badge.svg)](https://snapcraft.io/arcli)

### `arcli` - Awesome Redmine CLI
`arcli` is CLI for [Redmine](https://www.redmine.org/) that simplifies some actions such as checking for issue details and tracking time. 
It supports Redmine v3.3.1+ (tested with v3.3.1 and v4.1.0).

### Quick examples

Listing (ls) all assigned issues (i) to current user.
```
➜  ~ arcli i ls  
    ID  PROJECT           SUBJECT                   URL                                       
 20123  Webshop           Managing users            https://custom.url.com/issues/20123 
 20660  Webshop Android   Notification management   https://custom.url.com/issues/20460 
```

Log spent time (l) for issue (i) with ID 20123 with time duration (-t) of 1.5 hours.
```
➜  ~ arcli l i 20123 -t 1.5
Time entry created!
 ENTRY ID  PROJECT NAME  ISSUE ID  HOURS  ACTIVITY     COMMENT  SPENT ON        
 39458     Webshop       20123     1.5    programming           Thu, 2020-03-12 
```

Show tracking time status.
```
➜  ~ arcli status
[324] John Doe (john.doe@email.com)
PERIOD       HOURS   H/LOG   # of I   # of P  
Today        0       0       0        0       
Yesterday    0       0       0        0       
This Week    0       0       0        0       
Last Week    40      6.7     3        2       
This Month   40      6.7     3        2       
Last Month   160     5.5     8        6 
```


## Installation

### macOS

```
$ brew tap mightymatth/arcli https://github.com/mightymatth/arcli
$ brew install arcli
```

### Linux
[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/arcli)

```
snap install arcli
```

## Usage
```
➜  ~ arcli -h
Awesome Redmine CLI. Wrapper around Redmine API

Usage:
  arcli [flags]
  arcli [command]

Available Commands:
  aliases     Words that can be used instead of issue or project ids
  defaults    User session defaults
  help        Help about any command
  issues      Shows issue details
  log         Time entries on projects and issues
  login       Authenticate to Redmine server
  logout      Logout current user
  projects    Shows project details
  search      Search Redmine
  status      Overall account info
  view        Shows data in different views

Flags:
  -h, --help      help for arcli
  -v, --version   Current arcli and supported Redmine API version

Use "arcli [command] --help" for more information about a command.
```
