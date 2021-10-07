package core

import (
	"github.com/xanzy/go-gitlab"
)

type GitlabLite interface {
	GetGroupByName(groupName string) (*gitlab.Group, error)
	GetGroup(groupId int) (*gitlab.Group, error)
	ListSubgroups(groupId int) (ret []*gitlab.Group, err error)
}
