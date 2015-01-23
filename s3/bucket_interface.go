package s3

type Bucket interface {
	Files() ([]File, error)

	URL() string
}

type File interface {
	Key() string
	ETag() string

	Size() uint64
	LastModified() string

	URL() (string, error)
}
