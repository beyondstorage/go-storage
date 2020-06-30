// Code generated by go-bindata. DO NOT EDIT.
// sources:
// tmpl/function.tmpl (258B)
// tmpl/info.tmpl (1.442kB)
// tmpl/open.tmpl (1.146kB)
// tmpl/operation.tmpl (904B)
// tmpl/pair.tmpl (2.309kB)
// tmpl/service.tmpl (7.679kB)

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

var _functionTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x8e\x41\x4a\xc4\x40\x10\x45\xf7\x39\xc5\x67\xc8\x22\x33\x68\x1f\x40\x70\x15\x74\x25\x32\xa8\x17\x28\xda\x1e\x6d\x4c\x55\x9a\x54\x05\x07\xda\xba\xbb\xf4\x44\x04\x5d\x15\xd4\xff\xbc\xff\x6a\xbd\x46\x7f\x92\x0f\xdc\xdc\xa2\x0f\xf7\xab\xc4\xf0\x48\x9c\xf0\x05\x9b\x47\xe2\x34\xc1\xbd\x3b\xad\x12\x31\x28\x0e\xb5\xa2\xbf\xe4\x5a\x28\x6e\xa5\x23\x69\xa4\xd6\xda\xa3\xd6\x46\x72\x1f\xa2\x9d\x11\x67\xb1\x74\xb6\x30\x6e\xf7\xaa\xa5\x1b\xff\x48\x0b\xb1\x86\x97\x25\xf3\x03\xa9\x85\x67\x5b\xb2\xbc\xdd\xc9\xab\x7e\x66\x7b\x1f\x67\x66\x72\xc7\x5c\x0c\x87\x42\x79\xf9\xb7\xe9\xde\x1e\xcd\xf8\xef\xfa\xf0\xcb\x7f\x4a\xba\x4e\xa6\x3f\xdc\x8b\x58\x07\x00\x85\x24\xc7\x61\x97\xb9\x4c\x89\x93\x18\xb2\xed\xf6\x9d\x7f\x07\x00\x00\xff\xff\x5b\xb1\x71\xf7\x02\x01\x00\x00")

func functionTmplBytes() ([]byte, error) {
	return bindataRead(
		_functionTmpl,
		"function.tmpl",
	)
}

func functionTmpl() (*asset, error) {
	bytes, err := functionTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "function.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xf3, 0x17, 0x45, 0xc4, 0x8, 0xb1, 0x57, 0xf3, 0x99, 0x7b, 0x1f, 0xcd, 0xe4, 0x30, 0x7e, 0x3, 0x30, 0xba, 0x28, 0xed, 0x36, 0x1c, 0xd3, 0x4a, 0x24, 0xfc, 0x10, 0x73, 0xfe, 0xf7, 0xda, 0x75}}
	return a, nil
}

var _infoTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x92\xcf\xab\x9b\x4e\x14\xc5\xf7\xfe\x15\xe7\x1b\xf2\x40\xc1\xaf\xee\x5b\xb2\x28\x2f\xa5\xbc\x45\xf3\x0a\x09\x5d\xb4\x94\x32\xd1\xab\x48\xc6\x19\x99\x19\x2d\x69\x9e\xff\x7b\x19\x9d\xd4\x1f\x49\x4a\x5a\xde\x26\x44\xef\xb9\xe7\xde\xcf\xb9\xc6\x31\x1e\x65\x4a\xc8\x49\x90\x62\x86\x52\xec\x8f\xc8\xe5\xef\x67\x14\xc2\x90\x12\x8c\xc7\x49\x99\xbe\xc5\xfa\x19\x9b\xe7\x1d\xde\xaf\x9f\x76\x91\x57\xb1\xe4\xc0\x72\x2b\xc9\xa4\xe7\xc5\x31\xde\x71\x0e\xd6\xb0\x82\xb3\x3d\x27\x94\x64\x58\xca\x0c\x8b\xbc\x44\x0a\x6d\xe0\x7b\xa7\xd3\xff\x50\x4c\xe4\x84\xe5\xf7\x10\xcb\x06\x6f\x56\x88\x9e\x44\x26\x35\xda\xd6\x03\x00\xab\x58\x56\x82\x95\x64\x6b\xcb\x26\xda\xd8\xbf\x2f\x30\xf2\x13\xd3\x09\xe3\x63\x5d\x91\x59\xc1\xba\xd0\x15\x67\xc7\x4e\xe7\x8a\x33\xa3\xd5\x0d\x99\x95\x90\x48\x87\x47\xab\xdb\x26\xb2\x9a\x0d\xec\x0b\x8f\xcc\x50\x2e\xd5\xf1\xb2\xd6\x4f\x69\x5b\xac\xb0\xe8\xb5\x6e\xca\xc2\x1b\x8d\x08\xbc\xbb\xf0\x5f\x8b\xfe\x1f\xe0\x35\xa9\xa6\x48\x68\xe3\xc6\x57\xaa\x10\x26\xc3\xe2\x41\x3f\xe8\x05\xfc\x2b\xd9\x04\xdd\xdb\x2b\xc1\x04\x67\xdb\x38\xc6\x07\x32\x93\x90\x7e\x14\x9c\x23\x27\x83\x49\x54\x68\x18\xaf\x09\x99\x92\xe5\xe8\xbb\xb1\x16\x59\x2d\x12\xf8\xe5\xc5\x86\x6d\x1b\xcc\xbd\xfd\x00\x7e\xef\xba\x3b\x56\xf6\x45\x88\xbd\x94\x3c\x00\x4e\x9d\x55\x13\x42\x1e\x2c\x5a\x19\x95\x5f\x2f\xfd\xc6\x56\xdf\xba\x86\x22\xc3\x7f\xf2\xe0\xba\x87\xd0\xbf\x90\x92\x9f\xbb\x7d\x47\x91\x2b\x32\xb5\x12\x8e\x6a\xac\x08\x91\x31\xae\x69\x48\x9d\xeb\xdb\x8d\xdd\xe2\x2f\xf8\xf9\xa7\xfe\xe1\x6a\xfd\xaf\x33\x68\xa2\x29\x7c\x10\xc2\xa8\x9a\x9c\xf0\x7c\x8e\x8f\xb5\x36\x2e\xb6\x21\xfc\xd7\xbc\xc9\x68\xc0\xe8\x2e\x93\xcd\xce\xf7\x70\x8b\xdf\x71\x8d\x39\xda\x0c\x6a\x7b\x15\x48\xdf\x00\x2a\x84\x91\x7f\x01\xb4\x9d\xc1\x34\x53\x98\xe0\x4a\x93\xe3\xbb\x03\x0c\x2b\x34\x93\x2c\x1c\xd9\xe8\xd0\xbf\x02\x00\x00\xff\xff\x72\xdd\x2f\x33\xa2\x05\x00\x00")

