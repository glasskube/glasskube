package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"

	"go.uber.org/multierr"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/invopop/jsonschema"
)

var (
	types = map[string]interface{}{
		"package-manifest": &v1alpha1.PackageManifest{},
		"index":            &repo.PackageRepoIndex{},
		"versions":         &repo.PackageIndex{},
	}

	outBase    = "."
	idBase     = "https://glasskube.dev/"
	schemaPath = "schemas/v1"
	outDir     = path.Join(outBase, schemaPath)

	reflector = jsonschema.Reflector{
		ExpandedStruct: true,
	}
)

func run() (retErr error) {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	for k, v := range types {
		fileName := fmt.Sprintf("%v.json", k)
		schema := reflector.Reflect(v)

		if id, err := url.JoinPath(idBase, schemaPath, fileName); err != nil {
			return err
		} else {
			schema.ID = jsonschema.ID(id)
		}

		file, err := os.Create(path.Join(outDir, fileName))
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			closeErr := file.Close()
			if closeErr != nil {
				retErr = multierr.Append(err, closeErr)
			}
		}(file)

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		fmt.Fprintf(os.Stderr, "Writing to %v\n", file.Name())
		if err = encoder.Encode(schema); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
