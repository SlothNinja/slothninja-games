package name

import (
	"errors"
	"strings"

	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"go.chromium.org/gae/service/datastore"
	"golang.org/x/net/context"
)

var ErrNameInUse = errors.New("Name already in use.")

type Name struct {
	//Key       *datastore.Key `datastore:"-"`
	ID        string `gae:"$id"`
	Kind      string `gae:"$kind"`
	GoogleID  string
	CreatedAt restful.CTime
	UpdatedAt restful.UTime
}

const kind = "UserName"

func New() *Name {
	return &Name{Kind: kind}
}

func ByName(ctx context.Context, name string, n *Name) error {
	n.ID = strings.ToLower(name)
	return datastore.Get(ctx, n)
}

//func NewKey(ctx *restful.Context, name string) *datastore.Key {
//	return datastore.NewKey(ctx, kind, name, 0, nil)
//}

func IsUnique(ctx context.Context, name string) bool {
	return ByName(ctx, name, New()) == datastore.ErrNoSuchEntity
}
