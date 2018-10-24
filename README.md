# `tinybtree`

[![GoDoc](https://godoc.org/github.com/tidwall/tinybtree?status.svg)](https://godoc.org/github.com/tidwall/tinybtree)

Just an itsy bitsy b-tree.

## Usage

Keys are strings, values are interfaces.

### Functions

```
Get(key string) (value interface{}, gotten bool)
Set(key string, value interface{}) (prev interface{}, replaced bool)
Delete(key string) (prev interface{}, deleted bool)
Scan(iter func(key string, value interface{}) bool)
Ascend(pivot string, iter func(key string, value interface{}) bool)
Descend(pivot string, iter func(key string, value interface{}) bool)
```

### Example

```go
// Create a btree
var tr tinybtree.BTree

// Set a key. Returns the previous value and ok a previous value exists.
prev, ok := tr.Set("hello", "world")

// Get a key. Returns the value and ok if the value exists.
value, ok := tr.Get("hello")

// Delete a key. Returns the deleted value and ok if the previous value exists.
prev, ok := tr.Delete("hello")
```

## Contact

Josh Baker [@tidwall](http://twitter.com/tidwall)

## License

`tinybtree` source code is available under the MIT [License](/LICENSE).
