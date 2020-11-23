
pair "content_md5" {
  type = "string"
}
pair "content_type" {
  type = "string"
}
pair "context" {
  type        = "context.Context"
  description = "context in all request"
  default     = "context.Background()"
}
pair "continuation_token" {
  type        = "string"
  description = "specify the continuation token for list_dir or list_prefix."
}
pair "credential" {
  type        = "*credential.Provider"
  description = "specify how to provide credential for service or storage"
  parser      = "credential.Parse"
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
}
pair "interceptor" {
  type = "Interceptor"
}
pair "location" {
  type        = "string"
  description = "specify the location for service or storage"
}
pair "name" {
  type        = "string"
  description = "specify the storage name"
}
pair "offset" {
  type        = "int64"
  description = "specify offset for this request, storage will seek to this offset before read"
  parser      = "parseInt64"
}
pair "pair_policy" {
  type = "PairPolicy"
}
pair "read_callback_func" {
  type        = "func([]byte)"
  description = "specify what todo every time we read data from source"
}
pair "size" {
  type        = "int64"
  description = "specify size for this request, storage will only read limited content data"
  parser      = "parseInt64"
}
pair "storage_class" {
  type = "string"
}
pair "user_agent" {
  type        = "string"
  description = "specify the custom user-agent from client"
}
pair "work_dir" {
  type        = "string"
  description = "specify the work dir for service or storage, every operation will be relative to this dir. work_dir MUST start with / for every storage services. work_dir will be default to / if not set. \n For fs storage service on windows platform, the behavior is undefined."
}
