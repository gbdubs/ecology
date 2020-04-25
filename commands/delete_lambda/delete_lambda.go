package delete_lambda

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/util/output"
)

type DeleteLambdaCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Project         string
	Lambda          string
}

func (dlc DeleteLambdaCommand) Execute(o *output.Output) (err error) {
	o.Info("Deleting Lambda").Indent()
	o.Dedent().Done()
	return nil
}
