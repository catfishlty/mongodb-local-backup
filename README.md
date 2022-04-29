# MongoDB Local Backup


[![Release](https://img.shields.io/github/v/release/catfishlty/mongodb-local-backup)](https://github.com/catfishlty/mongodb-local-backup/releases/latest)
[![Download](https://img.shields.io/github/downloads/catfishlty/mongodb-local-backup/latest/total)](https://github.com/catfishlty/mongodb-local-backup/releases/latest)

[![Go Reference](https://pkg.go.dev/badge/github.com/catfishlty/mongodb-local-backup.svg)](https://pkg.go.dev/github.com/catfishlty/mongodb-local-backup)
[![Github Actions](https://github.com/catfishlty/mongodb-local-backup/actions/workflows/master.yml/badge.svg?branch=develop)](https://github.com/catfishlty/mongodb-local-backup/actions/workflows/master.yml)
[![codecov](https://codecov.io/gh/catfishlty/mongodb-local-backup/branch/develop/graph/badge.svg?token=79I70TJ9SF)](https://codecov.io/gh/catfishlty/mongodb-local-backup)
[![Go Report Card](https://goreportcard.com/badge/github.com/catfishlty/mongodb-local-backup)](https://goreportcard.com/report/github.com/catfishlty/mongodb-local-backup)

[![Mongo DB](https://img.shields.io/badge/MongoDB-4EA94B?style=for-the-badge&logo=mongodb&logoColor=white&text=MongoDB)](https://docs.mongodb.com/database-tools/mongoexport/)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)
![MacOS X](https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=apple&logoColor=white)
![Windows](https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)

A tool can back up your MongoDB data in Local File system.

English, [中文](https://github.com/catfishlty/mongodb-local-backup/blob/master/README_zh.md)

## Usage

### 1. define config file

#### 1.1 config file type

`MLB` support 3 types of config files, they're `json`, `toml` and `yaml`. And here're the example of the config
files: [json](#config_json), [toml](#config_toml), [yaml](#config_yaml)

#### 1.2 config file have several fields

| key | value | required | description |
| --- | :-- | :--: | :-- |
| mongo | C:\Program Files\MongoDB\Tools\100\bin\mongoexport.exe | Y | specific 'mongoexport' path |
| host | 127.0.0.1 | Y | MongoDB service host |
| port | 27017 | Y | MongoDB service port |
| username | test | N | MongoDB service username for authentication, use with password, unset or set to null means no authentication |
| password | test | N | MongoDBservice password for authentication, use with username, unset or set to null means no authentication |
| target | 'must bean array' | Y | define which db and collection to export |
| db | 'must be an object' | Y | define which db to export |
| collection | 'must be an array' | Y | define which collections in this db to export |
| prefix | mongodb-local-backup | N | define the prefix of the exported data file names |
| type | json/csv | Y | define the export data file format |
| output | E:\mongo_backup\ | Y | define the directory where store the export data files. |
| cron | */1 * * * * | N | define when run the export task, it will work only with command include '-d' option.|

#### 1.3 config file example

<a id="config_json" href="#">config.json</a>

```json
{
  "mongo": "C:\\Program Files\\MongoDB\\Tools\\100\\bin\\mongoexport.exe",
  "host": "127.0.0.1",
  "port": 27017,
  "username": null,
  "password": null,
  "target": [
    {
      "db": "test",
      "collection": [
        "test",
        "test1"
      ]
    }
  ],
  "prefix": "mongodb-local-backup-json",
  "type": "json",
  "output": "E:\\mongo_backup\\",
  "cron": "*/1 * * * *"
}
```

<a id="config_toml" href="#">config.toml</a>

```toml
# Define the path of mongoexport executable
mongo = "C:\\Program Files\\MongoDB\\Tools\\100\\bin\\mongoexport.exe"
# MongoDB host
host = "127.0.0.1"
# MongoDB port
port = 27017
# username of MongoDB connection
# username =
# password of MongoDB connection
# password =
# prefix is the prefix of exported data file name.
prefix = "mongodb-local-backup-toml"
#
type = "json"
output = "E:\\mongo_backup\\"
cron = "*/1 * * * *"
# define the MongoDB db and collections in target, mlb can do the backup for multiple dbs and collections.
[[target]]
db = "test"
collection = ["test", "test1"]
```

<a id="config_yaml" href="#">config.yaml</a>

```yaml
mongo: C:\\Program Files\\MongoDB\\Tools\\100\\bin\\mongoexport.exe
host: 127.0.0.1
port: 27017
target:
  -
    db: test
    collection:
      - test
      - test1
prefix: "mongodb-local-backup-yaml"
type: json
output: "E:\\mongo_backup\\"
cron: '*/1 * * * *'
```

### 2. Run command

#### 2.1 Run in Windows

```cmd
# one time
./mlb.exe -c config.json -f json
./mlb.exe -c config.toml -f toml
./mlb.exe -c config.yml -f yaml
./mlb.exe -c config.yaml -f yaml

# cron job
./mlb.exe -c config.json -f json -d
./mlb.exe -c config.toml -f toml -d
./mlb.exe -c config.yml -f yaml -d
./mlb.exe -c config.yaml -f yaml -d
```

#### 2.2 Run in Linux/Darwin

```bash
# one time
mlb -c config.json -f json
mlb -c config.toml -f toml
mlb -c config.yml -f yaml
mlb -c config.yaml -f yaml

# cron job
mlb -c config.json -f json -d
mlb -c config.toml -f toml -d
mlb -c config.yml -f yaml -d
mlb -c config.yaml -f yaml -d
```

## Feature plans

### 1. Service

The service is about the `Service` in `Windows`, `Linux` and `MacOS`. With the `Service`, `mlb` can run in background
and do the backup tasks automatically.

### 2. Notification

Aim to send message to IM software such as `Telegram`,`Bark`, `WXWorks`, `DingTalk` and so on.

## Contribution

If you have some great ideas, plz submit issue and let me know :D

Welcome to PR~

## Background

Why did I design this tool? A few months before, the data in MongoDB Standalone is broken by the hard disk failure, so I
lost my data and can't not backup. So if you want to keep your data safe, it's better to use MongoDB ReplicaSet or
backup more often.

## Thanks

1. [github.com/BurntSushi/toml](https://github.com/BurntSushi/toml) MIT Licence
2. [github.com/alexflint/go-arg](https://github.com/alexflint/go-arg) BSD-2-Clause License
3. [github.com/asaskevich/govalidator](https://github.com/asaskevich/govalidator) MIT Licence
4. [github.com/commander-cli/cmd](https://github.com/commander-cli/cmd) MIT Licence
5. [github.com/go-co-op/gocron](https://github.com/go-co-op/gocron) MIT Licence
6. [github.com/go-ozzo/ozzo-validation/v4](https://github.com/go-ozzo/ozzo-validation) MIT Licence
7. [github.com/gorhill/cronexpr](https://github.com/gorhill/cronexpr) GPL v3/APL v2
8. [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus) MIT Licence
9. [github.com/smartystreets/goconvey/convey](https://github.com/smartystreets/goconvey) [Licence](https://github.com/smartystreets/goconvey/blob/master/LICENSE.md)
10. [github.com/agiledragon/gomonkey](https://github.com/agiledragon/gomonkey) MIT Licence
