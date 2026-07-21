package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func init() {
	loadDotenv()
}

func loadDotenv() {
	f, err := os.Open(".env")
	if err != nil {
		return
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok || key == "" {
			continue
		}
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: migrate <apply|down|diff|hash|validate|status> [args...]")
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  apply              Apply pending up migrations")
		fmt.Println("  down [N]           Revert last N migrations (default: 1)")
		fmt.Println("  diff NAME=xxx      Generate migration from GORM models vs current state")
		fmt.Println("  hash               Re-hash migration files after manual edits")
		fmt.Println("  validate           Check migration directory integrity")
		fmt.Println("  status             Show applied/pending migrations")
		fmt.Println("")
		os.Exit(1)
	}

	cmd := args[0]
	cmdArgs := args[1:]

	ctx := context.Background()

	switch cmd {
	case "diff":
		name := extractName(cmdArgs)
		if name == "" {
			log.Fatal("NAME is required for diff (e.g., NAME=create_rules_table)")
		}
		runAtlas(ctx, "migrate", "diff", name, "--env", "gorm")
		runAtlas(ctx, "migrate", "hash", "--dir", "file://migrations")

	case "apply":
		runAtlas(ctx, "migrate", "apply", "--env", "gorm", "--allow-dirty")

	case "down":
		amount := "1"
		if len(cmdArgs) > 0 {
			amount = cmdArgs[0]
		}
		runAtlas(ctx, "migrate", "down", amount, "--env", "gorm", "--plan")

	case "hash":
		runAtlas(ctx, "migrate", "hash", "--dir", "file://migrations")

	case "validate":
		runAtlas(ctx, "migrate", "validate", "--env", "gorm")

	case "status":
		runAtlas(ctx, "migrate", "status", "--env", "gorm")

	default:
		log.Fatalf("unknown command: %s", cmd)
	}
}

func extractName(args []string) string {
	for _, a := range args {
		if strings.HasPrefix(a, "NAME=") {
			return strings.TrimPrefix(a, "NAME=")
		}
	}
	return os.Getenv("NAME")
}

func runAtlas(ctx context.Context, atlasArgs ...string) {
	log.Printf("atlas %s", strings.Join(atlasArgs, " "))
	cmd := exec.CommandContext(ctx, "atlas", atlasArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatalf("atlas failed: %v", err)
	}
}
