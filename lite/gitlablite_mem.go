package lite

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
)

type GitlabLiteMem struct {
	groups         map[int]*gitlab.Group
	groupsChildren map[int][]*gitlab.Group
	groupsByPath   map[string]*gitlab.Group
}

func NewGitlabLiteMem() (ret *GitlabLiteMem) {
	ret = &GitlabLiteMem{
		groups:         map[int]*gitlab.Group{},
		groupsChildren: map[int][]*gitlab.Group{},
		groupsByPath:   map[string]*gitlab.Group{},
	}
	return
}

func (o *GitlabLiteMem) GetGroupByName(groupName string) (ret *gitlab.Group, err error) {
	if ret, _ = o.groupsByPath[groupName]; ret == nil {
		err = fmt.Errorf("no group found, for name '%v'", groupName)
	}
	return
}

func (o *GitlabLiteMem) GetGroup(groupId int) (ret *gitlab.Group, err error) {
	if ret, _ = o.groups[groupId]; ret == nil {
		err = fmt.Errorf("no group found, for id '%v'", groupId)
	}
	return
}

func (o *GitlabLiteMem) ListSubgroups(groupId int) (ret []*gitlab.Group, err error) {
	if ret, _ = o.groupsChildren[groupId]; ret == nil {
		ret = []*gitlab.Group{}
	}
	return
}

func (o *GitlabLiteMem) AddGroup(group *gitlab.Group) {
	o.groups[group.ID] = group
	o.groupsByPath[group.Path] = group
	children := o.groupsChildren[group.ParentID]
	if children == nil {
		children = []*gitlab.Group{}
	}
	o.groupsChildren[group.ParentID] = append(children, group)
}
