# Jirb
Jirb is a simple tool that allows you to create a git branch from a Jira issue without messing up the ticket number in your branch name.

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
_jirb_ stores your settings locally so that following uses are faster. You will not need to enter your credentials or Jira base url. 

## Preferences
Preferences currently only include the base url of your Jira server. They are stored in the `$HOME/.jirb` file.

## Credentials storage
Your ceredentials are store inside the default keychain of your OS, and never anywhere else.

## Jira Authentication
Requests to your Jira server are authenticated using basic auth.

## How are issues selected ?
_jirb_ lists all the issues that are assigned to you and in progress.