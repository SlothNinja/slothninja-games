module github.com/SlothNinja/slothninja-games/sn/user_controller

require (
	bitbucket.org/SlothNinja/glicko v0.0.0-20130425132718-9615bf559204 // indirect
	cloud.google.com/go/datastore v1.1.0
	github.com/SlothNinja/log v0.0.2
	github.com/SlothNinja/slothninja-games/sn/game v0.0.0-20200213013905-9431a99919ae
	github.com/SlothNinja/slothninja-games/sn/log v0.0.0-20200213013905-9431a99919ae
	github.com/SlothNinja/slothninja-games/sn/misc v0.0.0-20200213013905-9431a99919ae
	github.com/SlothNinja/slothninja-games/sn/name v0.0.0-20200213013905-9431a99919ae
	github.com/SlothNinja/slothninja-games/sn/rating v0.0.0-20200212041301-c78b1da4fcdb // indirect
	github.com/SlothNinja/slothninja-games/sn/restful v0.0.0-20200213013905-9431a99919ae
	github.com/SlothNinja/slothninja-games/sn/type v0.0.0-20200212041301-c78b1da4fcdb
	github.com/SlothNinja/slothninja-games/sn/user v0.0.0-20200212041301-c78b1da4fcdb
	github.com/SlothNinja/slothninja-games/sn/user/stats v0.0.0-20200212041301-c78b1da4fcdb
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.5.0
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	google.golang.org/appengine v1.6.5
)

replace github.com/SlothNinja/slothninja-games/sn/user => ./user

replace github.com/SlothNinja/slothninja-games/sn/rating => ./rating

replace github.com/SlothNinja/slothninja-games/sn/contest => ./contest

replace github.com/SlothNinja/slothninja-games/sn/name => ./name
