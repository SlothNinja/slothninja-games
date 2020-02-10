package user_controller

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"bitbucket.org/SlothNinja/slothninja-games/sn"
	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/name"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"bitbucket.org/SlothNinja/slothninja-games/sn/user"
	"bitbucket.org/SlothNinja/slothninja-games/sn/user/stats"
	"github.com/gin-gonic/gin"
	"go.chromium.org/gae/service/datastore"
	"go.chromium.org/gae/service/info"
	"golang.org/x/net/context"
)

const (
	welcomePath = "/welcome"
	userNewPath = "/user"
	homePath    = "/"
)

func Index(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	c.HTML(http.StatusOK, "user/index", gin.H{
		"Context":   ctx,
		"VersionID": info.VersionID(ctx),
		"CUser":     user.CurrentFrom(c),
	})
}

func Show(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	u := user.From(ctx)
	c.HTML(http.StatusOK, "user/show", gin.H{
		"Context":   ctx,
		"VersionID": info.VersionID(ctx),
		"User":      u,
		"CUser":     user.CurrentFrom(ctx),
		"IsAdmin":   user.IsAdmin(ctx),
		"Stats":     stats.Fetched(ctx),
	})
}

func Edit(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	u := user.From(ctx)
	c.HTML(http.StatusOK, "user/edit", gin.H{
		"Context":   ctx,
		"VersionID": info.VersionID(ctx),
		"User":      u,
		"CUser":     user.CurrentFrom(ctx),
		"IsAdmin":   user.IsAdmin(ctx),
		"Stats":     stats.Fetched(ctx),
	})
}

//func Remote(ctx *restful.Context, render render.Render, params martini.Params) {
//	if u, err := user.ByGoogleID(ctx, params["uid"]); err == nil {
//		render.JSON(http.StatusOK, u)
//	} else {
//		render.HTML(http.StatusGone, "", "")
//	}
//}

type jUserIndex struct {
	Data            []*jUser `json:"data"`
	Draw            int      `json:"draw"`
	RecordsTotal    int64    `json:"recordsTotal"`
	RecordsFiltered int64    `json:"recordsFiltered"`
}

type omit *struct{}

type jUser struct {
	IntID         int64         `json:"id"`
	StringID      string        `json:"sid"`
	OldID         int64         `json:"oldid"`
	GoogleID      string        `json:"googleid"`
	Name          string        `json:"name"`
	Email         string        `json:"email"`
	Gravatar      template.HTML `json:"gravatar"`
	Joined        restful.CTime `json:"joined"`
	Updated       restful.UTime `json:"updated"`
	OmitCreatedAt omit          `json:"createdat,omitempty"`
	OmitUpdatedAt omit          `json:"updatedat,omitempty"`
}

func toUserTable(c *gin.Context, us []interface{}) (table *jUserIndex, err error) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	table = new(jUserIndex)
	l := len(us)
	table.Data = make([]*jUser, l)

	var (
		u  *user.User
		nu *user.NUser
		ok bool
	)

	for i, uinf := range us {
		if u, ok = uinf.(*user.User); ok {
			table.Data[i] = &jUser{
				IntID:    u.ID,
				StringID: "",
				OldID:    0,
				GoogleID: u.GoogleID,
				Name:     u.Name,
				Email:    u.Email,
				Gravatar: user.Gravatar(u),
				Joined:   u.CreatedAt,
				Updated:  u.UpdatedAt,
			}
		} else if nu, ok = uinf.(*user.NUser); ok {
			table.Data[i] = &jUser{
				IntID:    0,
				StringID: nu.ID,
				OldID:    nu.OldID,
				GoogleID: nu.GoogleID,
				Name:     nu.Name,
				Email:    nu.Email,
				Gravatar: user.NGravatar(nu),
				Joined:   nu.CreatedAt,
				Updated:  nu.UpdatedAt,
			}
		} else {
			err = fmt.Errorf("not user")
			return
		}
	}

	if draw, err := strconv.Atoi(c.PostForm("draw")); err != nil {
		return nil, err
	} else {
		table.Draw = draw
	}
	table.RecordsTotal = user.CountFrom(ctx)
	table.RecordsFiltered = user.CountFrom(ctx)
	return
}

func JSON(c *gin.Context) {
	us := user.UsersFrom(c)
	if data, err := toUserTable(c, us); err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("%v", err))
	} else {
		c.JSON(http.StatusOK, data)
	}
}

func NewAction(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	u := user.New(ctx)
	gu := user.GUserFrom(ctx)
	if gu == nil {
		restful.AddErrorf(ctx, "You must be logged in to access this page.")
		c.Redirect(http.StatusSeeOther, welcomePath)
		return
	}

	u.Name = strings.Split(gu.Email, "@")[0]
	u.LCName = strings.ToLower(u.Name)
	u.Email = gu.Email

	c.HTML(http.StatusOK, "user/new", gin.H{
		"Context": ctx,
		"User":    user.FromGUser(ctx, user.GUserFrom(ctx)),
	})
}

