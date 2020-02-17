package user

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/SlothNinja/slothninja-games/sn/user"
	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/slothninja-games/sn/restful"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gofrs/uuid"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func init() {
	gob.Register(sessionToken{})
}

const (
	HOST        = "HOST"
	authPath    = "/auth"
	sessionKey  = "session"
	userNewPath = "/user/new"
)

type User struct {
	isAdmin bool
	Key     *datastore.Key `datastore:"__key__"`
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

// func (u *User) CTX() context.Context {
// 	return u.ctx
// }

func RootKey() *datastore.Key {
	return datastore.NameKey("Users", "root", nil)
}

func newKey(id int64) *datastore.Key {
	return datastore.IDKey(kind, id, RootKey())
}

func New(id int64) *User {
	return &User{Key: newKey(id)}
}

func (u User) ID() int64 {
	if u.Key != nil {
		return u.Key.ID
	}
	return 0
}

type NUser struct {
	Key   *datastore.Key `datastore:"__key__"`
	OldID int64          `json:"oldid"`
	Data
}

func (u *NUser) ID() string {
	if u.Key != nil {
		return u.Key.Name
	}
	return ""
}

func NNew(id string) *NUser {
	return &NUser{Key: datastore.NameKey(kind, id, RootKey())}
}

func ToNUser(u *User) *NUser {
	nu := NNew(GenID(u.GoogleID))
	nu.OldID = u.ID()
	nu.Data = u.Data
	return nu
}

func GenID(gid string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(salt+gid)))
}

func NewKeyFor(id int64) *datastore.Key {
	return New(id).Key
}

// func FromGUser(ctx context.Context, gu *user.User) *User {
// 	if gu == nil {
// 		return nil
// 	} else {
// 		n := strings.Split(gu.Email, "@")[0]
// 		u := New(ctx)
// 		u.Name = n
// 		u.LCName = strings.ToLower(n)
// 		u.Email = gu.Email
// 		u.GoogleID = gu.ID
// 		return u
// 	}
// }

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

func ByGoogleID(c *gin.Context, gid string) (*NUser, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	u := New(0)
	u.GoogleID = gid
	nu := ToNUser(u)
	log.Debugf("nu: %#v", nu)

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return nil, err
	}

	err = dsClient.Get(c, nu.Key, nu)
	if err == datastore.ErrNoSuchEntity {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Warningf(err.Error())
		return nil, err
	}

	return nu, nil
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

func AllQuery() *datastore.Query {
	return datastore.NewQuery(kind).Ancestor(RootKey())
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
	return ""
	// return info.VersionID(ctx) + gid
}

//func setByGoogleID(c *gin.Context, gid string, u *User) error {
//	if v, err := u.encode(); err != nil {
//		return err
//	} else {
//		ctx := restful.ContextFrom(c)
//		return memcache.Set(ctx, memcache.NewItem(ctx, MCKey(c, gid)).SetValue(v))
//	}
//}

// func IsAdmin(ctx context.Context) bool {
// 	return user.IsAdmin(ctx)
// }

func (u *User) IsAdmin() bool {
	return u != nil && u.isAdmin
}

func (u *User) IsAdminOrCurrent(ctx context.Context) bool {
	return u.IsAdmin() || u.IsCurrent(ctx)
}

func (u *User) Gravatar(options ...string) template.URL {
	return template.URL(GravatarURL(u.Email, options...))
}

func (u *NUser) Gravatar(options ...string) template.URL {
	return template.URL(GravatarURL(u.Email, options...))
}

