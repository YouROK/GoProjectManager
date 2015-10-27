package build

import (
	"core/project"
	"core/run"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	StdOut io.Writer
)

func getArgs(p *project.Project) []string {
	cfg := p.GetConfig()
	cmdArgs := []string{"build"}

	if cfg.BuildOpt.BuildMode != "" {
		cmdArgs = append(cmdArgs, "-buildmode")
		cmdArgs = append(cmdArgs, cfg.BuildOpt.BuildMode)
	}

	if cfg.BuildOpt.LDFlags != "" {
		cmdArgs = append(cmdArgs, "-ldflags")
		cmdArgs = append(cmdArgs, cfg.BuildOpt.LDFlags)
	}

	cmdArgs = append(cmdArgs, "-v")

	if cfg.BuildOpt.ShowCommand == "1" {
		cmdArgs = append(cmdArgs, "-x")
	}

	if cfg.BuildOpt.Race == "1" {
		cmdArgs = append(cmdArgs, "-race")
	}

	if cfg.BuildOpt.Jobs != "" {
		cmdArgs = append(cmdArgs, "-p")
		cmdArgs = append(cmdArgs, cfg.BuildOpt.Jobs)
	}

	out := ""
	if cfg.BuildOpt.BuildOutDir != "" {
		out, _ = filepath.Abs(filepath.Join(p.GetProjectPath(), cfg.BuildOpt.BuildOutDir))
		out += "/"
		os.MkdirAll(out, 0755)
	}
	if cfg.BuildOpt.BuildOutFile != "" {
		out += cfg.BuildOpt.BuildOutFile
	} else {
		out = cfg.MainPkgPath
	}

	if out != "" {
		cmdArgs = append(cmdArgs, "-o")
		cmdArgs = append(cmdArgs, out)
	}

	if cfg.BuildOpt.BuildCMD != "" {
		cmdArgs = append(cmdArgs, strings.Fields(cfg.BuildOpt.BuildCMD)...)
	}

	cmdArgs = append(cmdArgs, cfg.MainPkgPath)
	return cmdArgs
}

func GetDepends(p *project.Project, downloadOnly bool) error {
	cmdArgs := []string{"get"}
	if downloadOnly {
		cmdArgs = append(cmdArgs, "-d")
	}
	cmdArgs = append(cmdArgs, "-v")
	cmdArgs = append(cmdArgs, p.GetConfig().MainPkgPath)
	cmd := run.NewExec(filepath.Join(p.GetConfig().GoBin+"go"), cmdArgs, p.GetEnvs())
	if StdOut == nil {
		StdOut = os.Stdout
	}

	cmd.StdOut = StdOut
	cmd.StdErr = StdOut

	fmt.Fprintln(StdOut, "Get project dependens:", p.GetConfig().ProjectName)
	fmt.Fprintln(StdOut, "Get command:", filepath.Join(p.GetConfig().GoBin+"go"), cmdArgs)
	fmt.Fprintln(StdOut, "Get:")

	return cmd.Run()
}

func Build(p *project.Project) error {

	cmdArgs := getArgs(p)
	cmd := run.NewExec(filepath.Join(p.GetConfig().GoBin+"go"), cmdArgs, p.GetEnvs())
	if StdOut == nil {
		StdOut = os.Stdout
	}

	cmd.StdOut = StdOut
	cmd.StdErr = StdOut

	fmt.Fprintln(StdOut, "Build project:", p.GetConfig().ProjectName)
	fmt.Fprintln(StdOut, "Build path:", p.GetProjectPath())
	fmt.Fprintln(StdOut, "Build command:", filepath.Join(p.GetConfig().GoBin+"go"), cmdArgs)
	fmt.Fprintln(StdOut, "Build:")

	return cmd.Run()
}
