
![Modzy Logo](https://www.modzy.com/wp-content/uploads/2020/06/MODZY-RGB-POS.png)
[![Go Report Card](https://goreportcard.com/badge/github.com/modzy/cli)](https://goreportcard.com/report/github.com/modzy/cli)

# Modzy Command Line Interface
Modzy's CLI provides terminal commands for some of our most useful API endpoints. Modzy's CLI is designed primarily for 
operations teams to monitor instance health, check on inference statuses, or to quickly query important operational 
information. As a result, not all of Modzy's API endpoints are available via the CLI.

For more detailed information, visit our [CLI Documentation page](https://docs.modzy.com/v1.0.6/docs/cli)
## Installation
As a prerequisite for installing the Modzy CLI, you must install have golang installed. Download and install golang 
here: https://golang.org/doc/install

Once you've successfully installed golang, open your terminal and run the following command

    go get -u github.com/modzy/cli/modzy

## Documentation

There exists documentation for all commands within the cli:

    modzy --help
    modzy models --help
    modzy models get --help
    etc.

The top level `modzy --help` has additional high level documentation on how to use the cli.

## Available Commands
| **Available Commands** | **Description** |
| --- | ---|
| [modzy alerts](https://docs.modzy.com/v1.0.6/docs/modzy-alerts) | List all alerts |
| [modzy completion](https://docs.modzy.com/v1.0.6/docs/modzy-completion) | Generate the autocompletion script for the specified shell ||
| [modzy help](https://docs.modzy.com/v1.0.6/docs/modzy-help) | Help about any command |
| [modzy jobs](https://docs.modzy.com/v1.0.6/docs/modzy-jobs) | Work with inference/preduction jobs |
| [modzy models](https://docs.modzy.com/v1.0.6/docs/modzy-models) | work with models |
| [modzy projects](https://docs.modzy.com/v1.0.6/docs/modzy-projects) | work with projects | 
| [modzy status](https://docs.modzy.com/v1.0.6/docs/status) | Returns top level dashboard information for your account | 
| [modzy whoami](https://docs.modzy.com/v1.0.6/docs/modzy-whoami) | Summarize effective authentication information | 

## Global Flags

|**Flag** | **Description** | 
|--- |  ---|
|`--api-key` | Modzy API key to use for authentication |
|`--base-url` | Modzy API base URL | 
|`-h`, `--help` | Help for Modzy CLI commands | 
|`-p`, `--profile` | profile under which command will be executed | 
|`--team-id` | Modzy API team ID to use for team authentication | 
|`--team-token` | Modzy API team token to use for team authentication | 
| `-v`, `--verbose` | Enable more verbose log output | 
| `--verbose-http` | Enable log output of http request and response data |

## Authentication
You can provide your authentication token through any mixture of command flags,
ENV variables, or configuration files.  Precedence is command flag > ENV > configuration file.

In all three examples below, you will need to replace placeholder values with valid inputs.
Replace ***`BASE_URL`***  with the URL of your instance of Modzy, such as https://app.modzy.com
Replace ***`API_KEY`*** with a valid API key string [Here's how to download an API Key from Modzy](doc:view-and-manage-api-keys)
Optionally, replace ***`TEAM_ID`*** with the ID for the Modzy team you'd like to access
Optionally, replace ***`TEAM_TOKEN`*** with a valid API key string

### Via Command Flags
To authenticate using command flags, you'll need to include `--base-url flag` with each command. You'll also need to include either the `--api-key` flag, or both the `--team-id` and `--team-token` flags.
```
$ modzy [command] --base-url BASE_URL {--api-key API_KEY | --team-id TEAM_ID --team-token TEAM_TOKEN}
```

### Via ENV Variables
Set the following ENV variables
MODZY_BASE_URL=***`BASE_URL`***
MODZY_API_KEY=***`API_KEY`***
MODZY_TEAM_ID=***`TEAM_ID`***
MODZY_TEAM_TOKEN=***`TEAM_TOKEN`***

### Via Configuration File
To use a configuration file, create a yaml file at any of these locations:
* `/etc/modzy/{profile}.yaml`
* `$HOME/.modzy/{profile}.yaml`

The default profile is called "default".

The yaml file you create should look something like this.
```
	> cat ~/.modzy/default.yaml
	base-url: BASE_URLf
	# use an api key:
	api-key: API_KEY
	# or use a team key:
	team-id: TEAM_ID
	team-token: TEAM_TOKEN
```
