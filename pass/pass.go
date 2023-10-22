package pass

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/jgsheppa/localPass/generator"
	"golang.org/x/term"
)

type Pass struct {
	URL      string
	password string
	reader   io.Reader
	writer   io.Writer
}

type option func(*Pass) error

func NewPass() *Pass {
	pass := &Pass{
		reader: os.Stdin,
		writer: os.Stdout,
	}

	return pass
}

func (p *Pass) WithInput(input io.Reader) *Pass {
	p.reader = input
	return p
}

func (p *Pass) WithOutput(output io.Writer) *Pass {
	p.writer = output
	return p
}

func (p *Pass) EnterURL(reader *bufio.Reader) error {
	fmt.Fprint(p.writer, "Enter a URL: ")
	url, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	trimURL := strings.TrimSpace(url)
	p.URL = trimURL

	return nil
}

func (p *Pass) EnterPassword(reader *bufio.Reader) error {
	fmt.Fprint(p.writer, "Generate password? (Y/n): ")
	answer, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	trimAnswer := strings.TrimSpace(answer)

	if trimAnswer == "Y" || trimAnswer == "y" {
		conditions := generator.NewPasswordConditions()
		pw, err := conditions.GeneratePassword()
		if err != nil {
			return err
		}
		p.password = pw
		fmt.Println("Pass successfully created!")
		return nil
	}

	fmt.Fprint(p.writer, "Enter your password: ")
	userPassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	p.password = string(userPassword)
	fmt.Println("Pass successfully created!")

	return nil
}

func (p *Pass) CreatePass() (*Pass, error) {
	reader := bufio.NewReader(p.reader)

	err := p.EnterURL(reader)
	if err != nil {
		return p, err
	}
	err = p.EnterPassword(reader)
	if err != nil {
		return p, err
	}

	return p, nil
}

func Run() (int, error) {
	_, err := NewPass().WithInput(os.Stdin).WithOutput(os.Stdout).CreatePass()
	if err != nil {
		return 1, err
	}

	return 0, nil
}
