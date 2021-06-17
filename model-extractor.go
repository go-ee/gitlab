package gitlab

import (
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

type ModelExtractor struct {
	client               *gitlab.Client
	alreadyHandledGroups map[int]bool
	ignoreGroupNames     map[string]bool
}

func Extract(params *ExtractParams) (ret *GroupNode, err error) {
	var client *gitlab.Client
	if client, err = gitlab.NewClient(params.Token, gitlab.WithBaseURL(params.Url)); err == nil {
		extractor := &ModelExtractor{
			client:               client,
			alreadyHandledGroups: make(map[int]bool, 0),
			ignoreGroupNames:     params.IgnoreGroupNames,
		}
		ret, err = extractor.ExtractByGroupName(params.GroupName)
	}
	return
}

func (o *ModelExtractor) ExtractByGroupName(groupName string) (ret *GroupNode, err error) {
	var group *gitlab.Group
	if group, _, err = o.client.Groups.GetGroup(groupName); err == nil {
		ret, err = o.extract(group)
	}
	return
}

func (o *ModelExtractor) handleChildGroup(parent *GroupNode, groupId int, groupName string) (err error) {
	if o.shallHandle(groupId, groupName) {
		o.alreadyHandledGroups[groupId] = true

		var group *gitlab.Group
		if group, _, err = o.client.Groups.GetGroup(groupId); err != nil {
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
	options := &gitlab.ListSubgroupsOptions{AllAvailable: new(bool)}
	if subGroups, _, err := o.client.Groups.ListSubgroups(groupId, options); err == nil {
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

type ExtractParams struct {
	Url              string
	Token            string
	GroupName        string
	IgnoreGroupNames map[string]bool
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
