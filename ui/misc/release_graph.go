package misc

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bhjobui "github.com/bosh-io/web/ui/job"
	bhrelui "github.com/bosh-io/web/ui/release"
)

type ReleaseGraph struct {
	pkgs []bhrelui.Package
	jobs []bhjobui.Job

	runner boshsys.CmdRunner

	logTag string
	logger boshlog.Logger
}

func NewReleaseGraph(
	pkgs []bhrelui.Package,
	jobs []bhjobui.Job,
	runner boshsys.CmdRunner,
	logger boshlog.Logger,
) ReleaseGraph {
	return ReleaseGraph{
		pkgs: pkgs,
		jobs: jobs,

		runner: runner,

		logTag: "ReleaseGraph",
		logger: logger,
	}
}

func (g ReleaseGraph) SVG() template.HTML {
	return g.render(nil)
}

func (g ReleaseGraph) FocusedSVG(focusedPkg bhrelui.Package) template.HTML {
	return g.render(&focusedPkg)
}

func (g ReleaseGraph) render(focusedPkg *bhrelui.Package) template.HTML {
	var in bytes.Buffer

	fmt.Fprintf(&in, "digraph packages {\n")   //nolint:errcheck
	fmt.Fprintf(&in, " size=\"11,1000\";\n")   //nolint:errcheck
	fmt.Fprintf(&in, " overlap=\"false\";\n")  //nolint:errcheck
	fmt.Fprintf(&in, " packMode=\"node\";\n")  //nolint:errcheck
	fmt.Fprintf(&in, " splines=\"spline\";\n") //nolint:errcheck
	fmt.Fprintf(&in, " sep=\"+20\";\n")        //nolint:errcheck

	for i, pkg := range g.pkgs {
		nodeI := i

		style := ", style=\"\", color=\"#ffffff\", fontname=\"Helvetica\""

		fmt.Fprintf( //nolint:errcheck
			&in,
			" n%d [label=\"%s\", URL=\"%s\", tooltip=\"%s\" %s];\n",
			nodeI, pkg.Name, pkg.URL(), pkg.Fingerprint, style,
		)

		for _, depPkg := range pkg.Dependencies {
			// Draw edge between pkg and its depPkg
			for allI, allPkg := range g.pkgs {
				if depPkg.Name == allPkg.Name {
					fmt.Fprintf(&in, " n%d -> n%d [style=\"dotted\"];\n", nodeI, allI) //nolint:errcheck
					break
				}
			}
		}
	}

	pkgsLen := len(g.pkgs)

	for i, job := range g.jobs {
		nodeI := pkgsLen + i

		style := ", style=\"filled\", color=\"#7E0180\", fontname=\"Helvetica\", fontcolor=\"white\""

		fmt.Fprintf(
			&in,
			" n%d [label=\"%s\", URL=\"%s\", nodesep=\"1\", shape=\"box\" %s];\n",
			nodeI, job.Name, job.URL(), style,
		)

		for _, depPkg := range job.Packages {
			// Draw edge between pkg and its depPkg
			for allI, allPkg := range g.pkgs {
				if depPkg.Name == allPkg.Name {
					fmt.Fprintf(&in, " n%d -> n%d [style=\"\"];\n", nodeI, allI) //nolint:errcheck
					break
				}
			}
		}
	}

	in.WriteString("}\n")

	stdout, _, _, err := g.runner.RunCommandWithInput(
		string(in.Bytes()), //nolint:staticcheck
		"sfdp",
		"-Tsvg",
	)
	if err != nil {
		g.logger.Error(g.logTag, "Failed to render: %s", err.Error())
		// Ignoring err from 'sfdp'
	}

	i := strings.Index(stdout, "<svg")
	if i < 0 {
		g.logger.Error(g.logTag, "Failed to find <svg")
		return template.HTML(errors.New("<svg not found").Error())
	}

	return template.HTML(stdout[i:])
}