func (u *OAuth) Gravatar(options ...string) template.URL {
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

func (u *User) Update(c *gin.Context) error {
	n := New(0)
	err := restful.BindWith(c, n, binding.FormPost)
	if err != nil {
		return err
	}

	log.Debugf("n: %#v", n)

	if u.IsAdmin() {
		if n.Email != "" {
			u.Email = n.Email
		}
	}

	if u.IsAdminOrCurrent(c) {
		if err := u.updateName(c, n.Name); err != nil {
			return err
		}
		u.EmailNotifications = n.EmailNotifications
	}

	return nil
}

func (u *User) updateName(c *gin.Context, n string) error {
	matcher := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9._%+\-]+$`)

	switch {
	case n == u.Name:
		return nil
	case len(n) > 15:
		return fmt.Errorf("%q is too long.", n)
	case !matcher.MatchString(n):
		return fmt.Errorf("%q is not a valid user name.", n)
	case !NameIsUnique(c, n):
		return fmt.Errorf("%q is not a unique user name.", n)
	}

	u.Name = n
	u.LCName = strings.ToLower(n)
	return nil
}

func NameIsUnique(c *gin.Context, name string) bool {
	LCName := strings.ToLower(name)

	q := datastore.NewQuery("User").Filter("LCName=", LCName)
	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return false
	}

	cnt, err := dsClient.Count(c, q)
	if err != nil {
		return false
	}
	return cnt == 0
}

func (u *User) Equal(u2 *User) bool {
	return u.ID() == u2.ID()
}

func (u *User) Link() template.HTML {
	if u == nil {
		return ""
	}
	return LinkFor(u.ID(), u.Name)
}

func LinkFor(uid int64, name string) template.HTML {
	return template.HTML(fmt.Sprintf("<a href=%q>%s</a>", PathFor(uid), name))
}

func PathFor(uid int64) template.HTML {
	return template.HTML(fmt.Sprintf("/user/show/%d", uid))
}

// func GetGUserHandler(c *gin.Context) {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Exiting")
//
// 	WithGUser(c, user.Current(ctx))
// }

// // Use after GetGUserHandler handler
// func GetCUserHandler(c *gin.Context) {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Exiting")
//
// 	var u *User
// 	if gu := GUserFrom(ctx); gu != nil {
// 		// Attempt to fetch and return stored User
// 		if nu, err := ByGoogleID(ctx, gu.ID); err != nil {
// 			log.Debugf(err.Error())
// 		} else {
// 			u = New(ctx)
// 			u.ID, u.Data, u.isAdmin = nu.OldID, nu.Data, user.IsAdmin(ctx)
// 		}
// 	}
// 	WithCurrent(c, u)
// }

func Login(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		session := sessions.Default(c)
		state := randToken()
		session.Set("state", state)
		session.Save()

		c.Redirect(http.StatusSeeOther, getLoginURL(c, path, state))
	}
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func getLoginURL(c *gin.Context, path, state string) string {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	// State can be some kind of random generated hash string.
	// See relevant RFC: http://tools.ietf.org/html/rfc6749#section-10.12
	return oauth2Config(c, path, scopes()...).AuthCodeURL(state)
}

func oauth2Config(c *gin.Context, path string, scopes ...string) *oauth2.Config {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	log.Debugf("request: %#v", c.Request)

	// protocol := "http"
	// if c.Request.TLS != nil {
	// 	protocol = "https"
	// }

	return &oauth2.Config{
		ClientID:     "435340145701-t5o50sjq7hsbilopgreobhvrv30e1tj4.apps.googleusercontent.com",
		ClientSecret: "Fe5f-Ht1V5_GohDEOS_TQOVc",
		Endpoint:     google.Endpoint,
		Scopes:       scopes,
		RedirectURL:  fmt.Sprintf("%s%s", getHost(), path),
	}
}

func scopes() []string {
	return []string{"email", "profile", "openid"}
}

func getHost() string {
	return os.Getenv(HOST)
}

type Info struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	LoggedIn      bool
	Admin         bool
}

const fqdn = "www.slothninja.com"

var namespaceUUID = uuid.NewV5(uuid.NamespaceDNS, fqdn)

// Generates ID for User from ID obtained from OAuth OpenID Connect
func id(s string) string {
	return uuid.NewV5(namespaceUUID, s).String()
}

type OAuth struct {
	isAdmin bool
	Key     *datastore.Key `datastore:"__key__"`
	OldID1  int64          `json:"oldID1"`
	OldID2  string         `json:"oldID2"`
	Data
}

func (u OAuth) ID() string {
	if u.Key != nil {
		return u.Key.Name
	}
	return ""
}

func (u *OAuth) Update(c *gin.Context) error {
	n := NewOAuth("")
	err := restful.BindWith(c, n, binding.FormPost)
	if err != nil {
		return err
	}

	log.Debugf("n: %#v", n)

	cu := user.CurrentFrom(c)
	if cu.IsAdmin() {
		if n.Email != "" {
			u.Email = n.Email
		}
	}

	if cu.IsAdminOrCurrent(c) {
		if err := u.updateName(c, n.Name); err != nil {
			return err
		}
		u.EmailNotifications = n.EmailNotifications
	}

	return nil
}

func (u *OAuth) updateName(c *gin.Context, n string) error {
	matcher := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9._%+\-]+$`)

	switch {
	case n == u.Name:
		return nil
	case len(n) > 15:
		return fmt.Errorf("%q is too long.", n)
	case !matcher.MatchString(n):
		return fmt.Errorf("%q is not a valid user name.", n)
	case !NameIsUnique(c, n):
		return fmt.Errorf("%q is not a unique user name.", n)
	}

	u.Name = n
	u.LCName = strings.ToLower(n)
	return nil
}

func NewKeyOAuth(id string) *datastore.Key {
	return datastore.NameKey(kind, id, RootKey())
}

func NewOAuth(id string) *OAuth {
	return &OAuth{Key: NewKeyOAuth(id)}
}

