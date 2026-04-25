package namespaces

import (
	"context"
	"time"

	grab_instagram "github.com/vladimirgolovanov/grab-proto/gen/instagram"
)

type InstagramChecker struct {
	client grab_instagram.InstagramClient
}

func NewInstagramChecker(client grab_instagram.InstagramClient) *InstagramChecker {
	return &InstagramChecker{client: client}
}

func (i *InstagramChecker) GetId() int {
	return 0
}

func (i *InstagramChecker) GetName() string {
	return "Instagram"
}

func (i *InstagramChecker) PrepareName(name string) string {
	return name
}

func (i *InstagramChecker) ValidateName(name string) error {
	return nil
}

func (i *InstagramChecker) Check(name string, params map[string]interface{}) CheckStatus {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	stream, err := i.client.CheckUsername(ctx, &grab_instagram.CheckUsernameRequest{
		Usernames: []string{name},
	})
	if err != nil {
		return StatusPending
	}

	resp, err := stream.Recv()
	if err != nil {
		return StatusFailed
	}

	if resp.Error != "" {
		return StatusFailed
	}

	if resp.Exists {
		return StatusUsed
	}

	return StatusFree
}
