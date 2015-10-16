package build

import (
	"core/project"
	"fmt"
	"io"
	"os"
	"os/exec"
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

func Build(p *project.Project) error {
	cfg := p.GetConfig()
	cmdArgs := getArgs(p)
	cmd := exec.Command(filepath.Join(cfg.GoBin+"go"), cmdArgs...)
	//	cmd := exec.Command("env")

	if StdOut == nil {
		StdOut = os.Stdout
	}

	cmd.Stdout = StdOut
	cmd.Stderr = StdOut

	fmt.Fprintln(StdOut, "Build project:", cfg.ProjectName)
	fmt.Fprintln(StdOut, "Build command:", filepath.Join(cfg.GoBin+"go"), cmdArgs)
	fmt.Fprintln(StdOut, "Build:")

	cmd.Env = append(cmd.Env, "PATH="+os.Getenv("PATH"))

	if cfg.GoArch != "" {
		cmd.Env = append(cmd.Env, "GOARCH="+cfg.GoArch)
	}
	if cfg.GoOS != "" {
		cmd.Env = append(cmd.Env, "GOOS="+cfg.GoOS)
	}
	if cfg.GoRoot != "" {
		cmd.Env = append(cmd.Env, "GOROOT="+cfg.GoRoot)
	}
	if cfg.GoPath != "" {
		cmd.Env = append(cmd.Env, "GOPATH="+cfg.GoPath)
	}

	return cmd.Run()
}
