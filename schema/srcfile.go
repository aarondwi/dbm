package schema

// Srcfile is the format of the source file
// in which the user will fill with `up` and `down` quereis
type Srcfile struct {
	Up   string `yaml:"up"`
	Down string `yaml:"down"`
}
