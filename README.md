## log

```golang
package main

import (
	"context"
	"net/http"

	"github.com/chanxuehong/log"
)

func main() {
	req, _ := http.NewRequest(http.MethodGet, "http://a.com", nil)
	httpHandler(nil, req)
}

func httpHandler(w http.ResponseWriter, req *http.Request) {
	// In general, it is a middleware
	{
		requestId := req.Header.Get(log.RequestIdHeaderKey)
		if requestId == "" {
			requestId = log.NewRequestId()
			req.Header.Set(log.RequestIdHeaderKey, requestId)
		}
		req = req.WithContext(log.NewContext(req.Context(), log.New(requestId)))
		// defer w.Header().Set(log.RequestIdHeaderKey, requestId)
	}

	l := log.FromRequest(req)
	l.Info("1.info message")
	l.Info("2.info message", "key1", 1, "key2", 2)

	l = l.WithField("key3", 3)
	l = l.WithFields("key4", 4, "key5", 5)
	l.Info("3.info message")

	ctx := log.NewContext(req.Context(), l)
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
time=2018-04-18 13:27:30.784, level=info, request_id=30b53e7442c911e8927fb4d5bdb21e16, file_line=test1/main.go:28, msg=1.info message
time=2018-04-18 13:27:30.784, level=info, request_id=30b53e7442c911e8927fb4d5bdb21e16, file_line=test1/main.go:29, msg=2.info message, key2=2, key1=1
time=2018-04-18 13:27:30.784, level=info, request_id=30b53e7442c911e8927fb4d5bdb21e16, file_line=test1/main.go:33, msg=3.info message, key3=3, key4=4, key5=5
time=2018-04-18 13:27:30.785, level=info, request_id=30b53e7442c911e8927fb4d5bdb21e16, file_line=test1/main.go:41, msg=4.info message, key3=3, key4=4, key5=5
time=2018-04-18 13:27:30.785, level=info, request_id=30b53e7442c911e8927fb4d5bdb21e16, file_line=test1/main.go:42, msg=5.info message, key6=6, key3=3, key4=4, key5=5
time=2018-04-18 13:27:30.785, level=info, request_id=30b53e7442c911e8927fb4d5bdb21e16, file_line=test1/main.go:51, msg=6.info message, key7=7, key3=3, key4=4, key8=8, key5=5
time=2018-04-18 13:27:30.785, level=info, request_id=30b53e7442c911e8927fb4d5bdb21e16, file_line=test1/main.go:58, msg=7.info message, key7=7, key9=9, key3=3, key4=4, key5=5
```
