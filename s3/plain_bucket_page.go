package s3

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data" // translations between char sets

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type PlainBucketPage struct {
	baseURL string
	bucket  *PlainBucket

	// Optionally added to the baseURL
	maxKeys int
	lastKey string

	logTag string
	logger boshlog.Logger

	// Used by fetchOnce
	fetchedResult *plainBucketPage_ListBucketResult
	fetchedErr    error
}

func NewPlainBucketPage(baseURL string, bucket *PlainBucket, logger boshlog.Logger) *PlainBucketPage {
	if bucket == nil {
		panic("PlainBucketPage must be associated with a S3 Bucket object")
	}

	return &PlainBucketPage{
		baseURL: baseURL,
		bucket:  bucket,

		logTag: "PlainBucketPage",
		logger: logger,
	}
}

func (p *PlainBucketPage) Next() (*PlainBucketPage, error) {
	result, err := p.fetchOnce()
	if err != nil {
		return nil, err
	}

	if result.IsTruncated && len(result.Contents) > 0 {
		// Return page with next maxKeys/lastKeys offset
		nextPage := &PlainBucketPage{
			baseURL: p.baseURL,
			bucket:  p.bucket,

			maxKeys: result.MaxKeys,
			lastKey: result.Contents[len(result.Contents)-1].Key,

			logTag: "PlainBucketPage",
			logger: p.logger,
		}

		p.logger.Debug(p.logTag, "Next bucket page url '%s'", nextPage.fullURL())

		return nextPage, nil
	}

	p.logger.Debug(p.logTag, "No more bucket pages after url '%s'", p.fullURL())

	// This was the last bucket page
	return nil, nil
}

func (p *PlainBucketPage) Files() ([]File, error) {
	var files []File

	result, err := p.fetchOnce()
	if err != nil {
		return files, err
	}

	for _, c := range result.Contents {
		file := NewPlainFile(c.Key, c.ETag, c.Size, c.LastModified, p.bucket.url)
		files = append(files, file)
	}

	return files, nil
}

func (p *PlainBucketPage) fetchOnce() (plainBucketPage_ListBucketResult, error) {
	// Return memoized results; not thread safe
	if p.fetchedResult != nil || p.fetchedErr != nil {
		return *p.fetchedResult, p.fetchedErr
	}

	fullURL := p.fullURL()

	p.logger.Debug(p.logTag, "Fetching bucket page from url '%s'", fullURL)

	var result plainBucketPage_ListBucketResult

	resp, err := http.Get(fullURL)
	if err != nil {
		p.fetchedErr = err
		return result, bosherr.WrapError(err, "Requesting bucket page")
	}

	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReader

	err = decoder.Decode(&result)
	if err != nil {
		p.fetchedErr = err
		return result, bosherr.WrapError(err, "Parsing bucket page")
	}

	p.logger.Debug(p.logTag, "Saving bucket page result '%v'", result)

	p.fetchedResult = &result

	return result, nil
}

// fullURL returns URL of the specific page with proper key offset
func (p PlainBucketPage) fullURL() string {
	fullURL := p.baseURL + "?"

	if p.maxKeys > 0 {
		fullURL += fmt.Sprintf("max-keys=%d&", p.maxKeys)
	}

	if len(p.lastKey) > 0 {
		fullURL += fmt.Sprintf("marker=%s&", p.lastKey)
	}

	return fullURL
}

/*
Expected XML (ListBucketResult is the root of the document):

<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
  <Name>bosh-jenkins-artifacts</Name>
  <Prefix/>
  <Marker>release/bosh-1639.tgz</Marker>
  <MaxKeys>1000</MaxKeys>
  <IsTruncated>false</IsTruncated>
  <Contents>
    <Key>release/bosh-1644.tgz</Key>
    <LastModified>2013-12-25T07:01:23.000Z</LastModified>
    <ETag>"302e3c1f30d571efc150c37b7442aaea"</ETag>
    <Size>129786667</Size>
    <StorageClass>STANDARD</StorageClass>
  </Contents>
  ...
</ListBucketResult>

ListBucketResult.Marker is empty unless marker query param is specified.
*/

type plainBucketPage_ListBucketResult struct {
	XMLName xml.Name `xml:"ListBucketResult"`
	Name    string   `xml:"Name"`

	Prefix string `xml:"Prefix"`
	Marker string `xml:"Marker"`

	MaxKeys     int  `xml:"MaxKeys"`
	IsTruncated bool `xml:"IsTruncated"`

	Contents []plainBucketPage_Content `xml:"Contents"`
}

type plainBucketPage_Content struct {
	XMLName xml.Name `xml:"Contents"`
	Key     string   `xml:"Key"`

	ETag string `xml:"ETag"`
	Size uint64 `xml:"Size"`

	LastModified string `xml:"LastModified"`

	StorageClass string `xml:"StorageClass"`
}
