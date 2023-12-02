package main

import (
	"fmt"
	"strings"

	"github.com/lusingander/kasane"
)

func main() {
	bgline := fmt.Sprintf("\x1b[3m\x1b[38;5;238m%s\x1b[0m", strings.Repeat("kasane", 7))
	bg := strings.Join([]string{bgline, bgline, bgline}, "\n")

	fg := "\x1b[1m\x1b[31m kasane - 重ね \x1b[0m"

	fmt.Println(kasane.OverlayString(bg, fg, 1, 13))
}
