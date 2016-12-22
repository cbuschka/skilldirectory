package controller

type GenericController struct {
	*BaseController
}

func (c GenericController) Base() *BaseController {
	return c.BaseController
}

func (c GenericController) Get() error    { return nil }
func (c GenericController) Post() error   { return nil }
func (c GenericController) Delete() error { return nil }
func (c GenericController) Put() error    { return nil }
