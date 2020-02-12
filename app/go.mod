module github.com/SlothNinja/slothninja-games

require (
	github.com/SlothNinja/slothninja-games/sn/restful v0.0.0-20200211052815-753b8615b1ed
	github.com/SlothNinja/slothninja-games/welcome v0.0.0-20200211045043-2e3fbf7658c5
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.5.0
	golang.org/x/sys v0.0.0-20191120155948-bd437916bb0e // indirect
	google.golang.org/appengine v1.6.5
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace github.com/SlothNinja/slothninja-games/welcome => ./welcome

replace github.com/SlothNinja/slothninja-games/sn/restful => ./sn/restful
