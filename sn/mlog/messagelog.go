package mlog

import (
	"net/http"
	"strconv"

	"github.com/SlothNinja/slothninja-games/sn/codec"
	"github.com/SlothNinja/slothninja-games/sn/color"
	"github.com/SlothNinja/slothninja-games/sn/log"
	"github.com/SlothNinja/slothninja-games/sn/restful"
	"github.com/SlothNinja/slothninja-games/sn/user"
	"github.com/gin-gonic/gin"
	"go.chromium.org/gae/service/datastore"
	"golang.org/x/net/context"
)

type MLog struct {
	ID        int64  `gae:"$id"`
	Kind      string `gae:"$kind"`
	Messages  `gae:"SavedState"`
	CreatedAt restful.CTime
	UpdatedAt restful.UTime
}

func (ms *Messages) ToProperty() (datastore.Property, error) {
	v, err := codec.Encode(ms)
	return datastore.MkPropertyNI(v), err
}

func (ms *Messages) FromProperty(p datastore.Property) error {
	return codec.Decode(ms, p.Value().([]byte))
}

func New() *MLog {
	return &MLog{Kind: kind}
}

const (
	kind     = "MessageLog"
	mlKey    = "MessageLog"
	homePath = "/"
)

//func NewKey(ctx *restful.Context, stringID string, intID int64, pk *datastore.Key) *datastore.Key {
//	return datastore.NewKey(ctx, kind, stringID, intID, pk)
//}

//func (ml *MLog) init(ctx *restful.Context) error {
//	for _, m := range ml.Messages {
//		m.log = ml
//		m.ctx = ctx
//	}
//	return nil
//}

func AddMessage(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.ContextFrom(c)
		log.Debugf(ctx, "Entering")
		defer log.Debugf(ctx, "Exiting")

		ml := From(ctx)
		if ml == nil {
			log.Errorf(ctx, "Missing messagelog.")
			restful.AddErrorf(ctx, "Missing messagelog.")
			c.HTML(http.StatusOK, "shared/flashbox", gin.H{
				"Notices": restful.NoticesFrom(ctx),
				"Errors":  restful.ErrorsFrom(ctx),
			})
			return
		}
		m := ml.NewMessage(ctx)
		m.Text = c.PostForm("message")
		creatorsid := c.PostForm("creatorid")
		if creatorsid != "" {
			intID, err := strconv.ParseInt(creatorsid, 10, 64)
			if err != nil {
				restful.AddErrorf(ctx, "Invalid value received for creatorsid: %v", creatorsid)
				c.HTML(http.StatusOK, "shared/flashbox", gin.H{
					"Notices": restful.NoticesFrom(ctx),
					"Errors":  restful.ErrorsFrom(ctx),
				})
				return
			}
			m.CreatorID = intID
		}
		if err := datastore.Put(ctx, ml); err != nil {
			restful.AddErrorf(ctx, "mlog::AddMessage gaelic.Put Error: %s", err)
			log.Errorf(ctx, "mlog::AddMessage gaelic.Put Error: %s", err)
			c.HTML(http.StatusOK, "shared/flashbox", gin.H{
				"Notices": restful.NoticesFrom(ctx),
				"Errors":  restful.ErrorsFrom(ctx),
			})
			return
		}
		log.Debugf(ctx, "m: %#v", m)
		log.Debugf(ctx, "ml: %#v", ml)
		c.HTML(http.StatusOK, "shared/message", gin.H{
			"message": m,
			"ctx":     ctx,
			"map":     color.MapFrom(ctx),
			"link":    user.CurrentFrom(ctx).Link(),
		})
	}
}

//func BySID(ctx *restful.Context, sid string) (*MLog, error) {
//	ml := New(ctx)
//	id, err := strconv.ParseInt(sid, 10, 64)
//	if err != nil {
//		return nil, err
//	}
//	ml.Key = NewKey(ctx, "", id, nil)
//	if err := gaelic.Get(ctx, ml.Key, ml); err != nil {
//		return nil, err
//	}
//	return ml, nil
//}

func getID(c *gin.Context) (int64, error) {
	sid := c.Param("hid")
	return strconv.ParseInt(sid, 10, 64)
}

func Get(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	id, err := getID(c)
	if err != nil {
		restful.AddErrorf(ctx, err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		c.Abort()
		return
	}

	ml := New()
	ml.ID = id
	if err := datastore.Get(ctx, ml); err != nil {
		restful.AddErrorf(ctx, "Unable to get message log with ID: %v", id)
		c.Redirect(http.StatusSeeOther, homePath)
	} else {
		with(c, ml)
	}
}

func From(ctx context.Context) (ml *MLog) {
	ml, _ = ctx.Value(mlKey).(*MLog)
	return
}

func with(c *gin.Context, ml *MLog) *gin.Context {
	c.Set(mlKey, ml)
	return c
}

//func (ml *MLog) AfterCache() error {
//	return ml.init(ml.ctx)
//}

//func (ml *MLog) Save(c chan<- datastore.Property) error {
//	// Time stamp
//	t := time.Now()
//	if ml.CreatedAt.IsZero() {
//		ml.CreatedAt = t
//	}
//	ml.UpdatedAt = t
//
//	// Encode and save messages
//	if encoded, err := codec.Encode(ml.Messages); err != nil {
//		return err
//	} else {
//		ml.SavedState = encoded
//		return datastore.SaveStruct(ml, c)
//	}
//}
//
//func (ml *MLog) Load(c <-chan datastore.Property) error {
//	if err := datastore.LoadStruct(ml, c); err != nil {
//		return err
//	}
//	messages := make(Messages, 0)
//	if err := codec.Decode(&messages, ml.SavedState); err != nil {
//		return err
//	}
//	ml.Messages = messages
//	return ml.init(ml.ctx)
//}
