package mapping

import (
	"embed"
	"io"
	"io/fs"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed configs
var mappingFiles embed.FS

type MappingConfig struct {
	Version  string       `yaml:"version"`
	Name     string       `yaml:"name"`
	Mappings MappedFields `yaml:"mappings"`
	OneOf    OneOf        `yaml:"oneOf"`
}

type MappedFields []string

func (f MappedFields) ToSourceDestination() map[string]string {
	m := make(map[string]string)
	for _, s := range f {
		vals := strings.Split(s, ":")
		m[vals[0]] = vals[1]
	}
	return m
}

type OneOf struct {
	FieldName      string          `yaml:"field_name"`
	Discriminators []Discriminator `yaml:"discriminators"`
}

type Discriminator struct {
	Type     string       `yaml:"type"`
	Mappings MappedFields `yaml:"mappings"`
}

type Mappings map[string]MappingConfig

func Load() (Mappings, error) {
	mappings := make(Mappings)

	err := fs.WalkDir(mappingFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		f, err := mappingFiles.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		mc := MappingConfig{}
		err = yaml.Unmarshal(b, &mc)
		if err != nil {
			return err
		}

		mappings[mc.Name] = mc
		return nil
	})
	return mappings, err
}
