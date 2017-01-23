# sheriff
[![GoDoc](https://godoc.org/github.com/liip/sheriff?status.svg)](https://godoc.org/github.com/liip/sheriff) [![Build Status](https://travis-ci.org/liip/sheriff.svg?branch=master)](https://travis-ci.org/liip/sheriff) [![Coverage Status](https://coveralls.io/repos/github/liip/sheriff/badge.svg?branch=master)](https://coveralls.io/github/liip/sheriff?branch=master)

```
go get github.com/liip/sheriff
```

Package sheriff transforms structs into a map based on specific tags on the
struct fields. A typical use is an API which marshals structs into JSON and
maintains different API versions. Using sheriff, struct fields can be annotated
with API version and group tags. By invoking sheriff with specific options,
those tags determine whether a field will be added to the output map or not. It
can then be marshalled using "encoding/json".

Example:
```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-version"
	"github.com/liip/sheriff"
)

type User struct {
	Username string   `json:"username" groups:"api"`
	Email    string   `json:"email" groups:"personal"`
	Name     string   `json:"name" groups:"api"`
	Roles    []string `json:"roles" groups:"api" since:"2"`
}

type UserList []User

func MarshalUsers(version *version.Version, groups []string, users UserList) ([]byte, error) {
	o := &sheriff.Options{
		Groups:     groups,
		ApiVersion: version,
	}

	data, err := sheriff.Marshal(o, users)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(data, "", "  ")
}

func main() {
	users := UserList{
		User{
			Username: "alice",
			Email:    "alice@example.org",
			Name:     "Alice",
			Roles:    []string{"user", "admin"},
		},
		User{
			Username: "bob",
			Email:    "bob@example.org",
			Name:     "Bob",
			Roles:    []string{"user"},
		},
	}

	v1, err := version.NewVersion("1.0.0")
	if err != nil {
		log.Panic(err)
	}
	v2, err := version.NewVersion("2.0.0")

	output, err := MarshalUsers(v1, []string{"api"}, users)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Version 1 output:")
	fmt.Printf("%s\n\n", output)

	output, err = MarshalUsers(v2, []string{"api"}, users)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Version 2 output:")
	fmt.Printf("%s\n\n", output)

	output, err = MarshalUsers(v2, []string{"api", "personal"}, users)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Version 2 output with personal group too:")
	fmt.Printf("%s\n\n", output)

}

// Output:
// Version 1 output:
// [
//   {
//     "name": "Alice",
//     "username": "alice"
//   },
//   {
//     "name": "Bob",
//     "username": "bob"
//   }
// ]
//
// Version 2 output:
// [
//   {
//     "name": "Alice",
//     "roles": [
//       "user",
//       "admin"
//     ],
//     "username": "alice"
//   },
//   {
//     "name": "Bob",
//     "roles": [
//       "user"
//     ],
//     "username": "bob"
//   }
// ]
//
// Version 2 output with personal group too:
// [
//   {
//     "email": "alice@example.org",
//     "name": "Alice",
//     "roles": [
//       "user",
//       "admin"
//     ],
//     "username": "alice"
//   },
//   {
//     "email": "bob@example.org",
//     "name": "Bob",
//     "roles": [
//       "user"
//     ],
//     "username": "bob"
//   }
// ]
```

## Benchmarks

There's a simple benchmark in `bench_test.go` which compares running sheriff -> JSON versus just marshalling into JSON 
and runs on every build. Just marshalling JSON itself takes usually between 3 and 5 times less nanoseconds per operation
compared to running sheriff and JSON.

Want to make sheriff faster? Please send us your pull request or open an issue discussing a possible improvement 🚀!

## Acknowledgements

- This idea and code has been created partially during a [Liip](https://liip.ch) hackday.
- Thanks to [@basgys](https://github.com/basgys) for reviews & improvements.
