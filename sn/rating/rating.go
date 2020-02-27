package rating

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/glicko"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/slothninja-games/sn/contest"
	gType "github.com/SlothNinja/slothninja-games/sn/type"
	"github.com/SlothNinja/slothninja-games/sn/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine"
	"google.golang.org/appengine/taskqueue"
)

const (
	currentRatingsKey = "CurrentRatings"
	projectedKey      = "Projected"
	homePath          = "/"
)

func CurrentRatingsFrom(c *gin.Context) (rs []*CurrentRating) {
	rs, _ = c.Value(currentRatingsKey).([]*CurrentRating)
	return
}

func ProjectedFrom(c *gin.Context) (r *Rating) {
	r, _ = c.Value(projectedKey).(*Rating)
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
// type Ratings []*Rating
type Rating struct {
	Key *datastore.Key `datastore:"__key__"`
	Common
}

// type CurrentRatings []*CurrentRating
type CurrentRating struct {
	Key *datastore.Key `datastore:"__key__"`
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
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *CurrentRating) Rank() *glicko.Rank {
	return &glicko.Rank{
		R:  r.R,
		RD: r.RD,
	}
}

func New(pk *datastore.Key, t gType.Type, params ...float64) *Rating {
	r, rd := defaultR, defaultRD
	if len(params) == 2 {
		r, rd = params[0], params[1]
	}

	rating := new(Rating)
	rating.Key = datastore.IncompleteKey(rKind, pk)
	rating.R = r
	rating.RD = rd
	rating.Low = r - (2.0 * rd)
	rating.High = r + (2.0 * rd)
	rating.Type = t
	return rating
}

func NewCurrent(pk *datastore.Key, t gType.Type, params ...float64) *CurrentRating {
	r, rd := defaultR, defaultRD
	if len(params) == 2 {
		r, rd = params[0], params[1]
	}

	rating := new(CurrentRating)
	rating.Key = datastore.NameKey(crKind, t.SString(), pk)
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
	if me, ok := e.(datastore.MultiError); ok {
		return me[0]
	}
	return e
}

// Get Current Rating for gType.Type and user associated with uKey
func Get(c *gin.Context, uKey *datastore.Key, t gType.Type) (*CurrentRating, error) {
	ratings, err := GetMulti(c, []*datastore.Key{uKey}, t)
	return ratings[0], singleError(err)
}

func GetMulti(c *gin.Context, uKeys []*datastore.Key, t gType.Type) ([]*CurrentRating, error) {
	ratings := make([]*CurrentRating, len(uKeys))
	ks := make([]*datastore.Key, len(uKeys))
	for i, uKey := range uKeys {
		ratings[i] = NewCurrent(uKey, t)
		ks[i] = ratings[i].Key
	}

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return nil, err
	}

	err = dsClient.GetMulti(c, ks, ratings)
	if err == nil {
		return ratings, nil
	}

	me := err.(datastore.MultiError)
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

func GetAll(c *gin.Context, uKey *datastore.Key) ([]*CurrentRating, error) {
	rs := make([]*CurrentRating, len(gType.Types))
	ks := make([]*datastore.Key, len(gType.Types))
	for i, t := range gType.Types {
		rs[i] = NewCurrent(uKey, t)
		ks[i] = rs[i].Key
	}

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return nil, err
	}

	err = dsClient.GetMulti(c, ks, rs)
	if err == nil {
		return nil, err
	}

	var (
		merr     datastore.MultiError
		ok, enil bool
	)

	merr, ok = err.(datastore.MultiError)
	if !ok {
		return nil, err
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
		return rs, nil
	}
	return nil, err
}

func GetHistory(c *gin.Context, uKey *datastore.Key, t gType.Type) ([]*Rating, error) {
	ratings := make([]*Rating, len(gType.Types))
	keys := make([]*datastore.Key, len(gType.Types))
	for i, t := range gType.Types {
		ratings[i] = New(uKey, t)
		keys[i] = ratings[i].Key
	}

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return nil, err
	}

	err = dsClient.GetMulti(c, keys, ratings)
	if err != nil {
		return nil, err
	}
	return ratings, nil
}

func GetFor(c *gin.Context, t gType.Type) ([]*CurrentRating, error) {
	q := datastore.NewQuery(crKind).Ancestor(user.RootKey()).Filter("Type=", t).Order("-Low").KeysOnly()

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return nil, err
	}

	ks, err := dsClient.GetAll(c, q, nil)
	if err != nil {
		return nil, err
	}

	length := len(ks)
	ratings := make([]*CurrentRating, length)
	for i := range ratings {
		ratings[i] = NewCurrent(nil, t)
	}

	err = dsClient.GetMulti(c, ks, ratings)
	if err != nil {
		return nil, err
	}
	return ratings, nil
}

