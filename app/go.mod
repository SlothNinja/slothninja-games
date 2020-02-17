module github.com/SlothNinja/slothninja-games

require (
	github.com/SlothNinja/slothninja-games/sn/restful v0.0.0-20200213013905-9431a99919ae
	github.com/SlothNinja/slothninja-games/sn/user_controller v0.0.0-20200213013905-9431a99919ae
	github.com/SlothNinja/slothninja-games/welcome v0.0.0-20200213013905-9431a99919ae
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.5.0
)

replace github.com/SlothNinja/slothninja-games/sn/user_controller => ./user_controller

replace github.com/SlothNinja/slothninja-games/sn/user => ./user
