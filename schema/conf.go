package schema

// Conf is the config to be passed
// to whatever database used
type Conf struct {
	Dialect  string `yaml:"dialect"`
	Host     string `yaml:"host"`
	Port     int16  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Sslmode  string `yaml:"sslmode"`
}
