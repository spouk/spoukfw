package spoukfw

import (
	"sync"
	"net/http"
	"log"
	"fmt"
	"runtime"
	"reflect"
)

const (
	prefixSpoukMix = "[spoukmix][logger]"
	defaultErrorMsgSpoukMux = "[spoukmux][error] `%v`"
//errors spoukmux
	errorSessionObject = "не найден объект сессии, надо подключить"
)

type (
	Spoukmux struct {
		config      *Spoukconfig
		router      *spoukrouter
		pool        sync.Pool
		middlewares spoukstockmiddlewares
		render      SpoukMuxRendering
		logger      *SpoukLogger
		RouteMapper map[string]spoukMapRoute
		session     *SpoukSession
		handlers    SpoukHandlerStock
	}
//интерфейс для рендера
	SpoukMuxRendering  interface {
		SpoukRenderIO(name string, data interface{}, resp http.ResponseWriter, reloadTemplate bool) (err error)
		AddUserFilter(name string, f interface{})
		AddFilters(stack map[string]interface{})
	}
//интерфейс для стека обработчиков
	SpoukHandlerStock interface {
		StockInit(s *SpoukCarry, args ...interface{})
	}

)

func (m *Spoukmux) SetStockHandlerStock(stock SpoukHandlerStock) {
	m.handlers = stock
}
func (m *Spoukmux) ShowRoutingMap() {
	for _, x := range m.RouteMapper {
		fmt.Printf("[spoukmux] [%7s] `%20s`  `%+v`\n", x.Method, x.Path, x.Handler)
	}
}

func (m *Spoukmux) SetRender(r SpoukMuxRendering) {
	m.render = r
}
func (m *Spoukmux) Session() *SpoukSession {
	return m.session
}

func (m Spoukmux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.router.ServeHTTP(w, r)
}
func (m *Spoukmux) getPool(w http.ResponseWriter, r *http.Request) (res *SpoukCarry) {
	newcarry := m.pool.Get().(*SpoukCarry)
	newcarry.response = &w
	newcarry.request = r.WithContext(r.Context())
	return newcarry
}
func (m *Spoukmux) putPool(res *SpoukCarry) {
	m.pool.Put(res)
}
func (m *Spoukmux) StaticFiles(realpath, wwwpath string) {
	m.router.StaticFiles(realpath, wwwpath)
}
//router.ServeFiles("/src/*filepath", http.Dir("/var/www"))

func (s *Spoukmux) catchErrors() {
	log.Printf("[spoukmux][catcher-errros]\n")
}
func (s *Spoukmux) Run() {

	defer s.catchErrors()
	srv := http.Server{
		Addr:         s.config.Address,
		Handler:      s,
		ReadTimeout:  s.config.HTTPReadTimeout,
		WriteTimeout: s.config.HTTPWriteTimeout,
	}
	log.Printf(fmt.Sprintf("[middlewares]%v\n[session] `%v`\n[sessionobject] `%v`\n", s.middlewares, s.session))
	log.Printf(fmt.Sprintf(runInfo, s.config.Address))
	log.Fatal(srv.ListenAndServe())
}

func (s *Spoukmux) AddMiddleware(mid Midfunc) {
	s.middlewares.setStockMiddlePrefix("", mid)
}
func (s *Spoukmux) ShowConfig() {
	fmt.Printf("[spoukmux]ConfigProxy: %v\nSpoukconfig: %v\n ", s.config.ConfigProxy, s.config)
}

