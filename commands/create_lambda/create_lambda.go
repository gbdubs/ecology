package create_lambda

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/util/output"
)

type CreateLambdaCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Project         string
	Lambda          string
	Platform        string
	Region          string
}

func (clc CreateLambdaCommand) Execute(o *output.Output) (err error) {
	o.Info("Creating Lambda").Indent()
	o.Dedent().Done()
	return nil
}
