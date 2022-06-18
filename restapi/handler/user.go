package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"maydere.com/opentel-labs/restapi/model"
	"maydere.com/opentel-labs/restapi/utils"
)

/*****************************************************************/
/* private methods */
/*****************************************************************/

func (h *Handler) getCurrentUser(c echo.Context) (*model.User, error) {
	userId, err := utils.GetCurrentUserIdFromToken(c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	user, err := h.UserStore.GetUserById(c.Request().Context(), userId)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if user == nil || user.Deleted {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// CurrentUser godoc
// @Summary Get current user session
// @Id user-me
// @Tags Users
// @Success 200 {object} SessionResponse
// @Failure 401 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /users/session [get]
func (h *Handler) CurrentUser(c echo.Context) error {
	ctx := c.Request().Context()
	logger := zerolog.Ctx(ctx)

	user, err := h.getCurrentUser(c)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return errors.WithStack(c.JSON(http.StatusOK, newSessionResponse(user)))
}

// SessionResponse ...
type SessionResponse struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Language    string `json:"language"`
	Avatar      string `json:"avatar"`
	Token       string `json:"token"`
}

func newSessionResponse(user *model.User) *SessionResponse {
	res := new(SessionResponse)
	res.Username = user.Username
	res.Email = user.Email
	res.DisplayName = user.DisplayName
	res.Language = user.Language
	res.Avatar = user.Avatar
	res.Token = utils.GenerateJWT(user.Id)
	return res
}