// user filter map[commandintemplate]func closure
func (s *Spoukmux) AddUserFlitersRender(userfilters  map[string]interface{}) {
	s.render.AddFilters(userfilters)
}
func (s *Spoukmux) AddUserFilterRender(name string, f interface{}) {
	s.render.AddUserFilter(name, f)
}
func NewSpoukmux(spcfg *Spoukconfig) *Spoukmux {
	m := Spoukmux{
		pool:   sync.Pool{},
		config:spcfg,
		router:&spoukrouter{},
	}
	m.middlewares = make(spoukstockmiddlewares)
	m.middlewares.setStockMiddlePrefix("", loggerMiddleware)
	m.router.middlewares = &m.middlewares

	m.RouteMapper = make(map[string]spoukMapRoute)
	if m.config.UseSession {
		m.session = newSpoukSession(&m)
		m.middlewares.setStockMiddlePrefix("", sessionMiddleware)
	}
	//m.middlewares.setStockMiddlePrefix("", initFuncMidleware)
	m.logger = NewSpoukLogger(prefixSpoukMix, nil)
	m.router.router.HandleMethodNotAllowed = false
	m.router.router.NotFound = m.router.wrapperforSpoukHandler(error404spoukhandler)
	//m.SetMehodNotAllowedHandler(error405methodNotAllowed)
	m.router.router.MethodNotAllowed = http.HandlerFunc(allow405)

	m.router.spoukmux = &m
	m.router.router.RedirectTrailingSlash = true
	m.pool.New = func() interface{} {
		return &SpoukCarry{
			spoukmux: &m,
			request: &http.Request{},
		}
	}

	//rendering default
	m.render = NewSpoukRender(m.config.TemplatePath, m.config.TemplateDebug)

	//fmt.Printf("[spoukmux] middlewaresstock: %v\n", m.middlewares)
	//m.ShowMiddlewares()
	return &m
}

func (s *Spoukmux) ShowMiddlewares() {
	for prefix, stock := range s.middlewares {
		fmt.Printf("[spoukmux][`%s`] \n", prefix)
		for _, x := range stock {
			nameHandler := runtime.FuncForPC(reflect.ValueOf(x).Pointer()).Name()
			fmt.Printf("[spoukmux][`%s`]  `%v`\n", prefix, nameHandler)
		}
	}
}
func (s *Spoukmux) Multi(methods []string, path string, h SpoukHandler) {
	s.router.Multi(methods, path, h)
}
func (s *Spoukmux) Get(path string, h SpoukHandler) {
	s.router.Get(path, h)
}
func (s *Spoukmux) Post(path string, h SpoukHandler) {
	s.router.Post(path, h)
}
func (s *Spoukmux) Config() *Spoukconfig {
	return s.config
}
func (s *Spoukmux) SetSessionObject(so SpoukSessionObject) {
	if s.session != nil {
		s.session.SessionObject = so
	} else {
		s.logger.Printf(makeErrorMessage(defaultErrorMsgSpoukMux, errorSessionObject).Error())
	}
}
func (s *Spoukmux ) SetMehodNotAllowedHandler(h  SpoukHandler) {
	s.router.router.MethodNotAllowed = s.router.wrapperforSpoukHandler(h)
}
//---------------------------------------------------------------------------
//  testing subroute
//---------------------------------------------------------------------------
func (s *Spoukmux) Subroute(prefix string) *spoukSubdomain {
	newss := &spoukSubdomain{spoukmux:*s, prefix:prefix}
	newss.prefix = prefix
	return newss
}

type spoukSubdomain struct {
	spoukmux Spoukmux
	prefix   string
}
func (ss *spoukSubdomain) Multi(methods []string, path string, s SpoukHandler) {
	for _, m := range methods {
		if ss.spoukmux.router.checkMethod(m) {
			ss.spoukmux.router.addRoute(m, path, ss.prefix, s)
		} else {
			err := makeErrorMessage(defRouterError, fmt.Sprintf("метод  `%s` не подходящий", m))
			log.Fatal(err.Error())
		}
	}
}
func (ss *spoukSubdomain) AddMiddleware(mid Midfunc) {
	ss.spoukmux.middlewares.setStockMiddlePrefix(ss.prefix, mid)
	fmt.Printf("[spoukSubdomain][AddMiddleware] %v\n", ss.spoukmux.middlewares)
}
func (ss *spoukSubdomain) Get(path string, h SpoukHandler) {
	ss.spoukmux.router.addRoute("GET", path, ss.prefix, h)
}
func (ss *spoukSubdomain) Post(path string, h SpoukHandler) {
	ss.spoukmux.router.addRoute("POST", path, ss.prefix, h)
}
func (ss *spoukSubdomain) Delete(path string, h SpoukHandler) {
	ss.spoukmux.router.addRoute("DELETE", path, ss.prefix, h)
}
func (ss *spoukSubdomain) Update(path string, h SpoukHandler) {
	ss.spoukmux.router.addRoute("UPDATE", path, ss.prefix, h)
}