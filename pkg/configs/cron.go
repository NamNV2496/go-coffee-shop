package configs

type ClearAllOrder struct {
	Hour   string `yaml:"hour"`
	Minute string `yaml:"minute"`
	// Schedule string `yaml:"schedule"`
}

type Cron struct {
	ClearAllOrder ClearAllOrder `yaml:"clear_all_order"`
}
