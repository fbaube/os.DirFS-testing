package main

import(
	"fmt"
	"io/fs"
	"os"
	"errors"
	S "strings"
	SU "github.com/fbaube/stringutils"
)

// func DirFS(dir string) fs.FS
// DirFS returns an fs.FS for the file tree rooted at "dir".
// The result implements
//  - StatFS:     Stat    (name string) (FileInfo, error)  // error: *PathError
//  - ReadDirFS:  ReadDir (name string) ([]DirEntry, error) // filenames sorted
//  - ReadFileFS: ReadFile(name string) ([]byte, error)
//    ReadFile reads the named file and returns its contents.
//    A successful call returns a nil error, not io.EOF.

func main() {

     var fsfs fs.FS
     fsfs = os.DirFS(".")
     _, listing, e := Walk(fsfs)
     fmt.Printf("ERROR: %#v \n", e)
     fmt.Printf(listing)
     // fmt.Printf("ITEMS: %#v \n", ss)
     }


func Walk(cab fs.FS) ([]string, string, error) {

     // Return values, built by the walk function 
	var entries []string
	var sb S.Builder

	var ok1, ok2, ok3 bool
	// var stfs fs.StatFS
	// var rdfs fs.ReadDirFS
	// var rffs fs.ReadFileFS

	/*stfs*/ _, ok1 = cab.(fs.StatFS)
	/*rdfs*/ _, ok2 = cab.(fs.ReadDirFS)
	/*rffs*/ _, ok3 = cab.(fs.ReadFileFS)
	if !(ok1&&ok2&&ok3) {
	   println("CONVERSION PROBLEM")
	   return nil, "", errors.New("OOPS")
	   }

	reterr := fs.WalkDir(cab, ".",
		func(path string, d fs.DirEntry, argerr error) error {
		
		if argerr != nil { return fmt.Errorf(
		   "Walk(%s): in-arg: %w", path, argerr) }
		   
		var e error
		
		/* fi, e = stfs.Stat(path)
		if e != nil {
		     return fmt.Errorf("Walk: Stat(%s): %w", path, e)
		     } */
		if d.IsDir() {
			name := d.Name()
			// if the directory is a dot return nil
			// this may be the root directory.
			if name == "." || name == ".." {
				return nil
			}
			// if the directory name is "testdata"
			// or it starts with "."
			// or it starts with "_"
			// then return filepath.SkipDir
			if name == "testdata" ||
			   S.HasPrefix(name, ".") || S.HasPrefix(name, "_") {
				return fs.SkipDir
			}
			return nil
		}
		// Info() (FileInfo, error)
		var fi fs.FileInfo
		fi, e = d.Info()
		if e != nil {
		   return fmt.Errorf("Walk: DirEntry.Info(%s): %w", path, e)
		   }
		var lsString string
		lsString = SU.LS_lh(fi, path)
		
		// append the entry to the list
		entries = append(entries, path)
		// append the file listing to the string builder
		sb.WriteString(lsString + "\n")
		// nil tells walk to continue
		return nil
	})
	return entries, sb.String(), reterr
}

