package user_controller

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/slothninja-games/sn/user"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
)

const (
	welcomePath = "/welcome"
	userNewPath = "/user"
	homePath    = "/"
	uidParam    = "uid"
)

func Index(c *gin.Context) {
	cu, found := user.Current(c)
	c.HTML(http.StatusOK, "user/index", gin.H{
		"Context":   c,
		"VersionID": appengine.VersionID(c),
		"CUser":     cu,
		"CUFound":   found,
	})
}

func Show(c *gin.Context) {
	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		log.Errorf("unable to connect to database: %v", err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	cu, found := user.Current(c)

	uid, err := strconv.ParseInt(c.Param(uidParam), 10, 64)
	if err != nil {
		log.Errorf(err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	u := user.New(uid)
	err = dsClient.Get(c, u.Key, &u)
	if err != nil {
		log.Errorf("unable to pull user and stats from database: %v", err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	st := user.NewStats(u)
	err = dsClient.Get(c, st.Key, &st)
	if err != nil && err != datastore.ErrNoSuchEntity {
		log.Errorf("unable to pull user and stats from database: %v", err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	log.Debugf("stats: %#v", st)

	c.HTML(http.StatusOK, "user/show", gin.H{
		"Context":   c,
		"VersionID": appengine.VersionID(c),
		"User":      u,
		"CUser":     cu,
		"CUFound":   found,
		"Stats":     st,
	})
}

func Edit(c *gin.Context) {
	cu, found := user.Current(c)

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		log.Errorf("unable to connect to database: %v", err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	uid, err := strconv.ParseInt(c.Param(uidParam), 10, 64)
	if err != nil {
		log.Errorf(err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	u := user.New(uid)
	err = dsClient.Get(c, u.Key, &u)
	if err != nil {
		log.Errorf("unable to pull user from database: %v", err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	st := user.NewStats(u)
	err = dsClient.Get(c, st.Key, &st)
	if err != nil && err != datastore.ErrNoSuchEntity {
		log.Errorf("unable to pull user and stats from database: %v", err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	log.Debugf("stats: %#v", st)
	c.HTML(http.StatusOK, "user/edit", gin.H{
		"Context":   c,
		"VersionID": appengine.VersionID(c),
		"User":      u,
		"CUser":     cu,
		"CUFound":   found,
		"Stats":     st,
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
	Joined        time.Time     `json:"joined"`
	Updated       time.Time     `json:"updated"`
	OmitCreatedAt omit          `json:"createdat,omitempty"`
	OmitUpdatedAt omit          `json:"updatedat,omitempty"`
}

func toUserTable(c *gin.Context, us []interface{}) (table *jUserIndex, err error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	table = new(jUserIndex)
	l := len(us)
	table.Data = make([]*jUser, l)

	var (
		u user.User
		// nu *user.NUser
		ok bool
	)

	for i, uinf := range us {
		if u, ok = uinf.(user.User); ok {
			table.Data[i] = &jUser{
				IntID:    u.ID(),
				StringID: "",
				OldID:    0,
				GoogleID: u.GoogleID,
				Name:     u.Name,
				Email:    u.Email,
				Gravatar: user.Gravatar(u),
				Joined:   u.CreatedAt,
				Updated:  u.UpdatedAt,
			}
			// } else if nu, ok = uinf.(*user.NUser); ok {
			// 	table.Data[i] = &jUser{
			// 		IntID:    0,
			// 		StringID: nu.ID(),
			// 		OldID:    nu.OldID,
			// 		GoogleID: nu.GoogleID,
			// 		Name:     nu.Name,
			// 		Email:    nu.Email,
			// 		Gravatar: user.NGravatar(nu),
			// 		Joined:   nu.CreatedAt,
			// 		Updated:  nu.UpdatedAt,
			// 	}
			// } else {
			// 	err = fmt.Errorf("not user")
			// 	return
		}
	}

	if draw, err := strconv.Atoi(c.PostForm("draw")); err != nil {
		return nil, err
	} else {
		table.Draw = draw
	}
	table.RecordsTotal = user.CountFrom(c)
	table.RecordsFiltered = user.CountFrom(c)
	return
}

func JSON(c *gin.Context) {
	// us := user.UsersFrom(c)
	// if data, err := toUserTable(c, us); err != nil {
	// 	c.JSON(http.StatusInternalServerError, fmt.Sprintf("%v", err))
	// } else {
	// 	c.JSON(http.StatusOK, data)
	// }
}

func NewAction(c *gin.Context) {
	session := sessions.Default(c)
	token, ok := user.SessionTokenFrom(session)
	if !ok {
		log.Errorf("missing token")
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	if token.Loaded {
		log.Warningf("%v already has account", token.Name)
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		log.Errorf("unable to connect to datastore: %s", err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	u := user.New(token.ID())

	err = dsClient.Get(c, u.Key, u)
	if err == nil {
		log.Warningf("%v already has account", u.Name)
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	log.Debugf(err.Error())

	// dsClient, err := datastore.NewClient(c, "")
	// if err != nil {
	// 	restful.AddErrorf(c, "unable to connect to datastore")
	// 	c.Redirect(http.StatusSeeOther, homePath)
	// 	return
	// }

	// u := user.NewOAuth(token.ID)
	// gu := user.GUserFrom(ctx)
	// if gu == nil {
	// 	restful.AddErrorf(ctx, "You must be logged in to access this page.")
	// 	c.Redirect(http.StatusSeeOther, welcomePath)
	// 	return
	// }

	u.Name = strings.Split(token.Email, "@")[0]
	u.LCName = strings.ToLower(u.Name)
	u.Email = token.Email

	c.HTML(http.StatusOK, "user/new", gin.H{
		"Context": c,
		"User":    u,
	})
}

func Create(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// u := user.FromGUser(ctx, user.GUserFrom(ctx))
		// switch existing, err := user.ByGoogleID(ctx, u.GoogleID); {
		// case err == user.ErrNotFound:
		// case err != nil:
		// 	restful.AddErrorf(ctx, err.Error())
		// 	c.Redirect(http.StatusSeeOther, userNewPath)
		// 	return
		// case existing != nil:
		// 	restful.AddErrorf(ctx, "You already have an account.")
		// 	c.Redirect(http.StatusSeeOther, homePath)
		// 	return
		// default:
		// 	log.Errorf(ctx, "Unexpected result for user.Create. err: %v existing: %v", err, existing)
		// 	c.Redirect(http.StatusSeeOther, userNewPath)
		// 	return
		// }

		session := sessions.Default(c)
		token, ok := user.SessionTokenFrom(session)
		if !ok {
			log.Errorf("missing token")
			c.Redirect(http.StatusSeeOther, homePath)
			return
		}

		if token.Loaded {
			log.Warningf("%v already has account", token.Name)
			c.Redirect(http.StatusSeeOther, homePath)
			return
		}

		dsClient, err := datastore.NewClient(c, "")
		if err != nil {
			log.Errorf("unable to connect to datastore: %s", err.Error())
			c.Redirect(http.StatusSeeOther, homePath)
			return
		}

		u := user.New(token.ID())
		err = dsClient.Get(c, u.Key, u)
		if err == nil {
			log.Warningf("%v already has account", u.Name)
			c.Redirect(http.StatusSeeOther, homePath)
			return
		}

		// Fell through 'switch' thus err == user.ErrNotFound
		u.Name = strings.Split(c.PostForm("user-name"), "@")[0]
		u.LCName = strings.ToLower(u.Name)
		u.Email = token.Email
		//u.Key = user.NewKey(ctx, 0)

		// n := name.New(u.LCName)
		// if !name.IsUnique(c, u.LCName) {
		// 	restful.AddErrorf(c, "%q is not a unique user name.", u.LCName)
		// 	c.Redirect(http.StatusSeeOther, userNewPath)
		// 	return
		// }

		// n.GoogleID = u.GoogleID

		t := time.Now()
		u.CreatedAt, u.Joined, u.UpdatedAt = t, t, t
		_, err = dsClient.Put(c, u.Key, u)
		if err != nil {
			log.Errorf(err.Error())
			c.Redirect(http.StatusSeeOther, homePath)
			return
		}

		c.Redirect(http.StatusSeeOther, showPath(prefix, u.ID()))
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
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		log.Errorf(err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	uid, err := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		log.Errorf(err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	// Get Resource
	u := user.New(uid)

	err = dsClient.Get(c, u.Key, &u)
	if err != nil {
		log.Errorf("User/Controller#Update user.BySID Error: %s", err)
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	log.Debugf("before update u: %#v", u)
	u, err = u.Update(c)
	log.Debugf("after update u: %#v", u)
	if err != nil {
		log.Errorf("User/Controller#Update u.update Error: %s", err)
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	_, err = dsClient.Put(c, u.Key, &u)
	if err != nil {
		log.Errorf(err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}
	route := fmt.Sprintf("/user/show/%s", c.Param("uid"))
	c.Redirect(http.StatusSeeOther, route)
}

// func GamesIndex(c *gin.Context) {
// 	ctx := restful.ContextFrom(c)
// 	log.Debugf(ctx, "Entering")
// 	defer log.Debugf(ctx, "Exiting")
//
// 	if status := game.StatusFrom(ctx); status != game.NoStatus {
// 		c.HTML(200, "shared/games_index", gin.H{})
// 	} else {
// 		c.HTML(200, "user/games_index", gin.H{})
// 	}
// }
