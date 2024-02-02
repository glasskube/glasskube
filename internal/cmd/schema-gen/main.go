package main

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"go.uber.org/multierr"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/invopop/jsonschema"
	"github.com/spf13/cobra"
)

var (
	types = map[string]interface{}{
		"package-manifest": &v1alpha1.PackageManifest{},
	}

	cmd = &cobra.Command{
		Use: "schema-gen",
		RunE: func(cmd *cobra.Command, args []string) (retErr error) {
			for k, v := range types {
				schema := jsonschema.Reflect(v)
				outPath := path.Join(output, k)
				if err := os.MkdirAll(outPath, 0755); err != nil {
					return err
				}
				file, err := os.Create(path.Join(outPath, fileName))
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
				encoder.SetIndent("", indent)
				if err = encoder.Encode(schema); err != nil {
					return err
				}
			}
			return nil
		},
	}

	output   string
	indent   string
	fileName string
)

func init() {
	cmd.Flags().StringVarP(&output, "output", "o", "", "root directory for output files")
	cmd.Flags().StringVar(&indent, "indent", "  ", "indent string")
	cmd.Flags().StringVar(&fileName, "file-name", "schema.json", "name of schema files")
	_ = cmd.MarkFlagRequired("output")
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
