package main

import (
	"bitbucket.org/SlothNinja/restful"
	"github.com/SlothNinja/slothninja-games/welcome"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
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

func main() {
	if appengine.IsDevAppServer() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	store := cookie.NewStore([]byte("secret123"))

	r := gin.New()
	r.Use(
		restful.CTXHandler(),
		restful.TemplateHandler(r),
		//		user.GetGUserHandler,
		//		user.GetCUserHandler,
		sessions.Sessions("sngsession", store),
	)
	// 	r := gin.Default()
	// 	store := cookie.NewStore([]byte("secret"))
	// 	r.Use(sessions.Sessions("mysession", store))

	r.GET("/hello", func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("hello") != "world" {
			session.Set("hello", "world")
			session.Save()
		}

		c.JSON(200, gin.H{"hello": session.Get("hello")})
	})

	// Welcome Page (index.html) route
	welcome.AddRoutes(r)

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
