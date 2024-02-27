package main

import(
	"fmt"
	"io/fs"
	"os"
	"errors"
	S "strings"
	HB "github.com/fbaube/humanbytes"
)

// func DirFS(dir string) fs.FS
// DirFS returns an fs.FS for the file tree rooted at "dir".
// The result implements
//  - StatFS:     Stat    (name string) (FileInfo, error)  // error: *PathError
//  - ReadDirFS:  ReadDir (name string) ([]DirEntry, error) // filenames sorted
//  - ReadFileFS: ReadFile(name string) ([]byte, error)
//    ReadFile reads the named file and returns its contents.
//    A successful call returns a nil error, not io.EOF.

var fsfs fs.FS

var lsTmFmt = "_2 Jan 2006 15:04"

func main() {

     fsfs = os.DirFS(".")
     _, e := Walk(fsfs)
     fmt.Printf("ERROR: %#v \n", e)
     // fmt.Printf("ITEMS: %#v \n", ss)
     }


func Walk(cab fs.FS) ([]string, error) {
	var entries []string
	// var e error
	var ok1, ok2, ok3 bool
	// var stfs fs.StatFS
	// var rdfs fs.ReadDirFS
	// var rffs fs.ReadFileFS

	/*stfs*/ _, ok1 = fsfs.(fs.StatFS)
	/*rdfs*/ _, ok2 = fsfs.(fs.ReadDirFS)
	/*rffs*/ _, ok3 = fsfs.(fs.ReadFileFS)
	if !(ok1&&ok2&&ok3) {
	   println("CONVERSION PROBLEM")
	   return nil, errors.New("OOPS")
	   }

	err := fs.WalkDir(cab, ".",
	func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} /*
		fi, e := stfs.Stat(path)
		if e != nil {
		     fmt.Printf("\nSTAT ERROR: %s : %s \n\n",
		     	path, e.Error())
		     } /* else {
		     fmt.Printf("\nFileInfo: %#v \n", fi)
		     } */
		// Info() (FileInfo, error)
		var info fs.FileInfo
		info, _ = d.Info()
		/* type FileInfo interface {
		   Name() string       // base name of the file
		   Size() int64        // length in bytes for regular files; system-dependent for others
		   Mode() FileMode     // file mode bits
		   ModTime() time.Time // modification time
		   IsDir() bool        // abbreviation for Mode().IsDir()
		   Sys() any           // underlying data source (can return nil)
		   } */
		var optSlash string 
		if info.IsDir() {
		   optSlash = "/"
		   }
		fmt.Printf("%s %4s  %s  %s%s\n",
			info.Mode(), HB.SizeLS(int(info.Size())), 
			info.ModTime().UTC().Format(lsTmFmt),
			path, optSlash) // info.Name())
		// fmt.Println(path, " :: ", fs.FormatFileInfo(info))

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
		// append the entry to the list
		entries = append(entries, path)
		// nil tells walk to continue
		return nil
	})
	return entries, err
}



/*
func DirFS(dir string) fs.FS

DirFS returns an fs.FS for the file tree rooted at "dir".

The directory dir must not be "".

The result implements
- io/fs.StatFS:     Stat    (name string) (FileInfo, error)  // error: *PathError
- io/fs.ReadDirFS:  ReadDir (name string) ([]DirEntry, error) // filenames sorted
- io/fs.ReadFileFS: ReadFile(name string) ([]byte, error)

ReadFile reads the named file and returns its contents.

A successful call returns a nil error, not io.EOF, cos
ReadFile reads the whole file, so the expected EOF from
the final Read is not an error to be reported.

The caller may modify the returned byte slice.
This method should return a copy of the underlying data.

= = = = = =

https://robthorne-26852.medium.com/a-tale-of-two-file-systems-in-go-b749038c7373

fSys := os.Dir("my-files")

// Want to read my-files/list.txt?
contents, err := fs.ReadFile(fSys, "list.txt")

// Read the directory?
items, err := fs.ReadDir(fSys, ".")
if err != nil {
    // handle the error...
}
for _, item := range items {
    fmt.Println(item.Name())
}

// Go Globbing for files?
matches, err := fs.Glob(fSys, "*.txt")
if err != nil {
    // guess we would bail here
}
for _, path := range matches {
    fmt.Println(path)
}

Serving files from this directory using 
http.FileServer() looks something like this:

    filesDir := os.DirFS("my-files")

    handler := http.FileServer(http.FS(filesDir))
    http.Handle("/", handler)

    log.Println("Serving static files at :5000")
    err := http.ListenAndServe(":5000", handler)
    if err != nil {
        log.Fatal(err)
    }

= = = = = =

You’d think that since embed.FS and os.DirFS() share
this same fs.FS interface, that you’d be able to use
them interexchangeably here. But sadly: YOU’D BE WRONG :-)

To fix the security holes you get with os.DirFS, there
are two approaches, both in my accompanying repo:
- Write middleware that filters by file name 
- Or, take the file system wrapper approach, and
  add that filtering logic to that code instead.

https://github.com/torenware/vite-go/blob/master/asset-server.go

// Open implements the fs.FS interface for wrapperFS
func (wrpr wrapperFS) Open(path string) (fs.File, error) {
	f, err := wrpr.FS.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	// CHECK THAT PATHS ARE KOSHER !!!!!

	if s.IsDir() {
		// Have an index file or go home!
		index := filepath.Join(path, "index.html")
		if _, err := wrpr.FS.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

= = = = = =

WALKING!

https://gopherguides.com/golang-fundamentals-book/14-files/fs

*/