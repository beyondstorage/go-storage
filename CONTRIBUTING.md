# Contributing to Storage

+1tada First off, thanks for taking the time to contribute! tada+1

## Did you find a bug?

- Ensure the bug was not already reported by searching on GitHub under [Issues](https://github.com/Xuanwo/storage/issues).
- Open a new issue with following things: 
  - bug description
  - lib commit id
  - minimal reproduction code

## Did you write a patch that fixes a bug?

- Open a new GitHub pull request with the patch.
- Ensure the PR description clearly describes the problem and solution. Include the relevant issue number if applicable.
- Add unittest for this bug.

## Do you intend to implement a new service?

- `Storager` must be implemented, others can be optional.
- Add support in `coreutils.Open`.
- Add unittests as best effort.

## Do you intend to change public API?

- Open a new Github Issue for discuss.
- After achieve consensus, add a proposal in `docs/design` and submit a PR.
- Implement a proposal and change status to `candidate`

> In next release, relevant proposal statue will be updated to `finished`
