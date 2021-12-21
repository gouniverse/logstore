# Log Store

Logs messages to a database table.

## Installation
```
go get -u github.com/gouniverse/logstore
```

## Setup

```
logStore = logstore.NewStore(logstore.WithGormDb(databaseInstance), logstore.WithTableName("log_meta"), logstore.WithAutoMigrate(true))
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