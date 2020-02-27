package main

import (
	got "github.com/SlothNinja/gt"
	"github.com/SlothNinja/slothninja-games/sn/game"
	"github.com/SlothNinja/slothninja-games/sn/restful"
	gType "github.com/SlothNinja/slothninja-games/sn/type"
	"github.com/SlothNinja/slothninja-games/sn/user_controller"
	"github.com/SlothNinja/slothninja-games/welcome"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

const (
	userPrefix   = "user"
	gamesPrefix  = "games"
	ratingPrefix = "rating"
	mailPrefix   = "mail"
	rootPath     = "/"
)

// func main() {
// 	// if appengine.IsDevAppServer() {
// 	// 	gin.SetMode(gin.DebugMode)
// 	// } else {
// 	// 	gin.SetMode(gin.ReleaseMode)
// 	// }
//
// 	store := cookie.NewStore([]byte("secret123"))
//
// 	r := gin.New()
// 	r.Use(
// 		// restful.CTXHandler(),
// 		// restful.TemplateHandler(r),
// 		// user.GetGUserHandler,
// 		// user.GetCUserHandler,
// 		sessions.Sessions("sngsession", store),
// 	)
//
// 	// Welcome Page (index.html) route
// 	// welcome.AddRoutes(r)
// }

const (
	hashKeyLength  = 64
	blockKeyLength = 32
	sessionName    = "sngsession"
)

func main() {
	// if appengine.IsDevAppServer() {
	// 	gin.SetMode(gin.DebugMode)
	// } else {
	// 	gin.SetMode(gin.ReleaseMode)
	// }

	hashKey := securecookie.GenerateRandomKey(hashKeyLength)
	if hashKey == nil {
		panic("generated hashKey was nil")
	}

	blockKey := securecookie.GenerateRandomKey(blockKeyLength)
	if blockKey == nil {
		panic("generated blockKey was nil")
	}

	store := cookie.NewStore(hashKey, blockKey)

	renderer := restful.ParseTemplates("templates", ".tmpl")
	r := gin.Default()
	r.Use(
		restful.TemplateHandler(r, renderer),
		sessions.Sessions(sessionName, store),
	)

	// Welcome Page (index.html) route
	welcome.AddRoutes(r)

	// User Routes
	user_controller.AddRoutes(userPrefix, r)

	// Games Routes
	game.AddRoutes(gamesPrefix, r)

	// Guild of Thieves
	got.Register(gType.GOT, r)

	r.Run()
}

// func init() {
// 	if appengine.IsDevAppServer() {
// 		gin.SetMode(gin.DebugMode)
// 	} else {
// 		gin.SetMode(gin.ReleaseMode)
// 	}
//
// 	store := cookie.NewStore([]byte("secret123"))
//
// 	r := gin.New()
// 	r.Use(
// 		restful.CTXHandler(),
// 		restful.TemplateHandler(r),
// 		user.GetGUserHandler,
// 		user.GetCUserHandler,
// 		sessions.Sessions("sngsession", store),
// 	)
//
// 	// Welcome Page (index.html) route
// 	welcome.AddRoutes(r)
//
// 	// Mail route
// 	send.AddRoutes(mailPrefix, r)
//
// 	// Games Routes
// 	game.AddRoutes(gamesPrefix, r)
//
// 	// User Routes
// 	user_controller.AddRoutes(userPrefix, r)
//
// 	// Rating Routes
// 	rating.AddRoutes(ratingPrefix, r)
//
// 	// After The Flood
// 	atf.Register(gType.ATF, r)
//
// 	// Guild of Thieves
// 	got.Register(gType.GOT, r)
//
// 	// Tammany Hall
// 	tammany.Register(gType.Tammany, r)
//
// 	// Indonesia
// 	indonesia.Register(gType.Indonesia, r)
//
// 	// Confucius
// 	confucius.Register(gType.Confucius, r)
//
// 	http.Handle(rootPath, r)
// }
