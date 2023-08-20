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
	lineFlag, wordFlag, tokenFlag, charFlag bool
	wordRegex                               = regexp.MustCompile(`\S+`)
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "tokenutil",
		Short: "Utility for token operations",
	}

	var countCmd = &cobra.Command{
		Use:   "count [file...]",
		Short: "Count lines, words, tokens, and characters in a file or from stdin",
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

	rootCmd.AddCommand(countCmd)
	countCmd.Flags().BoolVarP(&lineFlag, "lines", "l", false, "Count lines")
	countCmd.Flags().BoolVarP(&wordFlag, "words", "w", false, "Count words")
	countCmd.Flags().BoolVarP(&tokenFlag, "tokens", "t", true, "Count tokens")
	countCmd.Flags().BoolVarP(&charFlag, "chars", "c", false, "Count characters")
	countCmd.PersistentFlags().StringVarP(&encoding, "model", "m", "gpt-3.5-turbo", "model name to encode for")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
