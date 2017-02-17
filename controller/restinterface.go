package controller

type RESTController interface {
	Get() error
	Post() error
	Delete() error
	Put() error
	Options() error
	Base() *BaseController
}
