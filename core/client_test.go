package core

import (
	"github.com/xanzy/go-gitlab"
	"go.uber.org/zap"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	token := os.Getenv("TOKEN")
	url := os.Getenv("URL")

	if client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url)); err != nil {
		t.Fatalf(`can't connect %v`, err)
	} else {
		if users, _, err := client.Users.ListUsers(&gitlab.ListUsersOptions{}); err != nil {
			t.Fatalf("can't get users %v", err)
		} else {
			logrus.Infof("users %v", users)
		}
	}
}
