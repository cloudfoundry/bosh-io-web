### Prod configuration

bosh-hub is meant to be the destination for all BOSH related things: releases, stemcells, documentation, etc. 

#### Release tracking

Releases are tracked by adding a new [WatcherRec](release/watchersrepo/watcher.go). [Periodic release watcher](release/watcher/periodic_watcher.go) runs in a goroutine and downloads git repository for every watched release periodically. It compares if there are any _new_ final releases and adds them as new [ImportRec](release/importsrepo/import.go)s. 

[Periodic release importer](release/importer/periodic_importer.go) runs in a goroutine and picks up one release import at a time. It runs `bosh create release X --with-tarball` to build a release tarball for a specific release version. It then uploads created tarball and creates a [ReleaseTarballRec](release/releasetarsrepo/release_tarball.go). If importing of a release fails new [ImportErrRec](release/importerrsrepo/import_err.go) is created and release version is not tried to be imorted again until associated ImportErr is deleted.

In addition to creating ReleaseTarballRecs for each imported release version, release, jobs and packages details are saved.

#### Stemcell tracking

[Periodic S3 bucket importer](stemcell/importer/periodic_s3_bucket_importer.go) runs periodically and downloads all S3 items' metadata from all configured buckets. It then filters out non-stemcell items and saves links to found stemcell items as [s3StemcellRec](stemcell/stemsrepo/s3_stemcells_repository.go)s

[StemcellsController](controllers/stemcells_controller.go) sorts and shows saved results.

#### Docs

[DocsController](controllers/docs_controller.go) wraps pre-rendered HTML documents under `docs/` path. Doc HTML pages are rendered before deploying an app. Running `./build-docs.sh` uses Bookbinder project with configuration stored in `docs-bosh-io/` folder.
