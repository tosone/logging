# logging

[![Server Build CI](https://github.com/tosone/logging/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/tosone/logging/actions/workflows/ci.yaml) [![codecov](https://codecov.io/gh/tosone/logging/branch/main/graph/badge.svg?token=Y0l7RHluoS)](https://codecov.io/gh/tosone/logging)

Example Code:

``` go
package main

import "github.com/tosone/logging"

type Test struct {
	String string
	Int    int
}

func main() {
	var test = Test{String: "123123123", Int: 1}
	logging.Info("info level")
	logging.WithFields(logging.Fields{"field1": 1, "field2": "123"}).Info("info level")
	logging.Warn("warn info")
	logging.WithFields(logging.Fields{"field1": 1, "field2": "123"}).Warn("warn level")
	logging.Debug(test)
	logging.Debugf("%+v", test)
}
```

OutPut:

``` bash
INFO[10:39:49.100] info level                                file=main.go line=12
INFO[10:39:49.101] info level                                file=main.go line=13 field1=1 field2=123
WARN[10:39:49.101] warn info                                 file=main.go line=14
WARN[10:39:49.101] warn level                                file=main.go line=15 field1=1 field2=123
DEBU[10:39:49.101] {123123123 1}                             file=main.go line=16
DEBU[10:39:49.101] {String:123123123, Int:1}                 file=main.go line=17
```
