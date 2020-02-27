package user

import (
	"crypto/md5"
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

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gofrs/uuid"
	"github.com/gorilla/securecookie"
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
	Key *datastore.Key `datastore:"__key__"`
	Data
}

type Data struct {
	Name               string    `json:"name" form:"name"`
	LCName             string    `json:"lcname"`
	Email              string    `json:"email" form:"email"`
	GoogleID           string    `json:"googleid"`
	XMPPNotifications  bool      `json:"xmppnotifications"`
	EmailNotifications bool      `json:"emailnotifications" form:"emailNotifications"`
	EmailReminders     bool      `json:"emailreminders"`
	Joined             time.Time `json:"joined"`
	CreatedAt          time.Time `json:"createdat"`
	UpdatedAt          time.Time `json:"updatedat"`
	IsAdmin            bool      `json:"isAdmin"`
	Loaded             bool      `json:"loaded" datastore:"-"`
}

const (
	uKind            = "User"
	oauthsKind       = "OAuths"
	oauthKind        = "OAuth"
	root             = "root"
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
	return datastore.IDKey(uKind, id, RootKey())
}

func New(id int64) User {
	return User{Key: newKey(id)}
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
	return &NUser{Key: datastore.NameKey(uKind, id, RootKey())}
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

// func ByGoogleID(c *gin.Context, gid string) (*NUser, error) {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Exiting")
//
// 	u := New(0)
// 	u.GoogleID = gid
// 	nu := ToNUser(u)
// 	log.Debugf("nu: %#v", nu)
//
// 	dsClient, err := datastore.NewClient(c, "")
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	err = dsClient.Get(c, nu.Key, nu)
// 	if err == datastore.ErrNoSuchEntity {
// 		return nil, ErrNotFound
// 	}
//
// 	if err != nil {
// 		log.Warningf(err.Error())
// 		return nil, err
// 	}
//
// 	return nu, nil
// }

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
	return datastore.NewQuery(uKind).Ancestor(RootKey())
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
// func MCKey(ctx context.Context, gid string) string {
// 	return ""
// 	// return info.VersionID(ctx) + gid
// }

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

// func (u *User) IsAdmin() bool {
// 	return u != nil && u.isAdmin
// }

func (u User) IsAdminOrCurrent(c *gin.Context) bool {
	return u.IsAdmin || u.IsCurrent(c)
}

func (u User) Gravatar(options ...string) template.URL {
	return template.URL(GravatarURL(u.Email, options...))
}

func (u *NUser) Gravatar(options ...string) template.URL {
	return template.URL(GravatarURL(u.Email, options...))
}

// func (u OAuth) Gravatar(options ...string) template.URL {
// 	return template.URL(GravatarURL(u.Email, options...))
// }

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

func (u User) Update(c *gin.Context) (User, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	n := New(0)
	err := c.ShouldBind(n)
	if err != nil {
		return u, err
	}

	log.Debugf("n: %#v", n)

	if u.IsAdmin {
		if n.Email != "" {
			u.Email = n.Email
		}
	}

	if u.IsAdminOrCurrent(c) {
		err = u.updateName(c, n.Name)
		if err != nil {
			return u, err
		}
		u.EmailNotifications = n.EmailNotifications
	}

	return u, nil
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

func (u User) Equal(u2 User) bool {
	return u.ID() == u2.ID()
}

func (u User) Link() template.HTML {
	return LinkFor(u.ID(), u.Name)
}

func LinkFor(uid int64, name string) template.HTML {
	return template.HTML(fmt.Sprintf("<a href=%q>%s</a>", PathFor(uid), name))
}

func PathFor(uid int64) template.HTML {
	return template.HTML(fmt.Sprintf("/user/show/%d", uid))
}

// func For(c *gin.Context, oldID1 int64, oldID2, id string) (OAuth, error) {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Exiting")
//
// 	oa := NewOAuth(id)
//
// 	dsClient, err := datastore.NewClient(c, "")
// 	if err != nil {
// 		return oa, err
// 	}
//
// 	if id != "" {
// 		err = dsClient.Get(c, oa.Key, &oa)
// 		if err == nil {
// 			return oa, nil
// 		}
// 	}
//
// 	if oldID2 != "" {
// 		nu := NNew(oldID2)
// 		err = dsClient.Get(c, nu.Key, nu)
// 		if err == nil {
// 			return OAuth{
// 				Data:   nu.Data,
// 				OldID1: nu.OldID,
// 				OldID2: oldID2,
// 			}, nil
// 		}
// 	}
//
// 	if oldID1 != 0 {
// 		u := New(oldID1)
// 		err = dsClient.Get(c, u.Key, u)
// 		if err == nil {
// 			return OAuth{
// 				Data:   u.Data,
// 				OldID1: oldID1,
// 			}, nil
// 		}
// 	}
//
// 	return oa, fmt.Errorf("user not found")
// }

// func (u OAuth) Link() template.HTML {
// 	if u == None {
// 		return ""
// 	}
// 	return LinkFor(u.ID(), u.Name)
// }
//
// func LinkFor(uid, name string) template.HTML {
// 	return template.HTML(fmt.Sprintf("<a href=%q>%s</a>", PathFor(uid), name))
// }
//
// func PathFor(uid string) template.HTML {
// 	return template.HTML("/user/show/" + uid)
// }

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

const tokenLength = 32

func Login(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		session := sessions.Default(c)
		state := randToken(tokenLength)
		session.Set("state", state)
		session.Save()

		c.Redirect(http.StatusSeeOther, getLoginURL(c, path, state))
	}
}

func Logout(c *gin.Context) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	s := sessions.Default(c)
	s.Delete(sessionKey)
	err := s.Save()
	if err != nil {
		log.Warningf("unable to save session: %v", err)
	}
	c.Redirect(http.StatusSeeOther, homePath)
}

