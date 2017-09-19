package vagrant

import (
	"fmt"
	"strings"
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bpeventlog "github.com/cppforlife/bosh-provisioner/eventlog"
)

const (
	aptDepsProvisionerLogTag            = "AptDepsProvisioner"
	aptDepsProvisionerUnableToFetchMsg  = "E: Unable to fetch some archives, maybe run apt-get update"
	aptDepsProvisionerUnableToLocateMsg = "E: Unable to locate package"
)

// AptDepsProvisioner installs basic dependencies for running
// packaging scripts from BOSH packages. It also installs
// non-captured dependencies by few common BOSH releases.
// (e.g. cmake, quota)
type AptDepsProvisioner struct {
	fullStemcellCompatibility bool

	runner   boshsys.CmdRunner
	eventLog bpeventlog.Log
	logger   boshlog.Logger
}

func NewAptDepsProvisioner(
	fullStemcellCompatibility bool,
	runner boshsys.CmdRunner,
	eventLog bpeventlog.Log,
	logger boshlog.Logger,
) AptDepsProvisioner {
	return AptDepsProvisioner{
		fullStemcellCompatibility: fullStemcellCompatibility,

		runner:   runner,
		eventLog: eventLog,
		logger:   logger,
	}
}

func (p AptDepsProvisioner) Provision() error {
	pkgNames := aptDepsProvisionerPkgsForMinimumStemcellCompatibility

	if p.fullStemcellCompatibility {
		pkgNames = append(pkgNames, aptDepsProvisionerPkgsForFullStemcellCompatibility...)
	}

	stage := p.eventLog.BeginStage("Installing dependencies", len(pkgNames))

	installedPkgNames, err := p.listInstalledPkgNames()
	if err != nil {
		return bosherr.WrapError(err, "Listing installed packages")
	}

	for _, pkgName := range pkgNames {
		task := stage.BeginTask(fmt.Sprintf("Package %s", pkgName))

		if p.isPkgInstalled(pkgName, installedPkgNames) {
			p.logger.Debug(aptDepsProvisionerLogTag, "Package %s is already installed", pkgName)
			task.End(nil)
			continue
		}

		err := task.End(p.installPkg(pkgName))
		if err != nil {
			return bosherr.WrapErrorf(err, "Installing %s", pkgName)
		}
	}

	return nil
}

func (p AptDepsProvisioner) InstallRunit() error {
	p.logger.Info(aptDepsProvisionerLogTag, "Installing runit")

	// todo non-bash
	cmd := boshsys.Command{
		Name: "bash",
		Args: []string{
			"-c", "apt-get -q -y -o Dpkg::Options::='--force-confdef' -o Dpkg::Options::='--force-confold' install runit",
		},
		Env: map[string]string{
			"DEBIAN_FRONTEND": "noninteractive",
		},
	}

	_, _, _, err := p.runner.RunComplexCommand(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (p AptDepsProvisioner) installPkg(name string) error {
	p.logger.Debug(aptDepsProvisionerLogTag, "Installing package %s", name)

	_, _, _, err := p.runner.RunCommand("apt-get", "-y", "install", name)
	if err == nil {
		return nil
	}

	unableToFetch := strings.Contains(err.Error(), aptDepsProvisionerUnableToFetchMsg)
	unableToLocate := strings.Contains(err.Error(), aptDepsProvisionerUnableToLocateMsg)

	// Avoid running 'apt-get update' since it usually takes 30sec
	if unableToFetch || unableToLocate {
		_, _, _, err := p.runner.RunCommand("apt-get", "-y", "update")
		if err != nil {
			return bosherr.WrapError(err, "Updating sources")
		}

		var lastInstallErr error

		// For some reason libssl-dev was really hard to install on the first try
		for i := 0; i < 3; i++ {
			_, _, _, lastInstallErr = p.runner.RunCommand("apt-get", "-y", "install", name)
			if lastInstallErr == nil {
				return nil
			}

			time.Sleep(1 * time.Second)
		}

		return bosherr.WrapErrorf(lastInstallErr, "Installing %s after updating", name)
	}

	return err
}

func (p AptDepsProvisioner) listInstalledPkgNames() ([]string, error) {
	var installedPkgNames []string

	installedPkgStdout, _, _, err := p.runner.RunCommand("dpkg", "--get-selections")
	if err != nil {
		return nil, bosherr.WrapError(err, "dkpg query")
	}

	// e.g. 'zlib1g:amd64 install'
	//      'util-linux   install'
	for _, line := range strings.Split(installedPkgStdout, "\n") {
		pieces := strings.Fields(line)

		// Last line is empty
		if len(pieces) == 2 && pieces[1] == "install" {
			pkgName := strings.Split(pieces[0], ":")[0]
			installedPkgNames = append(installedPkgNames, pkgName)
		}
	}

	return installedPkgNames, nil
}

func (p AptDepsProvisioner) isPkgInstalled(pkgName string, installedPkgs []string) bool {
	for _, installedPkgName := range installedPkgs {
		if installedPkgName == pkgName {
			return true
		}
	}

	return false
}

var aptDepsProvisionerPkgsForMinimumStemcellCompatibility = []string{
	// Most BOSH releases require it for packaging
	"build-essential", // 16sec
	"cmake",           // 6sec

	"libcap2-bin",
	"libcap-dev",

	"libbz2-1.0",   // noop on precise64 Vagrant box
	"libbz2-dev",   // 2sec
	"libxslt1-dev", // 2sec
	"libxml2-dev",  // 2sec

	// Used by BOSH Agent
	"iputils-arping",

	// For warden
	"quota", // 1sec

	// Started needing that in saucy for building BOSH
	"libssl-dev",

	"bison",
	"flex",

	"gettext",
	"libreadline6-dev",
	"libncurses5-dev",

	// Needed to render job templates
	"ruby1.9.3",
}

// Taken from base_apt stemcell builder stage
var aptDepsProvisionerPkgsForFullStemcellCompatibility = []string{
	"libaio1",
	"uuid-dev",
	"nfs-common",
	"zlib1g-dev",
	"apparmor-utils",
	"openssh-server",

	"libgcrypt-dev",
	"ca-certificates",

	// CURL
	"libcurl3",
	"libcurl3-dev",

	// XML
	"libxml2",
	"libxml2-dev",
	"libxslt1.1",
	"libxslt1-dev",

	// Utils
	"bind9-host",
	"dnsutils",
	"zip",
	"unzip",
	"psmisc",
	"lsof",
	"strace",
	"curl",
	"wget",
	"gdb",
	"sysstat",
	"rsync",

	"iptables",
	"tcpdump",
	"traceroute",
}
