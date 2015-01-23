package s3

type CDNFile struct {
	File

	urlFactory URLFactory
}

func NewCDNFile(file File, urlFactory URLFactory) File {
	return CDNFile{File: file, urlFactory: urlFactory}
}

func (f CDNFile) URL() (string, error) {
	return f.urlFactory.New("/"+f.Key(), "").String()
}