func randToken(length int) string {
	key := securecookie.GenerateRandomKey(length)
	// b := make([]byte, 32)
	// rand.Read(b)
	return base64.StdEncoding.EncodeToString(key)
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
func GenOAuthID(s string) string {
	return uuid.NewV5(namespaceUUID, s).String()
}

type OAuth struct {
	Key       *datastore.Key `datastore:"__key__"`
	ID        int64
	UpdatedAt time.Time
}

// func (u OAuth) FromUser(old *User) OAuth {
// 	u.OldID1 = old.Key.ID
// 	u.Data = old.Data
// 	u.UpdatedAt = time.Now()
// 	return u
// }
//
// func (u OAuth) FromNUser(old *NUser) OAuth {
// 	u.OldID1 = old.OldID
// 	u.OldID2 = old.Key.Name
// 	u.Data = old.Data
// 	u.UpdatedAt = time.Now()
// 	return u
// }
//
// func (u OAuth) ID() string {
// 	if u.Key != nil {
// 		return u.Key.Name
// 	}
// 	return ""
// }

// func (u *OAuth) Update(c *gin.Context) error {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Exiting")
//
// 	n := struct {
// 		Name               string `form:"name"`
// 		Email              string `form:"email"`
// 		EmailNotifications bool   `form:"emailNotifications"`
// 	}{}
//
// 	err := c.ShouldBind(&n)
// 	if err != nil {
// 		return err
// 	}
//
// 	cu, found := Current(c)
// 	if !found {
// 		return errors.New("current user not found.")
// 	}
// 	if cu.IsAdmin {
// 		if n.Email != "" {
// 			u.Email = n.Email
// 		}
// 		// err = u.updateName(c, n.Name)
// 		// if err != nil {
// 		// 	return err
// 		// }
// 	}
//
// 	if cu.IsAdmin || (cu.ID() == u.ID()) {
// 		u.EmailNotifications = n.EmailNotifications
// 	}
//
// 	return nil
// }
//
// func (u *OAuth) updateName(c *gin.Context, n string) error {
// 	matcher := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9._%+\-]+$`)
//
// 	switch {
// 	case n == u.Name:
// 		return nil
// 	case len(n) > 15:
// 		return fmt.Errorf("%q is too long.", n)
// 	case !matcher.MatchString(n):
// 		return fmt.Errorf("%q is not a valid user name.", n)
// 		//	case !NameIsUnique(c, n):
// 		//		return fmt.Errorf("%q is not a unique user name.", n)
// 	}
//
// 	u.Name = n
// 	u.LCName = strings.ToLower(n)
// 	return nil
// }

func pk() *datastore.Key {
	return datastore.NameKey(oauthsKind, root, nil)
}

// func tmpPK() *datastore.Key {
// 	return datastore.NameKey(oauthKind, tmpRoot, nil)
// }

// func NewTmpKey(email string) *datastore.Key {
// 	return datastore.NameKey(oauthKind, email, tmpPK())
// }

func NewKeyOAuth(id string) *datastore.Key {
	return datastore.NameKey(oauthKind, id, pk())
}

func NewOAuth(id string) OAuth {
	return OAuth{Key: NewKeyOAuth(id)}
}

func ByEmail(c *gin.Context, email string) (OAuth, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return OAuth{}, err
	}

	q := datastore.NewQuery(oauthKind).
		Ancestor(pk()).
		Filter("Equal=", email)

	var oas []OAuth
	_, err = dsClient.GetAll(c, q, oas)
	if err != nil {
		return OAuth{}, err
	}
	l := len(oas)
	if l != 1 {
		return OAuth{}, fmt.Errorf("found %d, expect 1", l)
	}
	return oas[0], nil
}

