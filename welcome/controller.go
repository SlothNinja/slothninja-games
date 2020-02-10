package welcome

import (
	"net/http"

	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"bitbucket.org/SlothNinja/slothninja-games/sn/user"
	"github.com/gin-gonic/gin"
	"go.chromium.org/gae/service/info"
)

func Index(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering welcome#Index")
	defer log.Debugf(ctx, "Exiting welcome#Index")

	if cu, gu := user.CurrentFrom(ctx), user.GUserFrom(ctx); cu == nil && gu != nil {
		c.Redirect(http.StatusSeeOther, "/user/new")
	} else {
		log.Debugf(ctx, "cu: %#v", cu)
		c.HTML(http.StatusOK, "welcome/index", gin.H{
			"VersionID": info.VersionID(ctx),
			"CUser":     cu,
			"Context":   ctx})
	}
}
