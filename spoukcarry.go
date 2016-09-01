package spoukfw

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/url"
	"time"
	"strings"
	"fmt"
)

type (
	SpoukCarry struct {
		response *http.ResponseWriter
		request  *http.Request
		spoukmux *Spoukmux
		params   *httprouter.Params
		Spoukmux *Spoukmux
	}
)

func (m *SpoukCarry) Config() *Spoukconfig {
	return m.spoukmux.config
}
func (m *SpoukCarry) Render(name string, data interface{}) error{
	m.spoukmux.render.SpoukRenderIO(name, data, *m.response, m.spoukmux.config.TemplateDebug)
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
func (sr *SpoukCarry) WriteHTML(httpcode int, text string) {
	resp := *sr.response
	resp.Header().Set(ContentType, TextHTMLCharsetUTF8)
	resp.WriteHeader(httpcode)
	resp.Write([]byte(text))
	//return nil
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
	fmt.Sprintf("[middlewares]%v\n[session] `%v`\n[sessionobject] `%v`\n", s.spoukmux.middlewares, s.spoukmux.session, s.spoukmux.session.SessionObject)
}