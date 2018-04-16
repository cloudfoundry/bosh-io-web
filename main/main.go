package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"

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
)

func main() {
	flag.Parse()

	logger, fs, runner, _ := basicDeps(*debugOpt)
	defer logger.HandlePanic("Main")

	config, err := NewConfigFromPath(*configPathOpt, fs)
	ensureNoErr(logger, "Loading config", err)

	repos, err := NewRepos(config.Repos, fs, logger)
	ensureNoErr(logger, "Failed building repos", err)

	redirects, err := LoadRedirects(fs)
	ensureNoErr(logger, "Failed loading redirects", err)

	controllerFactory, err := bhctrls.NewFactory(redirects, repos, runner, logger)
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

	m.Use(controllerFactory.RedirectsController.ServeHTTP)
	m.NotFound(controllerFactory.NotFoundController.ServeHTTP)

	m.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusTemporaryRedirect)
	})

	m.Get("/docs/**",
		mart.Static("templates/docs", mart.StaticOptions{Prefix: "/docs"}),
		controllerFactory.NotFoundController.ServeHTTP,
	)

	{
		// Stemcell viewing
		stemcellsController := controllerFactory.StemcellsController
		m.Get("/stemcells", stemcellsController.Index)
		m.Get("/stemcells/**", stemcellsController.Index)

		m.Get("/d/stemcells/**", stemcellsController.Download)

		m.Get("/api/v1/stemcells/**", stemcellsController.APIV1Index)
	}

	{
		// Release viewing
		releasesController := controllerFactory.ReleasesController
		m.Get("/releases", releasesController.Index)
		m.Get("/releases/**", releasesController.Show)

		jobsController := controllerFactory.JobsController
		m.Get("/jobs/:name", jobsController.Show)

		packagesController := controllerFactory.PackagesController
		m.Get("/packages/:name", packagesController.Show)

		// ...make sure /d/** is after /d/stemcells/**
		releaseTarballsController := controllerFactory.ReleaseTarballsController
		m.Get("/d/**", releaseTarballsController.Download)

		m.Get("/api/v1/releases/**", releasesController.APIV1Index)
	}

	m.Run()
}

func configureAssets(m *mart.ClassicMartini, analyticsConfig AnalyticsConfig, logger boshlog.Logger) {
	themeStylesheetApplication, err := filepath.Glob("templates/docs/assets/stylesheets/application.*.css")
	if err != nil {
		logger.Error(mainLogTag, fmt.Sprintf("Failed to find application.*.css: %v", err))
		os.Exit(1)
	} else if len(themeStylesheetApplication) != 1 {
		logger.Error(mainLogTag, fmt.Sprintf("Failed to find exactly one application.*.css: found %d", len(themeStylesheetApplication)))
		os.Exit(1)
	}

	themeStylesheetApplicationPalette, err := filepath.Glob("templates/docs/assets/stylesheets/application-palette.*.css")
	if err != nil {
		logger.Error(mainLogTag, fmt.Sprintf("Failed to find application-palette.*.css: %v", err))
		os.Exit(1)
	} else if len(themeStylesheetApplicationPalette) != 1 {
		logger.Error(mainLogTag, fmt.Sprintf("Failed to find exactly one application-palette.*.css: found %d", len(themeStylesheetApplicationPalette)))
		os.Exit(1)
	}

	themeStylesheetExtra, err := filepath.Glob("templates/docs/assets/stylesheets/extra.*.css")
	if err != nil {
		logger.Error(mainLogTag, fmt.Sprintf("Failed to find extra.*.css: %v", err))
		os.Exit(1)
	} else if len(themeStylesheetExtra) != 1 {
		logger.Error(mainLogTag, fmt.Sprintf("Failed to find exactly one extra.*.css: found %d", len(themeStylesheetExtra)))
		os.Exit(1)
	}

	themeJavascriptModernizr, err := filepath.Glob("templates/docs/assets/javascripts/modernizr.*.js")
	if err != nil {
		logger.Error(mainLogTag, fmt.Sprintf("Failed to find modernizr.*.js: %v", err))
		os.Exit(1)
	} else if len(themeJavascriptModernizr) != 1 {
		logger.Error(mainLogTag, fmt.Sprintf("Failed to find exactly one modernizr.*.js: found %d", len(themeJavascriptModernizr)))
		os.Exit(1)
	}

	themeJavascriptApplication, err := filepath.Glob("templates/docs/assets/javascripts/application.*.js")
	if err != nil {
		logger.Error(mainLogTag, fmt.Sprintf("Failed to find application.*.js: %v", err))
		os.Exit(1)
	} else if len(themeJavascriptApplication) != 1 {
		logger.Error(mainLogTag, fmt.Sprintf("Failed to find exactly one application.*.js: found %d", len(themeJavascriptApplication)))
		os.Exit(1)
	}

	themeImageFavicon, err := filepath.Glob("templates/docs/assets/images/favicon.*.png")
	if err != nil {
		logger.Error(mainLogTag, fmt.Sprintf("Failed to find favicon.*.png: %v", err))
		os.Exit(1)
	} else if len(themeImageFavicon) != 1 {
		logger.Error(mainLogTag, fmt.Sprintf("Failed to find exactly one favicon.*.png: found %d", len(themeImageFavicon)))
		os.Exit(1)
	}

	themeImageLogo, err := filepath.Glob("templates/docs/assets/images/logo.*.png")
	if err != nil {
		logger.Error(mainLogTag, fmt.Sprintf("Failed to find logo.*.png: %v", err))
		os.Exit(1)
	} else if len(themeImageLogo) != 1 {
		logger.Error(mainLogTag, fmt.Sprintf("Failed to find exactly one logo.*.png: found %d", len(themeImageLogo)))
		os.Exit(1)
	}

	assetsFuncs := template.FuncMap{
		"themeStylesheetApplication": func() string {
			return fmt.Sprintf("/docs/assets/stylesheets/%s", path.Base(themeStylesheetApplication[0]))
		},
		"themeStylesheetApplicationPalette": func() string {
			return fmt.Sprintf("/docs/assets/stylesheets/%s", path.Base(themeStylesheetApplicationPalette[0]))
		},
		"themeStylesheetExtra": func() string {
			return fmt.Sprintf("/docs/assets/stylesheets/%s", path.Base(themeStylesheetExtra[0]))
		},
		"themeJavascriptModernizr": func() string {
			return fmt.Sprintf("/docs/assets/javascripts/%s", path.Base(themeJavascriptModernizr[0]))
		},
		"themeJavascriptApplication": func() string {
			return fmt.Sprintf("/docs/assets/javascripts/%s", path.Base(themeJavascriptApplication[0]))
		},
		"themeImageFavicon": func() string {
			return fmt.Sprintf("/docs/assets/images/%s", path.Base(themeImageFavicon[0]))
		},
		"themeImageLogo": func() string {
			return fmt.Sprintf("/docs/assets/images/%s", path.Base(themeImageLogo[0]))
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

	utilityFuncs := template.FuncMap{
		"add": func(i1 int, i2 int) int {
			return i1 + i2
		},
		"sub": func(i1 int, i2 int) int {
			return i1 - i2
		},
	}

	m.Use(martrend.Renderer(
		martrend.Options{
			Layout:     "layout",
			Directory:  "./templates",
			Extensions: []string{".tmpl"},
			Funcs:      []template.FuncMap{assetsFuncs, analyticsConfigFuncs, htmlFuncs, utilityFuncs},
		},
	))
}
