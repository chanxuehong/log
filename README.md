## log

```golang
package main

import (
	"context"
	"net/http"

	"github.com/chanxuehong/log"
)

func main() {
	var req = &http.Request{}
	var ctx = req.Context()

	l := log.FromRequest(req)
	l.Info("1.info message")
	l.Info("2.info message", "key1", 1, "key2", 2)

	l = l.WithField("key3", 3)
	l = l.WithFields("key4", 4, "key5", 5)
	l.Info("3.info message")

	ctx = log.NewContext(ctx, l)
	fn1(ctx)
}

func fn1(ctx context.Context) {
	l := log.FromContext(ctx)
	l.Info("4.info message")
	l.Info("5.info message", "key6", 6)

	l = l.WithField("key7", 7)
	ctx = log.NewContext(ctx, l)
	fn2(ctx)
}

func fn2(ctx context.Context) {
	l := log.FromContext(ctx)
	l.Info("6.info message", "key8", 8)

	fn3(ctx)
}

func fn3(ctx context.Context) {
	l := log.FromContext(ctx)
	l.Info("7.info message", "key9", 9)
}
```

```Text
time=2018-04-17 20:38:55.636, level=info, request_id=4ade7aeb423c11e8809cb4d5bdb21e16, file_line=test1/main.go:15, msg=1.info message
time=2018-04-17 20:38:55.636, level=info, request_id=4ade7aeb423c11e8809cb4d5bdb21e16, file_line=test1/main.go:16, msg=2.info message, key1=1, key2=2
time=2018-04-17 20:38:55.636, level=info, request_id=4ade7aeb423c11e8809cb4d5bdb21e16, file_line=test1/main.go:20, msg=3.info message, key3=3, key4=4, key5=5
time=2018-04-17 20:38:55.636, level=info, request_id=4ade7aeb423c11e8809cb4d5bdb21e16, file_line=test1/main.go:28, msg=4.info message, key3=3, key4=4, key5=5
time=2018-04-17 20:38:55.636, level=info, request_id=4ade7aeb423c11e8809cb4d5bdb21e16, file_line=test1/main.go:29, msg=5.info message, key5=5, key3=3, key4=4, key6=6
time=2018-04-17 20:38:55.636, level=info, request_id=4ade7aeb423c11e8809cb4d5bdb21e16, file_line=test1/main.go:38, msg=6.info message, key5=5, key7=7, key3=3, key4=4, key8=8
time=2018-04-17 20:38:55.636, level=info, request_id=4ade7aeb423c11e8809cb4d5bdb21e16, file_line=test1/main.go:45, msg=7.info message, key7=7, key9=9, key3=3, key4=4, key5=5
```
