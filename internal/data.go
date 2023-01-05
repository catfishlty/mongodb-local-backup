package internal

const (
	// Name define app name
	Name = "mongodb-local-backup"
	// Version define app version
	Version = "0.5.2"
	// BackupFileName define backup temp file name
	BackupFileName = "mongo_backup_temp"
	// TagCmd define variable for parse data in golang tags
	TagCmd = "cmd"
)

// BaseCmd base command struct
type BaseCmd struct {
	Config   string `arg:"-c,--config" help:"config file"`
	Format   string `arg:"-f,--format" help:"config file format(json, yaml, toml)"`
	Daemon   bool   `arg:"-d,--daemon" help:"non stop running"`
	LogLevel string `arg:"-l,--log" default:"info" help:"log level(debug, info, warn, error, fatal, panic)"`
}

// Args all commands struct
type Args struct {
	StartCmd *BaseCmd `arg:"subcommand:start" help:"start application"`
}

// MongoTarget command options for target MongoDB db name and collection name
type MongoTarget struct {
	Db         string   `json:"db,omitempty" yaml:"db,omitempty" cmd:"db" help:"mongodb db"`
	Collection []string `json:"collection,omitempty" yaml:"collection,omitempty" cmd:"collection" help:"mongodb collection"`
}

// Config config file struct
type Config struct {
	Mongo    string        `json:"mongo,omitempty" yaml:"mongo,omitempty" help:"path of mongoexport program"`
	Host     string        `json:"host,omitempty" yaml:"host,omitempty" cmd:"host" help:"mongodb host"`
	Port     int           `json:"port,omitempty" yaml:"port,omitempty" cmd:"port" help:"mongodb port"`
	Username string        `json:"username,omitempty" yaml:"username,omitempty" cmd:"username" help:"mongodb username"`
	Password string        `json:"password,omitempty" yaml:"password" cmd:"password" help:"mongodb password"`
	Target   []MongoTarget `json:"target,omitempty" yaml:"target,omitempty" help:"mongodb target"`
	Type     string        `json:"type,omitempty" yaml:"type,omitempty" cmd:"type" help:"backup file format"`
	Output   string        `json:"output,omitempty" yaml:"output,omitempty" help:"specific a path to store backup files"`
	Cron     string        `json:"cron,omitempty" yaml:"cron,omitempty" help:"CRON expression"`
	Prefix   string        `json:"prefix,omitempty" yaml:"prefix,omitempty" help:"backup file prefix"`
	Log      string        `json:"log,omitempty" yaml:"log,omitempty" help:"log file path"`
}
