package pass

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/jgsheppa/localPass/generator"
)

type Pass struct {
	URL      string
	password string
	reader   *bufio.Reader
}

func (p *Pass) EnterURL() error {
	fmt.Print("Enter a URL: ")
	url, err := p.reader.ReadString('\n')
	if err != nil {
		return err
	}

	trimURL := strings.TrimSpace(url)
	p.URL = trimURL

	return nil
}

func (p *Pass) EnterPassword() error {
	fmt.Print("Generate password? (Y/n): ")
	answer, err := p.reader.ReadString('\n')
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
		return nil
	}

	fmt.Print("Enter your password: ")
	userPassword, err := p.reader.ReadString('\n')
	if err != nil {
		return err
	}
	p.password = userPassword

	return nil
}

func CreatePass(rd io.Reader) (Pass, error) {
	var pass Pass

	pass.reader = bufio.NewReader(rd)
	pass.EnterURL()
	pass.EnterPassword()

	return pass, nil
}
