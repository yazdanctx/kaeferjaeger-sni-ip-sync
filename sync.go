package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	baseURL    = "https://kaeferjaeger.gay/sni-ip-ranges"
	dataDir    = "data"
	outputFile = "final.txt"
)

var providers = []string{
	"amazon",
	"digitalocean",
	"google",
	"microsoft",
	"oracle",
}

func main() {
	fmt.Println("Starting SNI IP range sync...")
	start := time.Now()

	for _, provider := range providers {
		if err := download(provider); err != nil {
			fmt.Fprintf(os.Stderr, "FATAL: %s download failed: %v\n", provider, err)
			os.Exit(1)
		}
	}

	if err := merge(); err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: merge failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done in %s\n", time.Since(start).Round(time.Second))
}

func download(provider string) error {
	url := fmt.Sprintf("%s/%s/ipv4_merged_sni.txt", baseURL, provider)
	dir := filepath.Join(dataDir, provider)
	path := filepath.Join(dir, "ipv4_merged_sni.txt")
	tmpPath := path + ".tmp"

	fmt.Printf("Downloading %s...\n", provider)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	f, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return err
	}

	if err := f.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return os.Rename(tmpPath, path)
}

func merge() error {
	fmt.Println("Merging into final.txt...")

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, provider := range providers {
		path := filepath.Join(dataDir, provider, "ipv4_merged_sni.txt")
		src, err := os.Open(path)
		if err != nil {
			return err
		}

		if _, err := io.Copy(f, src); err != nil {
			src.Close()
			return err
		}
		src.Close()
	}

	return nil
}
