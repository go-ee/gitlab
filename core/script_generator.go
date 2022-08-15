package core

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"os"
)

type ScriptGenerator struct {
	writers   []commandWriter
	devBranch string
}

func Generate(groupNode *GroupNode, scriptsDir string, reposDir string, devBranch string) (err error) {
	logrus.Infof("generate scripts for group '%v', scripts folder '%v', repos folder '%v', devBranch '%v'",
		groupNode.Group.Name, scriptsDir, reposDir, devBranch)
	commands := []commandWriter{
		&repoCommandWriter{&commandFileName{command: "clone --recurse-submodules -j8", fileName: "clone"}},
		&genericCommandWriter{&commandFileName{command: "pull", fileName: "pull"}},
		&genericCommandWriter{&commandFileName{command: "status", fileName: "status"}},
		&genericCommandWriter{&commandFileName{command: "checkout " + devBranch, fileName: "devBranch"}},
		&genericCommandWriter{&commandFileName{command: "checkout master", fileName: "master"}},
	}

	generator := &ScriptGenerator{
		writers:   commands,
		devBranch: devBranch,
	}

	if err = generator.createFileWriter(scriptsDir, reposDir); err != nil {
		return
	}

	defer func() {
		err = generator.close()
	}()

	if err = generator.generate(groupNode); err != nil {
		return
	}

	if err = generator.pause(); err != nil {
		return
	}

	if err = generator.flush(); err != nil {
		return
	}

	return
}

func (o *ScriptGenerator) createFileWriter(scriptsDir string, reposDir string) (err error) {
	for _, writer := range o.writers {
		if err = writer.createFileWriter(scriptsDir, reposDir); err != nil {
			return
		}
	}
	return
}

func (o *ScriptGenerator) flush() (err error) {
	for _, writer := range o.writers {
		if err = writer.flush(); err != nil {
			return
		}
	}
	return
}

func (o *ScriptGenerator) close() (err error) {
	for _, writer := range o.writers {
		if err = writer.close(); err != nil {
			return
		}
	}
	return
}

func (o *ScriptGenerator) generateDir(groupNode *GroupNode) (err error) {
	if err = o.ensureDir(groupNode.Group); err != nil {
		return
	}
	if err = o.cd(groupNode.Group); err != nil {
		return
	}

	if err = o.generate(groupNode); err != nil {
		return
	}

	if err = o.cdBack(); err != nil {
		return
	}

	return
}

func (o *ScriptGenerator) generate(groupNode *GroupNode) (err error) {
	logrus.Debugf("handle group '%v'", groupNode.Group.Name)
	for _, project := range groupNode.Group.Projects {
		if err = o.command(project); err != nil {
			return
		}
	}
	for _, subGroup := range groupNode.Children {
		if err = o.generateDir(subGroup); err != nil {
			logrus.Warn(err)
		}

	}
	return
}

func (o *ScriptGenerator) ensureDir(group *gitlab.Group) (err error) {
	for _, writer := range o.writers {
		if err = writer.ensureDir(group); err != nil {
			return
		}
	}
	return
}

func (o *ScriptGenerator) cd(group *gitlab.Group) (err error) {
	for _, writer := range o.writers {
		if err = writer.cd(group); err != nil {
			return
		}
	}
	return
}

func (o *ScriptGenerator) cdBack() (err error) {
	for _, writer := range o.writers {
		if err = writer.cdBack(); err != nil {
			return
		}
	}
	return
}

func (o *ScriptGenerator) command(project *gitlab.Project) (err error) {
	logrus.Debugf("handle project '%v'", project.Name)
	for _, writer := range o.writers {
		if err = writer.command(project); err != nil {
			return
		}
	}
	return
}

func (o *ScriptGenerator) pause() (err error) {
	for _, writer := range o.writers {
		if err = writer.pause(); err != nil {
			return
		}
	}
	return
}

type commandFileName struct {
	scriptsDir string
	reposDir   string

	command  string
	fileName string

	shFile  *os.File
	cmdFile *os.File

	shWriter  *bufio.Writer
	cmdWriter *bufio.Writer
}

