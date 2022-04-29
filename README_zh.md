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

一款易用的MongoDB本地备份工具。

[English](https://github.com/catfishlty/mongodb-local-backup#readme), 中文

## 使用指南

### 1. 配置文件

#### 1.1 配置文件格式

`MLB` 支持 `json`, `toml` and `yaml` 3 种格式的配置文件
示例文件: [json](#config_json), [toml](#config_toml), [yaml](#config_yaml)

#### 1.2 配置文件字段解析

| key | value | required | description |
| --- | :-- | :--: | :-- |
| mongo | C:\Program Files\MongoDB\Tools\100\bin\mongoexport.exe | Y | 'mongoexport'路径 |
| host | 127.0.0.1 | Y | MongoDB数据库 Host |
| port | 27017 | Y | MongoDB数据库 端口 |
| username | test | N | MongoDB数据库 用户名，如没有密码则保持为空 |
| password | test | N | MongoDB数据库 密码，如没有密码则保持为空 |
| target | 'must bean array' | Y | 需要导出的目标数据库和集合信息 |
| db | 'must be an object' | Y | 需要导出的目标数据库 |
| collection | 'must be an array' | Y | 需要导出的目标集合 |
| prefix | mongodb-local-backup | N | 导出保存文件的前缀 |
| type | json/csv | Y | 导出文件格式 |
| output | E:\mongo_backup\ | Y | 导出文件位置 |
| cron | */1 * * * * | N | 定时任务配置, 仅在启用 '-d' 后运行 |

#### 1.3 配置文件示例

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

### 2. 运行命令

#### 2.1 Windows平台

```cmd
# 仅运行一次
./mlb.exe -c config.json -f json
./mlb.exe -c config.toml -f toml
./mlb.exe -c config.yml -f yaml
./mlb.exe -c config.yaml -f yaml

# 定时任务
./mlb.exe -c config.json -f json -d
./mlb.exe -c config.toml -f toml -d
./mlb.exe -c config.yml -f yaml -d
./mlb.exe -c config.yaml -f yaml -d
```

#### 2.2 Linux/Darwin 平台

```bash
# 仅运行一次
mlb -c config.json -f json
mlb -c config.toml -f toml
mlb -c config.yml -f yaml
mlb -c config.yaml -f yaml

# 定时任务
mlb -c config.json -f json -d
mlb -c config.toml -f toml -d
mlb -c config.yml -f yaml -d
mlb -c config.yaml -f yaml -d
```

## 开发计划

### 1. Service

目标是 `mlb` 能够作为服务运行在 `Windows`, `Linux` and `MacOS`. 作为 `Service`, `mlb` 能够后台运行，并且根据配置自动备份数据。

### 2. 消息推送

能够使用系列的应用推送消息 `Telegram`,`Bark`, `WXWorks`, `DingTalk` .

## 代码贡献

如果你有更好的想法，欢迎提出issue

更欢迎大家能够多提PR~
Welcome to PR~

## 背景

之所以想做这样工具主要是满足自己的日常需求。几个月前，自己电脑中的MongoDB数据库突然就不能启动，也想官方咨询过，数据恢复无望。
考虑到当时上面的数据也没什么重要资料，后面就开始研究这样一个工具来完成自动备份的工作。

## 致谢

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
