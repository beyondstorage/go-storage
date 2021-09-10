- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-09-08
- RFC PR: [beyondstorage/go-storage#749](https://github.com/beyondstorage/go-storage/pull/749)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP-749: Unify Path Behavior

Previous Discussion:

- [Specify the behavior of Path](https://forum.beyondstorage.io/t/topic/195)
- [OS-specific path separator problem](https://github.com/beyondstorage/go-service-ftp/issues/22)

## Background

Unix and Windows all use a different syntax for filesystem paths. Normally, when specifying file paths on Windows, we would use backslashes(`\`). However, in Java, and many other places outside the Windows world, backslashes are the escape character, so we have to double them up, like `C:\\trash\\blah\\blah`. Forward slashes(`/`), on the other hand, do not need to be doubled up and work on both Windows and Unix. There is no harm in having double forward slashes. They do nothing to the path and just take up space (`//` is equivalent to `/./`).

In our services, we have to handle paths like path splicing and path conversion. And here comes the problems:

- Different path separators may exist in a path at the same time, is it possible or how to use cross-platform paths for users?
- For fs storage service on Windows platform, the behavior of `work_dir` is undefined. Drive letters in connection string cannot be handled for Windows at this stage.

## Proposal

### Absolute Path and Relative Path

All our services should support two kinds of path:

- Absolute Path:
  - For Unix, the absolute path starts with `/`.
  - For Windows, the absolute path starts with a drive letter.
- Relative Path: 
  - The relative path is for the working directory `WorkDir`.
  - Path with prefix `./` or `../` is allowed: they are relative paths for file system, but common strings for object storage service.

### WorkDir and path

`WorkDir` specifies the working directory of the process.

- For file system, `WorkDir` SHOULD be an absolute path.
  - `WorkDir` MUST start with `/` for Unix. For Windows, the passed in `WorkDir` with drive letter SHOULD be allowed.
  - The default value is `/`, that means the working directory is the root path for Unix, and the current volume of the process for Windows.
  - Services SHOULD set `WorkDir` to the path name after the evaluation of any symbolic links internal.
- For object storage, `WorkDir` is the simulated directory or prefix of the object key.
  - `WorkDir` SHOULD be Unix style with the prefix `/`.
  - The default value is `/`.
  - Form service side, the prefix `/` needs to be removed for internal processing, but `/` needs to be added to the fields returned to the user.
  
`path` is the file path for file system, or an object key for object storage. Also, it could be a prefix filter for `List` operation.

- For file system services on Windows, `path` MUST be the relative path for `WorkDir`. For file system on Unix, `path` could be an absolute path or a relative path.
- For object storage service, `path` MUST be Unix style.
- For the unique key `Object.ID` in storage, it SHOULD be an absolute path compatible with the target platform, unless there's a returned unique identifier like in dropbox.
- For the path `Object.Path` in storage, it SHOULD be Unix style.
- Users SHOULD follow the file and object naming of different services.

### Path Separator

Input path for users that contains system-related path separator SHOULD be allowed. When the drive letter is included for Windows, it should be something like `c:\\a\\b`.

From go-storage side:

- Replace each separator character in path with a slash (`/`) character could be generated for the input path.

From service side:

- Services SHOULD replace `/` in path to the system-related path separator at the beginning of the operation.

## Rationale

### filepath

Package `filepath` implements utility routines for manipulating filename paths in a way compatible with the target operating system-defined file paths. `FromSlash` returns the result of replacing each slash (`/`) character in path with a separator character.

### Object key

- For [object key in s3](https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-keys.html#object-key-guidelines)
  - Forward slash (`/`) is the safe character, and backslash (`\`) need to be avoided. There is no hierarchy of subfolders. However, users can infer logical hierarchy using key name prefixes and delimiters (`/`). 
  - Objects with a prefix of `./` or `../`, or with key names ending with period(s) (`.`) are allowed but should be aware of the prefix limitations.
- For [object key in oss](https://www.alibabacloud.com/help/doc-detail/87728.htm)
  - The name cannot start with a forward slash (`/`) or a backslash (`\`).

## Compatibility

- `work_dir` must start with `/` is only for Unix.
- The format for `name` and `work_dir` in connection string need to be reconsidered as we can't separate them by `/` anymore.

## Implementation

N/A
