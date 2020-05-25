
interface "copier" {
  description = "is the interface for Copy.\n"

  op "copy" {
    description = "will copy an Object or multiple object in the service.\n"
    params      = ["src", "dst"]
  }
}
interface "dir_lister" {
  description = "is used for directory based storage service to list objects under a dir.\n"

  op "list_dir" {
    description = "will return list a specific dir.\n"
    params      = ["dir"]
  }
}
interface "dir_segments_lister" {
  description = "is used for directory based storage service to list segments under a dir.\n"
  embed       = ["segmenter"]

  op "list_dir_segments" {
    description = "will list segments via dir.\n"
    params      = ["dir"]
  }
}
interface "index_segmenter" {
  description = "is the interface for index based segment.\n"
  embed       = ["segmenter"]

  op "init_index_segment" {
    description = "will init an index based segment.\n"
    params      = ["path"]
    results     = ["seg"]
  }
  op "write_index_segment" {
    description = "will write a part into an index based segment.\n"
    params      = ["seg", "r", "index", "size"]
  }
}
interface "mover" {
  description = "is the interface for Move.\n"

  op "move" {
    description = "will move an object or multiple object in the service.\n"
    params      = ["src", "dst"]
  }
}
interface "prefix_lister" {
  description = "is used for prefix based storage service to list objects under a prefix.\n"

  op "list_prefix" {
    description = "will return list a specific dir.\n\nCaller:\n  - prefix SHOULD NOT start with /, and SHOULD relative to workdir.\n"
    params      = ["prefix"]
  }
}
interface "prefix_segments_lister" {
  description = "is used for prefix based storage service to list segments under a prefix.\n"
  embed       = ["segmenter"]

  op "list_prefix_segments" {
    description = "will list segments.\n\nImplementer:\n  - If prefix == \"\", services should return all segments.\n"
    params      = ["prefix"]
  }
}
interface "reacher" {
  description = "is the interface for Reach.\n"

  op "reach" {
    description = "will provide a way, which can reach the object.\n\nImplementer:\n  - SHOULD return a publicly reachable http url.\n"
    params      = ["path"]
    results     = ["url"]
  }
}
interface "segmenter" {
  internal = true

  op "abort_segment" {
    description = "will abort a segment.\n\nImplementer:\n  - SHOULD return error while caller call AbortSegment without init.\nCaller:\n  - SHOULD call InitIndexSegment before AbortSegment.\n"
    params      = ["seg"]
  }
  op "complete_segment" {
    description = "will complete a segment and merge them into a File.\n\nImplementer:\n  - SHOULD return error while caller call CompleteSegment without init.\nCaller:\n  - SHOULD call InitIndexSegment before CompleteSegment.\n"
    params      = ["seg"]
  }
}
interface "servicer" {
  description = "can maintain multipart storage services.\n\nImplementer can choose to implement this interface or not.\n"

  op "create" {
    description = "will create a new storager instance.\n"
    params      = ["name"]
    results     = ["store"]
  }
  op "delete" {
    description = "will delete a storager instance.\n"
    params      = ["name"]
  }
  op "get" {
    description = "will get a valid storager instance for service.\n"
    params      = ["name"]
    results     = ["store"]
  }
  op "list" {
    description = "will list all storager instances under this service.\n"
  }
}
interface "statistician" {
  description = "is the interface for Statistical.\n"

  op "statistical" {
    description = "will count service's statistics, such as Size, Count.\n\nImplementer:\n  - Statistical SHOULD only return dynamic data like Size, Count.\nCaller:\n  - Statistical call COULD be expensive.\n"
    results     = ["url"]
  }
}
interface "storager" {
  description = "is the interface for storage service.\n\nCurrently, we support two different types of storage services: prefix based and directory based. Prefix based storage\nservice is usually an object storage service, such as AWS; And directory based service is often a POSIX file system.\nWe used to treat them as different abstract level services, but in this project, we will unify both of them to make a\nproduction ready high performance vendor lock free storage layer.\n\nEvery storager will implement the same interface but with different capability and operation pairs set.\n\nEverything in a storager is an Object with two types: File, Dir.\nFile is the smallest unit in service, it will have content and metadata. Dir is a container for File and Dir.\nIn prefix-based storage service, Dir is usually an empty key end with \"/\" or with special content type.\nFor directory-based service, Dir will be corresponded to the real directory on file system.\n\nIn the comments of every method, we will use following rules to standardize the Storager's behavior:\n\n  - The keywords \"MUST\", \"MUST NOT\", \"REQUIRED\", \"SHALL\", \"SHALL NOT\", \"SHOULD\", \"SHOULD NOT\", \"RECOMMENDED\", \"MAY\",\n    and \"OPTIONAL\" in this document are to be interpreted as described in RFC 2119.\n  - Implementer is the provider of the service, while trying to implement Storager interface, you need to follow.\n  - Caller is the user of the service, while trying to use the Storager interface, you need to follow.\n"

  op "delete" {
    description = "will delete an Object from service.\n"
    params      = ["path"]
  }
  op "metadata" {
    description = "will return current storager's metadata.\n\nImplementer:\n  - Metadata SHOULD only return static data without API call or with a cache.\nCaller:\n  - Metadata SHOULD be cheap.\n"
    results     = ["meta"]
  }
  op "read" {
    description = "will read the file's data.\n\nCaller:\n  - MUST close reader while error happened or all data read.\n"
    params      = ["path"]
    results     = ["rc"]
  }
  op "stat" {
    description = "will stat a path to get info of an object.\n"
    params      = ["path"]
    results     = ["o"]
  }
  op "write" {
    description = "will write data into a file.\n\nCaller:\n  - MUST close reader while error happened or all data written.\n"
    params      = ["path", "r"]
  }
}

field "dir" {
  type = "string"
}
field "dst" {
  type = "string"
}
field "err" {
  type = "error"
}
field "index" {
  type = "int"
}
field "meta" {
  type = "info.StorageMeta"
}
field "name" {
  type = "string"
}
field "o" {
  type = "*types.Object"
}
field "pairs" {
  type = "...*types.Pair"
}
field "path" {
  type = "string"
}
field "prefix" {
  type = "string"
}
field "r" {
  type = "io.Reader"
}
field "rc" {
  type = "io.ReadCloser"
}
field "seg" {
  type = "segment.Segment"
}
field "size" {
  type = "int64"
}
field "src" {
  type = "string"
}
field "statistic" {
  type = "info.StorageStatistic"
}
field "store" {
  type = "storage.Storager"
}
field "url" {
  type = "string"
}
