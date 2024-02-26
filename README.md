# go-immutable
[![Build Status](https://github.com/mkeeler/go-immutable/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/mkeeler/go-immutable/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/mkeeler/go-immutable)](https://goreportcard.com/report/github.com/mkeeler/go-immutable) [![PkgGoDev](https://pkg.go.dev/badge/github.com/mkeeler/go-immutable)](https://pkg.go.dev/github.com/mkeeler/go-immutable)

This library contains functions to operate on Go types in a generic and immutable way. Much of the overall behavior exists in the standard library already and you should really only use this library if you specifically require the immutability guarantees.

See it in action:

```go
package foo

import (
   "slices"
   
   "github.com/mkeeler/go-immutable/immutableslice"
)

func FindAndRemove(s []int, v int) []int {
   idx := slices.Index(s, v)
   return immutableslice.Delete(s, idx, idx+1)
}
```
