package prompt

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var placeholderRe = regexp.MustCompile(`\{\{(\w+)\}\}`)

func ExtractVars(content string) []string {
	matches := placeholderRe.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool)
	var vars []string
	for _, m := range matches {
		name := strings.TrimSpace(m[1])
		if name != "" && !seen[name] {
			seen[name] = true
			vars = append(vars, name)
		}
	}
	return vars
}

func ResolveVars(content string) (string, error) {
	vars := ExtractVars(content)
	if len(vars) == 0 {
		return content, nil
	}
	replaced := content
	scanner := bufio.NewScanner(os.Stdin)
	for _, v := range vars {
		fmt.Printf("%s: ", v)
		if !scanner.Scan() {
			break
		}
		value := strings.TrimSpace(scanner.Text())
		replaced = strings.ReplaceAll(replaced, "{{"+v+"}}", value)
	}
	return replaced, nil
}
