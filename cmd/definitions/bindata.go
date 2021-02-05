// Code generated by go-bindata. DO NOT EDIT.
// sources:
// cmd/definitions/tmpl/function.tmpl (332B)
// cmd/definitions/tmpl/info.tmpl (1.699kB)
// cmd/definitions/tmpl/object.tmpl (1.908kB)
// cmd/definitions/tmpl/operation.tmpl (1.02kB)
// cmd/definitions/tmpl/pair.tmpl (483B)
// cmd/definitions/tmpl/service.tmpl (8.245kB)

// +build tools

package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _cmdDefinitionsTmplFunctionTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x90\xcf\x4a\xc4\x30\x10\xc6\xef\x7d\x8a\x8f\x25\x87\xee\xb2\xe6\x01\x04\x4f\x45\x41\x10\x59\xd4\xbb\x0c\xd9\x74\x0d\x36\x93\x92\x4c\xeb\x42\xcc\xbb\x4b\x5a\x75\xf1\xe0\x29\x21\xbf\xc9\xf7\x67\x72\xbe\x42\x24\x3e\x59\xa8\xd7\x3d\xd4\x8c\xeb\x1b\x28\x7d\x37\xb1\x49\x28\xa5\xa9\xd8\xf5\xe0\x20\x50\xb3\xbe\xf7\xe3\x60\xbd\x65\xb1\xc7\x1f\xa8\x7a\x7e\x5f\xfe\xcc\xfa\x91\xbc\xc5\x27\x24\x74\xe4\xed\x50\x07\xfa\x89\x0d\xda\x84\x5d\xce\x50\x17\x7e\xa0\x64\xa8\x0e\x6c\x91\x73\x15\x28\xa5\x35\x72\x86\x09\x2c\xf6\x2c\xba\x5b\xcf\x7d\xa5\xb3\x3e\x50\x24\x9f\xf4\x4b\x74\xfe\x81\x92\xe8\x67\x89\x8e\x4f\xb7\x7c\x4c\x1f\x4e\xde\xba\xe0\x3d\x95\x82\x30\x0a\x76\x23\xb9\xf8\x8f\x55\x7d\xae\x51\xff\xfa\xb7\x8b\xc3\x93\x4d\xd3\x20\xe9\x5b\x79\x89\xd5\x00\xc0\x48\xec\x4c\xbb\xa9\xe5\xdd\xa5\xf9\x66\xdb\xac\xd5\x2d\xff\x6e\x61\xbd\x7e\x05\x00\x00\xff\xff\x3b\x1f\xad\x5f\x4c\x01\x00\x00")

func cmdDefinitionsTmplFunctionTmplBytes() ([]byte, error) {
	return bindataRead(
		_cmdDefinitionsTmplFunctionTmpl,
		"cmd/definitions/tmpl/function.tmpl",
	)
}

func cmdDefinitionsTmplFunctionTmpl() (*asset, error) {
	bytes, err := cmdDefinitionsTmplFunctionTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cmd/definitions/tmpl/function.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xd1, 0x53, 0x97, 0x24, 0x27, 0xb1, 0x48, 0xc, 0x9a, 0xff, 0xbf, 0xcc, 0x50, 0x35, 0x3e, 0xf2, 0xfd, 0x26, 0xfd, 0x33, 0xa9, 0x5, 0x3f, 0xa2, 0x4c, 0x9b, 0x7, 0xfa, 0x1a, 0x11, 0xb7, 0xe0}}
	return a, nil
}

var _cmdDefinitionsTmplInfoTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x94\x5f\x6f\xd3\x30\x14\xc5\xdf\xfd\x29\x0e\x55\x85\x1a\xd4\x35\x9b\x84\x78\x80\xf5\x69\x1b\x68\x42\xdb\x90\x36\xf1\x82\x10\x72\x93\x9b\xca\x34\xb6\x23\xdb\x89\x56\x32\x7f\x77\x64\x27\x74\x4d\xf7\x07\x8a\x78\xe1\xcd\x37\xf6\xbd\xe7\x77\x8f\xaf\x93\xa6\x38\xd1\x39\x61\x49\x8a\x0c\x77\x94\x63\xb1\xc6\x52\x6f\x62\x08\xe5\xc8\x28\x5e\xa6\x99\xcc\xdf\xe1\xf4\x0a\x97\x57\x37\x38\x3b\x3d\xbf\x99\xb1\x8a\x67\x2b\xbe\x24\xb8\x75\x45\x96\x31\x21\x2b\x6d\x1c\x26\x0c\x00\x46\x85\x74\x23\x96\x30\xd6\xb6\x07\x30\x5c\x2d\x09\xe3\xd5\x14\x63\xa1\x0a\x6d\xf1\x76\x8e\xd9\x79\x58\x5d\xf0\x0a\xde\xb3\xb6\xc5\xd8\x92\x69\x44\x46\x97\x5c\x52\xd8\x1f\xaf\x70\x07\xa7\x4f\xb8\xa4\x32\x1c\x61\x69\x8a\xf7\x82\xca\x1c\x42\xe5\x74\x0b\xa1\xd0\xb6\xdb\x49\xde\x63\x21\x1c\xcb\xb4\xb2\x01\x62\x47\xb7\x89\x35\x3b\x75\xef\x23\xe2\x6e\xfa\x79\xa8\x1b\x48\x9a\x59\x84\x08\xf2\x9f\xb8\xcd\x78\xd0\xc7\x1c\x47\xc7\xc7\x61\x77\xd5\x01\x1f\x80\x54\x1e\x96\x09\x63\xc1\x00\xec\xf6\x30\x4c\xb7\xce\xd4\x99\x43\xdb\x2b\x6f\xd8\xbe\x3d\xc5\x16\x30\x6e\xd6\x55\x57\xcb\xfb\xad\x2f\xf7\x67\x36\x0c\x31\x4e\xd3\x60\x00\x6a\x4b\x39\xb8\x05\x0f\x91\xe4\x15\x0a\x6d\xa0\x17\xdf\x29\x73\x68\x78\x59\xd3\x14\x87\x90\xc4\x95\x85\xd2\x0e\x96\xdc\x14\x47\xfd\x07\x4b\x2e\x96\x8a\x75\x84\x72\x6f\x5e\xc7\x50\x42\xf2\xea\x8b\x75\x46\xa8\xe5\xd7\x38\x10\x05\xcf\xa8\xf5\xac\x57\x7e\xde\xeb\xb0\x2b\x8a\x40\x7f\x76\x1b\x27\xc4\x7b\x56\xd4\x2a\xc3\x44\xe2\xd5\xb3\xae\x25\xf8\x40\xae\x6b\xfc\x54\xd8\xaa\xe4\xeb\xde\x8d\x49\x32\xf4\xa3\xf7\xd5\x90\xab\x8d\x82\x9c\x3d\xb0\x2f\x90\xfe\xa9\xe6\xf5\x13\x9a\xcd\x50\x33\xf9\x4d\xa1\x9e\xe9\x11\x18\xcc\xd1\x0c\x78\x59\x3f\x50\xa5\x8d\xac\xff\xc0\x9d\xc9\x00\x75\x8a\x85\xd6\x65\xd2\x13\x89\x02\x72\x16\x6e\xf8\xe5\x9e\x4f\xe0\xc5\x1c\x87\x7d\x8d\xe7\xdd\x9e\xc2\x99\x9a\xe2\x41\xbf\xdd\xe8\x16\xd4\x1d\x7e\x90\xd1\x9f\xc3\x3c\xc6\x84\x82\x97\x96\xf6\xb9\xa5\x8b\xda\xba\xfd\xa6\xe3\xaf\xfb\x9e\x0f\xfb\xae\xb8\x12\xd9\xa4\x90\x6e\x76\x5d\x19\xa1\x5c\x31\x19\x3d\xc6\xfa\x91\x16\x7c\x71\xff\x72\x7f\xdd\xbd\xd8\x3c\xbb\x51\x92\x3c\xb4\xe8\xbf\x99\xdd\xce\xca\xbb\xf9\x7e\x5e\x3e\x3a\xf6\xdd\x3f\x6c\xe7\x97\x76\xbf\xfc\x19\x00\x00\xff\xff\x1a\xaf\xa4\x3e\xa3\x06\x00\x00")

func cmdDefinitionsTmplInfoTmplBytes() ([]byte, error) {
	return bindataRead(
		_cmdDefinitionsTmplInfoTmpl,
		"cmd/definitions/tmpl/info.tmpl",
	)
}

func cmdDefinitionsTmplInfoTmpl() (*asset, error) {
	bytes, err := cmdDefinitionsTmplInfoTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cmd/definitions/tmpl/info.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xbe, 0xfc, 0x4d, 0xd5, 0xaa, 0x43, 0xd, 0x38, 0x2a, 0xe9, 0x1b, 0x28, 0x33, 0xa1, 0xf, 0x25, 0xb3, 0x72, 0xf8, 0xb5, 0xf, 0xf, 0x28, 0x3d, 0xdd, 0x4d, 0x3d, 0x95, 0x9e, 0xd1, 0x56, 0x3e}}
	return a, nil
}

var _cmdDefinitionsTmplObjectTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x55\x41\x8f\xeb\x34\x10\x3e\xd7\xbf\xe2\xa3\x5a\x41\x82\xba\xc9\x7b\x80\x38\x3c\x5e\x0f\xe8\x75\x81\x95\xd8\x16\xa9\x85\xbb\xeb\x4c\x5a\xb3\x89\x1d\xd9\x93\xd2\xd2\xcd\x7f\x47\x4e\xd2\x25\xdd\x2d\xa8\x2b\x2e\xdc\x3c\x9e\x99\xcf\xf3\x7d\x33\xb6\xd3\x14\x9f\x6c\x46\xd8\x90\x21\x27\x99\x32\xac\x0f\xd8\xd8\x67\x1b\xda\x30\x39\x23\x8b\x54\x95\xd9\x77\x98\x2d\x30\x5f\xac\x70\x37\xbb\x5f\x25\xa2\x92\xea\x51\x6e\x08\x7c\xa8\xc8\x0b\xa1\xcb\xca\x3a\x46\x24\x00\x60\x9c\x97\x3c\xee\x56\xac\x4b\xea\x97\xfe\x60\xd4\x58\xc4\x42\xa4\x29\x7e\xd0\x54\x64\xd0\x26\xa3\x3d\xb4\x81\x5d\xff\x4e\x8a\xb1\xd6\x2c\x94\x35\x3e\xe0\x1c\x8f\xb7\x70\xd2\x6c\x08\x37\x8f\x13\xdc\xec\xf0\x61\x8a\x64\xd1\xc6\x3d\x10\x4b\x34\x4d\x8b\xda\x65\xde\x07\xa0\xe3\x11\x37\xbb\x64\x2e\x4b\xc2\x13\xd8\xfe\x22\xbd\x92\x05\x9a\x06\xb5\x36\xfc\xed\x37\x98\xe2\xfd\xc7\x8f\x21\xe8\x31\x24\x07\x7c\x32\x59\x58\x76\x25\x75\xd8\xd0\x1e\xbc\x25\xf8\x52\x16\x05\x79\x46\x6d\x34\x87\x12\x37\xf6\xd6\xb3\x75\x72\x43\x89\x48\xd3\x90\x30\x5f\xac\xee\x96\x1f\xc2\x0a\xb8\xed\xd3\xbf\xf0\xc8\x03\x35\x8f\xe5\x4f\x8b\x5f\x7f\x9e\xc1\x58\xc6\x9a\xa0\xb6\x81\x4a\x06\x5b\xb3\xd7\x19\xc1\x93\xdb\x69\x45\x3e\x39\x4f\xc7\xa7\xef\xe7\x41\xe2\x90\x61\x2b\x4d\xd9\x0b\xb7\xf6\x50\xd6\xa8\xda\x39\x32\x0c\x2f\x73\x4a\x44\x68\xc0\xc9\xef\xd9\xd5\x8a\x71\xbc\x56\xbd\x10\xa6\xf3\x20\xdb\x8c\xbc\x72\xba\x62\x6d\xcd\xc9\x89\x34\x45\xa7\xe9\x05\xe7\x40\xbe\xce\x0c\x71\xab\x43\x45\xad\xfe\x4d\x33\xd8\x79\xa1\xb6\x18\xa5\x29\x54\xa1\x03\x83\x5e\xec\x93\x65\xf0\xc7\x56\xab\xed\x80\xad\x2c\xf4\x8e\x12\x31\xea\x23\x96\x5d\x07\x9c\x10\x7d\x7d\x6b\xcd\xa8\x3d\x65\x90\x1e\x32\x58\xa5\xac\x90\x5b\x77\x9a\xa8\x9d\x2c\x6a\x9a\xe0\x1d\x4a\x92\xc6\xb7\xdd\xf0\xc4\x13\xbc\xef\x37\x3c\x71\x0b\xd5\xe2\xb4\x63\x22\x46\x99\x35\xd4\x1a\x5f\x7f\x25\x46\x65\xf0\x86\xc9\x4d\x1e\x6a\xa6\xbd\x68\x84\xb8\x46\xdc\xbf\x85\xbd\xdb\xb7\x17\xa3\x69\x44\x5e\x1b\x85\xc8\xe2\xcb\x2e\x34\xc6\x8f\xc4\xbd\xbe\xda\x57\x85\x3c\xf4\xd2\x45\xf1\xb9\x78\x38\xb6\x25\x3a\xe2\xda\x19\xd8\xe4\x95\xd6\xa1\xa8\x57\xe0\xcb\x7f\x00\xdf\x9d\x83\xc7\xa7\x8c\xfe\x94\x0b\xf0\x98\x62\x77\x56\x81\xe8\xfb\x59\xf8\xf6\xf4\xb7\x10\x8b\xce\x0e\x9f\x60\x6d\x6d\x11\x3f\x1f\xed\x59\x72\x14\x77\xcd\xd5\x39\x6c\x12\xda\xf2\xf9\x55\x77\xfc\xb3\x29\xde\xf5\x38\xff\x2e\xd6\x04\xec\x6a\x6a\x03\x1b\x31\xa4\x35\xa8\xec\x09\x7f\x92\xb3\xbf\x85\xd9\x69\x33\x72\x59\x78\xba\xa8\xf2\x43\xed\xf9\x6d\x6d\xfc\x4f\x2c\xa7\xe7\x2c\x2b\x69\xb4\x8a\xf2\x92\x93\x65\xe5\xb4\xe1\x3c\x1a\xf7\x73\x3f\x40\x68\x9a\x70\x91\xfa\xc9\x1f\xc7\x71\x4f\xfd\xff\x31\x52\x1d\xf9\xa7\xe9\x35\xec\x2f\xce\x60\xf7\xa6\x0c\x9f\x97\x57\x65\xab\xc2\x1a\x8a\xf6\x83\x9d\xe3\xf3\x13\x76\xcd\x0f\x73\xb1\xfc\xfd\x25\xd1\x5e\x3c\x8c\x03\x86\x6d\x42\xf8\xde\x9a\xbf\x02\x00\x00\xff\xff\xe0\xaf\xa2\xb9\x74\x07\x00\x00")

