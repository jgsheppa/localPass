package generator_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jgsheppa/localPass/generator"
)

func TestPasswordGenerator_Length(t *testing.T) {
	t.Parallel()

	passwordGenerator := generator.NewPassword(generator.WithLength(24))

	generatedPw, err := passwordGenerator.GeneratePassword()
	if err != nil {
		t.Errorf("could not generate password in TestPasswordGenerator_Length: %v", err)
	}

	got := len(generatedPw)
	want := 24

	if !cmp.Equal(got, want) {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestPasswordGenerator_Random(t *testing.T) {
	t.Parallel()

	passwordGenerator := generator.NewPassword(generator.WithLength(24))

	got, err := passwordGenerator.GeneratePassword()
	if err != nil {
		t.Errorf("could not generate password for got in TestPasswordGenerator_Random: %v", err)
	}
	want, err := passwordGenerator.GeneratePassword()
	if err != nil {
		t.Errorf("could not generate password for want in TestPasswordGenerator_Random: %v", err)
	}

	if cmp.Equal(got, want) {
		t.Errorf("got %s, want %s", got, want)
	}
}
