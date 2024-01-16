package behavior

import (
	"fmt"
	"sync"

	"github.com/pojol/gobot/database"
)

type behaviorPool struct {
	sync.Mutex
	saved map[string][]*Tree
}

func _new_tree(name string) (*Tree, error) {

	dat, err := database.GetBehavior().Find(name)
	if err != nil {
		return nil, err
	}

	tree, err := Load(dat.File)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func Get(name string) *Tree {
	bp.Lock()
	defer bp.Unlock()

	lst := bp.saved[name]

	n := len(lst)
	if n == 0 {
		tree, err := _new_tree(name)
		if err != nil {
			fmt.Println("behaviorPool.Get", err.Error())
			return nil
		}

		return tree
	}

	x := lst[n-1]
	bp.saved[name] = bp.saved[name][0 : n-1]

	return x
}

func Put(name string, tree *Tree) {
	bp.Lock()
	defer bp.Unlock()

	tree.Reset()
	bp.saved[name] = append(bp.saved[name], tree)

}

var bp = &behaviorPool{
	saved: make(map[string][]*Tree),
}
