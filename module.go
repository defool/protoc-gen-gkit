package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

var (
	logger *log.Logger
)

type mod struct {
	*pgs.ModuleBase
	pgsgo.Context
}

func newMod() pgs.Module {
	return &mod{ModuleBase: &pgs.ModuleBase{}}
}

func (m *mod) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.Context = pgsgo.InitContext(c.Parameters())
}

func (mod) Name() string {
	return "gkit"
}

func (m mod) Execute(targets map[string]pgs.File, gpkgs map[string]pgs.Package) []pgs.Artifact {
	initLogger(enableLogger)

	for _, pbFile := range targets {
		pkgName := *pbFile.Descriptor().Package
		if idx := strings.LastIndexByte(pkgName, '.'); idx >= 0 {
			pkgName = pkgName[idx+1:]
		}
		info := &Config{
			Package: pkgName,
		}
		for _, svc := range pbFile.Services() {
			svcInfo := ServiceInfo{
				ServiceName:      svc.Name().String(),
				ServiceNameLower: firstLowger(svc.Name().String()),
			}
			info.Services = append(info.Services, svcInfo)
		}
		if len(info.Services) == 0 {
			continue
		}
		buf := bytes.NewBuffer(nil)
		tp := template.Must(template.New("").Parse(outTemplate))
		err := tp.Execute(buf, info)
		m.CheckErr(err)
		filename := m.Context.OutputPath(pbFile).SetExt(".gkit.go").String()
		m.AddGeneratorFile(filename, buf.String())
	}
	return m.Artifacts()
}

func firstLowger(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func initLogger(enable bool) {
	logger = log.Default()
	if enable {
		logger.SetOutput(os.Stderr)
	} else {
		logger.SetOutput(io.Discard)
	}
}
