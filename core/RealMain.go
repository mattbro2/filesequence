// Package core contains code to run the file sequencing
// This will create and return an array of File_seq objects
package core

import (
	"github.com/mattbro2/fileseq/expanders"
	"github.com/mattbro2/fileseq/filesys"
	"github.com/mattbro2/fileseq/reducers"
	"github.com/mattbro2/fileseq/seq_manip"
)

//Call the functions and return the data and errors
func ListMain(curdir string, verbose bool) ([]reducers.File_seq, error) {
	files, rec_err := filesys.Recurse(curdir, verbose)
	reduced, red_err := reducers.ReduceBase(files)
	file_seqs, fseq_err := reducers.ReduceFileseq(reduced)

	if rec_err != nil {
		return nil, rec_err
	}

	if red_err != nil {
		return nil, red_err
	}

	if fseq_err != nil {
		return nil, fseq_err
	}

	return file_seqs, nil
}

func ReverseMain(fs reducers.File_seq) ([]string, error) {
	files, err := expanders.Fseq_expand(fs)
	return files, err
}

func ReverseSeqMain(fs string) (reducers.File_seq, error) {
	fseq, err := expanders.Fseq_to_object(fs)
	return fseq, err
}

func CopySeqMain(fs string, fd string, force bool, verbose bool) error {
	err := seq_manip.CopySeq(fs, fd, force, verbose)
	return err
}

func MoveSeqMain(fs string, fd string, force bool, verbose bool) error {
	err := seq_manip.MoveSeq(fs, fd, force, verbose)
	return err
}

func ReSeqMain(fs string, fd string, verbose bool) error {
	err := seq_manip.ReSeq(fs, fd, verbose)
	return err
}

func DeleteSeqMain(fs string, force bool, verbose bool) error {
	err := seq_manip.DeleteSeq(fs, force, verbose)
	return err
}
