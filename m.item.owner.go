package gev

type IItemOwnerModel interface {
	IItemModel
}

type ItemOwnerModel struct {
	ItemModel `xorm:"extends"`
	OwnerId   int `json:"owner_id,omitempty" xorm:"not null"`
}

func (o *ItemOwnerModel) CanRead(user IUserModel) bool {
	if user == nil {
		return false
	}
	if o.OwnerId == user.GetId() {
		return true
	}
	return false
}

func (o *ItemOwnerModel) CanWrite(user IUserModel) bool {
	if user == nil {
		return false
	}
	if o.OwnerId == 0 {
		o.OwnerId = user.GetId()
		return true
	}
	if o.OwnerId == user.GetId() {
		return true
	}
	return false
}

func (m *ItemOwnerModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		m.ItemModel.Bind(g, m)
	} else {
		m.ItemModel.Bind(g, self)
	}
}
