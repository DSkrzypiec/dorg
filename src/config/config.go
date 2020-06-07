package config

// Config stores all of dorg configuration.
type Config struct {
	DownloadsPath  string `json:"downloads_path"`
	TargetBasePath string `json:"target_path"`
}

// Default creates default dorg configuration.
func Default(downloadsPath string) Config {
	return Config{
		DownloadsPath:  downloadsPath,
		TargetBasePath: downloadsPath,
	}
}

// TODO: Implement loading configuration from .json file.
func TryParseConfig(configPath string) (Config, error) {
	return Config{}, nil
}
