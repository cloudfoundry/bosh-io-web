package s3

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strconv"
	"strings"
	"time"
)

type CDNSignedURLFactory struct {
	// Base should include transport
	base string

	keyPairID  string
	privateKey *rsa.PrivateKey
}

type CDNSignedURL struct {
	base string

	keyPairID  string
	privateKey *rsa.PrivateKey

	// Path should start with a slash
	path     string
	fileName string
}

func NewCDNURLFactory(base, keyPairID, privateKey string) (CDNSignedURLFactory, error) {
	p, _ := pem.Decode([]byte(privateKey))
	if p == nil {
		return CDNSignedURLFactory{}, errors.New("no PEM block found")
	}

	privateKey_, err := x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		return CDNSignedURLFactory{}, errors.New("no RSA private key found")
	}

	f := CDNSignedURLFactory{
		base: base,

		keyPairID:  keyPairID,
		privateKey: privateKey_,
	}

	return f, nil
}

func (f CDNSignedURLFactory) New(path, fileName string) URL {
	return CDNSignedURL{
		base: f.base,

		keyPairID:  f.keyPairID,
		privateKey: f.privateKey,

		path:     path,
		fileName: fileName,
	}
}

func (u CDNSignedURL) String() (string, error) {
	urlToSign := u.base + u.path + u.buildQuery()

	// 10 mins into the future
	dateLessThanSecs := strconv.FormatInt(time.Now().Unix()+60*10, 10)

	// No spaces allowed
	policyJSON := []byte("{\"Statement\":[{\"Resource\":\"" + urlToSign +
		"\",\"Condition\":{\"DateLessThan\":{\"AWS:EpochTime\":" + dateLessThanSecs + "}}}]}")

	_, policySignature, err := u.hashAndSignPolicy(policyJSON)
	if err != nil {
		return "", err
	}

	signedURL := urlToSign
	signedURL += "&Policy=" + u.awsEncode(policyJSON)
	signedURL += "&Signature=" + u.awsEncode(policySignature)
	signedURL += "&Key-Pair-Id=" + u.keyPairID

	return signedURL, nil
}

func (u CDNSignedURL) buildQuery() string {
	query := ""

	// Headers can be passed in via query
	if len(u.fileName) > 0 {
		query = "response-content-disposition=attachment%3B%20filename%3D" + u.fileName
	}

	return "?" + query
}

func (u CDNSignedURL) hashAndSignPolicy(policy []byte) ([]byte, []byte, error) {
	h := sha1.New()
	h.Write(policy)
	sum := h.Sum(nil)

	sig, err := rsa.SignPKCS1v15(rand.Reader, u.privateKey, crypto.SHA1, sum)
	if err != nil {
		return nil, nil, err
	}

	return sum, sig, nil
}

// amzEncode encodes data with base64 and special char replacements for AWS
func (u CDNSignedURL) awsEncode(bytes []byte) string {
	in := base64.StdEncoding.EncodeToString(bytes)
	in = strings.Replace(in, "+", "-", -1)
	in = strings.Replace(in, "=", "_", -1)
	return strings.Replace(in, "/", "~", -1)
}
