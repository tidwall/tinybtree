package tinybtree

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"
)

func init() {
	seed := time.Now().UnixNano()
	fmt.Printf("seed: %d\n", seed)
	rand.Seed(seed)
}

func randKeys(N int) (keys []string) {
	format := fmt.Sprintf("%%0%dd", len(fmt.Sprintf("%d", N-1)))
	for _, i := range rand.Perm(N) {
		keys = append(keys, fmt.Sprintf(format, i))
	}
	return
}

const flatLeaf = true

func (tr *BTree) print() {
	tr.root.print(0, tr.height)
}

func (n *node) print(level, height int) {
	if n == nil {
		println("NIL")
		return
	}
	if height == 0 && flatLeaf {
		fmt.Printf("%s", strings.Repeat("  ", level))
	}
	for i := 0; i < n.numItems; i++ {
		if height > 0 {
			n.children[i].print(level+1, height-1)
		}
		if height > 0 || (height == 0 && !flatLeaf) {
			fmt.Printf("%s%v\n", strings.Repeat("  ", level), n.items[i].key)
		} else {
			if i > 0 {
				fmt.Printf(",")
			}
			fmt.Printf("%s", n.items[i].key)
		}
	}
	if height == 0 && flatLeaf {
		fmt.Printf("\n")
	}
	if height > 0 {
		n.children[n.numItems].print(level+1, height-1)
	}
}

func (tr *BTree) deepPrint() {
	fmt.Printf("%#v\n", tr)
	tr.root.deepPrint(0, tr.height)
}

func (n *node) deepPrint(level, height int) {
	if n == nil {
		fmt.Printf("%s %#v\n", strings.Repeat("  ", level), n)
		return
	}
	fmt.Printf("%s count: %v\n", strings.Repeat("  ", level), n.numItems)
	fmt.Printf("%s items: %v\n", strings.Repeat("  ", level), n.items)
	if height > 0 {
		fmt.Printf("%s child: %v\n", strings.Repeat("  ", level), n.children)
	}
	if height > 0 {
		for i := 0; i < n.numItems; i++ {
			n.children[i].deepPrint(level+1, height-1)
		}
		n.children[n.numItems].deepPrint(level+1, height-1)
	}
}

func stringsEquals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestDescend(t *testing.T) {
	var tr BTree
	var count int
	tr.Descend("1", func(key string, value interface{}) bool {
		count++
		return true
	})
	if count > 0 {
		t.Fatalf("expected 0, got %v", count)
	}
	var keys []string
	for i := 0; i < 1000; i += 10 {
		keys = append(keys, fmt.Sprintf("%03d", i))
		tr.Set(keys[len(keys)-1], nil)
	}
	var exp []string
	tr.Reverse(func(key string, _ interface{}) bool {
		exp = append(exp, key)
		return true
	})
	for i := 999; i >= 0; i-- {
		var key string
		key = fmt.Sprintf("%03d", i)
		var all []string
		tr.Descend(key, func(key string, value interface{}) bool {
			all = append(all, key)
			return true
		})
		for len(exp) > 0 && key < exp[0] {
			exp = exp[1:]
		}
		var count int
		tr.Descend(key, func(key string, value interface{}) bool {
			if count == (i+1)%maxItems {
				return false
			}
			count++
			return true
		})
		if count > len(exp) {
			t.Fatalf("expected 1, got %v", count)
		}

		if !stringsEquals(exp, all) {
			fmt.Printf("exp: %v\n", exp)
			fmt.Printf("all: %v\n", all)
			t.Fatal("mismatch")
		}
	}
}

func TestAscend(t *testing.T) {
	var tr BTree
	var count int
	tr.Ascend("1", func(key string, value interface{}) bool {
		count++
		return true
	})
	if count > 0 {
		t.Fatalf("expected 0, got %v", count)
	}
	var keys []string
	for i := 0; i < 1000; i += 10 {
		keys = append(keys, fmt.Sprintf("%03d", i))
		tr.Set(keys[len(keys)-1], nil)
	}
	exp := keys
	for i := -1; i < 1000; i++ {
		var key string
		if i == -1 {
			key = ""
		} else {
			key = fmt.Sprintf("%03d", i)
		}
		var all []string
		tr.Ascend(key, func(key string, value interface{}) bool {
			all = append(all, key)
			return true
		})

		for len(exp) > 0 && key > exp[0] {
			exp = exp[1:]
		}
		var count int
		tr.Ascend(key, func(key string, value interface{}) bool {
			if count == (i+1)%maxItems {
				return false
			}
			count++
			return true
		})
		if count > len(exp) {
			t.Fatalf("expected 1, got %v", count)
		}
		if !stringsEquals(exp, all) {
			t.Fatal("mismatch")
		}
	}
}

