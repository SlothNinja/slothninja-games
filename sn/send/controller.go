package send

import (
	"io/ioutil"
	"net/http"

	"github.com/SlothNinja/log"
	"github.com/SlothNinja/slothninja-games/sn/codec"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine/mail"
)

func Mail(c *gin.Context) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	encoded, err := ioutil.ReadAll(c.Request.Body)
	c.Request.Body.Close()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	m := new(mail.Message)
	if err := codec.Decode(m, encoded); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := mail.Send(c, m); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
