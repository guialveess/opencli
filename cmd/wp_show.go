package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/guialveess/opencli/internal/config"
	"github.com/guialveess/opencli/internal/openproject"
	"github.com/guialveess/opencli/internal/ui"
	"github.com/spf13/cobra"
)

var wpShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Exibe detalhes de um Work Package",
	Long:  "Exibe os detalhes de um Work Package específico pelo ID.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "ID inválido: %s\n", args[0])
			os.Exit(1)
		}

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao carregar configuração: %v\n", err)
			os.Exit(1)
		}

		client := openproject.NewClient(cfg.BaseURL, cfg.APIKey, cfg.Project)

		ui.StartSpinner("Carregando Work Package...")
		wp, err := client.GetWorkPackage(id)
		ui.StopSpinner()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
			os.Exit(1)
		}

		renderWorkPackage(wp)
	},
}

var wpAssignMeCmd = &cobra.Command{
	Use:   "assign-me <id>",
	Short: "Atribui o Work Package a você mesmo",
	Long:  "Atribui o Work Package especificado pelo ID ao usuário autenticado.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "ID inválido: %s\n", args[0])
			os.Exit(1)
		}

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao carregar configuração: %v\n", err)
			os.Exit(1)
		}

		client := openproject.NewClient(cfg.BaseURL, cfg.APIKey, cfg.Project)

		ui.StartSpinner("Obtendo informações do usuário...")
		user, err := client.GetCurrentUser()
		ui.StopSpinner()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao obter usuário: %v\n", err)
			os.Exit(1)
		}

		ui.StartSpinner("Atribuindo Work Package...")
		err = client.AssignTaskForMe(id, user.ID)
		ui.StopSpinner()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao atribuir Work Package: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Work Package #%d atribuído a você com sucesso!\n", id)
	},
}

func renderWorkPackage(wp *openproject.WorkPackage) {
	idText := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#A78BFA")).
		Render(fmt.Sprintf("#%d", wp.ID))

	subjectText := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F3F4F6")).
		Render(wp.Subject)

	headerContent := lipgloss.JoinVertical(lipgloss.Left, idText, subjectText)
	header := headerBox.Render(headerContent)

	fmt.Println()
	fmt.Println(header)
	fmt.Println()

	props := []struct {
		label string
		value string
		style lipgloss.Style
	}{
		{"Status", wp.Links.Status.Title, statusStyle(wp.Links.Status.Title)},
		{"Tipo", wp.Links.Type.Title, propValueStyle},
		{"Prioridade", wp.Links.Priority.Title, propValueStyle},
	}

	if wp.Links.Assignee.Title != "" {
		props = append(props, struct {
			label string
			value string
			style lipgloss.Style
		}{"Assignee", wp.Links.Assignee.Title, lipgloss.NewStyle().Foreground(lipgloss.Color("#60A5FA"))})
	}

	for _, prop := range props {
		label := propLabelStyle.Render(prop.label)
		value := prop.style.Render(prop.value)
		fmt.Printf("%s %s\n", label, value)
	}

	fmt.Println()

	dateStyle := lipgloss.NewStyle().Foreground(mutedColor)
	createdAt := formatDate(wp.CreatedAt)
	updatedAt := formatDate(wp.UpdatedAt)

	fmt.Println(dateStyle.Render(fmt.Sprintf("Criado: %s  •  Atualizado: %s", createdAt, updatedAt)))

	if wp.Description.Raw != "" {
		fmt.Println()
		descTitle := lipgloss.NewStyle().
			Bold(true).
			Foreground(mutedColor).
			Render("Descrição")
		fmt.Println(descTitle)

		desc := descriptionStyle.Render(strings.TrimSpace(wp.Description.Raw))
		fmt.Println(desc)
	}

	fmt.Println()
}

func formatDate(isoDate string) string {
	t, err := time.Parse(time.RFC3339, isoDate)
	if err != nil {
		return isoDate
	}
	return t.Format("02/01/2006 15:04")
}

func init() {
	wpCmd.AddCommand(wpShowCmd)
	wpCmd.AddCommand(wpAssignMeCmd)
}
