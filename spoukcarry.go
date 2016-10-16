package spoukfw

import (
	"net/http"
	"github.com/spouk/spoukfw/httprouter"
	"golang.org/x/net/context"
	"net/url"
	"time"
	"strings"
	"fmt"
	"encoding/json"
)

type (
	SpoukCarry struct {
		response *http.ResponseWriter
		request  *http.Request
		params   *httprouter.Params
		Spoukmux *Spoukmux
	}
)

func (m *SpoukCarry) RealPath() string {
	_, rp, _, _ := m.Spoukmux.router.router.LookupRoute(m.request.Method, m.request.URL.Path)
	return rp
}

func (m *SpoukCarry) Config() *Spoukconfig {
	return m.Spoukmux.config
}
func (m *SpoukCarry) Render(name string, data interface{}) error {
	m.Spoukmux.render.SpoukRenderIO(name, data, *m.response, m.Spoukmux.config.TemplateDebug)
	return nil
}
func (sr *SpoukCarry) Request() *http.Request {
	return sr.request
}
func (sr *SpoukCarry) Response() *http.ResponseWriter {
	return sr.response
}
func (sr *SpoukCarry) Params() *httprouter.Params {
	return sr.params
}
func (sr *SpoukCarry) Redirect(path string) error {
	http.Redirect(*sr.response, sr.request, path, http.StatusFound)
	return nil
}
func (sr *SpoukCarry) WriteHTML(httpcode int, text string) (error) {
	resp := *sr.response
	resp.Header().Set(ContentType, TextHTMLCharsetUTF8)
	resp.WriteHeader(httpcode)
	resp.Write([]byte(text))
	return nil
}
func (sr *SpoukCarry) JSONB(httpcode int, b []byte) (error) {
	resp := *sr.response
	resp.Header().Set(ContentType, ApplicationJavaScriptCharsetUTF8)
	resp.WriteHeader(httpcode)
	resp.Write(b)
	return nil
}
func (sr *SpoukCarry) JSON(code int, answer interface{}) (err error) {
	b, err := json.Marshal(answer)
	if err != nil {
		return err
	}
	return sr.JSONB(code, b)
}
func (sr *SpoukCarry) Set(key string, value interface{}) {
	ctx := context.WithValue(sr.request.Context(), key, value)
	sr.request = sr.request.WithContext(ctx)
}
func (sr *SpoukCarry) Get(key string) (value interface{}) {
	return sr.request.Context().Value(key)
}
func (sr *SpoukCarry) Path() (value string) {
	return sr.request.URL.Path
}
func (sr *SpoukCarry) GetParam(key string) string {
	return sr.params.ByName(key)
}
func (sr *SpoukCarry) GetQueryAll(key string) url.Values {
	return sr.request.URL.Query()
}
func (sr *SpoukCarry) GetQuery(key string) string {
	return sr.request.URL.Query().Get(key)
}
func (sr *SpoukCarry) GetFormValue(key string) string {
	return sr.request.FormValue(key)
}
func (sr *SpoukCarry) GetForm(key string) string {
	return sr.request.PostFormValue(key)
}
func (sr *SpoukCarry) GetCooks() []*http.Cookie {
	return sr.request.Cookies()
}
func (sr *SpoukCarry) GetCook(nameCook string) *http.Cookie {
	for _, c := range sr.request.Cookies() {
		if nameCook == c.Name {
			return c
		}
	}
	return nil
}
func (sr *SpoukCarry) SetCook(c http.Cookie) {
	_newCook := new(http.Cookie)
	_newCook.Name = c.Name
	_newCook.Value = c.Value
	//проверка на domainName из конфига ибо может быть различный доменое имя или IP
	//типа: localhost, 127.0.0.1, 0.0.0.0 - комп один, сетевуха одна, а различие в алиасах есть отсюда разные кукисы
	//установка времени истечения срока действия печеньки

	if time.Now().Sub(c.Expires) > 0 {
		t := time.Now()
		_newCook.Expires = t.Add(time.Duration(86000 * 30) * time.Minute)
	}
	//установка домена
	if c.Domain == "" {
		domainReal := strings.Split(sr.request.Host, ":")[0]
		if "localhost" != domainReal {
			_newCook.Domain = domainReal
		} else {
			_newCook.Domain = "localhost"
		}
	} else {
		_newCook.Domain = c.Domain
	}
	//path
	if c.Path == "" {
		_newCook.Path = "/"
	}
	http.SetCookie(*sr.response, _newCook)
}
func (s *SpoukCarry) ShowMiddlwaresAndSession() {
	fmt.Sprintf("[middlewares]%v\n[session] `%v`\n[sessionobject] `%v`\n", s.Spoukmux.middlewares, s.Spoukmux.session, s.Spoukmux.session.SessionObject)
}