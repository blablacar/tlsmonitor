package config

type Config struct {
	Metrics         Metrics  `yaml:"metrics"`
	NbDialRetry     int      `yaml:"nbDialRetry"`
	Hosts           []string `yaml:"hosts,omitempty"`
	CertsPaths      []string `yaml:"certsPaths,omitempty"`
	ChecksFrequency int      `yaml:"checksFrequency"      default:"1"`
}

type Metrics struct {
	Enabled      bool   `yaml:"enabled"      default:"false"`
	Path         string `yaml:"path"`
	Port         int    `yaml:"port"`
	WaitScraping bool   `yaml:"waitScraping" default:"false"`
}
