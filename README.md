# Log Store <a href="https://gitpod.io/#https://github.com/dracory/logstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

[![Tests Status](https://github.com/dracory/logstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/logstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/logstore)](https://goreportcard.com/report/github.com/dracory/logstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/logstore)](https://pkg.go.dev/github.com/dracory/logstore)

Logs messages to a database table.


## License

This project is licensed under the GNU General Public License version 3 (GPL-3.0). You can find a copy of the license at https://www.gnu.org/licenses/gpl-3.0.en.html

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

## Installation
```
go get -u github.com/dracory/logstore
```

## Setup

```golang
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

```golang
logStore.Info("Hello")

// with additional context
logStore.InfoWithContext("Hello", map[string]string{
    "name": "John Doe"
})
```

## Slog

As slog is the now official logger in golang, LogStore provides a SlogHandler.

```golang
Logger = *slog.New(logstore.NewSlogHandler(&LogStore))

logger.Info("Hello", "name", "John Doe")
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
2024.09.23 - Added a SlogHandler

2023.07.19 - Updated instance creation to use options struct

2022.06.26 - Updated dependencies

2021.12.21 - Added LICENSE

2021.12.21 - Added test badge

2021.12.21 - Added support for DB dialects

2021.12.21 - Removed GORM dependency and moved to the standard library
