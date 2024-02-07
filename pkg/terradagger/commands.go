package terradagger

import "strings"

const (
	BashEntrypoint = "bash"
	ShEntrypoint   = "sh"
)

func BuildCommand(binary string, command string, args []string) string {
	var sb strings.Builder

	sb.WriteString(binary)
	if command != "" {
		sb.WriteString(" " + command)
	}

	for _, arg := range args {
		sb.WriteString(" " + arg)
	}

	return sb.String()
}

func RunWithBash(command string) []string {
	return []string{BashEntrypoint, "-c", command}
}

func RunWithSh(command string) []string {
	return []string{ShEntrypoint, "-c", command}
}
