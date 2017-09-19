package vagrant

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

const (
	runitProvisionerLogTag = "RunitProvisioner"
	runitStopStepDuration  = 1 * time.Second
)

var (
	// Matches 'svlogd -tt /var/vcap/bosh/log'
	runitSvlogdRegex = regexp.MustCompile(`\s*svlogd\s+\-tt\s+(.+)\s*`)

	// Matches 'down: agent: 3s, normally up; run: log: (pid 15318) 7762s'
	runitStatusDownRegex = regexp.MustCompile(`\Adown: [a-z\/]+: \d+`)

	runitPossibleRunsvdirPaths = []string{"/usr/sbin/runsvdir-start", "/sbin/runsvdir-start"}
)

// RunitProvisioner installs runit and
// adds specified service under runit's control.
type RunitProvisioner struct {
	fs              boshsys.FileSystem
	cmds            SimpleCmds
	depsProvisioner DepsProvisioner
	runner          boshsys.CmdRunner
	assetManager    AssetManager
	logger          boshlog.Logger
}

func NewRunitProvisioner(
	fs boshsys.FileSystem,
	cmds SimpleCmds,
	depsProvisioner DepsProvisioner,
	runner boshsys.CmdRunner,
	assetManager AssetManager,
	logger boshlog.Logger,
) RunitProvisioner {
	return RunitProvisioner{
		fs:              fs,
		cmds:            cmds,
		depsProvisioner: depsProvisioner,
		runner:          runner,
		assetManager:    assetManager,
		logger:          logger,
	}
}

func (p RunitProvisioner) Provision(name string, stopTimeout time.Duration) error {
	p.logger.Info(runitProvisionerLogTag, "Provisioning %s service", name)

	err := p.depsProvisioner.InstallRunit()
	if err != nil {
		return bosherr.WrapError(err, "Installing runit")
	}

	servicePath, enableServicePath := p.buildServicePaths(name)

	err = p.stopRunAndLog(servicePath, enableServicePath, name, stopTimeout)
	if err != nil {
		return bosherr.WrapError(err, "Stopping run and log")
	}

	err = p.setUpRun(servicePath, name)
	if err != nil {
		return bosherr.WrapError(err, "Setting up run")
	}

	err = p.setUpLog(servicePath, name)
	if err != nil {
		return bosherr.WrapError(err, "Setting up log")
	}

	err = p.startRunAndLog(servicePath, enableServicePath, name)
	if err != nil {
		return bosherr.WrapError(err, "Starting run and log")
	}

	return nil
}

func (p RunitProvisioner) Deprovision(name string, stopTimeout time.Duration) error {
	p.logger.Info(runitProvisionerLogTag, "Deprovisioning %s service", name)

	servicePath, enableServicePath := p.buildServicePaths(name)

	err := p.stopRunAndLog(servicePath, enableServicePath, name, stopTimeout)
	if err != nil {
		return bosherr.WrapError(err, "Stopping run and log")
	}

	return nil
}

func (p RunitProvisioner) buildServicePaths(name string) (string, string) {
	servicePath := fmt.Sprintf("/etc/sv/%s", name)
	enableServicePath := fmt.Sprintf("/etc/service/%s", name)
	return servicePath, enableServicePath
}

// setUpRun sets up script that runit will execute for the primary process
func (p RunitProvisioner) setUpRun(servicePath, name string) error {
	p.logger.Info(runitProvisionerLogTag, "Setting up %s run", name)

	err := p.cmds.MkdirP(servicePath)
	if err != nil {
		return err
	}

	runPath := fmt.Sprintf("%s/run", servicePath)

	err = p.assetManager.Place(fmt.Sprintf("%s/%s-run", name, name), runPath)
	if err != nil {
		return err
	}

	return p.cmds.ChmodX(runPath)
}

