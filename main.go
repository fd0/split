package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/restic/chunker"
	"github.com/spf13/pflag"
)

var opts = struct {
	Polynomial string
	MinSize    uint
	MaxSize    uint
	Verbose    bool
	Template   string
	Dir        string
	InputFile  string
}{}

func warn(msg string, args ...interface{}) {
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprintf(os.Stderr, msg, args...)
}

func die(msg string, args ...interface{}) {
	warn(msg, args...)
	os.Exit(1)
}

func v(msg string, args ...interface{}) {
	if !opts.Verbose {
		return
	}

	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Printf(msg, args...)
}

// DefaultPolynomial is used if no other polynomial is specified in the flags
const DefaultPolynomial = "3DA3358B4DC173"

func main() {
	flags := pflag.NewFlagSet("split", pflag.ContinueOnError)
	flags.StringVarP(&opts.Polynomial, "polynomial", "p", DefaultPolynomial, "Use polynomial `p` for splitting (hex notation, no prefix)")
	flags.UintVarP(&opts.MinSize, "min-size", "l", chunker.MinSize, "Set minimal chunk size to `n` bytes")
	flags.UintVarP(&opts.MaxSize, "max-size", "u", chunker.MaxSize, "Set maximal chunk size to `n` bytes")
	flags.BoolVarP(&opts.Verbose, "verbose", "v", false, "Be verbose")
	flags.StringVarP(&opts.Template, "template", "t", "split-%03d", "Use `s` as the (printf-style) template for output files")
	flags.StringVarP(&opts.Dir, "output", "o", ".", "Write files to directory `dir` instead of the current directory")
	flags.StringVarP(&opts.InputFile, "input", "i", "", "Read from `file` instead of stdin")

	err := flags.Parse(os.Args)
	if err == pflag.ErrHelp {
		os.Exit(0)
	}

	var input = os.Stdin
	if opts.InputFile != "" {
		f, err := os.Open(opts.InputFile)
		if err != nil {
			die("unable to open input file: %v", opts.InputFile)
		}

		input = f
	}

	if opts.MinSize >= opts.MaxSize {
		die("invalid settings: minimal size larger or equal to maximal size, exiting")
	}

	p, err := strconv.ParseUint(opts.Polynomial, 16, 64)
	if err != nil {
		die("unable to parse polynomial from hex: %v", err)
	}

	pol := chunker.Pol(p)
	if !pol.Irreducible() {
		die("invalid polynomial specified (forgot the '0x' hexadecimal prefix?), exiting")
	}

	var buf []byte
	var filenum int
	var bytes uint64
	c := chunker.NewWithBoundaries(input, pol, opts.MinSize, opts.MaxSize)
	for {
		chunk, err := c.Next(buf)
		if err != nil {
			if err != io.EOF {
				warn("error: %v", err)
			}
			break
		}

		filename := filepath.Join(opts.Dir, fmt.Sprintf(opts.Template, filenum))
		err = ioutil.WriteFile(filename, chunk.Data, 0644)
		if err != nil {
			die("unable to write data: %v", err)
		}

		v("next chunk offset %d, %d bytes written to %v\n", chunk.Start, chunk.Length, filename)

		buf = chunk.Data
		filenum++
		bytes += uint64(chunk.Length)
	}

	v("wrote %d bytes to %d files", bytes, filenum)

	err = input.Close()
	if err != nil {
		die("unable to close input: %v", err)
	}
}
