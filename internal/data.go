package internal

const (
    Name           = "mongodb-local-backup"
    Version        = "0.3.0"
    BackupFileName = "mongo_backup_temp"
    TagCmd         = "cmd"
)

type BaseCmd struct {
    Config string `arg:"-c,--config,required" help:"config file"`
    Format string `arg:"-f,--format,required" help:"config file format(json, yaml, toml)"`
    Daemon bool   `arg:"-d,--daemon" help:"non stop running"`
}

type Args struct {
    StartCmd *BaseCmd `arg:"subcommand:start" help:"start application"`
}

type MongoTarget struct {
    Db         string   `json:"db,omitempty" yaml:"db,omitempty" cmd:"db" help:"mongodb db"`
    Collection []string `json:"collection,omitempty" yaml:"collection,omitempty" cmd:"collection" help:"mongodb collection"`
}

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
