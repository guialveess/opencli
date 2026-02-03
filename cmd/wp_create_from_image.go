package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/guialveess/opencli/internal/clipboard"
	"github.com/guialveess/opencli/internal/config"
	"github.com/guialveess/opencli/internal/ollama"
	"github.com/guialveess/opencli/internal/openproject"
	"github.com/guialveess/opencli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	ollamaModel   string
	autoConfirm   bool
	fromClipboard bool
)

var wpCreateFromImageCmd = &cobra.Command{
	Use:   "create-from-image [image-path]",
	Short: "Cria um Work Package a partir de uma imagem",
	Long: `Analisa uma imagem (screenshot de bug, erro, etc.) usando IA local (Ollama)
e cria um Work Package com título e descrição gerados automaticamente.

Requer Ollama rodando localmente com um modelo de visão (ex: llava, minicpm-v).

Exemplos:
  op wp create-from-image ./screenshot-bug.png
  op wp create-from-image ./erro.jpg --model llava
  op wp create-from-image ./bug.png -y  # cria sem confirmação
  op wp create-from-image --clipboard   # usa imagem do clipboard`,
	Args: cobra.MaximumNArgs(1),
	Run:  runCreateFromImage,
}

func runCreateFromImage(cmd *cobra.Command, args []string) {
	var imagePath string
	var cleanupPath string

	if fromClipboard {
		ui.StartSpinner("Obtendo imagem do clipboard...")
		path, err := clipboard.GetImageFromClipboard()
		ui.StopSpinner()

		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}
		imagePath = path
		cleanupPath = path
		ui.PrintSuccess("Imagem obtida do clipboard")
	} else {
		if len(args) == 0 {
			ui.PrintError("Forneça o caminho da imagem ou use --clipboard")
			fmt.Println()
			ui.PrintInfo("Uso: op wp create-from-image <image-path>")
			ui.PrintInfo("     op wp create-from-image --clipboard")
			os.Exit(1)
		}

		imagePath = args[0]

		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			ui.PrintError(fmt.Sprintf("Arquivo não encontrado: %s", imagePath))
			os.Exit(1)
		}
	}

	if cleanupPath != "" {
		defer clipboard.Cleanup(cleanupPath)
	}

	cfg, err := config.Load()
	if err != nil {
		ui.PrintError(fmt.Sprintf("Erro ao carregar configuração: %v", err))
		os.Exit(1)
	}

	ollamaClient := ollama.NewClient(ollamaModel)

	fmt.Println()
	ui.StartThinkingSpinner("IA analisando imagem...")
	analysis, err := ollamaClient.AnalyzeScreenshot(imagePath)
	ui.StopSpinner()

	if err != nil {
		ui.PrintError(fmt.Sprintf("Erro ao analisar imagem: %v", err))
		fmt.Println()
		ui.PrintInfo("Verifique se o Ollama está rodando: ollama serve")
		ui.PrintInfo(fmt.Sprintf("E se o modelo está instalado: ollama pull %s", ollamaModel))
		os.Exit(1)
	}

	fmt.Println(ui.RenderAnalysisResult(analysis.Title, analysis.Description))

	if !autoConfirm {
		promptStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A78BFA")).
			Bold(true)

		fmt.Print(promptStyle.Render("\nCriar Work Package? "))
		fmt.Print("[Y/n]: ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "" && response != "y" && response != "yes" && response != "s" && response != "sim" {
			ui.PrintInfo("Operação cancelada")
			return
		}
	}

	opClient := openproject.NewClient(cfg.BaseURL, cfg.APIKey, cfg.Project)

	ui.StartSpinner("Criando Work Package...")
	wp, err := opClient.CreateWorkPackage(&openproject.CreateWorkPackageRequest{
		Subject:     analysis.Title,
		Description: analysis.Description,
	})
	ui.StopSpinner()

	if err != nil {
		ui.PrintError(fmt.Sprintf("Erro ao criar Work Package: %v", err))
		os.Exit(1)
	}

	fmt.Println()
	successBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#10B981")).
		Padding(0, 2).
		Foreground(lipgloss.Color("#10B981"))

	idStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#A78BFA"))

	msg := fmt.Sprintf("Work Package %s criado com sucesso!", idStyle.Render(fmt.Sprintf("#%d", wp.ID)))
	fmt.Println(successBox.Render(msg))
}

func init() {
	wpCreateFromImageCmd.Flags().StringVarP(&ollamaModel, "model", "m", "llava", "Modelo Ollama para análise de imagem")
	wpCreateFromImageCmd.Flags().BoolVarP(&autoConfirm, "yes", "y", false, "Criar sem pedir confirmação")
	wpCreateFromImageCmd.Flags().BoolVarP(&fromClipboard, "clipboard", "c", false, "Usar imagem do clipboard")

	wpCmd.AddCommand(wpCreateFromImageCmd)
}
