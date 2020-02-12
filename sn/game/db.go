package game

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/SlothNinja/slothninja-games/sn/log"
	"github.com/SlothNinja/slothninja-games/sn/restful"
	gType "github.com/SlothNinja/slothninja-games/sn/type"
	"github.com/gin-gonic/gin"
	"go.chromium.org/gae/service/datastore"
	"golang.org/x/net/context"
)

const (
	statusKey = "Status"
	countKey  = "Count"
	NoCount   = -1
)

func getAllQuery(ctx context.Context) *datastore.Query {
	return datastore.NewQuery("Game").Ancestor(GamesRoot(ctx))
}

func getFiltered(ctx context.Context, status, sid, start, length string, t gType.Type) (gs Gamers, cnt int64, err error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	c := restful.GinFrom(ctx)
	q := getAllQuery(ctx).KeysOnly(true)
	if status != "" {
		st := ToStatus[strings.ToLower(status)]
		q = q.Eq("Status", st)
		WithStatus(c, st)
	}

	if sid != "" {
		if id, err := strconv.Atoi(sid); err == nil {
			q = q.Eq("UserIDS", id)
		}
	}

	if t != gType.All {
		q = q.Eq("Type", t).Order("-UpdatedAt")
	} else {
		q = q.Order("-UpdatedAt")
	}

	if cnt, err = datastore.Count(ctx, q); err != nil {
		log.Errorf(ctx, "sn/game#GetFiltered q.Count Error: %s", err)
		return
	}

	if start != "" {
		if st, err := strconv.ParseInt(start, 10, 32); err == nil {
			q = q.Offset(int32(st))
		}
	}

	if length != "" {
		if l, err := strconv.ParseInt(length, 10, 32); err == nil {
			q = q.Limit(int32(l))
		}
	}

	var ks []*datastore.Key
	if err = datastore.GetAll(ctx, q, &ks); err != nil {
		log.Errorf(ctx, "getFiltered GetAll Error: %s", err)
		return
	}

	l := len(ks)
	gs = make([]Gamer, l)
	hs := make([]*Header, l)
	for i := range gs {
		var ok bool
		if t == gType.All {
			k := strings.ToLower(ks[i].Parent().Kind())
			if t, ok = gType.ToType[k]; !ok {
				err = fmt.Errorf("Unknown Game Type For: %s", k)
				log.Errorf(ctx, err.Error())
				return
			}
		}
		gs[i] = factories[t](ctx)
		hs[i] = gs[i].GetHeader()
		if ok := datastore.PopulateKey(hs[i], ks[i]); !ok {
			err = fmt.Errorf("Unable to populate header with key.")
			log.Errorf(ctx, err.Error())
			return
		}
	}

	if err = datastore.Get(ctx, hs); err != nil {
		log.Errorf(ctx, "SN/Game#GetFiltered datastore.Get Error: %s", err)
		return
	}

	for i := range hs {
		if err = hs[i].AfterLoad(gs[i]); err != nil {
			log.Errorf(ctx, "SN/Game#GetFiltered h.AfterLoad Error: %s", err)
			return
		}
	}

	return
}

//func SetStatus(c *gin.Context) {
//	ctx := restful.ContextFrom(c)
//	log.Debugf(ctx, "Entering")
//	defer log.Debugf(ctx, "Exiting")
//
//	stat := c.Param("status")
//	status := ToStatus[stat]
//	WithStatus(c, status)
//}

func WithStatus(c *gin.Context, s Status) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	c.Set(statusKey, s)
}

func StatusFrom(ctx context.Context) (s Status) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	c := restful.GinFrom(ctx)
	if s = ToStatus[strings.ToLower(c.Param("status"))]; s != NoStatus {
		WithStatus(c, s)
	} else {
		s, _ = ctx.Value(statusKey).(Status)
	}
	return
}

func withCount(c *gin.Context, cnt int64) *gin.Context {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	c.Set(countKey, cnt)
	return c
}

func countFrom(ctx context.Context) (cnt int64) {
	cnt, _ = ctx.Value(countKey).(int64)
	return
}

//func SetType(ctx context.Context) {
//	ctx = gType.WithType(ctx, gType.All)
//	prefix := restful.PrefixFrom(ctx)
//	if t, ok := gType.ToType[prefix]; ok {
//		ctx = gType.WithType(ctx, t)
//		return
//	}
//}

