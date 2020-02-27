package contest

import (
	"time"

	"cloud.google.com/go/datastore"
	gType "github.com/SlothNinja/slothninja-games/sn/type"
	"github.com/gin-gonic/gin"
)

const kind = "Contest"

// type Contests []*Contest
type Contest struct {
	Key       *datastore.Key `datastore:"__key__"`
	GameID    int64
	Type      gType.Type
	R         float64
	RD        float64
	Outcome   float64
	Applied   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Result struct {
	GameID  int64
	Type    gType.Type
	R       float64
	RD      float64
	Outcome float64
}

type Results []*Result
type ResultsMap map[*datastore.Key]Results
type Places []ResultsMap

func New(pk *datastore.Key, gid int64, t gType.Type, r, rd, outcome float64) Contest {
	return Contest{
		Key:     datastore.IncompleteKey(kind, pk),
		GameID:  gid,
		Type:    t,
		R:       r,
		RD:      rd,
		Outcome: outcome,
	}
}

func GenContests(c *gin.Context, places Places) (cs []Contest) {
	for _, rmap := range places {
		for ukey, rs := range rmap {
			for _, r := range rs {
				cs = append(cs, New(ukey, r.GameID, r.Type, r.R, r.RD, r.Outcome))
			}
		}
	}
	return
}

//func (cs Contests) Keys() []*datastore.Key {
//	ks := make([]*datastore.Key, len(cs))
//	for i, c := range cs {
//		ks[i] = c.Key
//	}
//	return ks
//}

//func (cs Contests) Save(ctx *restful.Context) error {
//	keys := make([]*datastore.Key, len(cs))
//	now := time.Now()
//	for i, c := range cs {
//		keys[i] = c.Key
//		if c.CreatedAt.IsZero() {
//			c.CreatedAt = now
//		}
//		c.UpdatedAt = now
//	}
//	return datastore.RunInTransaction(ctx, func(tc appengine.Context) error {
//		_, err := gaelic.PutMulti(tc, keys, cs)
//		return err
//	}, nil)
//}

func UnappliedFor(c *gin.Context, ukey *datastore.Key, t gType.Type) ([]Contest, error) {
	q := datastore.NewQuery(kind).
		Ancestor(ukey).
		Filter("Applied=", false).
		Filter("Type=", t)

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return nil, err
	}

	var cs []Contest
	_, err = dsClient.GetAll(c, q, cs)
	if err != nil {
		return nil, err
	}

	return cs, nil
}

type ContestMap map[gType.Type][]Contest

func Unapplied(c *gin.Context, ukey *datastore.Key) (ContestMap, error) {
	q := datastore.NewQuery(kind).
		Ancestor(ukey).
		Filter("Applied=", false)

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return nil, err
	}

	var cs []Contest
	_, err = dsClient.GetAll(c, q, cs)
	if err != nil {
		return nil, err
	}

	cm := make(ContestMap, len(gType.Types))
	for _, c := range cs {
		c.Applied = true
		cm[c.Type] = append(cm[c.Type], c)
	}
	return cm, nil
}

//func (c *Contest) Save(ch chan<- datastore.Property) error {
//	// Time stamp
//	t := time.Now()
//	if c.CreatedAt.IsZero() {
//		c.CreatedAt = t
//	}
//	c.UpdatedAt = t
//	return datastore.SaveStruct(c, ch)
//}
//
//func (c *Contest) Load(ch <-chan datastore.Property) error {
//	return datastore.LoadStruct(c, ch)
//}
