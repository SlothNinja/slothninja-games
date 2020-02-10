package user

import (
	"crypto/md5"
	"crypto/sha1"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.chromium.org/gae/service/datastore"
	"go.chromium.org/gae/service/info"
	"go.chromium.org/gae/service/user"
	"golang.org/x/net/context"
)

type User struct {
	ctx     context.Context
	isAdmin bool
	ID      int64          `gae:"$id"`
	Parent  *datastore.Key `gae:"$parent"`
	Data
}

type Data struct {
	Name               string        `json:"name" form:"name"`
	LCName             string        `json:"lcname"`
	Email              string        `json:"email" form:"email"`
	GoogleID           string        `json:"googleid"`
	XMPPNotifications  bool          `json:"xmppnotifications"`
	EmailNotifications bool          `json:"emailnotifications" form:"emailNotifications"`
	EmailReminders     bool          `json:"emailreminders"`
	Joined             time.Time     `json:"joined"`
	CreatedAt          restful.CTime `json:"createdat"`
	UpdatedAt          restful.UTime `json:"updatedat"`
}

const (
	kind             = "User"
	uidParam         = "uid"
	guserKey         = "guser"
	currentKey       = "current"
	userKey          = "User"
	countKey         = "count"
	homePath         = "/"
	salt             = "slothninja"
	usersKey         = "Users"
	NotFound   int64 = -1
)

type Users []*User

type UserName struct {
	GoogleID string
}

var (
	ErrNotFound     = errors.New("User not found.")
	ErrTooManyFound = errors.New("Found too many users.")
)

func (u *User) CTX() context.Context {
	return u.ctx
}

func RootKey(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, "Users", "root", 0, nil)
}

func NewKey(ctx context.Context, id int64) *datastore.Key {
	return datastore.NewKey(ctx, kind, "", id, RootKey(ctx))
}

func New(ctx context.Context) *User {
	return &User{ctx: ctx, Parent: RootKey(ctx)}
}

type NUser struct {
	ID     string         `gae:"$id"`
	Parent *datastore.Key `gae:"$parent"`
	Kind   string         `gae:"$kind"`
	OldID  int64          `json:"oldid"`
	Data
}

func NNew(ctx context.Context) *NUser {
	return &NUser{Parent: RootKey(ctx), Kind: kind}
}

func ToNUser(ctx context.Context, u *User) (nu *NUser) {
	nu = NNew(ctx)
	nu.ID = GenID(u.GoogleID)
	nu.OldID = u.ID
	nu.Data = u.Data
	return
}

func GenID(gid string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(salt+gid)))
}

func NewKeyFor(ctx context.Context, id int64) *datastore.Key {
	u := New(ctx)
	u.ID = id
	return datastore.KeyForObj(ctx, u)
}

func FromGUser(ctx context.Context, gu *user.User) *User {
	if gu == nil {
		return nil
	} else {
		n := strings.Split(gu.Email, "@")[0]
		u := New(ctx)
		u.Name = n
		u.LCName = strings.ToLower(n)
		u.Email = gu.Email
		u.GoogleID = gu.ID
		return u
	}
}

//func ByGoogleID(ctx context.Context, gid string) (*User, error) {
//	q := datastore.NewQuery(kind).Ancestor(RootKey(ctx)).Eq("GoogleID", gid).KeysOnly(true)
//
//	var keys []*datastore.Key
//	err := datastore.GetAll(ctx, q, &keys)
//	if err != nil {
//		return nil, err
//	}
//
//	u := New(ctx)
//	//var key *datastore.Key
//	switch l := len(keys); l {
//	case 0:
//		return nil, ErrNotFound
//	case 1:
//		datastore.PopulateKey(u, keys[0])
//		if err = datastore.Get(ctx, u); err != nil {
//			return nil, err
//		}
//	default:
//		return nil, ErrTooManyFound
//	}
//
//	return u, nil
//}

