package models

import (
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-rock/rock"
	"github.com/jinzhu/gorm"
)

const (
	PAGESIZE = 12
	PAGE     = 1
)

type Pagination struct {
	TotalRecord int  `json:"total_record"`
	CurrentPage int  `json:"current_page"`
	PageSize    int  `json:"page_size"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
}
type DbError struct {
	Error   error  `json:"-"`
	Message string `json:"message"`
	Code    int    `json:"status_code"`
}

type condition struct {
	Field string
	Query string
	Args  interface{}
	IsOr  bool
}
type Repo struct {
	Ctx          rock.Context
	PageSize     int
	DB           *gorm.DB
	Result       interface{} //结果集，slice或struct
	Pagination   Pagination  //分页
	AutoResponse bool        //是否直接通过ctx输出JSON结果
	DisExpTag    bool        //是否禁止解析tag中的`exp`作为默认条件
	Error        DbError
	count        int               //结果统计
	IgnoreParam  bool              //是否忽略path上的参数，非query
	PathParamMap map[string]string //path 上的参数,如果IgnoreParam是 false,这个可以不填,格式[表中字段名]参数名
	ApplyWhere   bool
}

func isSlice(ret interface{}) bool {
	t := reflect.TypeOf(ret)
	k := t.Elem().Kind()
	return k.String() == "slice"
}
func genDbError(retErr error) (err DbError) {
	if gorm.IsRecordNotFoundError(retErr) {
		err = DbError{
			Error:   retErr,
			Message: retErr.Error(),
			Code:    http.StatusNotFound,
		}
	} else {
		err = DbError{
			Error:   retErr,
			Message: retErr.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	return
}
func parseModelTags(tags reflect.StructTag) map[string]string {
	setting := map[string]string{}
	for _, str := range []string{tags.Get("sql"), tags.Get("gorm")} {
		tags := strings.Split(str, ";")
		for _, value := range tags {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if len(v) >= 2 {
				setting[k] = strings.Join(v[1:], ":")
			} else {
				setting[k] = k
			}
		}
	}
	return setting
}
func (repo *Repo) Fetch() Repo {
	repo.preFetch()

	isSlice := isSlice(repo.Result)
	var ret *gorm.DB
	if isSlice {
		offset := repo.Pagination.PageSize * (repo.Pagination.CurrentPage - 1)
		limit := repo.Pagination.PageSize

		ret = repo.DB.Model(repo.Result).Count(&repo.count).Offset(offset).Limit(limit).Find(repo.Result)
	} else {
		ret = repo.DB.First(repo.Result)
	}

	if ret.Error != nil {
		repo.Error = genDbError(ret.Error)
	}

	repo.Pagination.TotalRecord = repo.count
	repo.Pagination.TotalPages = int(math.Ceil(float64(repo.count) / float64(repo.Pagination.PageSize)))
	repo.Pagination.HasNext = repo.Pagination.TotalPages > repo.Pagination.CurrentPage
	if repo.AutoResponse {
		repo.JsonResponse()
	}

	return *repo
}

func (repo Repo) JsonResponse() {
	if repo.Error.Error != nil {
		repo.Ctx.JSON(repo.Error.Code, repo.Error)
	} else {
		ret := make(map[string]interface{})
		ret["code"] = http.StatusOK
		ret["data"] = repo.Result
		if isSlice(repo.Result) {
			ret["pagination"] = repo.Pagination
		}

		repo.Ctx.JSON(200, ret)
	}
}

func (repo *Repo) preFetch() {

	if repo.ApplyWhere {
		repo.applyWhere()
	}
	// repo.applyBetween()
	repo.applyPage()
	repo.applyOrder()
	// repo.applyGroup()
	if !repo.IgnoreParam {
		for k, v := range repo.PathParamMap {
			repo.DB = repo.DB.Where("? = ?", k, v)
		}
	}
}

// ?where=name:test:like:or,age:30:>,status:1,id:3~5~7:not
// parsed: where age > 30 and status = 1 and id not in(3, 5, 7) or name like "%test%"
func (repo *Repo) applyWhere() {
	where, has := repo.Ctx.GetQuery("query")
	if has && where != "" {
		whereSlice := repo.ParseWhere(where, "")
		for _, cond := range whereSlice {
			if cond.IsOr {
				repo.DB = repo.DB.Or(cond.Query, cond.Args)
			} else {
				repo.DB = repo.DB.Where(cond.Query, cond.Args)
			}
		}
	}
}

// ?between=created_at~2018-05-20 15:37:25~2018-05-21 15:40:33,updated_at~2018-05-30~2018-06-08~or
// parsed: (created_at between 2018-05-20 15:37:25 and 2018-05-21 15:40:33) or (updated_at between 2018-05-30 and 2018-06-08)
func (repo *Repo) applyBetween() {
	between, has := repo.Ctx.GetQuery("between")
	if has && between != "" {
		betweenMap := strings.Split(between, ",")
		for _, item := range betweenMap {
			itemMap := strings.SplitN(item, "~", 4)
			mapLen := len(itemMap)
			if mapLen == 3 {
				repo.DB = repo.DB.Where(itemMap[0]+" BETWEEN ? AND ?", itemMap[1], itemMap[2])
			}

			if mapLen == 4 && strings.ToLower(itemMap[3]) == "or" {
				repo.DB = repo.DB.Or(itemMap[0]+" BETWEEN ? AND ?", itemMap[1], itemMap[2])
			}
		}
	}
}

// ?page=2&size=10
func (repo *Repo) applyPage() {
	if repo.Pagination.CurrentPage == 0 {
		page, hasPage := repo.Ctx.GetQuery("page")

		if hasPage && page != "" {
			p, e := strconv.Atoi(page)
			if e == nil {
				repo.Pagination.CurrentPage = p
			}
		}

		if repo.Pagination.CurrentPage < 1 {
			repo.Pagination.CurrentPage = PAGE
		}
	}

	if repo.Pagination.PageSize == 0 {
		size, hasSize := repo.Ctx.GetQuery("size")
		if hasSize && size != "" {
			s, e := strconv.Atoi(size)
			if e == nil {
				repo.Pagination.PageSize = s
			}
		}
		if repo.Pagination.PageSize < 1 { //默认10条
			repo.Pagination.PageSize = PAGESIZE
		}
	}
}

// ?order=id:desc,age:asc
func (repo *Repo) applyOrder() {
	sorter, has := repo.Ctx.GetQuery("order")
	if has && sorter != "" {
		orders := repo.parseOrder(sorter)

		for _, item := range orders {
			repo.DB = repo.DB.Order(strings.Join(item, " "))
		}
	} else if isSlice(repo.Result) { //添加一个默认的排序，防止分页时记录可能会重复出现的问题
		repo.DB = repo.DB.Order("id desc")
	}
}

func (repo Repo) parseOrder(sorter string) (sorterSlice [][]string) {
	sorterMap := strings.Split(sorter, ",")
	if len(sorterMap) > 0 {
		for _, item := range sorterMap {
			itemMap := strings.SplitN(item, ":", 2)
			itemLen := len(itemMap)
			if itemLen > 0 {
				s := []string{
					itemMap[0],
					"asc",
				}

				if itemLen > 1 {
					sortedBy := strings.ToLower(itemMap[1])
					if sortedBy == "desc" || sortedBy == "asc" {
						s[1] = sortedBy
					}
				}

				sorterSlice = append(sorterSlice, s)
			}
		}
	}
	return
}
func (repo *Repo) applyGroup() {
	groups, has := repo.Ctx.GetQuery("group")
	if has && groups != "" {
		groupMap := strings.Split(groups, ",")
		for _, item := range groupMap {
			repo.DB = repo.DB.Group(item)
		}
	}
}
func (repo Repo) ParseWhere(where string, preload string) (condSlice []condition) {
	searchMap := strings.Split(where, ",")

	for _, item := range searchMap {
		itemMap := strings.SplitN(item, ":", 4)
		itemLen := len(itemMap)
		if itemLen > 1 {
			var exp string
			if itemLen > 2 {
				exp = strings.ToLower(itemMap[2])

				if exp == "and" || exp == "or" {
					exp = ""
				}
			}

			if exp == "" && !repo.DisExpTag {
				scope := repo.DB.NewScope(repo.Result)

				if preload != "" {
					scope := repo.DB.NewScope(repo.Result)
					if fieldStruct, ok := scope.GetModelStruct().ModelType.FieldByName(preload); ok {
						// scope = repo.DB.NewScope(fieldStruct)
						repo.DB.NewScope(fieldStruct)
					}
				}
				for _, field := range scope.GetStructFields() {
					if gorm.ToDBName(field.Name) == itemMap[0] {
						exp = field.Tag.Get("exp")
					} else {
						tagMap := parseModelTags(field.Tag)
						if column, ok := tagMap["COLUMN"]; ok && column == itemMap[0] {
							exp = field.Tag.Get("exp")
						}
					}
				}

			}

			if exp == "" {
				exp = "="
			}

			var field string = itemMap[0]
			var query string
			var args interface{}
			var isOr bool
			switch exp {
			case "not":
				query = field + " NOT IN (?)"
				args = strings.Split(itemMap[1], "~")

			case "in":
				query = field + " IN (?)"
				args = strings.Split(itemMap[1], "~")

			case "like":
				query = strings.Join(append([]string{}, field, "LIKE ?"), " ")
				args = strings.Join(append([]string{}, "%", itemMap[1], "%"), "")

			default:
				query = strings.Join(append([]string{}, field, exp, "?"), " ")
				args = itemMap[1]
			}

			if (itemLen == 4 && strings.ToLower(itemMap[3]) == "or") || (itemLen == 3 && strings.ToLower(itemMap[2]) == "or") {
				isOr = true
			}

			condSlice = append(condSlice, condition{
				Field: field,
				Query: query,
				Args:  args,
				IsOr:  isOr,
			})
		}
	}
	return
}
