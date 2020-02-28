module github.com/SlothNinja/slothninja-games/welcome

require (
	github.com/SlothNinja/log v0.0.2
	github.com/SlothNinja/slothninja-games/sn/user v0.0.0-20200227212226-2cd0a8e24f94
	github.com/gin-gonic/gin v1.5.0
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	google.golang.org/appengine v1.6.5
)

replace github.com/SlothNinja/slothninja-games/sn/user => ./user