func infoTmplBytes() ([]byte, error) {
	return bindataRead(
		_infoTmpl,
		"info.tmpl",
	)
}

func infoTmpl() (*asset, error) {
	bytes, err := infoTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "info.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x62, 0x3c, 0xa5, 0x63, 0x49, 0xa1, 0x4f, 0xb6, 0xc6, 0xa1, 0x7d, 0xfa, 0x3c, 0x28, 0xef, 0x86, 0xa3, 0xaa, 0x1b, 0xf2, 0x4e, 0xdd, 0x7a, 0xeb, 0x3a, 0xbf, 0x42, 0xf, 0xcd, 0x6a, 0x44, 0xe3}}
	return a, nil
}

var _openTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x92\x4f\x8f\xd3\x30\x10\xc5\xef\xfe\x14\x4f\xd5\x1e\x5a\xd4\xda\xd2\x1e\x17\xed\x89\x05\x89\x03\x5d\x24\x7a\x43\x08\x79\x9d\x59\x63\xd1\xd8\x66\xec\x06\x55\x51\xbe\x3b\x4a\xf3\x87\x00\x91\x48\x39\x71\x88\x92\x8c\x9f\xdf\xcc\xfb\x69\x94\xc2\xab\x50\x10\x2c\x79\x62\x9d\xa9\xc0\xd3\x19\x36\x8c\xff\xa8\x9c\x86\xf3\x99\xd8\xeb\xa3\x32\x65\xa1\x4c\x60\x3a\x65\x77\x4c\x2f\xf1\xf0\x88\xfd\xe3\x01\xaf\x1f\xde\x1e\xa4\x88\xda\x7c\xd5\x96\x30\x9e\x0b\xe1\xca\x18\x38\x63\x2d\x00\x60\x65\x5d\xfe\x72\x7a\x92\x26\x94\x4a\x87\xb4\x2b\xa8\x52\x36\xec\x52\x0e\xac\x2d\xa9\xea\x76\xb5\x4c\xa6\xf2\x39\x52\x5a\x89\xba\xde\x81\xb5\xb7\x84\x9b\xcf\x5b\xdc\x24\xae\x70\x77\x0f\xf9\x81\xb8\x72\x86\x12\x9a\x66\xa1\x5f\xea\x6f\xa8\xba\xbe\xd8\xc8\xbd\x2e\x09\x4d\xd3\xb5\x20\x5f\xb4\x56\x1b\x21\xda\xbe\x08\x91\xfc\x9b\x93\x37\x78\x3e\x79\xb3\x0e\x31\x43\x4a\xf9\xe2\x32\x92\x7c\xaf\x1d\x6f\xb0\x6e\x27\xe9\xed\x87\x69\x78\x7b\xa9\xd0\xcf\x7a\xf7\xe6\x2d\x88\xb9\x7d\x02\x6f\x84\xa8\x34\x8f\x0d\xde\xe9\x88\x7b\x94\x3a\x7e\x4c\x99\x9d\xb7\x9f\x86\x83\xfa\xf7\xe4\x25\x65\x3d\x1b\xbd\xd5\xb9\x67\xd0\x37\xac\x8f\xe4\x3b\xe1\x25\x5c\x8a\xda\x50\xda\xe0\x76\x50\x76\xea\x89\x02\x4d\x23\x0f\xe7\x48\x77\x7f\x96\xf7\xf4\x7d\x3b\xda\xf7\x74\x26\x9f\xcd\x04\xd4\x90\xfe\x9f\x80\xcd\x81\x99\x1a\xce\x02\x9a\x0a\xfe\x03\x50\x63\x9a\xa5\xc0\xfa\xb5\xf8\x1b\xb0\x2b\x76\x69\x6a\x39\x8f\x6c\x22\xb8\x06\xd9\x72\x08\xc3\x7c\xbf\xa6\xfe\x11\x00\x00\xff\xff\xc4\x27\xa0\x37\x7a\x04\x00\x00")

