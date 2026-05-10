package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"promptr/internal/builtin"
	"promptr/internal/clipboard"
	"promptr/internal/prompt"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	if err := builtin.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "init error: %v\n", err)
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "search", "s":
		runSearch(args)
	case "list", "ls":
		runList(args)
	case "show", "cat":
		runShow(args)
	case "add", "new":
		runAdd(args)
	case "edit", "e":
		runEdit(args)
	case "rm", "delete":
		runRemove(args)
	case "categories", "tags":
		runCategories()
	case "cp", "copy":
		runCopy(args)
	case "--version", "-v":
		fmt.Println("promptr v" + version)
		fmt.Println()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func runSearch(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: promptr search <query>")
		os.Exit(1)
	}
	query := strings.Join(args, " ")
	results, err := prompt.Search(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if len(results) == 0 {
		fmt.Println("No prompts found")
		return
	}
	fmt.Printf("Found %d prompt(s) matching \"%s\":\n\n", len(results), query)
	for _, p := range results {
		fmt.Printf("  %s  %s  [%s]\n", p.ID, p.Title, strings.Join(p.Category, ", "))
	}
}

func runList(args []string) {
	categoryFilter := ""
	if len(args) > 0 {
		categoryFilter = args[0]
	}
	prompts, err := prompt.List(categoryFilter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if len(prompts) == 0 {
		fmt.Println("No prompts found")
		return
	}
	for _, p := range prompts {
		fmt.Printf("  %-25s %-35s [%s]\n", p.ID, p.Title, strings.Join(p.Category, ", "))
	}
}

func runShow(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: promptr show <id>")
		os.Exit(1)
	}
	id := args[0]
	p, _, err := prompt.FindByID(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prompt not found: %s\n", id)
		os.Exit(1)
	}
	fmt.Printf("ID:       %s\n", p.ID)
	fmt.Printf("Title:    %s\n", p.Title)
	fmt.Printf("Category: %s\n", strings.Join(p.Category, ", "))
	fmt.Println("---")
	fmt.Println(p.Content)
}

func runAdd(args []string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	idFlag := fs.String("id", "", "Prompt ID")
	titleFlag := fs.String("title", "", "Prompt title")
	categoryFlag := fs.String("category", "", "Comma-separated categories")
	contentFlag := fs.String("content", "", "Prompt content")
	fs.Parse(args)

	p := &prompt.Prompt{}

	if *idFlag != "" && *titleFlag != "" {
		p.ID = *idFlag
		p.Title = *titleFlag
		if *categoryFlag != "" {
			p.Category = strings.Split(*categoryFlag, ",")
		}
		p.Content = *contentFlag
	} else {
		p = interactiveAdd()
	}

	if _, _, err := prompt.FindByID(p.ID); err == nil {
		fmt.Fprintf(os.Stderr, "Error: prompt with ID \"%s\" already exists\n", p.ID)
		os.Exit(1)
	}

	path, err := prompt.UserFilePath(p.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if err := p.Save(path); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created user/%s.yaml\n", p.ID)
}

func interactiveAdd() *prompt.Prompt {
	p := &prompt.Prompt{}
	fmt.Print("ID: ")
	fmt.Scanln(&p.ID)
	fmt.Print("Title: ")
	p.Title = readLine()
	fmt.Print("Category (comma-separated): ")
	catStr := readLine()
	if catStr != "" {
		for _, c := range strings.Split(catStr, ",") {
			c = strings.TrimSpace(c)
			if c != "" {
				p.Category = append(p.Category, c)
			}
		}
	}
	fmt.Println("Content (press Enter on empty line to finish):")
	var lines []string
	for {
		line := readLine()
		if line == "" {
			break
		}
		lines = append(lines, line)
	}
	p.Content = strings.Join(lines, "\n")
	return p
}

func readLine() string {
	var sb strings.Builder
	buf := make([]byte, 1)
	for {
		n, _ := os.Stdin.Read(buf)
		if n == 0 || buf[0] == '\n' {
			break
		}
		if buf[0] != '\r' {
			sb.WriteByte(buf[0])
		}
	}
	return sb.String()
}

func runEdit(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: promptr edit <id>")
		os.Exit(1)
	}
	id := args[0]
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "notepad"
	}

	p, path, err := prompt.FindByID(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prompt not found: %s\n", id)
		os.Exit(1)
	}

	if prompt.IsBuiltin(id) {
		userPath, err := prompt.UserFilePath(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if err := p.Save(userPath); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		path = userPath
		fmt.Printf("Copied builtin prompt to user/%s.yaml\n", id)
	}

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "editor error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Edited successfully")
}

func runRemove(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: promptr rm <id>")
		os.Exit(1)
	}
	id := args[0]
	if prompt.IsBuiltin(id) {
		fmt.Fprintf(os.Stderr, "Cannot remove built-in prompt \"%s\". Use a user prompt instead.\n", id)
		os.Exit(1)
	}
	path, err := prompt.UserFilePath(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if err := os.Remove(path); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Removed %s\n", id)
}

func runCategories() {
	cats, err := prompt.AllCategories()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	for _, c := range cats {
		fmt.Println(c)
	}
}

func runCopy(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: promptr cp <id>")
		os.Exit(1)
	}
	id := args[0]
	p, _, err := prompt.FindByID(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prompt not found: %s\n", id)
		os.Exit(1)
	}

	resolved, err := prompt.ResolveVars(p.Content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := clipboard.Copy(resolved); err != nil {
		fmt.Fprintf(os.Stderr, "clipboard copy failed: %v\n", err)
		fmt.Println("--- Resolved text (stdout fallback) ---")
		fmt.Println(resolved)
		os.Exit(1)
	}
	fmt.Println("Copied to clipboard")
}

func printUsage() {
	fmt.Print(`promptr — AI prompt manager

Usage:
  promptr <command> [args]

Manage commands:
  search, s   <query>       Search prompts by query text
  list, ls     [category]   List prompts, optionally filtered by category
  show, cat    <id>         Show full prompt details
  cp, copy     <id>         Copy prompt to clipboard (with variable filling)
  categories, tags          List all category tags

Edit commands:
  add, new                  Add a new prompt (interactive or with --flags)
  edit, e      <id>         Edit a prompt in $EDITOR
  rm, delete   <id>         Remove a user prompt

Other:
  help                      Show this help
  --version, -v             Show version
`)
}
