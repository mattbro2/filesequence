//Package to reduce listing of files to File_seq objects
package reducers

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/mattbro2/fileseq/seq_definition"
)

// Struct for the File_seq object, contains the following:
// Base is the file with the file numbers replaced with '@'
//    Base detects files with numbering schema such as:
//       filename.1001.jpg
//       filename_1001.jpg
//       filename 1001.jpg
// File_num is a map of map[<int of file number>]<string of file num>
// File_list is ordered array of file number integers
// F_seq is the condensed listing of file sequence with file numbers listed
//  inside brackets such as:  test.[001-003].jpg
//  Non continuous sequence:  nonseq.[01,03-05,10,15-17].jpg
type File_seq struct {
	Base      string
	File_num  map[int]string
	File_list []int
	F_seq     string
	Force     bool
}

//Function to take listing of files and create the base and file list
func ReduceBase(files []string) (map[string]map[int]string, error) {
	fmt.Println()
	bases := make(map[string]map[int]string)

	seq_def, err := seq_definition.SeqDefinition()
	if err != nil {
		return bases, err
	}

	fr_number_regex, reg_err := regexp.Compile(seq_def.ReducerRegex)
	if reg_err != nil {
		return bases, reg_err
	}

	for _, f := range files {
		frnum := fr_number_regex.FindStringSubmatch(f)
		if len(frnum) == 0 {
			bases[f] = make(map[int]string)
			bases[f][0] = "0"
		} else {
			repl := strings.Replace(f, frnum[1], `@`, 1)
			ifrnum, _ := strconv.Atoi(frnum[1])
			_, ok := bases[repl]
			if !ok {
				bases[repl] = make(map[int]string)
			}
			bases[repl][ifrnum] = frnum[1]
		}
	}
	return bases, reg_err
}

//Function to retrieve the base and file list and convert it to a File_seq obj
func ReduceFileseq(bases map[string]map[int]string) ([]File_seq, error) {
	var file_seqs []File_seq
	for f, v := range bases {
		var keys []int
		for key, _ := range v {
			keys = append(keys, key)
		}
		sort.Ints(keys)
		f_seq_range_format := "["
		cont := ""
		for i, k := range keys {
			if keys[i] == keys[0] {
				f_seq_range_format += v[keys[i]]
				continue
			}
			if k-1 == keys[i-1] {
				cont = "-" + v[keys[i]]
				continue
			} else {
				f_seq_range_format += fmt.Sprintf("%s,%s", cont, v[keys[i]])
				cont = ""
			}
		}
		f_seq_range_format = fmt.Sprintf("%s%s]", f_seq_range_format, cont)
		f_seq := strings.Replace(f, `@`, f_seq_range_format, 1)

		//If the file matches the sequence regex, but there is only one
		if len(v) == 1 {
			f_seq = strings.Replace(f, `@`, v[keys[0]], 1)
		}

		fs := File_seq{
			Base:      f,
			File_list: keys,
			File_num:  v,
			F_seq:     f_seq,
		}
		file_seqs = append(file_seqs, fs)
	}
	return file_seqs, nil
}
