package main

import (
	"github.com/vladimirgolovanov/checker-backend/namespaces"
	grab_instagram "github.com/vladimirgolovanov/grab-proto/gen/instagram"
)

type Deps struct {
	InstagramClient grab_instagram.InstagramClient
	Fetcher         namespaces.Fetcher
}

func NewCheckerRegistry(deps Deps) map[int]func(map[string]interface{}) []namespaces.Checker {
	return map[int]func(map[string]interface{}) []namespaces.Checker{
		0: single(namespaces.NewInstagramChecker(deps.InstagramClient)),
		1: func(params map[string]interface{}) []namespaces.Checker {
			zones := []string{"com"}
			if params != nil {
				if z, ok := params["zones"]; ok {
					if zoneSlice, ok := z.([]interface{}); ok {
						if len(zoneSlice) > 50 {
							zoneSlice = zoneSlice[:50]
						}
						zones = make([]string, 0, len(zoneSlice))
						for _, zone := range zoneSlice {
							if s, ok := zone.(string); ok {
								zones = append(zones, s)
							}
						}
					}
				}
			}
			checkers := make([]namespaces.Checker, len(zones))
			for i, zone := range zones {
				checkers[i] = &namespaces.DomainChecker{Zone: zone}
			}
			return checkers
		},
		5:  single(namespaces.NewTiktokChecker(deps.Fetcher)),
		6:  single(namespaces.NewSnapchatChecker(deps.Fetcher)),
		7:  single(namespaces.NewNpmChecker(deps.Fetcher)),
		8:  single(namespaces.NewGithubChecker(deps.Fetcher)),
		9:  single(namespaces.NewTelegramChecker(deps.Fetcher)),
		10: single(namespaces.NewTelegramBotChecker(deps.Fetcher)),
		11: single(namespaces.NewEtsyChecker(deps.Fetcher)),
		12: single(namespaces.NewPinterestChecker(deps.Fetcher)),
		13: single(namespaces.NewYoutubeChecker(deps.Fetcher)),
	}
}

func single(c namespaces.Checker) func(map[string]interface{}) []namespaces.Checker {
	return func(params map[string]interface{}) []namespaces.Checker {
		return []namespaces.Checker{c}
	}
}
