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

func Generate(groupNode *GroupNode, targetDir string, devBranch string) (err error) {

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

	if err = generator.createFileWriter(targetDir); err != nil {
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

func (o *ScriptGenerator) createFileWriter(target string) (err error) {
	for _, writer := range o.writers {
		if err = writer.createFileWriter(target); err != nil {
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
	command  string
	fileName string

	shFile  *os.File
	cmdFile *os.File

	shWriter  *bufio.Writer
	cmdWriter *bufio.Writer
}

func (o *commandFileName) createFileWriter(target string) (err error) {
	if o.shFile, err = os.OpenFile(
		fmt.Sprintf("%v/%v.sh", target, o.fileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777); err != nil {
		return
	}

	if o.cmdFile, err = os.Create(
		fmt.Sprintf("%v/%v.cmd", target, o.fileName)); err != nil {
		return
	}

	o.shWriter = bufio.NewWriter(o.shFile)
	o.cmdWriter = bufio.NewWriter(o.cmdFile)

	_, err = o.shWriter.WriteString("#!/bin/bash\n")

	//change to directory of the script
	_, err = o.shWriter.WriteString("DIR=\"$( cd \"$( dirname \"${BASH_SOURCE[0]}\" )\" >/dev/null 2>&1 && pwd )\"\n")
	_, err = o.shWriter.WriteString("pushd \"$DIR\"\n")

	_, err = o.cmdWriter.WriteString("pushd \"%~dp0\"\r\n")

	return
}

func (o *commandFileName) flush() (err error) {
	if _, err = o.shWriter.WriteString("popd\n"); err != nil {
		return
	}

	if err = o.shWriter.Flush(); err != nil {
		return
	}

	if _, err = o.cmdWriter.WriteString("popd\r\n"); err != nil {
		return
	}

	if err = o.cmdWriter.Flush(); err != nil {
		return
	}
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

func (o *genericCommandWriter) ensureDir(group *gitlab.Group) (err error) {
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

	createFileWriter(target string) (err error)
	flush() (err error)
	close() (err error)
}
