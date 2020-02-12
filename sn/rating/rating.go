package rating

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/SlothNinja/glicko"
	"github.com/SlothNinja/slothninja-games/sn/contest"
	"github.com/SlothNinja/slothninja-games/sn/log"
	"github.com/SlothNinja/slothninja-games/sn/restful"
	gType "github.com/SlothNinja/slothninja-games/sn/type"
	"github.com/SlothNinja/slothninja-games/sn/user"
	"github.com/gin-gonic/gin"
	"go.chromium.org/gae/service/datastore"
	"go.chromium.org/gae/service/info"
	"go.chromium.org/gae/service/taskqueue"
	"go.chromium.org/luci/common/errors"
	"golang.org/x/net/context"
)

const (
	currentRatingsKey = "CurrentRatings"
	projectedKey      = "Projected"
	homePath          = "/"
)

func CurrentRatingsFrom(ctx context.Context) (rs CurrentRatings) {
	rs, _ = ctx.Value(currentRatingsKey).(CurrentRatings)
	return
}

func ProjectedFrom(ctx context.Context) (r *Rating) {
	r, _ = ctx.Value(projectedKey).(*Rating)
	return
}

func AddRoutes(prefix string, engine *gin.Engine) {
	g1 := engine.Group(prefix + "s")
	g1.POST("/userUpdate", updateUser)

	g1.GET("/update/:type", Update)

	g1.GET("/show/:type", Index)

	g1.POST("/show/:type/json", JSONFilteredAction)
}

// Ratings
type Ratings []*Rating
type Rating struct {
	ID     int64          `gae:"$id"`
	Parent *datastore.Key `gae:"$parent"`
	Common
}

type CurrentRatings []*CurrentRating
type CurrentRating struct {
	ID     string         `gae:"$id"`
	Parent *datastore.Key `gae:"$parent"`
	Common
}

type Common struct {
	generated bool
	Type      gType.Type
	R         float64
	RD        float64
	Low       float64
	High      float64
	Leader    bool
	CreatedAt restful.CTime
	UpdatedAt restful.UTime
}

func (r *CurrentRating) Rank() *glicko.Rank {
	return &glicko.Rank{
		R:  r.R,
		RD: r.RD,
	}
}

func New(ctx context.Context, pk *datastore.Key, t gType.Type, params ...float64) *Rating {
	r, rd := defaultR, defaultRD
	if len(params) == 2 {
		r, rd = params[0], params[1]
	}

	rating := new(Rating)
	rating.Parent = pk
	rating.R = r
	rating.RD = rd
	rating.Low = r - (2.0 * rd)
	rating.High = r + (2.0 * rd)
	rating.Type = t
	return rating
}

func NewCurrent(ctx context.Context, pk *datastore.Key, t gType.Type, params ...float64) *CurrentRating {
	r, rd := defaultR, defaultRD
	if len(params) == 2 {
		r, rd = params[0], params[1]
	}

	rating := new(CurrentRating)
	rating.ID = t.SString()
	rating.Parent = pk
	rating.R = r
	rating.RD = rd
	rating.Low = r - (2.0 * rd)
	rating.High = r + (2.0 * rd)
	rating.Type = t
	return rating
}

const (
	defaultR  float64 = 1500
	defaultRD float64 = 350
)

const (
	rKind  = "Rating"
	crKind = "CurrentRating"
)

func singleError(e error) error {
	if e == nil {
		return e
	}
	if me, ok := e.(errors.MultiError); ok {
		return me[0]
	}
	return e
}

// Get Current Rating for gType.Type and user associated with uKey
func Get(ctx context.Context, uKey *datastore.Key, t gType.Type) (*CurrentRating, error) {
	ratings, err := GetMulti(ctx, []*datastore.Key{uKey}, t)
	return ratings[0], singleError(err)
}

func GetMulti(ctx context.Context, uKeys []*datastore.Key, t gType.Type) (CurrentRatings, error) {
	ratings := make(CurrentRatings, len(uKeys))
	for i, uKey := range uKeys {
		ratings[i] = NewCurrent(ctx, uKey, t)
	}

	err := datastore.Get(ctx, ratings)
	if err == nil {
		return ratings, nil
	}

	me := err.(errors.MultiError)
	isNil := true
	for i := range uKeys {
		if me[i] == datastore.ErrNoSuchEntity {
			ratings[i].generated = true
			me[i] = nil
		} else {
			isNil = false
		}
	}
	if isNil {
		return ratings, nil
	} else {
		return ratings, me
	}
}

