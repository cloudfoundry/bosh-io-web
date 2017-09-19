package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
	mart "github.com/go-martini/martini"
	martrend "github.com/martini-contrib/render"

	bhctrls "github.com/bosh-io/web/controllers"
)

const mainLogTag = "main"

var (
	debugOpt      = flag.Bool("debug", false, "Output debug logs")
	configPathOpt = flag.String("configPath", "", "Path to configuration file")
	assetsIDOpt   = flag.String("assetsID", "", "Assets ID value")
)

func main() {
	flag.Parse()

	logger, fs, runner, _ := basicDeps(*debugOpt)
	defer logger.HandlePanic("Main")

	config, err := NewConfigFromPath(*configPathOpt, fs)
	ensureNoErr(logger, "Loading config", err)

	repos, err := NewRepos(config.Repos, fs, logger)
	ensureNoErr(logger, "Failed building repos", err)

	controllerFactory, err := bhctrls.NewFactory(repos, runner, logger)
	ensureNoErr(logger, "Failed building controller factory", err)

	runControllers(controllerFactory, config.Analytics, logger)
}

func basicDeps(debug bool) (boshlog.Logger, boshsys.FileSystem, boshsys.CmdRunner, boshuuid.Generator) {
	logLevel := boshlog.LevelInfo

	// Debug generates a lot of log activity
	if debug {
		logLevel = boshlog.LevelDebug
	}

	logger := boshlog.NewWriterLogger(logLevel, os.Stderr)
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

	// Release viewing
	releasesController := controllerFactory.ReleasesController
	m.Get("/releases", releasesController.Index)
	m.Get("/releases/**", releasesController.Show)
	m.Get("/api/v1/releases/**", releasesController.APIV1Index)

	jobsController := controllerFactory.JobsController
	m.Get("/jobs/:name", jobsController.Show)

	packagesController := controllerFactory.PackagesController
	m.Get("/packages/:name", packagesController.Show)

	// Stemcell viewing
	stemcellsController := controllerFactory.StemcellsController
	m.Get("/stemcells", stemcellsController.Index)
	m.Get("/stemcells/**", stemcellsController.Index)
	m.Get("/d/stemcells/**", stemcellsController.Download)
	m.Get("/api/v1/stemcells/**", stemcellsController.APIV1Index)

	// ...make sure /d/** is after /d/stemcells/**
	releaseTarballsController := controllerFactory.ReleaseTarballsController
	m.Get("/d/**", releaseTarballsController.Download)

	m.Run()
}

func configureAssets(m *mart.ClassicMartini, analyticsConfig AnalyticsConfig, logger boshlog.Logger) {
	assetsID := strings.TrimSpace(*assetsIDOpt)

	if len(assetsID) == 0 {
		logger.Error(mainLogTag, "Expected non-empty assets ID")
		os.Exit(1)
	}

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

	htmlFuncs := template.FuncMap{
		"href": func(s string) template.HTMLAttr {
			return template.HTMLAttr(fmt.Sprintf(" href='%s' ", s))
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
			Funcs:      []template.FuncMap{assetsFuncs, analyticsConfigFuncs, htmlFuncs},
		},
	))
}
