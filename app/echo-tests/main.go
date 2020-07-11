package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func (r *RequestError) Error() string {
	return fmt.Sprintf("status %d: err %v :msg %v", r.StatusCode, r.Err, r.Message)
}

type RequestError struct {
	StatusCode int
	Message    string
	Err        error
}

func main() {
	server := echo.New()
	server.Use(
		middleware.Recover(),   // Recover from all panics to always have your server up
		middleware.Logger(),    // Log everything to stdout
		middleware.RequestID(), // Generate a request id on the HTTP response headers for identification
	)
	server.Debug = false
	server.HideBanner = true
	server.HTTPErrorHandler = func(err error, c echo.Context) {
		// Take required information from error and context and send it to a service like New Relic
		fmt.Println(c.Path(), c.QueryParams(), err.Error())

		// Call the default handler to return the HTTP response
		server.DefaultHTTPErrorHandler(err, c)
	}

	server.GET("/users", func(c echo.Context) error {
		users, err := dbGetUsers()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, users)
	})

	log.Fatal(server.Start(":8088"))
}

func dbGetUsers() ([]string, error) {
	return nil, &RequestError{
		StatusCode: 404,
		Message:    "Message custom.",
		Err:        errors.New("unavailable"),
	}
}