func GetAll(ctx context.Context, uKey *datastore.Key) (rs CurrentRatings, err error) {
	rs = make(CurrentRatings, len(gType.Types))
	for i, t := range gType.Types {
		rs[i] = NewCurrent(ctx, uKey, t)
		rs[i].Parent = uKey
		rs[i].ID = t.SString()
	}

	if err = datastore.Get(ctx, rs); err == nil {
		return
	}

	var (
		merr     errors.MultiError
		ok, enil bool
	)

	if merr, ok = err.(errors.MultiError); !ok {
		return
	}

	enil = true
	for i, e := range merr {
		if e == datastore.ErrNoSuchEntity {
			rs[i].generated = true
			merr[i] = nil
		} else if e != nil {
			enil = false
		}
	}

	if enil {
		err = nil
	} else {
		err = merr
	}

	return
}

func GetHistory(ctx context.Context, uKey *datastore.Key, t gType.Type) (Ratings, error) {
	ratings := make(Ratings, len(gType.Types))
	keys := make([]*datastore.Key, len(gType.Types))
	for i, t := range gType.Types {
		ratings[i] = New(ctx, uKey, t)
		if ok := datastore.PopulateKey(ratings[i], keys[i]); !ok {
			return nil, fmt.Errorf("Unable to populate rating with key.")
		}
	}
	if err := datastore.Get(ctx, ratings); err != nil {
		return nil, err
	}
	return ratings, nil
}

func GetFor(ctx context.Context, t gType.Type) (CurrentRatings, error) {
	q := datastore.NewQuery(crKind).Ancestor(user.RootKey(ctx)).Eq("Type", t).Order("-Low").KeysOnly(true)

	var ks []*datastore.Key
	if err := datastore.GetAll(ctx, q, &ks); err != nil {
		return nil, err
	}

	length := len(ks)
	ratings := make(CurrentRatings, length)
	for i := range ratings {
		ratings[i] = NewCurrent(ctx, nil, t)
		if ok := datastore.PopulateKey(ratings[i], ks[i]); !ok {
			return nil, fmt.Errorf("Unable to populate current rating with key.")
		}
	}
	if err := datastore.Get(ctx, ratings); err != nil {
		return nil, err
	}
	return ratings, nil
}

func (rs CurrentRatings) Projected(ctx context.Context, cm contest.ContestMap) (CurrentRatings, error) {
	ratings := make(CurrentRatings, len(rs))
	for i, r := range rs {
		var err error
		if ratings[i], err = r.Projected(ctx, cm[r.Type]); err != nil {
			return nil, err
		}
	}
	return ratings, nil
}

func (r *CurrentRating) Projected(ctx context.Context, cs contest.Contests) (*CurrentRating, error) {
	log.Debugf(ctx, "Entering r.Projected")
	defer log.Debugf(ctx, "Exiting r.Projected")

	l := len(cs)
	if l == 0 && r.generated {
		log.Debugf(ctx, "rating: %#v generated: %v", r, r.generated)
		return r, nil
	}

	gcs := make(glicko.Contests, l)
	for i, c := range cs {
		gcs[i] = glicko.NewContest(c.Outcome, c.R, c.RD)
	}

	rating, err := glicko.UpdateRating(r.Rank(), gcs)
	if err != nil {
		return nil, err
	}

	return NewCurrent(ctx, r.Parent, r.Type, rating.R, rating.RD), nil
}

func (r *CurrentRating) Generated() bool {
	return r.generated
}

func Index(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	t := gType.ToType[c.Param("type")]
	c.HTML(http.StatusOK, "rating/index", gin.H{
		"Type":      t,
		"Heading":   "Ratings: " + t.String(),
		"Types":     gType.Types,
		"Context":   ctx,
		"VersionID": info.VersionID(ctx),
		"CUser":     user.CurrentFrom(ctx),
	})
}

func getAllQuery(ctx context.Context) *datastore.Query {
	return datastore.NewQuery(crKind).Ancestor(user.RootKey(ctx))
}

func getFiltered(ctx context.Context, t gType.Type, leader bool, offset, limit int32) (CurrentRatings, int64, error) {
	q := getAllQuery(ctx)

	if leader {
		q = q.Eq("Leader", true)
	}

	if t != gType.NoType {
		q = q.Eq("Type", t)
	}

	var cnt int64
	if count, err := datastore.Count(ctx, q); err != nil {
		return nil, 0, err
	} else {
		cnt = count
	}

	q = q.Offset(offset).Limit(limit).Order("-Low").KeysOnly(true)

	var ks []*datastore.Key
	if err := datastore.GetAll(ctx, q, &ks); err != nil {
		return nil, 0, fmt.Errorf("rating#getFiltered q.GetAll error: %v", err)
	}

	rs := make(CurrentRatings, len(ks))
	for i, k := range ks {
		rs[i] = NewCurrent(ctx, k.Parent(), t)
	}

	if err := datastore.Get(ctx, rs); err != nil {
		return nil, 0, fmt.Errorf("rating#getFiltered datastore.Get error: %v", err)
	} else {
		return rs, cnt, err
	}
}