func Auth(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		uInfo, err := getUInfo(c, path)
		if err != nil {
			log.Errorf(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		oaid := GenOAuthID(uInfo.Sub)
		oa, err := getOAuth(c, oaid)
		// Succesfully pulled oauth id from datastore

		u := New(oa.ID)
		if err == nil {
			dsClient, err := datastore.NewClient(c, "")
			if err != nil {
				log.Errorf("unable to connect to datastore")
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			err = dsClient.Get(c, u.Key, &u)
			if err != nil {
				log.Errorf(err.Error())
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			u.Loaded = true
			saveToSessionAndReturnTo(c, u, homePath)
			return
		}

		// Datastore error other than missing entity.
		if err != datastore.ErrNoSuchEntity {
			log.Errorf("unable to get user for key: %#v", u.Key)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// oauth id not present in datastore
		// Check to see if other entities exist for same email address.
		// If so, use old entities for user
		u, err = getByEmail(c, uInfo.Email)
		// if err != nil {
		// 	log.Errorf(err.Error())
		// 	c.AbortWithStatus(http.StatusBadRequest)
		// 	return
		// }

		log.Debugf("getByEmail => u: %v\nerr: %v", u, err)

		// u, err = u.migrateOld(c, ks)
		if err == nil {
			dsClient, err := datastore.NewClient(c, "")
			if err != nil {
				log.Errorf(err.Error())
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			oa := NewOAuth(oaid)
			oa.ID = u.ID()
			oa.UpdatedAt = time.Now()
			_, err = dsClient.Put(c, oa.Key, &oa)
			if err != nil {
				log.Errorf(err.Error())
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			saveToSessionAndReturnTo(c, u, homePath)
			return
		}

		u = New(0)
		u.Email = uInfo.Email
		saveToSessionAndReturnTo(c, u, userNewPath)
		return

		// l := len(ks)
		// switch l {
		// case 0: // If no keys, then no old entities for user
		// 	log.Warningf(
		// 		"creating new user for %q without migratition from prior users",
		// 		uInfo.Email,
		// 	)
		// 	u.Email = uInfo.Email
		// 	saveToSessionAndReturnTo(c, u, userNewPath)
		// 	return
		// case 1:
		// 	log.Infof("attempting to migrate user with key %#v and email %q", ks[0], uInfo.Email)
		// 	migrateOld(c, ks)
		// 	ou := New(0)
		// 	err = dsClient.Get(c, ks[0], ou)
		// 	if err == nil {
		// 		u = u.fromOld1(ou)
		// 		_, err = dsClient.Put(c, u.Key, &u)
		// 		if err == nil {
		// 			log.Infof("successfully migrated user for email %q", uInfo.Email)
		// 			u.Loaded = true
		// 			saveToSessionAndReturnTo(c, u, homePath)
		// 			return
		// 		}
		// 	}
		// 	log.Warningf("unable to load user for key %#v due to error: %s", ks[0], err.Error())
		// 	log.Warningf(
		// 		"initiating creation of new user for %q without migrating from prior user",
		// 		uInfo.Email,
		// 	)
		// 	u.Email = uInfo.Email
		// 	saveToSessionAndReturnTo(c, u, userNewPath)
		// 	return
		// case 2:
		// 	ok := ks[0]
		// 	if ok.Name == "" {
		// 		ok = ks[1]
		// 	}
		// 	log.Infof("attempting to migrate user with key %#v and email %q", ok, uInfo.Email)
		// 	ou := NNew(ok.Name)
		// 	err = dsClient.Get(c, ou.Key, ou)
		// 	if err == nil {
		// 		u = u.fromOld2(ou)
		// 		_, err = dsClient.Put(c, u.Key, &u)
		// 		if err == nil {
		// 			log.Infof("successfully migrated user for email %q", uInfo.Email)
		// 			u.Loaded = true
		// 			saveToSessionAndReturnTo(c, u, homePath)
		// 			return
		// 		}
		// 	}
		// 	log.Warningf("unable to load user for key %#v due to error: %s", ok, err.Error())
		// 	log.Warningf(
		// 		"initiating creation of new user for %q without migrating from prior user",
		// 		uInfo.Email,
		// 	)
		// 	u.Email = uInfo.Email
		// 	saveToSessionAndReturnTo(c, u, userNewPath)
		// 	return
		// default:
		// 	log.Warningf(
		// 		"initiating creation of new user for %q without migrating from prior user",
		// 		uInfo.Email,
		// 	)
		// 	u.Email = uInfo.Email
		// 	saveToSessionAndReturnTo(c, u, userNewPath)
		// 	return
		// }
	}
}

func getUInfo(c *gin.Context, path string) (Info, error) {
	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	if retrievedState != c.Query("state") {
		return Info{}, fmt.Errorf("Invalid session state: %s", retrievedState)
	}

	conf := oauth2Config(c, path, scopes()...)
	tok, err := conf.Exchange(c, c.Query("code"))
	if err != nil {
		return Info{}, fmt.Errorf("tok error: %#v", err)
	}

	client := conf.Client(c, tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return Info{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Info{}, err
	}

	uInfo := Info{}
	var b binding.BindingBody = binding.JSON
	err = b.BindBody(body, &uInfo)
	if err != nil {
		return Info{}, err
	}
	return uInfo, nil
}

func getOAuth(c *gin.Context, id string) (OAuth, error) {
	u := NewOAuth(id)
	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return u, err
	}
	err = dsClient.Get(c, u.Key, &u)
	return u, err
}

func saveToSessionAndReturnTo(c *gin.Context, u User, path string) {
	session := sessions.Default(c)
	err := u.SaveTo(session)
	if err != nil {
		log.Errorf(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Redirect(http.StatusSeeOther, path)
	return
}

func getByEmail(c *gin.Context, email string) (User, error) {
	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return None, err
	}

	q := datastore.NewQuery(uKind).
		Ancestor(RootKey()).
		Filter("Email=", email).
		KeysOnly()

	ks, err := dsClient.GetAll(c, q, nil)
	if err != nil {
		return None, err
	}

	log.Debugf("ks: %v", ks)
	for i := range ks {
		if ks[i].ID != 0 {
			return getByID(c, ks[i].ID)
		}
	}
	return None, errors.New("unable to find user")
}

func getByID(c *gin.Context, id int64) (User, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return None, err
	}
	u := New(id)
	err = dsClient.Get(c, u.Key, &u)
	return u, err
}

// func (u OAuth) migrateOld(c *gin.Context, ks []*datastore.Key) (OAuth, error) {
// 	if ks[0].Name != "" {
// 		return u.migrateOld2(c, ks[0])
// 	}
// 	return u.migrateOld2(c, ks[1])
// }

// func (u OAuth) migrateOld1(c *gin.Context, k *datastore.Key) (OAuth, error) {
// 	dsClient, err := datastore.NewClient(c, "")
// 	if err != nil {
// 		return u, err
// 	}
//
// 	ou := New(k.ID)
// 	err = dsClient.Get(c, ou.Key, ou)
// 	if err != nil {
// 		return u, err
// 	}
//
// 	u = u.fromOld1(ou)
// 	st := NewStats(u)
// 	st, err = st.migrate(c)
// 	if err != nil {
// 		_, err = dsClient.Put(c, u.Key, &u)
// 		return u, err
// 	}
//
// 	st.Key = datastore.NameKey(sKind, u.Key.Name, nil)
// 	_, err = dsClient.PutMulti(c, []*datastore.Key{st.Key, u.Key}, []interface{}{&st, &u})
// 	return u, err
// }

// func (u OAuth) migrateOld2(c *gin.Context, k *datastore.Key) (OAuth, error) {
// 	dsClient, err := datastore.NewClient(c, "")
// 	if err != nil {
// 		return u, err
// 	}
//
// 	nu := NNew(k.Name)
// 	err = dsClient.Get(c, nu.Key, nu)
// 	if err != nil {
// 		return u, err
// 	}
//
// 	u = u.FromNUser(nu)
// 	st := NewStats(u)
// 	st, err = st.migrate(c)
// 	if err != nil {
// 		_, err = dsClient.Put(c, u.Key, &u)
// 		return u, err
// 	}
//
// 	st.Key = datastore.NameKey(sKind, u.Key.Name, nil)
// 	_, err = dsClient.PutMulti(c, []*datastore.Key{st.Key, u.Key}, []interface{}{&st, &u})
// 	return u, err
// }
//
// func (st Stats) migrate(c *gin.Context) (Stats, error) {
// 	log.Debugf("Entering")
// 	defer log.Debugf("Exiting")
//
// 	dsClient, err := datastore.NewClient(c, "")
// 	if err != nil {
// 		return st, err
// 	}
//
// 	err = dsClient.Get(c, st.Key, &st)
// 	return st, err
// }

func (u User) SaveTo(s sessions.Session) error {
	s.Set(sessionKey, sessionToken{u})
	return s.Save()
}

type sessionToken struct {
	User
}

func SessionTokenFrom(s sessions.Session) (sessionToken, bool) {
	token, ok := s.Get(sessionKey).(sessionToken)
	return token, ok
}

var None = User{}

func current(c *gin.Context) (User, bool) {
	u, ok := c.Value(currentKey).(User)
	return u, ok
	// if !ok {
	// 	return None
	// }
	// return u
}

func withCurrent(c *gin.Context, u User) {
	c.Set(currentKey, u)
}

func Current(c *gin.Context) (User, bool) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	u, found := current(c)
	if found {
		return u, true
	}

	session := sessions.Default(c)
	token, ok := SessionTokenFrom(session)
	if !ok {
		log.Warningf("missing token")
		return None, false
	}

	if token.Loaded {
		withCurrent(c, token.User)
		return token.User, true
	}

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		log.Warningf("Client error: %v", err)
		return None, false
	}

	u = New(token.ID())
	err = dsClient.Get(c, u.Key, &u)
	if err != nil {
		log.Warningf("ByID err: %s", err)
		return None, false
	}
	log.Debugf("u: %#v", u)
	withCurrent(c, u)
	return u, true
}

func RequireLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		_, ok := SessionTokenFrom(session)
		if !ok {
			log.Warningf("RequireLogin failed.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
		}
	}
}

func RequireCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, found := Current(c)
		if !found {
			log.Warningf("RequireCurrentUser failed.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
		}
	}
}

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
// 	uid, err := getuid(c)
// 	if err != nil || uid == notfound {
// 		log.errorf("getuid error: %v", err.error())
// 		c.redirect(http.statusseeother, "/")
// 		c.abort()
// 		return
// 	}
//
// 	u := new(ctx)
// 	u.id = uid
// 	if err = datastore.get(ctx, u); err != nil {
// 		log.errorf("unable to get user for id: %v", c.param("uid"))
// 		c.redirect(http.statusseeother, "/")
// 		c.abort()
// 	} else {
// 		withuser(c, u)
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

// func Fetched(ctx context.Context) *User {
// 	return From(ctx)
// }

func Gravatar(u User) template.HTML {
	return template.HTML(fmt.Sprintf(`<a href="/user/show/%d"><img src=%q alt="Gravatar" class="black-border" /></a>`, u.ID, u.Gravatar()))
}

// func Gravatar(u OAuth) template.HTML {
// 	return template.HTML(fmt.Sprintf(`<a href="/user/show/%d"><img src=%q alt="Gravatar" class="black-border" /></a>`, u.ID, u.Gravatar()))
// }
//
func NGravatar(nu *NUser) template.HTML {
	return template.HTML(fmt.Sprintf(`<a href="/user/show/%s"><img src=%q alt="Gravatar" class="black-border" /></a>`, nu.ID, nu.Gravatar()))
}

// func from(ctx context.Context, key string) (u *User) {
// 	u, _ = ctx.Value(key).(*User)
// 	return
// }

// func From(ctx context.Context) *User {
// 	return from(ctx, userKey)
// }

// func CurrentFrom(ctx context.Context) *User {
// 	return from(ctx, currentKey)
// }

func (u User) IsCurrent(c *gin.Context) bool {
	cu, found := Current(c)
	if !found {
		return false
	}
	return u.ID() == cu.ID()
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

// func UsersFrom(ctx context.Context) (us []interface{}) {
// 	us, _ = ctx.Value(usersKey).([]interface{})
// 	return
// }

func withUsers(c *gin.Context, us []interface{}) {
	c.Set(usersKey, us)
}

func withCount(c *gin.Context, cnt int64) *gin.Context {
	c.Set(countKey, cnt)
	return c
}

func CountFrom(c *gin.Context) int64 {
	cnt, ok := c.Value(countKey).(int64)
	if ok {
		return cnt
	}
	return -1
}
