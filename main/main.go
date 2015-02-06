package main

import (
	"flag"
	// "net/http"
	// "net/http/pprof"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	boshuuid "github.com/cloudfoundry/bosh-agent/uuid"
	bpdload "github.com/cppforlife/bosh-provisioner/downloader"
	mart "github.com/go-martini/martini"
	martrend "github.com/martini-contrib/render"

	bhctrls "github.com/cppforlife/bosh-hub/controllers"
	bhfetcher "github.com/cppforlife/bosh-hub/release/fetcher"
	bhimporter "github.com/cppforlife/bosh-hub/release/importer"
	bhnoteimporter "github.com/cppforlife/bosh-hub/release/noteimporter"
	bhwatcher "github.com/cppforlife/bosh-hub/release/watcher"
	bhstemsimp "github.com/cppforlife/bosh-hub/stemcell/importer"
)

const mainLogTag = "main"

var (
	debugOpt      = flag.Bool("debug", false, "Output debug logs")
	configPathOpt = flag.String("configPath", "", "Path to configuration file")
)

func main() {
	flag.Parse()

	logger, fs, runner, uuidGen := basicDeps(*debugOpt)
	defer logger.HandlePanic("Main")

	config, err := NewConfigFromPath(*configPathOpt, fs)
	ensureNoErr(logger, "Loading config", err)

	repos, err := NewRepos(config.Repos, fs, logger)
	ensureNoErr(logger, "Failed building repos", err)

	controllerFactory, err := bhctrls.NewFactory(config.APIKey, repos, runner, logger)
	ensureNoErr(logger, "Failed building controller factory", err)

	downloader := bpdload.NewDefaultMuxDownloader(fs, runner, nil, logger)
	fetcher := bhfetcher.NewConcreteFetcher(fs, downloader, logger)

	{
		watcherFactory, err := bhwatcher.NewFactory(
			config.Watcher, repos, fetcher, logger)
		ensureNoErr(logger, "Failed building watcher factory", err)

		go watcherFactory.Watcher.Watch()
	}

	{
		importerFactory, err := bhimporter.NewFactory(
			config.Importer, repos, fetcher, fs, runner, downloader, uuidGen, logger)
		ensureNoErr(logger, "Failed building importer factory", err)

		go importerFactory.Importer.Import()
	}

	{
		noteImporterFactory, err := bhnoteimporter.NewFactory(config.NoteImporter, repos, logger)
		ensureNoErr(logger, "Failed building note importer factory", err)

		go noteImporterFactory.Importer.Import()
	}

	{
		stemcellImporterFactory := bhstemsimp.NewFactory(config.StemcellImporter, repos, logger)
		ensureNoErr(logger, "Failed building stemcell importer factory", err)

		go stemcellImporterFactory.Importer.Import()
	}

	if config.ActAsWorker {
		select {}
	} else {
		runControllers(controllerFactory, config.Analytics, logger)
	}
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

func runControllers(controllerFactory bhctrls.Factory, analyticsConfig AnalyticsConfig, logger boshlog.Logger) {
	m := mart.Classic()

	configureAssets(m, analyticsConfig, logger)

	homeController := controllerFactory.HomeController
	m.Get("/", homeController.Home)

	docsController := controllerFactory.DocsController
	m.Get("/docs", docsController.Page)
	m.Get("/docs/**", docsController.Page)

	// Watching release
	releaseWatchersController := controllerFactory.ReleaseWatchersController
	m.Get("/release_watchers", releaseWatchersController.Index)
	m.Post("/release_watchers", releaseWatchersController.WatchOrUnwatch) // actually Watch/Unwatch

	// Importing release
	releaseImportsController := controllerFactory.ReleaseImportsController
	m.Get("/release_imports", releaseImportsController.Index)
	m.Post("/release_imports", releaseImportsController.Delete) // actually Delete

	releaseImportErrsController := controllerFactory.ReleaseImportErrsController
	m.Get("/release_import_errs", releaseImportErrsController.Index)
	m.Post("/release_import_errs", releaseImportErrsController.Delete) // actually Delete

	// Release viewing
	releasesController := controllerFactory.ReleasesController
	m.Get("/releases", releasesController.Index)
	m.Get("/releases/**", releasesController.Show)

	releaseTarballsController := controllerFactory.ReleaseTarballsController
	m.Get("/d/**", releaseTarballsController.Download)

	jobsController := controllerFactory.JobsController
	m.Get("/jobs/:name", jobsController.Show)

	packagesController := controllerFactory.PackagesController
	m.Get("/packages/:name", packagesController.Show)

	// Stemcell viewing
	stemcellsController := controllerFactory.StemcellsController
	m.Get("/stemcells", stemcellsController.Index)

	// todo turn on based on config
	// m.Get("/debug/pprof", http.HandlerFunc(pprof.Index))
	// m.Get("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	// m.Get("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	// m.Get("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	// m.Get("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	// m.Get("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
	// m.Get("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	// m.Get("/debug/pprof/block", pprof.Handler("block").ServeHTTP)

	m.Run()
}

func configureAssets(m *mart.ClassicMartini, analyticsConfig AnalyticsConfig, logger boshlog.Logger) {
	assetsIDBytes, err := ioutil.ReadFile("./public/assets-id")
	ensureNoErr(logger, "Failed to find assets ID", err)

	assetsID := strings.TrimSpace(string(assetsIDBytes))

	assetsFuncs := template.FuncMap{
		"cssPath": func(fileName string) (string, error) {
			return "/" + assetsID + "/stylesheets/" + fileName, nil
		},
		"jsPath": func(fileName string) (string, error) {
			return "/" + assetsID + "/javascript/" + fileName, nil
		},
		"imgPath": func(fileName string) (string, error) {
			return "/" + assetsID + "/images/" + fileName, nil
		},
	}

	analyticsConfigFuncs := template.FuncMap{
		"analyticsConfig": func() AnalyticsConfig {
			return analyticsConfig
		},
	}

	// Use prefix to cache bust images, stylesheets, and js
	m.Use(mart.Static(
		"./public",
		mart.StaticOptions{
			Prefix: assetsID,
		},
	))

	// Make sure docs' images are available as `docs/images/X`
	m.Use(mart.Static(
		"./templates/docs/images",
		mart.StaticOptions{
			Prefix: "docs/images",
		},
	))

	m.Use(martrend.Renderer(
		martrend.Options{
			Layout:     "layout",
			Directory:  "./templates",
			Extensions: []string{".tmpl", ".html"},
			Funcs:      []template.FuncMap{assetsFuncs, analyticsConfigFuncs},
		},
	))
}
