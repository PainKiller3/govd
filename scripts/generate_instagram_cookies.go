package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	asciiLetters   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	asciiLowercase = "abcdefghijklmnopqrstuvwxyz"
	digits         = "0123456789"
)

func randomString(length int, charset string) string {
	b := make([]byte, length)
	max := big.NewInt(int64(len(charset)))
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		b[i] = charset[num.Int64()]
	}
	return string(b)
}

type CookieEntry struct {
	Comment string
	Fields  []string
}

func generateGuestCookies() []CookieEntry {
	csrfToken := randomString(32, asciiLetters+digits)
	igDid := randomString(24, asciiLowercase+digits)
	mid := randomString(24, asciiLetters+digits)

	expires := strconv.FormatInt(time.Now().Unix()+31536000, 10)

	return []CookieEntry{
		{Comment: "# Netscape HTTP Cookie File"},
		{Comment: "# http://curl.haxx.se/rfc/cookie_spec.html"},
		{Comment: "# This is a generated file!  Do not edit."},
		{Comment: ""},
		{Fields: []string{".instagram.com", "TRUE", "/", "TRUE", expires, "csrftoken", csrfToken}},
		{Fields: []string{".instagram.com", "TRUE", "/", "TRUE", expires, "ig_did", igDid}},
		{Fields: []string{".instagram.com", "TRUE", "/", "TRUE", expires, "mid", mid}},
		{Fields: []string{".instagram.com", "TRUE", "/", "TRUE", expires, "ig_nrcb", "1"}},
		{Fields: []string{".instagram.com", "TRUE", "/", "TRUE", expires, "wd", "1920x1080"}},
		{Fields: []string{".instagram.com", "TRUE", "/", "TRUE", expires, "dpr", "2"}},
	}
}

func writeNetscapeCookies(filePath string, cookies []CookieEntry) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	var sb strings.Builder
	for _, entry := range cookies {
		if len(entry.Fields) == 0 {
			sb.WriteString(entry.Comment + "\n")
		} else {
			sb.WriteString(strings.Join(entry.Fields, "\t") + "\n")
		}
	}

	if err := os.WriteFile(filePath, []byte(sb.String()), 0644); err != nil {
		return err
	}

	fmt.Printf("[+] Successfully wrote Netscape cookie file to: %s\n", filePath)
	return nil
}

func main() {
	targetPath := filepath.Join("private", "cookies", "instagram.txt")
	fmt.Println("==========================================")
	fmt.Println("      Instagram Cookie File Generator     ")
	fmt.Println("==========================================")

	cookies := generateGuestCookies()
	if err := writeNetscapeCookies(targetPath, cookies); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing cookies: %v\n", err)
		os.Exit(1)
	}
}
