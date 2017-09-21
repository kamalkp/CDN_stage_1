package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"flag"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

const usage = `
Usage:
  watch paths... [options]
Example:
  watch src --on-change 'make build'
Options:
      --on-change <arg>  Run command on any change
  -h, --halt             Exits on error (Default: false)
  -i, --interval <arg>   Run command once within this interval (Default: 1s)
  -n, --no-recurse       Skip subfolders (Default: false)
  -V, --version          Output the version number
  -q, --quiet            Suppress standard output (Default: false)
Intervals can be milliseconds(ms), seconds(s), minutes(m), or hours(h).
The format is the integer followed by the abbreviation.
`

var (
	last     time.Time
	interval time.Duration
	paths    []string
	err      error
)

// var opts struct {
//      // Help        bool   `short:"h" long:"help"       description:"Show this help message" default:false`
//      // Halt        bool   `short:"h" long:"halt"       description:"Exits on error (Default: false)" default:false`
//      Quiet       bool   `short:"q" long:"quiet"      description:"Suppress standard output (Default: false)" default:false`
//      Interval    string `short:"i" long:"interval"   description:"Run command once within this interval (Default: 1s)" default:"1s"`
//      NoRecurse   bool   `short:"n" long:"no-recurse" description:"Skip subfolders (Default: false)" default:false`
//      Version     bool   `short:"V" long:"version"    description:"Output the version number" default:false`
//      OnChange    string `long:"on-change"            description:"Run command on change."`
// }

var (
	Quiet   = flag.Bool("quiet", false, "Suppress standard output (Default: false)")
	Interval = flag.String("interval", "1s", "Run command once within this interval (Default: 1s)")
	NoRecurse = flag.Bool("no-recurse", false, "Skip subfolders (Default: false)")
	OnChange = flag.String("on-change", "echo hi", "Run command on change.")
	path = flag.String("paths", "testfile", "File to check on")
)

func init() {
	// args, err := flag.ParseArgs(&opts, os.Args)
	// if err != nil {
	//      fmt.Fprintln(os.Stderr, err)
	//      os.Exit(1)
	// }

	flag.Parse()
	paths, err = ResolvePaths(path)

	// if len(paths) <= 0 {
	//      fmt.Fprintln(os.Stderr, usage)
	//      os.Exit(2) // 2 for --help exit code
	// }

	interval, err = time.ParseDuration(*Interval)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	last = time.Now().Add(-interval)
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	done := make(chan bool)

	// clean-up watcher on interrupt (^C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		if !*Quiet {
			fmt.Fprintln(os.Stdout, "Interrupted. Cleaning up before exiting...")
		}
		watcher.Close()
		os.Exit(0)
	}()

	// process watcher events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if !*Quiet {
					fmt.Fprintln(os.Stdout, ev)
				}
				if time.Since(last).Nanoseconds() > interval.Nanoseconds() {
					last = time.Now()
					err = ExecCommand()
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						// if *Halt {
						//      os.Exit(1)
						// }
					}
				}
			case err := <-watcher.Error:
				fmt.Fprintln(os.Stderr, err)
				// if *Halt {
				//      os.Exit(1)
				// }
			}
		}
	}()

	// add paths to be watched
	for _, p := range paths {
		err = watcher.Watch(p)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	// wait and watch
	<-done
}

func ExecCommand() error {
	if *OnChange == "" {
		return nil
	} else {
		args := strings.Split(*OnChange, " ")
		cmd := exec.Command(args[0], args[1:]...)

		if !*Quiet {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		cmd.Stdin = os.Stdin
		return cmd.Run()
	}
}

// Resolve path arguments by walking directories and adding subfolders.
func ResolvePaths(args *string) ([]string, error) {
	var stat os.FileInfo
	resolved := make([]string, 0)

	var recurse error = nil

	if *NoRecurse {
		recurse = filepath.SkipDir
	}

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			resolved = append(resolved, path)
		}

		return recurse
	}
	path := args

	stat, err = os.Stat(*path)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		resolved = append(resolved, *path)
	}

	err = filepath.Walk(*path, walker)

	return resolved, nil
}
