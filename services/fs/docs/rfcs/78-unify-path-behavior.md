- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-09-14
- RFC PR: [beyondstorage/go-service-fs#78](https://github.com/beyondstorage/go-service-fs/issues/78)
- Tracking Issue: [beyondstorage/go-service-fs#83](https://github.com/beyondstorage/go-storage/issues/83)

# RFC-78: Unify Path Behavior

Previous Discussion:

- [Specify the behavior of Path](https://forum.beyondstorage.io/t/topic/195)

## Background

A path is a string of characters used to uniquely identify a location in a directory structure.

- Windows supports multiple root names. A file system consists of a forest of trees, each with its own root directory, such as `c:\` or `\\network_name\`, and each with its own current directory. POSIX supports a single tree, with no root name, the single root directory /, and a single current directory.
- The directory hierarchy is separated by delimiters. The delimiting character is most commonly the slash (`/`), the backslash (`\`), or colon (`:`), though some operating systems may use a different delimiter.

Currently, the undefined behavior for `work_dir` and path makes our service unable to preserve compatibility for platforms and familiarity for users.

## Proposal

### Absolute path and relative path

Service SHOULD support two kinds of path:

- Absolute path:
    - For Unix, the absolute path starts with `/`.
    - For Windows, the absolute path starts with a drive letter, or the directory separator to represent an absolute path from the root of the current drive.
- Relative path:
    - The relative path is based on the working directory.

### Directory separator

From user side:

- System-related directory separator SHOULD be allowed. 
  - When the drive letter is included for Windows, it should be something like `c:\\a\\b`.

From service side:

- Service SHOULD translate passed in paths into platform-specific paths at the beginning of operations.

### Implementation

**work_dir**

The pair `work_dir` specifies the working directory of the process.

From service side:

- `work_dir` SHOULD be an absolute path.
  - For Unix, `work_dir` MUST start with `/`.
  - For Windows, `work_dir` with prefix `/` means an absolute path from the root of the current drive. And `work_dir` starts with drive letter SHOULD be allowed.
- The default value is `/`.
- Service SHOULD set `work_dir` to the path name after the evaluation of any symbolic links internal.

**path**

The field `path` is the file or directory path for operations.

From service side:

- For Unix, `path` could be an absolute path or a relative path based on `work_dir`.
- For Windows, `path` could be a relative path based on `work_dir` or an absolute path from the root of the current drive like `\Program Files\Custom Utilities\StringFinder.exe`. But absolute path with driver letter SHOULD NOT be allowed.
  
From user side:

- Users SHOULD follow the file naming of file system.

**Object:{Path, ID}**

`Object` will be returned by creating/stat/retrieve a file or object to identify the operation object.

- For the unique key `Object.ID` in storage, it SHOULD be an absolute path compatible with the target platform.
  - For Unix, it should be like `/path/to/workdir/hello.txt`.
  - For Windows, it should be like `c:\\path\\to\\workdir\\hello.txt`.
- For the file path `Object.Path` in storage, it SHOULD always be separated by slash (`/`).

## Rationale

### Unix style path name

The following example shows a Unix-style path:

```txt
/users/mark/
```

Mac OS X, as a derivative of UNIX, uses UNIX paths internally.

### MS-DOS/Microsoft Windows style

Contrary to popular belief, the Windows system API accepts slash, and thus all the Unix style path should work. But many applications on Windows interpret a slash for other purposes or treat it as an invalid character, and thus require users to enter backslash.

In addition, `\` does not indicate a single root, but instead the root of the "current disk". Indicating a file on a disk other than the current one requires prefixing a drive letter and colon.

The following examples show MS-DOS/Windows-style paths, with backslashes used to match the most common syntax:

```txt
A:\Temp\File.txt
```

## Compatibility

The format for `work_dir` in connection string need to be reconsidered as we can't separate it with `name` by `/` anymore.

## Implementation

Support driver letter and directory separator on Windows.
