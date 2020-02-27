package game

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/log"
	gType "github.com/SlothNinja/slothninja-games/sn/type"
	"github.com/gin-gonic/gin"
)

const (
	statusKey = "Status"
	countKey  = "Count"
	NoCount   = -1
)

func getAllQuery() *datastore.Query {
	return datastore.NewQuery("Game").Ancestor(GamesRoot())
}

func getFiltered(c *gin.Context, status, sid, start, length string, t gType.Type) (Gamers, int, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	q := getAllQuery().KeysOnly()
	if status != "" {
		st := ToStatus[strings.ToLower(status)]
		q = q.Filter("Status=", int(st))
		WithStatus(c, st)
	}

	if sid != "" {
		if id, err := strconv.Atoi(sid); err == nil {
			q = q.Filter("UserIDS=", id)
		}
	}

	if t != gType.All {
		q = q.
			Filter("Type=", int(t)).
			Order("-UpdatedAt")
	} else {
		q = q.Order("-UpdatedAt")
	}

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		log.Errorf("sn/game#GetFiltered q.Count Error: %s", err)
		return nil, -1, err
	}

	cnt, err := dsClient.Count(c, q)
	if err != nil {
		log.Errorf("sn/game#GetFiltered q.Count Error: %s", err)
		return nil, -1, err
	}

	if start != "" {
		st, err := strconv.ParseInt(start, 10, 32)
		if err == nil {
			q = q.Offset(int(st))
		}
	}

	if length != "" {
		l, err := strconv.ParseInt(length, 10, 32)
		if err == nil {
			q = q.Limit(int(l))
		}
	}

	ks, err := dsClient.GetAll(c, q, nil)
	if err != nil {
		log.Errorf("getFiltered GetAll Error: %s", err)
		return nil, -1, err
	}

	l := len(ks)
	gs := make([]Gamer, l)
	hs := make([]*Header, l)
	for i := range gs {
		var ok bool
		if t == gType.All {
			k := strings.ToLower(ks[i].Parent.Kind)
			if t, ok = gType.ToType[k]; !ok {
				err = fmt.Errorf("Unknown Game Type For: %s", k)
				log.Errorf(err.Error())
				return nil, -1, err
			}
		}
		gs[i] = factories[t](ks[i].ID)
		hs[i] = gs[i].GetHeader()
	}

	err = dsClient.GetMulti(c, ks, hs)
	if err != nil {
		log.Errorf("SN/Game#GetFiltered datastore.Get Error: %s", err)
		return nil, -1, err
	}

	for i := range hs {
		err = hs[i].AfterLoad(c, gs[i])
		if err != nil {
			log.Errorf("h.AfterLoad error: %s", err)
			return nil, -1, err
		}
	}

	return gs, cnt, nil
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
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	c.Set(statusKey, s)
}

func StatusFrom(c *gin.Context) (s Status) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	if s = ToStatus[strings.ToLower(c.Param("status"))]; s != NoStatus {
		WithStatus(c, s)
	} else {
		s, _ = c.Value(statusKey).(Status)
	}
	return
}

func withCount(c *gin.Context, cnt int) *gin.Context {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	c.Set(countKey, cnt)
	return c
}

func countFrom(c *gin.Context) int {
	cnt, _ := c.Value(countKey).(int)
	return cnt
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
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		gs, cnt, err := getFiltered(
			c,
			c.Param("status"),
			c.Param("uid"),
			c.PostForm("start"),
			c.PostForm("length"),
			t,
		)

		log.Debugf("gs: %#v\ncnt: %d\nerr: %v", gs, cnt, err)

		if err != nil {
			log.Errorf(err.Error())
			c.Redirect(http.StatusSeeOther, homePath)
			c.Abort()
		}
		withGamers(withCount(c, cnt), gs)
	}
}

func GetRunning(c *gin.Context) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	gs, cnt, err := getFiltered(c, c.Param("status"), "", "", "", gType.All)

	if err != nil {
		log.Errorf(err.Error())
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
