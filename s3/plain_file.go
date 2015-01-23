package s3

import (
	"strings"
)

type PlainFile struct {
	key  string
	eTag string

	size         uint64
	lastModified string

	baseURL string
}

func NewPlainFile(key, eTag string, size uint64, lastModified string, baseURL string) File {
	return PlainFile{
		key:  key,
		eTag: eTag,

		size:         size,
		lastModified: lastModified,

		baseURL: baseURL,
	}
}

func (f PlainFile) Key() string  { return f.key }
func (f PlainFile) ETag() string { return f.eTag }

func (f PlainFile) Size() uint64         { return f.size }
func (f PlainFile) LastModified() string { return f.lastModified }

// URL returns full URL of the S3 object
func (f PlainFile) URL() (string, error) {
	// do not use path since it will collapse http://
	return strings.TrimSuffix(f.baseURL, "/") + "/" + f.key, nil
}
