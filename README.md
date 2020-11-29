# Go Struct Patch

## Get and Import
To use with **modules**, add to import list:
```
import "github.com/sztelzer/structpatch"
```
Then update your module:
`$ go mod download` or `$ go mod tidy`

Or without modules **go get** to add to GOPATH:
```
$ go get "github.com/sztelzer/structpatch"
```

## Use
The only public function is:
```
func Patch(src interface{}, dst interface{}, lockTag string) error {
```

Example:
```
type T1 struct {
	A	string              // skip if zero value
	B	*string             // skip if nil
	C	string `lock:""`    // skip if tag used
	D	*string `lock:""`   // skip if tag used
}

src = &T1{...}
dst = &T1{...}
err := structpatch.Patch(src, dst, "lock")
```

Patch will range over public source fields of any types. If field name is present in the destination, it may be overwritten following these rules:

```
   a == zero	 skip
   a != (*)b	 skip
   a 	         b =>    a
   a            *b =>   *a
  *a == nil      skip
(*)a != b        skip
(*)a != (*)b     skip
  *a             b => (*)a
  *a           	*b =>   *a
```
Legend:
- a: represents source field value type
- b: represents destination field value type
- \*a: represents source field pointer to value type
- (\*)a: represents underlying value type of pointer
- zero is the zero value for the type (see reflection.IsZero())
- nil is a nil pointer
- b is equivalent to 'a' but on destiny

Source and Destination must be pointers to struct, but don't need to be of the same Type.

**Nil Pointer Source Fields will be skipped.** This is probably the most useful rule, as nil pointer may differentiate if we want to set a value to it's zero value or not set. As rule of thumb, if data comes in a bunch, as json.Unmarshal(), consider using a src struct with only pointer fields.

Pointer Source Fields will set destination respecting if it is pointer or value. Destination pointer can be nil as it will be set.

Only the pointer will be copied, so values are referenced, not deep copied.

You can pass a field tag of the destination to lock it and skip.

Beware that Zero values (on naked type fields) will skip. If some new value is equal to the type Zero value (string:"", int:0, etc.) it will not be used! Use pointer field source if you want to set to zero value.