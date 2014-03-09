retencoder
==========

Allows you to return values from your handlers and have it encoded based on the Accept header. Supports JSON, XML, and fmt.Stringer.

I made this to answer a question on the [martini googlegroups](https://groups.google.com/forum/#!topic/martini-go/Ppu3v1wxYg4).

### Usage

```
m.Map(retencoder.ReturnHandler())
```

### Example

```
package main

import (
	"encoding/xml"
	"log"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/kdar/martini-retencoder"
)

type Some struct {
	XMLName  xml.Name `json:"-" xml:"response"`
	Login    string   `json:"login" xml:"login"`
	Password string   `json:"password" xml:"password"`
}

func (s *Some) String() string {
	return "Login: " + s.Login + ", Password: " + s.Password
}

func main() {
	m := martini.New()
	route := martini.NewRouter()

	m.Map(retencoder.ReturnHandler())

	route.Get("/test", func() (int, interface{}) {
		result := &Some{Login: "awesome", Password: "hidden"}
		return http.StatusOK, result
	})

	m.Action(route.Handle)

	log.Println("Waiting for connections...")

	if err := http.ListenAndServe(":8000", m); err != nil {
		log.Fatal(err)
	}
}

```

### Results

```
curl http://localhost:8000/test -H "Accept: application/json"
{"login":"awesome","password":"hidden"}

curl http://localhost:8000/test -H "Accept: application/text"
Login: awesome, Password: hidden

curl http://localhost:8000/test -H "Accept: application/xml"
<response><login>awesome</login><password>hidden</password></response>
```