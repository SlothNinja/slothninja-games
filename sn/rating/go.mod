module github.com/SlothNinja/slothninja-games/sn/rating

require (
	cloud.google.com/go/datastore v1.1.0
	github.com/SlothNinja/glicko v0.0.0-20200212023639-f1ccdd954e0a
	github.com/SlothNinja/log v0.0.2
	github.com/SlothNinja/slothninja-games/sn/contest v0.0.0-20200212040002-f57bad0ae251
	github.com/SlothNinja/slothninja-games/sn/restful v0.0.0-20200213013905-9431a99919ae
	github.com/SlothNinja/slothninja-games/sn/type v0.0.0-20200212040712-1ee26e37e1e4
	github.com/SlothNinja/slothninja-games/sn/user v0.0.0-20200212040712-1ee26e37e1e4
	github.com/gin-gonic/gin v1.5.0
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	google.golang.org/api v0.17.0
	google.golang.org/appengine v1.6.5
)

replace github.com/SlothNinja/slothninja-games/sn/user => ./user

replace github.com/SlothNinja/slothninja-games/sn/contest => ./contest
