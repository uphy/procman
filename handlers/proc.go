package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/uphy/procman/process"
)

type (
	Proc struct {
		process *process.Process
	}
)

func NewProc(p *process.Process) *Proc {
	return &Proc{p}
}

func (p *Proc) result(c echo.Context, err error) error {
	if err == nil {
		return c.JSON(http.StatusOK, echo.Map{
			"result": "success",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"result":  "failure",
		"message": err.Error(),
	})
}

func (p *Proc) Start(c echo.Context) error {
	return p.result(c, p.process.Start())
}

func (p *Proc) Stop(c echo.Context) error {
	return p.result(c, p.process.Stop())
}

func (p *Proc) Restart(c echo.Context) error {
	return p.result(c, p.process.Restart())
}

func (p *Proc) Status(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"status": p.process.State(),
	})
}
