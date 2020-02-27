module github.com/SlothNinja/slothninja-games/sn/mlog

require (
	cloud.google.com/go v0.53.0
	cloud.google.com/go/datastore v1.1.0
	github.com/SlothNinja/log v0.0.2
	github.com/SlothNinja/slothninja-games/sn/codec v0.0.0-20200212041531-4ad6845c4545
	github.com/SlothNinja/slothninja-games/sn/color v0.0.0-20200212041531-4ad6845c4545
	github.com/SlothNinja/slothninja-games/sn/log v0.0.0-20200213013905-9431a99919ae
	github.com/SlothNinja/slothninja-games/sn/restful v0.0.0-20200213013905-9431a99919ae
	github.com/SlothNinja/slothninja-games/sn/user v0.0.0-20200213013905-9431a99919ae
	github.com/gin-gonic/gin v1.5.0
	github.com/googleapis/gax-go v2.0.2+incompatible // indirect
	go.chromium.org/gae v0.0.0-20190826183307-50a499513efa
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
)

replace github.com/SlothNinja/slothninja-games/sn/color => ./color

replace github.com/SlothNinja/slothninja-games/sn/restful => ./restful

replace github.com/SlothNinja/slothninja-games/sn/user => ./user
