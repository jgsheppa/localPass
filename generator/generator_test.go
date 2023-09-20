package generator_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jgsheppa/localPass/generator"
)

func TestPasswordGenerator_Length(t *testing.T) {
	t.Parallel()

	passwordConditions := generator.NewPasswordConditions(generator.WithLength(24))

	generated, _ := passwordConditions.GeneratePassword()

	got := len(generated)
	want := 24

	if !cmp.Equal(got, want) {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestPasswordGenerator_Random(t *testing.T) {
	t.Parallel()

	passwordConditions := generator.NewPasswordConditions(generator.WithLength(24))

	got, _ := passwordConditions.GeneratePassword()
	want, _ := passwordConditions.GeneratePassword()

	if cmp.Equal(got, want) {
		t.Errorf("got %s, want %s", got, want)
	}
}