func TestBTree(t *testing.T) {
	N := 10000
	var tr BTree
	keys := randKeys(N)

	// insert all items
	for _, key := range keys {
		value, replaced := tr.Set(key, key)
		if replaced {
			t.Fatal("expected false")
		}
		if value != nil {
			t.Fatal("expected nil")
		}
	}

	// check length
	if tr.Len() != len(keys) {
		t.Fatalf("expected %v, got %v", len(keys), tr.Len())
	}

	// get each value
	for _, key := range keys {
		value, gotten := tr.Get(key)
		if !gotten {
			t.Fatal("expected true")
		}
		if value == nil || value.(string) != key {
			t.Fatalf("expected '%v', got '%v'", key, value)
		}
	}

	// scan all items
	var last string
	all := make(map[string]interface{})
	tr.Scan(func(key string, value interface{}) bool {
		if key <= last {
			t.Fatal("out of order")
		}
		if value.(string) != key {
			t.Fatalf("mismatch")
		}
		last = key
		all[key] = value
		return true
	})
	if len(all) != len(keys) {
		t.Fatalf("expected '%v', got '%v'", len(keys), len(all))
	}

	// reverse all items
	var prev string
	all = make(map[string]interface{})
	tr.Reverse(func(key string, value interface{}) bool {
		if prev != "" && key >= prev {
			t.Fatal("out of order")
		}
		if value.(string) != key {
			t.Fatalf("mismatch")
		}
		prev = key
		all[key] = value
		return true
	})
	if len(all) != len(keys) {
		t.Fatalf("expected '%v', got '%v'", len(keys), len(all))
	}

	// try to get an invalid item
	value, gotten := tr.Get("invalid")
	if gotten {
		t.Fatal("expected false")
	}
	if value != nil {
		t.Fatal("expected nil")
	}

	// scan and quit at various steps
	for i := 0; i < 100; i++ {
		var j int
		tr.Scan(func(key string, value interface{}) bool {
			if j == i {
				return false
			}
			j++
			return true
		})
	}

	// reverse and quit at various steps
	for i := 0; i < 100; i++ {
		var j int
		tr.Reverse(func(key string, value interface{}) bool {
			if j == i {
				return false
			}
			j++
			return true
		})
	}

	// delete half the items
	for _, key := range keys[:len(keys)/2] {
		value, deleted := tr.Delete(key)
		if !deleted {
			t.Fatal("expected true")
		}
		if value == nil || value.(string) != key {
			t.Fatalf("expected '%v', got '%v'", key, value)
		}
	}

	// check length
	if tr.Len() != len(keys)/2 {
		t.Fatalf("expected %v, got %v", len(keys)/2, tr.Len())
	}

	// try delete half again
	for _, key := range keys[:len(keys)/2] {
		value, deleted := tr.Delete(key)
		if deleted {
			t.Fatal("expected false")
		}
		if value != nil {
			t.Fatalf("expected nil")
		}
	}

	// try delete half again
	for _, key := range keys[:len(keys)/2] {
		value, deleted := tr.Delete(key)
		if deleted {
			t.Fatal("expected false")
		}
		if value != nil {
			t.Fatalf("expected nil")
		}
	}

	// check length
	if tr.Len() != len(keys)/2 {
		t.Fatalf("expected %v, got %v", len(keys)/2, tr.Len())
	}

	// scan items
	last = ""
	all = make(map[string]interface{})
	tr.Scan(func(key string, value interface{}) bool {
		if key <= last {
			t.Fatal("out of order")
		}
		if value.(string) != key {
			t.Fatalf("mismatch")
		}
		last = key
		all[key] = value
		return true
	})
	if len(all) != len(keys)/2 {
		t.Fatalf("expected '%v', got '%v'", len(keys), len(all))
	}

	// replace second half
	for _, key := range keys[len(keys)/2:] {
		value, replaced := tr.Set(key, key)
		if !replaced {
			t.Fatal("expected true")
		}
		if value == nil || value.(string) != key {
			t.Fatalf("expected '%v', got '%v'", key, value)
		}
	}

	// delete next half the items
	for _, key := range keys[len(keys)/2:] {
		value, deleted := tr.Delete(key)
		if !deleted {
			t.Fatal("expected true")
		}
		if value == nil || value.(string) != key {
			t.Fatalf("expected '%v', got '%v'", key, value)
		}
	}

	// check length
	if tr.Len() != 0 {
		t.Fatalf("expected %v, got %v", 0, tr.Len())
	}

	// do some stuff on an empty tree
	value, gotten = tr.Get(keys[0])
	if gotten {
		t.Fatal("expected false")
	}
	if value != nil {
		t.Fatal("expected nil")
	}
	tr.Scan(func(key string, value interface{}) bool {
		t.Fatal("should not be reached")
		return true
	})
	tr.Reverse(func(key string, value interface{}) bool {
		t.Fatal("should not be reached")
		return true
	})

	var deleted bool
	value, deleted = tr.Delete("invalid")
	if deleted {
		t.Fatal("expected false")
	}
	if value != nil {
		t.Fatal("expected nil")
	}
}

