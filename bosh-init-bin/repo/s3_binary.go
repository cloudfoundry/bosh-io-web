package repo

import (
	"fmt"
	"regexp"
	"strings"

	semiver "github.com/cppforlife/go-semi-semantic/version"
)

var (
	s3FileRegexp = regexp.MustCompile(`\Abosh-init-(?P<version>\d+\.\d+\.\d+)-(?P<platform>linux|darwin)-(?P<arch>amd64)\z`)

	prettyPlatformNames = map[string]string{
		"linux":  "Linux",
		"darwin": "Mac OS X",
	}
)

type S3Binary struct {
	name string

	version   semiver.Version
	updatedAt string

	size uint64
	etag string

	platform string // e.g. linux, darwin
	arch     string // e.g. amd64

	url string
}

func NewS3Binary(key, etag string, size uint64, lastModified, url string) *S3Binary {
	m := matchS3FileKey(key)

	if len(m) == 0 {
		return nil
	}

	version, err := semiver.NewVersionFromString(m["version"])
	if err != nil {
		return nil
	}

	binary := &S3Binary{
		name: key,

		version:   version,
		updatedAt: lastModified,

		size: size,
		etag: strings.Trim(etag, "\""),

		platform: m["platform"],
		arch:     m["arch"],

		url: url,
	}

	return binary
}

func (f S3Binary) Description() string {
	// bosh-init for Linux (amd64)
	// bosh-init for Mac OS X (amd64)
	// bosh-init is not supported on Windows
	return fmt.Sprintf("bosh-init for %s (%s)", prettyPlatformNames[f.platform], f.arch)
}

func (f S3Binary) Version() semiver.Version { return f.version }
func (f S3Binary) UpdatedAt() string        { return f.updatedAt }

func (f S3Binary) Size() uint64 { return f.size }
func (f S3Binary) MD5() string  { return f.etag }

func (f S3Binary) Platform() string { return f.platform }
func (f S3Binary) Arch() string     { return f.arch }

func (f S3Binary) URL() string { return f.url }

func matchS3FileKey(key string) map[string]string {
	match := s3FileRegexp.FindStringSubmatch(key)
	if match == nil {
		return nil
	}

	result := make(map[string]string)

	for i, name := range s3FileRegexp.SubexpNames() {
		if len(name) > 0 {
			result[name] = match[i]
		}
	}

	return result
}
