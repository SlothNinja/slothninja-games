package restful

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"strings"

	"github.com/SlothNinja/inflect"
	"github.com/gin-gonic/gin"
)

const (
	nKey = "Notices"
	eKey = "Errors"
)

//
//func FlashesFrom(c *gin.Context) Flashes {
//	if v, ok := c.Get(flashKey); ok {
//		if f, ok := v.(Flashes); ok {
//			return f
//		}
//	}
//	f := make(Flashes)
//	c = WithFlashes(c, f)
//	return f
//}
//
//func WithFlashes(c *gin.Context, f Flashes) *gin.Context {
//	c.Set(flashKey, f)
//	return c
//}
//
//type Flashes map[string]interface{}
//
//func (f Flashes) add(t string, msg ...interface{}) {
//	if msgs, ok := f[t]; ok {
//		f[t] = append(msgs.([]interface{}), msg...)
//	} else {
//		f[t] = msg
//	}
//}
//
//func ClearFlash(ctx *Context) {
//	ctx.Data["Flashes"] = make(Flashes, 0)
//}

func init() {
	gob.Register(new(template.HTML))
}

type Notices []template.HTML

func NoticesFrom(c *gin.Context) (ns Notices) {
	ns, _ = c.Value(nKey).(Notices)
	return
}

func withNotices(c *gin.Context, ns Notices) *gin.Context {
	c.Set(nKey, ns)
	return c
}

func AddNoticef(c *gin.Context, format string, args ...interface{}) {
	withNotices(c, append(NoticesFrom(c), HTML(format, args...)))
}

type Errors []template.HTML

func ErrorsFrom(c *gin.Context) (es Errors) {
	es, _ = c.Value(eKey).(Errors)
	return
}

func withErrors(c *gin.Context, es Errors) *gin.Context {
	c.Set(eKey, es)
	return c
}

func AddErrorf(c *gin.Context, format string, args ...interface{}) {
	withErrors(c, append(ErrorsFrom(c), HTML(format, args...)))
}

func HTML(format string, args ...interface{}) template.HTML {
	return template.HTML(fmt.Sprintf(format, args...))
}

func ToSentence(strings []string) (sentence string) {
	switch length := len(strings); length {
	case 0:
	case 1:
		sentence = strings[0]
	case 2:
		sentence = strings[0] + " and " + strings[1]
	default:
		for i, s := range strings {
			switch i {
			case 0:
				sentence += s
			case length - 1:
				sentence += ", and " + s
			default:
				sentence += ", " + s
			}
		}
	}
	return sentence
}

func Camelize(ss ...string) string {
	return strings.Replace(Titlize(ss...), " ", "", -1)
}

func Titlize(ss ...string) string {
	return strings.Title(strings.TrimSpace(combine(ss...)))
}

func combine(ss ...string) string {
	result := ""
	for _, s := range ss {
		result += " " + strings.TrimSpace(s)
	}
	return result
}

func IDize(ss ...string) string {
	return strings.Replace(strings.ToLower(combine(ss...)), " ", "-", -1)
}

//func HTML(format string, args ...interface{}) template.HTML {
//	return template.HTML(fmt.Sprintf(format, args...))
//}

func IDString(s string) string {
	return strings.Replace(strings.ToLower(s), " ", "-", -1)
}

func JSONString(ss ...string) string {
	switch s := Camelize(ss...); len(s) {
	case 0:
		return ""
	case 1:
		return strings.ToLower(s[0:1])
	default:
		return strings.ToLower(s[0:1]) + s[1:]
	}
	//	var temp string
	//	for i, s := range strings.Split(s, " ") {
	//		if trimmed := strings.TrimSpace(s); i == 0 {
	//			temp += strings.ToLower(trimmed)
	//		} else {
	//			temp += trimmed
	//		}
	//	}
	//	return temp
}

func Pluralize(label string, value int) string {
	if value != 1 {
		return inflect.Pluralize(label)
	}
	return label
}

//func debugf(c *gin.Context, tmpl string, args ...interface{}) {
//	l(c).Debugf(caller()+tmpl, args...)
//}
//
//func errorf(c *gin.Context, tmpl string, args ...interface{}) {
//	l(c).Errorf(caller()+tmpl, args...)
//}
//
//func infof(c *gin.Context, tmpl string, args ...interface{}) {
//	l(c).Infof(caller()+tmpl, args...)
//}
//
//func warningf(c *gin.Context, tmpl string, args ...interface{}) {
//	l(c).Warningf(caller()+tmpl, args...)
//}
//
//func l(c *gin.Context) logging.Logger {
//	ctx := ContextFrom(c)
//	return logging.Get(ctx)
//}
//
//func caller() string {
//	pc, file, line, _ := runtime.Caller(2)
//	files := strings.Split(file, "/")
//	if lenFiles := len(files); lenFiles > 1 {
//		file = files[lenFiles-1]
//	}
//	fun := runtime.FuncForPC(pc).Name()
//	funs := strings.Split(fun, "/")
//	if lenFuns := len(funs); lenFuns > 2 {
//		fun = strings.Join(funs[len(funs)-2:], "/")
//	}
//	return fmt.Sprintf("%v#%v(L: %v) => ", file, fun, line)
//}
//
//func Printf(fmt string, args ...interface{}) {
//	log.Printf(caller()+fmt, args...)
//}
