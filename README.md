# Go Struct Patch

To use with module, add to import list:
```
import "github.com/sztelzer/structpatch"
```
Or go get:
```
$ go get "github.com/sztelzer/structpatch"
```

Use:
```
package main
import (
    "log"
    sp "github.com/sztelzer/structpatch"
)

type T struct {
    A string    // this will patch if A is not ""
    B *string   // this will patch if referenced value is not ""
    C string `lock:""` // this will never be touched
    D *string `lock:""` // this too
}

func main() {
    var s1 = "a"
    var s2 = "z"
    var sE = ""

    var dst = &T1 {A:s1, B:&s1, C:s1, D:&s1}
    var srcNew = &T1 {A:s2, B:&s2, C:s2, D:&s2}    
    err := sp.StructPatch(srcNew, dst, "lock")
    log.Print(dst)
    // "z", "z", "a", "a"

    dst = &T1 {A:s1, B:&s1, C:s1, D:&s1}
    var srcEmpty = &T1 {A:sE, B:&sE, C:sE, D:&sE}
    err := sp.StructPatchTag(srcNew, dst, "lock")
    log.Print(dst)
    // "a", "a", "a", "a"
}
```

Source and destination must be pointers to struct.

Pointer fields will be copied if underlying source value is not type zero (reflect IsZero()).

Only the pointer will be copied, so structs are only referenced, not deep copied.

Any tag can be used as lock.

Struct Types need to be the same.

Beware that Zero values will not be used to patch. So, if some new value is equal to the type Zero value (string:"", int:0, etc.) it will not be used!
To overcome this limitation identifying not changed values, you could use pointer fields:

```
type T2 struc {
    A *string
    B *string
}
```

