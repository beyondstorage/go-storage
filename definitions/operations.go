package definitions

import (
	"sort"
)

type Operation struct {
	Name        string
	Namespace   string
	Local       bool
	Params      []Field
	Results     []Field
	Pairs       []Pair
	Description string
}

func SortOperations(ops []Operation) []Operation {
	sort.SliceStable(ops, func(i, j int) bool {
		return ops[i].Name < ops[j].Name
	})
	return ops
}

var OperationsServiceMap = make(map[string]Operation)

var OperationsService = []Operation{
	{
		Name:      "create",
		Namespace: NamespaceService,
		Params: []Field{
			getField("name"),
		},
		Results: []Field{
			getField("store"),
		},
		Description: "will create a new storager instance.",
	},
	{
		Name:      "delete",
		Namespace: NamespaceService,
		Params: []Field{
			getField("name"),
		},
		Description: "will delete a storager instance.",
	},
	{
		Name:      "get",
		Namespace: NamespaceService,
		Params: []Field{
			getField("name"),
		},
		Results: []Field{
			getField("store"),
		},
		Description: "will get a valid storager instance for service.",
	},
	{
		Name:      "list",
		Namespace: NamespaceService,
		Results: []Field{
			getField("sti"),
		},
		Description: "will list all storager instances under this service.",
	},
}

var OperationsStorageMap = make(map[string]Operation)

