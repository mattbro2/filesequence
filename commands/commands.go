// Package commands contains the code to process the service flags and commands
package commands

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mattbro2/fileseq/filesys"
)

type Options struct {
	Curdir  string
	Reverse string
	Copy    string
	Move    string
	Delete  string
	Reseq   string
	Nocolor bool
	Force   bool
}

// InitCommands parses command line flags
// updates global options
func InitCommands(out io.Writer) Options {
	printUsage := false
	curdir := filesys.Curdir()
	reverse := ""
	copyf := ""
	move := ""
	deletef := ""
	reseq := ""
	nocolor := false
	force := false

	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagset.SetOutput(out)
	flagset.BoolVar(&printUsage, "h", false, "Print Help")
	flagset.BoolVar(&printUsage, "help", false, "Print Help")
	flagset.StringVar(&curdir, "p", curdir, "Set directory to search")
	flagset.StringVar(&reverse, "r", "", "Take a F_seq and expand to list of files (offline files are printed to terminal in red)")
	flagset.StringVar(&copyf, "c", copyf, "Copy ie: fseq1.[01-10].jpg:fseq2.[01-10].jpg - cannot be a resequencing of same files")
	flagset.StringVar(&move, "m", move, "Move ie: fseq1.[01-10].jpg:fseq2.[01-10].jpg\n\t"+
		"Move will result in original files being renamed. Source and dest must be different")
	flagset.StringVar(&reseq, "q", reseq, "Renumber a sequence of files ie: fseq1.[001-009].jpg:fseq1.[101-109].jpg")
	flagset.StringVar(&deletef, "d", deletef, "Remove all files in sequence")
	flagset.BoolVar(&nocolor, "n", false, "Do not add colors to printed output")
	flagset.BoolVar(&force, "f", force, "Allow for overwriting of exiting files (destination cannot overwrite source unless using 'q' flag)")
	flagset.Parse(os.Args[1:])

	if printUsage {
		fmt.Fprintf(out, "\nA tool to condense sequences of files into a compact format. ie: fseq1.[1-10].jpg \n"+
			"where the first file in sequence is fseq1.01.jpg\n\n"+
			"Will run from the current directory, or can be passed a directory to search and \n"+
			"it will recursively gather all the files and format them if possible.\n\n"+
			"Additionally, it will copy, move or delete sequences of files.\n\n"+
			"Lastly, you may pass the tool a 'reverse' lookup where you give it a compacted file sequence\n"+
			"and it will turn it into a file_seq object and give you a list of files, indicating if any are\n"+
			"offline in red.\n\n"+
			"File sequences are detected by a file that ends with either a '.#.ext, ' #.ext', or '_#.ext' \n"+
			"This may not match your naming convnetion, see the documentation in repo. \n\n"+
			"Sequential files are noted by a dash '-' and non sequential are noted by commas ','."+
			"\n\n%s [options]\n\n  options\n  -------\n\n", os.Args[0])
		flagset.PrintDefaults()
		fmt.Fprintln(out, "")
		os.Exit(1)
	}

	o := Options{
		Curdir:  strings.TrimRight(curdir, "/"),
		Reverse: reverse,
		Copy:    copyf,
		Move:    move,
		Delete:  deletef,
		Reseq:   reseq,
		Nocolor: nocolor,
		Force:   force,
	}

	return o
}
