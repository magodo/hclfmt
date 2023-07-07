package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

const Version = "0.1.0"

var (
	write      = flag.Bool("w", false, "write result to (source) file instead of stdout")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to this file")
	version    = flag.Bool("version", false, "version information")
)

func main() {
	if err := realMain(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func realMain() error {
	flag.Usage = usage
	flag.Parse()

	if *version {
		fmt.Println(Version)
		return nil
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			return fmt.Errorf("creating cpu profile: %s\n", err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if flag.NArg() == 0 {
		if *write {
			return errors.New("error: cannot use -w with standard input")
		}

		return processFile("<standard input>", os.Stdin, os.Stdout, true)
	}

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		switch dir, err := os.Stat(path); {
		case err != nil:
			return err
		case dir.IsDir():
			walkDir(path)
		default:
			if err := processFile(path, nil, os.Stdout, false); err != nil {
				return err
			}
		}
	}

	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: hclfmt [flags] [path ...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func isHclFile(f os.FileInfo) bool {
	// ignore non-hcl files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".hcl")
}

func walkDir(path string) {
	filepath.Walk(path, visitFile)
}

func visitFile(path string, f os.FileInfo, err error) error {
	if err == nil && isHclFile(f) {
		err = processFile(path, nil, os.Stdout, false)
	}

	return err
}

// If in == nil, the source is the contents of the file with the given filename.
func processFile(filename string, in io.Reader, out io.Writer, stdin bool) error {
	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	res := hclwrite.Format(src)

	if *write {
		err = ioutil.WriteFile(filename, res, 0644)
	} else {
		_, err = out.Write(res)
	}

	return err
}
