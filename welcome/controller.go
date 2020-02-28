package welcome

import (
	"net/http"

	"github.com/SlothNinja/log"
	"github.com/SlothNinja/slothninja-games/sn/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
)

func Index(c *gin.Context) {
	log.Debugf("Entering welcome#Index")
	defer log.Debugf("Exiting welcome#Index")

	cu := user.CurrentFrom(c)
	log.Debugf("cu: %#v", cu)
	c.HTML(http.StatusOK, "welcome/index", gin.H{
		"VersionID": appengine.VersionID(c),
		"CUser":     cu,
		"Context":   c})
	// if cu, gu := user.CurrentFrom(c), user.GUserFrom(c); cu == nil && gu != nil {
	// 	c.Redirect(http.StatusSeeOther, "/user/new")
	// } else {
	// 	log.Debugf("cu: %#v", cu)
	// 	c.HTML(http.StatusOK, "welcome/index", gin.H{
	// 		"VersionID": appengine.VersionID(c),
	// 		"CUser":     cu,
	// 		"Context":   c})
	// }
}
