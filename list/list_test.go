package list_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jgsheppa/localPass/list"
)

func TestList(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer

	list.NewList().WithOutput(&output).WithSites([]string{"www.test.com", "www.test2.com"}).PrintList()

	got := output.String()
	want := "www.test.com \nwww.test2.com \n"
	if !cmp.Equal(got, want) {
		t.Errorf("got %s, want %s", got, want)
	}
}
