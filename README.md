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
		traceId := req.Header.Get(log.TraceIdHeaderKey)
		if traceId == "" {
			traceId = log.NewTraceId()
			req.Header.Set(log.TraceIdHeaderKey, traceId)
		}
		req = req.WithContext(log.NewContext(req.Context(), log.New(traceId)))
		// defer w.Header().Set(log.TraceIdHeaderKey, traceId)
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
time=2018-04-19 10:24:57.851, level=info, request_id=daaa14a3437811e884ddb4d5bdb21e16, location=main.httpHandler(test2/main.go:28), msg=1.info message
time=2018-04-19 10:24:57.851, level=info, request_id=daaa14a3437811e884ddb4d5bdb21e16, location=main.httpHandler(test2/main.go:29), msg=2.info message, key1=1, key2=2
time=2018-04-19 10:24:57.851, level=info, request_id=daaa14a3437811e884ddb4d5bdb21e16, location=main.httpHandler(test2/main.go:33), msg=3.info message, key5=5, key3=3, key4=4
time=2018-04-19 10:24:57.851, level=info, request_id=daaa14a3437811e884ddb4d5bdb21e16, location=main.fn1(test2/main.go:41), msg=4.info message, key3=3, key4=4, key5=5
time=2018-04-19 10:24:57.851, level=info, request_id=daaa14a3437811e884ddb4d5bdb21e16, location=main.fn1(test2/main.go:42), msg=5.info message, key3=3, key4=4, key5=5, key6=6
time=2018-04-19 10:24:57.851, level=info, request_id=daaa14a3437811e884ddb4d5bdb21e16, location=main.fn2(test2/main.go:51), msg=6.info message, key4=4, key5=5, key3=3, key7=7, key8=8
time=2018-04-19 10:24:57.851, level=info, request_id=daaa14a3437811e884ddb4d5bdb21e16, location=main.fn3(test2/main.go:58), msg=7.info message, key4=4, key5=5, key3=3, key7=7, key9=9
```
