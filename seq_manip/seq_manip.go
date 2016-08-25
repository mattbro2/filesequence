//Package to manage manipulation of sequences, ie copy, move, delete
package seq_manip

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"encoding/hex"

	"github.com/mattbro2/fileseq/expanders"
	"github.com/mattbro2/fileseq/filesys"
	"github.com/mattbro2/fileseq/reducers"
)

//Copy one sequence of files to another, force will allow overwriting.
//  Will perform md5 checksum validation post copy
func CopySeq(fs string, fd string, force bool, verbose bool) error {
	fs_source, fs_err := expanders.Fseq_to_object(fs)
	if fs_err != nil {
		return fs_err
	}
	fs_dest, fd_err := expanders.Fseq_to_object(fd)
	if fd_err != nil {
		return fd_err
	}

	files_source, files_dest, fs_err := FormatFileLists(fs_source, fs_dest, force)
	if fs_err != nil {
		return fs_err
	}

	mk_err := MakeDir(fd)
	if mk_err != nil {
		return mk_err
	}

	for i, _ := range files_source {
		in, err := os.Open(files_source[i])
		if err != nil {
			return err
		}
		defer in.Close()
		source_md5, hash_err := hash_file_md5(in)
		if hash_err != nil {
			return fmt.Errorf("Unable to generate checksum for source: %v", hash_err)
		}

		out, err := os.Create(files_dest[i])
		if err != nil {
			return err
		}
		if verbose {
			fmt.Printf("%s -> %s\n", files_source[i], files_dest[i])
		}
		if _, err := io.Copy(out, in); err != nil {
			return err
		}
		sync_error := out.Sync()
		if sync_error != nil {
			return fmt.Errorf("Error syncing file %v", sync_error)
		}

		dest_md5, hash_err := hash_file_md5(out)
		if hash_err != nil {
			return fmt.Errorf("Unable to generate checksum for source: %v", hash_err)
		}

		close_err := out.Close()
		if close_err != nil {
			return close_err
		}

		if source_md5 != dest_md5 {
			if verbose {
				fmt.Printf("source checksum %v", source_md5)
				fmt.Printf("dest checksum %v", dest_md5)
			}
			DeleteSeq(fd, true, verbose)
			return errors.New("Destination files are not valid, backing out of copy.\n" +
				"Please validate destination to ensure it is writable, ie disk full.\n" +
				"Or possibly source file is larger than destination FS permits.")
		}
	}
	return nil
}

//Rename one sequence to another (not copy).  Original file names will not
//  exist after the move
func MoveSeq(fs string, fd string, force bool, verbose bool) error {
	fs_source, fs_err := expanders.Fseq_to_object(fs)
	if fs_err != nil {
		return fs_err
	}
	fs_dest, fd_err := expanders.Fseq_to_object(fd)
	if fd_err != nil {
		return fd_err
	}

	mk_err := MakeDir(fd)
	if mk_err != nil {
		return mk_err
	}

	files_source, files_dest, fs_err := FormatFileLists(fs_source, fs_dest, force)
	if fs_err != nil {
		return fs_err
	}

	for i, _ := range files_source {
		if verbose {
			fmt.Printf("%s -> %s\n", files_source[i], files_dest[i])
		}
		mv_err := os.Rename(files_source[i], files_dest[i])
		if mv_err != nil {
			return mv_err
		}
	}
	return nil
}

//Renumber a sequence of files, performs a copy to a temp dir, deletes the
//  original files then renames them to the renumbered sequence
func ReSeq(fs string, fd string, verbose bool) error {
	fs_source, fs_err := expanders.Fseq_to_object(fs)
	if fs_err != nil {
		return fs_err
	}
	fs_dest, fd_err := expanders.Fseq_to_object(fd)
	if fd_err != nil {
		return fd_err
	}

	if fs_source.Base != fs_dest.Base {
		return errors.New("Source and destination must be the same name and location\n" +
			"This option is only to renumber the files in place.\n" +
			"You should use copy or move instead.\n")
	}

	temp_dir := os.TempDir()
	temp_file := filepath.Base(fd)
	temp_fd := fmt.Sprintf("%s/%s", temp_dir, temp_file)

	cperr := CopySeq(fs, temp_fd, true, verbose)
	if cperr != nil {
		return fmt.Errorf("Unable to create temp files: %v\n", cperr)
	}

	rmerr := DeleteSeq(fs, true, verbose)
	if rmerr != nil {
		return fmt.Errorf("Unable to remove original files: %v\n", rmerr)
	}

	mverr := MoveSeq(temp_fd, fd, false, verbose)
	if mverr != nil {
		return fmt.Errorf("Unable to move renumbered files back to original location\nError was %v\nYou may recover original files here: %s", mverr, temp_fd)
	}
	return nil

}

//Delete the files from disk
func DeleteSeq(fs string, force bool, verbose bool) error {
	fs_source, fs_err := expanders.Fseq_to_object(fs)
	if fs_err != nil {
		return fs_err
	}
	files_source, files_err := expanders.Fseq_expand(fs_source)
	if files_err != nil {
		return files_err
	}

	if !force {
		for _, x := range files_source {
			isfile, _ := filesys.IsFile(x)
			if !isfile {
				return errors.New(fs_source.F_seq + " files to delete are not completely online\n")
			}
		}
	}

	for _, x := range files_source {
		if verbose {
			fmt.Printf("deleting %s\n", x)
		}
		rm_err := os.Remove(x)
		if rm_err != nil {
			return rm_err
		}
	}
	return nil
}

//Abstraction of making directory and checking for errors
func MakeDir(pth string) error {
	dest_path := filepath.Dir(pth)
	mk_err := os.MkdirAll(dest_path, 0777)
	if mk_err != nil {
		return mk_err
	}
	return nil
}

//Take in File_seq objects and return slices of individual files, also check for inconsistencies between file_seqs
//  such as differet lengths, source files being offline or destition files being online.  Tool does not allow for overwriting
func FormatFileLists(fs_source reducers.File_seq, fs_dest reducers.File_seq, force bool) ([]string, []string, error) {
	if len(fs_source.File_list) != len(fs_dest.File_list) {
		return []string{}, []string{}, errors.New(fs_source.F_seq + " and " + fs_dest.F_seq +
			" do not contain the same number of files")
	}

	files_source, files_err := expanders.Fseq_expand(fs_source)
	if files_err != nil {
		return []string{}, []string{}, files_err
	}
	for _, x := range files_source {
		isfile, _ := filesys.IsFile(x)
		if !isfile {
			return []string{}, []string{}, errors.New(fs_source.F_seq + " source is not completely online\n")
		}
	}

	files_dest, files_err := expanders.Fseq_expand(fs_dest)
	if files_err != nil {
		return []string{}, []string{}, files_err
	}
	if fs_source.Base == fs_dest.Base {
		force = true
	}
	if !force {
		for _, x := range files_dest {
			isfile, _ := filesys.IsFile(x)
			if isfile {
				return []string{}, []string{}, errors.New(fs_dest.F_seq + " some or all destination files already exist\n")
			}
		}
	}
	return files_source, files_dest, nil
}

//Create a md5 hash, used for validating copy
func hash_file_md5(file *os.File) (string, error) {
	var md5string string
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return md5string, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	md5string = hex.EncodeToString(hashInBytes)
	return md5string, nil

}