func ByGoogleID(ctx context.Context, gid string) (nu *NUser, err error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	u := New(ctx)
	u.GoogleID = gid
	nu = ToNUser(ctx, u)
	log.Debugf(ctx, "nu: %#v", nu)
	if err = datastore.Get(ctx, nu); err != nil {
		if err == datastore.ErrNoSuchEntity {
			err = ErrNotFound
		}
		log.Warningf(ctx, err.Error())
		return
	}

	return
}

//func GetMulti(ctx context.Context, ks []*datastore.Key) (Users, error) {
//	us := make([]*User, len(ks))
//	for i, k := range ks {
//		us[i] = new(User)
//		datastore.PopulateKey(us[i], k)
//	}
//	err := datastore.Get(ctx, us)
//	return us, err
//}
//
//func ByID(ctx context.Context, id int64) (u *User, err error) {
//	u = New(ctx)
//	u.ID = id
//	err = datastore.Get(ctx, u)
//	return
//}
//
//func BySID(ctx context.Context, sid string) (*User, error) {
//	id, err := strconv.ParseInt(sid, 10, 64)
//	if err != nil {
//		return nil, err
//	}
//	return ByID(ctx, id)
//}
//
//func ByIDS(ctx context.Context, ids []int64) (Users, error) {
//	ks := make([]*datastore.Key, len(ids))
//	for i, id := range ids {
//		ks[i] = NewKey(ctx, id)
//	}
//	return GetMulti(ctx, ks)
//}

func AllQuery(ctx context.Context) *datastore.Query {
	return datastore.NewQuery(kind).Ancestor(RootKey(ctx))
}

//func getByGoogleID(ctx context.Context, gid string) (*User, error) {
//	itm := memcache.NewItem(ctx, MCKey(ctx, gid))
//	if err := memcache.Get(ctx, itm); err != nil {
//		return nil, err
//	}
//
//	u := New(ctx)
//	if err := decode(u, itm.Value()); err != nil {
//		return nil, err
//	}
//	return u, nil
//}

//func (u User) encode() ([]byte, error) {
//	return codec.Encode(u)
//}
//
//func decode(dst *User, v []byte) error {
//	return codec.Decode(dst, v)
//}
//
func MCKey(ctx context.Context, gid string) string {
	return info.VersionID(ctx) + gid
}

//func setByGoogleID(c *gin.Context, gid string, u *User) error {
//	if v, err := u.encode(); err != nil {
//		return err
//	} else {
//		ctx := restful.ContextFrom(c)
//		return memcache.Set(ctx, memcache.NewItem(ctx, MCKey(c, gid)).SetValue(v))
//	}
//}

func IsAdmin(ctx context.Context) bool {
	return user.IsAdmin(ctx)
}

func (u *User) IsAdmin() bool {
	return u != nil && u.isAdmin
}

func (u *User) IsAdminOrCurrent(ctx context.Context) bool {
	return IsAdmin(ctx) || u.IsCurrent(ctx)
}

func (u *User) Gravatar(options ...string) template.URL {
	return template.URL(GravatarURL(u.Email, options...))
}

func (u *NUser) Gravatar(options ...string) template.URL {
	return template.URL(GravatarURL(u.Email, options...))
}

func GravatarURL(email string, options ...string) string {
	size := "80"
	if len(options) == 1 {
		size = options[0]
	}

	email = strings.ToLower(strings.TrimSpace(email))
	hash := md5.New()
	hash.Write([]byte(email))
	md5string := fmt.Sprintf("%x", hash.Sum(nil))
	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?s=%s&d=monsterid", md5string, size)
}

func (u *User) Update(ctx context.Context) error {
	n := New(ctx)
	if err := restful.BindWith(ctx, n, binding.FormPost); err != nil {
		return err
	}

	log.Debugf(ctx, "n: %#v", n)

	if IsAdmin(ctx) {
		if n.Email != "" {
			u.Email = n.Email
		}
	}

	if u.IsAdminOrCurrent(ctx) {
		if err := u.updateName(ctx, n.Name); err != nil {
			return err
		}
		u.EmailNotifications = n.EmailNotifications
	}

	return nil
}

