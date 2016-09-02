package spoukfw

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"fmt"
	"strings"
//"runtime"
//"reflect"
	"log"
    
"runtime"
	"reflect"
)

type (
	spoukrouter struct {
		router      httprouter.Router
		//middlewares *spoukmiddlewares
		middlewares *spoukstockmiddlewares
		spoukmux    *Spoukmux
	}
	SpoukHandler func(*SpoukCarry) error
	spoukMapRoute struct {
		Path    string
		Method  string
		Handler string
	}
)

var (
	validMethods map[string]string = map[string]string{
		"GET" : "GET",
		"POST" : "POST",
		"UPDATE": "UPDATE",
		"DELETE" : "DELETE",
	}
)

func newSpoukRouter(mux *Spoukmux) *spoukrouter {
	sp := &spoukrouter{
		spoukmux:mux,
		router:httprouter.Router{HandleMethodNotAllowed:true, RedirectFixedPath:true, RedirectTrailingSlash:true},
	}
	return sp
}
//---------------------------------------------------------------------------
//  SPOUKROUTER
//---------------------------------------------------------------------------
//     router.ServeFiles("/src/*filepath", http.Dir("/var/www"))
func (sr *spoukrouter) StaticFiles(realpath, wwwpath string) {
	sr.router.ServeFiles(realpath, http.Dir(wwwpath))
}
func (sr spoukrouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sr.router.ServeHTTP(w, r)
}

func (sr *spoukrouter) fixPrefix(prefix string) string {
	//убирает слеш последний заебанный
	if len(strings.TrimSpace(prefix)) == 1 && prefix == "/" {
		return ""
	}
	if len(strings.TrimSpace(prefix)) > 1 {
		if prefix[len(prefix) - 1:] == "/" {
			return prefix[:len(prefix) - 1]
		}
	}
	return prefix
}
func (sr *spoukrouter) addRoute(method, path, prefix string, s SpoukHandler) {
	nameHandler := runtime.FuncForPC(reflect.ValueOf(s).Pointer()).Name()
	sr.spoukmux.RouteMapper[prefix + path] = spoukMapRoute{Path:prefix + path, Method:method, Handler:nameHandler}
	//prefixMiddle
	hu := sr.middlewares.getStockMiddlesPrefix(sr.fixPrefix(prefix)).wrappermidfunc(s)
	sr.router.Handle(strings.ToUpper(method), prefix + path, httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		newcarry := sr.spoukmux.getPool(w, r)
		newcarry.params = &ps
		//if sr.spoukmux.handlers != nil {
		//	sr.spoukmux.handlers.StockInit(newcarry)
		//}
		//sr.wrapperForStockInitInterface(hu)(newcarry)
		hu(newcarry)
		sr.spoukmux.putPool(newcarry)
	}))
}
func (sr *spoukrouter) checkMethod(method string) bool {
	if _, found := validMethods[method]; found {
		return true
	}
	return false
}
func (sr *spoukrouter) Multi(methods []string, path string, s SpoukHandler) {
	for _, m := range methods {
		if sr.checkMethod(m) {
			sr.addRoute(m, path, "", s)
		} else {
			err := makeErrorMessage(defRouterError, fmt.Sprintf("метод  `%s` не подходящий", m))
			log.Fatal(err.Error())
		}
	}
}
func (sr *spoukrouter) Get(path string, s SpoukHandler) {
	sr.addRoute("GET", path, "", s)
}
func (sr *spoukrouter) Post(path string, s SpoukHandler) {
	sr.addRoute("POST", path, "", s)
}
func (sr *spoukrouter) Delete(path string, s SpoukHandler) {
	sr.addRoute("DELETE", path, "", s)
}
func (sr *spoukrouter) Update(path string, s SpoukHandler) {
	sr.addRoute("UPDATE", path, "", s)
}
func (sr *spoukrouter) wrapperforSpoukHandler(h SpoukHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newmux := sr.spoukmux.getPool(w, r)
		h(newmux)
		sr.spoukmux.putPool(newmux)
	})
}
//func (sr *spoukrouter) SetMethodNotAllowed(h SpoukHandler) {
//	sr.router.MethodNotAllowed = sr.wrapperforSpoukHandler(h)
//}
//---------------------------------------------------------------------------
//  SPOUKHANDLER
//---------------------------------------------------------------------------
func (m SpoukHandler) ServeHTTPSpouker(s *SpoukCarry) {
	m(s)
}
//func error404spoukhandler(s *SpoukCarry) error {
func error404spoukhandler(s *SpoukCarry) error {
	fmt.Printf("[ERROR][CODE] %v:%v\n", s.request.Method, s.request.RequestURI)
	switch s.request.Method {
	case "POST", "DELETE", "UPDATE":
		s.WriteHTML(http.StatusMethodNotAllowed, fmt.Sprintf(notallowed405, s.request.Method))
	case "GET":
		s.WriteHTML(http.StatusNotFound, fmt.Sprintf(notFound404, s.request.RequestURI))
	}
	return nil
}
func error405methodNotAllowed(s *SpoukCarry) error {
	fmt.Printf("NOT ALLOWED METHOD\n")
	s.WriteHTML(http.StatusMethodNotAllowed, notallowed405)
	return nil
}
//func error405methodNotAllowedTest() http.Handler {
//	return httprouter.Handle(func(w http.ResponseWriter, r *http.Response, ps httprouter.Params) {
//		fmt.Printf("NOT ALLOWED METHOD\n")
//		w.Write([]byte("NOT ALLOWED METHOD\n"))
//		//s.WriteHTML(http.StatusMethodNotAllowed, notallowed405)
//	})
//}

func allow405(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("NOT ALLOWED METHOD\n")
	w.Write([]byte("NOT ALLOWED METHOD\n"))
}
//func error405methodNotAllowedTest() http.Handle {
//	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//		fmt.Printf("NOT ALLOWED METHOD\n")
//		w.Write([]byte("NOT ALLOWED METHOD\n"))
//		//s.WriteHTML(http.StatusMethodNotAllowed, notallowed405)
//	})
//}