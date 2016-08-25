//Package fileseq will compact a sequence of files into a single line representing the sequence
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mattbro2/fileseq/commands"
	"github.com/mattbro2/fileseq/core"
	"github.com/mattbro2/fileseq/filesys"

	"github.com/daviddengcn/go-colortext"
)

//main entry point for the file sequencer comand line.  Will print to stdout
func main() {
	options := commands.InitCommands(os.Stdout)
	reader := bufio.NewReader(os.Stdin)
	//If the user wants a file list from a File_seq object
	if options.Reverse != "" {
		fseq, rvseq_err := core.ReverseSeqMain(options.Reverse)
		if rvseq_err != nil {
			fmt.Printf("Unable to create sequence from %s - %s\n", options.Reverse, rvseq_err)
			os.Exit(1)
			return
		}

		reverse, rev_err := core.ReverseMain(fseq)
		if rev_err != nil {
			fmt.Printf("Unable to list files from sequence %s - %s\n", options.Reverse, rev_err)
			os.Exit(1)
			return
		}

		for _, x := range reverse {
			isfile, _ := filesys.IsFile(x)
			if isfile {
				fmt.Println(x)
				continue
			}
			if options.Nocolor {
				fmt.Println(x)
				continue
			}
			ct.Foreground(ct.Red, true)
			fmt.Println(x)
			ct.ResetColor()

		}
		return
	}

	//Copy one File_seq to another
	if options.Copy != "" {
		fs_split := strings.Split(options.Copy, "::")
		if len(fs_split) != 2 {
			fmt.Printf("-c param %s not two fseqs separated by '::'\n", options.Copy)
			os.Exit(1)
			return
		}
		err := core.CopySeqMain(fs_split[0], fs_split[1], options.Force, options.Verbose)
		if err != nil {
			fmt.Printf("Unable to copy files %s\n", err)
			os.Exit(1)
			return
		}
		return
	}

	//Move from one file_seq to another
	if options.Move != "" {
		fs_split := strings.Split(options.Move, "::")
		if len(fs_split) != 2 {
			fmt.Printf("-c param %s not two fseqs separated by '::'\n", options.Move)
			os.Exit(1)
			return
		}
		err := core.MoveSeqMain(fs_split[0], fs_split[1], options.Force, options.Verbose)
		if err != nil {
			fmt.Printf("Unable to move files %s\n", err)
			os.Exit(1)
		}
		return
	}

	//Renumber a sequence of files, will allow overwriting
	if options.Reseq != "" {
		fs_split := strings.Split(options.Reseq, "::")
		if len(fs_split) != 2 {
			fmt.Printf("-q param %s not two fseqs separated by '::'\n", options.Reseq)
			os.Exit(1)
			return
		}
		err := core.ReSeqMain(fs_split[0], fs_split[1], options.Verbose)
		if err != nil {
			fmt.Printf("Unable to resequence files %s\n", err)
			os.Exit(1)
		}
		return
	}

	//Delete a file seq
	if options.Delete != "" {
		if !options.Force {
			fmt.Println("This will remove your data, are you sure? [y/n]: ")
			response, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("error occurred %s", err)
				os.Exit(1)
				return
			}
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" {
				fmt.Println("Not continuing with delete, reponse was not 'y'")
				os.Exit(1)
				return
			}
		}
		err := core.DeleteSeqMain(options.Delete, options.Force, options.Verbose)
		if err != nil {
			fmt.Printf("Error occurred %s ", err)
			os.Exit(1)
			return
		}
		return
	}

	//Default behavior of doing a file_seq listing
	file_seqs, err := core.ListMain(options.Curdir, options.Verbose)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	var fmt_seqs []string

	for _, x := range file_seqs {
		fmt_seqs = append(fmt_seqs, x.F_seq)
	}

	sort.Strings(fmt_seqs)

	for _, x := range fmt_seqs {
		fmt.Println(x)
	}

	return
}
