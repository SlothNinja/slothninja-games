module github.com/SlothNinja/slothninja-games

require (
	github.com/SlothNinja/gt v0.0.0-20200211000447-e0d0da4579a3
	github.com/SlothNinja/slothninja-games/sn/game v0.0.0-20200217151912-ca6570104a61
	github.com/SlothNinja/slothninja-games/sn/restful v0.0.0-20200217151912-ca6570104a61
	github.com/SlothNinja/slothninja-games/sn/schema v0.0.0-20200224004439-b903a4bc4557 // indirect
	github.com/SlothNinja/slothninja-games/sn/type v0.0.0-20200217151912-ca6570104a61
	github.com/SlothNinja/slothninja-games/sn/user_controller v0.0.0-20200217151912-ca6570104a61
	github.com/SlothNinja/slothninja-games/welcome v0.0.0-20200217151912-ca6570104a61
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.5.0
	github.com/gorilla/securecookie v1.1.1
)

replace github.com/SlothNinja/slothninja-games/sn/user_controller => ./user_controller

replace github.com/SlothNinja/slothninja-games/sn/user => ./user

replace github.com/SlothNinja/slothninja-games/sn/restful => ./restful

replace github.com/SlothNinja/slothninja-games/welcome => ./welcome

replace github.com/SlothNinja/slothninja-games/sn/game => ./game

replace github.com/SlothNinja/slothninja-games/sn/rating => ./rating

replace github.com/SlothNinja/slothninja-games/sn/contest => ./contest

replace github.com/SlothNinja/slothninja-games/sn/color => ./color

replace github.com/SlothNinja/slothninja-games/sn/send => ./send

replace github.com/SlothNinja/slothninja-games/sn/mlog => ./mlog

replace github.com/SlothNinja/slothninja-games/sn/type => ./type

replace github.com/SlothNinja/gt => ./got
