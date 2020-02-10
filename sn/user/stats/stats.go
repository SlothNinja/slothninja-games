package stats

import (
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"bitbucket.org/SlothNinja/slothninja-games/sn/user"
	"github.com/gin-gonic/gin"
	"go.chromium.org/gae/service/datastore"
	"go.chromium.org/luci/common/errors"
	"golang.org/x/net/context"
)

const (
	kind     = "Stats"
	name     = "root"
	statsKey = "Stats"
	homePath = "/"
)

func From(ctx context.Context) (s *Stats) {
	s, _ = ctx.Value(statsKey).(*Stats)
	return
}

func With(c *gin.Context, s *Stats) {
	c.Set(statsKey, s)
}

type Stats struct {
	ID        string         `gae:"$id"`
	Parent    *datastore.Key `gae:"$parent"`
	Turns     int
	Duration  time.Duration
	Longest   time.Duration
	CreatedAt restful.CTime
	UpdatedAt restful.UTime
}

type MultiStats []*Stats

func (s *Stats) Average() time.Duration {
	if s.Turns == 0 {
		return 0
	}
	return (s.Duration / time.Duration(s.Turns))
}

// last is time associated with last move in game.
func (s *Stats) Update(c *gin.Context, last time.Time) {
	With(c, s.update(c, last))
}

func (s *Stats) GetUpdate(ctx context.Context, last time.Time) *Stats {
	return s.update(ctx, last)
}

func (s *Stats) update(ctx context.Context, last time.Time) *Stats {
	since := time.Since(last)

	s.Turns += 1
	s.Duration += since
	if since > s.Longest {
		s.Longest = s.Duration
	}

	return s
}

func (s *Stats) AverageString() string {
	switch d := s.Average(); {
	case d.Minutes() < 60:
		return fmt.Sprintf("%.f minutes", d.Minutes())
	case d.Hours() < 48:
		return fmt.Sprintf("%.1f hours", d.Hours())
	default:
		return fmt.Sprintf("%.1f days", d.Hours()/24)
	}
}

func (s *Stats) LongestString() string {
	switch d := s.Longest; {
	case d.Minutes() < 60:
		return fmt.Sprintf("%.f minutes", d.Minutes())
	case d.Hours() < 48:
		return fmt.Sprintf("%.1f hours", d.Hours())
	default:
		return fmt.Sprintf("%.1f days", d.Hours()/24)
	}
}

func (s *Stats) SinceLastString() string {
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

func New(ctx context.Context, u *user.User) *Stats {
	return &Stats{ID: name, Parent: datastore.KeyForObj(ctx, u)}
}

func singleError(err error) error {
	if err == nil {
		return err
	}
	if me, ok := err.(errors.MultiError); ok {
		return me[0]
	}
	return err
}

func ByUser(ctx context.Context, u *user.User) (s *Stats, err error) {
	s = New(ctx, u)
	if err = datastore.Get(ctx, s); err == datastore.ErrNoSuchEntity {
		err = nil
	}
	return
}

func ByUsers(ctx context.Context, us user.Users) (ss []*Stats, err error) {
	ss = make([]*Stats, len(us))
	for i := range ss {
		ss[i] = New(ctx, us[i])
	}

	if err = datastore.Get(ctx, ss); err == nil {
		return
	}

	var (
		me errors.MultiError
		ok bool
	)
	if me, ok = err.(errors.MultiError); !ok {
		return
	}

	// filter out ErrNoSuchEntity since the entity will not exist if the player has yet to take a turn.
	isNil := true
	for i, e := range me {
		if e != nil {
			if e == datastore.ErrNoSuchEntity {
				me[i] = nil
			} else {
				isNil = false
			}
		}
	}

	if isNil {
		err = nil
	}
	return
}

func Fetch(getUser func(context.Context) *user.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.ContextFrom(c)
		log.Debugf(ctx, "Entering")
		defer log.Debugf(ctx, "Exiting")

		if From(ctx) != nil {
			return
		}

		u := getUser(ctx)
		log.Debugf(ctx, "u: %#v", u)
		if u == nil {
			restful.AddErrorf(ctx, "missing user.")
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("missing user."))
			return
		}

		if s, err := ByUser(ctx, u); err != nil {
			restful.AddErrorf(ctx, err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			log.Debugf(ctx, "stats: %#v", s)
			With(c, s)
		}
	}
}

func Fetched(ctx context.Context) *Stats {
	return From(ctx)
}
