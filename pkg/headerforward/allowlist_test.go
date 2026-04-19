package headerforward_test

import (
	"bfm-example/pkg/headerforward"
	"net/http"
	"reflect"
	"testing"
)

func TestFilterHeaders(t *testing.T) {
	t.Parallel()

	src := http.Header{
		"Authorization":    []string{"Bearer x"},
		"X-Tenant-Id":      []string{"t1"},
		"Cookie":           []string{"session=secret"},
		"X-Evil-Forwarded": []string{"no"},
	}

	got := headerforward.FilterHeaders(src)
	want := http.Header{
		"Authorization": []string{"Bearer x"},
		"X-Tenant-Id":   []string{"t1"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FilterHeaders() = %#v, want %#v", got, want)
	}
}

func TestFilterHeaders_empty(t *testing.T) {
	t.Parallel()
	if got := headerforward.FilterHeaders(nil); got != nil {
		t.Fatalf("FilterHeaders(nil) = %#v, want nil", got)
	}
	if got := headerforward.FilterHeaders(http.Header{}); got != nil {
		t.Fatalf("FilterHeaders(empty) = %#v, want nil", got)
	}
}
