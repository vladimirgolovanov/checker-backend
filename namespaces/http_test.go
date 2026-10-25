package namespaces

import "context"

type stubFetcher struct {
	resp   *HttpResponse
	status CheckStatus
}

func (s *stubFetcher) Get(_ context.Context, _ string, _ map[string]string) (*HttpResponse, CheckStatus) {
	return s.resp, s.status
}
