# git-owners

**git-owners** is a system (inspired by the [Chromium OWNERS system](https://chromium.googlesource.com/chromium/src/+/master/docs/code_reviews.md#OWNERS-files))
to assign ownership to directories in a repository for the purposes of checking that a PR has been reviewed by sufficiently knowledgeable reviewers.

This project defines an OWNERS file format and includes tooling for finding reviewers for a set of files and given a set of files and acceptances verifying that all ownership rules are satisfied.

## OWNERS file format

An example OWNERS file:
```
# flags can be set to change resolving behaviour
@set noparent

# alice can approve any change in this directory
alice

# a review from both bob and carol together is sufficient to approve a change to a go file in this directory
(bob & carol) *.go

# BUILD files do not need a reviewer
@ignore BUILD

# dave and eve are not required to review changes but should be notified
@watchers (dave, eve@example.com)

@watchers fred *.go
```
Full specification of available directives is given below.

## CLI

### Find owners for a file
`git owners [-a] pathToFile`

Prints out the owners of a file, one per line.
If `-a` is given then all owners will be resolved up to the root of the repo else the resolver will stop after the first OWNERS file.

### Find reviewers for a PR
`git owners [--base-branch <name=master>]`

Finds the commit on base-branch (default master) that this was branched from and gets the list of files changed since that commit.
Outputs a list of reviewers that satisfies the property:
> For every modified file there is at least one owner (or group of owners) in the list of reviewers that is an owner for that file.

This will attempt to find a minimal set of reviewers (i.e. minimise the amount of redundant reviews where two owners of a file review the same file) however this is only best effort.
This process is non-deterministic so suggestions will be load balanced between owners.

## Specification

#### Flags
A flag is of the form `@set flagname` and sets the given flag name to true.
Any string is allowed for a flag name however this implementation only recognises the `noparent` flag.

#### Owners
An owner directive is of the form of a user specification followed by an optional filename pattern, separated by whitespace.

A user specification is either a single owner or a bracketed group of owners separated by `&`.
A user may be any of a username, email, or other identifier.

A filename pattern is any valid golang match [pattern](https://golang.org/pkg/path/filepath/#Match).

#### Ignores
An ignore directive is of the form `@ignore` followed by a filename pattern.

#### Watchers
A watchers directive is of the form `@watchers` followed by a valid watchers specification and an optional filename pattern.

A watchers specification is either a single user or a bracketed, comma separated list of users.

### Resolution algorithm

Given a path to a file in a repo the set of owners is constructed as follow:
1. If there is not and OWNERS file in the current directory move to the parent and GOTO 1.
2. For each owner/group of owners in the OWNERS file, if they match the given file then add them to the set of owners.
3. If the `noparent` flag is not set then move to the parent directory and GOTO 1.

The resulting set is all of the users/groups of users who can approve a change to this file.
