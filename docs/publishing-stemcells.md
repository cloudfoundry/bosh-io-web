All teams building stemcells should document their stemcell artifacts by recording checksums and where those artifacts can be downloaded from. We use a [metalink file](https://tools.ietf.org/html/rfc5854) to record this information and a [`meta4` binary](https://github.com/dpb587/metalink) to help generate and reference the metalink files.

Git repositories are used to record and audit the stemcell references, which are then indexed and rendered by [bosh.io](https://bosh.io/). These indexing repositories are named based on the team maintaining them, and deployment keys can be used within pipelines to push new commits.

 * [bosh-io/stemcells-core-index](https://github.com/bosh-io/stemcells-core-index)
 * [bosh-io/stemcells-cpi-index](https://github.com/bosh-io/stemcells-cpi-index)
 * [bosh-io/stemcells-windows-index](https://github.com/bosh-io/stemcells-windows-index)

The bosh.io site has a hard-coded list of these repositories and subpaths (e.g. `published`) for which stemcells are shown. It regularly pulls the stemcell index repositories and updates its pages.


## In Practice

Using the [bosh-linux-stemcell-builder](https://github.com/cloudfoundry/bosh-linux-stemcell-builder)'s stemcell [pipeline](https://github.com/cloudfoundry/bosh-linux-stemcell-builder/blob/master/ci/pipeline.yml) as an example...

**Development Stemcells** - stemcells which have been built (but are not fully tested, nor released to the public) are found in the [`dev`](https://github.com/bosh-io/stemcells-core-index/tree/master/dev) directory with the convention of `{os_name}-{os_version}/{version}/{iaas}-{hypervisor}-go_agent.meta4`.

**Published Stemcells** - stemcells which are fully tested and ready for public consumption can be found in the [`published`](https://github.com/bosh-io/stemcells-core-index/tree/master/published) directory with the convention of `{os_name}-{os_version}/{version}/stemcells.meta4`. This file references the stemcell tarballs for all IaaSes. This is the directory that bosh.io watches.


### Pipeline Walkthrough

We build stemcells like normal (e.g. [`rake stemcell:build`](https://github.com/cloudfoundry/bosh-linux-stemcell-builder/blob/687ac6998792c791f9d780d79527d2a1640987fa/ci/tasks/build.sh#L63-L71)). Once a stemcell is built, we generate a `meta4` file as part of the same [build task](https://github.com/cloudfoundry/bosh-linux-stemcell-builder/blob/687ac6998792c791f9d780d79527d2a1640987fa/ci/tasks/build.sh#L77-L106)...

    # paths
    metalink=stemcells-index-output/dev/$OS_NAME-$OS_VERSION/$CANDIDATE_BUILD_NUMBER/$IAAS-$HYPERVISOR-go_agent.meta4
    stemcell="${stemcell_name}.tgz"

    # create an empty metalink file for us to add the files
    mkdir -p "$( dirname "$metalink" )"
    meta4 create --metalink="$metalink"

    # import the stemcell tarball (this automatically calculates checksums)
    meta4 import-file --metalink="$metalink" --file="$stemcell" \
      --version="$CANDIDATE_BUILD_NUMBER" \
      "stemcell-output/$stemcell"

    # define where the stemcell can be downloaded
    meta4 file-set-url --metalink="$metalink" --file="$stemcell" \
      "https://s3.amazonaws.com/bosh-core-stemcells/${IAAS}/${stemcell}"

That `meta4` file in the `stemcells-index-output` output directory should then be committed...

    git add -A
    git config --global user.email "ci@localhost"
    git config --global user.name "CI Bot"
    git commit -m "dev: $OS_NAME-$OS_VERSION/$CANDIDATE_BUILD_NUMBER ($IAAS-$HYPERVISOR)"

Once committed, Concourse can take care of pushing the repository and uploading the the tarball to s3. Your pipeline can then run whatever system tests it needs to verify quality. For example, we run BATS and BRATS against a subset of the stemcells.

When a stemcell is ready to publish, part of the [build task](https://github.com/cloudfoundry/bosh-linux-stemcell-builder/blob/687ac6998792c791f9d780d79527d2a1640987fa/ci/tasks/publish.sh#L16-L20) will merge all the stemcell metalinks which are about to be published into a single file in the `published` directory (the directory bosh.io reads from).

    # paths
    metalink=stemcells-index-output/published/$OS_NAME-$OS_VERSION/$VERSION/stemcells.meta4

    # create an empty metalink file for us to merge everything
    mkdir -p "$( dirname "$metalink" )"
    meta4 create --metalink="$metalink"

    # merge in all our built stemcells
    find stemcells-index-output/dev/$OS_NAME-$OS_VERSION/$VERSION -name *.meta4 \
      | xargs -n1 -- meta4 import-metalink --metalink="$metalink"

Once merged, it should be committed and Concourse can take care of pushing the repository (along with whatever other stemcell promotion steps may be involved).