func (rs CurrentRatings) getUsers(ctx context.Context) (user.Users, error) {
	log.Debugf(ctx, "Entering getUsers")
	defer log.Debugf(ctx, "Exiting getUsers")

	log.Debugf(ctx, "rs: %#v", rs)
	us := make(user.Users, len(rs))
	for i := range rs {
		log.Debugf(ctx, "rs[i]: %#v", rs[i])
		us[i] = user.New(ctx)
		if ok := datastore.PopulateKey(us[i], rs[i].Parent); !ok {
			log.Debugf(ctx, "Unable to populate user with key: %v", rs[i].Parent)
			return nil, fmt.Errorf("Unable to populate user with key.")
		}
	}

	if err := datastore.Get(ctx, us); err != nil {
		log.Debugf(ctx, "Get error: %v", err)
		return nil, err
	}
	return us, nil
}

func (rs CurrentRatings) getProjected(ctx context.Context) (ps CurrentRatings, err error) {
	log.Debugf(ctx, "Entering getProjected")
	defer log.Debugf(ctx, "Exiting getProjected")

	ps = make(CurrentRatings, len(rs))
	var cs contest.Contests
	for i, r := range rs {
		uKey := r.Parent

		if cs, err = contest.UnappliedFor(ctx, uKey, r.Type); err != nil {
			return
		}

		if ps[i], err = r.Projected(ctx, cs); err != nil {
			return
		}

		if r.generated && r.generated && len(cs) == 0 {
			ps[i].generated = true
		}
	}
	return
}

func For(ctx context.Context, u *user.User, t gType.Type) (*CurrentRating, error) {
	return Get(ctx, datastore.KeyForObj(ctx, u), t)
}

func MultiFor(ctx context.Context, u *user.User) (CurrentRatings, error) {
	return GetAll(ctx, datastore.KeyForObj(ctx, u))
}

