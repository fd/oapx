package handlers

import (
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"code.google.com/p/goauth2/oauth"
	"github.com/RangelReale/osin"
	"github.com/codegangsta/martini"
	"github.com/jmoiron/sqlx"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"

	"github.com/fd/oauth2-proxy/data"
	"github.com/fd/oauth2-proxy/providers"
)

type FlowType uint8

const (
	LoginFlow FlowType = 1 + iota
	LinkFlow
	AuthorizeFlow
)

type FlowState struct {
	Type       FlowType
	Source     string
	IdentityId int64
	State      string
	Provider   string
	StartAt    time.Time
}

func init() {
	gob.RegisterName("flow", FlowState{})
}

func GET_login(c martini.Context, sess sessions.Session, r *http.Request) {
	var (
		identity = ActiveIdentity(c)
		source   = r.Referer()
		handler  martini.Handler
	)

	if identity != nil {
		sess.Delete("flow")
		handler = redirect_to(source)
	} else {
		sess.Set("flow", FlowState{
			Type:    LoginFlow,
			Source:  source,
			StartAt: time.Now(),
		})
		handler = show_provider_chooser()
	}

	c.Invoke(handler)
}

func GET_link(c martini.Context, sess sessions.Session, r *http.Request) {
	var (
		identity = ActiveIdentity(c)
		source   = r.Referer()
		handler  martini.Handler
	)

	if identity == nil {
		sess.Delete("flow")
		handler = forbidden()
	} else {
		sess.Set("flow", FlowState{
			Type:       LinkFlow,
			Source:     source,
			IdentityId: identity.Id,
			StartAt:    time.Now(),
		})
		handler = show_provider_chooser()
	}

	c.Invoke(handler)
}

func GET_authorize(c martini.Context, sess sessions.Session, w http.ResponseWriter, r *http.Request, s *osin.Server) {
	resp := s.NewResponse()
	if ar := s.HandleAuthorizeRequest(resp, r); ar != nil {
		if !inner_GET_authorize(c, sess, r, ar) {
			return
		}

		ar.Authorized = true
		s.FinishAuthorizeRequest(resp, r, ar)
	}
	osin.OutputJSON(resp, w, r)
}

func inner_GET_authorize(c martini.Context, sess sessions.Session, r *http.Request, ar *osin.AuthorizeRequest) bool {
	var (
		identity = ActiveIdentity(c)
		source   = current_url(r)
		handler  martini.Handler
	)

	if identity != nil {
		ar.UserData = identity
		sess.Delete("flow")
		return true
	} else {
		sess.Set("flow", FlowState{
			Type:    AuthorizeFlow,
			Source:  source,
			StartAt: time.Now(),
		})

		if provider := r.URL.Query().Get("p"); provider == "" {
			handler = show_provider_chooser()
		} else {
			handler = redirect_to_provider(provider)
		}
	}

	c.Invoke(handler)
	return false
}

func GET_token(w http.ResponseWriter, r *http.Request, s *osin.Server) {
	resp := s.NewResponse()
	if ar := s.HandleAccessRequest(resp, r); ar != nil {
		// always true
		ar.Authorized = true
		s.FinishAccessRequest(resp, r, ar)
	}
	osin.OutputJSON(resp, w, r)
}

func GET_info(w http.ResponseWriter, r *http.Request, s *osin.Server) {
	resp := s.NewResponse()
	if ir := s.HandleInfoRequest(resp, r); ir != nil {
		s.FinishInfoRequest(resp, r, ir)
	}
	osin.OutputJSON(resp, w, r)
}

func GET_continue(c martini.Context, params martini.Params) {
	var (
		provider = params["provider"]
		handler  martini.Handler
	)

	if provider == "" {
		handler = show_provider_chooser()
	} else {
		handler = redirect_to_provider(provider)
	}

	c.Invoke(handler)
}

