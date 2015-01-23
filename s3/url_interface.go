package s3

type URLFactory interface {
	New(path, fileName string) URL
}

type URL interface {
	String() (string, error)
}
