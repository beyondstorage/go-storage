- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-04-23
- RFC PR: [beyondstorage/specs#41](https://github.com/beyondstorage/specs/issues/41)
- Tracking Issue: N/A

# AOS-41: Turn Pair Expire into Duration

## Background

Pair `expire` is `int` for now. It's hard to understand what's it real meaning.

Is it the deadline of this URL? Or is it the duration of this URL, in this situation, the deadline should be `time.Now() + expire`. 

## Proposal

So I propose to turn `expire` into `Duration`.

For service that accept a duration, we can pass `expire` directly.
For service that accept a deadline, we can pass `time.Now() + expire`.

## Rationale

We choose to use `Duration` because it's simple and well-defined in most languages. At the same time, `Datetime` is much more complicated than it, we need to handle time zone, offset and others.

## Compatibility

Only `qingstor` implement `Reach` and used `expire` and no user for now, it's OK to change the type.

## Implementation

The proposal will be implemented by proposer.