func openTmplBytes() ([]byte, error) {
	return bindataRead(
		_openTmpl,
		"open.tmpl",
	)
}

func openTmpl() (*asset, error) {
	bytes, err := openTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "open.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x1f, 0x23, 0x46, 0xd2, 0x77, 0xf5, 0xa4, 0xbe, 0x4f, 0xbb, 0x20, 0x3e, 0xa9, 0x65, 0x60, 0x44, 0x31, 0xd0, 0x32, 0x68, 0xae, 0x41, 0xb8, 0x3e, 0x8a, 0x48, 0xcf, 0x4b, 0x4a, 0x90, 0xa9, 0x58}}
	return a, nil
}

var _operationTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x52\x4d\x6b\xe3\x30\x10\xbd\xeb\x57\x0c\xc2\x07\x1b\x12\x0b\xf6\xb8\xb0\xa7\xdd\x2d\xf4\xd2\x86\xf6\xd0\x63\x51\x94\x89\x32\x24\x96\x54\x49\x09\x09\xae\xff\x7b\x91\x2d\xa7\x49\x68\x09\x2d\xbd\x8d\x66\xde\x7b\x7a\xf3\xe1\xa4\x5a\x4b\x8d\x10\xa2\xf5\x52\x23\x63\xd4\x38\xeb\x23\x94\x0c\x00\x80\x2b\x6b\x22\xee\x23\x1f\x5e\x64\x39\x1b\x22\x4d\x71\xb5\x9d\xd7\xca\x36\x42\xda\x30\x5d\xe0\x4e\x68\x3b\xcd\x22\x62\xf7\x4b\xb8\xb5\x16\x01\x75\x83\x66\x24\x5f\xa5\xc4\x83\xc3\xf0\x25\xb0\x20\xb3\xb4\x9c\x55\x8c\xb5\xed\x14\xbc\x34\x1a\xa1\x78\x9e\x40\x41\xf0\xfb\x0f\xd4\xb7\x26\xa2\x5f\x4a\x85\x01\xba\x8e\xb5\x2d\x14\x54\xff\xc3\xa0\x3c\xb9\x48\xd6\xa4\x64\x92\x81\x5c\xa1\xe0\x36\xf2\x70\x27\x1b\x84\xae\x03\x1a\xc9\xd0\xf6\x96\xd2\x0f\xb4\x04\xeb\xa1\xc4\x97\x84\xef\x81\x3c\xa0\xdf\x91\x42\xcf\xab\x8b\xfc\x60\xd5\xf3\xaa\xeb\x7a\xfe\x63\xf4\x64\x74\x59\x41\xe8\x83\xa3\x26\x9a\x45\x32\x72\x7c\x9f\x74\x81\xcd\x1c\x17\xa9\x93\x82\xea\xff\x7d\x9c\xb5\x06\x6c\x06\x5c\x18\xff\x44\xf8\x54\xd7\xba\x2c\x7a\xef\xc2\xc8\x10\xa2\x57\xb4\x6e\xf0\xff\x0a\xd1\xce\x64\x50\x72\x93\x66\x91\x2b\x17\xb3\x1b\x5d\x7c\xc8\x29\x73\xe5\xc6\xfa\x46\xc6\x99\xf4\xb2\x49\x7f\x55\x70\x5e\x78\xc0\xb0\xdd\xc4\xf0\x44\x71\x35\x1b\x2e\xf1\x6c\x7c\x3c\x51\xae\xfa\x4b\xec\xbf\xc3\xa1\x7e\xc7\xea\x09\xbd\x54\x71\x0f\xf9\xe6\xeb\x9c\x9b\xfc\x70\x27\x6d\x3b\xae\xa6\x63\xef\x7b\x7a\x0b\x00\x00\xff\xff\xff\xa8\x61\x8f\x88\x03\x00\x00")

func operationTmplBytes() ([]byte, error) {
	return bindataRead(
		_operationTmpl,
		"operation.tmpl",
	)
}

