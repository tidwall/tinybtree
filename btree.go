package tinybtree

import "github.com/tidwall/btree"

// BTree is an ordered set of key/value pairs where the key is a string
// and the value is an interface{}
type BTree struct {
	base *btree.BTree
}

type item struct {
	key   string
	value interface{}
}

func (tr *BTree) init() {
	tr.base = btree.NewNonConcurrent(func(a, b interface{}) bool {
		return a.(*item).key < b.(*item).key
	})
}

// Set or replace a value for a key
func (tr *BTree) Set(key string, value interface{}) (prev interface{}, replaced bool) {
	if tr.base == nil {
		tr.init()
	}
	if v := tr.base.Set(&item{key, value}); v != nil {
		return v.(*item).value, true
	}
	return nil, false
}

// Get a value for key
func (tr *BTree) Get(key string) (value interface{}, gotten bool) {
	if tr.base == nil {
		return nil, false
	}
	if v := tr.base.Get(&item{key, value}); v != nil {
		return v.(*item).value, true
	}
	return nil, false
}

// Delete a value for a key
func (tr *BTree) Delete(key string) (prev interface{}, deleted bool) {
	if tr.base == nil {
		tr.init()
	}
	if v := tr.base.Delete(&item{key: key}); v != nil {
		return v.(*item).value, true
	}
	return nil, false
}

// Len returns the number of items in the tree
func (tr *BTree) Len() int {
	if tr.base == nil {
		return 0
	}
	return tr.base.Len()
}

// Ascend the tree within the range [pivot, last]
func (tr *BTree) Ascend(
	pivot string,
	iter func(key string, value interface{}) bool,
) {
	if tr.base == nil {
		return
	}
	tr.base.Ascend(&item{key: pivot}, func(v interface{}) bool {
		return iter(v.(*item).key, v.(*item).value)
	})
}

// Scan all items in tree
func (tr *BTree) Scan(iter func(key string, value interface{}) bool) {
	if tr.base == nil {
		return
	}
	tr.base.Ascend(nil, func(v interface{}) bool {
		return iter(v.(*item).key, v.(*item).value)
	})
}

// Descend the tree within the range [pivot, first]
func (tr *BTree) Descend(
	pivot string,
	iter func(key string, value interface{}) bool,
) {
	if tr.base == nil {
		return
	}
	tr.base.Descend(&item{key: pivot}, func(v interface{}) bool {
		return iter(v.(*item).key, v.(*item).value)
	})
}

// Reverse interates over all items in tree, in reverse.
func (tr *BTree) Reverse(iter func(key string, value interface{}) bool) {
	if tr.base == nil {
		return
	}
	tr.base.Descend(nil, func(v interface{}) bool {
		return iter(v.(*item).key, v.(*item).value)
	})
}
