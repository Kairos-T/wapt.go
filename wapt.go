package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func fetchFirst15Lines(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var lines []string
	scanner := bufio.NewScanner(resp.Body)
	for i := 0; i < 15 && scanner.Scan(); i++ {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}

func checkForFalsePositive(fileURL string, homepageContent string) {
	fileContent, err := fetchFirst15Lines(fileURL)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", fileURL, err)
		return
	}

	if fileContent == homepageContent {
		fmt.Printf("%s: False Positive \n", fileURL)
	} else {
		fmt.Printf("%s: Content Found.\n", fileURL)
	}
}

func main() {
	homepageURL := flag.String("s", "", "The homepage URL to check against")
	flag.Parse()

	if *homepageURL == "" {
		fmt.Println("Please provide a homepage URL using the -s flag.")
		os.Exit(1)
	}

	paths := []string{
		"/.wwwacl",
		"/.www_acl",
		"/.htpasswd",
		"/.access",
		"/.addressbook",
		"/.bashrc",
		"/.forward",
		"/.history",
		"/.htaccess",
		"/.lynx_cookies",
		"/.passwd",
		"/.pinerc",
		"/.plan",
		"/.proclog",
		"/.procmailrc",
		"/.profile",
		"/.rhosts",
		"/.ssh",
		"/.nsconfig",
		"/.gitignore",
		"/.hgignore",
		"/.env",
		"/.dockerignore",
	}

	homepageContent, err := fetchFirst15Lines(*homepageURL)
	if err != nil {
		fmt.Printf("Error fetching homepage: %v\n", err)
		return
	}

	for _, path := range paths {
		fullURL := *homepageURL + path
		checkForFalsePositive(fullURL, homepageContent)
	}
}