package game

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/slothninja-games/sn/color"
	"github.com/SlothNinja/slothninja-games/sn/restful"
	"github.com/SlothNinja/slothninja-games/sn/send"
	gType "github.com/SlothNinja/slothninja-games/sn/type"
	"github.com/SlothNinja/slothninja-games/sn/user"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"google.golang.org/appengine/mail"
)

// Header provides fields common to all games.
type Header struct {
	gamer interface{}

	Key     *datastore.Key `datastore:"__key__"`
	Creator user.User      `datastore:"-" json:"-"`
	Users   []user.User    `datastore:"-" json:"users"`
	Stats   []user.Stats   `datastore:"-" json:"-"`
	// ID      int64          `gae:"$id"`
	// Parent  *datastore.Key `gae:"$parent"`
	// Kind    string         `gae:"$kind"`

	Type          gType.Type  `json:"type"`
	Title         string      `form:"title" json:"title"`
	Turn          int         `form:"turn" json:"turn" binding:"min=0"`
	Phase         Phase       `form:"phase" json:"phase" binding:"min=0"`
	SubPhase      SubPhase    `form:"sub-phase" json:"subPhase" binding:"min=0"`
	Round         int         `form:"round" json:"round" binding:"min=0"`
	NumPlayers    int         `form:"num-players" json:"numPlayers" binding"min=0,max=5"`
	Password      string      `form:"password" json:"-"`
	CreatorID     int64       `form:"creator-id" json:"creatorId"`
	CreatorSID    string      `form:"creator-sid" json:"creatorSId"`
	CreatorName   string      `form:"creator-name" json:"creatorName"`
	UserIDS       []int64     `form:"user-ids" json:"userIds"`
	UserSIDS      []string    `form:"user-sids" json:"userSIds"`
	UserNames     []string    `form:"user-names" json:"userNames"`
	UserEmails    []string    `form:"user-emails" json:"userEmails"`
	OrderIDS      UserIndices `form:"order-ids" json:"-"`
	CPUserIndices UserIndices `form:"cp-user-indices" json:"cpUserIndices"`
	WinnerIDS     UserIndices `form:"winner-ids" json:"winnerIndices"`
	Status        Status      `form:"status" json:"status"`
	Progress      string      `form:"progress" json:"progress"`
	Options       []string    `form:"options" json:"options"`
	OptString     string      `form:"opt-string" json:"optString"`
	SavedState    []byte      `gae:"SavedState,noindex" json:"-"`
	CreatedAt     time.Time   `form:"created-at" json:"createdAt"`
	UpdatedAt     time.Time   `form:"updated-at" json:"updatedAt"`
	UpdateCount   int         `json:"-"`
}

// func (h *Header) CTX() context.Context {
// 	return h.ctx
// }

func (h Header) ID() int64 {
	if h.Key != nil {
		return h.Key.ID
	}
	return 0
}

type headerer interface {
	GetHeader() *Header
	GetAcceptDialog() bool
	AcceptedPlayers() int
	//CurrentUserIsCurrentPlayerOrAdmin() bool
	PlayererByID(int) Playerer
	PlayererByUserID(int64) Playerer
	PlayererByIndex(int) Playerer
	Winnerers() Playerers
	User(int) user.User
	Stat(int) user.Stats
	CurrentPlayerers() Playerers
	NextPlayerer(*gin.Context, ...Playerer) Playerer
	DefaultColorMap() color.Colors
	UserLinks() template.HTML
	//PlayerLinks() template.HTML
	Private() bool
	CanAdd(user.User) bool
	CanDropout(user.User) bool
	Stub() string
	//CTX() context.Context
	//CurrentUser() *user.User
	Accept(*gin.Context, user.User) (bool, error)
	Drop(user.User) error
	IsCurrentPlayer(user.User) bool
}

func (h *Header) GetHeader() *Header {
	return h
}

type UserIndices []int

