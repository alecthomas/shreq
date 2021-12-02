package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
	"mvdan.cc/sh/v3/syntax"
)

var validCommands = func() map[string]bool {
	cmds := []string{
		// POSIX utilites (http://pubs.opengroup.org/onlinepubs/9699919799/idx/utilities.html)
		"admin", "alias", "ar", "asa", "at", "awk", "basename", "batch", "bc", "bg", "c99",
		"cal", "cat", "cd", "cflow", "chgrp", "chmod", "chown", "cksum", "cmp", "comm", "command",
		"compress", "cp", "crontab", "csplit", "ctags", "cut", "cxref", "date", "dd", "delta", "df",
		"diff", "dirname", "du", "echo", "ed", "env", "ex", "expand", "expr", "false", "fc",
		"fg", "file", "find", "fold", "fort77", "fuser", "gencat", "get", "getconf", "getopts", "grep",
		"hash", "head", "iconv", "id", "ipcrm", "ipcs", "jobs", "join", "kill", "lex", "link",
		"ln", "locale", "localedef", "logger", "logname", "lp", "ls", "m4", "mailx", "make", "man",
		"mesg", "mkdir", "mkfifo", "more", "mv", "newgrp", "nice", "nl", "nm", "nohup", "od",
		"paste", "patch", "pathchk", "pax", "pr", "printf", "prs", "ps", "pwd", "qalter", "qdel",
		"qhold", "qmove", "qmsg", "qrerun", "qrls", "qselect", "qsig", "qstat", "qsub", "read", "renice",
		"rm", "rmdel", "rmdir", "sact", "sccs", "sed", "sh", "sleep", "sort", "split", "strings",
		"strip", "stty", "tabs", "tail", "talk", "tee", "test", "time", "touch", "tput", "tr",
		"true", "tsort", "tty", "type", "ulimit", "umask", "unalias", "uname", "uncompress", "unexpand", "unget",
		"uniq", "unlink", "uucp", "uudecode", "uuencode", "uustat", "uux", "val", "vi", "wait", "wc",
		"what", "who", "write", "xargs", "yacc",
		// Bash builtins
		":", ".", "[", "alias", "bg", "bind", "break", "builtin", "case", "cd", "command", "compgen",
		"complete", "continue", "declare", "dirs", "disown", "echo", "enable",
		"eval", "exec", "exit", "export", "fc", "fg", "getopts", "hash", "help",
		"history", "if", "jobs", "kill", "let", "local", "logout", "popd", "printf",
		"pushd", "pwd", "read", "readonly", "return", "set", "shift", "shopt",
		"source", "suspend", "test", "times", "trap", "type", "typeset", "ulimit",
		"umask", "unalias", "unset", "until", "wait", "while",
		// Other
		"bash",
	}
	out := make(map[string]bool, len(cmds))
	for _, cmd := range cmds {
		out[cmd] = true
	}
	return out
}()

var cli struct {
	Allow  []string `short:"a" enum:"none,relative,var-relative" help:"Enable optional features (${enum})." default:"none"`
	Cmds   []string `short:"c" placeholder:"CMD" help:"Extra commands to allow."`
	Script []string `arg:"" placeholder:"SCRIPT" type:"existingfile" help:"Shell scripts to validate."`
}

func main() {
	kctx := kong.Parse(&cli, kong.Description(`
This utility verifies all commands used by a shell script against an allow list:

	`+strings.ReplaceAll(builtins(70), "\n", "\n  ")))
	pwd, err := os.Getwd()
	kctx.FatalIfErrorf(err)
	for _, cmd := range cli.Cmds {
		validCommands[cmd] = true
	}
	allow := map[string]bool{}
	for _, feature := range cli.Allow {
		if feature == "none" {
			allow = map[string]bool{}
		} else {
			allow[feature] = true
		}
	}
	parser := syntax.NewParser()
	warnings := []string{}
	for _, path := range cli.Script {
		path, err = filepath.Rel(pwd, path)
		kctx.FatalIfErrorf(err)
		r, err := os.Open(path)
		kctx.FatalIfErrorf(err)
		ast, err := parser.Parse(r, path)
		_ = r.Close()
		kctx.FatalIfErrorf(err)
		localFunctions := map[string]bool{}

		// Collection forward declarations
		syntax.Walk(ast, func(node syntax.Node) bool {
			switch node := node.(type) {
			case *syntax.FuncDecl:
				localFunctions[node.Name.Value] = true
			}
			return true
		})
		// repr.Println(ast)
		syntax.Walk(ast, func(node syntax.Node) bool {
			switch node := node.(type) {
			case *syntax.CallExpr:
				if len(node.Args) == 0 {
					break
				}
				cmd := stringify(node.Args[0])
				if strings.HasPrefix(cmd, "\"") {
					uqcmd, err := strconv.Unquote(cmd) // FIXME: this is a hack
					if err == nil {
						cmd = uqcmd
					}
				}
				if allow["var-relative"] && strings.HasPrefix(cmd, "$") {
					return true
				}
				if !validCommands[cmd] && !localFunctions[cmd] {
					warnings = append(warnings, fmt.Sprintf("%s:%s: unsupported external command: %s", ast.Name, node.Pos(), cmd))
				}
			}
			return true
		})
	}

	for _, warning := range warnings {
		fmt.Fprintln(os.Stderr, warning)
	}
	kctx.Exit(len(warnings))
}

func stringify(node syntax.Node) string {
	out := &strings.Builder{}
	syntax.NewPrinter().Print(out, node)
	return out.String()
}

func builtins(maxWidth int) string {
	w := &strings.Builder{}
	cmds := make([]string, 0, len(validCommands))
	for cmd := range validCommands {
		cmds = append(cmds, cmd)
	}
	sort.Strings(cmds)
	width := 0
	for _, cmd := range cmds {
		if width > 0 && width+1+len(cmd) > maxWidth {
			fmt.Fprintln(w)
			width = 0
		}
		if width > 0 {
			fmt.Fprint(w, " ")
			width++
		}
		fmt.Fprint(w, cmd)
		width += len(cmd)
	}
	return w.String()
}
