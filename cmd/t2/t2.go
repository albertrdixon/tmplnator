package main

import (
	"flag"
	"fmt"
	"github.com/albertrdixon/tmplnator/template"
	"github.com/albertrdixon/tmplnator/version"
	"io/ioutil"
	"os"
	"sync"
)

func exitErr(err error) {
	fmt.Printf("There's been a problem: %v", err)
	os.Exit(1)
}

func main() {
	var (
		V       = flag.Bool("version", false, "Print version")
		src     = flag.String("td", "", "template dir to render")
		del     = flag.Bool("d", false, "Remove templates after processing")
		pre     = flag.String("P", "", "Dir prefix")
		threads = flag.Int("T", 4, "Render threads")
	)

	flag.Parse()
	if *V {
		fmt.Printf("Version: %s\n", version.RuntimeVersion(version.CodeVersion, version.Build))
		os.Exit(0)
	}

	if *src == "" {
		fmt.Printf("Sorry, you need to provide the src dir.\n")
		os.Exit(1)
	}
	if *threads < 1 {
		*threads = 1
	}

	s, p := *src, *pre
	var err error

	if p == "tmp" {
		if p, err = ioutil.TempDir("", "tnator"); err != nil {
			fmt.Printf("Couldn't get tmp dir.")
			exitErr(err)
		}
		fmt.Println(pre)
	}

	fmt.Printf("==> Parsing Templates in %q\n", s)
	tStack, err := template.ParseDirectory(s, p)
	if err != nil {
		fmt.Printf("Failed to parse templates.")
		exitErr(err)
	}

	var wg sync.WaitGroup
	wg.Add(*threads)
	for i := 0; i < *threads; i++ {
		go func() {
			defer wg.Done()
			for tStack.Len() > 0 {
				if t, ok := tStack.Pop().(*template.Template); ok {
					if err := t.Write(); err == nil {
						if *del {
							os.Remove(t.Src)
						}
					} else {
						fmt.Printf("Error writing template: %v", err)
					}
				} else {
					fmt.Printf("Could not cast stack item as template: %v", t)
				}
			}
		}()
	}

	wg.Wait()
	os.Exit(0)
}
