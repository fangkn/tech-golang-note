package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w AccessLogWriter) Write(p []byte) (int, error) {
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

// AccessLog 调试情况下url请求的详细信息
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyWriter := &AccessLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyWriter

		reqBody, err := c.GetRawData()
		if err == nil {
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		}

		fmt.Printf("[REQUEST] %s path[%s] query[%v] clientip[%v] user-agent[%s] body[%v]",
			c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery,
			c.Request.RemoteAddr, c.Request.UserAgent(), string(reqBody))

		start := time.Now()
		c.Next()
		fmt.Printf("[RESPONSE] %s path[%s] query[%v] request[%v] response[%v] timecost[%v]",
			c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery,
			string(reqBody), bodyWriter.body.String(), time.Since(start).Seconds(),
		)
	}
}
