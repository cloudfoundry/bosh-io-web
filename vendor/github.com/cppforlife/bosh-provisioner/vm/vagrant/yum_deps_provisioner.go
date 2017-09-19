package vagrant

import (
	"fmt"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bpeventlog "github.com/cppforlife/bosh-provisioner/eventlog"
)

const (
	yumDepsProvisionerLogTag = "YumDepsProvisioner"
)

// YumDepsProvisioner installs basic dependencies for running
// packaging scripts from BOSH packages. It also installs
// non-captured dependencies by few common BOSH releases.
// (e.g. cmake, quota)
type YumDepsProvisioner struct {
	fullStemcellCompatibility bool

	runner   boshsys.CmdRunner
	eventLog bpeventlog.Log
	logger   boshlog.Logger
}

func NewYumDepsProvisioner(
	fullStemcellCompatibility bool,
	runner boshsys.CmdRunner,
	eventLog bpeventlog.Log,
	logger boshlog.Logger,
) YumDepsProvisioner {
	return YumDepsProvisioner{
		fullStemcellCompatibility: fullStemcellCompatibility,

		runner:   runner,
		eventLog: eventLog,
		logger:   logger,
	}
}

func (p YumDepsProvisioner) Provision() error {
	groupNames := []string{"Base", "Development Tools"}

	pkgNames := yumDepsProvisionerPkgsForMinimumStemcellCompatibility

	if p.fullStemcellCompatibility {
		pkgNames = append(pkgNames, yumDepsProvisionerPkgsForFullStemcellCompatibility...)
	}

	stage := p.eventLog.BeginStage("Installing dependencies", len(groupNames)+len(pkgNames))

	for _, groupName := range groupNames {
		task := stage.BeginTask(fmt.Sprintf("Group %s", groupName))

		_, _, _, err := p.runner.RunCommand("yum", "--assumeyes", "groupinstall", groupName)
		if task.End(err) != nil {
			return err
		}
	}

	installedPkgNames, err := p.listInstalledPkgNames()
	if err != nil {
		return bosherr.WrapError(err, "Listing installed packages")
	}

	for _, pkgName := range pkgNames {
		task := stage.BeginTask(fmt.Sprintf("Package %s", pkgName))

		if p.isPkgInstalled(pkgName, installedPkgNames) {
			p.logger.Debug(yumDepsProvisionerLogTag, "Package %s is already installed", pkgName)
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

func (p YumDepsProvisioner) InstallRunit() error {
	p.logger.Info(yumDepsProvisionerLogTag, "Installing runit")

	installedPkgNames, err := p.listInstalledPkgNames()
	if err != nil {
		return bosherr.WrapError(err, "Listing installed packages")
	}

	if !p.isPkgInstalled("runit", installedPkgNames) {
		cmds := [][]string{
			{"curl", "-L", "https://github.com/opscode-cookbooks/runit/archive/v1.2.0.tar.gz", "-o", "/tmp/v1.2.0.tar.gz"},
			{"tar", "-C", "/tmp", "-xvf", "/tmp/v1.2.0.tar.gz"},
			{"tar", "-C", "/tmp", "-xvf", "/tmp/runit-1.2.0/files/default/runit-2.1.1.tar.gz"},
			{"yum", "--assumeyes", "install", "rpmdevtools"},
			{"/tmp/runit-2.1.1/build.sh"},
			{"rpm", "-i", "/root/rpmbuild/RPMS/runit-2.1.1.rpm"},
		}

		for _, cmd := range cmds {
			_, _, _, err := p.runner.RunCommand(cmd[0], cmd[1:]...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p YumDepsProvisioner) installPkg(name string) error {
	p.logger.Debug(yumDepsProvisionerLogTag, "Installing package %s", name)

	_, _, _, err := p.runner.RunCommand("yum", "--assumeyes", "install", name)
	if err == nil {
		return nil
	}

	return err
}

func (p YumDepsProvisioner) listInstalledPkgNames() ([]string, error) {
	var installedPkgNames []string

	installedPkgStdout, _, _, err := p.runner.RunCommand("yum", "list", "installed")
	if err != nil {
		return nil, bosherr.WrapError(err, "Listing yum packages")
	}

	var packagesStarted bool

	// $ yum list installed
	// Loaded plugins: fastestmirror, langpacks
	// Loading mirror speeds from cached hostfile
	//  * base: mirror.tocici.com
	//  * extras: mirror.tocici.com
	//  * updates: mirror.keystealth.org
	// Installed Packages
	// ModemManager-glib.x86_64      1.1.0-6.git20130913.el7                    @anaconda
	// NetworkManager.x86_64         1:0.9.9.1-25.git20140326.4dba720.el7_0     @updates
	// ...
	for _, line := range strings.Split(installedPkgStdout, "\n") {
		if line == "Installed Packages" {
			packagesStarted = true
		} else if !packagesStarted {
			continue
		}

		pieces := strings.Fields(line)

		// Last line is empty
		if len(pieces) == 3 {
			pkgName := strings.Split(pieces[0], ".")[0]
			installedPkgNames = append(installedPkgNames, pkgName)
		}
	}

	return installedPkgNames, nil
}

func (p YumDepsProvisioner) isPkgInstalled(pkgName string, installedPkgs []string) bool {
	for _, installedPkgName := range installedPkgs {
		if installedPkgName == pkgName {
			return true
		}
	}

	return false
}

var yumDepsProvisionerPkgsForMinimumStemcellCompatibility = []string{
	"glibc-static",

	// "libcap2-bin", // todo needed
	"libcap-devel",

	// Used by BOSH Agent
	"iputils",

	// For warden
	"quota",

	"openssl-devel",

	"bison",
	"flex",

	"gettext",
	"readline-devel",
	"ncurses-devel",

	// Needed to render job templates
	"ruby",
}

// Taken from base_apt stemcell builder stage
var yumDepsProvisionerPkgsForFullStemcellCompatibility = []string{
	"libaio",
	"libuuid-devel",
	// "nfs-common", // todo
	// "apparmor-utils",
	"openssh-server",

	"libgcrypt-devel",
	"ca-certificates",

	// CURL
	// "libcurl3", // todo
	// "libcurl3-dev",

	// XML
	"libxml2",
	"libxml2-devel",
	"libxslt",
	"libxslt-devel",

	// Utils
	// "bind9-host", // todo
	// "dnsutils",
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