func (uis *UserIndices) Append(indices ...int)             { *uis = uis.AppendS(indices...) }
func (uis UserIndices) AppendS(indices ...int) UserIndices { return append(uis, indices...) }

func (uis UserIndices) Include(index int) bool {
	for _, i := range uis {
		if i == index {
			return true
		}
	}
	return false
}

func (uis UserIndices) RemoveS(indices ...int) UserIndices {
	for _, index := range indices {
		uis = uis.remove(index)
	}
	return uis
}

func (uis UserIndices) remove(index int) UserIndices {
	for i, indx := range uis {
		if indx == index {
			return uis.removeAt(i)
		}
	}
	return uis
}

func (uis UserIndices) removeAt(i int) UserIndices { return append(uis[:i], uis[i+1:]...) }

func NewHeader(g Gamer) *Header {
	return &Header{gamer: g}
}

type Strings []string

type ColorMaps map[gType.Type]color.Colors

var defaultColorMaps = ColorMaps{
	gType.Confucius:  color.Colors{color.Yellow, color.Purple, color.Green, color.White, color.Black},
	gType.Tammany:    color.Colors{color.Red, color.Yellow, color.Purple, color.Black, color.Brown},
	gType.ATF:        color.Colors{color.Red, color.Green, color.Purple},
	gType.GOT:        color.Colors{color.Yellow, color.Purple, color.Green, color.Black},
	gType.Indonesia:  color.Colors{color.White, color.Black, color.Green, color.Purple, color.Orange},
	gType.Gettysburg: color.Colors{color.White, color.Black},
}

func (h *Header) DefaultColorMap() color.Colors {
	return defaultColorMaps[h.Type]
}

func (h Header) ColorMapFor(c *gin.Context, u user.User) color.Map {
	cm := h.DefaultColorMap()
	if u.Key != nil {
		if p := h.PlayererByUserID(u.ID()); p != nil {
			cm = p.ColorMap(c)
		}
	}
	ps := h.gamer.(GetPlayerers).GetPlayerers()
	cMap := make(color.Map, len(ps))
	for i, u2 := range h.Users {
		cMap[u2.ID()] = cm[i]
	}
	return cMap
}

func (ss Strings) Include(s string) bool {
	for _, value := range ss {
		if s == value {
			return true
		}
	}
	return false
}

func actionPath(r *http.Request) string {
	s := strings.Split(r.URL.String(), "/")
	return s[len(s)-1]
}

func (h *Header) FromParams(c *gin.Context, t gType.Type) error {
	return h.FromForm(c, t)
}

func (h *Header) FromForm(c *gin.Context, t gType.Type) (err error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	h2 := NewHeader(nil)
	if err = restful.BindWith(c, h2, binding.FormPost); err != nil {
		return
	}

	cu, _ := user.Current(c)
	log.Debugf("h2: %#v", h2)

	if h2.Title == "" {
		h.Title = cu.Name + "'s Game"
	} else {
		h.Title = h2.Title
	}

	if h2.NumPlayers == 0 {
		h.NumPlayers = 4
		log.Debugf("h: %#v", h)
	} else {
		h.NumPlayers = h2.NumPlayers
	}

	h.Password = h2.Password
	h.Creator = cu
	h.CreatorID = cu.ID()
	h.AddUser(cu)
	h.Status = Recruiting
	h.Type = t
	return nil
}

func getType(form url.Values) gType.Type {
	sType := form.Get("game-type")
	iType, err := strconv.Atoi(sType)
	if err != nil {
		return gType.NoType
	}

	t := gType.Type(iType)
	if _, ok := gType.TypeStrings[t]; !ok {
		return gType.NoType
	}
	return t
}

func (h Header) User(index int) user.User {
	i := index
	if l := len(h.UserIDS); l > 0 {
		i = i % l
	}
	return h.Users[i]
}

func (h Header) Stat(i int) user.Stats {
	l := len(h.Stats)
	if l == 0 {
		return user.Stats{}
	}
	return h.Stats[i%l]
}