func cmdDefinitionsTmplObjectTmplBytes() ([]byte, error) {
	return bindataRead(
		_cmdDefinitionsTmplObjectTmpl,
		"cmd/definitions/tmpl/object.tmpl",
	)
}

func cmdDefinitionsTmplObjectTmpl() (*asset, error) {
	bytes, err := cmdDefinitionsTmplObjectTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cmd/definitions/tmpl/object.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xea, 0xa8, 0x6e, 0xe2, 0x32, 0xf5, 0x94, 0xf, 0x33, 0xa9, 0x21, 0x65, 0x6b, 0xd, 0xec, 0xdc, 0x2d, 0xd5, 0xda, 0x1d, 0xb8, 0xc7, 0xec, 0x5a, 0xa5, 0x69, 0xb8, 0x45, 0x20, 0x16, 0x5, 0x35}}
	return a, nil
}

var _cmdDefinitionsTmplOperationTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x52\xcd\x6a\xe3\x30\x10\xbe\xeb\x29\x06\xe3\x83\x0d\x59\xfb\xbe\xb0\x87\x65\x97\x42\x2f\xad\x69\x0f\x3d\x16\x55\x95\xdd\xa1\xb6\xa5\x4a\x4a\x48\x50\xf5\xee\x65\x64\x25\x38\xc6\x25\xa5\xf4\x66\xcf\xcf\xf7\xa3\xf9\x34\x17\xaf\xbc\x93\xe0\x0e\x5a\x5a\xc6\x70\xd0\xca\x38\x28\x18\x00\x40\x26\xd4\xe8\xe4\xde\x65\xd3\x1f\xaa\x8c\x95\x8c\x79\xff\x0b\x0c\x1f\x3b\x09\xf9\xe3\x06\x72\x84\xdf\x7f\xa0\xba\x1e\x9d\x34\x2d\x17\xd2\x42\x08\xcc\x7b\xc8\xb1\xfa\x2f\xad\x30\xa8\x1d\xaa\x91\x8a\xc4\x00\xa9\x83\x56\xf7\xfc\x70\xc3\x07\x09\x21\x00\x1e\x97\xc1\x47\x26\x62\xc0\x16\x94\x81\x42\xbe\xd1\x7c\x1c\xcc\xac\x34\x3b\x14\xd2\x64\xe5\xa2\xee\x94\xe1\x1d\xd5\x43\x88\xfb\xf7\xce\xe0\xd8\x15\x25\xd8\xf8\x71\xc2\x94\xe3\x33\x09\x49\xff\x73\x13\x4a\x93\x8b\x1c\xab\x5b\x1d\x0d\xd0\x44\x5d\x47\xb5\x4a\x4f\x34\xef\xe0\x54\xc3\xad\xe0\x3d\x49\x4e\x9d\x85\xc5\x04\xbc\xbe\x53\xa4\xce\x95\x32\x03\x77\x0d\x37\x7c\x20\xae\x12\xce\x1b\x77\xd2\x6e\x7b\x67\x1f\xd0\xbd\x34\xd3\x6d\xce\x5c\x66\xb4\x72\x51\x1f\x6d\xff\x9b\x8e\xf7\x1d\xa9\xb3\xf5\x42\xb8\x3d\xa4\x1c\x54\xa9\xb6\xf9\x61\x27\xde\x1f\x4f\x13\xd8\xfc\x4e\x31\x31\x0d\x47\xd3\xa8\x1e\xc5\x81\xce\xb9\x15\x2e\x85\xe4\x6f\xdf\xc3\x93\x52\xfd\x17\xf3\x88\x2d\x8c\x72\x25\x4b\xd4\xae\x6b\xd0\x1c\x8d\x85\x56\x99\x79\x18\xfd\x69\xfe\xec\x75\x96\x8c\x6b\xe1\xf9\x3c\x39\xab\xa2\x77\x11\x41\xe9\xaa\x89\x3a\x2e\x81\x50\x63\xb7\xc0\x3d\xc1\xa6\xd7\xbb\xf8\x19\xd8\x47\x00\x00\x00\xff\xff\xf0\x4f\x23\x75\xfc\x03\x00\x00")

func cmdDefinitionsTmplOperationTmplBytes() ([]byte, error) {
	return bindataRead(
		_cmdDefinitionsTmplOperationTmpl,
		"cmd/definitions/tmpl/operation.tmpl",
	)
}

