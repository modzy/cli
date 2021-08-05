# cli

## Installation

    go get -u github.com/modzy/cli/modzy

## Documentation

There is documentation for all commands within the cli:

    modzy --help
    modzy models --help
    modzy models get --help
    etc.

The top level `modzy --help` has additional high level documentation on how to use the cli.

## Temporary Notes
### Requirements

0. General Requirements:
    - If applicable, commands should all support the `--profile` flag to denote which Modzy installation or team to run the command under
    - Where applicable, data should be returned to std out in a tabular format 
    - Commands should support a `-v` or `-verbose` flag that would return the full api response 
        - An additional flag should be supported to specify output in JSON or YAML
    - Commands with larger responses (describe job details or describing a model) should output in a manner similar to `kubectl describe pod <pod_name>` 


1. CLI Configuration  
    - Ability to configure CLI with modzy installation properties
        - Modzy URL and API Key
        - Support of multiple configurations via profiles
    - Support for profiles will require the global `--profile` flag 
        - modzy list models vs. modzy --profile dev list models
    - Ability to run "whoami" command to view specified profile 
2. Models 
    - List the models available in the specified Modzy installation
        - Should include model ID, model name, and available versions
        - Should be able to specify author by flag
        - Should be able to specify number of results by flag (default 10)
        - Sort by ASC or DESC via flag
    - List model details 
        - Specify model ID as an arg
        - "images" array not to be displayed
3. Jobs
    - Cancel a job
        - If my API key allows for it, I'd like to specify a job ID and have that job be cancelled
    - View Job Status
        - Specify a job ID as an arg 
        - View status of the specified job
    - View Job Details
        - Job ID as arg
        - View details as returned by the API
    - List Jobs
        - Should include Job ID, submitted by, team the job belongs to (if possible), status
        - Should be able to filter by --status flag
        - Should be able to specify number of results by flag (default 10)
        - If possible, sort by date submitted via flag

### Potential deviations from requirements:

- `-v` is for verbose output, for debugging purpose.  "complete output" will be achieved when supplying a --output format that isn't the default (like json or yaml)
- instead of `modzy jobs status <job_id>` -- what about `modzy jobs get <job_id> --output jsonpath='{.status}'`?  This is a more generic method at grabbing any property from the response.
- I added `modzy jobs wait <job_id>`

### Command examples

    modzy --profile dev -v -o json ...
    modzy --api-key APIKEY ...
    MODZY_API_KEY=APIKEY modzy ...
    modzy ... whoami

    modzy models list --filter author=modzy --page 1+10 --sort somefield:asc
    modzy models get <model_id>
    
    modzy jobs list --filter status=TIMEDOUT --page 1+10 --sort date:desc
    modzy jobs get <job_id> --output jsonpath='{.status}'
    modzy jobs wait <job_id>
    modzy jobs cancel <job_id>