// func (h *Header) CurrentUser() *user.User {
// 	return user.CurrentFrom(h.CTX())
// }

func (h *Header) AfterLoad(c *gin.Context, gamer Gamer) error {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	var err error
	h.Users, err = h.getUsers(c)
	if err != nil {
		return err
	}

	h.Creator, err = h.getCreator(c)
	if err != nil {
		return err
	}

	// l := len(h.UserIDS)

	// ids := make([]string, l)
	// copy(ids, h.UserSIDS)

	// addID := true
	// cIndex := len(h.UserSIDS)
	// for i, id := range ids {
	// 	if id == h.CreatorSID {
	// 		addID = false
	// 		cIndex = i
	// 	}
	// }

	// if addID {
	// 	ids = append(ids, h.CreatorSID)
	// }

	// us := make([]user.User, len(ids))
	// ks := make([]*datastore.Key, len(ids))
	// for i := range us {
	// 	us[i] = user.NewUser(ids[i])
	// 	ks[i] = us[i].Key
	// 	log.Debugf("us[%d]: %#v", i, us[i])
	// 	log.Debugf("ks[%d]: %v", i, us[i].Key)
	// }

	// dsClient, err := datastore.NewClient(c, "")
	// if err != nil {
	// 	return err
	// }

	// err = dsClient.GetMulti(c, ks, us)
	// if err != nil {
	// 	return err
	// }

	// h.Users = us[:l]
	// h.Creator = us[cIndex]
	return nil
}

func (h Header) getCreator(c *gin.Context) (user.User, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return user.User{}, err
	}

	u := user.New(h.CreatorID)
	err = dsClient.Get(c, u.Key, &u)
	return u, err
}

func (h Header) getUsers(c *gin.Context) ([]user.User, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return nil, err
	}

	us := make([]user.User, len(h.UserIDS))
	ks := make([]*datastore.Key, len(h.UserIDS))
	for i := range h.UserIDS {
		us[i] = user.New(h.UserIDS[i])
		ks[i] = us[i].Key
	}

	err = dsClient.GetMulti(c, ks, us)
	return us, err
}

func include(ints []int64, i1 int64) bool {
	for _, i2 := range ints {
		if i1 == i2 {
			return true
		}
	}
	return false
}

func remove(ints []int64, i int64) []int64 {
	for index, j := range ints {
		if j == i {
			return append(ints[:index], ints[index+1:]...)
		}
	}
	return ints
}

func (h Header) CanAdd(u user.User) bool {
	return len(h.UserIDS) < h.NumPlayers && !include(h.UserIDS, u.ID())
}

func (h Header) CanDropout(u user.User) bool {
	return h.Status == Recruiting && include(h.UserIDS, u.ID())
}

func (h *Header) Stub() string {
	return strings.ToLower(h.Type.SString())
}

func (h *Header) Private() bool {
	return h.Password != ""
}

func (h Header) HasUser(u user.User) bool {
	return include(h.UserIDS, u.ID())
}

func (h *Header) RemoveUser(u2 user.User) {
	for i, u := range h.Users {
		if u.ID() == u2.ID() {
			h.Users = append(h.Users[:i], h.Users[i+1:]...)
		}
	}
	h.updateUserFields()
}

func (h *Header) updateUserFields() {
	l := len(h.Users)

	h.UserIDS = make([]int64, l)
	h.UserNames = make([]string, l)
	// h.UserSIDS = make([]string, l)
	for i, u := range h.Users {
		h.UserIDS[i] = u.ID()
		h.UserNames[i] = u.Name
		//	h.UserSIDS[i] = user.GenID(u.GoogleID)
	}
}

func (h *Header) AddUser(u user.User) {
	h.AddUsers(u)
}

func (h *Header) AddUsers(us ...user.User) {
	h.Users = append(h.Users, us...)
	h.updateUserFields()
}

