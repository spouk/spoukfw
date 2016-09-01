package spoukfw

import (
	"time"
	"log"
	"fmt"
	"strings"
	//"runtime"
	//"reflect"
)

type (
//---------------------------------------------------------------------------
//  testing block
//---------------------------------------------------------------------------
	spoukstockmiddlewares  map[string]spoukmiddlewares//map[prefix_subdoimain][]MidFunc
	spoukmiddlewares []Midfunc
	Midfunc func(SpoukHandler) SpoukHandler
//testing othert format `spoukmiddlewares` = map[string]Midfunc
//	spoukmiddlewares map[string]Midfunc

)


//func newSpoukStockMiddlewares() *spoukstockmiddlewares {
//	return make(spoukstockmiddlewares)
//}
func (ss spoukstockmiddlewares) getStockMiddlesPrefix(prefix string) *spoukmiddlewares {
	//fmt.Printf("[middleware] `%v`: `%v` : `%v`\n", ss, prefix, ss[prefix])
	stock := ss[prefix]
	return &stock
}
//func (ss spoukstockmiddlewares) setStockMiddlePrefix(prefix string, m ...Midfunc) {
//	var tmp spoukmiddlewares
//	var stock spoukmiddlewares
//	var realprefix string
//
//	if len(strings.TrimSpace(prefix)) == 0 {
//
//		realprefix = ""
//		stock = ss[""]
//
//	} else {
//		realprefix = prefix
//		//получаю текущий по префиксу список
//		if oldstock, found := ss[prefix]; found {
//			//список есть
//			stock = oldstock
//
//		} else {
//			//префикс не найден
//			//копирую уже имеющиеся миддлы
//			def := ss[""]
//			mstock := make(spoukmiddlewares, len(def))
//			copy(mstock, def)
//			stock = mstock
//		}
//	}
//
//	//общая часть
//	//исключаю initStockMiddleware функцию
//	tmp = parseinitFuncMidleware(stock)
//	//добавляю новые к имеющимся
//	for _, x := range m {
//		nameHandler := runtime.FuncForPC(reflect.ValueOf(x).Pointer()).Name()
//		if !strings.Contains(nameHandler, "initFuncMidleware") {
//			tmp = append(tmp, x)
//		}
//	}
//	//добавляю в последнюю очередь мидд initStockMiddleware
//	tmp = append(tmp, initFuncMidleware)
//	//добавляю итог к префикусу в общую кучу
//	ss[realprefix] = tmp
//}

//func parseinitFuncMidleware(stock spoukmiddlewares) spoukmiddlewares {
//	var tmp spoukmiddlewares
//	//исключаю initStockMiddleware функцию
//	for _, x := range stock {
//		nameHandler := runtime.FuncForPC(reflect.ValueOf(x).Pointer()).Name()
//		if !strings.Contains(nameHandler, "initFuncMidleware") {
//			fmt.Printf("[parseinitFuncMidleware][append]  %v\n", nameHandler)
//			tmp = append(tmp, x)
//		}
//	}
//	return tmp
//}
func (sm spoukmiddlewares) wrappermidfunc(h SpoukHandler) SpoukHandler {
	for x := len(sm) - 1; x >= 0; x-- {
		h = sm[x](h)
	}
	return h
}
//---------------------------------------------------------------------------
//  default middlewares
//---------------------------------------------------------------------------
func loggerMiddleware(h SpoukHandler) SpoukHandler {

	fu := SpoukHandler(func(c *SpoukCarry) error{
		start := time.Now()
		h(c)
		log.Printf(infoLoggerApp, c.request.Method, fmt.Sprintf("`%v`:`%v`", c.request.RequestURI, time.Since(start)))
		return nil
	})
	return fu
}

type SessObj struct {
	Spoukmux *Spoukmux
	SpoukCarry  *SpoukCarry
}
func sessionMiddleware(h SpoukHandler) SpoukHandler {
	fu := SpoukHandler(func(c *SpoukCarry) error {
		//fmt.Printf("[sessionMiddleware] run...\n")

		s := c.spoukmux.session
		if s.SessionObject != nil {
			//start := time.Now()
			//создание нового динамического объекта для передачи по контексту
			newsessionobject := s.SessionObject.NewSpoukSessionObject(s,  c)
			//log.Printf("[sessionMiddleware][sessionobject] `%x` \n", &newsessionobject)
			c.Set("session", newsessionobject)
			//инициализация этого объекта
			s.SessionObject.InitSpoukSessionObject(newsessionobject)
			//log.Printf("RESULT TIMER: %v\b", time.Since(start))
		} else {
			log.Printf(infoLoggerApp, fmt.Sprintf("[spouksession] не найден `SpoukSessionObject`"))
		}
		//s := &SessObj{Spoukmux:c.spoukmux, SpoukCarry:c}
		//c.Set("sos", s)
		h(c)
		return nil
	})
	return fu
}
func initFuncMidleware(h SpoukHandler) SpoukHandler {
	fu := SpoukHandler(func(c *SpoukCarry) error {
		if c.spoukmux.handlers != nil {
			c.spoukmux.handlers.StockInit(c)
		}
		//fmt.Printf("[wrapperForStockInitInterface] run...\n")
		h(c)
		return nil
	})
	return fu
}
//---------------------------------------------------------------------------
//  простой вариант без оконечного миддла, префикс = список миддлов
//---------------------------------------------------------------------------
func (ss spoukstockmiddlewares) setStockMiddlePrefix(prefix string, m ...Midfunc) {

	//var tmp spoukmiddlewares
	var stock spoukmiddlewares
	var realprefix string

	//отсутствие префикса = основной сайт
	if len(strings.TrimSpace(prefix)) == 0 {

		realprefix = ""
		stock = ss[""]

	} else {
		//есть поддомен
		realprefix = prefix
		//получаю текущий по префиксу список
		if oldstock, found := ss[prefix]; found {
			//список есть
			stock = oldstock

		} else {
			//префикс не найден
			//копирую уже имеющиеся миддлы
			def := ss[""]
			mstock := make(spoukmiddlewares, len(def))
			copy(mstock, def)
			stock = mstock
		}
	}

	//общая часть
	//исключаю initStockMiddleware функцию
	//tmp = parseinitFuncMidleware(stock)
	//добавляю новые к имеющимся
	for _, x := range m {
		stock = append(stock, x)
		//nameHandler := runtime.FuncForPC(reflect.ValueOf(x).Pointer()).Name()
		//if !strings.Contains(nameHandler, "initFuncMidleware") {
		//	tmp = append(tmp, x)
		//}
	}
	//добавляю в последнюю очередь мидд initStockMiddleware
	//tmp = append(tmp, initFuncMidleware)
	//добавляю итог к префикусу в общую кучу
	ss[realprefix] = stock
}
