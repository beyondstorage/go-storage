
pair "project" {
  type        = "string"
  description = "specify project name/id for this service or storage"
  parser      = ""
}
pair "read_callback_func" {
  type        = "func([]byte)"
  description = "specify what todo every time we read data from source"
  parser      = ""
}
pair "storager_func" {
  type        = "storage.StoragerFunc"
  description = "specify what todo with a storager"
  parser      = ""
}
pair "checksum" {
  type        = "string"
  description = "specify checksum for this request, could be used as content md5 or etag"
  parser      = ""
}
pair "context" {
  type        = "context.Context"
  description = "context in all request"
  parser      = ""
}
pair "name" {
  type        = "string"
  description = "specify the storage name"
  parser      = ""
}
pair "offset" {
  type        = "int64"
  description = "specify offset for this request, storage will seek to this offset before read"
  parser      = "parseInt64"
}
pair "work_dir" {
  type        = "string"
  description = "specify the work dir for service or storage, every operation will be relative to this dir. work_dir MUST start with / for every storage services. work_dir will be default to / if not set. \n For fs storage service on windows platform, the behavior is undefined."
  parser      = ""
}
pair "credential" {
  type        = "*credential.Provider"
  description = "specify how to provide credential for service or storage"
  parser      = "credential.Parse"
}
pair "file_func" {
  type        = "types.ObjectFunc"
  description = "specify what todo with a file object"
  parser      = ""
}
pair "object_func" {
  type        = "types.ObjectFunc"
  description = "specify what todo with an object"
  parser      = ""
}
pair "size" {
  type        = "int64"
  description = "specify size for this request, storage will only read limited content data"
  parser      = "parseInt64"
}
pair "dir_func" {
  type        = "types.ObjectFunc"
  description = "specify what todo with a dir object"
  parser      = ""
}
pair "segment_func" {
  type        = "segment.Func"
  description = "specify what todo with a segment"
  parser      = ""
}
pair "location" {
  type        = "string"
  description = "specify the location for service or storage"
  parser      = ""
}
pair "storage_class" {
  type        = "storageclass.Type"
  description = "specify checksum for this request, could be used as storage class"
  parser      = ""
}
pair "endpoint" {
  type        = "endpoint.Provider"
  description = "specify how to provide endpoint for service or storage"
  parser      = "endpoint.Parse"
}
pair "expire" {
  type        = "int"
  description = "specify when the url returned by reach will expire"
  parser      = "parseInt"
}
pair "http_client_options" {
  type        = "*httpclient.Options"
  description = "sepcify the options for the http client"
  parser      = ""
}
pair "index" {
  type        = "int"
  description = "specify the index of this segment"
  parser      = "parseInt"
}
