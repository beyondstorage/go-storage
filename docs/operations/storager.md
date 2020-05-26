# Storager

Storager is the interface for storage service.

Everything in a storager is an Object with two types: File, Dir.

File is the smallest unit in service, it will have content and metadata. Dir is a container for File and Dir.
In prefix-based storage service, Dir is usually an empty key end with "/" or with special content type.
For directory-based service, Dir will be corresponded to the real directory on file system.

In the comments of every method, we will use following rules to standardize the Storager's behavior:

- The keywords "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY",
and "OPTIONAL" in this document are to be interpreted as described in RFC 2119.
- Implementer is the provider of the service, while trying to implement Storager interface, you need to follow.
- Caller is the user of the service, while trying to use the Storager interface, you need to follow.

## Delete

Delete will delete an Object from service.

## Metadata

Metadata will return current storager's metadata.

Implementer:

- Metadata SHOULD only return static data without API call or with a cache.

Caller:

- Metadata SHOULD be cheap.

## Read

Read will read the file's data.

Caller:

- MUST close reader while error happened or all data read.

## Stat

Stat will stat a path to get info of an object.

## Write

Write will write data into a file.

Caller:

- MUST close reader while error happened or all data written.
