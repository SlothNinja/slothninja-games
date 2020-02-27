package game

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/slothninja-games/sn/color"
	sn "github.com/SlothNinja/slothninja-games/sn/misc"
	gType "github.com/SlothNinja/slothninja-games/sn/type"
	"github.com/SlothNinja/slothninja-games/sn/user"
	"github.com/gin-gonic/gin"
)

type Gamers []Gamer
type Gamer interface {
	//Start(context.Context) error
	PhaseName() string
	FromParams(*gin.Context, gType.Type) error
	ColorMapFor(*gin.Context, user.User) color.Map
	//NewKey(context.Context, int64) *datastore.Key
	headerer
}

type GetPlayerers interface {
	GetPlayerers() Playerers
}

type hasUpdate interface {
	Update(context.Context) (string, ActionType, error)
}

func GamesRoot() *datastore.Key {
	return datastore.NameKey("Games", "root", nil)
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
func (h *Header) Accept(c *gin.Context, u user.User) (bool, error) {
	log.Debugf("Entering")
	defer log.Debugf("Entering")

	err := h.validateAccept(c, u)
	if err != nil {
		return false, err
	}

	h.AddUser(u)

	if len(h.Users) == h.NumPlayers {
		return true, nil
	}
	return false, nil
}

func (h Header) validateAccept(c *gin.Context, u user.User) error {
	switch {
	case len(h.UserIDS) >= h.NumPlayers:
		return sn.NewVError("Game already has the maximum number of players.")
	case h.HasUser(u):
		return sn.NewVError("%s has already accepted this invitation.", u.Name)
	case h.Password != "" && c.PostForm("password") != h.Password:
		return sn.NewVError("%s provided incorrect password for Game #%d: %s.", u.Name, h.ID, h.Title)
	}
	return nil
}

func (h *Header) Drop(u user.User) error {
	err := h.validateDrop(u)
	if err != nil {
		return err
	}

	h.RemoveUser(u)
	return nil
}

func (h Header) validateDrop(u user.User) error {
	switch {
	case h.Status != Recruiting:
		return sn.NewVError("Game is no longer recruiting, thus %s can't drop.", u.Name)
	case !h.HasUser(u):
		return sn.NewVError("%s has not joined this game, thus %s can't drop.", u.Name, u.Name)
	}
	return nil
}

func RequireCurrentPlayerOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		g := GamerFrom(c)
		if g == nil {
			log.Warningf("Missing Gamer")
			c.Abort()
			return
		}

		if !g.GetHeader().CUserIsCPlayerOrAdmin(c) {
			log.Warningf("Current User is Not Current Player or Admin")
			c.Abort()
			return
		}
	}
}
