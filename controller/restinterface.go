package controller

type RESTController interface {
	Get() error
	Post() error
	Delete() error
	Put() error
	Base() *BaseController
}
