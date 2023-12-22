package namespaces

type CheckStatus int

const (
	StatusFree CheckStatus = iota + 1
	StatusUsed
	StatusPending
	StatusFailed
)

type Checker interface {
	Check(name string) CheckStatus
	GetId() int
}
