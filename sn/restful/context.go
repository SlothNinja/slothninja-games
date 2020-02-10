package restful

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"bitbucket.org/SlothNinja/gin-render"
	"github.com/gin-gonic/gin"
	"go.chromium.org/gae/filter/dscache"
	"go.chromium.org/gae/impl/memory"
	"go.chromium.org/gae/impl/prod"
	"go.chromium.org/gae/service/datastore"
	"go.chromium.org/luci/common/logging"
	"golang.org/x/net/context"
)

const (
	contextKey = "Context"
	prefixKey  = "Prefix"
	requestKey = 0
	ginKey     = -1
	tmplKey    = "Templates"
)

func ContextFrom(c *gin.Context) (ctx context.Context) {
	if v, ok := c.Get(contextKey); ok {
		ctx, _ = v.(context.Context)
	}
	return
}

func WithContext(c *gin.Context, ctx context.Context) *gin.Context {
	c.Set(contextKey, ctx)
	return c
}

func GinFrom(ctx context.Context) (c *gin.Context) {
	c, _ = ctx.Value(ginKey).(*gin.Context)
	return
}

func RequestFrom(ctx context.Context) (r *http.Request) {
	r, _ = ctx.Value(requestKey).(*http.Request)
	return
}

func SetRequest(ctx context.Context, r *http.Request) *http.Request {
	c := GinFrom(ctx)
	c.Request = r
	return c.Request
}

func withTemplates(c *gin.Context, ts map[string]*template.Template) *gin.Context {
	c.Set(tmplKey, ts)
	return c
}

func TemplatesFrom(ctx context.Context) (ts map[string]*template.Template) {
	ts, _ = ctx.Value(tmplKey).(map[string]*template.Template)
	return
}

//func WithGin(ctx context.Context, c *gin.Context) context.Context {
//	return context.WithValue(ctx, ginKey, c)
//}

const (
	initialPoolSize = 32
)

func CTXHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c = WithContext(c, CreateContext(c))
	}
}

func CreateContext(c *gin.Context) context.Context {
	return logging.SetLevel(dscache.FilterRDS(prod.Use(c, c.Request)), logging.Debug)
}

func CreateTestContext(c *gin.Context) context.Context {
	return logging.SetLevel(dscache.FilterRDS(memory.Use(c)), logging.Debug)
}

func TemplateHandler(engine *gin.Engine) gin.HandlerFunc {
	r := render.New()
	r.TemplatesDir = "templates/"
	r.Exts = []string{".tmpl"}
	if gin.IsDebugging() {
		r.Debug = true
	}
	r.TemplateFuncMap = Builtins
	r = r.Init()
	return func(c *gin.Context) {
		engine.HTMLRender = r
		withTemplates(c, r.Templates)
	}
}

// Type for auto-updating time upon save
type UTime time.Time

// Implements datastore.PropertyConverter Interface
// To Auto Load/Save UTime property
func (ut *UTime) ToProperty() (p datastore.Property, err error) {
	p = datastore.Property{}
	err = p.SetValue(time.Now().UTC(), datastore.ShouldIndex)
	return
}

func (t *UTime) FromProperty(p datastore.Property) error {
	if it, err := p.Project(datastore.PTTime); err != nil {
		return err
	} else if tm, ok := it.(time.Time); !ok {
		return fmt.Errorf("Invalid UTime Value.")
	} else {
		*t = UTime(tm)
		return nil
	}
}

// GobEncode implements the gob.GobEncoder interface.
func (ut UTime) GobEncode() ([]byte, error) {
	t := time.Time(ut)
	return t.MarshalBinary()
}

// GobDecode implements the gob.GobDecoder interface.
func (ut *UTime) GobDecode(data []byte) error {
	t := new(time.Time)
	if err := t.UnmarshalBinary(data); err != nil {
		return err
	} else {
		*ut = UTime(*t)
		return nil
	}
}

// Type for that auto-sets time upon first save.
// No updates on subsequent saves.
type CTime time.Time

// Implements datastore.PropertyConverter Interface
// To Auto Load/Save CTime property
func (ct *CTime) ToProperty() (p datastore.Property, err error) {
	t := time.Time(*ct)
	if t.IsZero() {
		t = time.Now().UTC()
	}

	p = datastore.Property{}
	err = p.SetValue(t, datastore.ShouldIndex)
	return
}

func (t *CTime) FromProperty(p datastore.Property) error {
	if it, err := p.Project(datastore.PTTime); err != nil {
		return err
	} else if tm, ok := it.(time.Time); !ok {
		return fmt.Errorf("Invalid CTime Value.")
	} else {
		*t = CTime(tm)
		return nil
	}
}

// GobEncode implements the gob.GobEncoder interface.
func (t CTime) GobEncode() ([]byte, error) {
	return time.Time(t).MarshalBinary()
}

// GobDecode implements the gob.GobDecoder interface.
func (ct *CTime) GobDecode(data []byte) error {
	t := new(time.Time)
	if err := t.UnmarshalBinary(data); err != nil {
		return err
	} else {
		*ct = CTime(*t)
		return nil
	}
}
