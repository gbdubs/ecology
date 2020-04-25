package push_lambda

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/util/output"
)

type PushLambdaCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Project         string
	Lambda          string
}

func (plc PushLambdaCommand) Execute(o *output.Output) (err error) {
	o.Info("Pushing Lambda").Indent()
	o.Dedent().Done()
	return nil
}
