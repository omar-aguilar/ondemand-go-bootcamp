package config

type Config struct {
	CSVSource string
	Port      int
}

var c = Config{
	CSVSource: "./data/db.csv",
	Port:      1234,
}

func GetConfig() Config {
	return c
}