func (u *User) updateName(ctx context.Context, n string) error {
	matcher := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9._%+\-]+$`)

	switch {
	case n == u.Name:
		return nil
	case len(n) > 15:
		return fmt.Errorf("%q is too long.", n)
	case !matcher.MatchString(n):
		return fmt.Errorf("%q is not a valid user name.", n)
	case !NameIsUnique(ctx, n):
		return fmt.Errorf("%q is not a unique user name.", n)
	}

	u.Name = n
	u.LCName = strings.ToLower(n)
	return nil
}

func NameIsUnique(ctx context.Context, name string) bool {
	LCName := strings.ToLower(name)

	q := datastore.NewQuery("User").Eq("LCName", LCName)
	if cnt, err := datastore.Count(ctx, q); err != nil {
		return false
	} else {
		return cnt == 0
	}
}

func (u *User) Equal(u2 *User) bool {
	return u2 != nil && u != nil && u.ID == u2.ID
}

func (u *User) Link() template.HTML {
	if u == nil {
		return ""
	}
	return LinkFor(u.ID, u.Name)
}

func LinkFor(uid int64, name string) template.HTML {
	return template.HTML(fmt.Sprintf("<a href=%q>%s</a>", PathFor(uid), name))
}

func PathFor(uid int64) template.HTML {
	return template.HTML(fmt.Sprintf("/user/show/%d", uid))
}

func GetGUserHandler(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	WithGUser(c, user.Current(ctx))
}

// Use after GetGUserHandler handler
func GetCUserHandler(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	var u *User
	if gu := GUserFrom(ctx); gu != nil {
		// Attempt to fetch and return stored User
		if nu, err := ByGoogleID(ctx, gu.ID); err != nil {
			log.Debugf(ctx, err.Error())
		} else {
			u = New(ctx)
			u.ID, u.Data, u.isAdmin = nu.OldID, nu.Data, user.IsAdmin(ctx)
		}
	}
	WithCurrent(c, u)
}

// Use after GetGUserHandler and GetUserHandler handlers
func RequireLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.ContextFrom(c)

		if gu := GUserFrom(ctx); gu == nil {
			log.Warningf(ctx, "RequireLogin failed.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
		}
	}
}

func RequireCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.ContextFrom(c)

		if cu := CurrentFrom(ctx); cu == nil {
			log.Warningf(ctx, "RequireCurrentUser failed.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
		}
	}
}

func RequireAdmin(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	if !user.IsAdmin(ctx) {
		log.Warningf(ctx, "user not admin.")
		c.Redirect(http.StatusSeeOther, "/")
		c.Abort()
	}
}

func Fetch(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering user#Fetch")
	defer log.Debugf(ctx, "Exiting user#Fetch")

	uid, err := getUID(c)
	if err != nil || uid == NotFound {
		log.Errorf(ctx, "getUID error: %v", err.Error())
		c.Redirect(http.StatusSeeOther, "/")
		c.Abort()
		return
	}

	u := New(ctx)
	u.ID = uid
	if err = datastore.Get(ctx, u); err != nil {
		log.Errorf(ctx, "Unable to get user for id: %v", c.Param("uid"))
		c.Redirect(http.StatusSeeOther, "/")
		c.Abort()
	} else {
		WithUser(c, u)
	}
}

func FetchAll(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	us, cnt, err := getFiltered(ctx, c.PostForm("start"), c.PostForm("length"))

	if err != nil {
		log.Errorf(ctx, err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		c.Abort()
	}
	withUsers(withCount(c, cnt), us)
}

func getFiltered(ctx context.Context, start, length string) (us []interface{}, cnt int64, err error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	q := AllQuery(ctx).Order("GoogleID").KeysOnly(true)
	if cnt, err = datastore.Count(ctx, q); err != nil {
		return
	}

	if start != "" {
		if st, err := strconv.ParseInt(start, 10, 32); err == nil {
			q = q.Offset(int32(st))
		}
	}

	if length != "" {
		if l, err := strconv.ParseInt(length, 10, 32); err == nil {
			q = q.Limit(int32(l))
		}
	}

	var ks []*datastore.Key
	if err = datastore.GetAll(ctx, q, &ks); err != nil {
		return
	}

	l := len(ks)
	us = make([]interface{}, l)
	for i := range us {
		if id := ks[i].IntID(); id != 0 {
			u := New(ctx)
			u.ID = id
			us[i] = u
		} else {
			u := NNew(ctx)
			u.ID = ks[i].StringID()
			us[i] = u
		}
	}

	err = datastore.Get(ctx, us)
	return
}

//func FetchAll(c *gin.Context) {
//	ctx := restful.ContextFrom(c)
//	log.Debugf(ctx, "Entering user#Fetch")
//	defer log.Debugf(ctx, "Exiting user#Fetch")
//
//	var (
//		ks  []*datastore.Key
//		err error
//	)
//
//	if err = datastore.GetAll(ctx, q, &ks); err != nil {
//		log.Errorf(ctx, "unable to get users, error %v", err.Error())
//		c.Redirect(http.StatusSeeOther, "/")
//		c.Abort()
//	}
//
//	u := New(ctx)
//	u.ID = uid
//	if err = datastore.Get(ctx, u); err != nil {
//		log.Errorf(ctx, "Unable to get user for id: %v", c.Param("uid"))
//		c.Redirect(http.StatusSeeOther, "/")
//		c.Abort()
//	} else {
//		WithUser(c, u)
//	}
//}

func getUID(c *gin.Context) (id int64, err error) {
	if id, err = strconv.ParseInt(c.Param(uidParam), 10, 64); err != nil {
		id = NotFound
	}
	return
}

func Fetched(ctx context.Context) *User {
	return From(ctx)
}

func Gravatar(u *User) template.HTML {
	return template.HTML(fmt.Sprintf(`<a href="/user/show/%d"><img src=%q alt="Gravatar" class="black-border" /></a>`, u.ID, u.Gravatar()))
}

func NGravatar(nu *NUser) template.HTML {
	return template.HTML(fmt.Sprintf(`<a href="/user/show/%s"><img src=%q alt="Gravatar" class="black-border" /></a>`, nu.ID, nu.Gravatar()))
}

func from(ctx context.Context, key string) (u *User) {
	u, _ = ctx.Value(key).(*User)
	return
}

func From(ctx context.Context) *User {
	return from(ctx, userKey)
}

func CurrentFrom(ctx context.Context) *User {
	return from(ctx, currentKey)
}

func (u *User) IsCurrent(ctx context.Context) bool {
	return u.Equal(CurrentFrom(ctx))
}

func GUserFrom(ctx context.Context) (u *user.User) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	u, _ = ctx.Value(guserKey).(*user.User)
	log.Debugf(ctx, "u: %#v", u)
	return
}

func WithUser(c *gin.Context, u *User) {
	c.Set(userKey, u)
}

func WithGUser(c *gin.Context, u *user.User) {
	c.Set(guserKey, u)
}

func WithCurrent(c *gin.Context, u *User) {
	c.Set(currentKey, u)
}

func UsersFrom(ctx context.Context) (us []interface{}) {
	us, _ = ctx.Value(usersKey).([]interface{})
	return
}

func withUsers(c *gin.Context, us []interface{}) {
	c.Set(usersKey, us)
}

func withCount(c *gin.Context, cnt int64) *gin.Context {
	c.Set(countKey, cnt)
	return c
}

func CountFrom(ctx context.Context) (cnt int64) {
	cnt, _ = ctx.Value(countKey).(int64)
	return
}
