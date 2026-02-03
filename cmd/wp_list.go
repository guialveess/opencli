package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/guialveess/opencli/internal/config"
	"github.com/guialveess/opencli/internal/openproject"
	"github.com/guialveess/opencli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	listPage     int
	listPageSize int
	listAll      bool

	assigneeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60A5FA")).
			Background(lipgloss.Color("#1E3A5F")).
			Padding(0, 1).
			MarginRight(2)

	emptyAssigneeStyle = lipgloss.NewStyle().Width(20)
)

func renderAssignee(name string) string {
	if name == "" {
		return emptyAssigneeStyle.Render("")
	}
	return assigneeStyle.Render(name)
}

var wpListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista os Work Packages do projeto",
	Long:  "Lista todos os Work Packages do projeto configurado no OpenProject.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao carregar configuração: %v\n", err)
			os.Exit(1)
		}

		client := openproject.NewClient(cfg.BaseURL, cfg.APIKey, cfg.Project)

		header := lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

		if listAll {
			ui.StartSpinner("Carregando Work Packages...")
			workPackages, err := client.ListAllWorkPackages()
			ui.StopSpinner()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao listar Work Packages: %v\n", err)
				os.Exit(1)
			}

			fmt.Println(header.Render(fmt.Sprintf("Work Packages (%d)", len(workPackages))))
			fmt.Println()

			for _, wp := range workPackages {
				id := idStyle.Render(fmt.Sprintf("#%-5d", wp.ID))
				status := statusStyle(wp.Links.Status.Title).Render(fmt.Sprintf("%-12s", wp.Links.Status.Title))
				assignee := renderAssignee(wp.Links.Assignee.Title)
				subject := subjectStyle.Render(wp.Subject)

				fmt.Printf("%s  %s  %s  %s\n", id, status, assignee, subject)
			}
			return
		}

		ui.StartSpinner("Carregando Work Packages...")
		page, err := client.ListWorkPackages(listPage, listPageSize)
		ui.StopSpinner()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao listar Work Packages: %v\n", err)
			os.Exit(1)
		}

		pageInfo := lipgloss.NewStyle().Foreground(warningColor)
		fmt.Println(header.Render(fmt.Sprintf("Work Packages (%d total)", page.Total)))
		fmt.Println(pageInfo.Render(fmt.Sprintf("Página %d de %d", page.Page, page.TotalPages)))
		fmt.Println()

		for _, wp := range page.Items {
			id := idStyle.Render(fmt.Sprintf("#%-5d", wp.ID))
			status := statusStyle(wp.Links.Status.Title).Render(fmt.Sprintf("%-12s", wp.Links.Status.Title))
			assignee := renderAssignee(wp.Links.Assignee.Title)
			subject := subjectStyle.Render(wp.Subject)

			fmt.Printf("%s  %s  %s  %s\n", id, status, assignee, subject)
		}

		if page.HasNextPage {
			fmt.Println()
			fmt.Println(pageInfo.Render(fmt.Sprintf("Use --page %d para ver mais", page.Page+1)))
		}
	},
}

func init() {
	wpListCmd.Flags().IntVarP(&listPage, "page", "p", 1, "Número da página")
	wpListCmd.Flags().IntVarP(&listPageSize, "size", "s", 70, "Itens por página")
	wpListCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Lista todos os Work Packages")
	wpCmd.AddCommand(wpListCmd)
}
