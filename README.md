# Jirb
Jirb is a simple tool that allows you to create a git branch from a Jira issue without messing up the ticket number in your branch name.

Currently only useful if your branch names follow a syntax similar to `feature/ABC-123_branch-description` or `hotfix/ABC-123_branch-description`

## First Use

1. Move to a git project you want to work on
1. Run `jirb`
1. Enter the base URL of your jira server (e.g. `https://example.net/jira` or `https://jira.example.net`)
1. Enter your Jira username
1. Enter your password
1. _jirb_ will show a prompt with issues your are working on. Pick the one you want to create a branch for.
1. Choose the branch prefix (`feature/` or `hotfix/`)
1. _jirb_ will prefil a branch name in your prompt. Modify it as you prefer.
1. _jirb_ will create a branch using the branch name you entered
    * _jirb_ does so with the command `git checkout -b <name>`

## Following uses
_jirb_ stores your settings locally. Following uses will not required you to enter your credentials or Jira URL. 

## Options
Option | Description
--- | ---
`-config` | Cycle through all configurations and optionally change values.
`-reset` | Remove all configurations files and keychain entries making it as if you had never run this tool.
`-help` | Prints a usage message describing these options.


## Preferences
Preferences currently only include the base url of your Jira server. They are stored in the `$HOME/.jirb` file.

## Credentials storage
Your credentials are stored inside the default keychain and loaded in memory during execution.

## Jira Authentication
Requests to your Jira server are authenticated using basic auth.

## How are issues selected ?
_jirb_ lists all the issues that are assigned to you and currently in progress.