var OperationsStorage = []Operation{
	// Basic operations.
	{
		Name:      "create",
		Namespace: NamespaceStorage,
		Local:     true,
		Params: []Field{
			getField("path"),
		},
		Results: []Field{
			getField("o"),
		},
		Pairs: []Pair{
			PairObjectMode,
		},
		Description: `will create a new object without any api call.

## Behavior

- Create SHOULD NOT send any API call.
- Create SHOULD accept ObjectMode pair as object mode.`,
	},
	{
		Name:      "delete",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("path"),
		},
		Pairs: []Pair{
			PairObjectMode,
		},
		Description: `will delete an object from service.

## Behavior

- Delete only delete one and only one object.
  - Service DON'T NEED to support remove all.
  - User NEED to implement remove_all by themself.
- Delete is idempotent.
  - Successful delete always return nil error.
  - Delete SHOULD never return ObjectNotExist
  - Delete DON'T NEED to check the object exist or not.`,
	},
	{
		Name:      "metadata",
		Namespace: NamespaceStorage,
		Local:     true,
		Results: []Field{
			getField("meta"),
		},
		Description: `will return current storager metadata.`,
	},
	{
		Name:      "list",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("path"),
		},
		Results: []Field{
			getField("oi"),
		},
		Pairs: []Pair{
			PairListMode,
		},
		Description: `will return list a specific path.

## Behavior

- Service SHOULD support default ListMode.
- Service SHOULD implement ListModeDir without the check for VirtualDir.
- Service DON'T NEED to Stat while in List.`,
	},
	{
		Name:      "read",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("path"),
			getField("w"),
		},
		Results: []Field{
			getField("n"),
		},
		Pairs: []Pair{
			PairSize,
			PairOffset,
			PairIoCallback,
		},
		Description: `will read the file's data.`,
	},
	{
		Name:      "stat",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("path"),
		},
		Results: []Field{
			getField("o"),
		},
		Pairs: []Pair{
			PairObjectMode,
		},
		Description: `will stat a path to get info of an object.

## Behavior

- Stat SHOULD accept ObjectMode pair as hints.
  - Service COULD have different implementations for different object mode.
  - Service SHOULD check if returning ObjectMode is match`,
	},
	{
		Name:      "write",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("path"),
			getField("r"),
			getField("size"),
		},
		Results: []Field{
			getField("n"),
		},
		Pairs: []Pair{
			PairIoCallback,
		},
		Description: `will write data into a file.

## Behavior

- Write SHOULD support users pass nil io.Reader.
  - Service that has native support for pass nil io.Reader doesn't NEED to check the io.Reader is nil or not.
  - Service that doesn't have native support for pass nil io.Reader SHOULD check and create an empty io.Reader if it is nil.
- Write SHOULD NOT return an error as the object exist.
  - Service that has native support for overwrite doesn't NEED to check the object exists or not.
  - Service that doesn't have native support for overwrite SHOULD check and delete the object if exists.
- A successful write operation SHOULD be complete, which means the object's content and metadata should be the same as specified in write request.
`,
	},

	// Append related operations.
	{
		Name:      "create_append",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("path"),
		},
		Results: []Field{
			getField("o"),
		},
		Pairs: []Pair{},
		Description: `will create an append object.

## Behavior

- CreateAppend SHOULD create an appendable object with position 0 and size 0.
- CreateAppend SHOULD NOT return an error as the object exist.
  - Service SHOULD check and delete the object if exists.`,
	},
	{
		Name:      "write_append",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("o"),
			getField("r"),
			getField("size"),
		},
		Results: []Field{
			getField("n"),
		},
		Description: `will append content to an append object.`,
	},
	{
		Name:      "commit_append",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("o"),
		},
		Description: `will commit and finish an append process.`,
	},

	// Block related operations.
	{
		Name:      "create_block",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("path"),
		},
		Results: []Field{
			getField("o"),
		},
		Description: `will create a new block object.

## Behavior
- CreateBlock SHOULD NOT return an error as the object exist.
  - Service that has native support for overwrite doesn't NEED to check the object exists or not.
  - Service that doesn't have native support for overwrite SHOULD check and delete the object if exists.`,
	},
	{
		Name:      "write_block",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("o"),
			getField("r"),
			getField("size"),
			getField("bid"),
		},
		Results: []Field{
			getField("n"),
		},
		Description: `will write content to a block.`,
	},
	{
		Name:      "combine_block",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("o"),
			getField("bids"),
		},
		Description: `will combine blocks into an object.`,
	},
	{
		Name:      "list_block",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("o"),
		},
		Results: []Field{
			getField("bi"),
		},
		Description: `will list blocks belong to this object.`,
	},

	// Copy operations
	{
		Name:      "copy",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("src"),
			getField("dst"),
		},
		Description: `will copy an Object or multiple object in the service.

## Behavior

- Copy only copy one and only one object.
  - Service DON'T NEED to support copy a non-empty directory or copy files recursively.
  - User NEED to implement copy a non-empty directory and copy recursively by themself.
  - Copy a file to a directory SHOULD return ErrObjectModeInvalid.
- Copy SHOULD NOT return an error as dst object exists.
  - Service that has native support for overwrite doesn't NEED to check the dst object exists or not.
  - Service that doesn't have native support for overwrite SHOULD check and delete the dst object if exists.
- A successful copy opration should be complete, which means the dst object's content and metadata should be the same as src object.`,
	},

	// Dir related operations
	{
		Name:      "create_dir",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("path"),
		},
		Results: []Field{
			getField("o"),
		},
		Description: `will create a new dir object`,
	},

	// Fetch related operations.
	{
		Name:      "fetch",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("path"),
			getField("url"),
		},
		Description: `will fetch from a given url to path.

## Behavior

- Fetch SHOULD NOT return an error as the object exists.
- A successful fetch operation should be complete, which means the object's content and metadata should be the same as requiring from the url.`,
	},

	// Link related operations
	{
		Name:      "create_link",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("path"),
			getField("target"),
		},
		Results: []Field{
			getField("o"),
		},
		Description: `Will create a link object.

# Behavior

- path and target COULD be relative or absolute path.
- If target not exists, CreateLink will still create a link object to path.
- If path exists:
  - If path is a symlink object, CreateLink will remove the symlink object and create a new link object to path.
  - If path is not a symlink object, CreateLink will return an ErrObjectModeInvalid error when the service does not support overwrite.
- A link object COULD be returned in Stat or List.
- CreateLink COULD implement virtual_link feature when service without native support.
  - Users SHOULD enable this feature by themselves.`,
	},

	// Move related operations
	{
		Name:      "move",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("src"),
			getField("dst"),
		},
		Description: `will move an object in the service.

## Behavior

- Move only move one and only one object.
  - Service DON'T NEED to support move a non-empty directory.
  - User NEED to implement move a non-empty directory by themself.
  - Move a file to a directory SHOULD return ErrObjectModeInvalid.
- Move SHOULD NOT return an error as dst object exists.
  - Service that has native support for overwrite doesn't NEED to check the dst object exists or not.
  - Service that doesn't have native support for overwrite SHOULD check and delete the dst object if exists.
- A successful move operation SHOULD be complete, which means the dst object's content and metadata should be the same as src object.`,
	},

	// Multipart related operations.
	{
		Name:      "create_multipart",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("path"),
		},
		Results: []Field{
			getField("o"),
		},
		Description: `will create a new multipart.

## Behavior

- CreateMultipart SHOULD NOT return an error as the object exists.`,
	},
	{
		Name:      "write_multipart",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("o"),
			getField("r"),
			getField("size"),
			getField("index"),
		},
		Results: []Field{
			getField("n"),
			getField("part"),
		},
		Description: `will write content to a multipart.`,
	},
	{
		Name:      "complete_multipart",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("o"),
			getField("parts"),
		},
		Description: `will complete a multipart upload and construct an Object.`,
	},
	{
		Name:      "list_multipart",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("o"),
		},
		Results: []Field{
			getField("pi"),
		},
		Description: `will list parts belong to this multipart.`,
	},

	// Page related operations
	{
		Name:      "create_page",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("path"),
		},
		Results: []Field{
			getField("o"),
		},
		Description: `will create a new page object.

## Behavior

- CreatePage SHOULD NOT return an error as the object exists.`,
	},
	{
		Name:      "write_page",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("o"),
			getField("r"),
			getField("size"),
			getField("offset"),
		},
		Results: []Field{
			getField("n"),
		},
		Description: `will write content to specific offset.`,
	},

	// Query sign HTTP related operations.
	{
		Name:      "query_sign_http_read",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("path"),
			getField("expire"),
		},
		Results: []Field{
			getField("req"),
		},
		Description: "will read data from the file by using query parameters to authenticate requests.",
	},
	{
		Name:      "query_sign_http_write",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("path"),
			getField("size"),
			getField("expire"),
		},
		Results: []Field{
			getField("req"),
		},
		Description: "will write data into a file by using query parameters to authenticate requests.",
	},
	{
		Name:      "query_sign_http_delete",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("path"),
			getField("expire"),
		},
		Results: []Field{
			getField("req"),
		},
		Description: "will delete an object from service by using query parameters to authenticate requests.",
	},
	{
		Name:      "query_sign_http_create_multipart",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("path"),
			getField("expire"),
		},
		Results: []Field{
			getField("req"),
		},
		Description: "will create a new multipart by using query parameters to authenticate requests.",
	},
	{
		Name:      "query_sign_http_write_multipart",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("o"),
			getField("size"),
			getField("index"),
			getField("expire"),
		},
		Results: []Field{
			getField("req"),
		},
		Description: "will write content to a multipart by using query parameters to authenticate requests.",
	},
	{
		Name:      "query_sign_http_list_multipart",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("o"),
			getField("expire"),
		},
		Results: []Field{
			getField("req"),
		},
		Description: "will list parts belong to this multipart by using query parameters to authenticate requests.",
	},
	{
		Name:      "query_sign_http_complete_multipart",
		Namespace: NamespaceStorage,
		Params: []Field{
			getField("o"),
			getField("parts"),
			getField("expire"),
		},
		Results: []Field{
			getField("req"),
		},
		Description: "will complete a multipart upload and construct an Object by using query parameters to authenticate requests.",
	},
}

func init() {
	for idx, ops := range [][]Operation{OperationsService, OperationsStorage} {
		for k := range ops {
			ops[k].Params = append(ops[k].Params, getField("pairs"))
			if !ops[k].Local {
				ops[k].Results = append(ops[k].Results, getField("err"))
			}

			if idx == 0 {
				OperationsServiceMap[ops[k].Name] = ops[k]
			} else {
				OperationsStorageMap[ops[k].Name] = ops[k]
			}
		}
	}
}
