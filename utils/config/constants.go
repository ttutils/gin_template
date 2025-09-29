package config

// GetDefaultAuthConfig 返回默认的AuthConfig
func GetDefaultAuthConfig() AuthConfig {
	return AuthConfig{
		ExcludedPaths: []string{
			"/api/user/login",
			"/nacos/v1/auth/login",
			"/api/ping",
			"/api/metrics",
			"/nacos/v1/cs/configs",
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
