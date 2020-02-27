package mlog

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/slothninja-games/sn/color"
	"github.com/SlothNinja/slothninja-games/sn/restful"
	"github.com/SlothNinja/slothninja-games/sn/user"
	"github.com/gin-gonic/gin"
)

type MLog struct {
	Key        *datastore.Key `datastore:"__key__"`
	Messages   `datastore:"-"`
	SavedState []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
	// ID        int64  `gae:"$id"`
	// Kind      string `gae:"$kind"`
	// Messages  `gae:"SavedState"`
	// CreatedAt restful.CTime
	// UpdatedAt restful.UTime
}

// func (ms *Messages) ToProperty() (datastore.Property, error) {
// 	v, err := codec.Encode(ms)
// 	return datastore.MkPropertyNI(v), err
// }
//
// func (ms *Messages) FromProperty(p datastore.Property) error {
// 	return codec.Decode(ms, p.Value().([]byte))
// }

func New(id int64) *MLog {
	return &MLog{Key: datastore.IDKey(kind, id, nil)}
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
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		ml := From(c)
		if ml == nil {
			log.Errorf("Missing messagelog.")
			restful.AddErrorf(c, "Missing messagelog.")
			c.HTML(http.StatusOK, "shared/flashbox", gin.H{
				"Notices": restful.NoticesFrom(c),
				"Errors":  restful.ErrorsFrom(c),
			})
			return
		}
		m := ml.NewMessage(c)
		m.Text = c.PostForm("message")
		creatorsid := c.PostForm("creatorid")
		if creatorsid != "" {
			intID, err := strconv.ParseInt(creatorsid, 10, 64)
			if err != nil {
				restful.AddErrorf(c, "Invalid value received for creatorsid: %v", creatorsid)
				c.HTML(http.StatusOK, "shared/flashbox", gin.H{
					"Notices": restful.NoticesFrom(c),
					"Errors":  restful.ErrorsFrom(c),
				})
				return
			}
			m.CreatorID = intID
		}
		dsClient, err := datastore.NewClient(c, "")
		if err != nil {
			log.Errorf(err.Error())
			c.Abort()
		}

		_, err = dsClient.Put(c, ml.Key, ml)
		if err != nil {
			restful.AddErrorf(c, "mlog::AddMessage gaelic.Put Error: %s", err)
			log.Errorf("mlog::AddMessage gaelic.Put Error: %s", err)
			c.HTML(http.StatusOK, "shared/flashbox", gin.H{
				"Notices": restful.NoticesFrom(c),
				"Errors":  restful.ErrorsFrom(c),
			})
			return
		}
		log.Debugf("m: %#v", m)
		log.Debugf("ml: %#v", ml)
		cu, found := user.Current(c)
		var link template.HTML
		if found {
			link = cu.Link()
		}
		c.HTML(http.StatusOK, "shared/message", gin.H{
			"message": m,
			"ctx":     c,
			"map":     color.MapFrom(c),
			"link":    link,
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
	id, err := getID(c)
	if err != nil {
		restful.AddErrorf(c, err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		log.Errorf(err.Error())
		c.Abort()
		return
	}

	ml := New(id)
	err = dsClient.Get(c, ml.Key, ml)
	if err != nil {
		restful.AddErrorf(c, "Unable to get message log with ID: %v", id)
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}
	with(c, ml)
}

func From(c *gin.Context) (ml *MLog) {
	ml, _ = c.Value(mlKey).(*MLog)
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