// func (h *Header) IsAdmin() bool {
// 	return user.IsAdmin(h.CTX())
// }

func (h Header) CurrentPlayerer(c *gin.Context) Playerer {
	switch cps := h.CurrentPlayerers(); len(cps) {
	case 0:
		return nil
	case 1:
		return cps[0]
	default:
		return h.CurrentUserPlayerer(c)
	}
}

// CurrentPlayererFrom provides the first current player from players ps.
func (h *Header) CurrentPlayerFrom(ps Playerers) (cp Playerer) {
	if cps := h.CurrentPlayersFrom(ps); len(cps) > 0 {
		cp = cps[0]
	}
	return
}

func (h Header) CurrentUserPlayerer(c *gin.Context) Playerer {
	switch cps := h.CurrentUserPlayerers(c); len(cps) {
	case 0:
		return nil
	case 1:
		return cps[0]
	default:
		log.Warningf("CurrentUserPlayerer found %d current user players.  Returned only the first.")
		return cps[0]
	}
}

func (h Header) CurrentUserPlayerers(c *gin.Context) Playerers {
	var cps Playerers
	for _, cp := range h.CurrentPlayerers() {
		cu, _ := user.Current(c)
		if cp.User().ID() == cu.ID() {
			cps = append(cps, cp)
		} else if cu.IsAdmin {
			return append(cps, cp)
		}
	}
	return cps
}

// CurrentPlayererFor returns the current player from players ps associated with the user u.
// If no player is associated with the user, but user is admin, then returns default current player.
func (h *Header) CurrentPlayerFor(ps Playerers, u user.User) Playerer {
	for _, p := range h.CurrentPlayersFrom(ps) {
		if p.User().ID() == u.ID() {
			return p
		}
	}

	if u.IsAdmin {
		return h.CurrentPlayerFrom(ps)
	}
	return nil
}

func (h *Header) CurrentPlayerers() Playerers {
	if h.Status == Completed {
		return nil
	}

	var playerers Playerers
	for _, index := range h.CPUserIndices {
		playerers = append(playerers, h.PlayerByUserIndex(index))
	}
	return playerers
}

// CurrentPlayerers returns the current players in players.
func (h *Header) CurrentPlayersFrom(players Playerers) (ps Playerers) {
	if h.Status != Completed {
		for _, index := range h.CPUserIndices {
			ps = append(ps, PlayerByUserIndex(players, index))
		}
	}
	return
}

// ps is an optional parameter.
// If no player is provided, assume current player.
func (h *Header) NextPlayerer(c *gin.Context, ps ...Playerer) Playerer {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	cp := h.CurrentPlayerer(c)
	i := cp.Index() + 1
	log.Debugf("cp: %#v", cp)
	log.Debugf("i: %#v", i)
	if len(ps) == 1 {
		i = ps[0].Index() + 1
		log.Debugf("len(ps) == 1 => i: %#v", i)
	}
	p := h.PlayererByIndex(i)
	log.Debugf("p: %#v", p)
	return p
}

// ps is an optional parameter.
// If no player is provided, assume current player.
func (h *Header) PreviousPlayerer(c *gin.Context, ps ...Playerer) Playerer {
	cp := h.CurrentPlayerer(c)
	i := cp.Index() - 1
	if len(ps) == 1 {
		i = ps[0].Index() - 1
	}
	return h.PlayererByIndex(i)
}

func (h *Header) Winnerers() Playerers {
	if len(h.WinnerIDS) == 0 || h.WinnerIDS[0] == -1 {
		return nil
	}

	var playerers Playerers
	for _, index := range h.WinnerIDS {
		playerers = append(playerers, h.PlayerByUserIndex(index))
	}
	return playerers
}

func (h *Header) SetCurrentPlayerers(players ...Playerer) {
	switch length := len(players); {
	case length > 0:
		h.CPUserIndices = make(UserIndices, length)
		for i, p := range players {
			h.CPUserIndices[i] = p.ID()
		}
	default:
		h.CPUserIndices = nil
	}
}