func CurrentProjected(c *gin.Context, rs []*CurrentRating, cm contest.ContestMap) ([]*CurrentRating, error) {
	var err error
	ratings := make([]*CurrentRating, len(rs))
	for i, r := range rs {
		ratings[i], err = r.Projected(c, cm[r.Type])
		if err != nil {
			return nil, err
		}
	}
	return ratings, nil
}

func (r *CurrentRating) Projected(c *gin.Context, cs []contest.Contest) (*CurrentRating, error) {
	log.Debugf("Entering r.Projected")
	defer log.Debugf("Exiting r.Projected")

	l := len(cs)
	if l == 0 && r.generated {
		log.Debugf("rating: %#v generated: %v", r, r.generated)
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

	return NewCurrent(r.Key.Parent, r.Type, rating.R, rating.RD), nil
}

func (r *CurrentRating) Generated() bool {
	return r.generated
}

func Index(c *gin.Context) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	cu, found := user.Current(c)
	t := gType.ToType[c.Param("type")]
	c.HTML(http.StatusOK, "rating/index", gin.H{
		"Type":      t,
		"Heading":   "Ratings: " + t.String(),
		"Types":     gType.Types,
		"Context":   c,
		"VersionID": appengine.VersionID(c),
		"CUser":     cu,
		"CUFound":   found,
	})
}

func getAllQuery(c *gin.Context) *datastore.Query {
	return datastore.NewQuery(crKind).Ancestor(user.RootKey())
}

func getFiltered(c *gin.Context, t gType.Type, leader bool, offset, limit int) ([]*CurrentRating, int, error) {
	q := getAllQuery(c)

	if leader {
		q = q.Filter("Leader=", true)
	}

	if t != gType.NoType {
		q = q.Filter("Type=", t)
	}

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return nil, 0, err
	}

	count, err := dsClient.Count(c, q)
	if err != nil {
		return nil, 0, err
	}

	q = q.Offset(offset).Limit(limit).Order("-Low").KeysOnly()

	ks, err := dsClient.GetAll(c, q, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("rating#getFiltered q.GetAll error: %v", err)
	}

	rs := make([]*CurrentRating, len(ks))
	ks2 := make([]*datastore.Key, len(ks))
	for i, k := range ks {
		rs[i] = NewCurrent(k.Parent, t)
		ks2[i] = rs[i].Key
	}

	err = dsClient.GetMulti(c, ks2, rs)
	if err != nil {
		return nil, 0, fmt.Errorf("rating#getFiltered datastore.Get error: %v", err)
	}
	return rs, count, err
}

func getUsers(c *gin.Context, rs []*CurrentRating) ([]user.User, error) {
	log.Debugf("Entering getUsers")
	defer log.Debugf("Exiting getUsers")

	log.Debugf("rs: %#v", rs)
	us := make([]user.User, len(rs))
	ks := make([]*datastore.Key, len(rs))
	for i := range rs {
		log.Debugf("rs[i]: %#v", rs[i])
		us[i] = user.New(rs[i].Key.Parent.ID)
		ks[i] = rs[i].Key.Parent
	}

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return nil, err
	}

	err = dsClient.GetMulti(c, ks, us)
	if err != nil {
		log.Debugf("Get error: %v", err)
		return nil, err
	}
	return us, nil
}

func getProjected(c *gin.Context, rs []*CurrentRating) ([]*CurrentRating, error) {
	log.Debugf("Entering getProjected")
	defer log.Debugf("Exiting getProjected")

	ps := make([]*CurrentRating, len(rs))
	for i, r := range rs {
		uKey := r.Key.Parent

		cs, err := contest.UnappliedFor(c, uKey, r.Type)
		if err != nil {
			return nil, err
		}

		ps[i], err = r.Projected(c, cs)
		if err != nil {
			return nil, err
		}

		if r.generated && r.generated && len(cs) == 0 {
			ps[i].generated = true
		}
	}
	return ps, nil
}

func For(c *gin.Context, u user.User, t gType.Type) (*CurrentRating, error) {
	return Get(c, u.Key, t)
}

func MultiFor(c *gin.Context, u user.User) ([]*CurrentRating, error) {
	return GetAll(c, u.Key)
}

