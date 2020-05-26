name = "s3"

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
  }
}
namespace "storage" {
  implement = ["dir_lister", "index_segmenter", "prefix_lister", "prefix_segments_lister"]

  new {
    required = ["location", "name"]
    optional = ["work_dir"]
  }

  op "list_dir" {
    optional = ["dir_func", "file_func"]
  }
  op "list_prefix" {
    required = ["object_func"]
  }
  op "list_prefix_segments" {
    required = ["segment_func"]
  }
  op "write" {
    required = ["size"]
    optional = ["checksum", "storage_class"]
  }
}

pairs {

  pair "storage_class" {
    type = "string"
  }
}

infos {

  info "object" "meta" "storage-class" {
    type = "string"
  }
}
