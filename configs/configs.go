package configs

import (
	"encoding/json"
	"flag"
	"log"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const appName = "at"

var options = []option{
	// config section
	{"config", "string", "", "config file"},

	// server config section
	{"server.http.port", "int", 8080, "Server http port"},
	{"server.http.timeoutsec", "int", 20, "Server http timeout"},
	{"server.http.limiter.enabled", "bool", false, "Enables or disables limiter"},
	{"server.http.limiter.limit", "float64", 10000.0, "Limit tokens per second"},
	{"server.grpc.port", "int", 8081, "Server grpc port"},
	{"server.grpc.timeoutsec", "int", 20, "Server grpc timeout"},

	// database config section
	{"database.driver", "string", "mongodb", "database driver"},
	{"database.host", "string", "localhost", "database host"},
	{"database.port", "int", 27017, "database port"},
	{"database.user", "string", "root", "database user"},
	{"database.password", "string", "empty", "database password"},
	{"database.databasename", "string", appName, "database name"},
	{"database.secure", "string", "disable", "database SSL support"},
	{"database.args", "string", "", "database args"},

	//crawler
	{"crawler.concurrency", "int", 4, "Crawler concurency request"},

	// logger config section
	{"logger.level", "string", "debug", "LogLevel is global log level:  EMERG(0), ALERT(1), CRIT(2), ERR(3), WARNING(4), NOTICE(5), INFO(6), DEBUG(7)"},
	{"logger.timeformat", "string", "2006-01-02T15:04:05.999999999Z07:00", "LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00"},
}

type Config struct {
	Server   Server `yaml:"Servers"`
	Database Database
	Crawler  Crawler
	Logger   Logger
}

type Crawler struct {
	Concurrency int
}

type Server struct {
	GRPC Grpc
	HTTP Http
}

type Grpc struct {
	Port       int
	TimeoutSec int
}

type Http struct {
	Port       int
	TimeoutSec int
	Limiter    struct {
		Enabled bool
		Limit   float64
	}
}

type Database struct {
	Driver       string
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	Secure       string
	Args         string
}

type Logger struct {
	Level      string
	TimeFormat string
}

type option struct {
	name        string
	typing      string
	value       interface{}
	description string
}

func NewConfig() *Config {
	return &Config{}
}

// Read read parameters for config.
// Read from environment variables, flags or file.
func (c *Config) Read() error {
	viper.SetEnvPrefix(appName)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	for _, o := range options {
		switch o.typing {
		case "string":
			pflag.String(o.name, o.value.(string), o.description)
		case "int":
			pflag.Int(o.name, o.value.(int), o.description)
		case "bool":
			pflag.Bool(o.name, o.value.(bool), o.description)
		default:
			viper.SetDefault(o.name, o.value)
		}
	}

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	viper.BindPFlags(pflag.CommandLine)
	pflag.Parse()

	if fileName := viper.GetString("config"); fileName != "" {
		viper.SetConfigFile(fileName)
		viper.SetConfigType("toml")

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}

	if err := viper.Unmarshal(c); err != nil {
		return err
	}

	return nil
}

// Print print config structure
func (c *Config) Print() error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	log.Println(string(b))
	return nil
}
