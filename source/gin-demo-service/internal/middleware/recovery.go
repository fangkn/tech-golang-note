package middleware

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

// Recovery 对url中的panic信息统一在这里上报，目前只打印错误日志
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				s := printCallersFrames()
				fmt.Println(strings.Join(s, "\n"))
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "panic"})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func printCallersFrames() []string {
	maxCallerDepth := 25
	minCallerDepth := 1
	callers := []string{}
	pcs := make([]uintptr, maxCallerDepth)
	depth := runtime.Callers(minCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		s := fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function)
		callers = append(callers, s)
		if !more {
			break
		}
	}
	return callers
}
