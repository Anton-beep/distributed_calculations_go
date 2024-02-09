package ExpressionLogger

import (
	"strings"
	"sync"
	"time"
)

type ExpLogger struct {
	data string
	mu   sync.Mutex
}

func timeNowString() string {
	return time.Now().Format("01-02-2006 15:04:05")
}

func New() *ExpLogger {
	return &ExpLogger{}
}

func (e *ExpLogger) Add(in ...string) {
	e.mu.Lock()
	e.data += "[" + timeNowString() + "] " + strings.Join(in, " ") + "\n"
	e.mu.Unlock()
}

func (e *ExpLogger) Get() string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.data
}

func (e *ExpLogger) Reset() {
	e.mu.Lock()
	e.data = ""
	e.mu.Unlock()
}