func GET_callback(c martini.Context, sess sessions.Session, r *http.Request, db *sqlx.DB) {
	flow, ok := sess.Get("flow").(FlowState)
	if !ok {
		c.Invoke(redirect_to("/login"))
		return
	}
	if flow.StartAt.Before(time.Now().Add(-10 * time.Minute)) {
		c.Invoke(redirect_to("/login"))
		return
	}
	if flow.State == "" {
		c.Invoke(redirect_to("/login"))
		return
	}
	if r.URL.Query().Get("code") == "" {
		c.Invoke(redirect_to("/login"))
		return
	}
	if flow.State != r.URL.Query().Get("state") {
		c.Invoke(redirect_to("/login"))
		return
	}

	var (
		provider  = tmp_new_provider(flow.Provider)
		transport = provider.Transport(nil)
	)

	token, err := transport.Exchange(r.URL.Query().Get("code"))
	if err != nil {
		panic(err)
	}

	profile, err := provider.GetProfile(transport)
	if err != nil {
		panic(err)
	}

	var (
		tx      = db.MustBegin()
		success bool
	)

	defer func() {
		if success {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	account, err := data.GetAccountWithRemoteId(tx, profile.RemoteId())
	if err != nil {
		panic(err)
	}

	c.MapTo(profile, (*providers.Profile)(nil))
	c.Map(token)
	c.Map(tx)
	c.Map(account)

	if account != nil {
		c.Invoke(GET_callback_A)
	} else {
		c.Invoke(GET_callback_B)
	}

	success = true
}

func GET_callback_A(c martini.Context, sess sessions.Session) {
	flow := sess.Get("flow").(FlowState)

	switch flow.Type {
	case LoginFlow:
		c.Invoke(GET_callback_AA)
	case LinkFlow:
		c.Invoke(GET_callback_AB)
	case AuthorizeFlow:
		c.Invoke(GET_callback_AC)
	default:
		panic("unknown flow type")
	}
}

func GET_callback_B(c martini.Context, sess sessions.Session) {
	flow := sess.Get("flow").(FlowState)

	switch flow.Type {
	case LoginFlow:
		c.Invoke(GET_callback_BA)
	case LinkFlow:
		c.Invoke(GET_callback_BB)
	case AuthorizeFlow:
		c.Invoke(GET_callback_BC)
	default:
		panic("unknown flow type")
	}
}

func GET_callback_AA(c martini.Context, sess sessions.Session) {
	flow := sess.Get("flow").(FlowState)

	c.Invoke(update_account)
	c.Invoke(activate_session)
	c.Invoke(redirect_to(flow.Source))
}

func GET_callback_AB(c martini.Context, sess sessions.Session) {
	flow := sess.Get("flow").(FlowState)

	c.Invoke(match_session_identity_with_account)
	c.Invoke(match_session_identity_with_flow)
	c.Invoke(update_account)
	c.Invoke(redirect_to(flow.Source))
}

func GET_callback_AC(c martini.Context, sess sessions.Session) {
	flow := sess.Get("flow").(FlowState)

	c.Invoke(update_account)
	c.Invoke(activate_session)
	c.Invoke(redirect_to(flow.Source))
}

func GET_callback_BA(c martini.Context, sess sessions.Session) {
	flow := sess.Get("flow").(FlowState)

	c.Invoke(create_identity)
	c.Invoke(create_account)
	c.Invoke(activate_session)
	c.Invoke(redirect_to(flow.Source))
}

func GET_callback_BB(c martini.Context, sess sessions.Session) {
	flow := sess.Get("flow").(FlowState)

	c.Invoke(match_session_identity_with_flow)
	c.Invoke(create_account)
	c.Invoke(redirect_to(flow.Source))
}

func GET_callback_BC(c martini.Context, sess sessions.Session) {
	flow := sess.Get("flow").(FlowState)

	c.Invoke(create_identity)
	c.Invoke(create_account)
	c.Invoke(activate_session)
	c.Invoke(redirect_to(flow.Source))
}

func match_session_identity_with_account(c martini.Context, account *data.Account) {
	identity := ActiveIdentity(c)
	if identity == nil {
		panic("account belongs to a different identity (1)")
	}

	if account.IdentityId != identity.Id {
		panic("account belongs to a different identity (2)")
	}
}

func match_session_identity_with_flow(c martini.Context, sess sessions.Session) {
	identity := ActiveIdentity(c)
	flow := sess.Get("flow").(FlowState)

	if identity == nil {
		panic("account belongs to a different identity (3)")
	}

	if flow.IdentityId != identity.Id {
		panic("account belongs to a different identity (4)")
	}
}

func update_account(c martini.Context, tx *sqlx.Tx,
	account *data.Account, profile providers.Profile, token *oauth.Token) {

	account = &data.Account{
		Id:         account.Id,
		IdentityId: account.IdentityId,
		RemoteId:   profile.RemoteId(),
		Name:       profile.Name(),
		Email:      profile.Email(),
		Picture:    profile.PictureURL(),
		RawProfile: profile.RawData(),
		RawToken:   encode_token(token),
	}

	err := data.UpdateAccount(tx, account)
	if err != nil {
		panic(err)
	}
}

func create_account(c martini.Context, tx *sqlx.Tx,
	identity *data.Identity, profile providers.Profile, token *oauth.Token) {

	account := &data.Account{
		IdentityId: identity.Id,
		RemoteId:   profile.RemoteId(),
		Name:       profile.Name(),
		Email:      profile.Email(),
		Picture:    profile.PictureURL(),
		RawProfile: profile.RawData(),
		RawToken:   encode_token(token),
	}
	err := data.CreateAccount(tx, account)
	if err != nil {
		panic(err)
	}
}

func create_identity(c martini.Context, tx *sqlx.Tx) {
	identity := &data.Identity{}
	err := data.CreateIdentity(tx, identity)
	if err != nil {
		panic(err)
	}

	c.Map(identity)
}

func activate_session(sess sessions.Session, account *data.Account) {
	sess.Set("identity_id", account.IdentityId)
}

// TODO: Handle https
func current_url(r *http.Request) string {
	return "http://" + r.Host + r.RequestURI
}

func forbidden() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func redirect_to(target string) martini.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target, http.StatusFound)
	}
}

