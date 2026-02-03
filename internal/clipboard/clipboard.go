package clipboard

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func GetImageFromClipboard() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return getImageMacOS()
	case "linux":
		return getImageLinux()
	default:
		return "", fmt.Errorf("sistema operacional não suportado: %s", runtime.GOOS)
	}
}

func getImageMacOS() (string, error) {
	tmpFile := filepath.Join(os.TempDir(), "opcli-clipboard.png")

	script := `
		use framework "AppKit"
		set pb to current application's NSPasteboard's generalPasteboard()
		set imgData to pb's dataForType:(current application's NSPasteboardTypePNG)
		if imgData is missing value then
			error "Nenhuma imagem no clipboard"
		end if
		set filePath to "` + tmpFile + `"
		imgData's writeToFile:filePath atomically:true
		return filePath
	`

	cmd := exec.Command("osascript", "-l", "AppleScript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return getImageMacOSFallback()
	}

	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		return "", fmt.Errorf("falha ao salvar imagem: %s", string(output))
	}

	return tmpFile, nil
}

func getImageMacOSFallback() (string, error) {
	tmpFile := filepath.Join(os.TempDir(), "opcli-clipboard.png")

	if _, err := exec.LookPath("pngpaste"); err == nil {
		cmd := exec.Command("pngpaste", tmpFile)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("nenhuma imagem no clipboard")
		}
		return tmpFile, nil
	}

	return "", fmt.Errorf("nenhuma imagem no clipboard (instale pngpaste: brew install pngpaste)")
}

func getImageLinux() (string, error) {
	tmpFile := filepath.Join(os.TempDir(), "opcli-clipboard.png")

	if _, err := exec.LookPath("xclip"); err == nil {
		cmd := exec.Command("xclip", "-selection", "clipboard", "-t", "image/png", "-o")
		output, err := cmd.Output()
		if err != nil || len(output) == 0 {
			return "", fmt.Errorf("nenhuma imagem no clipboard")
		}
		if err := os.WriteFile(tmpFile, output, 0644); err != nil {
			return "", err
		}
		return tmpFile, nil
	}

	if _, err := exec.LookPath("xsel"); err == nil {
		cmd := exec.Command("xsel", "--clipboard", "--output")
		output, err := cmd.Output()
		if err != nil || len(output) == 0 {
			return "", fmt.Errorf("nenhuma imagem no clipboard")
		}
		if err := os.WriteFile(tmpFile, output, 0644); err != nil {
			return "", err
		}
		return tmpFile, nil
	}

	return "", fmt.Errorf("instale xclip ou xsel para usar clipboard no Linux")
}

func Cleanup(path string) {
	if path != "" && filepath.Dir(path) == os.TempDir() {
		os.Remove(path)
	}
}