func Auth(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		// Handle the exchange code to initiate a transport.
		session := sessions.Default(c)
		retrievedState := session.Get("state")
		if retrievedState != c.Query("state") {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
			return
		}

		log.Debugf("retrievedState: %#v", retrievedState)
		// ac := appengine.NewContext(c.Request)
		ac := context.Background()
		conf := oauth2Config(c, path, scopes()...)
		log.Debugf("conf: %#v", conf)
		tok, err := conf.Exchange(ac, c.Query("code"))
		if err != nil {
			log.Errorf("tok error: %#v", err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		log.Debugf("tok: %#v", tok)

		client := conf.Client(ac, tok)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		log.Debugf("body: %s", body)

		uinfo := Info{}
		var b binding.BindingBody = binding.JSON
		err = b.BindBody(body, &uinfo)
		if err != nil {
			log.Errorf("BindBody error: %v", err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		log.Debugf("info: %#v", uinfo)

		u := NewOAuth(id(uinfo.Sub))
		dsClient, err := datastore.NewClient(c, "")
		if err != nil {
			log.Errorf("Client error: %v", err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err = dsClient.Get(c, u.Key, u)
		if err == datastore.ErrNoSuchEntity {
			u.Email = uinfo.Email
			err = u.SaveTo(session)
			if err != nil {
				log.Errorf("session.Save error: %v", err)
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			log.Debugf("session saved")
			c.Redirect(http.StatusSeeOther, userNewPath)
			return
		}

		if err != nil {
			log.Errorf("ByID => \n\t id: %s\n\t error: %s", id, err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err = u.SaveTo(session)
		if err != nil {
			log.Errorf("session.Save error: %v", err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		log.Debugf("session saved")

		c.Redirect(http.StatusSeeOther, "/")
	}
}

func (u OAuth) SaveTo(s sessions.Session) error {
	s.Set(sessionKey, sessionToken{ID: u.ID(), Email: u.Email})
	return s.Save()
}

type sessionToken struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func SessionTokenFrom(s sessions.Session) (sessionToken, bool) {
	token, ok := s.Get(sessionKey).(sessionToken)
	return token, ok
}

// // Use after GetGUserHandler and GetUserHandler handlers
// func RequireLogin() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if gu := GUserFrom(ctx); gu == nil {
// 			log.Warningf("RequireLogin failed.")
// 			c.Redirect(http.StatusSeeOther, "/")
// 			c.Abort()
// 		}
// 	}
// }
//
// func RequireCurrentUser() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx := restful.ContextFrom(c)
//
// 		if cu := CurrentFrom(ctx); cu == nil {
// 			log.Warningf("RequireCurrentUser failed.")
// 			c.Redirect(http.StatusSeeOther, "/")
// 			c.Abort()
// 		}
// 	}
// }

// func RequireAdmin(c *gin.Context) {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Exiting")
//
// 	if !user.IsAdmin(ctx) {
// 		log.Warningf("user not admin.")
// 		c.Redirect(http.StatusSeeOther, "/")
// 		c.Abort()
// 	}
// }

// func Fetch(c *gin.Context) {
// 	log.Debugf("Entering user#Fetch")
// 	defer log.Debugf("Exiting user#Fetch")
//
// 	uid, err := getUID(c)
// 	if err != nil || uid == NotFound {
// 		log.Errorf("getUID error: %v", err.Error())
// 		c.Redirect(http.StatusSeeOther, "/")
// 		c.Abort()
// 		return
// 	}
//
// 	u := New(ctx)
// 	u.ID = uid
// 	if err = datastore.Get(ctx, u); err != nil {
// 		log.Errorf("Unable to get user for id: %v", c.Param("uid"))
// 		c.Redirect(http.StatusSeeOther, "/")
// 		c.Abort()
// 	} else {
// 		WithUser(c, u)
// 	}
// }
//
// func FetchAll(c *gin.Context) {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Exiting")
//
// 	us, cnt, err := getFiltered(ctx, c.PostForm("start"), c.PostForm("length"))
//
// 	if err != nil {
// 		log.Errorf(err.Error())
// 		c.Redirect(http.StatusSeeOther, homePath)
// 		c.Abort()
// 	}
// 	withUsers(withCount(c, cnt), us)
// }

func getFiltered(c *gin.Context, start, length string) ([]interface{}, int, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return nil, 0, err
	}

	q := AllQuery().Order("GoogleID").KeysOnly()
	cnt, err := dsClient.Count(c, q)
	if err != nil {
		return nil, 0, err
	}

	if start != "" {
		st, err := strconv.ParseInt(start, 10, 32)
		if err == nil {
			q = q.Offset(int(st))
		}
	}

	if length != "" {
		l, err := strconv.ParseInt(length, 10, 32)
		if err == nil {
			q = q.Limit(int(l))
		}
	}

	var ks []*datastore.Key
	ks, err = dsClient.GetAll(c, q, nil)
	if err != nil {
		return nil, 0, err
	}

	l := len(ks)
	us := make([]interface{}, l)
	for i := range us {
		if id := ks[i].ID; id != 0 {
			u := New(id)
			us[i] = u
		} else {
			u := NNew(ks[i].Name)
			us[i] = u
		}
	}

	err = dsClient.GetMulti(c, ks, us)
	return us, cnt, err
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

// func GUserFrom(ctx context.Context) (u *user.User) {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Exiting")
//
// 	u, _ = ctx.Value(guserKey).(*user.User)
// 	log.Debugf("u: %#v", u)
// 	return
// }

func WithUser(c *gin.Context, u *User) {
	c.Set(userKey, u)
}

// func WithGUser(c *gin.Context, u *user.User) {
// 	c.Set(guserKey, u)
// }

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
