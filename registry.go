package main

import "github.com/vladimirgolovanov/checker-backend/namespaces"

var CheckerRegistry = map[int]func(params map[string]interface{}) []namespaces.Checker{
	0: single(&namespaces.InstagramChecker{}),
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
	5:  single(&namespaces.TiktokChecker{}),
	6:  single(&namespaces.SnapchatChecker{}),
	7:  single(&namespaces.NpmChecker{}),
	8:  single(&namespaces.GithubChecker{}),
	9:  single(&namespaces.TelegramChecker{}),
	10: single(&namespaces.TelegramBotChecker{}),
	11: single(&namespaces.EtsyChecker{}),
	12: single(&namespaces.PinterestChecker{}),
}

func single(c namespaces.Checker) func(map[string]interface{}) []namespaces.Checker {
	return func(params map[string]interface{}) []namespaces.Checker {
		return []namespaces.Checker{c}
	}
}
