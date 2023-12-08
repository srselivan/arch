package consolepresenter

import (
	"arch/internal/viewer/entity"
	"fmt"
	"os"
	"os/exec"
)

type Presenter struct {
}

func New() *Presenter {
	return &Presenter{}
}

func (p *Presenter) UpdateView(messages []entity.Message) error {
	if err := p.cleanTerminal(); err != nil {
		return fmt.Errorf("p.cleanTerminal: %w", err)
	}

	for _, message := range messages {
		fmt.Println(message.String())
	}

	return nil
}

func (p *Presenter) cleanTerminal() error {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cmd.Run: %w", err)
	}
	return nil
}
