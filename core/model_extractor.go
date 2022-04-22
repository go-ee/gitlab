package core

import (
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"strconv"
)

type ModelExtractor struct {
	client               GitlabLite
	alreadyHandledGroups map[int]bool
	ignoreGroupNames     map[string]bool
}

type ExtractParams struct {
	Group            string
	IgnoreGroupNames map[string]bool
}

func ExtractFromServer(params *ExtractParams, access *ServerAccess) (ret *GroupNode, err error) {
	var gitlabLiteByServer *GitlabLiteByAPI
	if gitlabLiteByServer, err = NewGitlabLiteByAPI(access); err == nil {
		ret, err = Extract(params, gitlabLiteByServer)
	}
	return
}

func Extract(params *ExtractParams, client GitlabLite) (ret *GroupNode, err error) {
	extractor := &ModelExtractor{
		client:               client,
		alreadyHandledGroups: make(map[int]bool, 0),
		ignoreGroupNames:     params.IgnoreGroupNames,
	}
	ret, err = extractor.ExtractByGroup(params.Group)
	return
}

func (o *ModelExtractor) ExtractByGroup(groupNameOrId string) (ret *GroupNode, err error) {
	var group *gitlab.Group
	if group, err = o.client.GetGroupByName(groupNameOrId); err == nil {
		ret, err = o.extract(group)
	} else {
		logrus.Errorf("can't find group by name: %v => %v", groupNameOrId, err)
		if groupId, numErr := strconv.Atoi(groupNameOrId); numErr == nil {
			if group, err = o.client.GetGroup(groupId); err == nil {
				ret, err = o.extract(group)
			} else {
				logrus.Errorf("can't find group by ID: %v => %v", groupNameOrId, err)
			}
		}
	}
	return
}

func (o *ModelExtractor) handleChildGroup(parent *GroupNode, groupId int, groupName string) (err error) {
	if o.shallHandle(groupId, groupName) {
		o.alreadyHandledGroups[groupId] = true

		var group *gitlab.Group
		if group, err = o.client.GetGroup(groupId); err != nil {
			return
		}
		var groupNode *GroupNode
		if groupNode, err = o.extract(group); err == nil {
			parent.AddChild(groupNode)
		}
	}
	return
}

func (o *ModelExtractor) extract(group *gitlab.Group) (ret *GroupNode, err error) {
	logrus.Debugf("handle group '%v'", group.Name)
	o.alreadyHandledGroups[group.ID] = true

	ret = NewGroupNode(group)

	logrus.Debugf("%v projects in %v", len(group.Projects), group.Name)
	for _, project := range group.Projects {
		o.handleSharedGroups(ret, project)
	}
	o.handleSubGroups(ret, group.ID)

	return
}

func (o *ModelExtractor) handleSubGroups(parent *GroupNode, groupId int) {
	if subGroups, err := o.client.ListSubgroups(groupId); err == nil {
		for _, subGroup := range subGroups {
			if err := o.handleChildGroup(parent, subGroup.ID, subGroup.Name); err != nil {
				logrus.Warn(err)
			}
		}
	} else {
		logrus.Warn(err)
	}
	return
}

func (o *ModelExtractor) shallHandle(groupId int, groupName string) bool {
	return !o.alreadyHandledGroups[groupId] && !o.ignoreGroupNames[groupName]
}

func (o *ModelExtractor) handleSharedGroups(parent *GroupNode, project *gitlab.Project) {
	logrus.Debugf("handle group of project '%v'", project.Name)
	for _, sharedGroup := range project.SharedWithGroups {
		if err := o.handleChildGroup(parent, sharedGroup.GroupID, sharedGroup.GroupName); err != nil {
			logrus.Warn(err)
		}
	}
	return
}

type GroupNode struct {
	Group    *gitlab.Group `json:"group"`
	Children []*GroupNode  `json:"children"`
}

func NewGroupNode(group *gitlab.Group) *GroupNode {
	return &GroupNode{group, []*GroupNode{}}
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