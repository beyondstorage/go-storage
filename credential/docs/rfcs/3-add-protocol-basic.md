- Author: npofsi <npofsi@outlook.com>
- Start Date: 2021-08-02
- RFC PR: https://github.com/beyondstorage/go-credential/pull/3
- Tracking Issue: [beyondstorage/go-credential#1](https://github.com/beyondstorage/go-credential/issues/1)

# RFC-3: Add protocol basic


## Background

Some services don't support other credentials, but user or password.


## Proposal

Add protocol `basic`. Like `hmac`, `basic` have two parameters, corresponding to user and password of an account.

For example, go-service-ftp need a account to sign in, like

`ftp://xxx?credential=basic:user:password`

## Rationale

- Account is the only certification to some platform.
- Account is a basic method to identify quests.

## Compatibility

Will just add a choose to use protocol `basic`.

## Implementation

Just need to parse `basic` like `hmac`.