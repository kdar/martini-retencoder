package retencoder

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/codegangsta/martini"
)

type Greeting struct {
	XMLName xml.Name `json:"-" xml:"response"`
	One     string   `json:"one"`
	Two     string   `json:"two"`
}

func (g *Greeting) String() string {
	return "One: " + g.One + ", Two: " + g.Two
}

func TestEncodeJSON(t *testing.T) {
	m := martini.Classic()

	m.Map(ReturnHandler())

	// routing
	m.Get("/foobar", func() (int, interface{}) {
		result := &Greeting{One: "hello", Two: "world"}
		return http.StatusOK, result
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)
	req.Header.Set("Accept", "application/json")

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get("Content-Type"), "application/json; charset=utf-8")
	expect(t, res.Body.String(), `{"one":"hello","two":"world"}`)
}

func TestEncodeXML(t *testing.T) {
	m := martini.Classic()

	m.Map(ReturnHandler())

	// routing
	m.Get("/foobar", func() (int, interface{}) {
		result := &Greeting{One: "hello", Two: "world"}
		return http.StatusOK, result
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)
	req.Header.Set("Accept", "application/xml")

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get("Content-Type"), "application/xml; charset=utf-8")
	expect(t, res.Body.String(), `<response><One>hello</One><Two>world</Two></response>`)
}

func TestEncodeText(t *testing.T) {
	m := martini.Classic()

	m.Map(ReturnHandler())

	// routing
	m.Get("/foobar", func() (int, interface{}) {
		result := &Greeting{One: "hello", Two: "world"}
		return http.StatusOK, result
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Body.String(), `One: hello, Two: world`)
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
