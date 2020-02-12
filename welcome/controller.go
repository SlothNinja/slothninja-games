package welcome

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	// ctx := restful.ContextFrom(c)
	// log.Debugf(ctx, "Entering welcome#Index")
	// defer log.Debugf(ctx, "Exiting welcome#Index")

	// if cu, gu := user.CurrentFrom(ctx), user.GUserFrom(ctx); cu == nil && gu != nil {
	// 	c.Redirect(http.StatusSeeOther, "/user/new")
	// } else {
	// 	log.Debugf(ctx, "cu: %#v", cu)
	c.HTML(http.StatusOK, "welcome/index", gin.H{
		// "VersionID": info.VersionID(ctx),
		// "CUser":     cu,
		// "Context": ctx})
		"Context": c})
	// }
}
