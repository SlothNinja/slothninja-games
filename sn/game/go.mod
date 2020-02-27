module github.com/SlothNinja/slothninja-games/sn/game

require (
	cloud.google.com/go/datastore v1.1.0
	github.com/SlothNinja/gt v0.0.0-20200211000447-e0d0da4579a3 // indirect
	github.com/SlothNinja/log v0.0.2
	github.com/SlothNinja/slothninja-games v0.0.0-20200217151912-ca6570104a61 // indirect
	github.com/SlothNinja/slothninja-games/sn/color v0.0.0-20200212041531-4ad6845c4545
	github.com/SlothNinja/slothninja-games/sn/log v0.0.0-20200213013905-9431a99919ae
	github.com/SlothNinja/slothninja-games/sn/misc v0.0.0-20200213013905-9431a99919ae
	github.com/SlothNinja/slothninja-games/sn/mlog v0.0.0-20200217151912-ca6570104a61 // indirect
	github.com/SlothNinja/slothninja-games/sn/rating v0.0.0-20200212041301-c78b1da4fcdb
	github.com/SlothNinja/slothninja-games/sn/restful v0.0.0-20200213013905-9431a99919ae
	github.com/SlothNinja/slothninja-games/sn/send v0.0.0-20200212040712-1ee26e37e1e4
	github.com/SlothNinja/slothninja-games/sn/type v0.0.0-20200212041301-c78b1da4fcdb
	github.com/SlothNinja/slothninja-games/sn/user v0.0.0-20200213013905-9431a99919ae
	github.com/SlothNinja/slothninja-games/sn/user/stats v0.0.0-20200217151912-ca6570104a61
	github.com/gin-gonic/gin v1.5.0
	go.chromium.org/gae v0.0.0-20190826183307-50a499513efa
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	google.golang.org/appengine v1.6.5
)

replace github.com/SlothNinja/slothninja-games/sn/user => ./user

replace github.com/SlothNinja/slothninja-games/sn/rating => ./rating

replace github.com/SlothNinja/slothninja-games/sn/contest => ./contest

replace github.com/SlothNinja/slothninja-games/sn/color => ./color

replace github.com/SlothNinja/slothninja-games/sn/send => ./send
