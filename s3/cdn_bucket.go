package s3

type CDNBucket struct {
	urlFactory URLFactory
	bucket     Bucket
}

func NewCDNBucket(urlFactory URLFactory, bucket Bucket) CDNBucket {
	return CDNBucket{urlFactory, bucket}
}

func (b CDNBucket) Files() ([]File, error) {
	files, err := b.bucket.Files()
	if err != nil {
		return nil, err
	}

	cdnFiles := []File{}

	for _, file := range files {
		cdnFile := NewCDNFile(file, b.urlFactory)
		cdnFiles = append(cdnFiles, cdnFile)
	}

	return cdnFiles, nil
}

func (b CDNBucket) URL() string { return b.bucket.URL() }
