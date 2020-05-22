
interface "copier" {

  op "copy" {
    params = ["dst", "src"]
  }
}
interface "dir_lister" {

  op "list_dir" {
    params = ["dir"]
  }
}
interface "dir_segments_lister" {
  embed = ["segmenter"]

  op "list_dir_segments" {
    params = ["dir"]
  }
}
interface "index_segmenter" {
  embed = ["segmenter"]

  op "init_index_segment" {
    params  = ["path"]
    results = ["seg"]
  }
  op "write_index_segment" {
    params = ["index", "r", "seg", "size"]
  }
}
interface "mover" {

  op "move" {
    params = ["dst", "src"]
  }
}
interface "prefix_lister" {

  op "list_prefix" {
    params = ["prefix"]
  }
}
interface "prefix_segments_lister" {
  embed = ["segmenter"]

  op "list_prefix_segments" {
    params = ["prefix"]
  }
}
interface "reacher" {

  op "reach" {
    params  = ["path"]
    results = ["url"]
  }
}
interface "segmenter" {

  op "abort_segment" {
    params = ["seg"]
  }
  op "complete_segment" {
    params = ["seg"]
  }
}
interface "servicer" {

  op "create" {
    params  = ["name"]
    results = ["store"]
  }
  op "delete" {
    params = ["name"]
  }
  op "get" {
    params  = ["name"]
    results = ["store"]
  }
  op "list" {
    description = "x"
  }
}
interface "statistician" {

  op "statistical" {
    results = ["url"]
  }
}
interface "storager" {

  op "delete" {
    params = ["path"]
  }
  op "metadata" {
    results = ["meta"]
  }
  op "read" {
    results = ["path", "rc"]
  }
  op "stat" {
    params  = ["path"]
    results = ["o"]
  }
  op "write" {
    params = ["path", "r"]
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
  type = "Storager"
}
field "url" {
  type = "url"
}
