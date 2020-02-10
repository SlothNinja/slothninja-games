package log

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"go.chromium.org/luci/common/logging"
	"golang.org/x/net/context"
)

func Debugf(ctx context.Context, tmpl string, args ...interface{}) {
	l(ctx).Debugf(caller()+tmpl, args...)
}

func Errorf(ctx context.Context, tmpl string, args ...interface{}) {
	l(ctx).Errorf(caller()+tmpl, args...)
}

func Infof(ctx context.Context, tmpl string, args ...interface{}) {
	l(ctx).Infof(caller()+tmpl, args...)
}

func Warningf(ctx context.Context, tmpl string, args ...interface{}) {
	l(ctx).Warningf(caller()+tmpl, args...)
}

func l(ctx context.Context) logging.Logger {
	return logging.Get(ctx)
}

func caller() string {
	pc, file, line, _ := runtime.Caller(2)
	files := strings.Split(file, "/")
	if lenFiles := len(files); lenFiles > 1 {
		file = files[lenFiles-1]
	}
	fun := runtime.FuncForPC(pc).Name()
	funs := strings.Split(fun, "/")
	if lenFuns := len(funs); lenFuns > 2 {
		fun = strings.Join(funs[len(funs)-2:], "/")
	}
	return fmt.Sprintf("%v#%v(L: %v) => ", file, fun, line)
}

func Printf(fmt string, args ...interface{}) {
	log.Printf(caller()+fmt, args...)
}
