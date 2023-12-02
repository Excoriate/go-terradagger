package tui

import (
	"strings"

	"github.com/pterm/pterm"
)

type MessageWriter interface {
	ShowError(opt MessageOptions)
	ShowInfo(opt MessageOptions)
	ShowSuccess(opt MessageOptions)
	ShowWarning(opt MessageOptions)
}

type Message struct{}

type MessageOptions struct {
	Title   string
	Message string
	Error   error
}

func (t *Message) ShowError(opt MessageOptions) {
	if opt.Title != "" {
		pterm.Error.Prefix = pterm.Prefix{
			Text:  strings.ToUpper(opt.Title),
			Style: pterm.NewStyle(pterm.BgCyan, pterm.FgRed),
		}
	}

	if opt.Error == nil {
		pterm.Error.Println(opt.Message)
		return
	}

	var errMsg string
	if opt.Message != "" {
		errMsg = opt.Message + ": " + opt.Error.Error()
	} else {
		errMsg = opt.Error.Error()
	}

	pterm.Error.Println(errMsg)
}

func (t *Message) ShowInfo(opt MessageOptions) {
	if opt.Title != "" {
		pterm.Info.Prefix = pterm.Prefix{
			Text:  strings.ToUpper(opt.Title),
			Style: pterm.NewStyle(pterm.BgCyan, pterm.FgBlack),
		}
	}

	pterm.Info.Println(opt.Message)
}

func (t *Message) ShowSuccess(opt MessageOptions) {
	if opt.Title != "" {
		pterm.Success.Prefix = pterm.Prefix{
			Text:  strings.ToUpper(opt.Title),
			Style: pterm.NewStyle(pterm.BgCyan, pterm.FgBlack),
		}
	}

	pterm.Success.Println(opt.Message)
}

func (t *Message) ShowWarning(opt MessageOptions) {
	if opt.Title != "" {
		pterm.Warning.Prefix = pterm.Prefix{
			Text:  strings.ToUpper(opt.Title),
			Style: pterm.NewStyle(pterm.BgCyan, pterm.FgBlack),
		}
	}

	pterm.Warning.Println(opt.Message)
}

func NewMessageWriter() MessageWriter {
	return &Message{}
}