func (h *Header) RemoveCurrentPlayers(ps ...Playerer) {
	if len(ps) > 0 {
		players := h.CurrentPlayerers()
		for _, rp := range ps {
			for i, p := range players {
				if p.ID() == rp.ID() {
					players = append(players[:i], players[i+1:]...)
					break
				}
			}
		}
		h.SetCurrentPlayerers(players...)
	}
}

func (h *Header) isCP(uIndex int) bool {
	if len(h.CPUserIndices) == 0 || h.CPUserIndices[0] == -1 || uIndex == -1 {
		return false
	}

	for _, cpi := range h.CPUserIndices {
		if cpi == uIndex {
			return true
		}
	}
	return false
}

// IsCurrentPlayer returns true if the specified user is the current player.
func (h Header) IsCurrentPlayer(u user.User) bool {
	return h.isCP(h.indexFor(u))
}

// IsCurrentPlayer returns ture if the user is the current player or an admin.
func (h Header) IsCurrentPlayerOrAdmin(u user.User) bool {
	return u.IsAdmin || h.IsCurrentPlayer(u)
}

func (h Header) isCurrentPlayerOrAdmin(c *gin.Context, u user.User) bool {
	return u.IsAdmin || h.IsCurrentPlayer(u)
}

// CurrentUserIsCurrentPlayerOrAdmin returns true if current user is the current player or is an administrator.
// Deprecated in favor of CUserIsCPlayerOrAdmin
func (h Header) CurrentUserIsCurrentPlayerOrAdmin(c *gin.Context) bool {
	log.Warningf("CurrentUserIsCurrentPlayerOrAdmin is deprecated in favor of CUserIsCPlayerOrAdmin.")
	cu, _ := user.Current(c)
	return h.isCurrentPlayerOrAdmin(c, cu)
}

// CUserIsCPlayerOrAdmin returns true if current user is the current player or is an administrator.
func (h Header) CUserIsCPlayerOrAdmin(c *gin.Context) bool {
	cu, _ := user.Current(c)
	return h.isCurrentPlayerOrAdmin(c, cu)
}

func (h Header) PlayerIsUser(p Playerer, u user.User) bool {
	return p != nil && h.UserIDFor(p) == u.ID()
}

func (h *Header) IsW(uIndex int) bool {
	return h.WinnerIDS.Include(uIndex)
}

func (h Header) IsWinner(u user.User) bool {
	for _, p := range h.PlayerersByUser(u) {
		if h.WinnerIDS.Include(p.ID()) {
			return true
		}
	}
	return false
}

func (h *Header) UserLinks() template.HTML {
	links := make([]string, len(h.UserIDS))
	for i, uid := range h.UserIDS {
		links[i] = string(h.UserLinkFor(uid))
	}
	return template.HTML(restful.ToSentence(links))
}

func (h *Header) UserLinkFor(uid int64) template.HTML {
	return user.LinkFor(uid, h.NameByUID(uid))
}

func (h Header) PlayerLinkByID(c *gin.Context, pid int) template.HTML {
	i := pid % len(h.UserIDS)
	uid := h.UserIDS[i]

	cu, found := user.Current(c)
	cp := h.isCP(pid)

	var me bool
	if !found {
		me = cu.ID() == uid
	}

	w := h.IsW(pid)
	n := h.NameByPID(pid)

	result := fmt.Sprintf(`<a href="/user/show/%d" >%s</a>`, uid, n)
	switch h.Status {
	case Running:
		switch {
		case cp && me:
			result = fmt.Sprintf(`<a href="/user/show/%d" class="current-player me">%s</a>`, uid, n)
		case cp:
			result = fmt.Sprintf(`<a href="/user/show/%d" class="current-player">%s</a>`, uid, n)
		}
	case Completed:
		switch {
		case w && me:
			result = fmt.Sprintf(`<a href="/user/show/%d" class="winner me">%s</a>`, uid, n)
		case w:
			result = fmt.Sprintf(`<a href="/user/show/%d" class="winner">%s</a>`, uid, n)
		}
	}
	return template.HTML(result)
}

