package main

import (
	"go/parser"
	"go/token"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/docopt/docopt-go"
	gardener_api "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type declaration struct {
	Name  string
	Value string
}

type config struct {
	PackageName  string `docopt:"--pkg"`
	Declarations []declaration
	Out          string `docopt:"--out"`
}

const usage = `Shoot example helper
Usage:
	shoot-helper --pkg=<pkg> --out=<out>
Options:
	--pkg=<pkg>	the name of the package shoot declarations will be generated to
	--out=<out>	the name of the generated file (should end with _test postfix
`

const shootExampleTpl = `package {{ .PackageName }}
	
import (
	"github.com/kyma-project/infrastructure-manager/internal/gardener/shoot/testing"
)

var (
	{{- range .Declarations }}
	ShootExample_{{ .Name }} testing.ShootExample = "{{ .Value }}"
	{{- end }}
)
`

func exit1(err error) {
	slog.Error(err.Error())
	os.Exit(1)
}

func prepareDeclarations(path string) ([]declaration, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []declaration
	for _, entry := range entries {
		if entry.Type().Type().IsDir() {
			continue
		}

		declaration, err := toDeclaration(entry, path)
		if err != nil {
			return nil, err
		}

		result = append(result, declaration)
	}
	return result, nil
}

func toDeclaration(entry os.DirEntry, path string) (declaration, error) {
	var result declaration

	info, err := entry.Info()
	if err != nil {
		return result, err
	}

	filePath := filepath.Join(path, info.Name())
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}

	var shoot gardener_api.Shoot
	if err := yaml.NewYAMLOrJSONDecoder(file, 2048).Decode(&shoot); err != nil {
		return result, err
	}

	name := shoot.GetName()
	name = strings.ToUpper(name)
	name = strings.Join(strings.Split(name, "-"), "_")

	result.Name = name
	result.Value = filePath

	return result, nil
}

func generate(w io.Writer, data config) error {
	tpl, err := template.New("shoot-helper").Parse(shootExampleTpl)
	if err != nil {
		return err
	}
	return tpl.Execute(w, data)
}

func main() {
	// get grguments
	args, err := docopt.ParseDoc(usage)
	if err != nil {
		exit1(err)
	}
	// bind arguments to app configuration
	var cfg config
	if err := args.Bind(&cfg); err != nil {
		exit1(err)
	}
	// create out file
	file, err := os.Create(cfg.Out)
	if err != nil {
		exit1(err)
	}
	// prepare data
	cfg.Declarations, err = prepareDeclarations("./testdata")
	if err != nil {
		exit1(err)
	}
	// generate content
	err = generate(file, cfg)
	if err != nil {
		exit1(err)
	}
	// validate generated content
	fset := token.NewFileSet()
	_, err = parser.ParseFile(fset, cfg.Out, nil, parser.DeclarationErrors)
	if err != nil {
		exit1(err)
	}
}
