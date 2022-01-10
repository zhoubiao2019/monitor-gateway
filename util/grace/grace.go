package grace

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"go.uber.org/zap"

	"github.com/zhoubiao2019/monitor-gateway/util/log"
)

type panicError struct {
	panicObject interface{}
}

//
func (pe panicError) Error() string {
	return fmt.Sprintf("panic(%s): %v", IdentifyPanic(), pe.panicObject)
}

func (pe panicError) PanicObject() interface{} {
	return pe.panicObject
}

type Grace struct {
	cleanups []func()
}

// Register 注册清理函数，确保在Run执行传入的函数后（即使收到中断信号）会执行之前注册过的清理函数
func (g *Grace) Register(f func()) {
	g.cleanups = append(g.cleanups, f)
}

// Run 执行传入的函数，并运行之前注册过的清理函数
func (g *Grace) Run(f func() error) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	quit := make(chan error)
	go func() {
		var err error
		defer func() {
			if e := recover(); e != nil {
				quit <- panicError{e}
			} else {
				quit <- err
			}
		}()

		err = f()
	}()

	select {
	case sig := <-signals:
		log.Info("receive signal", zap.String("signal", sig.String()))
	case err := <-quit:
		log.Error("exit with error", zap.Error(err))
	}

	for _, cleanup := range g.cleanups {
		cleanup()
	}
}

func IdentifyPanic() string {
	var name, file string
	var line int
	var pc [16]uintptr

	n := runtime.Callers(3, pc[:])
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line = fn.FileLine(pc)
		name = fn.Name()
		if !strings.HasPrefix(name, "runtime.") {
			break
		}
	}

	if name != "" {
		return fmt.Sprintf("%v:%v", name, line)
	} else if file != "" {
		return fmt.Sprintf("%v:%v", file, line)
	}

	return fmt.Sprintf("pc:%x", pc)
}
