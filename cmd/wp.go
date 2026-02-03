package cmd

import (
	"github.com/spf13/cobra"
)

var wpCmd = &cobra.Command{
	Use:   "wp",
	Short: "Gerencia Work Packages",
	Long:  "Comandos para listar, criar e atualizar Work Packages no OpenProject.",
}

func init() {
	rootCmd.AddCommand(wpCmd)
}
