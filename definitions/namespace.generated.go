// Code generated by go generate cmd/definitions; DO NOT EDIT.
package definitions

import (
	"fmt"

	"go.beyondstorage.io/v5/types"
)

type Service struct {
	Features types.ServiceFeatures

	Create []Pair
	Delete []Pair
	Get    []Pair
	List   []Pair
}

func (s Service) ListPairs(name string) []Pair {
	switch name {
	case "create":
		return SortPairs(s.Create)
	case "delete":
		return SortPairs(s.Delete)
	case "get":
		return SortPairs(s.Get)
	case "list":
		return SortPairs(s.List)
	default:
		panic(fmt.Errorf("invalid op: %s", name))
	}
}

type Storage struct {
	Features types.StorageFeatures

	CombineBlock                   []Pair
	CommitAppend                   []Pair
	CompleteMultipart              []Pair
	Copy                           []Pair
	Create                         []Pair
	CreateAppend                   []Pair
	CreateBlock                    []Pair
	CreateDir                      []Pair
	CreateLink                     []Pair
	CreateMultipart                []Pair
	CreatePage                     []Pair
	Delete                         []Pair
	Fetch                          []Pair
	List                           []Pair
	ListBlock                      []Pair
	ListMultipart                  []Pair
	Metadata                       []Pair
	Move                           []Pair
	QuerySignHTTPCompleteMultipart []Pair
	QuerySignHTTPCreateMultipart   []Pair
	QuerySignHTTPDelete            []Pair
	QuerySignHTTPListMultipart     []Pair
	QuerySignHTTPRead              []Pair
	QuerySignHTTPWrite             []Pair
	QuerySignHTTPWriteMultipart    []Pair
	Read                           []Pair
	Stat                           []Pair
	Write                          []Pair
	WriteAppend                    []Pair
	WriteBlock                     []Pair
	WriteMultipart                 []Pair
	WritePage                      []Pair
}

func (s Storage) ListPairs(name string) []Pair {
	switch name {
	case "combine_block":
		return SortPairs(s.CombineBlock)
	case "commit_append":
		return SortPairs(s.CommitAppend)
	case "complete_multipart":
		return SortPairs(s.CompleteMultipart)
	case "copy":
		return SortPairs(s.Copy)
	case "create":
		return SortPairs(s.Create)
	case "create_append":
		return SortPairs(s.CreateAppend)
	case "create_block":
		return SortPairs(s.CreateBlock)
	case "create_dir":
		return SortPairs(s.CreateDir)
	case "create_link":
		return SortPairs(s.CreateLink)
	case "create_multipart":
		return SortPairs(s.CreateMultipart)
	case "create_page":
		return SortPairs(s.CreatePage)
	case "delete":
		return SortPairs(s.Delete)
	case "fetch":
		return SortPairs(s.Fetch)
	case "list":
		return SortPairs(s.List)
	case "list_block":
		return SortPairs(s.ListBlock)
	case "list_multipart":
		return SortPairs(s.ListMultipart)
	case "metadata":
		return SortPairs(s.Metadata)
	case "move":
		return SortPairs(s.Move)
	case "query_sign_http_complete_multipart":
		return SortPairs(s.QuerySignHTTPCompleteMultipart)
	case "query_sign_http_create_multipart":
		return SortPairs(s.QuerySignHTTPCreateMultipart)
	case "query_sign_http_delete":
		return SortPairs(s.QuerySignHTTPDelete)
	case "query_sign_http_list_multipart":
		return SortPairs(s.QuerySignHTTPListMultipart)
	case "query_sign_http_read":
		return SortPairs(s.QuerySignHTTPRead)
	case "query_sign_http_write":
		return SortPairs(s.QuerySignHTTPWrite)
	case "query_sign_http_write_multipart":
		return SortPairs(s.QuerySignHTTPWriteMultipart)
	case "read":
		return SortPairs(s.Read)
	case "stat":
		return SortPairs(s.Stat)
	case "write":
		return SortPairs(s.Write)
	case "write_append":
		return SortPairs(s.WriteAppend)
	case "write_block":
		return SortPairs(s.WriteBlock)
	case "write_multipart":
		return SortPairs(s.WriteMultipart)
	case "write_page":
		return SortPairs(s.WritePage)
	default:
		panic(fmt.Errorf("invalid op: %s", name))
	}
}
