package engine

type Config struct {
	configPath string
}

var APPLICATION_CONFIG_PATH = CONFIG_PATH + "/paws"

func LoadConfig(uid string) *Config {
	return &Config{
		configPath: APPLICATION_CONFIG_PATH + "/" + uid,
	}
}
