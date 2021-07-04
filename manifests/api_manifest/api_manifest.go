package api_manifest

import ()

type ApiConfigInfo struct {
}

type ApiDeployInfo struct {
	Platform string
	Region   string
}

type ApiManifest struct {
	Config ApiConfigInfo
	Deploy ApiDeployInfo
}
