module github.com/SlothNinja/slothninja-games/welcome

require (
	github.com/SlothNinja/log v0.0.2
	github.com/SlothNinja/slothninja-games/sn/user v0.0.0-20200213013905-9431a99919ae
	github.com/gin-gonic/gin v1.5.0
	google.golang.org/appengine v1.6.5
)

replace github.com/SlothNinja/slothninja-games/sn/user => ./user
