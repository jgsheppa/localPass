package pass

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/jgsheppa/localPass/generator"
	"github.com/jgsheppa/localPass/models"
	"golang.org/x/term"
)

type Pass struct {
	URL         string
	password    string
	reader      io.Reader
	writer      io.Writer
	PassService models.PassService
}

type option func(*Pass) error

func NewPass() *Pass {
	pass := &Pass{
		reader: os.Stdin,
		writer: os.Stdout,
	}

	return pass
}

func (p *Pass) WithService(service models.PassService) *Pass {
	p.PassService = service
	return p
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

	cleanAnswer := strings.ToLower(strings.TrimSpace(answer))

	if cleanAnswer == "" || cleanAnswer == "y" {
		conditions := generator.NewPasswordConditions()
		pw, err := conditions.GeneratePassword()
		if err != nil {
			return err
		}
		p.password = pw
		return nil
	}

	fmt.Fprint(p.writer, "Enter your password: ")
	userPassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	p.password = string(userPassword)

	return nil
}

func (p *Pass) SavePass() (*Pass, error) {
	pass := &models.Pass{
		URL:      p.URL,
		Password: p.password,
	}
	err := p.PassService.Create(pass)
	if err == models.ErrURLInvalid {
		return p, err
	} else if err != nil {
		return p, err
	}

	fmt.Println("Pass successfully created!")
	return p, nil
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

func Run(service models.PassService) (int, error) {
	newPass, err := NewPass().WithInput(os.Stdin).WithOutput(os.Stdout).WithService(service).CreatePass()
	if err != nil {
		return 1, err
	}
	_, err = newPass.SavePass()
	if err != nil {
		return 1, err
	}

	return 0, nil
}
