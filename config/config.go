package config

import "embed"

//go:embed *.yaml
var Default embed.FS
