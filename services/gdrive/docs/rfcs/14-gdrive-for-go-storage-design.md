- Author: Jun jun@junz.org
- Start Date: 2021-7-18
- RFC PR: [beyondstorage/go-service-gdrive#14](https://github.com/beyondstorage/go-service-gdrive/issues/14)
- Tracking Issue: [beyondstorage/go-service-gdrive#15](https://github.com/beyondstorage/go-service-gdrive/issues/15)

# RFC-14: Gdrive for go-storage design

## Background

Google drive API has so many different notions that differs from `go-storage`, and we have briefly discussed in [Gdrive use FileId to manipulate data instead of file name #11](https://github.com/beyondstorage/go-service-gdrive/issues/11). Now I would like to start a RFC so that we can make all things more clear.

In Google drive API, `FileID` is a critical attribute of a file(or directory). We will use it to manipulate data instead of by path. In fact, path is very trivial in gdrive, and we can create files with the same name in the same location. In other words, path can be duplicate in gdrive. This behavior can cause some problems to our path based API.

## Proposal

**We manually stipulate that every path is unique.**

 When users try to call `Write` to an existing file, we update it's content instead of creating another file with the same name.

**We will do a conversion between path and `FileID`.**

In this way, every path can be converted to `FileID`, so we are able to build a good bridge between `go-storage` API and gdrive API.

**We will cache `path -> id` in memory with TTL.**

For performance reasons, we will cache the ids of the files as they are created, and we will only look up their ids when the cache expires.

## Implementation

When users try to call `Write("foo/bar/test.txt")`, we will do this:

First, we look up the `FileID` of `foo` in cache, and try to search it's `FileID` if it is expired. Then, we will do the same thing to `bar` and `test.txt`.  Be aware that when we can not find the `FileId` of a directory, we won't continue to do the search to it's subdirectories. In this case, we can consider the file doesn't exist.

After that, there are two possibilities:

When `foo/bar/test.txt` doesn't exist, we will create folders one by one. At the same time, we will cache their `FileId`.

When `foo/bar/test.txt` already exist, then we will update it's content instead of creating another one.

Our significant point `pathToId` can be implement like this:

If the file is in the root folder, then we just do a simple search by using `drive.service.Files.List().Q(searchArgs).Do()`. The return value type is `*drive.File`, and it's attribute `ID` is what we need.

But if the file path is like `foo/bar/demo.txt`, it would be a little complex.

First, we get the `FileID` of directory `foo` like what we previously do, then we can use this `FileID` to list all of it's content. By this way, we can find a directory named `bar` and it's `FileID`. At last, we just repeat what we did before, and get the `FileID` we want. 

