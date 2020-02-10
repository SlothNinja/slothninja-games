package send

import (
	"github.com/gin-gonic/gin"
)

//func AddRoutes(prefix string, r martini.Router) {
//	group := "/" + prefix
//	r.Group(group, func(r martini.Router) {
//		r.Post("mail"). //			Mail(),
//				Name(prefix)
//	})
//}

func AddRoutes(prefix string, engine *gin.Engine) {
	engine.Group(prefix).POST("", Mail)
}
