package database

import "container/list"

// LRUCache lru结构
type LRUCache struct {
	cap      int
	mapCache map[string]*list.Element
	list     *list.List
}

// LRUNode 存放数据
type LRUNode struct {
	key string
	val interface{}
}

// Constructor 构造
func Constructor(capacity int) LRUCache {
	return LRUCache{
		cap:      capacity,
		mapCache: make(map[string]*list.Element),
		list:     list.New(),
	}
}

// Get 获取数据
func (lp *LRUCache) Get(key string) (bool, interface{}) {

	if elem, ok := lp.mapCache[key]; ok {
		lp.list.MoveToFront(elem)
		return true, lp.mapCache[key].Value.(LRUNode).val
	}

	return false, nil
}

// Put 存放数据
func (lp *LRUCache) Put(key string, value interface{}) {

	if elem, ok := lp.mapCache[key]; ok {
		lp.list.MoveToFront(elem)
		elem.Value = LRUNode{
			key: key,
			val: value,
		}
	} else {

		if lp.list.Len() >= lp.cap { // 淘汰尾部元素
			delete(lp.mapCache, lp.list.Back().Value.(LRUNode).key)
			lp.list.Remove(lp.list.Back())
		}

		nod := LRUNode{
			key: key,
			val: value,
		}
		lp.list.PushFront(nod)
		lp.mapCache[nod.key] = lp.list.Front()
	}

}
