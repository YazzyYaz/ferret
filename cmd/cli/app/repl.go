package app

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/chzyer/readline"
	"strings"
)

func Repl(version string, opts Options) {
	ferret := compiler.New()

	fmt.Printf("Welcome to Ferret REPL %s\n", version)
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})

	if err != nil {
		panic(err)
	}

	defer rl.Close()

	var commands []string

	timer := NewTimer()

	for {
		line, err := rl.Readline()

		if err != nil {
			break
		}

		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		if strings.HasSuffix(line, "\\") {
			commands = append(commands, line[:len(line)-1])
			continue
		}

		commands = append(commands, line)
		query := strings.Join(commands, "\n")

		commands = make([]string, 0, 10)

		program, err := ferret.Compile(query)

		if err != nil {
			fmt.Println("Failed to parse the query")
			fmt.Println(err)
			continue
		}

		timer.Start()

		out, err := program.Run(
			context.Background(),
			runtime.WithBrowser(opts.Cdp),
		)

		timer.Stop()
		fmt.Println(timer.Print())

		if err != nil {
			fmt.Println("Failed to execute the query")
			fmt.Println(err)
			continue
		}

		fmt.Println(string(out))
	}
}
