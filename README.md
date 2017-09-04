## hlog4g
hlog4g is a Simplified implementation of Leveled logs for the Go language ([refer google glog](https://github.com/google/glog)).

## Example

create a new test file main.go, and input the following codes

```go
package main

import (
    "github.com/hooto/hlog4g/hlog"
)

func main() {

    // API:Print
    hlog.Print("info", "started")
    hlog.Print("error", "the error code/message: ", 400, "/", "bad request")

    // API::Printf
    hlog.Printf("error", "the error code/message: %d/%s", 400, "bad request")

    select {}
}
```

build the main.go file, run it and output the log into stderr console

```shell
go build main.go
./main -logtostderr=true
I 2016-04-30 22:24:26.463605 main.go:10] started
E 2016-04-30 22:24:26.463620 main.go:11] the error code/message: 400/bad request
E 2016-04-30 22:24:26.463628 main.go:14] the error code/message: 400/bad request
```

or run it and output the log into file
```shell
./main -log_dir="/var/log/"
```

the output file name will formated like this

```
/var/log/{program name}.{hostname}.{user name}.log.{level tag}.{date}-{time}.{pid}
```

## log levels
the default log levels are:

| Level | Tag | Description |
| --- | --- | --- |
| 0 | debug | Designates fine-grained informational events that are most useful to debug an application.|
| 1 | info | Designates informational messages that highlight the progress of the application at coarse-grained level|
| 2 | warn | Designates potentially harmful situations|
| 3 | error | Designates error events that might still allow the application to continue running|
| 4 | fatal | Designates very severe error events that will presumably lead the application to abort|


You can also define your custom levels:
```go
hlog.LevelConfig([]string{"warn", "error", "fatal"})
```

## Setting Flags

The flags influence hlog4g's output behavior by passing on the command line. For example, if you want to turn the flag --logtostderr on, you can start your application with the following command line:

``` shell
./your_application --logtostderr=true
```

The following flags are most commonly used:

| flag | type,default | Description |
| --- | --- | --- |
| log_dir | string, default="" | If specified, logfiles are written into this directory; if not specified or the directory does not valid, the hlog4g will not output messages to any logfiles|
| logtostderr | bool, default=false | If output log messages to stderr|
| minloglevel | int, default=1(which is INFO) | Log messages at or above this level. Again, the numbers of severity levels DEBUG, INFO, WARN, ERROR, and FATAL are 0, 1, 2, 3, and 4, respectively.|
| logtolevels | bool, default=false | If output messages to multi leveled logfiles from minloglevel to the max; or output messages to the minloglevel logfile.|


## Licensing
Licensed under the Apache License, Version 2.0


