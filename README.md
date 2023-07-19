# Log Store

[![Tests Status](https://github.com/gouniverse/logstore/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/gouniverse/logstore/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/logstore)](https://goreportcard.com/report/github.com/gouniverse/logstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/logstore)](https://pkg.go.dev/github.com/gouniverse/logstore)

Logs messages to a database table.

## Installation
```
go get -u github.com/gouniverse/logstore
```

## Setup

```
logStore, err = logstore.NewStore(logstore.NewStoreOptions{
    DB: databaseInstance,
    LogTableName: "log",
    AutomigrateEnabled: true,
})

if err != nil {
    panic(error.Error())
}
```

## Usage

```
logStore.Info("Hello")
```


# Log Levels

1. LevelTrace - Something very low level
2. LevelDebug - Useful debugging information
3. LevelInfo - Something noteworthy happened!
4. LevelWarn - You should probably take a look at this
5. LevelError - Something failed but I'm not quitting
6. LevelFatal - Bye. Calls os.Exit(1) after logging
7. LevelPanic - I'm bailing. Calls panic() after logging

## Change Log
2023.07.19 - Updated instance creation to use options struct

2022.06.26 - Updated dependencies

2021.12.21 - Added LICENSE

2021.12.21 - Added test badge

2021.12.21 - Added support for DB dialects

2021.12.21 - Removed GORM dependency and moved to the standard library
