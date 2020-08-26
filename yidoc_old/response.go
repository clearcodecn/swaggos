package yidoc

type Response interface {
	IsObject() bool
	IsArray() bool
	ArrayElement() *Item
	Object() []*Item
	Item() *Item
}

type basicItemResponse struct {
	item *Item
}

func (i *basicItemResponse) IsObject() bool { return false }

func (i *basicItemResponse) IsArray() bool { return false }

func (i *basicItemResponse) ArrayElement() *Item { return nil }

func (i *basicItemResponse) Object() []*Item { return nil }

func (i *basicItemResponse) Item() *Item { return i.item }

type objectResponse struct {
	items []*Item
}

func (o *objectResponse) IsObject() bool { return true }

func (o *objectResponse) IsArray() bool { return false }

func (o *objectResponse) ArrayElement() *Item { return nil }

func (o *objectResponse) Object() []*Item { return o.items }

func (o *objectResponse) Item() *Item { return nil }

type arrayResponse struct {
	item *Item
}

func (o *arrayResponse) IsObject() bool { return false }

func (o *arrayResponse) IsArray() bool { return true }

func (o *arrayResponse) ArrayElement() *Item { return o.item }

func (o *arrayResponse) Object() []*Item { return nil }

func (o *arrayResponse) Item() *Item { return nil }
