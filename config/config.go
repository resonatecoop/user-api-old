package config


var Config *config

func init() {
	addr := "127.0.0.1:5432"
	DevConfig := db_config{addr, "resonate_dev_user", "password", "resonate_dev"}
	TestingConfig := db_config{addr, "resonate_testing_user", "", "resonate_testing"}

	Config = &config{TestingConfig, DevConfig}
}

type db_config struct {
  Addr         string `json:"addr"`
  User     string `json:"user"`
  Password     string `json:"password"`
  Database string `json:"database"`
}

type config struct {
  Testing db_config `json:"testing"`
	Dev     db_config `json:"dev"`
}
