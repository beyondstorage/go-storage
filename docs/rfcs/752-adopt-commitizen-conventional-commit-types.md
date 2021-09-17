- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-09-10
- RFC PR: [beyondstorage/go-storage#752](https://github.com/beyondstorage/go-storage/pull/752)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP-752: Adopt commitizen/conventional-commit-types

Previous Discussion:

- [Adopt commitizen/conventional-commit-types](https://forum.beyondstorage.io/t/topic/193)

## Background

Generally, the commit message should be clear and unambiguous, stating the purpose of this commit, what specific operations have been done... But in daily development, fix bug and other kinds of generalized messages are commonplace. And sometimes we don't even know what problem we are modifying with the commit, or the impact of a particular commit. This leads to the cost of subsequent code maintenance.

Angular specification is currently the most widely used writing style, more reasonable and systematic, and has supporting tools.

## Proposal

`commitizen/conventional-commit-types` should be adopted in our commits.

Each commit message consists of a **header**, a **body**, and a **footer**.

```txt
<header>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

`header` and `body` are mandatory, `footer` is optional.

### Commit Message Header

```txt
<type>(<scope>): <short summary>
  │       │             │
  │       │             └─⫸ Summary in present tense. Not capitalized. No period at the end.
  │       │
  │       └─⫸ Commit Scope: cmd|definitions|internal|pairs|pkg|services|tests|types
  │
  └─⫸ Commit Type: build|ci|docs|feat|fix|perf|refactor|test
```

The `<type>` and `<summary>` fields are mandatory, the `(<scope>)` field is optional.

**Type** 

Must be one of the following:

- build: Changes that affect the build system or external dependencies
- ci: Changes to our CI configuration files and scripts
- chore: Changes to the build process or supporting tools
- docs: Documentation only changes
- feat: A new feature
- fix: A bug fix
- perf: A code change that improves performance
- refactor: A code change that neither fixes a bug nor adds a feature
- revert: Reverts a previous commit
- style：Changes that do not affect the meaning of the code
- test: Adding missing tests or correcting existing tests

**Scope**

The scope should be the name of the package affected in go-storage.

The following is the current list of supported scopes:

- cmd
- definitions
- internal
- pairs
- pkg
- services
- tests
- types

`*` can be used if the modification affects more than one scope.

There are currently a few exceptions to the "use package name" rule:

- packaging: used for changes that change the package layout in all of our packages
- changelog: used for updating the release notes in CHANGELOG.md
- dev-infra: used for dev-infra related changes
- docs-infra: used for docs related changes within the /docs directory
- migrations: used for changes to the update migrations

**Summary** 

Use the summary field to provide a succinct description of the change:

- use the imperative, present tense: "change" not "changed" nor "changes"
- don't capitalize the first letter
- no dot (.) at the end

### Commit Message Body

Just as in the summary, use the imperative, present tense: "fix" not "fixed" nor "fixes".

Explain the motivation for the change in the commit message body. This commit message should explain why we are making the change.

### Commit Message Footer

The footer can contain information about breaking changes and is also the place to reference GitHub issues, other PRs that this commit closes or is related to.

#### Breaking changes

Breaking change section should start with the phrase "BREAKING CHANGE: " followed by a summary of the breaking change, a blank line, and a detailed description of the breaking change that also includes migration instructions.

```txt
BREAKING CHANGE: <breaking change summary>
<BLANK LINE>
<breaking change description + migration instructions>
<BLANK LINE>
<BLANK LINE>
Fixes #<issue number>
```

#### Close Issue

If the current commit is for an issue, then we can close the issue in the `Footer` section:

```txt
Closes #123
```

It is also possible to close multiple issues at once:

```txt
Closes #123, #456, #789
```

### Revert Commit

If the commit reverts a previous commit, it should begin with `revert: `, followed by the header of the reverted commit.

The content of the commit message body should contain:

- information about the SHA of the commit being reverted in the following format: `This reverts commit <SHA>`,
- a clear description of the reason for reverting the commit message.

### Implementation

#### Commitizen

[Commitizen](http://commitizen.github.io/cz-cli/) is a tool for writing qualified commit messages.

We can make our repo Commitizen-friendly. Then we can simply use `git cz` or just `cz` instead of `git commit` when committing, the option appears to generate a formatted commit message.

#### commitlint

[commitlint](https://github.com/conventional-changelog/commitlint) is used to check the format of commit message. We can set rules to check the commit message.

#### Generate changelog

[conventional-changelog](https://github.com/conventional-changelog/conventional-changelog) is a tool for generating changelogs and release notes from a project's commit messages and metadata.

If all the commits are in Angular format, then the changelog can be automatically generated with a script when releasing a new version.

## Rationale

Commit types originally from: [Angular Git Commit Message Conventions](https://github.com/angular/angular/blob/master/CONTRIBUTING.md#type)

## Compatibility

New additional utility, no break change.

## Implementation

N/A