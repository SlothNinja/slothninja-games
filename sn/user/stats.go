package user

import (
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/slothninja-games/sn/restful"
	"github.com/gin-gonic/gin"
)

const (
	sKind    = "Stats"
	sName    = "root"
	statsKey = "Stats"
)

func From(c *gin.Context) (Stats, bool) {
	s, ok := c.Value(statsKey).(Stats)
	return s, ok
}

func With(c *gin.Context, s Stats) {
	c.Set(statsKey, s)
}

type Stats struct {
	Key       *datastore.Key `datastore:"__key__"`
	Turns     int
	Duration  time.Duration
	Longest   time.Duration
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s Stats) ID() string {
	if s.Key != nil {
		return s.Key.Name
	}
	return ""
}

type MultiStats []*Stats

func (s Stats) Average() time.Duration {
	if s.Turns == 0 {
		return 0
	}
	return (s.Duration / time.Duration(s.Turns))
}

// // last is time associated with last move in game.
// func (s Stats) Update(c *gin.Context, last time.Time) {
// 	With(c, s.update(c, last))
// }
//
func (s Stats) GetUpdate(c *gin.Context, last time.Time) Stats {
	return s.update(c, last)
}

func (s Stats) update(c *gin.Context, last time.Time) Stats {
	since := time.Since(last)

	s.Turns += 1
	s.Duration += since
	if since > s.Longest {
		s.Longest = s.Duration
	}

	return s
}

func (s Stats) AverageString() string {
	switch d := s.Average(); {
	case d.Minutes() < 60:
		return fmt.Sprintf("%.f minutes", d.Minutes())
	case d.Hours() < 48:
		return fmt.Sprintf("%.1f hours", d.Hours())
	default:
		return fmt.Sprintf("%.1f days", d.Hours()/24)
	}
}

func (s Stats) LongestString() string {
	switch d := s.Longest; {
	case d.Minutes() < 60:
		return fmt.Sprintf("%.f minutes", d.Minutes())
	case d.Hours() < 48:
		return fmt.Sprintf("%.1f hours", d.Hours())
	default:
		return fmt.Sprintf("%.1f days", d.Hours()/24)
	}
}

func (s Stats) SinceLastString() string {
	switch d := time.Since(time.Time(s.UpdatedAt)); {
	case d.Minutes() < 60:
		return fmt.Sprintf("%.f minutes", d.Minutes())
	case d.Hours() < 48:
		return fmt.Sprintf("%.1f hours", d.Hours())
	default:
		return fmt.Sprintf("%.1f days", d.Hours()/24)
	}
}

//func key(c *gin.Context, u *user.User) *datastore.Key {
//	return datastore.NewKey(ctx, kind, name, 0, u.Key)
//}

func NewStats(u User) Stats {
	k := New(u.ID()).Key
	return Stats{Key: datastore.NameKey(sKind, sName, k)}
}

func singleError(err error) error {
	if err == nil {
		return err
	}
	if me, ok := err.(datastore.MultiError); ok {
		return me[0]
	}
	return err
}

var NoStatus = Stats{}

func (u User) Stats(c *gin.Context) (Stats, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return NoStatus, nil
	}

	s := NewStats(u)
	err = dsClient.Get(c, s.Key, &s)
	if err == datastore.ErrNoSuchEntity {
		return s, nil
	}
	return s, err
}

// func ByUsers(c *gin.Context, us []user.OAuth) ([]Stats, error) {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Exiting")
//
// 	dsClient, err := datastore.NewClient(c, "")
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	ss := make([]Stats, len(us))
// 	sk := make([]*datastore.Key, len(us))
// 	for i := range ss {
// 		ss[i] = New(us[i])
// 		sk[i] = ss[i].Key
// 	}
//
// 	err = dsClient.GetMulti(c, sk, ss)
// 	if err == nil {
// 		return ss, nil
// 	}
//
// 	me, ok := err.(errors.MultiError)
// 	if !ok {
// 		return nil, err
// 	}
//
// 	// filter out ErrNoSuchEntity since the entity will not exist if the player has yet to take a turn.
// 	isNil := true
// 	for i, e := range me {
// 		if e != nil {
// 			if e == datastore.ErrNoSuchEntity {
// 				me[i] = nil
// 			} else {
// 				isNil = false
// 			}
// 		}
// 	}
//
// 	if isNil {
// 		return ss, nil
// 	}
// 	return ss, me
// }

func FetchStats(getUser func(*gin.Context) (User, bool)) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		_, found := From(c)
		if found {
			return
		}

		u, found := getUser(c)
		log.Debugf("u: %#v", u)
		if !found {
			restful.AddErrorf(c, "missing user.")
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("missing user."))
			return
		}

		s, err := u.Stats(c)
		if err != nil {
			restful.AddErrorf(c, err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		log.Debugf("stats: %#v", s)
		With(c, s)
	}
}

func FetchedStats(c *gin.Context) (Stats, bool) {
	return From(c)
}
