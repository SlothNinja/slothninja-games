package restful

import (
	"fmt"
	"html/template"
	"math"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var Builtins = template.FuncMap{
	"today":       today,
	"parity":      parity,
	"odd":         odd,
	"even":        even,
	"zero":        zero,
	"inc":         inc,
	"equal":       equal,
	"greater":     greater,
	"less":        less,
	"ints":        ints,
	"round2":      round2,
	"LastUpdated": LastUpdated,
	"Time":        Time,
	"Date":        Date,
	"ToLower":     ToLower,
	"LoginURL":    LoginURL,
	"LogoutURL":   LogoutURL,
	"noescape":    noescape,
	"comment":     comment,
	"data":        data,
	"add":         add,
	"toSentence":  ToSentence,
}

func today() string {
	return time.Now().UTC().Format("January 2, 2006")
}

func Time(t time.Time) string {
	return t.Format("3:04PM MST")
}

func Date(t time.Time) string {
	return t.Format("Jan 2, 2006")
}

func parity(i int) string {
	if i%2 == 0 {
		return "even"
	}
	return "odd"
}

func odd(i int) bool {
	return i%2 != 0
}

func even(i int) bool {
	return !odd(i)
}

func zero(v interface{}) (result bool) {
	switch v.(type) {
	case int:
		result = v.(int) == 0
	case int64:
		result = v.(int64) == 0
	case time.Time:
		result = v.(time.Time).IsZero()
	}
	return
}

func inc(v int) int {
	return v + 1
}

func equal(v1, v2 interface{}) (result bool) {
	switch v1.(type) {
	case *int:
		switch v2.(type) {
		case int:
			result = *v1.(*int) == v2.(int)
		case *int:
			result = *v1.(*int) == *v2.(*int)
		}
	case int:
		return v1.(int) == v2.(int)
	case int64:
		return v1.(int64) == v2.(int64)
	case string:
		return v1.(string) == v2.(string)
	}
	return false
}

func less(v1, v2 interface{}) bool {
	switch v1.(type) {
	case int:
		return v1.(int) < v2.(int)
	case int64:
		return v1.(int64) < v2.(int64)
	}
	return false
}

func greater(v1, v2 interface{}) bool {
	switch v1.(type) {
	case int:
		return v1.(int) > v2.(int)
	case int64:
		return v1.(int64) > v2.(int64)
	}
	return false
}

func ints(start, stop int) []int {
	var s []int
	for i := start; i <= stop; i++ {
		s = append(s, i)
	}
	return s
}

const minute = 60
const hour = 60 * minute
const day = 24 * hour
const month = 31 * day
const year = 365 * day

func LastUpdated(v time.Time) string {
	duration := time.Since(v)
	seconds := duration.Seconds()
	switch {
	case seconds < minute:
		return fmt.Sprintf("%d sec", int(seconds))
	case seconds < hour:
		return fmt.Sprintf("%d min", int(seconds/minute))
	case seconds < day:
		return fmt.Sprintf("%d hour", int(seconds/hour))
	case seconds < month:
		return fmt.Sprintf("%d day", int(seconds/day))
	case seconds < year:
		return fmt.Sprintf("%d month", int(seconds/month))
	}
	return fmt.Sprintf("%d year", int(seconds/year))
}

func round2(f float64) string {
	i := math.Floor(f + 0.005)
	d := int(math.Floor(((f + 0.005) - i) * 100))
	if d < 10 {
		return fmt.Sprintf("%d.0%d", int(i), d)
	} else {
		return fmt.Sprintf("%d.%d", int(i), d)
	}
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func LogoutURL(c *gin.Context, redirect, label string) (template.HTML, error) {
	return "", nil
	// url, err := user.LogoutURL(ctx, redirect)
	// if err != nil {
	// 	return "", err
	// }
	// return template.HTML(fmt.Sprintf(`<a href=%q>%s</a>`, url, label)), nil
}

func LoginURL(c *gin.Context, redirect, label string) (tmpl template.HTML, err error) {
	return "", nil
	// url, err := user.LoginURL(ctx, redirect)
	// if err != nil {
	// 	return
	// }
	// tmpl = template.HTML(fmt.Sprintf(`<a href=%q>%s</a>`, url, label))
	// return
}

func noescape(s string) template.HTML {
	return template.HTML(s)
}

func comment(s string) template.HTML {
	return template.HTML("<!-- " + s + " -->")
}

// BindWith binds the passed struct pointer using the specified binding engine.
// Similar to gin#context.BindWith but does not write status code 400 to Header.
// This method provides more flexibility, e.g., redirection on error.
func BindWith(c *gin.Context, obj interface{}, b binding.Binding) error {
	if err := b.Bind(c.Request, obj); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		c.Abort()
		return err
	}
	return nil
}

func data(args ...interface{}) map[string]interface{} {
	d := make(map[string]interface{}, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		d[args[i].(string)] = args[i+1]
	}
	return d
}

func add(args ...int) (sum int) {
	for i := range args {
		sum += args[i]
	}
	return
}