func cmdDefinitionsTmplOperationTmpl() (*asset, error) {
	bytes, err := cmdDefinitionsTmplOperationTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cmd/definitions/tmpl/operation.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xac, 0x77, 0x38, 0x72, 0x67, 0x12, 0xac, 0x40, 0x5c, 0xb6, 0x92, 0x9, 0x6e, 0xd8, 0xf4, 0xde, 0x4d, 0x9e, 0x30, 0x85, 0x95, 0xa4, 0xa1, 0x26, 0x9b, 0x9c, 0xa9, 0xe4, 0xd9, 0x43, 0xc0, 0x37}}
	return a, nil
}

var _cmdDefinitionsTmplPairTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x8f\x31\x6b\xfb\x30\x10\xc5\x77\x7d\x8a\x87\xf1\x90\x40\x62\x0d\xff\x2d\x7f\x3a\x35\x1d\x4a\x21\xc9\x10\xda\xb1\x28\xf2\x55\x11\xb1\x25\x21\x9f\xdd\x1a\xd7\xdf\xbd\xc8\x4e\x03\xa5\x4b\x35\xe9\x7e\xef\xde\xe3\x9d\x94\xb8\xf7\x25\xc1\x90\xa3\xa8\x98\x4a\x9c\x7a\x18\x7f\x9b\x61\x1d\x53\x74\xaa\x92\xba\x2e\xff\x63\xbb\xc7\x6e\x7f\xc4\xc3\xf6\xf1\x58\x88\xa0\xf4\x45\x19\x42\x50\x36\x36\x42\xd8\x3a\xf8\xc8\x58\x08\x00\xc8\xb4\x77\x4c\x1f\x9c\x89\x79\x34\x96\xcf\xed\xa9\xd0\xbe\x96\xca\x37\xeb\x92\x3a\x69\xfc\xba\x61\x1f\x95\x21\xd9\xfd\x93\xe1\x62\xe4\x99\x39\xe8\xca\x92\xe3\x6c\x72\x15\x7f\xf0\x71\x1f\xa8\xc9\x84\x58\x0a\x31\x0c\x6b\x44\xe5\x0c\x21\x7f\x5d\x21\xef\xb0\xb9\x43\x71\x48\xe5\x30\x8e\x93\x9a\x07\xa7\x6a\x4a\x3c\xef\x8a\x5d\xfa\x7e\x82\xfd\x41\x35\x5a\x55\x69\x47\x4a\xbc\x58\x3e\x0f\xc3\xf7\xe6\x38\xe2\xdd\x56\x15\x54\x08\x55\x8f\xc4\xaf\xbe\x71\x44\xa7\xaa\x96\xc0\x1e\xfb\xc0\xd6\xbb\x46\xcc\xf2\x96\x1a\x1d\xed\x84\x52\xe4\x5b\xeb\xf4\xaf\xd0\x45\x77\xcd\x3a\xf6\x21\xcd\x4b\xa4\x9a\x18\xa6\xab\x23\x71\x1b\xdd\x44\x66\x90\xde\x13\xf5\x1b\x64\x3f\x0a\x64\xab\x9b\xfa\x9c\xaa\x6c\xd0\xcd\x64\x14\xf3\xb5\xe4\xca\xd4\xe0\x2b\x00\x00\xff\xff\x26\xdf\x18\x87\xe3\x01\x00\x00")

func cmdDefinitionsTmplPairTmplBytes() ([]byte, error) {
	return bindataRead(
		_cmdDefinitionsTmplPairTmpl,
		"cmd/definitions/tmpl/pair.tmpl",
	)
}

func cmdDefinitionsTmplPairTmpl() (*asset, error) {
	bytes, err := cmdDefinitionsTmplPairTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cmd/definitions/tmpl/pair.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x59, 0x2f, 0x7a, 0x10, 0x8b, 0x7, 0xa5, 0x4, 0xf7, 0x8a, 0x7d, 0x60, 0x9f, 0xd0, 0xf9, 0xb1, 0xc1, 0xc9, 0x4b, 0x3d, 0x79, 0xdb, 0xc, 0xf8, 0x51, 0xfd, 0xe7, 0xd, 0x79, 0x8d, 0x96, 0x48}}
	return a, nil
}

var _cmdDefinitionsTmplServiceTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x59\xdd\x6e\xdb\x36\x14\xbe\xd7\x53\x9c\x0a\xc6\x20\x15\x8e\xb4\x61\x77\x19\x7c\xb1\xa5\x5d\x57\x74\x6b\x8d\x26\x5b\x2f\xd6\x21\x60\x24\xca\x26\x42\x93\x2a\x49\x2b\x31\x5c\xbd\xfb\xc0\x1f\x59\x94\x2c\x39\x76\x57\xac\x1b\x1a\xdf\x58\x26\xcf\xef\xc7\xef\x1c\x1d\xc9\x69\x0a\x17\x3c\xc7\xb0\xc0\x0c\x0b\xa4\x70\x0e\x37\x1b\x58\xf0\xdd\x6f\xa8\x08\x02\xc2\x14\x16\x0c\xd1\x34\x5b\xe5\xa9\xc4\xa2\x22\x19\xfe\x01\x9e\xbd\x81\xd7\x6f\xae\xe0\xf9\xb3\x97\x57\x49\x50\xa2\xec\x16\x2d\x30\x6c\xb7\x90\xbc\x46\x2b\x0c\x75\x1d\x04\x64\x55\x72\xa1\x20\x0a\x00\x00\xc2\x8c\x33\x85\xef\x55\x68\x7f\x11\x1e\x06\xf6\x6a\x41\xd4\x72\x7d\x93\x64\x7c\x95\x22\x2e\xcf\x72\x5c\xa5\x0b\x7e\x26\x15\x17\x68\x81\xd3\xea\xfb\xb4\xbc\x5d\xa4\x98\xe5\x25\x27\xac\xd1\x3e\x4a\x27\x13\x38\xc7\x4c\x11\x44\x4f\xd1\x5a\x2a\x55\x66\x94\xe0\xe3\x7d\x39\x44\xa4\x95\x4f\x8e\xd0\x50\x9b\x52\x8b\xc7\x41\x50\x21\x01\xd7\xd0\x46\x9a\xcc\x05\xaf\x48\x8e\x85\xdb\x69\xf2\x4e\xfe\x40\x74\x8d\xdd\xe2\xa5\xb5\xd4\xc8\x34\xfe\x93\x4b\x7b\xf1\x5c\x08\xde\xec\xb5\xd9\x24\x6f\x4a\x45\x38\x93\x41\x90\xa6\x70\xb5\x29\x31\x10\x09\x6a\x89\x41\x07\x03\x05\x17\x9d\xb3\xcb\x38\x93\xca\x8a\xcd\x20\xf4\x76\x42\xa3\xef\x3c\x01\xaa\x10\xa1\xe8\x86\x62\x28\x11\x11\x32\x71\x7a\x51\xb0\xdd\x9e\x81\x40\x6c\x81\x61\x72\x3d\x85\x49\x05\xe7\x33\x48\xe6\x5a\x46\x5b\xd7\x38\x69\x09\x52\x00\xe3\x0a\x26\x55\xf2\x82\xf2\x1b\x44\xdb\x3d\xbd\xf6\x0c\xcb\x4c\x10\x13\x75\x77\xe3\xe7\x35\xa5\x2e\x1c\x17\xdd\xa4\x09\xef\xda\x4a\xec\x82\x6d\x3c\x61\x96\x6b\x1b\xde\x65\x3c\x92\xc8\x0a\x2b\x94\x23\x85\x1e\xc8\xe5\x25\x2b\xf8\xd1\xb9\x94\x4c\xc7\x73\x3e\xdb\x85\xf6\x11\x14\x9f\x23\x99\xf9\x62\xc6\x84\x4e\x9b\xc8\x92\xa2\x4d\x73\x12\xe0\x3e\xad\x9d\xd9\x88\x94\x97\xdc\x41\x10\x7f\x73\x19\xb6\x16\xf7\x71\x3c\x3b\x05\xc7\x07\x4e\x7b\x0c\x1d\xbd\xfe\x10\x34\x69\x0a\xef\x88\x5a\x76\x42\xbd\x23\x94\x02\x2a\x4b\xba\x81\x4e\x94\x50\xe9\x1a\x01\xc5\xa1\xa1\xfa\x20\x04\xc5\x9a\x65\x7b\x46\xa3\xca\xd9\x32\x94\xaf\xeb\x18\x74\xf8\xb0\xdd\xc1\x2f\xb0\x5a\x0b\x66\x56\xdb\xc5\x57\x78\x73\xbe\xcf\xc9\xe9\x6e\xdf\x14\xed\x39\x54\x76\xa5\x0e\x3a\xc0\x79\x97\xc3\x00\x6a\x73\xb2\x44\x19\xee\xf0\x4c\xe1\x55\x49\x75\x67\x0e\x75\xc9\x5d\x33\x7c\x17\xc2\x0a\xdd\xe2\x4b\xaa\x79\x1c\x0d\xe1\x18\x9b\x45\x7c\x37\x6e\x46\xfa\x36\x1a\x0b\x26\x2b\x96\x8d\x78\xd7\x30\x1a\x94\x1f\x52\xed\xe7\x99\xe3\x82\xb0\xd6\xaf\x67\x7c\x52\x32\x9d\x38\x61\x39\xbe\x87\x04\xbe\x1d\x29\x93\x89\x26\xaf\x2f\xf8\x9d\xb1\xdd\x6c\xf7\x71\xb4\xe2\x9d\x4a\xea\xa7\xdf\x49\xa1\x64\x30\xa9\xf6\x4b\xca\x5d\x0e\x65\x61\x0f\xa1\x65\x74\x37\x89\x66\xbd\x60\xfb\x31\xbb\x8d\x5b\x13\x67\xc1\x86\x2a\x40\x97\x80\x76\x62\xd9\x0a\x75\xad\x2f\xb4\x4a\x5d\x37\x1d\xbc\x44\x42\xe2\x1c\xa4\x12\xeb\x4c\x05\xa6\x9f\xf7\x34\xb4\x42\x5d\x3b\x09\xc7\x6a\x83\x3f\xfc\xf9\x97\xa6\xb4\x45\x2f\x4d\xe1\x2d\xfe\xb0\x26\x02\xe7\x76\x77\x08\x53\xbd\xd1\x84\xbb\x93\x76\x70\xfd\x82\xa4\x71\x8a\x88\x18\x48\x05\x00\x6e\x38\xa7\x6d\x57\x1c\x15\x6b\xb7\x5d\x39\x0e\xf5\xb7\x34\x75\x75\x8e\xe8\x71\xd1\xee\xa4\x3f\x77\xb4\x47\xc7\xfb\x62\x37\x65\x1d\x15\x70\x2b\xfe\x6f\xe3\xdb\xd0\x4e\x48\x3c\x1f\xe1\x9e\xe9\xc2\x46\xc2\xb6\x4a\x69\xaa\x87\x30\xc5\xe1\xe9\x08\x5f\x6d\xeb\x3d\x64\x35\xe2\xa5\x6a\x38\x19\x43\x34\x66\x68\x0a\x58\xcf\x38\xb1\x63\xb2\xc0\x72\x4d\x95\x46\xed\x9b\x11\x85\xb6\x65\x1b\xe0\xcf\x41\xfb\x69\xba\xb2\xf9\xd2\x03\xd0\xf5\x14\x4c\xc7\xb0\xc7\x61\x42\x69\x15\xe5\x1d\x51\xd9\x12\xaa\xe4\x15\xde\x78\xcb\xc3\x35\x73\x62\xdd\xe8\x4f\x86\x24\x6e\xcf\xc5\xbb\x9d\x9c\xef\x44\xda\x5c\x93\x07\x98\x30\x03\x25\xd6\x78\x48\xf1\xb0\x56\x65\xe7\xcc\x24\xea\x13\x24\xee\xa4\xe5\xb1\x7a\xbc\x12\x4f\xac\xc6\xff\x3f\x06\x43\xd5\x7d\x6a\x85\x9f\x80\x82\x99\x76\xa6\xc0\xcd\xcd\x23\xe3\xac\xc2\x42\x79\x8c\x1f\xcf\x31\x1a\xcf\x30\xee\x78\x20\x85\xb6\xbe\xed\xac\xfd\x23\xfc\x8f\x3f\x03\xf3\xb4\xe3\xab\xd5\x80\xa9\x86\x65\xcf\xde\xd1\xd1\x14\x88\xca\x9e\xcd\x43\x07\x5a\xfb\xdd\xe1\x84\x5a\x26\x05\x3c\x39\x2e\xa4\xbd\xe1\x92\x11\x3a\x6d\x9f\xe6\x5e\xe3\x3b\xdd\x06\x1b\xfb\xe6\xa9\x2e\xb2\x13\x46\x87\x17\x67\x0d\x33\xf7\xba\x78\xe0\xd9\xb6\x21\x4d\xb5\x8f\xe0\xf0\x4c\xf6\xdf\x9e\x64\x76\x1a\x8f\xb3\xcc\xe3\x2c\xf3\xc5\x67\x99\x48\xc2\xd3\xdd\x66\xfc\x38\xd9\x3c\x4e\x36\x8f\x93\xcd\x67\x9e\x6c\x64\xf2\x15\xcf\x36\xdb\x2d\x29\x80\xd9\x97\x12\xa1\x7b\x57\x19\xfa\x88\x7a\xe9\xc8\x44\xfb\x9c\x73\x4a\xb2\x4d\xf2\x23\xa5\xf0\xf1\x63\x77\xcd\x6b\x78\xe3\x5b\xc7\x0d\x2c\xdd\xa4\xc7\x87\x97\xdf\x99\x5c\x97\x25\x17\xaa\x99\x5f\xaa\x78\xcf\xc6\x7e\x32\xdb\xad\xc6\x63\x20\xcb\x4f\x1f\xf5\x9c\x59\xf3\xfe\x66\xd8\xf9\x00\xf3\x73\x5c\xa0\x35\x55\x5d\x72\x1e\x75\x24\x27\x1e\xc7\x10\x0d\x3f\x1d\xd4\x3a\x38\x0a\xcc\x8c\x33\x45\x58\x8f\x8a\x7d\x8c\xbe\xca\x29\xd8\x7b\xb5\xf9\xa5\xdf\x4e\x36\xa1\x9c\xf0\x86\x72\x34\x9d\xc3\xd9\xf8\x7b\x43\x03\x3e\xb8\x51\xeb\x81\x19\x1f\xec\x1d\xc3\x27\xb7\xb9\x1c\xfa\x23\x22\x4d\x1b\xf1\xab\x25\x91\xd0\xc4\x69\xc7\xb1\x4c\x60\x8d\x00\x02\xf7\xbf\x21\xdc\x6c\x9a\x82\x4c\xec\x5c\xd3\x0e\x60\x93\x92\xd5\x75\xec\x39\x8d\xcc\x2b\xcf\x64\x8e\x04\x5a\xc9\xe4\x52\x09\xc2\x16\x5a\xc2\xad\xbf\x35\x34\xf0\x37\x5a\x16\x66\xea\xde\x3d\x50\x6b\xaf\xc9\x4f\x28\xbb\x5d\x08\xbe\x66\x79\x14\xf7\x99\x2a\xbd\x1a\x7e\x47\xd4\xf2\xc2\xea\x44\x99\xba\x9f\x42\x27\x82\x0b\x44\x29\x16\x2d\x45\xf7\x61\xf2\xf4\x0f\x20\x76\x28\xe7\x5e\x04\xbb\x0c\xdc\x5a\x2f\xa2\x93\x30\xc9\x71\x81\x85\x71\x1e\xc5\xdb\x5e\xbf\x30\x95\xa8\x96\x0d\x27\x9c\xf9\x39\x52\x4b\x9b\x74\xbf\x35\xba\x7f\x83\x10\xcb\x21\xc2\x1f\x9c\x72\x18\xc6\xee\x97\xd7\x54\xe3\xa1\x1b\x9d\x6d\x25\x6a\x09\x33\x08\xa7\xef\xc3\xf7\xe1\x5e\xef\x1d\x68\xe4\xfa\x83\x85\x00\x3d\x4d\x14\x5c\xac\x90\xb2\x7d\x24\xb4\xa9\x6b\x1e\xd7\x75\x68\x66\xef\xd6\x41\x5d\x43\x7b\xe2\xb5\x77\xfa\x15\x12\x7a\x84\x1e\x7f\x34\x68\x04\x79\xa9\xa6\x3b\xbf\x07\x1f\x0e\xcc\x58\xd5\x7a\x20\x85\x51\x7b\x32\xd3\x8d\xaa\x77\x77\xb0\xdc\xf3\x1a\xf4\x20\x2b\xbb\xf5\x79\x81\x56\xd8\x4c\x47\xfb\xcc\xbc\x12\x64\xf5\x2b\x92\xca\x51\xf4\x39\xcb\xf5\xe3\xc0\xf2\x82\xaf\x56\xa8\xae\x75\x0a\x0d\x67\x5b\x5c\xff\x0e\x00\x00\xff\xff\x43\xd5\xc1\xd4\x35\x20\x00\x00")

