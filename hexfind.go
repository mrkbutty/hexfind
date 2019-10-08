/*
search files within a dir tree for files containing specified bytes

Usage: hexfind <directory> <hexbytes>

	<directory> = list of directories to process.
	<hexbytes>  = hex bytes to find, e.g. 1A34F8

  -d    Follow hidden dot directorys
  -q    No output apart from errors
  -v    Prints detailed operations

*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

var flagVerbose bool
var flagQuiet bool
var flagDotted bool
var dateFormat = time.RFC1123

func walktree(droot string, hexstr string) (int, error) {

	droot, err := filepath.Abs(droot)
	if err != nil {
		log.Fatal(err)
	}
	drootinfo, err := os.Stat(droot)
	if err != nil {
		log.Fatal(err)
	}
	if !drootinfo.Mode().IsDir() {
		log.Fatal("Is not a directory")
	}

	count := 0

	if flagVerbose {
		fmt.Printf("Processing %s\n", droot)
	}
	flist, err := ioutil.ReadDir(droot)
	if err != nil {
		log.Fatal(err)
	}

	//search files
	checked := 0
	for _, i := range flist {
		if i.Mode().IsRegular() {
			if flagVerbose {
				fmt.Printf("Searching %s\n", filepath.Join(droot, i.Name()))
			}

		}
		checked++
	}

	//Next look in subdirectories
	for _, i := range flist {

		if i.Name()[0] == '.' {
			if !flagDotted {
				if flagVerbose {
					fmt.Println("Skipping hidden", i.Name())
				}
				continue
			}
		}
		count++

		//fmt.Println(i.Name(), i.Mode().IsDir())

		if i.Mode().IsDir() {
			dapath, err := filepath.Abs(filepath.Join(droot, i.Name()))
			count2, err := walktree(dapath, hexstr)
			if err != nil {
				log.Fatal(err)
			}
			count += count2
		}
	}

	return count, nil
}

func main() {
	flag.BoolVar(&flagVerbose, "v", false, "Prints detailed operations")
	flag.BoolVar(&flagQuiet, "q", false, "No output apart from errors")
	flag.BoolVar(&flagDotted, "d", false, "Follow hidden dot directorys")
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Printf("Usage: %s [options] <directory> <hexstring>\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		return
	}

	dirname := flag.Arg(0)
	hexstr := flag.Arg(1)

	start := time.Now()

	total, err := walktree(dirname, hexstr)
	if err != nil {
		log.Fatal(err)
	}
	if !flagQuiet {
		elapsed := time.Since(start)
		fmt.Printf("Processed %d total items in %v\n", total, elapsed)
	}
}
