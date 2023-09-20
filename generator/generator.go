package generator

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

const (
	lowerCharSet   string = "abcdedfghijklmnopqrst"
	upperCharSet   string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet string = "!@#$%&*"
	numberSet      string = "0123456789"
	allCharSet     string = lowerCharSet + upperCharSet + specialCharSet + numberSet
)

type PasswordConditions struct {
	length   int
	Password string
}

type PasswordConditionsOption func(*PasswordConditions)

func NewPasswordConditions(opts ...PasswordConditionsOption) *PasswordConditions {
	pc := &PasswordConditions{
		length: 24,
	}

	for _, opt := range opts {
		opt(pc)
	}

	return pc
}

func WithLength(length int) PasswordConditionsOption {
	return func(pc *PasswordConditions) {
		pc.length = length
	}
}

func (pc *PasswordConditions) GeneratePassword() (string, error) {
	var password strings.Builder

	if pc.length < 8 {
		return "", errors.New("Password must have a length greater than or equal to 8")
	}

	for i := 0; i < pc.length; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}

	return password.String(), nil
}

func Run(length *int) (int, error) {
	conditions := NewPasswordConditions(WithLength(*length))
	password, err := conditions.GeneratePassword()
	if err != nil {
		return 1, err
	}

	if *length != 0 {
		fmt.Println(password)
	}
	return 0, nil
}
