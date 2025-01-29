package entities

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"
)

const (
	timeFormat = time.RFC3339
)

func stringHandler(w io.Writer, v any) {
	fmt.Fprint(w, v)
}

func timeHandler(w io.Writer, v any) {
	t := v.(time.Time)
	fmt.Fprintf(w, "%s", t.Format(timeFormat))
}

func mapListHandler(w io.Writer, v any) {
	if v == nil {
		fmt.Fprint(w, "<none>")
		return
	}

	mapVals, ok := v.(map[string]string)
	if !ok || len(mapVals) == 0 {
		fmt.Fprint(w, "<empty>")
		return
	}

	pairs := make([]string, 0, len(mapVals))
	for k, val := range mapVals {
		pairs = append(pairs, fmt.Sprintf("%s=%v", k, val))
	}

	fmt.Fprint(w, strings.Join(pairs, ", "))
}

func mapPageHandler(w io.Writer, v interface{}) {
	if v == nil {
		fmt.Fprint(w, "<none>")
		return
	}

	m, ok := v.(map[string]string)
	if !ok {
		fmt.Fprint(w, "<none>")
		return
	}

	if len(m) == 0 {
		fmt.Fprint(w, "<empty>")
		return
	}

	pairs := make([]string, 0, len(m))
	iter := reflect.ValueOf(m).MapRange()
	for iter.Next() {
		pairs = append(pairs, fmt.Sprintf("%v=%v", iter.Key(), iter.Value()))
	}

	fmt.Fprint(w, strings.Join(pairs, "\n\t"))
}
