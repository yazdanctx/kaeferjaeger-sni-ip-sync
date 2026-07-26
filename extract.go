package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

const defaultInput = "/root/kaeferyeager/sni-ip-sync/final.txt"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: extract [-f file] <domain>")
		os.Exit(1)
	}

	input := defaultInput
	args := os.Args[1:]

	if len(args) >= 2 && args[0] == "-f" {
		input = args[1]
		args = args[2:]
	}

	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: extract [-f file] <domain>")
		os.Exit(1)
	}

	domain := args[0]

	f, err := os.Open(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	seen := make(map[string]bool)
	var results []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, domain) {
			continue
		}

		parts := strings.Split(line, "][")
		if len(parts) < 2 {
			continue
		}

		content := parts[0]
		if idx := strings.LastIndex(content, "["); idx != -1 {
			content = content[idx+1:]
		}

		snis := strings.Fields(content)
		for _, sni := range snis {
			sni = strings.TrimLeft(sni, "*.")
			if sni == "" || !strings.Contains(sni, domain) {
				continue
			}
			if !seen[sni] {
				seen[sni] = true
				results = append(results, sni)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	sort.Strings(results)

	outname := fmt.Sprintf("x_%s_%s.txt", domain, time.Now().Format("2006-01-02_150405"))
	out, err := os.Create(outname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	for _, r := range results {
		fmt.Fprintln(w, r)
	}
	w.Flush()

	fmt.Printf("Extracted %d unique domains -> %s\n", len(results), outname)
}
