package config


var Configs *Config

func init() {
	addr := "127.0.0.1:5432"
	devConfig := DbConfig{addr, "resonate_dev_user", "password", "resonate_dev"}
	testingConfig := DbConfig{addr, "resonate_testing_user", "", "resonate_testing"}

	Configs = &Config{testingConfig, devConfig}
}

type DbConfig struct {
  Addr         string `json:"addr"`
  User     string `json:"user"`
  Password     string `json:"password"`
  Database string `json:"database"`
}

type Config struct {
  Testing DbConfig `json:"testing"`
	Dev     DbConfig `json:"dev"`
}
