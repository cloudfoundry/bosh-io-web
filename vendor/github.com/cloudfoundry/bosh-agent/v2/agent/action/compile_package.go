package action

import (
	"errors"

	boshcrypto "github.com/cloudfoundry/bosh-utils/crypto"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	boshmodels "github.com/cloudfoundry/bosh-agent/v2/agent/applier/models"
	boshcomp "github.com/cloudfoundry/bosh-agent/v2/agent/compiler"
)

type CompilePackageAction struct {
	compiler boshcomp.Compiler
}

func NewCompilePackage(compiler boshcomp.Compiler) (compilePackage CompilePackageAction) {
	compilePackage.compiler = compiler
	return
}

func (a CompilePackageAction) IsAsynchronous(_ ProtocolVersion) bool {
	return true
}

func (a CompilePackageAction) IsPersistent() bool {
	return false
}

func (a CompilePackageAction) IsLoggable() bool {
	return true
}

func (a CompilePackageAction) Run(blobID string, multiDigest boshcrypto.MultipleDigest, name, version string, deps boshcomp.Dependencies) (map[string]interface{}, error) {
	val := map[string]interface{}{}

	pkg := boshcomp.Package{
		BlobstoreID: blobID,
		Name:        name,
		Sha1:        multiDigest,
		Version:     version,
	}

	modelsDeps := []boshmodels.Package{}

	for _, dep := range deps {
		modelsDeps = append(modelsDeps, boshmodels.Package{
			Name:    dep.Name,
			Version: dep.Version,
			Source: boshmodels.Source{
				Sha1:        dep.Sha1,
				BlobstoreID: dep.BlobstoreID,
			},
		})
	}

	uploadedBlobID, uploadedDigest, err := a.compiler.Compile(pkg, modelsDeps)
	if err != nil {
		return val, bosherr.WrapErrorf(err, "Compiling package %s", pkg.Name)
	}

	result := map[string]string{
		"blobstore_id": uploadedBlobID,
		"sha1":         uploadedDigest.String(),
	}

	val = map[string]interface{}{
		"result": result,
	}
	return val, nil
}

func (a CompilePackageAction) Resume() (interface{}, error) {
	return nil, errors.New("not supported")
}

func (a CompilePackageAction) Cancel() error {
	return errors.New("not supported")
}