// AddMulti has a limit of 100 tasks.  Thus, the batching.
func Update(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	var tk *taskqueue.Task

	tp := c.Param("type")
	q := user.AllQuery(ctx).KeysOnly(true)
	path := "/ratings/userUpdate"
	ts := make([]*taskqueue.Task, 0, 100)
	o := taskqueue.RetryOptions{RetryLimit: 5}

	tparams := make(url.Values)
	tparams.Set("type", tp)

	if err := datastore.Run(ctx, q, func(k *datastore.Key) {
		log.Debugf(ctx, "k: %s", k)
		tparams.Set("uid", fmt.Sprintf("%v", k.IntID()))
		tk = taskqueue.NewPOSTTask(path, tparams)
		tk.RetryOptions = &o
		ts = append(ts, tk)
	}); err != nil {
		log.Errorf(ctx, err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := taskqueue.Add(ctx, "", ts...); err != nil {
		log.Errorf(ctx, err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	}

}

func updateUser(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	var err error

	u := user.New(ctx)
	if u.ID, err = strconv.ParseInt(c.PostForm("uid"), 10, 64); err != nil {
		log.Errorf(ctx, "Invalid uid: %s received", c.PostForm("uid"))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err = datastore.Get(ctx, u); err != nil {
		log.Errorf(ctx, "Unable to find user for id: %s", c.PostForm("uid"))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	t := gType.ToType[c.PostForm("type")]

	var r *CurrentRating
	if r, err = For(ctx, u, t); err != nil {
		log.Errorf(ctx, "Unable to find rating for userid: %s", c.PostForm("uid"))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var cs contest.Contests
	if cs, err = contest.UnappliedFor(ctx, datastore.KeyForObj(ctx, u), t); err != nil {
		log.Errorf(ctx, "Ratings update error when getting unapplied contests for user ID: %v.\n Error: %s", u.ID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return

	}

	var p *CurrentRating
	if p, err = r.Projected(ctx, cs); err != nil {
		log.Errorf(ctx, "Ratings update error when getting projected rating for user ID: %v\n Error: %s", u.ID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if time.Since(time.Time(r.UpdatedAt)) < 504*time.Hour {
		log.Debugf(ctx, "Did not update rating for user ID: %v", u.ID)
		log.Debugf(ctx, "Rating updated %s ago.", time.Since(time.Time(r.UpdatedAt)))
		return
	}

	log.Debugf(ctx, "Processing rating update for user ID: %v", u.ID)
	log.Debugf(ctx, "Projected rating: %#v", p)
	log.Debugf(ctx, "Unapplied contest count: %v", len(cs))

	es := make([]interface{}, 0)

	log.Debugf(ctx, "Reached for that ranges projected")
	const threshold float64 = 200.0
	// Update leader value to indicate whether present on leader boards
	p.Leader = p.RD < threshold

	if !p.Generated() {
		r := New(ctx, p.Parent, p.Type, p.R, p.RD)
		es = append(es, p, r)
	}

	log.Debugf(ctx, "Reached for that unapplied contests")
	for _, c := range cs {
		c.Applied = true
		es = append(es, c)
	}

	if err := datastore.RunInTransaction(ctx, func(tc context.Context) error {
		return datastore.Put(tc, es)
	}, &datastore.TransactionOptions{XG: true}); err != nil {
		log.Errorf(ctx, "Ratings update err when saving updated ratings for user ID: %v\n Error: %s", u.ID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	log.Debugf(ctx, "Reached RunInTransaction")
}

func Fetch(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	if CurrentRatingsFrom(ctx) != nil {
		return
	}

	u := user.Fetched(ctx)
	if u == nil {
		restful.AddErrorf(ctx, "Unable to get ratings.")
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	if rs, err := MultiFor(ctx, u); err != nil {
		restful.AddErrorf(ctx, err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
	} else {
		c.Set(currentRatingsKey, rs)
	}
}

func Fetched(ctx context.Context) CurrentRatings {
	return CurrentRatingsFrom(ctx)
}

func FetchProjected(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	if ProjectedFrom(ctx) != nil {
		return
	}

	rs := Fetched(ctx)
	if rs == nil {
		restful.AddErrorf(ctx, "Unable to get projected ratings")
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	cm, err := contest.Unapplied(ctx, datastore.KeyForObj(ctx, user.Fetched(ctx)))
	if err != nil {
		restful.AddErrorf(ctx, err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	if pr, err := rs.Projected(ctx, cm); err != nil {
		restful.AddErrorf(ctx, err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
	} else {
		context.WithValue(ctx, projectedKey, pr)
	}
}

func Projected(ctx context.Context) (pr Ratings) {
	pr, _ = ctx.Value("Projected").(Ratings)
	return
}

//func (r *Rating) Save(ch chan<- datastore.Property) error {
//	// Time stamp
//	t := time.Now()
//	if r.CreatedAt.IsZero() {
//		r.CreatedAt = t
//	}
//	r.UpdatedAt = t
//	return datastore.SaveStruct(r, ch)
//}
//
//func (r *Rating) Load(ch <-chan datastore.Property) error {
//	return datastore.LoadStruct(r, ch)
//}

type jRating struct {
	Type template.HTML `json:"type"`
	R    float64       `json:"r"`
	RD   float64       `json:"rd"`
	Low  float64       `json:"low"`
	High float64       `json:"high"`
}

type jCombined struct {
	Rank      int           `json:"rank"`
	Gravatar  template.HTML `json:"gravatar"`
	Name      template.HTML `json:"name"`
	Type      template.HTML `json:"type"`
	Current   template.HTML `json:"current"`
	Projected template.HTML `json:"projected"`
}

func JSONIndexAction(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	var u *user.User
	uid, err := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		log.Errorf(ctx, "rating#JSONIndexAction BySID Error: %s", err)
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}
	u2 := user.New(ctx)
	u2.ID = uid
	if err := datastore.Get(ctx, u2); err != nil {
		log.Errorf(ctx, "rating#JSONIndexAction unable to find user for uid: %d", uid)
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}
	u = u2

	rs, err := MultiFor(ctx, u)
	if err != nil {
		log.Errorf(ctx, "rating#JSONIndexAction MultiFor Error: %s", err)
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	ps, err := rs.getProjected(ctx)
	if err != nil {
		log.Errorf(ctx, "rating#getProjected Error: %s", err)
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	if data, err := singleUser(ctx, u, rs, ps); err != nil {
		c.JSON(http.StatusOK, fmt.Sprintf("%v", err))
	} else {
		c.JSON(http.StatusOK, data)
	}
}

func JSONFilteredAction(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	t := gType.Get(ctx)

	var offset, limit int32 = 0, -1
	if o, err := strconv.ParseInt(c.PostForm("start"), 10, 64); err == nil && o >= 0 {
		offset = int32(o)
	}

	if l, err := strconv.ParseInt(c.PostForm("length"), 10, 64); err == nil {
		limit = int32(l)
	}

	rs, cnt, err := getFiltered(ctx, t, true, offset, limit)
	if err != nil {
		log.Errorf(ctx, "rating#getFiltered Error: %s", err)
		return
	}
	log.Debugf(ctx, "rs: %#v", rs)

	us, err := rs.getUsers(ctx)
	if err != nil {
		log.Errorf(ctx, "rating#getUsers Error: %s", err)
		return
	}
	log.Debugf(ctx, "us: %#v", us)

	ps, err := rs.getProjected(ctx)
	if err != nil {
		log.Errorf(ctx, "rating#getProjected Error: %s", err)
		return
	}
	log.Debugf(ctx, "ps: %#v", ps)

	if data, err := toCombined(ctx, us, rs, ps, offset, cnt); err != nil {
		log.Debugf(ctx, "toCombined error: %v", err)
		c.JSON(http.StatusOK, fmt.Sprintf("%v", err))
	} else {
		c.JSON(http.StatusOK, data)
	}
}

type jCombinedRatingsIndex struct {
	Data            []*jCombined `json:"data"`
	Draw            int          `json:"draw"`
	RecordsTotal    int          `json:"recordsTotal"`
	RecordsFiltered int          `json:"recordsFiltered"`
}

func (r *CurrentRating) String() string {
	return fmt.Sprintf("%.f (%.f : %.f)", r.Low, r.R, r.RD)
}

func singleUser(ctx context.Context, u *user.User, rs, ps CurrentRatings) (table *jCombinedRatingsIndex, err error) {
	log.Debugf(ctx, "Entering singleUser")
	defer log.Debugf(ctx, "Exiting singleUser")

	table = new(jCombinedRatingsIndex)
	l1, l2 := len(rs), len(ps)
	if l1 != l2 {
		err = fmt.Errorf("Length mismatch between ratings and projected ratings l1: %d l2: %d.", l1, l2)
		return
	}

	table.Data = make([]*jCombined, 0)
	for i, r := range rs {
		if p := ps[i]; !p.generated {
			table.Data = append(table.Data, &jCombined{
				Gravatar:  user.Gravatar(u),
				Name:      u.Link(),
				Type:      template.HTML(r.Type.String()),
				Current:   template.HTML(r.String()),
				Projected: template.HTML(p.String()),
			})
		}
	}

	c := restful.GinFrom(ctx)
	var draw int
	if draw, err = strconv.Atoi(c.PostForm("draw")); err != nil {
		log.Debugf(ctx, "strconv.Atoi error: %v", err)
		return
	}

	table.Draw = draw
	table.RecordsTotal = l1
	table.RecordsFiltered = l2
	return
}
func toCombined(ctx context.Context, us user.Users, rs, ps CurrentRatings, o int32, cnt int64) (*jCombinedRatingsIndex, error) {
	table := new(jCombinedRatingsIndex)
	l1, l2 := len(rs), len(ps)
	if l1 != l2 {
		return nil, fmt.Errorf("Length mismatch between ratings and projected ratings l1: %d l2: %d.", l1, l2)
	}
	table.Data = make([]*jCombined, 0)
	for i, r := range rs {
		if !r.generated {
			table.Data = append(table.Data, &jCombined{
				Rank:      i + int(o) + 1,
				Gravatar:  user.Gravatar(us[i]),
				Name:      us[i].Link(),
				Type:      template.HTML(r.Type.String()),
				Current:   template.HTML(r.String()),
				Projected: template.HTML(ps[i].String()),
			})
		}
	}

	c := restful.GinFrom(ctx)
	if draw, err := strconv.Atoi(c.PostForm("draw")); err != nil {
		return nil, err
	} else {
		table.Draw = draw
	}

	table.RecordsTotal = int(cnt)
	table.RecordsFiltered = int(cnt)
	return table, nil
}

func IncreaseFor(ctx context.Context, u *user.User, t gType.Type, cs contest.Contests) (cr, nr *CurrentRating, err error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	k := datastore.KeyForObj(ctx, u)

	var ucs contest.Contests
	if ucs, err = contest.UnappliedFor(ctx, k, t); err != nil {
		return
	}

	var r *CurrentRating
	if r, err = For(ctx, u, t); err != nil {
		return
	}

	if cr, err = r.Projected(ctx, ucs); err != nil {
		return
	}

	nr, err = r.Projected(ctx, append(ucs, filterContestsFor(cs, k)...))
	return
}

func filterContestsFor(cs contest.Contests, pk *datastore.Key) (fcs contest.Contests) {
	for _, c := range cs {
		if c.Parent.Equal(pk) {
			fcs = append(fcs, c)
		}
	}
	return
}
