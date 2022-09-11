package lite

import "github.com/xanzy/go-gitlab"

type GitlabLiteByAPI struct {
	client *gitlab.Client
}

func NewGitlabLiteByAPI(access *ServerAccess) (ret *GitlabLiteByAPI, err error) {
	var client *gitlab.Client
	if client, err = gitlab.NewClient(access.Token, gitlab.WithBaseURL(access.Url)); err == nil {
		ret = &GitlabLiteByAPI{client: client}
	}
	return
}

func (o *GitlabLiteByAPI) GetGroupByName(groupName string) (ret *gitlab.Group, err error) {
	ret, _, err = o.client.Groups.GetGroup(groupName, nil)
	return
}

func (o *GitlabLiteByAPI) GetGroup(groupId int) (ret *gitlab.Group, err error) {
	ret, _, err = o.client.Groups.GetGroup(groupId, nil)
	return
}

func (o *GitlabLiteByAPI) ListSubgroups(groupId int) (ret []*gitlab.Group, err error) {
	options := &gitlab.ListSubgroupsOptions{AllAvailable: new(bool)}
	ret, _, err = o.client.Groups.ListSubgroups(groupId, options)
	return
}

type ServerAccess struct {
	Url   string
	Token string
}
