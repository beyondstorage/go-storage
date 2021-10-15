- Author: bokket  <bokkett@gmail.com>
- Start Date: 2021-07-17
- RFC PR: [beyondstorage/go-endpoint#8](https://github.com/beyondstorage/go-endpoint/pull/8)
- Tracking Issue: [beyondstorage/go-endpoint/issues/9](https://github.com/beyondstorage/go-endpoint/issues/9)

# RFC-8: Add TCP protocol

Releated issue: [beyondstorage/go-endpoint/issues/7](https://github.com/beyondstorage/go-endpoint/issues/7)

## Background

hdfs usually use the `New(address string)` method to access a namenode node, the user will be the user running the code. If the address is an empty string, it will try to get the NameNode address from the Hadoop configuration file.

## Proposal

I suggest adding a tcp protocol to allow the user to specify the address.

It likes `tcp:<host>:<port>`

- The `type` of `tcp` should be `String` and is a `const`
- The `value` of `endpoint` should be parsed into `ProtocolTCP`  and   `args include <host>:<port>`

## Rationale

Now we don't have a pair operation on the `hdfs address` or tcp-like operation

## Compatibility

No compatibility issues at this time.

## Implementation

- Add protocol `tcp`
- Implement protocol tcp formatted (`func (p Endpoint) TCP() (addr, host string, port int)`)
- Implement protocol tcp parser (`func Parse(cfg string) (p Endpoint, err error)`)
- Implement protocol tcp object (`func NewTCP(host string,port int) Endpoint `)
