package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/shapor/tiktoken-go"
	"github.com/spf13/cobra"
)

var (
	name                                    = "tiktokenutil"
	encoding                                string
	separatorFlag                           string
	lineFlag, wordFlag, tokenFlag, charFlag bool
	statFlag                                bool
	wordRegex                               = regexp.MustCompile(`\S+`)
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "tokenutil",
		Short: "Utility for token operations",
	}

	var countCmd = &cobra.Command{
		Use:   "count [file...]",
		Short: "Count lines, words, tokens, and characters in file(s) or from stdin",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			totalLines, totalWords, totalTokens, totalChars := 0, 0, 0, 0
			if len(args) == 0 {
				lines, words, tokens, chars, _ := tokenCount(os.Stdin)
				display("stdin", lines, words, tokens, chars)
			} else {
				for _, filePath := range args {
					file, err := os.Open(filePath)
					if err != nil {
						fmt.Println("Error opening file:", err)
						continue
					}

					lines, words, tokens, chars, _ := tokenCount(file)
					display(filePath, lines, words, tokens, chars)
					file.Close()

					totalLines += lines
					totalWords += words
					totalTokens += tokens
					totalChars += chars
				}
			}

			if len(args) > 1 {
				display("total", totalLines, totalWords, totalTokens, totalChars)
			}
			return nil
		},
	}

	var encodeCmd = &cobra.Command{
		Use:   "encode [file...]",
		Short: "Tokenizes and encodes file(s) or from stdin",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			totalTokens := 0
			if len(args) == 0 {
				totalTokens += encode(os.Stdin)
			} else {
				for _, filePath := range args {
					file, err := os.Open(filePath)
					if err != nil {
						fmt.Println("Error opening file:", err)
						continue
					}

					totalTokens += encode(file)
					file.Close()
				}
			}
			if statFlag {
				fmt.Fprintf(os.Stderr, "Encoded %v tokens.\n", totalTokens)
			}
			return nil
		},
	}

	rootCmd.AddCommand(countCmd)
	rootCmd.AddCommand(encodeCmd)
	rootCmd.PersistentFlags().StringVarP(&encoding, "model", "m", "gpt-3.5-turbo", "model name to encode for")
	countCmd.Flags().BoolVarP(&lineFlag, "lines", "l", false, "Count lines")
	countCmd.Flags().BoolVarP(&wordFlag, "words", "w", false, "Count words")
	countCmd.Flags().BoolVarP(&tokenFlag, "tokens", "t", true, "Count tokens")
	countCmd.Flags().BoolVarP(&charFlag, "chars", "c", false, "Count characters")
	encodeCmd.Flags().BoolVarP(&statFlag, "tokens", "t", false, "Output token count stats to stderr")
	encodeCmd.Flags().StringVarP(&separatorFlag, "separator", "s", "\n", "Separator string between tokens")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func encode(r io.Reader) int {
	bytes, _ := io.ReadAll(r)
	contents := string(bytes)
	tkm, err := tiktoken.EncodingForModel(encoding)
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
		return 0
	}
	token := tkm.Encode(contents, nil, nil)
	for n, id := range token {
		if n > 0 {
			fmt.Print(separatorFlag)
		}
		fmt.Print(id)
	}
	fmt.Println()
	return len(token)
}

func tokenCount(r io.Reader) (int, int, int, int, error) {
	bytes, _ := io.ReadAll(r)
	contents := string(bytes)
	lines := strings.Count(contents, "\n")
	words := len(wordRegex.FindAllString(contents, -1))
	chars := len(contents)
	tkm, err := tiktoken.EncodingForModel(encoding)
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
		return lines, words, 0, chars, err
	}
	token := tkm.Encode(contents, nil, nil)
	return lines, words, len(token), chars, nil
}

func display(label string, lines, words, tokens, chars int) {
	if lineFlag {
		fmt.Printf("%8d", lines)
	}
	if wordFlag {
		fmt.Printf("%8d", words)
	}
	if tokenFlag {
		fmt.Printf("%8d", tokens)
	}
	if charFlag {
		fmt.Printf("%8d", chars)
	}
	fmt.Printf(" %s\n", label)
}