func operationTmpl() (*asset, error) {
	bytes, err := operationTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "operation.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xf4, 0x1e, 0x4e, 0xc4, 0x20, 0x1a, 0xa3, 0x76, 0x72, 0x5b, 0x5e, 0x3, 0xb8, 0x63, 0xb8, 0x23, 0x8a, 0x12, 0xe0, 0x6, 0x29, 0x40, 0x0, 0xcc, 0xe8, 0x6a, 0x3f, 0x28, 0x3b, 0x35, 0x70, 0xc8}}
	return a, nil
}

var _pairTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x56\xc1\x6e\xe3\x36\x10\xbd\xeb\x2b\xa6\x82\x11\x48\x85\x2c\x15\x3d\xaa\xf0\xa1\x68\x7c\x28\x8a\xc6\x39\x18\xed\x21\x08\x16\x0c\x35\x96\x09\x53\x24\x41\xd2\xcc\x1a\x5a\xfd\xfb\x82\x94\xed\xc8\x96\x9d\xc4\x87\x05\x96\x27\x91\xf3\x66\xe6\xcd\x1b\x0e\xed\xa2\x80\xbf\x64\x85\x50\xa3\x40\x4d\x2c\x56\xf0\xb2\x83\x5a\x1e\xf7\xc0\x84\x45\x2d\x08\x2f\x68\x53\xfd\x01\xf7\x0b\x78\x58\x2c\x61\x7e\xff\xf7\x32\x8f\x14\xa1\x1b\x52\x23\x28\xc2\xb4\x89\x22\xd6\x28\xa9\x2d\x24\x11\x00\x40\x4c\xa5\xb0\xf8\xd5\xc6\x51\xbf\xad\x99\x5d\x6f\x5f\x72\x2a\x9b\x82\x48\x33\xad\xd0\x15\xb5\x9c\x1a\x2b\x35\xa9\xb1\x70\xbf\xc7\x9f\x83\x15\x6a\x53\x17\x06\xeb\x06\x85\xbd\xc5\x05\x45\xa5\x24\xbb\xcd\x87\x6a\xac\x50\x58\x46\xf8\x2d\x5e\x6b\x6b\x15\xe5\xec\x06\x7e\x76\xa7\xd0\xc4\x51\x1a\x45\x45\x01\x7f\x72\x0e\xc4\x11\xc6\xc9\x0b\xdf\x2b\x9b\x47\x54\x0a\xe3\x85\x6d\xdb\x29\x68\x22\x6a\x84\xc9\x97\x0c\x26\x0e\xca\x19\xe4\x8f\x1e\x03\x5d\x17\xb2\x79\xc4\x44\x09\xd2\xa0\xb7\x4d\x5c\xfe\xe0\x3f\xbf\x81\x95\x8f\xc4\x50\xc2\x0f\xb8\xa2\x80\xb6\x3d\x20\xbb\x0e\x5e\x19\xe7\xe1\xc4\xe5\xf7\x68\xa8\x66\xca\x32\x29\xde\xa2\x0e\xa0\x33\x88\x7b\xe0\x43\x7f\x10\x07\x5a\x28\x2a\x8f\x4e\xa3\x0f\x48\x7e\x86\x60\x51\xc0\xff\xcc\xae\xc7\x04\x89\x52\x7c\x07\x27\xd9\xc1\x11\xbe\x45\xb0\x12\x16\x81\xb2\xf1\xde\xcb\x35\x33\x41\x3b\x60\x06\xb6\x06\x2b\x6f\xbf\x58\xdd\x6a\x2b\xe8\x28\x59\xe2\xf6\xe0\xe5\x4e\xf9\x7d\x0a\xbf\x86\x1e\x85\x2a\xa0\x0d\x92\x68\xb4\x5b\x2d\xe0\xee\xcd\xd0\x9f\xfb\xf5\x0f\xee\xca\x13\xcd\xb2\xa3\xe9\x3f\xcf\xb6\x04\xd7\x9f\x74\x51\x37\x14\xcf\x53\x7f\x24\xda\x60\x5f\xad\x0a\x9f\x04\x36\x53\x07\x0d\x51\xbe\x86\x70\x1f\xc0\x70\x46\x31\xef\xb9\x07\x7c\xd2\x78\xc0\x93\xb1\x9a\x89\xfa\x39\x4c\xec\x8a\x50\x6c\xbb\x14\x92\xa7\xe7\x01\xf9\x0c\x50\x6b\xa9\xd3\x7d\x11\x7d\xb8\x72\x06\x0d\xd9\xe0\x39\xf2\xb7\x0c\x38\x8a\xa4\x49\xd3\x7e\x84\x1d\xd1\xde\xbb\x8f\xd0\x1f\xad\xa4\x86\x4d\x06\xa1\xc5\x7d\xcb\x1b\x78\x93\xc1\x3b\x28\x07\x03\x3a\x47\x93\x79\x65\x96\xae\x61\x33\x40\x7f\x7c\xb7\x87\xc8\xcf\xdc\xf1\xc3\xa2\xc4\xe0\x49\x3b\xca\x13\xf3\x80\x8f\x0e\x59\x5d\x9e\x78\x1d\xd2\x01\xb9\x51\xa8\xe3\xdd\x18\xc7\x0a\xc2\x3a\x98\x81\x76\x23\x9b\xa7\xce\x56\xde\x3d\xf4\x4d\x9f\x53\x3d\xe6\xe8\x5b\x79\x2d\x78\x68\x23\xcc\xa0\x6d\x8f\x91\xba\x2e\xd1\x2e\xbd\x88\x67\xab\x00\xff\x65\x06\x82\xf1\x0b\x45\x1d\xd6\xfe\x4a\x0b\xc6\x33\xb8\x9b\xfb\x2e\x5f\xc7\x2e\x54\x09\x71\xb8\xa0\x71\x76\x15\x34\xd7\xba\xf4\xa9\xaf\x23\xde\x99\x94\xf3\xe5\xf5\x2e\x0f\x4f\xcf\x5e\xfc\x77\x72\x9f\x0e\xda\xf9\x1a\xab\x3e\x3e\x19\x0c\xe6\xb9\xa9\xc2\x15\xd9\x72\x7b\xb9\x3d\x3f\x46\xc6\xb9\xd6\x7e\x18\x7c\xe1\xff\x32\xd3\x10\x4b\xd7\x3f\xbb\xac\xdd\xc9\x6c\x5f\x90\xf2\xaa\x8c\xfe\xbf\x03\x13\x5b\x8c\xc6\xa1\xfa\x27\x6b\xe6\x7f\x0a\x50\x54\x49\xd8\x66\x27\xcf\x70\x28\x7e\x93\x1d\xa8\x2a\xd7\xa5\xfb\xc7\x76\xf8\x70\xef\x1d\x05\xe3\x51\x17\x7d\x0f\x00\x00\xff\xff\x99\xad\x0d\x69\x05\x09\x00\x00")

