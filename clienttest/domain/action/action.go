package action

type Action interface {
	Run() error
	Success() bool
	Error() error
}
