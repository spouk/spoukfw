package spoukfw

import (
	"fmt"
	"strings"
	"strconv"
	"github.com/labstack/gommon/log"
	"github.com/labstack/echo"
	"math"
	//dbs "spoukframework/spoukfw/dbs"
)
type (
	SpoukPaginate struct {
		Result *SpoukPaginateResult
		//Dbs *dbs.Database
		spoukmux *Spoukmux
	}
	SpoukPaginateResult struct {
		Total  interface{}
		Result interface{}
		Links  string
	}
)
func NewSpoukPaginate(s *Spoukmux) *SpoukPaginate {
	return &SpoukPaginate{
		spoukmux:s,
	}
}
//func (s *SpoukPaginate) Paginate(typeTable interface{}, page int64, path string) (*SpoukPaginateResult) {
//	countpage := s.spoukmux.config.CountOnPage
//	countlinks  := s.spoukmux.config.CountLinks
//	if answer := s.Dbs.Paginate(typeTable ,page, int(countpage)); answer.Status == false {
//		return nil
//	} else {
//		links := s.PaginateHTML(int(page), int(countpage), int(countlinks), path, answer.CountTotal)
//		return &SpoukPaginateResult{Total:answer.Total,Result:answer.Result,Links:links}
//	}
//}
func (s *SpoukPaginate) PaginateHTML(current_page int, count_on_page int, count_links int, path string, total_len int) string {
	// start,middle,end; middle = stacklinks, start,end = управляющие кнопки
	if current_page == 0 {
		current_page = 1
	}
	current_page = int(current_page)
	var tmp string = `<ul class="pagination"> %s </ul>`
	var START string = `<li class="%s"><a href="%s/page/%d" %s> < </a></li>`
	var END string = `<li class="%s"><a href="%s/page/%d" %s> > </a></li>`
	var LINK string = `	<li><a class="%s" href="%s/page/%d"><span> %d</a></li>`
	var LINK_ACTIVE string = `<li class="%s"><a href="%s/page/%d"><span> %d</a></li>`
	var out string
	var stacklinks string
	var start, end string
	var plinks []string
	//общее количество страниц исходя из количества элементов на странице
	totalpages := int(math.Ceil(float64(total_len) / float64(count_on_page)))

	for i := 1; i <= totalpages; i++ {
		if i == current_page {
			plinks = append(plinks, fmt.Sprintf(LINK_ACTIVE, "active", path, i, i))
		} else {
			plinks = append(plinks, fmt.Sprintf(LINK, "", path, i, i))
		}
	}
	// если общее количество страниц `num_pages` <= количеству страниц показывемых в пагинации
	//start,end страницы выставляются в disabled
	if totalpages <= count_links {
		if current_page == 1 {
			start = fmt.Sprintf(START, "disabled", path, current_page, "disabled")
		} else {
			start = fmt.Sprintf(START, "", path, current_page - 1, "")
		}

		if current_page < totalpages {
			end = fmt.Sprintf(END, "", path, current_page + 1, "")
		} else {
			end = fmt.Sprintf(END, "disabled", path, current_page, "disabled")
		}
	}
	if totalpages > count_links {
		if totalpages - current_page >= count_links {
			plinks = plinks[current_page - 1: current_page + count_links - 1]
			fmt.Printf("-- total > countlinks\n", plinks)
		} else {
			fmt.Printf("-- total < countlinks\n", plinks)
			plinks = plinks[totalpages - count_links:]
		}
		if current_page == 1 {
			start = fmt.Sprintf(START, "disabled", path, current_page, "disabled")
			end = fmt.Sprintf(END, "", path, current_page + 1, "")
		}
		if current_page > 1 {
			start = fmt.Sprintf(START, "", path, current_page - 1, "")
			if totalpages == current_page {
				end = fmt.Sprintf(END, "disabled", path, current_page + 1, "disabled")
			} else {
				end = fmt.Sprintf(END, "", path, current_page + 1, "")
			}
		}
	}
	//compose
	stacklinks += start
	stacklinks += strings.Join(plinks, "")
	stacklinks += end
	out = fmt.Sprintf(tmp, stacklinks)
	return out
}

func GetPage(c *echo.Context) (int) {
	var page int = 1
	q := c.Query("page")
	if q != "" {
		p, err := strconv.Atoi(q)
		if err != nil {
			log.Error(err.Error())
		}
		page = p
	}
	return page
}
func GetBoolForm(nameElement string, c *echo.Context) bool {
	active := c.Form(nameElement)
	if len(active) != 0 {
		return true
	} else {
		return false
	}
}
