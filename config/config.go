package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Db struct {
	Host      string `yaml:"host" env:"DB_HOST" env-default:"postgres"`
	Port      string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Pass      string `yaml:"pass" env:"DB_PASS" env-default:"test"`
	Name      string `yaml:"name" env:"DB_NAME" env-default:"postgres"`
	User      string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	MaxConn   int    `yaml:"maxConn" env:"DB_MAX_CONN" env-default:"50"`
	MaxIdle   int    `yaml:"maxIdle" env:"DB_MAX_IDLE" env-default:"10"`
	Schema    string `yaml:"schema" env:"DB_SCHEMA" env-default:"words"`
	BatchSize int    `yaml:"batchSize" env:"DB_BATCH" env-default:"1000"`
}

type Log struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"INFO"`
}

type App struct {
	Port           string `yaml:"port" env:"APP_PORT" env-default:"8080"`
	TestPort       string `yaml:"testPort" env:"APP_TEST_PORT" env-default:"8888"`
	Env            string `yaml:"env" env:"APP_ENV" env-default:"local"`
	Dicts          string `yaml:"dicts" env:"APP_DICTS_PATH" env-default:"./data/"`
	MigrationsPath string `yaml:"migrationsPath" env:"APP_MIGRATIONS_PATH" env-default:"./migrations"`
}

type Config struct {
	Db  `yaml:"db"`
	App `yaml:"app"`
	Log `yaml:"log"`
}

func New() (Config, error) {
	cfg := Config{}
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config/config.yaml" // путь по умолчанию
	}
	err := cleanenv.ReadConfig(path, &cfg)
	return cfg, err
}

func (c Db) GetDsn() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s search_path=%s sslmode=disable",
		c.User,
		c.Pass,
		c.Host,
		c.Port,
		c.Name,
		c.Schema,
	)
}

func (c Db) GetMigrateDsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		c.User,
		c.Pass,
		c.Host,
		c.Port,
		c.Name,
		c.Schema,
	)
}
