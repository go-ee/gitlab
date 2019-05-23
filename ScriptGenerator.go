package gitlab

import (
	"bufio"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"os"
)

type Params struct {
	Url       string
	GroupName string
	Target    string
	Token     string
	Ignores   map[string]bool
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

type ScriptGenerator struct {
	client                 *gitlab.Client
	alreadyHandledGroupIds map[int]bool
	writers                []commandWriter
	ignores                map[string]bool
}

func Generate(params *Params) (err error) {
	client := gitlab.NewClient(nil, params.Token)
	if err = client.SetBaseURL(params.Url); err != nil {
		return
	}

	var group *gitlab.Group
	if group, _, err = client.Groups.GetGroup(params.GroupName); err == nil {

		commands := []commandWriter{
			&repoCommandWriter{&commandFileName{command: "clone --recurse-submodules -j8", fileName: "clone"}},
			&genericCommandWriter{&commandFileName{command: "pull", fileName: "pull"}},
			&genericCommandWriter{&commandFileName{command: "status", fileName: "status"}},
			&genericCommandWriter{&commandFileName{command: "checkout development", fileName: "development"}},
			&genericCommandWriter{&commandFileName{command: "checkout master", fileName: "master"}},
		}

		generator := &ScriptGenerator{
			client:                 client,
			alreadyHandledGroupIds: make(map[int]bool, 0),
			writers:                commands,
			ignores:                params.Ignores,
		}

		if err = generator.createFileWriter(params.Target); err != nil {
			return
		}

		defer func() {
			err = generator.close()
		}()

		if err = generator.generate(group); err != nil {
			return
		}

		if err = generator.pause(); err != nil {
			return
		}

		if err = generator.flush(); err != nil {
			return
		}
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

func (o *ScriptGenerator) generateDir(group *gitlab.Group) (err error) {
	if err = o.ensureDir(group); err != nil {
		return
	}
	if err = o.cd(group); err != nil {
		return
	}

	if err = o.generate(group); err != nil {
		return
	}

	if err = o.cdBack(); err != nil {
		return
	}

	return
}

func (o *ScriptGenerator) generate(group *gitlab.Group) (err error) {
	o.alreadyHandledGroupIds[group.ID] = true

	if group.Projects == nil {
		if group, _, err = o.client.Groups.GetGroup(group.ID); err != nil {
			return
		}
	}

	for _, project := range group.Projects {
		if err = o.command(project); err != nil {
			return
		}
		err = o.handleSharedGroups(project)
	}
	err = o.handleSubGroups(group.ID)

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

func (o *ScriptGenerator) handleSubGroups(groupId int) (err error) {
	var subGroups []*gitlab.Group
	options := &gitlab.ListSubgroupsOptions{AllAvailable: new(bool)}
	subGroups, _, err = o.client.Groups.ListSubgroups(groupId, options)
	for _, subGroup := range subGroups {
		if !o.alreadyHandledGroupIds[subGroup.ID] && !o.ignores[subGroup.Name] {
			if err = o.generateDir(subGroup); err != nil {
				logrus.Warn(err)
			}
		}
	}
	return err
}

func (o *ScriptGenerator) handleSharedGroups(project *gitlab.Project) (err error) {
	var loadedGroup *gitlab.Group
	for _, sharedGroup := range project.SharedWithGroups {
		if !o.alreadyHandledGroupIds[sharedGroup.GroupID] && !o.ignores[sharedGroup.GroupName] {
			if loadedGroup, _, err = o.client.Groups.GetGroup(sharedGroup.GroupID); err == nil {
				if err = o.generateDir(loadedGroup); err != nil {
					logrus.Warn(err)
				}
			} else {
				logrus.Warn(err)
			}
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
	if o.shFile, err = os.Create(
		fmt.Sprintf("%v/%v.sh", target, o.fileName)); err != nil {
		return
	}

	if o.cmdFile, err = os.Create(
		fmt.Sprintf("%v/%v.cmd", target, o.fileName)); err != nil {
		return
	}

	o.shWriter = bufio.NewWriter(o.shFile)
	o.cmdWriter = bufio.NewWriter(o.cmdFile)

	return
}

func (o *commandFileName) flush() (err error) {
	if err = o.shWriter.Flush(); err != nil {
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
	_, err = o.shWriter.WriteString("cd ..\n")
	_, err = o.cmdWriter.WriteString("cd ..\r\n")
	return
}

func (o *commandFileName) cd(group *gitlab.Group) (err error) {
	_, err = o.shWriter.WriteString(fmt.Sprintf("cd \"%v\"\n", group.Path))
	_, err = o.cmdWriter.WriteString(fmt.Sprintf("cd \"%v\"\r\n", group.Path))
	return
}

func (o *commandFileName) pause() (err error) {
	_, err = o.shWriter.WriteString("\nread -n1 -r -p \"Press any key to continue...\" key\n")
	_, err = o.cmdWriter.WriteString("\r\npause\r\n")
	return
}

type repoCommandWriter struct {
	*commandFileName
}

func (o *repoCommandWriter) command(project *gitlab.Project) (err error) {
	_, err = o.shWriter.WriteString(fmt.Sprintf("git %v %v\n", o.commandFileName.command, project.SSHURLToRepo))
	_, err = o.cmdWriter.WriteString(fmt.Sprintf("git %v %v\r\n", o.commandFileName.command, project.SSHURLToRepo))

	return
}

func (o *repoCommandWriter) ensureDir(group *gitlab.Group) (err error) {
	_, err = o.shWriter.WriteString(fmt.Sprintf("\nmkdir \"%v\"\n", group.Path))
	_, err = o.cmdWriter.WriteString(fmt.Sprintf("\r\nmkdir \"%v\"\r\n", group.Path))
	return
}

type genericCommandWriter struct {
	*commandFileName
}

func (o *genericCommandWriter) command(project *gitlab.Project) (err error) {
	_, err = o.shWriter.WriteString(fmt.Sprintf("git -C %v %v\n", project.Path, o.commandFileName.command))
	_, err = o.cmdWriter.WriteString(fmt.Sprintf("git -C %v %v\r\n", project.Path, o.commandFileName.command))
	return
}

func (o *genericCommandWriter) ensureDir(group *gitlab.Group) (err error) {
	_, err = o.shWriter.WriteString(fmt.Sprintf("\n"))
	_, err = o.cmdWriter.WriteString(fmt.Sprintf("\r\n"))
	return
}
