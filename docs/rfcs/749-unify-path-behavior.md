- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-09-08
- RFC PR: [beyondstorage/go-storage#749](https://github.com/beyondstorage/go-storage/pull/749)
- Tracking Issue: [beyondstorage/go-storage#758](https://github.com/beyondstorage/go-storage/issues/758)

# GSP-749: Unify Path Behavior

Previous Discussion:

- [Specify the behavior of Path](https://forum.beyondstorage.io/t/topic/195)
- [OS-specific path separator problem](https://github.com/beyondstorage/go-service-ftp/issues/22)

## Background

For file system, paths are used extensively to represent the directory/file relationships. Resources can be represented by either absolute or relative paths. And the directory separator is platform specific.

For object storage, the key of an object uniquely identifies the object in a bucket. Key name prefixes and delimiters could be used to group objects or simulate directory.

In our services, path could be a file path for file system, or a file hosting service, also it could be an object for object storage. So we need to unify path behavior for different platforms and storages.

## Proposal

### Absolute path and relative path

All our services should support two kinds of path:

- Absolute path: must include the root directory.
- Relative path: a relative path based on the working directory.

### Directory separator

From user side:

- System-related directory separator including `/` for Unix and `\\` for Windows SHOULD be allowed.

From go-storage side:

- go-storage SHOULD be able to tolerate mixed use of directory separator and replace each separator with a slash (`/`) character could be generated for the input path.

From service side:

- Services SHOULD replace separator in the passed-in path to the current system-related directory separator at the beginning of operations.
  
### Implementation

**work_dir**

The pair `work_dir` is used to specific the working directory for service or storage. For object storage, it's the simulated directory or prefix of the object key.

- `work_dir` SHOULD be an absolute path.
- `work_dir` SHOULD be default to `/` if not set.
- `work_dir` SHOULD be Unix style path for object storage services.

Take the example of setting working directory: 

```go
// users could set work_dir for init
store, _ := s3.NewStorager(
	// ...
	pairs.WithWorkDir("/path/to/workdir"),
)

// services should set default value for `work_dir`
store := &Storage{
	// ... 
	workDir: "/",
}
if opt.HasWorkDir {
	store.workDir = opt.WorkDir
}
```

**path**

The field `path` is used as an input parameter for operation, indicating the file path for file system, object key for object storage or prefix filter for `List` operator.

- `path` could be absolut path or relative path based on `work_dir`.
- `path` SHOULD be consistent with the path style of `work_dir`.

Take the example of calling `Read`:

```go
// absolute path for Read 
n, err := store.Read("/path/to/workdir/hello.txt", w)

// relative path for Read 
n. err := store.Read("hello.txt", w)
```

**Object:{Path, ID}**

`Object` is the smallest unit in go-storage, returned by creating/stat/retrieve a file or object to identify the operation object, and the fields should not be changed outside services. `Object.ID` is the unique key in storage, and `Object.Path` is either the absolute path or the relative path based on the working directory depends on user's input.

- `Object.ID` SHOULD be an absolute path compatible with the target platform.
  - If a service like Dropbox returns a unique identifier, then `Object.ID` will be the returned identifier.
  - For object storage services, the prefix `/` of the full path needs to be removed for internal processing, and added back to `Object.ID`.
- `Object.Path` SHOULD be Unix style.

Take s3 as an example:

```go
// get absolute path for object
func (s *Storage) getAbsPath(path string) string {
	// remove prefix `/` 
	prefix := strings.TrimPrefix(s.workDir, "/")
	return prefix + path
}

// return Object for stat
func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {
	// use absolute path for request 
	rp := s.getAbsPath(path)
	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.name), 
		Key:    aws.String(rp),
	}
	// ... 
	o = s.newObject(true)
	o.ID = rp
	o.Path = path 
	// ...
}
```

## Rationale

### Object key

- For [object key in s3](https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-keys.html#object-key-guidelines)
  - Forward slash (`/`) is the safe character, and backslash (`\`) need to be avoided. There is no hierarchy of subfolders. However, users can infer logical hierarchy using key name prefixes and delimiters (`/`). 
  - Objects with a prefix of `./` or `../`, or with key names ending with period(s) (`.`) are allowed but should be aware of the prefix limitations.
- For [object key in oss](https://www.alibabacloud.com/help/doc-detail/87728.htm)
  - The name cannot start with a forward slash (`/`) or a backslash (`\`).
  
### How to be compatible with file system on Windows?

Exception rules could be described in the RFC of the corresponding service. Path behavior for local file system is defined [here](https://github.com/beyondstorage/go-service-fs/pull/78).

## Compatibility

This change will not break services and users.

- Services should implement or update the path behavior.
- `work_dir` in connection string need to be reconsidered.

## Implementation

- Update descriptions for `work_dir` and `ID`, `Path` of `Object`.
- Generate code that replaces separator with `/` in go-storage.
- Add test cases for path:
  - Basic operation with absolute path.
  - Basic operation with path using different separator.
  - Add `Object.ID` check for the current tests.
- Update path behavior in services.
