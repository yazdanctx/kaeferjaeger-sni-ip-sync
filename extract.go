package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const defaultInput = "./final.txt"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: extract [-f file] <domain>")
		os.Exit(1)
	}

	input := defaultInput
	args := os.Args[1:]
	var domain string

	for i := 0; i < len(args); i++ {
		if args[i] == "-f" && i+1 < len(args) {
			input = args[i+1]
			i++
		} else if domain == "" {
			domain = args[i]
		}
	}

	if domain == "" {
		fmt.Fprintln(os.Stderr, "Usage: extract [-f file] <domain>")
		os.Exit(1)
	}

	seen := make(map[string]bool)
	var results []string

	add := func(s string) {
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		if !seen[s] {
			seen[s] = true
			results = append(results, s)
		}
	}

	// --- Pass 1: Go parser on the input file ---
	if err := parseFile(input, domain, add); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: parse error: %v\n", err)
	}

	// --- Pass 2: bash pipeline on *.txt in the same directory ---
	dir := filepath.Dir(input)
	if err := runPipeline(dir, domain, add); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: pipeline error: %v\n", err)
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

func parseFile(path, domain string, add func(string)) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.Contains(line, domain) {
			continue
		}

		start := strings.Index(line, "[")
		end := strings.LastIndex(line, "]")
		if start == -1 || end == -1 || end <= start {
			continue
		}

		snis := strings.Fields(line[start+1 : end])
		for _, sni := range snis {
			sni = cleanSNI(sni)
			if sni == "" || !strings.Contains(sni, domain) {
				continue
			}
			add(sni)
		}
	}

	return scanner.Err()
}

func runPipeline(dir, domain string, add func(string)) error {
	grepPattern := "." + domain

	cmd := exec.Command("bash", "-c",
		fmt.Sprintf(
			`grep -F %q %s/*.txt | awk -F'-- ' '{print $2}' | tr ' [' $'\n''\n' | sed 's/\]//g' | grep -F %q | sort -u`,
			grepPattern, dir, grepPattern,
		),
	)

	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil // grep found nothing — not an error
		}
		return fmt.Errorf("pipeline failed: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := cleanSNI(strings.TrimSpace(scanner.Text()))
		if line != "" && strings.Contains(line, domain) {
			add(line)
		}
	}

	return scanner.Err()
}

func cleanSNI(sni string) string {
	if strings.HasPrefix(sni, "*.") {
		sni = sni[2:]
	} else if strings.HasPrefix(sni, "*") {
		sni = sni[1:]
	}

	sni = strings.TrimLeft(sni, ".")

	return strings.ToLower(sni)
}
