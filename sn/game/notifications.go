package game

import (
	"bytes"

	"github.com/SlothNinja/slothninja-games/sn/log"
	"github.com/SlothNinja/slothninja-games/sn/restful"
	"github.com/SlothNinja/slothninja-games/sn/send"
	gType "github.com/SlothNinja/slothninja-games/sn/type"
	"github.com/SlothNinja/slothninja-games/sn/user"
	"github.com/gin-gonic/gin"
	"go.chromium.org/gae/service/datastore"
	"go.chromium.org/gae/service/mail"
)

const (
	sender  = "webmaster@slothninja.com"
	subject = "SlothNinja Games: Daily Turn Notifications"
)

type inf struct {
	GameID int64
	Type   gType.Type
	Title  string
}

type infs []*inf
type notifications map[int64]infs

func DailyNotifications(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	gs := GamersFrom(ctx)

	notifications := make(notifications, 0)
	for _, g := range gs {
		h := g.GetHeader()
		gameInfo := &inf{GameID: h.ID, Type: h.Type, Title: h.Title}
		for _, index := range h.CPUserIndices {
			uid := h.UserIDS[index]
			notifications[uid] = append(notifications[uid], gameInfo)
		}
	}

	msg := &mail.Message{Sender: sender, Subject: subject}
	tmpl := restful.TemplatesFrom(ctx)["shared/daily_notification"]
	buf := new(bytes.Buffer)

	var err error
	for uid, gameInfos := range notifications {
		u := user.New(ctx)
		u.ID = uid

		if err = datastore.Get(ctx, u); err != nil {
			log.Errorf(ctx, "get user error: %s", err.Error())
		}

		if err = tmpl.Execute(buf, gin.H{
			"Info": gameInfos,
			"User": u,
		}); err != nil {
			log.Errorf(ctx, "template execution for %s generated error: %s", u.Name, err.Error())
		}

		msg.To = []string{u.Email}
		msg.HTMLBody = buf.String()

		if err = send.Message(ctx, msg); err != nil {
			log.Errorf(ctx, "enqueuing email message: %#v geneerated error: %s", msg, err.Error())
		}

		// Reset buffer for next message
		buf.Reset()
	}
}

//func DailyNotifications(c *gin.Context) {
//	ctx := restful.ContextFrom(c)
//	log.Debugf(ctx, "Entering")
//	defer log.Debugf(ctx, "Exiting")
//
//	gs := GamersFrom(ctx)
//
//	notifications := make(notifications, 0)
//	for _, g := range gs {
//		h := g.GetHeader()
//		gameInfo := &inf{GameID: h.ID, Type: h.Type, Title: h.Title}
//		for _, index := range h.CPUserIndices {
//			uid := h.UserIDS[index]
//			notifications[uid] = append(notifications[uid], gameInfo)
//		}
//	}
//	for uid, gameInfos := range notifications {
//		u := user.New(ctx)
//		u.ID = uid
//		if err := datastore.Get(ctx, u); err != nil {
//			log.Errorf(ctx, "DailyNotifications Get User Error: %s", err)
//			return
//		}
//		msg := &mail.Message{
//			Sender:  "webmaster@slothninja.com",
//			To:      []string{u.Email},
//			Subject: "SlothNinja Games: Daily Turn Notifications",
//		}
//		msg.HTMLBody += `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
//		<html>
//			<head>
//				<meta http-equiv="content-type" content="text/html; charset=ISO-8859-1">
//			</head>
//			<body bgcolor="#ffffff" text="#000
//				<div>It is your turn in the following games:</div>>
//				<table>
//					<thead>
//						<tr><th>Game ID</th><th>Type</th><th>Title</th></tr>
//					</thead>
//				<tbody>`
//		for _, gameInfo := range gameInfos {
//			url := fmt.Sprintf(`<a href="http:/www.slothninja.com/%s/show/%d">%s</a>`,
//				gameInfo.Type.Prefix(), gameInfo.GameID, gameInfo.Title)
//			msg.HTMLBody += fmt.Sprintf("<tr><td>%v</td><td>%v</td><td>%v</td></tr>",
//				gameInfo.GameID, gameInfo.Type, url)
//		}
//		msg.HTMLBody += `</tbody>
//				<table>
//			</body>
//		</html>`
//		if err := send.Message(ctx, msg); err != nil {
//			log.Errorf(ctx, "Enqueuing email message: %#v Error: %v", msg, err)
//		}
//	}
//}

//func gamers(c *gin.Context) (gs []Gamer) {
//	if v, ok := c.Get(gamersKey); ok {
//		if gs, ok = v.([]Gamer); ok {
//			return gs
//		}
//	}
//	return nil
//}
