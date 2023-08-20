# tokenutil

`tokenutil` is a command line utility for tokenizing data using the tiktoken algorithm.

## Build & Run

To build:
```
go build
```

To run:
```
./tokenutil
```

See usage details below.

## Usage

### Count

`tokenutil count` is a utility mimicking the Unix `wc` command, with enhanced functionality to count tokens in addition to words, lines, and characters.

```
tokenutil count [options] [file...]
```

#### Options

- `-t`: Count tokens. This option is enabled by default.
- `-l`: Count lines.
- `-w`: Count words.
- `-c`: Count characters.

If no options are provided, `tokenutil` counts lines, words, characters, and tokens by default.

If no file is specified, `tokenutil` reads from standard input.

#### Examples

Count tokens of `file.txt`:

```
tokenutil count file.txt
```

Count only lines, words and tokens of `file1.txt` and `file2.txt`:

```
tokenutil count -ltw file1.txt file2.txt
```

Count tokens of input received from standard input (only `-t` defaults true):

```
echo "This is a test" | tokenutil count
```

### Encode

`tokenutil encode` will take input file(s) or read from stdin and output tokens to stdout.

```
tokenutil encode [options] [file...]
```

#### Options

- `-t`: Output total token stats to stderr.
- `-s`: Separator string between tokens (default: `\n`, one token per line)

If no file is specified, `tokenutil` reads from standard input.

#### Examples

Tokenize `file.txt` outputting one token id per line:

```
tokenutil encode file.txt
```

Tokenize files `file1.txt` and `file2.txt` and find the top ten tokens:

```
tokenutil encode -t file1.txt file2.txt | sort | uniq -c | sort -nr | head -10
```

## Tiktoken

In `tokenutil` tokenization is performed using the [tiktoken-go](https://github.com/shapor/tiktoken-go) library.

You can control the tokenizer by specifying an OpenAI model name with the -m parameter

= `-m`: Model name (e.g. `gpt-3.5-turbo` or `text-davinci-003`).

