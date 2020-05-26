
interface "copier" {
  description = "is the interface for Copy."

  op "copy" {
    description = "will copy an Object or multiple object in the service."
    params      = ["src", "dst"]
  }
}
interface "dir_lister" {
  description = "is used for directory based storage service to list objects under a dir."

  op "list_dir" {
    description = "will return list a specific dir."
    params      = ["dir"]
  }
}
interface "dir_segments_lister" {
  description = "is used for directory based storage service to list segments under a dir."
  embed       = ["segmenter"]

  op "list_dir_segments" {
    description = "will list segments via dir."
    params      = ["dir"]
  }
}
interface "index_segmenter" {
  description = "is the interface for index based segment."
  embed       = ["segmenter"]

  op "init_index_segment" {
    description = "will init an index based segment."
    params      = ["path"]
    results     = ["seg"]
  }
  op "write_index_segment" {
    description = "will write a part into an index based segment."
    params      = ["seg", "r", "index", "size"]
  }
}
interface "mover" {
  description = "is the interface for Move."

  op "move" {
    description = "will move an object in the service."
    params      = ["src", "dst"]
  }
}
interface "prefix_lister" {
  description = "is used for prefix based storage service to list objects under a prefix."

  op "list_prefix" {
    description = "will return list a specific dir."
    params      = ["prefix"]
  }
}
interface "prefix_segments_lister" {
  description = "is used for prefix based storage service to list segments under a prefix."
  embed       = ["segmenter"]

  op "list_prefix_segments" {
    description = "will list segments."
    params      = ["prefix"]
  }
}
interface "reacher" {
  description = "is the interface for Reach."

  op "reach" {
    description = "will provide a way, which can reach the object."
    params      = ["path"]
    results     = ["url"]
  }
}
interface "segmenter" {
  internal = true

  op "abort_segment" {
    description = "will abort a segment."
    params      = ["seg"]
  }
  op "complete_segment" {
    description = "will complete a segment and merge them into a File."
    params      = ["seg"]
  }
}
interface "servicer" {
  description = "can maintain multipart storage services."

  op "create" {
    description = "will create a new storager instance."
    params      = ["name"]
    results     = ["store"]
  }
  op "delete" {
    description = "will delete a storager instance."
    params      = ["name"]
  }
  op "get" {
    description = "will get a valid storager instance for service."
    params      = ["name"]
    results     = ["store"]
  }
  op "list" {
    description = "will list all storager instances under this service."
  }
}
interface "statistician" {
  description = "is the interface for Statistical."

  op "statistical" {
    description = "will count service's statistics, such as Size, Count."
    results     = ["url"]
  }
}
interface "storager" {
  description = "is the interface for storage service."

  op "delete" {
    description = "will delete an Object from service."
    params      = ["path"]
  }
  op "metadata" {
    description = "will return current storager's metadata."
    results     = ["meta"]
  }
  op "read" {
    description = "will read the file's data."
    params      = ["path"]
    results     = ["rc"]
  }
  op "stat" {
    description = "will stat a path to get info of an object."
    params      = ["path"]
    results     = ["o"]
  }
  op "write" {
    description = "will write data into a file."
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