func cmdDefinitionsTmplServiceTmplBytes() ([]byte, error) {
	return bindataRead(
		_cmdDefinitionsTmplServiceTmpl,
		"cmd/definitions/tmpl/service.tmpl",
	)
}

func cmdDefinitionsTmplServiceTmpl() (*asset, error) {
	bytes, err := cmdDefinitionsTmplServiceTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cmd/definitions/tmpl/service.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x8e, 0x74, 0xbe, 0x35, 0xb5, 0x74, 0xb6, 0xc2, 0x1c, 0xab, 0xfd, 0x93, 0xee, 0x67, 0xa5, 0x2f, 0x16, 0x7e, 0x56, 0xe6, 0x16, 0xa3, 0xd9, 0x76, 0x34, 0xa3, 0x15, 0x17, 0x5a, 0xff, 0x3, 0xf4}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"cmd/definitions/tmpl/function.tmpl":  cmdDefinitionsTmplFunctionTmpl,
	"cmd/definitions/tmpl/info.tmpl":      cmdDefinitionsTmplInfoTmpl,
	"cmd/definitions/tmpl/object.tmpl":    cmdDefinitionsTmplObjectTmpl,
	"cmd/definitions/tmpl/operation.tmpl": cmdDefinitionsTmplOperationTmpl,
	"cmd/definitions/tmpl/pair.tmpl":      cmdDefinitionsTmplPairTmpl,
	"cmd/definitions/tmpl/service.tmpl":   cmdDefinitionsTmplServiceTmpl,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"cmd": {nil, map[string]*bintree{
		"definitions": {nil, map[string]*bintree{
			"tmpl": {nil, map[string]*bintree{
				"function.tmpl": {cmdDefinitionsTmplFunctionTmpl, map[string]*bintree{}},
				"info.tmpl": {cmdDefinitionsTmplInfoTmpl, map[string]*bintree{}},
				"object.tmpl": {cmdDefinitionsTmplObjectTmpl, map[string]*bintree{}},
				"operation.tmpl": {cmdDefinitionsTmplOperationTmpl, map[string]*bintree{}},
				"pair.tmpl": {cmdDefinitionsTmplPairTmpl, map[string]*bintree{}},
				"service.tmpl": {cmdDefinitionsTmplServiceTmpl, map[string]*bintree{}},
			}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory.
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
