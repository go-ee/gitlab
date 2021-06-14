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
	if client, err = gitlab.NewClient(params.Token, gitlab.WithBaseURL(params.Url)); err != nil {
		return
	}

	var group *gitlab.Group
	if group, _, err = client.Groups.GetGroup(params.GroupName); err == nil {

		extractor := &ModelExtractor{
			client:               client,
			alreadyHandledGroups: make(map[int]bool, 0),
			ignoreGroupNames:     params.IgnoreGroupNames,
		}

		ret, err = extractor.extract(group)
	}
	return
}

func (o *ModelExtractor) extract(group *gitlab.Group) (ret *GroupNode, err error) {
	o.alreadyHandledGroups[group.ID] = true
	ret = NewGroupNode(group)

	// refresh projects
	if group.Projects == nil {
		if group, _, err = o.client.Groups.GetGroup(group.ID); err != nil {
			return
		}
	}

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
			if o.shallHandle(subGroup.ID, subGroup.Name) {
				if group, err := o.extract(subGroup); err == nil {
					parent.AddChild(group)
				} else {
					logrus.Warn(err)
				}
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
	for _, sharedGroup := range project.SharedWithGroups {
		if o.shallHandle(sharedGroup.GroupID, sharedGroup.GroupName) {
			if loadedGroup, _, err := o.client.Groups.GetGroup(sharedGroup.GroupID); err == nil {
				if group, err := o.extract(loadedGroup); err != nil {
					parent.AddChild(group)
					logrus.Warn(err)
				}
			} else {
				logrus.Warn(err)
			}
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
	Group    *gitlab.Group `json:group"`
	Children []*GroupNode  `json:children"`
}

func NewGroupNode(group *gitlab.Group) *GroupNode {
	return &GroupNode{group, []*GroupNode{}}
}

func (o *GroupNode) AddChild(group *GroupNode) {
	o.Children = append(o.Children, group)
}