func pairTmplBytes() ([]byte, error) {
	return bindataRead(
		_pairTmpl,
		"pair.tmpl",
	)
}

func pairTmpl() (*asset, error) {
	bytes, err := pairTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "pair.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x83, 0xe8, 0xeb, 0x7e, 0x52, 0x12, 0xbb, 0x36, 0x7f, 0x63, 0x14, 0x35, 0xe2, 0xd, 0xc3, 0xe5, 0x6a, 0x72, 0xae, 0x3e, 0x30, 0xef, 0x94, 0xbe, 0x81, 0x28, 0xc, 0x95, 0x55, 0xe9, 0xf9, 0x47}}
	return a, nil
}

var _serviceTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x59\x5b\x6f\xdb\xbe\x15\x7f\xf7\xa7\x38\x7f\xc1\x2d\xa4\xc2\x91\xb6\x3d\xba\xf0\xcb\x92\x36\x2d\xba\x25\xc1\x92\xad\xc0\xda\x20\x60\x24\xca\x26\x2c\x91\x2a\x49\x2b\xf1\x5c\x7d\xf7\x81\x17\x49\x94\x2c\xb9\x76\x97\x76\x7e\x68\x6d\xf1\x5c\x7f\xe7\x4a\x25\x8a\xe0\x9c\x25\x18\x96\x98\x62\x8e\x24\x4e\xe0\x71\x0b\x4b\xd6\xfc\x86\x92\x20\x20\x54\x62\x4e\x51\x16\xc5\x79\x12\x09\xcc\x4b\x12\xe3\xb7\x70\x71\x0d\x57\xd7\x77\xf0\xee\xe2\xe3\x5d\x38\x29\x50\xbc\x46\x4b\x0c\xbb\x1d\x84\x57\x28\xc7\x50\x55\x93\x09\xc9\x0b\xc6\x25\xf8\x13\x00\x00\x2f\x66\x54\xe2\x67\xe9\x99\x5f\x84\x79\x13\xf3\x6d\x49\xe4\x6a\xf3\x18\xc6\x2c\x8f\x10\x13\x67\x09\x2e\xa3\x25\x3b\x13\x92\x71\xb4\xc4\x51\xf9\x17\xef\x38\xb2\xa8\x58\x2f\x23\x81\x97\x39\xa6\xf2\x14\x16\x4c\x93\x82\x91\xd3\x78\x62\x8e\x13\x4c\x25\x41\xd9\x29\x5c\x2b\x29\x8b\x38\x23\x27\xd8\x67\xb1\x16\xc7\xd2\xcb\x6d\x51\x13\x17\xe2\x58\xfa\xa8\x40\x84\x9f\xa6\x22\x22\x34\x65\xde\x24\x98\x4c\x4a\xc4\xe1\x01\x5a\x3c\xc2\x1b\xce\x4a\x92\x60\x6e\x4f\x6a\x74\xfb\xcf\x6d\xa0\xc2\x5b\xf3\x7f\xfd\xd4\xe8\x09\x6f\xcd\xff\x2d\xb1\x81\x21\xbc\x35\x5f\xde\x71\xce\xea\xb3\x16\xd4\xf0\xba\x90\x84\x51\x31\x99\x44\x11\xdc\x6d\x0b\x0c\x44\x80\x5c\x61\x50\x26\x43\xca\x78\x27\x39\x63\x46\x85\x34\x64\x0b\xf0\x9c\x13\x4f\xf3\x5b\x4d\x80\x4a\x44\x32\xf4\x98\x61\xd0\x28\x85\x96\xcf\x9f\xec\x76\x67\xc0\x11\x5d\x62\x98\x3e\xcc\x60\x5a\xc2\x7c\x01\xe1\x8d\xa2\x51\xd2\x15\x96\x8a\x82\xa4\x40\x99\x84\x69\x19\x5e\x66\xec\x11\x65\xee\xd9\xb4\xa0\x4a\xe3\x7c\xa1\x8e\xb5\xf2\xef\x20\xd9\x0d\x12\x71\x4b\x17\x45\xca\x68\x4b\x59\x55\xf0\x44\xb2\x4c\x3f\x29\xc3\x0b\x2c\x62\x4e\xb4\xcb\x35\xb5\x52\xdf\x21\x37\x9e\x4d\x6b\xd7\x1e\x0c\x67\xe3\x68\x6d\x09\xa6\x89\x12\xe1\x7c\x0d\x46\x40\x50\x71\xff\x11\x08\x1f\x15\xcd\x4b\x82\x60\x65\x28\x9f\x89\x28\x32\xb4\xad\x63\x08\xf6\xe3\x08\x5a\x8c\x90\x39\xbe\xa9\x9f\xca\x46\x03\xc6\x6d\xcc\x8a\x9e\x52\x73\x70\x8e\x24\x5e\x32\xbe\xdd\x3f\x1b\x83\xf7\xec\x14\x78\x0f\xe7\xcf\x18\x6c\xc7\x40\x16\x45\xf0\x99\xc8\xd5\x7e\xde\xa0\xa2\xc8\xb6\xd0\x31\x12\x4a\x94\x6d\x30\x48\x06\x75\xf1\xa8\xda\x59\x11\xa1\xd3\x5d\x15\xd0\x46\xe0\x44\x9d\x0f\x26\x5d\xba\xa1\xf1\x9e\x32\xbf\xb4\xc4\xba\xb8\xaa\x2a\x80\x37\xba\x69\x68\xef\x60\xd7\x44\x8d\x63\xb9\xe1\x14\x5e\xb7\x87\xed\xd9\x27\xbc\x9d\xef\xa5\xf3\xac\x39\xfe\x97\xb2\x7a\x0e\xa5\x79\x52\x4d\x3a\xe0\x3a\x5f\x47\x60\x3e\x21\x43\x5f\x2a\x41\x4f\xcc\x4f\xc5\x60\xbb\xde\x95\x55\x5f\x70\x42\x65\x0a\xde\x2b\xf1\x4a\x78\xe0\x0f\xa4\x6e\xa0\x9f\x0e\xe4\x6d\xe0\xb4\x93\x4b\x2c\xf7\x33\x63\x89\xe5\x60\x5e\xa4\x9c\xe5\x90\x63\x89\x12\x24\x51\xa8\x45\xe8\x90\xf7\x84\xf8\xb9\x6e\x0c\x61\xdf\x6a\x15\x7b\xbf\x93\x0b\x33\x78\x64\x2c\x0b\xc0\xa6\x41\x39\x03\xb6\x56\xce\xe5\xe1\x25\x96\x7e\x5d\x95\x5d\x19\xae\xa6\x40\xb3\x91\x14\xfe\x60\x6b\x2b\xa3\x05\xff\xdf\x98\x33\x9d\x19\x2e\xf4\x36\xcb\x8c\x15\x2e\xc5\x0c\x52\x94\x09\xdc\xa2\x9f\x89\x71\x46\x6d\xfe\x77\xf8\xcf\x21\xfe\x36\x7a\xe6\x5f\x2b\xa0\x0c\xbb\x10\x04\x33\x90\x7c\x83\x2d\x61\x1d\x16\x61\x10\x6d\x03\xa0\xe3\x22\x46\xe2\x42\xa8\x64\x43\x71\x11\xc7\xc6\x65\x06\x7b\x45\x3a\x42\x69\x51\xb6\xce\xe4\xe1\xed\x51\x71\x9a\x41\x19\x38\x50\x8c\x54\xe7\x60\x71\x2a\x69\xa2\x40\x31\xee\x54\xa8\xc4\x79\x91\xa9\x7d\xd4\x53\x8d\xc9\x83\x1c\xad\xf1\x6d\xa6\x86\x93\x3f\x54\x98\x81\x7e\x88\x9f\xc6\x45\x08\x57\x46\x2d\x61\x5a\x86\xef\x37\x34\x1e\xd1\xac\x20\xd6\x4d\xf2\x47\xac\x7d\x1f\x13\x9c\x12\xda\xea\xed\x8e\x3f\xe5\x34\xa1\x09\x7e\x86\x10\xfe\x34\xd2\x5c\xa6\x2a\xce\x2e\xe1\x9f\xb5\xec\xfa\xb8\x8f\xa1\x21\xef\x0d\xc8\x71\x04\x95\x11\xd3\x72\xbf\x11\xd9\xaf\x43\x5e\x1c\x76\xc2\x3d\x4b\xe9\xb8\xdd\xd3\x94\xea\xea\x9f\xa6\x74\xa8\xb5\xd6\xb5\x51\x34\x73\xc0\x66\x99\x62\xab\xaa\xbf\xa3\x02\x56\x2c\x4b\x04\x20\x35\xdc\xba\x8b\x9a\x69\x2f\x88\x1f\x62\x5e\x40\x8e\x8a\x2f\x42\x72\x42\x97\xf7\x42\xf2\x4d\x2c\x77\x55\x3b\x83\xa2\x08\xfe\x81\xbf\x6d\x08\xc7\x89\x23\x73\x00\x73\x3d\x2a\xad\x1b\x0d\x47\xbf\xfb\x23\xc2\xc3\xf7\x9b\x2c\xb3\xc5\x32\x87\x46\xa1\x33\xd7\x7a\x7d\xc4\x5a\x61\x46\x33\xca\x8e\xb7\xa2\xe1\x78\x41\x2b\x2e\x9b\x1b\xe2\xd1\x66\xb4\x2c\x2f\x61\xc7\x0f\x13\xa2\xde\xf7\x0b\xc4\xd5\xd2\x62\x44\x6a\x1e\x7d\x03\x18\xe3\x32\x74\xce\x62\xa2\xfd\x83\x2f\xf7\xce\xda\x32\xf9\x85\x59\xb1\x9f\xf9\x9d\x63\xdb\xa0\x7f\x4b\x8e\x7c\x40\xe2\xb0\x59\x6a\x7a\x1f\xeb\xc1\xc9\x3e\xbc\x5c\x86\xb5\x1b\x9d\x56\x7f\x81\x53\xb4\xc9\xe4\xff\xe2\xe9\x9e\xbd\x2f\x14\x3e\x37\xa9\xb9\xc0\x37\x23\x39\xaa\xb7\x01\x4d\xd1\x59\xa6\x85\xee\xdf\x7a\x1b\x78\x33\x92\xde\xed\x72\x70\x48\x81\xcf\x0a\xd9\x4b\xf9\x00\xfc\x31\x99\x33\xc0\xea\x06\x1e\x34\xbb\x81\x50\xe8\xce\x17\xf0\x7a\x84\xa1\x57\x5b\x73\x50\xea\x66\x2e\x00\x7a\xaf\x11\x7a\x15\x44\x6b\xec\x3b\x8d\x59\xbf\x7c\x4a\x51\x8c\x77\x76\xfb\x53\x37\xf9\x07\xb5\xc0\xcc\x17\x36\x2d\xb4\xf1\xbb\xa1\xf8\xfb\x2b\x24\x6e\x38\x4e\xc9\xb3\x31\xc5\xbb\xc2\x4f\x5e\xe0\xc6\x83\xa4\x4a\x98\xd9\x42\xc7\x67\xc5\x97\x32\xfc\x84\xb7\xf7\x6f\x9d\xbd\xb3\xfe\xd8\xcd\x88\x92\x6c\xd6\xbe\xa8\xb8\xc2\x4f\x0a\xc3\x7f\x52\xb1\x29\x0a\xc6\x25\x4e\xf4\x3b\x0b\xdf\x6e\x45\xd0\x6c\x46\x23\xd9\x65\xe0\xb0\x5a\x61\x01\x65\xa8\xd7\x4e\x17\x31\xeb\x25\xe3\xdd\x0e\xd3\x29\xed\xc1\x12\x51\x53\xb1\x04\x07\xd6\xe6\x29\x5b\xb7\x69\xef\xee\x30\x36\x43\x3f\x20\x9a\x64\x18\xf8\x7e\xfb\x3b\xb9\xf5\xa9\xa5\x61\x8d\xb7\x9a\x68\xa4\x8a\x5a\x2c\x74\x78\x16\x35\x26\x86\xb9\x33\x40\xce\xaa\xea\xde\x0d\xe8\x69\x41\xaa\x2d\x34\x11\x1a\x16\x3f\x14\x36\x05\xfe\xbe\x1e\x55\x0a\x7a\x8d\xd6\xfe\xe9\x97\x04\xf6\x06\xe0\xb6\x83\xbe\xbc\x5e\x0a\xb4\x60\xb3\xfd\xbe\x7e\xf2\xdc\xff\xb5\x60\x77\x30\xb0\xfe\xd7\xdd\xb5\x81\xa0\xb9\xf1\xfc\x12\x90\x96\x03\x83\xe3\x67\x86\xc6\x6f\x83\xe9\x98\x09\x75\x34\x9c\x03\xed\xe3\x67\x11\x76\x6c\x3b\x68\xd7\x48\x8d\x0d\x28\x6d\x34\x8e\x48\x1a\xed\x82\xfd\xfe\x63\x2b\xd8\x28\x99\xa9\x4a\xb6\xcd\x70\xec\xb2\xe5\xdc\xd6\xfe\xdf\x17\xae\xda\x94\x13\x2e\x5d\xa3\xee\xbc\xc4\xcd\xeb\xd8\x8b\x97\xb3\x80\x98\xaf\x43\x2f\x9d\xa3\xa8\x26\xd7\xef\x0d\x6b\x3b\xcd\xca\x12\x73\xac\x10\x40\x60\xff\xf8\x03\x8f\x5b\xe5\x91\xca\x04\xe7\xb5\x85\x2f\xe0\xcd\x6e\x37\x2d\x68\x55\x05\x8e\x52\x7f\xb7\x53\x3a\x6f\x10\x47\xb9\x08\x6f\xf5\x3a\x60\x5f\x29\x99\xc9\xa2\x52\xc1\x3d\x68\xf3\x31\x96\xcf\xca\x3f\xab\x35\xfc\x2b\x8a\xd7\x4b\xce\x36\x34\xf1\x83\xfe\x0b\x1e\x11\xb6\x1a\x3f\x13\xb9\x3a\x37\x3c\x7e\x2c\x9f\x67\xd0\xb1\xe0\x1c\x65\x19\xe6\x75\xd5\x0c\xc1\xe4\xf0\x1f\x40\xec\x90\xcf\x3d\x0b\x1a\x0f\xec\xb3\x9e\x45\x27\x61\x92\xe0\x14\x73\xad\xdc\x0f\xba\x95\x6b\xdb\x97\x5c\xd5\x39\x61\xc5\xdf\x20\xb9\x32\x4e\xf7\xcb\xd6\x36\x09\x44\x13\xf0\xf1\x37\xcb\xec\x79\x81\xfd\x45\xc1\xb3\x7f\x3b\xe8\x2e\x5b\x2d\xbb\x65\x59\x80\x37\xfb\xea\x7d\xf5\xbc\x21\x05\xbd\xa6\xa6\x3e\x98\x73\x58\x80\x08\x53\xc6\x73\x24\xcd\xc4\x6e\xe6\xf9\x75\xd1\x5b\x51\x5b\x3d\x55\x05\x4e\xaf\x73\x92\x40\xaf\x3d\x85\x3c\xbc\x3e\xab\x0f\x2b\xe4\xcc\xaa\x3f\xb8\x48\xeb\x09\x14\xb8\x9d\x5f\x31\xfd\xb1\x50\x1d\x6b\x70\x29\x71\x5a\xf0\x60\x6a\x76\x8b\xf4\x1c\xe5\x58\xd5\xe8\x40\x7a\xde\x71\x92\xff\x0d\x09\x69\xf3\xf4\x1d\x4d\xc4\x93\xce\xa6\x3c\x47\x55\xa5\x1c\x08\xf6\x5a\xe6\x7f\x03\x00\x00\xff\xff\x02\xe0\x15\x1d\xff\x1d\x00\x00")

