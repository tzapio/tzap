package resolver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tzapio/tzap/pkg/tzapaction/action"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
)

func LocalRunI(route string, i interface{}) (string, error) {
	jsonBytes, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return LocalRun(route, string(jsonBytes))
}

func LocalRun(route string, jsonBody string) (string, error) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, route, strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	t := initializeSession()
	e.POST("/search", func(c echo.Context) error {
		var searchArgs actionpb.SearchArgs
		if err := json.NewDecoder(c.Request().Body).Decode(&searchArgs); err != nil {
			return err
		}
		result, err := action.Search(t, &actionpb.SearchRequest{SearchArgs: &searchArgs})
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, result)
	})
	// Handle "/prompt" endpoint
	e.POST("/prompt", func(c echo.Context) error {
		var promptArgs actionpb.PromptArgs
		if err := json.NewDecoder(c.Request().Body).Decode(&promptArgs); err != nil {
			return err
		}

		result, err := action.Prompt(t, &actionpb.PromptRequest{PromptArgs: &promptArgs})
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, result)
	})

	refactor := func(c echo.Context) error {
		var basicConfig actionpb.RefactorArgs
		if err := json.NewDecoder(c.Request().Body).Decode(&basicConfig); err != nil {
			return err
		}
		result, err := action.Refactor(t, &actionpb.RefactorRequest{RefactorArgs: &basicConfig})
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, result)
	}
	// Handle "/refactor" endpoint
	e.POST("/refactor", refactor)
	// Handle "/refactor" endpoint
	e.POST("/test", refactor)
	// Handle "/documentation" endpoint
	e.POST("/documentation", refactor)

	e.POST("/code", func(c echo.Context) error {
		var implementCodeArgs actionpb.ImplementArgs
		if err := json.NewDecoder(c.Request().Body).Decode(&implementCodeArgs); err != nil {
			return err
		}
		response, err := action.Implement(t, &actionpb.ImplementRequest{ImplementArgs: &implementCodeArgs})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, response)
	})
	edit := func(c echo.Context) error {
		var editArgs actionpb.EditArgs
		if err := json.NewDecoder(c.Request().Body).Decode(&editArgs); err != nil {
			return err
		}
		response, err := action.Edit(t, &actionpb.EditRequest{EditArgs: &editArgs})
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, response)
	}
	e.POST("/edit", edit)
	e.POST("/add", edit)

	e.ServeHTTP(rec, req)
	localRunResult := rec.Body.String()
	return localRunResult, nil
}