func redirect_to_provider(provider_code string) martini.Handler {
	return func(c martini.Context, sess sessions.Session) {
		var (
			flow = sess.Get("flow").(FlowState)
			err  error
		)

		provider := tmp_new_provider(provider_code)
		if provider == nil {
			panic("unknown provider")
		}

		flow.Provider = provider_code
		flow.State, err = make_state()
		if err != nil {
			panic(err)
		}

		sess.Set("flow", flow)

		target := provider.Config().AuthCodeURL(flow.State)
		c.Invoke(redirect_to(target))
	}
}

func show_provider_chooser() martini.Handler {
	return func(render render.Render) {
		render.HTML(200, "provider_chooser", nil)
	}
}

func encode_token(token *oauth.Token) []byte {
	data, err := json.Marshal(token)
	if err != nil {
		panic(err)
	}
	return data
}

func make_state() (string, error) {
	buf := make([]byte, 20)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}

func tmp_new_provider(typ string) *providers.Provider {
	var (
		p   *providers.Provider
		err error
	)

	switch typ {
	case "go":
		p, err = providers.New(typ, "536228200101.apps.googleusercontent.com", "TvF9KQB6iMMEjj5oagTEsdqY")
	case "gh":
		p, err = providers.New(typ, "75d06562f297cf2cc0a5", "00e3c788d083c12678a29bf87e0d374c6a1d5bc2")
	case "fb":
		p, err = providers.New(typ, "303274556491300", "4bdf767541d8fab54f149ebceb0b114c")
	default:
		p, err = providers.New(typ, "", "")
	}

	if err != nil {
		panic(err)
	}

	return p
}
