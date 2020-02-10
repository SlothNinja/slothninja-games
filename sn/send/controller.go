package send

import (
	"io/ioutil"
	"net/http"

	"bitbucket.org/SlothNinja/slothninja-games/sn/codec"
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"github.com/gin-gonic/gin"
	"go.chromium.org/gae/service/mail"
)

func Mail(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

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

	if err := mail.Send(ctx, m); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