func BenchmarkTidwallSequentialSet(b *testing.B) {
	var tr BTree
	keys := randKeys(b.N)
	sort.Strings(keys)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Set(keys[i], nil)
	}
}

func BenchmarkTidwallSequentialGet(b *testing.B) {
	var tr BTree
	keys := randKeys(b.N)
	sort.Strings(keys)
	for i := 0; i < b.N; i++ {
		tr.Set(keys[i], nil)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Get(keys[i])
	}
}

func BenchmarkTidwallRandomSet(b *testing.B) {
	var tr BTree
	keys := randKeys(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Set(keys[i], nil)
	}
}

func BenchmarkTidwallRandomGet(b *testing.B) {
	var tr BTree
	keys := randKeys(b.N)
	for i := 0; i < b.N; i++ {
		tr.Set(keys[i], nil)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Get(keys[i])
	}
}

// type googleKind struct {
// 	key string
// }

// func (a *googleKind) Less(b btree.Item) bool {
// 	return a.key < b.(*googleKind).key
// }

// func BenchmarkGoogleSequentialSet(b *testing.B) {
// 	tr := btree.New(32)
// 	keys := randKeys(b.N)
// 	sort.Strings(keys)
// 	gkeys := make([]*googleKind, len(keys))
// 	for i := 0; i < b.N; i++ {
// 		gkeys[i] = &googleKind{keys[i]}
// 	}
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		tr.ReplaceOrInsert(gkeys[i])
// 	}
// }

// func BenchmarkGoogleSequentialGet(b *testing.B) {
// 	tr := btree.New(32)
// 	keys := randKeys(b.N)
// 	gkeys := make([]*googleKind, len(keys))
// 	for i := 0; i < b.N; i++ {
// 		gkeys[i] = &googleKind{keys[i]}
// 	}
// 	for i := 0; i < b.N; i++ {
// 		tr.ReplaceOrInsert(gkeys[i])
// 	}
// 	sort.Strings(keys)
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		tr.Get(gkeys[i])
// 	}
// }

// func BenchmarkGoogleRandomSet(b *testing.B) {
// 	tr := btree.New(32)
// 	keys := randKeys(b.N)
// 	gkeys := make([]*googleKind, len(keys))
// 	for i := 0; i < b.N; i++ {
// 		gkeys[i] = &googleKind{keys[i]}
// 	}
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		tr.ReplaceOrInsert(gkeys[i])
// 	}
// }

// func BenchmarkGoogleRandomGet(b *testing.B) {
// 	tr := btree.New(32)
// 	keys := randKeys(b.N)
// 	gkeys := make([]*googleKind, len(keys))
// 	for i := 0; i < b.N; i++ {
// 		gkeys[i] = &googleKind{keys[i]}
// 	}
// 	for i := 0; i < b.N; i++ {
// 		tr.ReplaceOrInsert(gkeys[i])
// 	}
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		tr.Get(gkeys[i])
// 	}
// }

func TestBTreeOne(t *testing.T) {
	var tr BTree
	tr.Set("1", "1")
	tr.Delete("1")
	tr.Set("1", "1")
	tr.Delete("1")
	tr.Set("1", "1")
	tr.Delete("1")
}

func TestBTree256(t *testing.T) {
	var tr BTree
	var n int
	for j := 0; j < 2; j++ {
		for _, i := range rand.Perm(256) {
			tr.Set(fmt.Sprintf("%d", i), i)
			n++
			if tr.Len() != n {
				t.Fatalf("expected 256, got %d", n)
			}
		}
		for _, i := range rand.Perm(256) {
			v, ok := tr.Get(fmt.Sprintf("%d", i))
			if !ok {
				t.Fatal("expected true")
			}
			if v.(int) != i {
				t.Fatalf("expected %d, got %d", i, v.(int))
			}
		}
		for _, i := range rand.Perm(256) {
			tr.Delete(fmt.Sprintf("%d", i))
			n--
			if tr.Len() != n {
				t.Fatalf("expected 256, got %d", n)
			}
		}
		for _, i := range rand.Perm(256) {
			_, ok := tr.Get(fmt.Sprintf("%d", i))
			if ok {
				t.Fatal("expected false")
			}
		}
	}
}