func Create(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.ContextFrom(c)
		u := user.FromGUser(ctx, user.GUserFrom(ctx))
		switch existing, err := user.ByGoogleID(ctx, u.GoogleID); {
		case err == user.ErrNotFound:
		case err != nil:
			restful.AddErrorf(ctx, err.Error())
			c.Redirect(http.StatusSeeOther, userNewPath)
			return
		case existing != nil:
			restful.AddErrorf(ctx, "You already have an account.")
			c.Redirect(http.StatusSeeOther, homePath)
			return
		default:
			log.Errorf(ctx, "Unexpected result for user.Create. err: %v existing: %v", err, existing)
			c.Redirect(http.StatusSeeOther, userNewPath)
			return
		}

		// Fell through 'switch' thus err == user.ErrNotFound
		u.Name = strings.Split(c.PostForm("user-name"), "@")[0]
		u.LCName = strings.ToLower(u.Name)
		//u.Key = user.NewKey(ctx, 0)

		n := name.New()
		if !name.IsUnique(ctx, u.LCName) {
			restful.AddErrorf(ctx, "%q is not a unique user name.", u.LCName)
			c.Redirect(http.StatusSeeOther, userNewPath)
			return
		}

		n.GoogleID = u.GoogleID
		n.ID = u.LCName

		err := datastore.RunInTransaction(ctx, func(tc context.Context) (terr error) {
			entities := []interface{}{u, n}
			if terr = datastore.Put(tc, entities); terr != nil {
				return
			}
			nu := user.ToNUser(ctx, u)
			return datastore.Put(tc, nu)
		}, &datastore.TransactionOptions{XG: true})

		if err != nil {
			log.Errorf(ctx, "User/Controller#Create datastore.RunInTransaction Error: %v", err)
			c.Redirect(http.StatusSeeOther, homePath)
			return
		}

		c.Redirect(http.StatusSeeOther, showPath(prefix, u.ID))
	}
}

func showPath(prefix string, id int64) string {
	return fmt.Sprintf("%s/show/%d", prefix, id)
}

//func SendTestMessage(ctx *restful.Context, render render.Render, routes martini.Routes, params martini.Params) {
//	u := user.Fetched(ctx)
//	m := new(xmpp.Message)
//	m.To = []string{u.Email}
//	m.Body = fmt.Sprintf("Test message from SlothNinja Games for %s", u.Name)
//	send.XMPP(ctx, m)
//	ctx.AddNoticef("Test IM sent to %s", u.Name)
//	render.Redirect(routes.URLFor("user_show", params["uid"]), http.StatusSeeOther)
//}
//
//func SendIMInvite(ctx *restful.Context, render render.Render, routes martini.Routes, params martini.Params) {
//	u := user.Fetched(ctx)
//	send.Invite(ctx, u.Email)
//	ctx.AddNoticef("IM Invite sent to %s", u.Name)
//	render.Redirect(routes.URLFor("user_show", params["uid"]), http.StatusSeeOther)
//}
func Update(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	log.Debugf(ctx, "Request: %#v", c.Request)
	var err error

	// Get Resource
	u := user.New(ctx)
	if uid, err := strconv.ParseInt(c.Param("uid"), 10, 64); err != nil {
		log.Errorf(ctx, err.Error())
		return
	} else {
		u.ID = uid
	}

	if err := datastore.Get(ctx, u); err != nil {
		log.Errorf(ctx, "User/Controller#Update user.BySID Error: %s", err)
		return
	}
	oldName := name.New()
	oldName.ID = u.LCName
	if err := u.Update(ctx); err != nil {
		log.Errorf(ctx, "User/Controller#Update u.update Error: %s", err)
		restful.AddErrorf(ctx, err.Error())
		route := fmt.Sprintf("/user/show/%s", c.Param("uid"))
		c.Redirect(http.StatusSeeOther, route)
		return
	}
	newName := name.New()
	newName.GoogleID = u.GoogleID
	newName.ID = u.LCName

	log.Debugf(ctx, "Before datastore.RunInTransaction")
	err = datastore.RunInTransaction(ctx, func(tc context.Context) (err error) {
		nu := user.ToNUser(ctx, u)
		entities := []interface{}{u, nu, newName, oldName}
		if err = datastore.Put(tc, entities); err != nil {
			return
		}

		return datastore.Delete(tc, oldName)
	}, &datastore.TransactionOptions{XG: true})

	log.Debugf(ctx, "error: %v", err)

	switch {
	case sn.IsVError(err):
		restful.AddErrorf(ctx, err.Error())
	case err != nil:
		log.Errorf(ctx, err.Error())
	}

	route := fmt.Sprintf("/user/show/%s", c.Param("uid"))
	c.Redirect(http.StatusSeeOther, route)
}

func GamesIndex(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	if status := game.StatusFrom(ctx); status != game.NoStatus {
		c.HTML(200, "shared/games_index", gin.H{})
	} else {
		c.HTML(200, "user/games_index", gin.H{})
	}
}
