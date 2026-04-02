package namespaces

type CheckStatus int

const (
	StatusFree CheckStatus = iota + 1
	StatusUsed
	StatusPending
	StatusFailed
)

type Checker interface {
	PrepareName(name string) string
	Check(name string, params map[string]interface{}) CheckStatus
	ValidateName(name string) error
	GetId() int
	GetName() string
}
