package pass_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jgsheppa/localPass/pass"
)

func TestPassEntry_URLAndGeneratedPassword(t *testing.T) {
	t.Parallel()

	expectedPass := pass.Pass{
		URL: "www.test.com",
	}

	input := "www.test.com\n Y \n"
	reader := strings.NewReader(input)

	var output bytes.Buffer

	gotPass, err := pass.NewPass().WithInput(reader).WithOutput(&output).CreatePass()
	if err != nil {
		t.Fatalf("could not create pass: %e", err)
	}

	got := expectedPass.URL
	want := gotPass.URL

	if !cmp.Equal(got, want) {
		t.Errorf("got %s, want %s", got, want)
	}

	got = output.String()
	want = "Enter a URL: Generate password? (Y/n):"

	if !strings.Contains(got, want) {
		t.Errorf("got %s, want %s", got, want)
	}
}
