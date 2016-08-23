//Package expanders used to generate Fseq objects from a Fseq listing
// or to take an fseq and return the expanded list of files
package expanders

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/mattbro2/fileseq/filesys"
	"github.com/mattbro2/fileseq/reducers"
)

// Function to take a File_seq listing ie: "test.[001-005].jpg" and create a
// File_seq object out of it.
func Fseq_to_object(files string) (reducers.File_seq, error) {
	fs_regex, _ := regexp.Compile(`.*[\.\_\ ](\[[0-9-,]+\])\.`)
	fs_listing := fs_regex.FindStringSubmatch(files)
	file_num := make(map[int]string)
	var file_list []int

	//if the listing is not a sequence
	if len(fs_listing) == 0 {
		isfile, err := filesys.IsFile(files)
		if err != nil {
			return reducers.File_seq{}, err
		}

		if isfile != true {
			return reducers.File_seq{}, errors.New(files + " is not a file or sequence of files")
		}

		file_num[0] = "0"
		file_list = append(file_list, 0)
		fseq := reducers.File_seq{
			Base:      files,
			File_num:  file_num,
			File_list: file_list,
			F_seq:     files,
		}
		return fseq, nil
	}

	//if the listing is a sequence ie: test.[001-003].jpg
	repl := strings.Replace(files, fs_listing[1], `@`, 1)
	num_regex, _ := regexp.Compile(`([0-9]+|-|,)`)
	num_results := num_regex.FindAllStringSubmatch(fs_listing[1], -1)
	fp := len(num_results[0][0])
	fp_string := fmt.Sprintf("%%0%vd", fp)

	for index, reg_slice := range num_results {
		if index == 0 {
			int, _ := strconv.Atoi(reg_slice[0])
			file_num[int] = reg_slice[0]
			file_list = append(file_list, int)
			continue
		}

		if reg_slice[0] == "," {
			continue
		}

		if reg_slice[0] == `-` {
			start, _ := strconv.Atoi(num_results[index-1][0])
			end, _ := strconv.Atoi(num_results[index+1][0])
			for n := start + 1; n < end; n++ {
				file_num[n] = fmt.Sprintf(fp_string, n)
				file_list = append(file_list, n)
			}
			continue
		}

		int, _ := strconv.Atoi(reg_slice[0])
		file_num[int] = reg_slice[0]
		file_list = append(file_list, int)
	}

	fseq := reducers.File_seq{
		Base:      repl,
		File_num:  file_num,
		File_list: file_list,
		F_seq:     files,
	}

	return fseq, nil
}

//Given a File_seq object, expand to the list of files in sequence
func Fseq_expand(fs reducers.File_seq) ([]string, error) {
	var files []string

	for _, f := range fs.File_list {
		file := strings.Replace(fs.Base, `@`, fs.File_num[f], 1)
		files = append(files, file)
	}

	return files, nil
}
