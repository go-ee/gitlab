package core

import "github.com/xanzy/go-gitlab"

type GitlabLiteByServer struct {
	client *gitlab.Client
	access ServerAccess
}

func NewGitlabLiteServer(access *ServerAccess) (ret *GitlabLiteByServer, err error) {
	var client *gitlab.Client
	if client, err = gitlab.NewClient(access.Token, gitlab.WithBaseURL(access.Url)); err == nil {
		ret = &GitlabLiteByServer{client: client}
	}
	return
}

func (o *GitlabLiteByServer) GetGroupByName(groupName string) (ret *gitlab.Group, err error) {
	ret, _, err = o.client.Groups.GetGroup(groupName, nil)
	return
}

func (o *GitlabLiteByServer) GetGroup(groupId int) (ret *gitlab.Group, err error) {
	ret, _, err = o.client.Groups.GetGroup(groupId, nil)
	return
}

func (o *GitlabLiteByServer) ListSubgroups(groupId int) (ret []*gitlab.Group, err error) {
	options := &gitlab.ListSubgroupsOptions{AllAvailable: new(bool)}
	ret, _, err = o.client.Groups.ListSubgroups(groupId, options)
	return
}

type ServerAccess struct {
	Url   string
	Token string
}