//func GetAll(c *gin.Context) {
//	q := getAllQuery(c).Order("-UpdatedAt").KeysOnly(true)
//
//	if stat := c.Param("status"); stat != "" {
//		status := ToStatus[stat]
//		q = q.Eq("Status", status)
//		c = WithStatus(c, status)
//	}
//
//	if sid := c.Param("uid"); sid != "" {
//		if id, err := strconv.Atoi(sid); err == nil {
//			q = q.Eq("UserIDS", id)
//		}
//	}
//
//	var limit int32 = 100
//	prefix := restful.PrefixFrom(c)
//	if t, ok := gType.ToType[prefix]; ok {
//		c = gType.WithType(c, t)
//		q = q.Eq("Type", t).Order("-UpdatedAt").Limit(limit)
//	} else {
//		q = q.Order("-UpdatedAt").Limit(limit)
//	}
//
//	var ks []*datastore.Key
//	ctx := restful.ContextFrom(c)
//	if err := datastore.GetAll(ctx, q, &ks); err != nil {
//		log.Errorf(c, "SN/Game#All Error: %s", err)
//		c.Redirect(http.StatusSeeOther, homePath)
//		return
//	}
//
//	length := len(ks)
//	gs := make([]Gamer, length)
//	for i := range gs {
//		var (
//			t  gType.Type
//			ok bool
//		)
//		if t, ok = gType.ToType[prefix]; !ok {
//			k := strings.ToLower(ks[i].Parent().Kind())
//			if t, ok = gType.ToType[k]; !ok {
//				log.Errorf(c, "Unknown Game Type For: %s", k)
//				c.Redirect(http.StatusSeeOther, homePath)
//				return
//			}
//		}
//		gs[i] = factories[t](c)
//		if ok := datastore.PopulateKey(gs[i], ks[i]); !ok {
//			log.Errorf(c, "Unable to populate gamer with key")
//			c.Redirect(http.StatusSeeOther, homePath)
//			return
//		}
//	}
//
//	if err := datastore.Get(ctx, gs); err != nil {
//		log.Errorf(c, "SN/Game#All gaelic.GetMulti Error: %s", err)
//		c.Redirect(http.StatusSeeOther, homePath)
//		return
//	}
//
//	WithGamers(c, gs)
//}

func GetFiltered(t gType.Type) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.ContextFrom(c)
		log.Debugf(ctx, "Entering")
		defer log.Debugf(ctx, "Exiting")

		gs, cnt, err := getFiltered(ctx, c.Param("status"), c.Param("uid"), c.PostForm("start"), c.PostForm("length"), t)

		if err != nil {
			log.Errorf(ctx, err.Error())
			c.Redirect(http.StatusSeeOther, homePath)
			c.Abort()
		}
		withGamers(withCount(c, cnt), gs)
	}
}

func GetRunning(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	gs, cnt, err := getFiltered(ctx, c.Param("status"), "", "", "", gType.All)

	if err != nil {
		log.Errorf(ctx, err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		c.Abort()
	}
	withGamers(withCount(c, cnt), gs)
}

//	q := datastore.NewQuery("Game").Ancestor(GamesRoot(c)).Eq("Status", Running).KeysOnly(true)
//
//	var ks []*datastore.Key
//	ctx := restful.ContextFrom(c)
//	if err := datastore.GetAll(ctx, q, &ks); err != nil {
//		log.Errorf(c, "SN/Game#All Error: %s", err)
//		c.Redirect(http.StatusSeeOther, homePath)
//		return
//	}
//
//	length := len(ks)
//	gs := make([]Gamer, length)
//	for i := range gs {
//		var (
//			t  gType.Type
//			ok bool
//		)
//		k := strings.ToLower(ks[i].Parent().Kind())
//		if t, ok = gType.ToType[k]; !ok {
//			log.Errorf(c, "Unknown Game Type For: %s", k)
//			c.Redirect(http.StatusSeeOther, homePath)
//			return
//		}
//		gs[i] = factories[t](c)
//		if ok := datastore.PopulateKey(gs[i], ks[i]); !ok {
//			log.Errorf(c, "Unable to populate gamer with key.")
//			c.Redirect(http.StatusSeeOther, homePath)
//			return
//		}
//	}
//
//	if err := datastore.Get(ctx, gs); err != nil {
//		log.Errorf(c, "SN/Game#All datastore.Get Error: %s", err)
//		c.Redirect(http.StatusSeeOther, homePath)
//		return
//	}
//
//	WithGamers(c, gs)
//}
