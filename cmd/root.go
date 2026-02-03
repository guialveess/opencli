package cmd

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "op",
	Short: "CLI para interagir com o OpenProject",
	Long:  "op é uma CLI para gerenciar Work Packages e outras entidades do OpenProject via API REST.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.AddTemplateFunc("split", strings.Split)
	cobra.AddTemplateFunc("styleHeading", styleHeading)
	cobra.AddTemplateFunc("styleCommand", styleCommand)
	cobra.AddTemplateFunc("styleFlag", styleFlag)
	cobra.AddTemplateFunc("styleDescription", styleDescription)

	rootCmd.SetUsageTemplate(usageTemplate)
	rootCmd.SetHelpTemplate(helpTemplate)
}

func styleHeading(s string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		Render(s)
}

func styleCommand(s string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A78BFA")).
		Render(s)
}

func styleFlag(s string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#22D3EE")).
		Render(s)
}

func styleDescription(s string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Render(s)
}

var helpTemplate = `{{ styleHeading .Short }}

{{ if .Long }}{{ .Long }}

{{ end }}{{ styleHeading "Uso:" }}
  {{ styleCommand .UseLine }}{{ if .HasAvailableSubCommands }}
  {{ styleCommand (print .CommandPath " [comando]") }}{{ end }}{{ if gt (len .Aliases) 0 }}

{{ styleHeading "Aliases:" }}
  {{ .NameAndAliases }}{{ end }}{{ if .HasAvailableSubCommands }}

{{ styleHeading "Comandos:" }}{{ range .Commands }}{{ if .IsAvailableCommand }}
  {{ styleCommand (rpad .Name .NamePadding) }} {{ styleDescription .Short }}{{ end }}{{ end }}{{ end }}{{ if .HasAvailableLocalFlags }}

{{ styleHeading "Flags:" }}
{{ range $line := split .LocalFlags.FlagUsages "\n" }}{{ if $line }}  {{ styleFlag $line }}
{{ end }}{{ end }}{{ end }}{{ if .HasAvailableInheritedFlags }}
{{ styleHeading "Flags Globais:" }}
{{ range $line := split .InheritedFlags.FlagUsages "\n" }}{{ if $line }}  {{ styleFlag $line }}
{{ end }}{{ end }}{{ end }}{{ if .HasHelpSubCommands }}
{{ styleHeading "Tópicos Adicionais:" }}{{ range .Commands }}{{ if .IsAdditionalHelpTopicCommand }}
  {{ styleCommand (rpad .CommandPath .CommandPathPadding) }} {{ styleDescription .Short }}{{ end }}{{ end }}{{ end }}{{ if .HasAvailableSubCommands }}
Use "{{ styleCommand (print .CommandPath " [comando] --help") }}" para mais informações.{{ end }}
`

var usageTemplate = `{{ styleHeading "Uso:" }}
  {{ styleCommand .UseLine }}{{ if .HasAvailableSubCommands }}
  {{ styleCommand (print .CommandPath " [comando]") }}{{ end }}{{ if gt (len .Aliases) 0 }}

{{ styleHeading "Aliases:" }}
  {{ .NameAndAliases }}{{ end }}{{ if .HasExample }}

{{ styleHeading "Exemplos:" }}
{{ .Example }}{{ end }}{{ if .HasAvailableSubCommands }}

{{ styleHeading "Comandos:" }}{{ range .Commands }}{{ if .IsAvailableCommand }}
  {{ styleCommand (rpad .Name .NamePadding) }} {{ styleDescription .Short }}{{ end }}{{ end }}{{ end }}{{ if .HasAvailableLocalFlags }}

{{ styleHeading "Flags:" }}
{{ range $line := split .LocalFlags.FlagUsages "\n" }}{{ if $line }}  {{ styleFlag $line }}
{{ end }}{{ end }}{{ end }}{{ if .HasAvailableInheritedFlags }}
{{ styleHeading "Flags Globais:" }}
{{ range $line := split .InheritedFlags.FlagUsages "\n" }}{{ if $line }}  {{ styleFlag $line }}
{{ end }}{{ end }}{{ end }}{{ if .HasHelpSubCommands }}
{{ styleHeading "Tópicos Adicionais:" }}{{ range .Commands }}{{ if .IsAdditionalHelpTopicCommand }}
  {{ styleCommand (rpad .CommandPath .CommandPathPadding) }} {{ styleDescription .Short }}{{ end }}{{ end }}{{ end }}{{ if .HasAvailableSubCommands }}
Use "{{ styleCommand (print .CommandPath " [comando] --help") }}" para mais informações.{{ end }}
`
