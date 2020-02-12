package contest

import (
	"fmt"

	"github.com/SlothNinja/slothninja-games/sn/restful"
	gType "github.com/SlothNinja/slothninja-games/sn/type"
	"go.chromium.org/gae/service/datastore"
	"golang.org/x/net/context"
)

const kind = "Contest"

type Contests []*Contest
type Contest struct {
	ctx       context.Context
	ID        int64          `gae:"$id"`
	Parent    *datastore.Key `gae:"$parent"`
	GameID    int64
	Type      gType.Type
	R         float64
	RD        float64
	Outcome   float64
	Applied   bool
	CreatedAt restful.CTime
	UpdatedAt restful.UTime
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

func New(ctx context.Context, pk *datastore.Key, gid int64, t gType.Type, r, rd, outcome float64) *Contest {
	return &Contest{
		Parent:  pk,
		GameID:  gid,
		Type:    t,
		R:       r,
		RD:      rd,
		Outcome: outcome,
	}
}

func GenContests(ctx context.Context, places Places) (cs Contests) {
	for _, rmap := range places {
		for ukey, rs := range rmap {
			for _, r := range rs {
				cs = append(cs, New(ctx, ukey, r.GameID, r.Type, r.R, r.RD, r.Outcome))
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

func UnappliedFor(ctx context.Context, ukey *datastore.Key, t gType.Type) (Contests, error) {
	q := datastore.NewQuery(kind).Ancestor(ukey).Eq("Applied", false).Eq("Type", t).KeysOnly(true)

	var ks []*datastore.Key
	if err := datastore.GetAll(ctx, q, &ks); err != nil {
		return nil, err
	}

	length := len(ks)
	if length == 0 {
		return nil, nil
	}

	cs := make(Contests, length)
	for i := range cs {
		cs[i] = new(Contest)
		if ok := datastore.PopulateKey(cs[i], ks[i]); !ok {
			return nil, fmt.Errorf("Unable to populate contest with key.")
		}
	}
	if err := datastore.Get(ctx, cs); err != nil {
		return nil, err
	}
	return cs, nil
}

type ContestMap map[gType.Type]Contests

func Unapplied(ctx context.Context, ukey *datastore.Key) (ContestMap, error) {
	q := datastore.NewQuery(kind).Ancestor(ukey).Eq("Applied", false).KeysOnly(true)

	var ks []*datastore.Key
	if err := datastore.GetAll(ctx, q, &ks); err != nil {
		return nil, err
	}

	length := len(ks)
	if length == 0 {
		return nil, nil
	}

	cs := make(Contests, length)
	for i := range cs {
		cs[i] = new(Contest)
		if ok := datastore.PopulateKey(cs[i], ks[i]); !ok {
			return nil, fmt.Errorf("Unable to populate contest with key.")
		}
	}

	if err := datastore.Get(ctx, cs); err != nil {
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
