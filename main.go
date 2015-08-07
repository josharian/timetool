// Timetool works with toolexec to print information about tool runs -- time, files processed and their sizes, etc.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	cmd  []string
	tool string // name of tool: "go", "6g", etc
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("timetool: ")

	flag.Parse()
	cmd = flag.Args()

	tool = cmd[0]
	if i := strings.LastIndex(tool, "/"); i >= 0 {
		tool = tool[i+1:]
	}

	suffix := ""
	switch tool {
	case "asm":
		suffix = ".s"
	case "5g", "6g", "7g", "8g", "9g":
		suffix = ".go"
	case "5l", "6l", "7l", "8l", "9l":
		suffix = ".a"
	}

	var in int
	var size int64
	if suffix != "" {
		for _, arg := range cmd[1:] {
			if strings.HasSuffix(arg, suffix) {
				fi, err := os.Stat(arg)
				if err != nil {
					log.Fatalln(arg, err)
					continue
				}
				in++
				size += fi.Size()
			}
		}
	}

	xcmd := exec.Command(cmd[0], cmd[1:]...)
	xcmd.Stdin = os.Stdin
	xcmd.Stdout = os.Stdout
	xcmd.Stderr = os.Stderr
	start := time.Now()
	err := xcmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(start)
	if suffix != "" {
		// TODO: use tab writer
		// TODO: reduce time precision
		// fmt.Printf("%s %v n=%d sz=%d t/n=%v t/sz=%v\n",
		// 	tool, elapsed,
		// 	in, size,
		// 	time.Duration(float64(elapsed)/float64(in)),
		// 	time.Duration(float64(elapsed)/float64(size)))
		fmt.Printf("%s %d %d\n",
			tool, int(elapsed.Seconds()*1000),
			// int(1000*(float64(elapsed)/float64(in))/float64(time.Second)),
			int(1000000*(float64(elapsed)/float64(size))/float64(time.Second)),
		)

	}
	os.Exit(0)
}
