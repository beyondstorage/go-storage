name = "qingstor"

namespace "service" {

  new {
    required = ["credential"]
    optional = ["endpoint"]
  }

  op "create" {
    required = ["location"]
  }
  op "delete" {
    optional = ["location"]
  }
  op "get" {
    optional = ["location"]
  }
  op "list" {
    required = ["storager_func"]
    optional = ["location"]
  }
}
namespace "storage" {
  implement = ["copier", "dir_lister", "index_segmenter", "mover", "prefix_lister", "prefix_segments_lister", "reacher", "segmenter", "statistician"]

  new {
    required = ["name"]
    optional = ["disable_uri_cleaning", "location", "work_dir"]
  }

  op "list_dir" {
    optional = ["dir_func", "file_func"]
  }
  op "list_prefix" {
    required = ["object_func"]
  }
  op "list_prefix_segments" {
    optional = ["segment_func"]
  }
  op "reach" {
    required = ["expire"]
  }
  op "read" {
    optional = ["offset", "size"]
  }
  op "write" {
    required = ["size"]
    optional = ["checksum", "storage_class"]
  }
}

pairs {

  pair "disable_uri_cleaning" {
    type = "bool"
  }
  pair "storage_class" {
    type = "string"
  }
}

infos {

  info "object" "meta" "storage-class" {
    type = "string"
  }
}
