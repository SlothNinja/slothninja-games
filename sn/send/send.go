package send

import (
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/slothninja-games/sn/codec"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine/mail"
	"google.golang.org/appengine/taskqueue"
)

//func init() {
//	gob.Register(new(mail.Message))
//	//gob.Register(new(xmpp.Message))
//}

//var qXMPP = delay.Func("xmpp", func(ctx context.Context, m *xmpp.Message) error {
//	return m.Send(ctx)
//})

//var qInvite = delay.Func("xmpp-invite", func(ctx context.Context, jid string) error {
//	return xmpp.Invite(ctx, jid, "")
//})

//var qMail = delay.Func("mail", func(ctx context.Context, m *mail.Message) error {
//	return mail.Send(ctx, m)
//})

//func XMPP(ctx context.Context, ms ...*xmpp.Message) error {
//	l := len(ms)
//	if l < 1 {
//		return fmt.Errorf("You must provide at least one message to send.")
//	}
//
//	args := make([]interface{}, len(ms))
//	for i, m := range ms {
//		args[i] = m
//	}
//	return enqueue(ctx, qXMPP, "xmpp", args)
//}
//
//func Invite(ctx context.Context, jids ...string) error {
//	l := len(jids)
//	if l < 1 {
//		return fmt.Errorf("You must provide at least one jid to invite.")
//	}
//	args := make([]interface{}, len(jids))
//	for i, jid := range jids {
//		args[i] = jid
//	}
//	return enqueue(ctx, qInvite, "xmpp", args)
//}

func Message(c *gin.Context, ms ...*mail.Message) error {
	log.Debugf("Entering")
	defer log.Debugf("Entering")

	ts := make([]*taskqueue.Task, len(ms))
	for i := range ms {
		encoded, err := codec.Encode(ms[i])
		if err != nil {
			return err
		}

		ts[i] = new(taskqueue.Task)
		ts[i].Path = "/mail"
		ts[i].Payload = encoded
	}

	ts, err := taskqueue.AddMulti(c, ts, "mail")
	return err
}

//func Mail(ctx context.Context, ms ...*mail.Message) error {
//	l := len(ms)
//	if l < 1 {
//		return fmt.Errorf("You must provide at least one message to mail.")
//	}
//	args := make([]interface{}, len(ms))
//	for i, m := range ms {
//		args[i] = m
//	}
//	return enqueue(ctx, qMail, "mail", args)
//}
//
//func enqueue(ctx context.Context, f *delay.Function, q string, args []interface{}) error {
//	l := len(args)
//	isNil := true
//	me := make(errors.MultiError, l)
//	ts := make([]*taskqueue.Task, l)
//
//	for i, arg := range args {
//		var err error
//		if ts[i], err = f.Task(arg); err != nil {
//			me[i] = err
//			isNil = false
//		}
//	}
//
//	if !isNil && l == 1 {
//		return me[0]
//	}
//
//	if !isNil {
//		return me
//	}
//
//	err := taskqueue.Add(ctx, q, ts...)
//	if me, ok := err.(errors.MultiError); ok {
//		if len(me) == 1 {
//			return me[0]
//		}
//		return me
//	}
//	return err
//}
