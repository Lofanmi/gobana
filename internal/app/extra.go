package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Just-maple/structgraph"
	"github.com/sirupsen/logrus"
)

func GenStructGraph(in interface{}, filename string) {
	dot := structgraph.Draw(in)
	if len(dot) == 0 {
		return
	}
	_ = os.MkdirAll(filepath.Dir(filename), 0775)
	cmd := fmt.Sprintf("dot -o%s -T%s -K%s", filename, "svg", "dot")
	if err := system(cmd, dot); err != nil {
		logrus.Errorf("gen structure graph error: %v", err)
	} else {
		logrus.Infof("gen structure graph [ %s ] success", filename)
	}
}

func system(c string, dot string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command(`cmd`, `/C`, c)
	} else {
		cmd = exec.Command(`/bin/sh`, `-c`, c)
	}
	cmd.Stdin = strings.NewReader(dot)
	return cmd.Run()
}
