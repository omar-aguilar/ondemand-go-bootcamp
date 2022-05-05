package config

import "fmt"

type Config struct {
	DBFile      string
	StoreFolder string
	Port        int
}

func (c Config) GetDBPath() string {
	return fmt.Sprintf("%s/%s", c.StoreFolder, c.DBFile)
}

var c = Config{
	DBFile:      "db.csv",
	StoreFolder: "./data",
	Port:        1234,
}

func GetConfig() Config {
	return c
}
