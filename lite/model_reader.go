package lite

import (
	"github.com/go-ee/utils/lg"
	"github.com/xanzy/go-gitlab"
	"path/filepath"
	"strconv"
)

type ModelReader struct {
	Client           GitlabLite
	IgnoreGroupNames map[string]bool

	alreadyHandledGroups map[int]bool
}

func (o *ModelReader) ReadGroupModelByGroup(groupNameOrId string) (ret *GroupNode, err error) {
	o.alreadyHandledGroups = make(map[int]bool, 0)

	var group *gitlab.Group
	if group, err = o.Client.GetGroupByName(groupNameOrId); err == nil {
		ret, err = o.readModelForGroup(group)
	} else {
		lg.LOG.Errorf("can't find group by name: %v => %v", groupNameOrId, err)
		if groupId, numErr := strconv.Atoi(groupNameOrId); numErr == nil {
			if group, err = o.Client.GetGroup(groupId); err == nil {
				ret, err = o.readModelForGroup(group)
			} else {
				lg.LOG.Errorf("can't find group by ID: %v => %v", groupNameOrId, err)
			}
		}
	}

	if err == nil {
		ret.SetAsRelativeRootPath()
	}
	return
}

func (o *ModelReader) readChildGroup(parent *GroupNode, groupId int, groupName string) (err error) {
	if o.shallReadGroup(groupId, groupName) {
		o.alreadyHandledGroups[groupId] = true

		var group *gitlab.Group
		if group, err = o.Client.GetGroup(groupId); err != nil {
			return
		}
		var groupNode *GroupNode
		if groupNode, err = o.readModelForGroup(group); err == nil {
			parent.AddChild(groupNode)
		}
	}
	return
}

func (o *ModelReader) readModelForGroup(group *gitlab.Group) (ret *GroupNode, err error) {
	lg.LOG.Infof("readModelForGroup group '%v(%v)'", group.Name, group.ID)
	o.alreadyHandledGroups[group.ID] = true

	ret = NewGroupNode(group)

	lg.LOG.Debugf("%v projects in %v", len(group.Projects), group.Name)
	for _, project := range group.Projects {
		o.readSharedGroups(ret, project)
	}
	o.readSubGroups(ret, group.ID)

	return
}

func (o *ModelReader) readSubGroups(parent *GroupNode, groupId int) {
	if subGroups, err := o.Client.ListSubgroups(groupId); err == nil {
		for _, subGroup := range subGroups {
			if err = o.readChildGroup(parent, subGroup.ID, subGroup.Name); err != nil {
				lg.LOG.Warn(err)
			}
		}
	} else {
		lg.LOG.Warn(err)
	}
	return
}

func (o *ModelReader) shallReadGroup(groupId int, groupName string) bool {
	return !o.alreadyHandledGroups[groupId] && !o.IgnoreGroupNames[groupName]
}

func (o *ModelReader) readSharedGroups(parent *GroupNode, project *gitlab.Project) {
	lg.LOG.Debugf("handle group of project '%v'", project.Name)
	for _, sharedGroup := range project.SharedWithGroups {
		if err := o.readChildGroup(parent, sharedGroup.GroupID, sharedGroup.GroupName); err != nil {
			lg.LOG.Warn(err)
		}
	}
	return
}

type GroupNode struct {
	Group            *gitlab.Group `json:"group"`
	Children         []*GroupNode  `json:"children"`
	RelativeRootPath string        `json:"-"`
}

func NewGroupNode(group *gitlab.Group) *GroupNode {
	return &GroupNode{Group: group, Children: []*GroupNode{}}
}

func (o *GroupNode) AddChild(group *GroupNode) {
	o.Children = append(o.Children, group)
}

func (o *GroupNode) ChildrenGroups() (ret []*gitlab.Group) {
	ret = make([]*gitlab.Group, len(o.Children))
	for i, groupNode := range o.Children {
		ret[i] = groupNode.Group
	}
	return
}

func (o *GroupNode) SetAsRelativeRootPath() {
	o.RelativeRootPath = o.Group.Path
	o.setChildrenRelativeRootPath()
}

func (o *GroupNode) setChildrenRelativeRootPath() {
	for _, child := range o.Children {
		child.RelativeRootPath = filepath.Join(o.RelativeRootPath, child.Group.Path)
		if len(child.Children) > 0 {
			child.setChildrenRelativeRootPath()
		}
	}
}
