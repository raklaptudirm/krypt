package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func main() {
	var args []string

	// set env variables from args
	for _, arg := range os.Args[1:] {
		idx := strings.IndexRune(arg, '=')
		if idx >= 0 { // var=value
			os.Setenv(arg[:idx], arg[idx+1:])
			continue
		}

		args = append(args, arg)
	}

	// remaining args are tasks
	for _, arg := range args {
		f, ok := tasks[arg]
		if !ok {
			fmt.Fprintf(os.Stderr, "Invalid task %v.\n", arg)
		}

		err := f()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

var tasks = map[string]func() error{
	"build": build,
}

func build() error {
	fmt.Fprintln(os.Stderr, "starting krypt build...")
	project := "github.com/raklaptudirm/krypt"

	fmt.Fprintln(os.Stderr, "setting ldflags...")
	ldflags := fmt.Sprintf("-X %s/internal/build.Version=%s", project, version())
	ldflags = fmt.Sprintf("-X %s/internal/build.Date=%s %s", project, date(), ldflags)

	fmt.Fprintln(os.Stderr, "building executable...")
	return run("go", "build", "-trimpath", "-ldflags", ldflags, "-o", exe(), "./cmd/krypt")
}

func exe() string {
	exe := "bin/krypt"
	if runtime.GOOS == "windows" {
		exe += ".exe" // bin/krypt.exe
	}

	return exe
}

func version() string {
	desc, err := runOutput("git", "describe", "--tags")
	if err == nil {
		return desc
	}

	rev, _ := runOutput("git", "rev-parse", "--short", "HEAD")
	return rev
}

func date() string {
	t := time.Now()
	return t.Format("2006-01-02") // YYYY-MM-DD
}

func run(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func runOutput(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stderr = ioutil.Discard // discard the stderr
	out, err := cmd.Output()    // copy the stdout

	return strings.TrimSuffix(string(out), "\n"), err
}
