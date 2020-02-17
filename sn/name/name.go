package name

import (
	"errors"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/slothninja-games/sn/restful"
	"github.com/gin-gonic/gin"
)

var ErrNameInUse = errors.New("Name already in use.")

type Name struct {
	Key       *datastore.Key `datastore:"__key__"`
	GoogleID  string
	CreatedAt restful.CTime
	UpdatedAt restful.UTime
}

const kind = "UserName"

func New(id string) *Name {
	return &Name{Key: datastore.NameKey(kind, id, nil)}
}

func ByName(c *gin.Context, n *Name) error {
	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return err
	}
	return dsClient.Get(c, n.Key, &n)
}

//func NewKey(ctx *restful.Context, name string) *datastore.Key {
//	return datastore.NewKey(ctx, kind, name, 0, nil)
//}

func IsUnique(c *gin.Context, name string) bool {
	return ByName(c, New(name)) == datastore.ErrNoSuchEntity
}