func serviceTmplBytes() ([]byte, error) {
	return bindataRead(
		_serviceTmpl,
		"service.tmpl",
	)
}

func serviceTmpl() (*asset, error) {
	bytes, err := serviceTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "service.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x79, 0xc9, 0xca, 0x58, 0xc1, 0xaa, 0x11, 0xcb, 0x3d, 0x80, 0x3d, 0x4a, 0xdc, 0xd5, 0x30, 0x53, 0xc8, 0xa2, 0xc2, 0xb, 0x70, 0x89, 0xf5, 0x87, 0x8c, 0x19, 0x81, 0x52, 0xb, 0x9, 0x5, 0x51}}
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
	"function.tmpl":  functionTmpl,
	"info.tmpl":      infoTmpl,
	"open.tmpl":      openTmpl,
	"operation.tmpl": operationTmpl,
	"pair.tmpl":      pairTmpl,
	"service.tmpl":   serviceTmpl,
}

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
	"function.tmpl":  &bintree{functionTmpl, map[string]*bintree{}},
	"info.tmpl":      &bintree{infoTmpl, map[string]*bintree{}},
	"open.tmpl":      &bintree{openTmpl, map[string]*bintree{}},
	"operation.tmpl": &bintree{operationTmpl, map[string]*bintree{}},
	"pair.tmpl":      &bintree{pairTmpl, map[string]*bintree{}},
	"service.tmpl":   &bintree{serviceTmpl, map[string]*bintree{}},
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
