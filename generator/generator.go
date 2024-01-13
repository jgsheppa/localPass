package generator

import (
	"errors"
	"flag"
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

type Password struct {
	length   int
	Password string
}

type PasswordOption func(*Password)

func NewPassword(opts ...PasswordOption) *Password {
	pc := &Password{
		length: 24,
	}

	for _, opt := range opts {
		opt(pc)
	}

	return pc
}

func WithLength(length int) PasswordOption {
	return func(pc *Password) {
		pc.length = length
	}
}

func (pc *Password) GeneratePassword() (string, error) {
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
	conditions := NewPassword(WithLength(*length))
	password, err := conditions.GeneratePassword()
	if err != nil {
		return 1, err
	}

	if *length != 0 {
		fmt.Println(password)
	}
	return 0, nil
}

type PasswordFlags struct {
	FlagSet *flag.FlagSet
	Length  *int
}

func Flag() *PasswordFlags {
	pwdFlag := flag.NewFlagSet("password", flag.ExitOnError)

	pwdFlag.Usage = func() {
		fmt.Println("Usage: [password] [...flags] ")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	length := pwdFlag.Int("length", 24, "Length of password")

	return &PasswordFlags{
		FlagSet: pwdFlag,
		Length:  length,
	}
}
