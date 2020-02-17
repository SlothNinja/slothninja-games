package user_controller

import (
	"github.com/SlothNinja/slothninja-games/sn/user"
	"github.com/gin-gonic/gin"
)

const (
	authPath  = "auth"
	loginPath = "login"
)

func AddRoutes(prefix string, engine *gin.Engine) {
	// User Group
	g1 := engine.Group(prefix)
	g1.GET("new",
		// user.RequireLogin(),
		// gType.SetTypes(),
		NewAction,
	)

	// Create
	g1.POST("",
		// user.RequireLogin(),
		Create(prefix),
	)

	// Show User
	g1.GET("show/:uid",
		// user.Fetch,
		// stats.Fetch(user.From),
		// gType.SetTypes(),
		Show,
	)

	// Edit User
	g1.GET("edit/:uid",
		// user.RequireLogin(),
		// user.Fetch,
		// stats.Fetch(user.From),
		// gType.SetTypes(),
		Edit,
	)

	// Update User
	g1.POST("update/:uid",
		// user.RequireLogin(),
		// user.Fetch,
		// gType.SetTypes(),
		Update,
	)

	// User Ratings
	// g1.POST("show/:uid/ratings/json",
	// 	user.Fetch,
	// 	rating.JSONIndexAction,
	// )

	// g1.POST("edit/:uid/ratings/json",
	// 	user.RequireLogin(),
	// 	user.Fetch,
	// 	rating.JSONIndexAction,
	// )

	// User Games
	// g1.POST("show/:uid/games/json",
	// 	gType.SetTypes(),
	// 	game.GetFiltered(gType.All),
	// 	game.JSONIndexAction,
	// )

	// g1.POST("edit/:uid/games/json",
	// 	user.RequireLogin(),
	// 	gType.SetTypes(),
	// 	game.GetFiltered(gType.All),
	// 	game.JSONIndexAction,
	// )

	g1.GET("login",
		user.Login("/user/auth"),
	)

	g1.GET("auth", user.Auth("/user/auth"))

	// Users group
	g2 := engine.Group(prefix + "s")

	// Index
	g2.GET("",
		// user.RequireAdmin,
		// gType.SetTypes(),
		Index,
	)

	// json data for Index
	g2.POST("/json",
		// user.RequireAdmin,
		// gType.SetTypes(),
		// user.FetchAll,
		JSON,
	)

}
