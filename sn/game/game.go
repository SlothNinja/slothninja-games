package game

import (
	"github.com/SlothNinja/slothninja-games/sn/color"
	"github.com/SlothNinja/slothninja-games/sn/log"
	sn "github.com/SlothNinja/slothninja-games/sn/misc"
	"github.com/SlothNinja/slothninja-games/sn/restful"
	gType "github.com/SlothNinja/slothninja-games/sn/type"
	"github.com/SlothNinja/slothninja-games/sn/user"
	"github.com/gin-gonic/gin"
	"go.chromium.org/gae/service/datastore"
	"golang.org/x/net/context"
)

type Gamers []Gamer
type Gamer interface {
	//Start(context.Context) error
	PhaseName() string
	FromParams(context.Context, gType.Type) error
	ColorMapFor(*user.User) color.Map
	//NewKey(context.Context, int64) *datastore.Key
	headerer
}

type GetPlayerers interface {
	GetPlayerers() Playerers
}

type hasUpdate interface {
	Update(context.Context) (string, ActionType, error)
}

func GamesRoot(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, "Games", "root", 0, nil)
}

func (h *Header) GetAcceptDialog() bool {
	return h.Private()
}

//func (h *Header) RandomTurnOrder(ps ...Playerer) {
//	for i := 0; i < h.NumPlayers; i++ {
//		ri := sn.MyRand.Intn(h.NumPlayers)
//		ps[i], ps[ri] = ps[ri], ps[i]
//	}
//	h.SetCurrentPlayerers(ps[0])
//}

func (h *Header) RandomTurnOrder() {
	ps := h.gamer.(GetPlayerers).GetPlayerers()
	for i := 0; i < h.NumPlayers; i++ {
		ri := sn.MyRand.Intn(h.NumPlayers)
		ps[i], ps[ri] = ps[ri], ps[i]
	}
	h.SetCurrentPlayerers(ps[0])

	h.OrderIDS = make(UserIndices, len(ps))
	for i, p := range ps {
		h.OrderIDS[i] = p.ID()
	}
}

// Returns (true, nil) if game should be started
func (h *Header) Accept(ctx context.Context, u *user.User) (start bool, err error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Entering")

	if err = h.validateAccept(ctx, u); err != nil {
		return
	}

	h.AddUser(u)

	if len(h.Users) == h.NumPlayers {
		start = true
	}
	return
}

func (h *Header) validateAccept(ctx context.Context, u *user.User) (err error) {
	switch c := restful.GinFrom(ctx); {
	case len(h.UserIDS) >= h.NumPlayers:
		err = sn.NewVError("Game already has the maximum number of players.")
	case h.HasUser(u):
		err = sn.NewVError("%s has already accepted this invitation.", u.Name)
	case h.Password != "" && c.PostForm("password") != h.Password:
		err = sn.NewVError("%s provided incorrect password for Game #%d: %s.", u.Name, h.ID, h.Title)
	}
	return
}

func (h *Header) Drop(u *user.User) (err error) {
	if err = h.validateDrop(u); err != nil {
		return
	}

	h.RemoveUser(u)
	return
}

func (h *Header) validateDrop(u *user.User) (err error) {
	switch {
	case h.Status != Recruiting:
		err = sn.NewVError("Game is no longer recruiting, thus %s can't drop.", u.Name)
	case !h.HasUser(u):
		err = sn.NewVError("%s has not joined this game, thus %s can't drop.", u.Name, u.Name)
	}
	return
}

func RequireCurrentPlayerOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf(c, "Entering RequireCurrentPlayerOrAdmin")
		defer log.Debugf(c, "Exiting RequireCurrentPlayerOrAdmin")

		g := GamerFrom(c)
		if g == nil {
			log.Warningf(c, "Missing Gamer")
			c.Abort()
			return
		}

		ctx := restful.ContextFrom(c)
		if !g.GetHeader().CUserIsCPlayerOrAdmin(ctx) {
			log.Warningf(c, "Current User is Not Current Player or Admin")
			c.Abort()
			return
		}
	}
}
