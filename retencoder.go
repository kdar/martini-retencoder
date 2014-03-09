package retencoder

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/codegangsta/inject"
	"github.com/codegangsta/martini"
)

// ReturnHandler is a service that is called when a route handler returns
// something. The ReturnHandler is responsible for writing to the
// ResponseWriter based on the values that are passed into this function.
func ReturnHandler() martini.ReturnHandler {
	return func(ctx martini.Context, vals []reflect.Value) {
		rv := ctx.Get(reflect.TypeOf((*http.Request)(nil)))
		req := rv.Interface().(*http.Request)
		acceptType := req.Header.Get("Accept")

		rv = ctx.Get(inject.InterfaceOf((*http.ResponseWriter)(nil)))
		res := rv.Interface().(http.ResponseWriter)
		var responseVal reflect.Value
		if len(vals) > 1 && vals[0].Kind() == reflect.Int {
			res.WriteHeader(int(vals[0].Int()))
			responseVal = vals[1]
		} else if len(vals) > 0 {
			responseVal = vals[0]
		}

		if strings.Contains(acceptType, "json") {
			res.Header().Set("Content-Type", "application/json; charset=utf-8")
			data, err := json.Marshal(responseVal.Interface())
			if err == nil {
				res.Write(data)
			}
		} else if strings.Contains(acceptType, "xml") {
			res.Header().Set("Content-Type", "application/xml; charset=utf-8")
			data, err := xml.Marshal(responseVal.Interface())
			if err == nil {
				res.Write(data)
			}
		} else {
			if stringer, ok := responseVal.Interface().(fmt.Stringer); ok {
				res.Write([]byte(stringer.String()))
			}
		}
	}
}