func (o *commandFileName) createFileWriter(scriptsDir string, reposDir string) (err error) {
	o.scriptsDir = scriptsDir
	o.reposDir = reposDir

	if o.shFile, err = os.OpenFile(
		fmt.Sprintf("%v/%v.sh", scriptsDir, o.fileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777); err != nil {
		return
	}

	if o.cmdFile, err = os.Create(
		fmt.Sprintf("%v/%v.cmd", scriptsDir, o.fileName)); err != nil {
		return
	}

	o.shWriter = bufio.NewWriter(o.shFile)
	o.cmdWriter = bufio.NewWriter(o.cmdFile)

	_, _ = o.shWriter.WriteString("#!/bin/bash\n")
	_, _ = o.shWriter.WriteString("# This file is generated, do not update manually\n\n")

	//change to directory of the script
	_, _ = o.shWriter.WriteString("DIR=\"$( cd \"$( dirname \"${BASH_SOURCE[0]}\" )\" >/dev/null 2>&1 && pwd )\"\n")
	_, _ = o.shWriter.WriteString("pushd \"$DIR\"\n")

	_, _ = o.cmdWriter.WriteString("REM This file is generated, do not update manually\n\n")
	_, _ = o.cmdWriter.WriteString("pushd \"%~dp0\"\r\n")

	if reposDir != "" {
		_, _ = o.shWriter.WriteString("mkdir \"" + reposDir + "\"\n")
		_, _ = o.shWriter.WriteString(fmt.Sprintf("pushd \"%v\"\n", reposDir))

		_, _ = o.cmdWriter.WriteString("REM This file is generated, do not update manually\n\n")
		_, _ = o.cmdWriter.WriteString(fmt.Sprintf("pushd \"%v\"\n", reposDir))
	}

	return
}

func (o *commandFileName) flush() (err error) {
	_, _ = o.shWriter.WriteString("popd\n")
	if o.reposDir != "" {
		_, _ = o.shWriter.WriteString("popd\n")
	}
	_ = o.shWriter.Flush()

	_, _ = o.cmdWriter.WriteString("popd\r\n")
	if o.reposDir != "" {
		_, _ = o.cmdWriter.WriteString("popd\n")
	}
	_ = o.cmdWriter.Flush()

	return
}

func (o *commandFileName) close() (err error) {
	if o.shFile != nil {
		if err = o.shFile.Close(); err != nil {
			return
		}
	}

	if o.cmdFile != nil {
		if err = o.cmdFile.Close(); err != nil {
			return
		}
	}
	return
}

func (o *commandFileName) cdBack() (err error) {
	if _, err = o.shWriter.WriteString("cd ..\n"); err != nil {
		return
	}
	_, err = o.cmdWriter.WriteString("cd ..\r\n")
	return
}

func (o *commandFileName) cd(group *gitlab.Group) (err error) {
	if _, err = o.shWriter.WriteString(fmt.Sprintf("cd \"%v\"\n", group.Path)); err != nil {
		return
	}
	_, err = o.cmdWriter.WriteString(fmt.Sprintf("cd \"%v\"\r\n", group.Path))
	return
}

func (o *commandFileName) pause() (err error) {
	if _, err = o.shWriter.WriteString("\nread -n1 -r -p \"Press any key to continue...\" key\n"); err != nil {
		return
	}
	_, err = o.cmdWriter.WriteString("\r\npause\r\n")
	return
}

type repoCommandWriter struct {
	*commandFileName
}

func (o *repoCommandWriter) command(project *gitlab.Project) (err error) {
	if _, err = o.shWriter.WriteString(fmt.Sprintf("git %v %v\n",
		o.commandFileName.command, project.SSHURLToRepo)); err != nil {
		return
	}
	_, err = o.cmdWriter.WriteString(fmt.Sprintf("git %v %v\r\n",
		o.commandFileName.command, project.SSHURLToRepo))

	return
}

func (o *repoCommandWriter) ensureDir(group *gitlab.Group) (err error) {
	if _, err = o.shWriter.WriteString(fmt.Sprintf("\nmkdir \"%v\"\n", group.Path)); err != nil {
		return
	}
	_, err = o.cmdWriter.WriteString(fmt.Sprintf("\r\nmkdir \"%v\"\r\n", group.Path))
	return
}

type genericCommandWriter struct {
	*commandFileName
}

func (o *genericCommandWriter) command(project *gitlab.Project) (err error) {
	if err = o.echo(project); err != nil {
		return
	}

	if _, err = o.shWriter.WriteString(fmt.Sprintf("git -C %v %v\n",
		project.Path, o.commandFileName.command)); err != nil {
		return
	}
	_, err = o.cmdWriter.WriteString(fmt.Sprintf("git -C %v %v\r\n",
		project.Path, o.commandFileName.command))
	return
}

func (o *genericCommandWriter) ensureDir(_ *gitlab.Group) (err error) {
	_, err = o.shWriter.WriteString(fmt.Sprintf("\n"))
	_, err = o.cmdWriter.WriteString(fmt.Sprintf("\r\n"))
	return
}

func (o *genericCommandWriter) echo(project *gitlab.Project) (err error) {
	if _, err = o.shWriter.WriteString(fmt.Sprintf("echo \"%v %v\"\n",
		o.commandFileName.command, project.PathWithNamespace)); err != nil {
		return
	}
	_, err = o.cmdWriter.WriteString(fmt.Sprintf("echo %v %v\r\n",
		o.commandFileName.command, project.PathWithNamespace))
	return
}

type commandWriter interface {
	command(project *gitlab.Project) (err error)

	ensureDir(group *gitlab.Group) (err error)
	cd(group *gitlab.Group) (err error)
	cdBack() (err error)
	pause() (err error)

	createFileWriter(scriptsDir string, reposDir string) (err error)
	flush() (err error)
	close() (err error)
}