// setUpLog sets up logging destination for the service
func (p RunitProvisioner) setUpLog(servicePath, name string) error {
	p.logger.Info(runitProvisionerLogTag, "Setting up %s log", name)

	logPath := fmt.Sprintf("%s/log", servicePath)

	err := p.cmds.MkdirP(logPath)
	if err != nil {
		return err
	}

	logRunPath := fmt.Sprintf("%s/run", logPath)

	err = p.assetManager.Place(fmt.Sprintf("%s/%s-log", name, name), logRunPath)
	if err != nil {
		return err
	}

	err = p.cmds.ChmodX(logRunPath)
	if err != nil {
		return err
	}

	contens, err := p.fs.ReadFileString(logRunPath)
	if err != nil {
		return err
	}

	// First match is the whole string
	svlogdPaths := runitSvlogdRegex.FindStringSubmatch(contens)

	// Create log file destination so that runit process can properly log
	if len(svlogdPaths) == 2 {
		err = p.cmds.MkdirP(svlogdPaths[1])
		if err != nil {
			return err
		}
	}

	return nil
}

func (p RunitProvisioner) stopRunAndLog(servicePath, enableServicePath, name string, stopTimeout time.Duration) error {
	p.logger.Info(runitProvisionerLogTag, "Stopping %s run", name)

	err := p.stopRunsv(name, stopTimeout)
	if err != nil {
		return bosherr.WrapError(err, "Stopping service")
	}

	err = p.stopRunsv(fmt.Sprintf("%s/log", name), stopTimeout)
	if err != nil {
		return bosherr.WrapError(err, "Stopping log service")
	}

	err = p.fs.RemoveAll(enableServicePath)
	if err != nil {
		return err
	}

	// Clear out all service state kept in supervise/ and control/ dirs
	return p.fs.RemoveAll(servicePath)
}

func (p RunitProvisioner) startRunAndLog(servicePath, enableServicePath, name string) error {
	if name == "monit" {
		return p.configureMonitService(enableServicePath)
	}

	// Enabling service will kick in monitoring
	_, _, _, err := p.runner.RunCommand("ln", "-sf", servicePath, enableServicePath)
	return err
}

func (p RunitProvisioner) configureMonitService(enableServicePath string) error {
	var runsvdirPath string

	for _, path := range runitPossibleRunsvdirPaths {
		if p.fs.FileExists(path) {
			runsvdirPath = path
		}
	}

	if len(runsvdirPath) == 0 {
		return bosherr.Error("Failed to find runsvdir-start binary")
	}

	runsvdirStr, err := p.fs.ReadFileString(runsvdirPath)
	if err != nil {
		return bosherr.WrapErrorf(err, "Reading '%s'", runsvdirPath)
	}

	// Already checks sys/run directory; nothing to do
	if strings.Contains(runsvdirStr, "/var/vcap/data/sys/run") {
		return nil
	}

	// from https://github.com/cloudfoundry/bosh/blob/master/stemcell_builder/stages/delay_monit_start/apply.sh:
	// > if /var/vcap/data/sys/run is not already mounted, the agent must not have been started yet
	// > in that case remove /etc/services/monit in order to prevent runsvdir from starting monit
	// > (the agent will do that during boostrapping).
	sedScript := fmt.Sprintf(
		"sed -i '2i if [ x`mount | grep -c /var/vcap/data/sys/run` = x0 ] ; then rm -f %s ; fi' %s",
		enableServicePath,
		runsvdirPath,
	)

	_, _, _, err = p.runner.RunCommand("bash", "-c", sedScript)
	return err
}

func (p RunitProvisioner) stopRunsv(name string, stopTimeout time.Duration) error {
	p.logger.Info(runitProvisionerLogTag, "Stopping runsv")

	// Potentially tried to deprovision before ever provisioning
	if !p.runner.CommandExists("sv") {
		return nil
	}

	downStdout, _, _, err := p.runner.RunCommand("sv", "down", name)
	if err != nil {
		p.logger.Error(runitProvisionerLogTag, "Ignoring down error %s", err)
	}

	// If runsv configuration does not exist, service was never started
	if strings.Contains(downStdout, "file does not exist") {
		return nil
	}

	var lastStatusStdout string
	var passedDuration time.Duration

	for ; passedDuration < stopTimeout; passedDuration += runitStopStepDuration {
		lastStatusStdout, _, _, _ = p.runner.RunCommand("sv", "status", name)

		if runitStatusDownRegex.MatchString(lastStatusStdout) {
			return nil
		}

		time.Sleep(runitStopStepDuration)
	}

	return bosherr.Errorf("Failed to stop runsv for %s. Output: %s", name, lastStatusStdout)
}
