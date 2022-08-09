package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem) *ListItem
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len  int
	head *ListItem
	tail *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	if l.len == 0 {
		return nil
	}
	return l.head
}

func (l *list) Back() *ListItem {
	if l.len == 0 {
		return nil
	}
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	tmp := &ListItem{Value: v}
	tmp.Next = l.head
	tmp.Prev = nil
	if l.head != nil {
		l.head.Prev = tmp
	}
	l.head = tmp
	if l.tail == nil {
		l.tail = tmp
	}
	l.len++
	return tmp
}

func (l *list) PushBack(v interface{}) *ListItem {
	tmp := &ListItem{Value: v}
	tmp.Next = nil
	tmp.Prev = l.tail
	if l.tail != nil {
		l.tail.Next = tmp
	}
	l.tail = tmp
	if l.head == nil {
		l.head = tmp
	}
	l.len++
	return tmp
}

func (l *list) Remove(i *ListItem) *ListItem {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if l.tail == i {
		l.tail = i.Prev
	}
	if l.head == i {
		l.head = i.Next
	}
	i.Next = nil
	i.Prev = nil
	l.len--
	return i
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.head {
		return
	}
	i.Prev.Next = i.Next
	if i.Prev == l.head {
		i.Prev.Prev = i
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i == l.tail {
		l.tail = i.Prev
	}
	i.Prev = nil
	i.Next = l.head
	l.head = i
}

func (l *list) MoveToFront2(i *ListItem) {
	if i == l.head {
		return
	}
	l.head = i

}

func NewList() List {
	newList := new(list)
	newList.head = nil
	newList.tail = nil
	newList.len = 0
	return newList
}
