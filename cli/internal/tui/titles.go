package tui

import (
	"strings"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

type TitleWriter interface {
	ShowTitle(title string)
	ShowTitleAndDescription(title, description string)
	ShowSubTitle(mainTitle, subtitle string)
}

type Title struct {
}

func (t *Title) ShowTitle(title string) {
	titleNormalised := strings.TrimSpace(strings.ToUpper(title))
	s, _ := pterm.DefaultBigText.WithLetters(putils.LettersFromString(titleNormalised)).
		Srender()
	pterm.DefaultCenter.Println(s)
}

func (t *Title) ShowTitleAndDescription(title, description string) {
	titleNormalised := strings.TrimSpace(strings.ToUpper(title))
	s, _ := pterm.DefaultBigText.WithLetters(putils.LettersFromString(titleNormalised)).
		Srender()
	pterm.DefaultCenter.Println(s)

	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(description)
}

func (t *Title) ShowSubTitle(title, subTitle string) {
	_ = pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle(strings.ToUpper(title), pterm.NewStyle(pterm.FgCyan)),
		putils.LettersFromStringWithStyle(strings.ToUpper(subTitle), pterm.NewStyle(pterm.FgLightMagenta))).
		Render()
}

func NewTitleWriter() TitleWriter {
	return &Title{}
}
