package user

import (
	"appengine/aetest"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"bytes"
	"encoding/gob"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var ctx *restful.Context
var req *http.Request

func getContext() *restful.Context {
	return ctx
}

func TestUser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Suite")
}

var _ = BeforeSuite(func() {
        aectx, err := aetest.NewContext(nil)
	Ω(err).ShouldNot(HaveOccurred())
	Ω(aectx).ShouldNot(BeNil())
        ctx = &restful.Context{Context: aectx}

	req, _ = http.NewRequest("GET", "http://localhost", nil)
	req.Header.Set("App-Testing", "1")
})

var _ = AfterSuite(func() {
	ctx.Context.(aetest.Context).Close()
})

func encode(src interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(src); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decode(dst interface{}, value []byte) error {
	if len(value) > 0 {
		buf := bytes.NewBuffer(value)
		dec := gob.NewDecoder(buf)
		return dec.Decode(dst)
	}
	return nil
}
