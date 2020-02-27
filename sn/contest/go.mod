module github.com/SlothNinja/slothninja-games/sn/contest

require (
	cloud.google.com/go v0.38.0
	github.com/SlothNinja/slothninja-games/sn/restful v0.0.0-20200212035615-29d0ccf53fcb
	github.com/SlothNinja/slothninja-games/sn/type v0.0.0-20200212035615-29d0ccf53fcb
	github.com/gin-gonic/gin v1.5.0
	github.com/googleapis/gax-go v2.0.2+incompatible // indirect
	go.chromium.org/gae v0.0.0-20190826183307-50a499513efa
	go.opencensus.io v0.22.3 // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	google.golang.org/api v0.17.0 // indirect
)

replace github.com/SlothNinja/slothninja-games/sn/type => ./type
