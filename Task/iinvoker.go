package task

type IInvoker interface {
	Run() bool
	Stop() bool
}
