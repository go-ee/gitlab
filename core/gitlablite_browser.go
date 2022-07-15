package core

import (
	"encoding/json"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"net/url"
	"strings"
)

type GitlabLiteByBrowser struct {
	access *BrowserAccess

	browser playwright.Browser
	page    playwright.Page
}

func NewGitlabLiteByBrowser(access *BrowserAccess, installBrowsers bool) (ret *GitlabLiteByBrowser, err error) {
	ret = &GitlabLiteByBrowser{access: access}
	err = ret.Init(installBrowsers)
	return
}

func (o *GitlabLiteByBrowser) InstallBrowsers() (err error) {
	err = playwright.Install()
	return
}

func (o *GitlabLiteByBrowser) Init(installBrowsers bool) (err error) {
	if installBrowsers {
		if err = playwright.Install(); err != nil {
			return
		}
	}

	var pw *playwright.Playwright
	if pw, err = playwright.Run(); err != nil {
		return
	}

	o.browser, err = pw.Chromium.Launch(
		playwright.BrowserTypeLaunchOptions{Headless: playwright.Bool(false)})
	return
}

func (o *GitlabLiteByBrowser) AuthInteractive(waitForAuth int) (err error) {
	if o.page, err = o.browser.NewPage(); err != nil {
		return
	}

	if _, err = o.page.Goto(o.access.UrlAuth); err != nil {
		return
	}

	o.page.WaitForTimeout(float64(waitForAuth))

	var resp playwright.Response
	groupsUrl := o.access.GroupsUrl()
	resp, err = o.page.Goto(groupsUrl)
	logrus.Debugf("response of %v: %v", groupsUrl, resp)
	return
}

func (o *GitlabLiteByBrowser) GetGroupByName(groupName string) (*gitlab.Group, error) {
	return o.getGroupByNameOrId(pathEscape(groupName))
}

func (o *GitlabLiteByBrowser) GetGroup(groupId int) (ret *gitlab.Group, err error) {
	return o.getGroupByNameOrId(groupId)
}

func (o *GitlabLiteByBrowser) getGroupByNameOrId(
	groupNameOrId interface{}) (ret *gitlab.Group, err error) {

	var resp playwright.Response
	if resp, err = o.page.Goto(o.access.GroupUrl(groupNameOrId)); err != nil {
		return
	}

	var jsonResponse []byte
	if jsonResponse, err = resp.Body(); err == nil {
		err = json.Unmarshal(jsonResponse, &ret)
	}
	return
}

func (o *GitlabLiteByBrowser) ListSubgroups(groupNameOrId int) (ret []*gitlab.Group, err error) {
	var resp playwright.Response
	if resp, err = o.page.Goto(o.access.SubGroupsUrl(groupNameOrId)); err != nil {
		return
	}

	var jsonResponse []byte
	if jsonResponse, err = resp.Body(); err == nil {
		err = json.Unmarshal(jsonResponse, &ret)
	}
	return
}

func pathEscape(s string) string {
	return strings.Replace(url.PathEscape(s), ".", "%2E", -1)
}

type BrowserAccess struct {
	UrlAuth string
	UrlApi  string
}

func (o *BrowserAccess) GroupsUrl() string {
	return fmt.Sprintf("%v/groups", o.UrlApi)
}

func (o *BrowserAccess) GroupUrl(groupNameOrId interface{}) string {
	return fmt.Sprintf("%v/groups/%v", o.UrlApi, groupNameOrId)
}

func (o *BrowserAccess) SubGroupsUrl(groupNameOrId interface{}) string {
	return fmt.Sprintf("%v/groups/%v/subgroups", o.UrlApi, groupNameOrId)
}