func (h Header) PlayerLinks(c *gin.Context) template.HTML {
	if h.Status == Recruiting {
		return h.UserLinks()
	}

	links := make([]string, len(h.OrderIDS))
	for i, index := range h.OrderIDS {
		links[i] = string(h.PlayerLinkByID(c, index))
	}
	return template.HTML(restful.ToSentence(links))
}

func (h Header) CurrentPlayerLinks(c *gin.Context) template.HTML {
	cps := h.CPUserIndices
	if len(cps) == 0 || h.Status != Running {
		return "None"
	}

	links := make([]string, len(cps))
	for j, i := range cps {
		links[j] = string(h.PlayerLinkByID(c, i))
	}
	return template.HTML(restful.ToSentence(links))
}

func (h *Header) NoCurrentPlayer() bool {
	return len(h.CPUserIndices) == 0
}

func (h *Header) CurrentPlayerLabel() string {
	if length := len(h.CPUserIndices); length == 1 {
		return "Current Player"
	}
	return "Current Players"
}

func (h *Header) AcceptedPlayers() int {
	return len(h.UserIDS)
}

// PlayererByID returns the player having the id.
func (h *Header) PlayererByID(id int) (p Playerer) {
	return PlayererByID(h.gamer.(GetPlayerers).GetPlayerers(), id)
}

// PlayererByID returns the player from ps having the id.
func PlayererByID(ps Playerers, id int) (p Playerer) {
	for _, p2 := range ps {
		if p2.ID() == id {
			p = p2
			return
		}
	}
	return
}

func (h *Header) PlayererByColor(ctx *gin.Context, c color.Color) Playerer {
	for _, p := range h.gamer.(GetPlayerers).GetPlayerers() {
		if p.Color(ctx) == c {
			return p
		}
	}
	return nil
}

// PlayerBySID provides the player having the id represented by the string.
func (h *Header) PlayerBySID(sid string) (p Playerer) {
	if id, err := strconv.Atoi(sid); err == nil {
		p = h.PlayererByID(id)
	}
	return
}

// PlayerBySID provides the player in ps having the id represented by the string.
func PlayerBySID(ps Playerers, sid string) (p Playerer) {
	if id, err := strconv.Atoi(sid); err == nil {
		p = PlayererByID(ps, id)
	}
	return
}

// PlayererByUserID returns the player associated with the user id
func (h *Header) PlayererByUserID(id int64) Playerer {
	return PlayererByUserID(h.gamer.(GetPlayerers).GetPlayerers(), id)
}

// PlayererByUserID returns the player from ps associated with the user id
func PlayererByUserID(ps Playerers, id int64) Playerer {
	for _, p2 := range ps {
		if p2.User().ID() == id {
			return p2
		}
	}
	return nil
}

func (h Header) PlayerersByUser(user user.User) Playerers {
	var ps Playerers
	for _, p := range h.gamer.(GetPlayerers).GetPlayerers() {
		if p.User().ID() == user.ID() {
			ps = append(ps, p)
		}
	}
	return ps
}

func (h *Header) PlayerByUserIndex(index int) Playerer {
	for _, p := range h.gamer.(GetPlayerers).GetPlayerers() {
		if p.ID() == index {
			return p
		}
	}
	return nil
}

// PlayerByUserIndex returns the player from players ps having the provided user index.
func PlayerByUserIndex(ps Playerers, index int) (p Playerer) {
	for _, p2 := range ps {
		if p2.ID() == index {
			p = p2
			return
		}
	}
	return
}