// AddMulti has a limit of 100 tasks.  Thus, the batching.
func Update(c *gin.Context) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	var tk *taskqueue.Task

	tp := c.Param("type")
	q := user.AllQuery().KeysOnly()
	path := "/ratings/userUpdate"
	ts := make([]*taskqueue.Task, 0, 100)
	o := taskqueue.RetryOptions{RetryLimit: 5}

	tparams := make(url.Values)
	tparams.Set("type", tp)

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		log.Errorf(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	it := dsClient.Run(c, q)
	for {
		var u *user.User
		k, err := it.Next(&u)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Errorf(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		log.Debugf("k: %s", k)
		tparams.Set("uid", fmt.Sprintf("%v", u.ID()))
		tk = taskqueue.NewPOSTTask(path, tparams)
		tk.RetryOptions = &o
		ts = append(ts, tk)
	}

	_, err = taskqueue.AddMulti(c, ts, "")
	if err != nil {
		log.Errorf(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func updateUser(c *gin.Context) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	id, err := strconv.ParseInt(c.PostForm("uid"), 10, 64)
	if err != nil {
		log.Errorf(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		log.Errorf(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	u := user.New(id)
	err = dsClient.Get(c, u.Key, &u)
	if err != nil {
		log.Errorf("Unable to find user for id: %s", c.PostForm("uid"))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	t := gType.ToType[c.PostForm("type")]

	var r *CurrentRating
	if r, err = For(c, u, t); err != nil {
		log.Errorf("Unable to find rating for userid: %s", c.PostForm("uid"))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	cs, err := contest.UnappliedFor(c, u.Key, t)
	if err != nil {
		log.Errorf("Ratings update error when getting unapplied contests for user ID: %v.\n Error: %s", u.ID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return

	}

	var p *CurrentRating
	if p, err = r.Projected(c, cs); err != nil {
		log.Errorf("Ratings update error when getting projected rating for user ID: %v\n Error: %s", u.ID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if time.Since(time.Time(r.UpdatedAt)) < 504*time.Hour {
		log.Debugf("Did not update rating for user ID: %v", u.ID)
		log.Debugf("Rating updated %s ago.", time.Since(time.Time(r.UpdatedAt)))
		return
	}

	log.Debugf("Processing rating update for user ID: %v", u.ID)
	log.Debugf("Projected rating: %#v", p)
	log.Debugf("Unapplied contest count: %v", len(cs))

	es := make([]interface{}, 0)
	ks := make([]*datastore.Key, 0)

	log.Debugf("Reached for that ranges projected")
	const threshold float64 = 200.0
	// Update leader value to indicate whether present on leader boards
	p.Leader = p.RD < threshold

	if !p.Generated() {
		r := New(p.Key.Parent, p.Type, p.R, p.RD)
		es = append(es, p, r)
		ks = append(ks, p.Key, r.Key)
	}

	log.Debugf("Reached for that unapplied contests")
	for _, c := range cs {
		c.Applied = true
		es = append(es, c)
		ks = append(ks, c.Key)
	}

	_, err = dsClient.RunInTransaction(c, func(tx *datastore.Transaction) error {
		_, err = tx.PutMulti(ks, es)
		return err
	})

	if err != nil {
		log.Errorf("Ratings update err when saving updated ratings for user ID: %v\n Error: %s", u.ID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	log.Debugf("Reached RunInTransaction")
}

// func Fetch(c *gin.Context) {
// 	if CurrentRatingsFrom(c) != nil {
// 		return
// 	}
//
// 	u := user.Fetched(c)
// 	if u == nil {
// 		restful.AddErrorf(c, "Unable to get ratings.")
// 		c.Redirect(http.StatusSeeOther, homePath)
// 		return
// 	}
//
// 	if rs, err := MultiFor(c, u); err != nil {
// 		restful.AddErrorf(c, err.Error())
// 		c.Redirect(http.StatusSeeOther, homePath)
// 	} else {
// 		c.Set(currentRatingsKey, rs)
// 	}
// }
//
// func Fetched(c *gin.Context) []*CurrentRating {
// 	return CurrentRatingsFrom(c)
// }
//
// func FetchProjected(c *gin.Context) {
// 	if ProjectedFrom(c) != nil {
// 		return
// 	}
//
// 	rs := Fetched(c)
// 	if rs == nil {
// 		restful.AddErrorf(c, "Unable to get projected ratings")
// 		c.Redirect(http.StatusSeeOther, homePath)
// 		return
// 	}
//
// 	cm, err := contest.Unapplied(c, user.Fetched(c).Key)
// 	if err != nil {
// 		restful.AddErrorf(c, err.Error())
// 		c.Redirect(http.StatusSeeOther, homePath)
// 		return
// 	}
//
// 	if pr, err := CurrentProjected(c, rs, cm); err != nil {
// 		restful.AddErrorf(c, err.Error())
// 		c.Redirect(http.StatusSeeOther, homePath)
// 	} else {
// 		context.WithValue(c, projectedKey, pr)
// 	}
// }

func Projected(c *gin.Context) (pr []*Rating) {
	pr, _ = c.Value("Projected").([]*Rating)
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
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	var u user.User
	uid, err := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		log.Errorf("rating#JSONIndexAction BySID Error: %s", err)
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		log.Errorf(err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	u2 := user.New(uid)
	err = dsClient.Get(c, u2.Key, u2)
	if err != nil {
		log.Errorf("rating#JSONIndexAction unable to find user for uid: %d", uid)
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}
	u = u2

	rs, err := MultiFor(c, u)
	if err != nil {
		log.Errorf("rating#JSONIndexAction MultiFor Error: %s", err)
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	ps, err := getProjected(c, rs)
	if err != nil {
		log.Errorf("rating#getProjected Error: %s", err)
		c.Redirect(http.StatusSeeOther, homePath)
		return
	}

	if data, err := singleUser(c, u, rs, ps); err != nil {
		c.JSON(http.StatusOK, fmt.Sprintf("%v", err))
	} else {
		c.JSON(http.StatusOK, data)
	}
}

func JSONFilteredAction(c *gin.Context) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	t := gType.Get(c)

	var offset, limit int = 0, -1
	o, err := strconv.ParseInt(c.PostForm("start"), 10, 64)
	if err == nil && o >= 0 {
		offset = int(o)
	}

	l, err := strconv.ParseInt(c.PostForm("length"), 10, 64)
	if err == nil {
		limit = int(l)
	}

	rs, cnt, err := getFiltered(c, t, true, offset, limit)
	if err != nil {
		log.Errorf("rating#getFiltered Error: %s", err)
		return
	}
	log.Debugf("rs: %#v", rs)

	us, err := getUsers(c, rs)
	if err != nil {
		log.Errorf("rating#getUsers Error: %s", err)
		return
	}
	log.Debugf("us: %#v", us)

	ps, err := getProjected(c, rs)
	if err != nil {
		log.Errorf("rating#getProjected Error: %s", err)
		return
	}
	log.Debugf("ps: %#v", ps)

	data, err := toCombined(c, us, rs, ps, offset, cnt)
	if err != nil {
		log.Debugf("toCombined error: %v", err)
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

func singleUser(c *gin.Context, u user.User, rs, ps []*CurrentRating) (table *jCombinedRatingsIndex, err error) {
	log.Debugf("Entering singleUser")
	defer log.Debugf("Exiting singleUser")

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

	var draw int
	if draw, err = strconv.Atoi(c.PostForm("draw")); err != nil {
		log.Debugf("strconv.Atoi error: %v", err)
		return
	}

	table.Draw = draw
	table.RecordsTotal = l1
	table.RecordsFiltered = l2
	return
}

func toCombined(c *gin.Context, us []user.User, rs, ps []*CurrentRating, o, cnt int) (*jCombinedRatingsIndex, error) {
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

	if draw, err := strconv.Atoi(c.PostForm("draw")); err != nil {
		return nil, err
	} else {
		table.Draw = draw
	}

	table.RecordsTotal = int(cnt)
	table.RecordsFiltered = int(cnt)
	return table, nil
}

func IncreaseFor(c *gin.Context, u user.User, t gType.Type, cs []contest.Contest) (*CurrentRating, *CurrentRating, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	ucs, err := contest.UnappliedFor(c, u.Key, t)
	if err != nil {
		return nil, nil, err
	}

	r, err := For(c, u, t)
	if err != nil {
		return nil, nil, err
	}

	cr, err := r.Projected(c, ucs)
	if err != nil {
		return nil, nil, err
	}

	nr, err := r.Projected(c, append(ucs, filterContestsFor(cs, u.Key)...))
	return cr, nr, nil
}

func filterContestsFor(cs []contest.Contest, pk *datastore.Key) (fcs []contest.Contest) {
	for _, c := range cs {
		if c.Key.Parent.Equal(pk) {
			fcs = append(fcs, c)
		}
	}
	return
}
