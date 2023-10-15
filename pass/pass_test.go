package pass_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jgsheppa/localPass/pass"
)

func TestPassEntry_URLAndGeneratedPassword(t *testing.T) {
	t.Parallel()

	expectedPass := pass.Pass{
		URL: "www.google.com",
	}

	input := "www.google.com\n Y \n"
	reader := strings.NewReader(input)

	gotPass, err := pass.CreatePass(reader)
	if err != nil {
		t.Fatalf("could not create pass: %e", err)
	}

	got := expectedPass.URL
	want := gotPass.URL

	if !cmp.Equal(got, want) {
		t.Errorf("got %s, want %s", got, want)
	}
}
