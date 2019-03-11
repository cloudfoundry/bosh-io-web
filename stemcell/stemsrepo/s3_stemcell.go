package stemsrepo

import (
	"regexp"
	"strings"

	semiver "github.com/cppforlife/go-semi-semantic/version"

	bhnotesrepo "github.com/bosh-io/web/stemcell/notesrepo"
)

var (
	s3StemcellAgentRegexp = regexp.MustCompile(`ruby|go|agent`)
	s3StemcellRegexp      = regexp.MustCompile(`\A(([\w-]+/)?\w+/)?(?P<flavor>[\w-]+)-stemcell-(?P<version>[\.\d]+)-(?P<name>(?P<inf_name>\w+)-(?P<hv_name>\w+(-\w+)?)-(?P<os_name>centos|ubuntu|windows)(?P<os_version>-trusty|-xenial|-lucid|2019|1803|2016|2012R2|-\d+)?(?P<agent_type>-go_agent)?(?P<disk_fmt>-raw)?)\.tgz\z`)

	// Previous verisons derived checksums from other locations instead of DB
	minLinuxChecksumedVersion, _   = semiver.NewVersionFromString("3262.2")
	minWindowsChecksumedVersion, _ = semiver.NewVersionFromString("1079.0")
)

type S3Stemcell struct {
	notesRepo bhnotesrepo.NotesRepository

	name   string
	flavor string // e.g. bosh vs light-bosh

	version   semiver.Version
	updatedAt string

	size uint64
	etag string
	sha1 string

	infName    string // e.g. aws
	hvName     string // e.g. kvm
	diskFormat string // e.g. raw

	osName    string // e.g. Ubuntu
	osVersion string // e.g. Trusty

	agentType string // e.g. Ruby

	url string
}

func NewS3Stemcell(key, etag, sha1 string, size uint64, lastModified, url string) *S3Stemcell {
	m := matchS3FileKey(key)

	if len(m) == 0 {
		return nil
	}

	version, err := semiver.NewVersionFromString(m["version"])
	if err != nil {
		return nil
	}

	var osName, osVersion, agentType string

	osName = m["os_name"]

	if len(m["os_version"]) > 0 {
		osVersion = strings.Trim(m["os_version"], "-")
	}

	if len(m["agent_type"]) > 0 {
		agentType = strings.Trim(m["agent_type"], "-")
	} else {
		agentType = "ruby_agent"
	}

	if s3StemcellAgentRegexp.MatchString(osVersion) {
		if osName == "ubuntu" {
			agentType = osVersion
			osVersion = "lucid"
		} else {
			agentType = osVersion
			osVersion = ""
		}
	}

	if len(osVersion) == 0 && osName == "ubuntu" {
		osVersion = "lucid"
	}

	s3Stemcell := &S3Stemcell{
		// todo assume that piece of the stemcell file name
		// matches actual stemcell name used in a manifest
		name:   "bosh-" + m["name"],
		flavor: m["flavor"],

		version:   version,
		updatedAt: lastModified,

		size: size,
		etag: strings.Trim(etag, "\""),
		sha1: sha1,

		infName:    m["inf_name"],
		hvName:     m["hv_name"],
		diskFormat: strings.Trim(m["disk_fmt"], "-"),

		osName:    osName,
		osVersion: osVersion,

		agentType: strings.Replace(agentType, "_agent", "", 1),

		url: url,
	}

	return s3Stemcell
}

func (f S3Stemcell) Name() string { return f.name }

func (f S3Stemcell) Version() semiver.Version { return f.version }
func (f S3Stemcell) UpdatedAt() string        { return f.updatedAt }

func (f S3Stemcell) Size() uint64 { return f.size }
func (f S3Stemcell) MD5() string  { return f.etag }
func (f S3Stemcell) SHA1() string { return f.sha1 }

func (f S3Stemcell) InfName() string    { return f.infName }
func (f S3Stemcell) HvName() string     { return f.hvName }
func (f S3Stemcell) DiskFormat() string { return f.diskFormat }

func (f S3Stemcell) OSName() string    { return f.osName }
func (f S3Stemcell) OSVersion() string { return f.osVersion }

func (f S3Stemcell) AgentType() string { return f.agentType }

func (f S3Stemcell) IsLight() bool {
	return strings.Index(f.flavor, "light-") == 0
}

func (f S3Stemcell) IsForChina() bool {
	return strings.Index(f.flavor, "-china-") != -1
}

func (f S3Stemcell) IsDeprecated() bool {
	// softlayer actually uses xen stemcells
	if f.name == "bosh-softlayer-esxi-ubuntu-trusty-go_agent" {
		return true
	}

	return f.osVersion == "lucid" || f.agentType == "ruby"
}

func (f S3Stemcell) URL() string { return f.url }

func (f S3Stemcell) MustHaveSHA1() bool {
	if f.osName == "windows" {
		return f.version.IsGt(minWindowsChecksumedVersion)
	}

	return f.version.IsGt(minLinuxChecksumedVersion)
}

func (f S3Stemcell) Notes() (bhnotesrepo.NoteRec, bool, error) {
	return f.notesRepo.Find(f.Version().String())
}

func matchS3FileKey(key string) map[string]string {
	match := s3StemcellRegexp.FindStringSubmatch(key)
	if match == nil {
		return nil
	}

	result := make(map[string]string)

	for i, name := range s3StemcellRegexp.SubexpNames() {
		if len(name) > 0 {
			result[name] = match[i]
		}
	}

	return result
}
