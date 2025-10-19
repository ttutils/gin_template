package config

// GetDefaultAuthConfig 返回默认的AuthConfig
func GetDefaultAuthConfig() AuthConfig {
	return AuthConfig{
		ExcludedPaths: []string{
			"/api/user/login",
			"/api/ping",
			"/api/metrics",
			"/api/server_info",
			"/api/is_demo",
		},
	}
}

// GetDefaultConfkeeperConfig 返回默认的ConfkeeperConfig
func GetDefaultConfkeeperConfig() ConfkeeperConfig {
	return ConfkeeperConfig{
		ConfigType: []string{
			"text",
			"json",
			"xml",
			"yaml",
			"html",
			"properties",
			"toml",
			"ini",
		},
		ActionType: []string{
			"r",
			"rw",
		},
	}
}
