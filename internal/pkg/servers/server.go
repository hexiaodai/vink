package servers

type Server interface {
	Run() error
	Stop() error
}