// PlayererByIndex returns the player at the index i in the ring of players ps
// Convenience method that automatically wraps-around based on number of players.
// TODO: Deprecated
func (h *Header) PlayererByIndex(i int) Playerer {
	return PlayererByIndex(h.gamer.(GetPlayerers).GetPlayerers(), i)
}

// PlayererByIndex returns the player at the index i in the ring of players ps
// Wraps-around based on number of players.
func PlayererByIndex(ps Playerers, i int) (p Playerer) {
	l := len(ps)
	if r := i % l; r < 0 {
		p = ps[l+r]
	} else {
		p = ps[r]
	}
	return
}

type Phase int

func (p Phase) Int() int {
	return int(p)
}

type PhaseNameMap map[Phase]string
type PhaseNameMaps map[gType.Type]PhaseNameMap

func registerPhaseNames(t gType.Type, names PhaseNameMap) {
	if phaseNameMaps == nil {
		phaseNameMaps = make(PhaseNameMaps, len(gType.Types))
	}
	phaseNameMaps[t] = names
}

func registerSubPhaseNames(t gType.Type, names SubPhaseNameMap) {
	if subPhaseNameMaps == nil {
		subPhaseNameMaps = make(SubPhaseNameMaps, len(gType.Types))
	}
	subPhaseNameMaps[t] = names
}

type factoryMap map[gType.Type]Factory

var factories factoryMap

type Factory func(int64) Gamer

func Register(t gType.Type, f Factory, p PhaseNameMap, sp SubPhaseNameMap) {
	if factories == nil {
		factories = make(factoryMap, len(gType.Types))
	}
	factories[t] = f
	registerPhaseNames(t, p)
	registerSubPhaseNames(t, sp)
}

func (h *Header) PhaseName() string {
	if phaseNameMaps == nil {
		return ""
	}
	if names, ok := phaseNameMaps[h.Type]; ok {
		return names[h.Phase]
	}
	return ""
}

type SubPhase int
type SubPhaseNameMap map[SubPhase]string
type SubPhaseNameMaps map[gType.Type]SubPhaseNameMap

func (h *Header) SubPhaseName() string {
	if subPhaseNameMaps == nil {
		return ""
	}
	if names, ok := subPhaseNameMaps[h.Type]; ok {
		return names[h.SubPhase]
	}
	return ""
}

func (s SubPhase) Int() int {
	return int(s)
}

var phaseNameMaps PhaseNameMaps
var subPhaseNameMaps SubPhaseNameMaps

func (h *Header) ValidateHeader() error {
	if len(h.UserIDS) > h.NumPlayers {
		return fmt.Errorf("UserIDS can't be greater than the number of players.")
	}
	return nil
}

func (h Header) SendTurnNotificationsTo(c *gin.Context, ps ...Playerer) error {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	subject := fmt.Sprintf("It's your turn in %s (%s #%d).", h.Type, h.Title, h.ID)
	url := fmt.Sprintf(`<a href="http://www.slothninja.com/%s/game/show/%d">here</a>`, h.Type.Prefix(), h.ID)
	body := fmt.Sprintf(`<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
		<html>
			<head>
				<meta http-equiv="content-type" content="text/html; charset=ISO-8859-1">
			</head>
			<body bgcolor="#ffffff" text="#000000">
				<p>%s</p>
				<p>You can take your turn %s.</p>
			</body>
		</html>`, subject, url)

	m := &mail.Message{
		Sender:   "webmaster@slothninja.com",
		Subject:  subject,
		HTMLBody: body,
	}

	for _, p := range ps {
		if p.User().EmailNotifications {
			m.To = []string{p.User().Email}
			if err := send.Message(c, m); err != nil {
				log.Errorf("sending notification: %#v\n generated error: %v", m, err)
			}
		}
	}

	return nil
}

func (h Header) indexFor(u user.User) int {
	sid := user.GenID(u.GoogleID)
	for i := range h.UserSIDS {
		if h.UserSIDS[i] == sid {
			return i
		}
	}

	return -1
}
