package tinybox

import "container/list"

// LRU is a LRU container.
// It is not safe for concurrent access for now.
type LRU[T any] struct {
	maxBytes  int
	curBytes  int
	list      *list.List
	elements  map[string]*list.Element
	onEvicted func(key string, value T)
}

func NewLRU[T any](maxBytes int, onEvicted func(key string, value T)) *LRU[T] {
	return &LRU[T]{
		maxBytes:  maxBytes,
		curBytes:  0,
		list:      list.New(),
		elements:  map[string]*list.Element{},
		onEvicted: onEvicted,
	}
}

type Item[T any] struct {
	key  string
	val  T
	size int
}

// Add adds a new key-value pair to the LRU.
func (l *LRU[T]) Add(key string, val T, size int) {
	item := Item[T]{
		key:  key,
		val:  val,
		size: size,
	}

	if ele, ok := l.elements[key]; ok {
		old := ele.Value.(Item[T])
		oldSize := old.size
		l.curBytes += size - oldSize

		l.list.MoveToFront(ele)
		ele.Value = item
	} else {
		l.curBytes += size

		ele := l.list.PushFront(item)
		l.elements[key] = ele
	}

	for l.maxBytes != 0 && l.curBytes > l.maxBytes {
		l.removeLast()
	}
}

// Get lookups the mapped value by key.
func (l *LRU[T]) Get(key string) (value T, ok bool) {
	ele, ok := l.elements[key]
	if !ok {
		return value, false
	}

	l.list.MoveToFront(ele)

	i := ele.Value.(Item[T])
	return i.val, true
}

func (l *LRU[T]) removeLast() {
	ele := l.list.Back()
	if ele == nil {
		return
	}

	l.list.Remove(ele)
	item, _ := ele.Value.(Item[T])
	delete(l.elements, item.key)

	l.curBytes -= item.size

	if l.onEvicted != nil {
		l.onEvicted(item.key, item.val)
	}
}
