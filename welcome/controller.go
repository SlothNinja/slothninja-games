package welcome

import (
	"net/http"

	"github.com/SlothNinja/log"
	"github.com/SlothNinja/slothninja-games/sn/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
)

func Index(c *gin.Context) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	cu, found := user.Current(c)
	c.HTML(http.StatusOK, "welcome/index", gin.H{
		"VersionID": appengine.VersionID(c),
		"CUser":     cu,
		"CUFound":   found,
		"Context":   c})
}
