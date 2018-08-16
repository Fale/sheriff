package sheriff_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-version"
	"github.com/liip/sheriff"
)

type UserMultiple struct {
	Username string   `json:"username" type:"read,list" groups:"api"`
	Email    string   `json:"email" type:"read" groups:"personal"`
	Name     string   `json:"name" type:"read,list" groups:"api"`
	Roles    []string `json:"roles" type:"read" groups:"api" since:"2"`
}

type UserMultipleList []UserMultiple

func MarshalUserMultiples(version *version.Version, groups []string, types []string, users UserMultipleList) ([]byte, error) {
	o := &sheriff.Options{
		Groups: []sheriff.Group{
			{
				Values: groups,
			},
			{
				Values: types,
				Name:   "type",
			},
		},
		ApiVersion: version,
	}

	data, err := sheriff.Marshal(o, users)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(data, "", "  ")
}

func ExampleMultiple() {
	users := UserMultipleList{
		UserMultiple{
			Username: "alice",
			Email:    "alice@example.org",
			Name:     "Alice",
			Roles:    []string{"user", "admin"},
		},
		UserMultiple{
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

	output, err := MarshalUserMultiples(v1, []string{"api"}, []string{"read"}, users)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Version 1 output:")
	fmt.Printf("%s\n\n", output)

	output, err = MarshalUserMultiples(v2, []string{"api"}, []string{"read"}, users)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Version 2 output:")
	fmt.Printf("%s\n\n", output)

	output, err = MarshalUserMultiples(v2, []string{"api", "personal"}, []string{"read"}, users)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Version 2 output with personal group too:")
	fmt.Printf("%s\n\n", output)

	output, err = MarshalUserMultiples(v2, []string{"api", "personal"}, []string{"list"}, users)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Version 2 output with personal group too in list mode:")
	fmt.Printf("%s\n\n", output)

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
	//
	// Version 2 output with personal group too in list mode:
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
}
