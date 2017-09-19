package main

import (
	"flag"
	"os"

	boshblob "github.com/cloudfoundry/bosh-utils/blobstore"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"

	bpdep "github.com/cppforlife/bosh-provisioner/deployment"
	bpdload "github.com/cppforlife/bosh-provisioner/downloader"
	bpeventlog "github.com/cppforlife/bosh-provisioner/eventlog"
	bpinstance "github.com/cppforlife/bosh-provisioner/instance"
	bptplcomp "github.com/cppforlife/bosh-provisioner/instance/templatescompiler"
	bpinstupd "github.com/cppforlife/bosh-provisioner/instance/updater"
	bppkgscomp "github.com/cppforlife/bosh-provisioner/packagescompiler"
	bpprov "github.com/cppforlife/bosh-provisioner/provisioner"
	bprel "github.com/cppforlife/bosh-provisioner/release"
	bpreljob "github.com/cppforlife/bosh-provisioner/release/job"
	bptar "github.com/cppforlife/bosh-provisioner/tar"
	bpvagrantvm "github.com/cppforlife/bosh-provisioner/vm/vagrant"
)

const mainLogTag = "main"

var (
	configPathOpt = flag.String("configPath", "", "Path to configuration file")
)

func main() {
	logger, fs, runner, uuidGen := basicDeps()

	defer logger.HandlePanic("Main")

	config := mustLoadConfig(fs, logger)

	eventLogFactory := bpeventlog.NewFactory(config.EventLog, logger)

	eventLog := eventLogFactory.NewLog()

	mustSetTmpDir(config, fs, eventLog)

	mustCreateReposDir(config, fs, eventLog)

	localBlobstore := boshblob.NewLocalBlobstore(
		fs,
		uuidGen,
		config.Blobstore.Options,
	)

	blobstore := boshblob.NewSHA1VerifiableBlobstore(localBlobstore)

	downloader := bpdload.NewDefaultMuxDownloader(fs, runner, blobstore, logger)

	extractor := bptar.NewCmdExtractor(runner, fs, logger)

	compressor := bptar.NewCmdCompressor(runner, fs, logger)

	renderedArchivesCompiler := bptplcomp.NewRenderedArchivesCompiler(
		fs,
		runner,
		compressor,
		logger,
	)

	jobReaderFactory := bpreljob.NewReaderFactory(
		downloader,
		extractor,
		fs,
		logger,
	)

	reposFactory := NewReposFactory(config.ReposDir, fs, downloader, blobstore, logger)

	blobstoreProvisioner := bpprov.NewBlobstoreProvisioner(
		fs,
		config.Blobstore,
		logger,
	)

	err := blobstoreProvisioner.Provision()
	if err != nil {
		eventLog.WriteErr(bosherr.WrapError(err, "Provisioning blobstore"))
		os.Exit(1)
	}

	templatesCompiler := bptplcomp.NewConcreteTemplatesCompiler(
		renderedArchivesCompiler,
		jobReaderFactory,
		reposFactory.NewJobsRepo(),
		reposFactory.NewTemplateToJobRepo(),
		reposFactory.NewRuntimePackagesRepo(),
		reposFactory.NewTemplatesRepo(),
		blobstore,
		logger,
	)

	packagesCompilerFactory := bppkgscomp.NewConcretePackagesCompilerFactory(
		reposFactory.NewPackagesRepo(),
		reposFactory.NewCompiledPackagesRepo(),
		blobstore,
		eventLog,
		logger,
	)

	updaterFactory := bpinstupd.NewFactory(
		templatesCompiler,
		packagesCompilerFactory,
		eventLog,
		logger,
	)

	releaseReaderFactory := bprel.NewReaderFactory(
		downloader,
		extractor,
		fs,
		logger,
	)

	deploymentReaderFactory := bpdep.NewReaderFactory(fs, logger)

	vagrantVMProvisionerFactory := bpvagrantvm.NewVMProvisionerFactory(
		fs,
		runner,
		config.AssetsDir,
		config.Blobstore.AsMap(),
		config.VMProvisioner,
		eventLog,
		logger,
	)

	vagrantVMProvisioner := vagrantVMProvisionerFactory.NewVMProvisioner()

	releaseCompiler := bpprov.NewReleaseCompiler(
		releaseReaderFactory,
		packagesCompilerFactory,
		templatesCompiler,
		vagrantVMProvisioner,
		eventLog,
		logger,
	)

	instanceProvisioner := bpinstance.NewProvisioner(
		updaterFactory,
		logger,
	)

	singleVMProvisionerFactory := bpprov.NewSingleVMProvisionerFactory(
		deploymentReaderFactory,
		config.DeploymentProvisioner,
		vagrantVMProvisioner,
		releaseCompiler,
		instanceProvisioner,
		eventLog,
		logger,
	)

	deploymentProvisioner := singleVMProvisionerFactory.NewSingleVMProvisioner()

	err = deploymentProvisioner.Provision()
	if err != nil {
		eventLog.WriteErr(bosherr.WrapError(err, "Provisioning deployment"))
		os.Exit(1)
	}
}

func basicDeps() (boshlog.Logger, boshsys.FileSystem, boshsys.CmdRunner, boshuuid.Generator) {
	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr, os.Stderr)

	fs := boshsys.NewOsFileSystem(logger)

	runner := boshsys.NewExecCmdRunner(logger)

	uuidGen := boshuuid.NewGenerator()

	return logger, fs, runner, uuidGen
}

func mustLoadConfig(fs boshsys.FileSystem, logger boshlog.Logger) Config {
	flag.Parse()

	config, err := NewConfigFromPath(*configPathOpt, fs)
	if err != nil {
		logger.Error(mainLogTag, "Failed to load config %s", err)
		os.Exit(1)
	}

	return config
}

func mustSetTmpDir(config Config, fs boshsys.FileSystem, eventLog bpeventlog.Log) {
	// todo leaky abstraction?
	if len(config.TmpDir) == 0 {
		return
	}

	err := fs.MkdirAll(config.TmpDir, os.ModeDir)
	if err != nil {
		eventLog.WriteErr(bosherr.WrapError(err, "Creating tmp dir"))
		os.Exit(1)
	}

	err = os.Setenv("TMPDIR", config.TmpDir)
	if err != nil {
		eventLog.WriteErr(bosherr.WrapError(err, "Setting TMPDIR"))
		os.Exit(1)
	}
}

func mustCreateReposDir(config Config, fs boshsys.FileSystem, eventLog bpeventlog.Log) {
	err := fs.MkdirAll(config.ReposDir, os.ModeDir)
	if err != nil {
		eventLog.WriteErr(bosherr.WrapError(err, "Creating repos dir"))
		os.Exit(1)
	}
}
