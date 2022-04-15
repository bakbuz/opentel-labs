package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"

	"maydere.com/opentel-labs/restapi/pb"
)

// Handler ...
type Handler struct {
	CommonClient pb.CommonServiceClient
}

func newHandler() *Handler {
	h := &Handler{}

	return h
}

func (h *Handler) RegisterHandlers(v1 *echo.Group) {
	common := v1.Group("/common")
	common.GET("/countries", h.HandleCountries)
	common.GET("/countries/:id", h.HandleCountry)
	common.GET("/languages", h.HandleLanguages)
	common.GET("/languages/:id", h.HandleLanguage)
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// GetCountries godoc
// @Summary Get countries
// @Description Get countries.
// @Id get-countries
// @Other countries
// @Tags Common
// @Accept json
// @Produce json
// @Success 200 {object} pb.CountriesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /common/countries [get]
func (h *Handler) HandleCountries(c echo.Context) error {
	fmt.Println("********************* HandleCountries 1 ***********************************")

	ctx := c.Request().Context()
	ctx, span := otel.Tracer("restapi").Start(ctx, "HandleCountries")
	defer span.End()

	fmt.Println("********************* HandleCountries 2 ***********************************")

	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("HandleCountries")

	fmt.Println("********************* HandleCountries 3 ***********************************")

	if h.CommonClient == nil {
		fmt.Println("********************* HandleCountries CommonClient IS NIL ***********************************")
	}

	data, err := h.CommonClient.GetCountries(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Error().Err(err).Msg("")
		return errors.WithStack(err)
	}

	fmt.Println("********************* HandleCountries 4 ***********************************")

	return errors.WithStack(c.JSON(http.StatusOK, data))
}

// GetCountry godoc
// @Summary Get country
// @Description Get country.
// @Id get-country
// @Other country
// @Tags Common
// @Accept json
// @Produce json
// @Success 200 {object} pb.CountryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /common/countries/{id} [get]
func (h *Handler) HandleCountry(c echo.Context) error {
	ctx := c.Request().Context()
	ctx, span := otel.Tracer("restapi").Start(ctx, "HandleCountry")
	defer span.End()

	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("HandleCountry")

	data, err := h.CommonClient.GetCountry(ctx, &pb.Identifier{Id: 1})
	if err != nil {
		logger.Error().Err(err).Msg("")
		return errors.WithStack(err)
	}
	return errors.WithStack(c.JSON(http.StatusOK, data))
}

// GetLanguages godoc
// @Summary Get languages
// @Description Get languages.
// @Id get-languages
// @Other languages
// @Tags Common
// @Accept json
// @Produce json
// @Success 200 {object} pb.LanguagesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /common/languages [get]
func (h *Handler) HandleLanguages(c echo.Context) error {
	ctx := c.Request().Context()
	ctx, span := otel.Tracer("restapi").Start(ctx, "HandleLanguages")
	defer span.End()

	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("HandleLanguages")

	data, err := h.CommonClient.GetLanguages(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Error().Err(err).Msg("")
		return errors.WithStack(err)
	}
	return errors.WithStack(c.JSON(http.StatusOK, data))
}

// GetLanguage godoc
// @Summary Get language
// @Description Get language.
// @Id get-language
// @Other language
// @Tags Common
// @Accept json
// @Produce json
// @Success 200 {object} pb.LanguageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /common/languages/{id} [get]
func (h *Handler) HandleLanguage(c echo.Context) error {
	ctx := c.Request().Context()
	ctx, span := otel.Tracer("restapi").Start(ctx, "HandleLanguage")
	defer span.End()

	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("HandleLanguage")

	data, err := h.CommonClient.GetLanguage(ctx, &pb.Identifier{Id: 1})
	if err != nil {
		logger.Error().Err(err).Msg("")
		return errors.WithStack(err)
	}
	return errors.WithStack(c.JSON(http.StatusOK, data))
}
