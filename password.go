package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	pw "github.com/sethvargo/go-password/password"
)

type GenConfig struct {
	Length      int
	NumDigits   int
	NumSymbols  int
	NoUpper     bool
	AllowRepeat bool
}

func randInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

func Generate(cfg GenConfig) string {
	// A lib já valida limites razoáveis internamente.
	s, err := pw.Generate(cfg.Length, cfg.NumDigits, cfg.NumSymbols, cfg.NoUpper, cfg.AllowRepeat)
	if err != nil {
		// fallback mínimo: tenta um padrão seguro caso algo dê errado
		s, _ = pw.Generate(16, 4, 2, false, false)
	}
	return s
}

func SavePassword(pwd string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	file := filepath.Join(home, ".passwords")

	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04")
	line := fmt.Sprintf("[%s] %s\n", timestamp, pwd)
	_, err = f.WriteString(line)
	return err
}
