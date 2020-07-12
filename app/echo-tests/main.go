package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

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
		// Print to stdout
		fmt.Println("Message from HTTPErrorHandler", c.Path(), c.QueryParams(), err)

		// Call the default handler to return the HTTP response
		server.DefaultHTTPErrorHandler(err, c)
	}

	server.GET("/users", func(c echo.Context) error {
		users, err := dbGetUsers()
		if err != nil {
			return c.JSON(err.Code, err)
		}

		return c.JSON(http.StatusOK, users)
	})

	server.GET("/posts", func(c echo.Context) error {
		users, err := dbPostUsers()
		if err != nil {
			return echo.NewHTTPError(err.Code, err)
		}

		return c.JSON(http.StatusOK, users)
	})

	log.Fatal(server.Start(":8088"))
}

func dbGetUsers() ([]string, *ServiceError) {

	err := &ServiceError{
		Code:    http.StatusBadRequest,
		Message: "Error for get endpoint..",
		Err:     errors.New("unavailable"),
	}

	return nil, err
}

func dbPostUsers() ([]string, *ServiceError) {

	return nil, &ServiceError{http.StatusTeapot,
		"Error for post endpoint.",
		errors.New("Another error."),
	}
}
