package gev

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

type ISearchModel interface {
	IModel
	GetCondition() ISearch
	SearchSession(user IUserModel, session *xorm.Session, condition ISearch)
	Search(user IUserModel, condition ISearch) (interface{}, error)
}

// M.search Entity
type SearchModel struct {
	Model `xorm:"extends"`
}

func (m *SearchModel) GetCondition() ISearch {
	return &SearchBody{}
}

func (s *SearchModel) SearchSession(user IUserModel, session *xorm.Session, condition ISearch) {
	condition.MakeSession(session)
}

func GetSearchData2(bean interface{}, condition ISearch, sessionFunc func(session *xorm.Session)) (*SearchData, error) {
	session := Db.NewSession()
	defer session.Close()
	sessionFunc(session)
	total, _ := session.Count(bean)
	sessionFunc(session)
	session.Limit(condition.GetSize(), condition.GetBegin())
	condition.GetOrder(session)
	data := make([]interface{}, condition.GetSize())
	n := 0
	err := session.Iterate(bean, func(i int, item interface{}) error {
		data[i] = item
		n++
		return nil
	})
	return &SearchData{data[:n], total}, err
}

func GetSearchData(bean IModel, condition ISearch, sessionFunc func(session *xorm.Session)) (*SearchData, error) {
	session := Db.NewSession()
	defer session.Close()
	sessionFunc(session)
	total, _ := session.Count(bean)
	sessionFunc(session)
	session.Limit(condition.GetSize(), condition.GetBegin())
	condition.GetOrder(session)
	data := make([]interface{}, condition.GetSize())
	n := 0
	err := session.Iterate(bean, func(i int, item interface{}) error {
		model := item.(IModel)
		model.SetSelf(model)
		data[i] = model.GetSearch()
		n++
		return nil
	})
	return &SearchData{data[:n], total}, err
}

func (m *SearchModel) Search(user IUserModel, condition ISearch) (interface{}, error) {
	bean := m.Self().(ISearchModel)
	return GetSearchData(bean, condition, func(session *xorm.Session) {
		bean.SearchSession(user, session, condition)
	})
}

func (m *SearchModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		self = m
	}
	m.Model.Bind(g, self)
	g.Info("搜索").Body(
		self.(ISearchModel).GetCondition(),
	).Data(
		NewSearchData(10, []interface{}{self.(ISearchModel).GetSearch()}),
	).POST("/search", func(c *gin.Context) {
		condition := m.Self().(ISearchModel).GetCondition()
		err := c.BindJSON(condition)
		if err != nil {
			Err(c, 1, err)
			return
		}
		if user, ok := c.Get("user"); ok {
			data, err := m.New().(ISearchModel).Search(user.(IUserModel), condition)
			Api(c, data, err)
		} else {
			data, err := m.New().(ISearchModel).Search(nil, condition)
			Api(c, data, err)
		}
	})
}
func NewSearchData(total int64, content []interface{}) map[string]interface{} {
	return map[string]interface{}{"content": content, "total": total}
}

type SearchData struct {
	Content []interface{} `json:"content" xorm:""`
	Total   int64         `json:"total" xorm:""`
}
