package main

import (
	"flag"
	"fmt"
	"os"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	boshuuid "github.com/cloudfoundry/bosh-agent/uuid"

	// bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

const mainLogTag = "main"

var (
	debugOpt      = flag.Bool("debug", false, "Output debug logs")
	configPathOpt = flag.String("configPath", "", "Path to configuration file")
)

func main() {
	flag.Parse()

	logger, fs, _, _ := basicDeps(*debugOpt)
	defer logger.HandlePanic("Main")

	config, err := NewConfigFromPath(*configPathOpt, fs)
	ensureNoErr(logger, "Loading config", err)

	repos, err := NewRepos(config.Repos, fs, logger)
	ensureNoErr(logger, "Failed building repos", err)

	customRun(repos, logger)
}

func customRun(repos Repos, logger boshlog.Logger) {
	releasesRepo := repos.ReleasesRepo()
	releaseTarsRepo := repos.ReleaseTarsRepo()

	sources, err := releasesRepo.ListAll()
	ensureNoErr(logger, "ListAll", err)

	var totalVers int

	for _, src := range sources {
		relVerRecs, _, err := releasesRepo.FindAll(src.Full)
		ensureNoErr(logger, "FindAll", err)

		vers := map[string]struct{}{}

		for _, relVerRec := range relVerRecs {
			vers[relVerRec.VersionRaw] = struct{}{}
		}

		if len(relVerRecs) != len(vers) {
			fmt.Printf("-----> %s: %d -> %d\n", src.Full, len(relVerRecs), len(vers))

			err := releasesRepo.RemoveDups(src.Full)
			ensureNoErr(logger, "RemoveDups", err)
		}

		totalVers += len(vers)
	}

	fmt.Printf("Total non-dup versions: %d\n", totalVers)

	keys, err := releaseTarsRepo.GetAll()
	ensureNoErr(logger, "GetAll", err)

	fmt.Printf("Total release tarballs: %d\n", len(keys))
}

func basicDeps(debug bool) (boshlog.Logger, boshsys.FileSystem, boshsys.CmdRunner, boshuuid.Generator) {
	logLevel := boshlog.LevelInfo

	// Debug generates a lot of log activity
	if debug {
		logLevel = boshlog.LevelDebug
	}

	logger := boshlog.NewWriterLogger(logLevel, os.Stderr, os.Stderr)
	fs := boshsys.NewOsFileSystem(logger)
	runner := boshsys.NewExecCmdRunner(logger)
	uuidGen := boshuuid.NewGenerator()
	return logger, fs, runner, uuidGen
}

func ensureNoErr(logger boshlog.Logger, errPrefix string, err error) {
	if err != nil {
		logger.Error(mainLogTag, "%s: %s", errPrefix, err)
		os.Exit(1)
	}
}
