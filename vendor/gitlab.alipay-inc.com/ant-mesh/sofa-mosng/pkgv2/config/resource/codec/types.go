package codec

type metadata struct {
	Version  string      `yaml:"apiVersion"`
	Kind     string      `yaml:"kind"`
	Metadata interface{} `yaml:"metadata,omitempty"`
	Import   []string    `yaml:"import,omitempty"`
	Spec     interface{} `yaml:"spec,omitempty"`
}
