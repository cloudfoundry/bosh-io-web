package releasetarsrepo

import (
	"encoding/xml"
	"time"
)

type Metalink struct {
	XMLName   xml.Name   `xml:"urn:ietf:params:xml:ns:metalink metalink" json:"-" yaml:"-"`
	Files     []File     `xml:"file" json:"files,omitempty" yaml:"files,omitempty"`
	Generator string     `xml:"generator,,omitempty" json:"generator,omitempty" yaml:"generator,omitempty"`
	Origin    *Origin    `xml:"origin,,omitempty" json:"origin,omitempty" yaml:"origin,omitempty"`
	Published *time.Time `xml:"published,,omitempty" json:"published,omitempty" yaml:"published,omitempty"`
	Updated   *time.Time `xml:"updated,,omitempty" json:"updated,omitempty" yaml:"updated,omitempty"`
}

type Origin struct {
	XMLName xml.Name `xml:"origin" json:"-" yaml:"-"`
	Dynamic *bool    `xml:"dynamic,attr,omitempty" json:"dynamic,omitempty" yaml:"dynamic,omitempty"`
	URL     string   `xml:",chardata" json:"url" yaml:"url"`
}

type File struct {
	XMLName     xml.Name   `xml:"file" json:"-" yaml:"-"`
	Name        string     `xml:"name,attr" json:"name" yaml:"name"`
	Copyright   string     `xml:"copyright,,omitempty" json:"copyright,omitempty" yaml:"copyright,omitempty"`
	Description string     `xml:"description,,omitempty" json:"description,omitempty" yaml:"description,omitempty"`
	Hashes      []Hash     `xml:"hash,,omitempty" json:"hashes,omitempty" yaml:"hashes,omitempty"`
	Identity    string     `xml:"identity,,omitempty" json:"identity,omitempty" yaml:"identity,omitempty"`
	Language    []string   `xml:"language,,omitempty" json:"language,omitempty" yaml:"language,omitempty"`
	Logo        string     `xml:"logo,,omitempty" json:"logo,omitempty" yaml:"logo,omitempty"`
	MetaURLs    []MetaURL  `xml:"metaurl,,omitempty" json:"metaurl,omitempty" yaml:"metaurl,omitempty"`
	OS          []string   `xml:"os,,omitempty" json:"os,omitempty" yaml:"os,omitempty"`
	Pieces      []Piece    `xml:"pieces,,omitempty" json:"piece,omitempty" yaml:"piece,omitempty"`
	Publisher   *Publisher `xml:"publisher" json:"publisher,omitempty" yaml:"publisher,omitempty"`
	Signature   *Signature `xml:"signature" json:"signature,omitempty" yaml:"signature,omitempty"`
	Size        uint64     `xml:"size,,omitempty" json:"size,omitempty" yaml:"size,omitempty"`
	URLs        []URL      `xml:"url,,omitempty" json:"url,omitempty" yaml:"url,omitempty"`
	Version     string     `xml:"version,omitempty" json:"version,omitempty" yaml:"version,omitempty"`
}

type URL struct {
	XMLName  xml.Name `xml:"url" json:"-" yaml:"-"`
	Location string   `xml:"location,attr,omitempty" json:"location,omitempty" yaml:"location,omitempty"`
	Priority *uint    `xml:"priority,attr,omitempty" json:"priority,omitempty" yaml:"priority,omitempty"`
	URL      string   `xml:",chardata" json:"url" yaml:"url"`
}

type Signature struct {
	XMLName   xml.Name `xml:"signature" json:"-" yaml:"-"`
	MediaType string   `xml:"mediatype,attr" json:"mediatype" yaml:"mediatype"`
	Signature string   `xml:",cdata" json:"signature" yaml:"signature"`
}

type Publisher struct {
	XMLName xml.Name `xml:"publisher" json:"-" yaml:"-"`
	Name    string   `xml:"name,attr" json:"name" yaml:"name"`
	URL     string   `xml:"url,attr,omitempty" json:"url,omitempty" yaml:"url,omitempty"`
}

type Hash struct {
	XMLName xml.Name `xml:"hash" json:"-" yaml:"-"`
	Type    string   `xml:"type,attr" json:"type" yaml:"type"`
	Hash    string   `xml:",chardata" json:"hash" yaml:"hash"`
}

type Piece struct {
	XMLName xml.Name `xml:"pieces" json:"-" yaml:"-"`
	Type    string   `xml:"type,attr" json:"type" yaml:"type"`
	Length  string   `xml:"length,attr" json:"length" yaml:"length"`
	Hash    []string `xml:"hash,chardata" json:"hash" yaml:"hash"`
}

type MetaURL struct {
	XMLName   xml.Name `xml:"metaurl" json:"-" yaml:"-"`
	Priority  *uint    `xml:"priority,attr,omitempty" json:"priority,omitempty" yaml:"priority,omitempty"`
	MediaType string   `xml:"mediatype,attr" json:"mediatype" yaml:"mediatype"`
	Name      string   `xml:"name,attr,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
	URL       string   `xml:",chardata" json:"url" yaml:"url"`
}

type Extra_ struct {
	XMLName xml.Name `json:"-" yaml:"-"`
	Data    string   `xml:",innerxml"`
}