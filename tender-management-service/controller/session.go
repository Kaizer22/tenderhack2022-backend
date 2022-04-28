package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"main/logging"
	"main/model/entity"
	"main/model/response"
	"main/repository"
	"main/utils"
	"net/http"
	"strconv"
	"time"
)

// Controller for quotation session entity
type SessionController struct {
	SessionRepo        repository.QuotationSessionRepository
	ProductJournalRepo repository.ProductJournalRepository
	ctx                context.Context
}

// NewProductController example
func NewSessionController(ctx context.Context, repo repository.QuotationSessionRepository,
	pJRepo repository.ProductJournalRepository) *SessionController {
	return &SessionController{
		ProductJournalRepo: pJRepo,
		SessionRepo:        repo,
		ctx:                ctx,
	}
}

// NewQuotationSession godoc
// @Summary            Add new session
// @Description    Add new quotation session
// @Tags                          sessions
// @Accept                        json
// @Produce                       json
// @Param        Authorization  header    string  true  "Authentication header"
// @Param               session   body            entity.QuotationSessionData  true  "Session info"
// @Success             201             {string}  string          "New session successfully added"
// @Failure        400  {object}  utils.HTTPError
// @Failure        404  {object}  utils.HTTPError
// @Failure        500  {object}  utils.HTTPError
// @Router                        /api/v1/sessions [post]
func (c SessionController) NewQuotationSession(ctx *gin.Context) {
	var s entity.QuotationSessionData
	var pJ []*entity.ProductJournal
	if err := ctx.ShouldBindJSON(&s); err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	//TODO add struct validation
	session := entity.QuotationSession{
		Name:                   s.Name,
		CreatorId:              s.CreatorId,
		Status:                 entity.StatusActive,
		SessionDuration:        s.SessionDuration,
		StartPrice:             s.StartPrice,
		CurrentPrice:           s.StartPrice,
		SessionStepPercent:     s.SessionStepPercent,
		StartTime:              time.Now(),
		IsInAdditionalPurchase: false,
	}
	for _, product := range s.Products {
		pJ = append(pJ, &entity.ProductJournal{
			RecordID:           0,
			ProductId:          product.ProductId,
			Product:            nil,
			QuotationSessionId: session.ID,
			Count:              product.Count,
		})
	}
	session.Products = pJ
	c2 := context.Background()
	id, err := c.SessionRepo.NewQuotationSession(c2, session)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	err = c.ProductJournalRepo.AddProductJournal(c.ctx, id, s.Products)
	if err != nil {
		logging.ErrorFormat("Error inserting product journal %s", err)
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, response.SessionCreated{
		Msg: "Session created",
		Id:  id,
	})
}

// GetSessionById godoc
// @Summary            Get session by ID
// @Description    Returns session by ID
// @Tags                     sessions
// @Accept                   json
// @Produce                  json
// @Param        Authorization  header    string  true  "Authentication header"
// @Param               id             path        int  true  "Session ID"
// @Success             200            {object}  entity.QuotationSession
// @Failure        400       {object}  utils.HTTPError
// @Failure        404       {object}  utils.HTTPError
// @Failure        500       {object}  utils.HTTPError
// @Router                   /api/v1/sessions/{id} [get]
func (c *SessionController) GetSessionById(ctx *gin.Context) {
	id := ctx.Param("id")
	iID, err := strconv.Atoi(id)
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	//TODO deal with contexts correctly
	c2 := context.Background()
	session, err := c.SessionRepo.GetSessionById(c2, int64(iID))
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, session)
}

// GetAllSessions godoc
// @Summary            Get all quotation sessions
// @Description    		Returns all short sessions
// @Tags                     sessions
// @Accept                   json
// @Produce                  json
// @Param        Authorization  header    string  true  "Authentication header"
// @Success             200            {object}  entity.QuotationSessionShort
// @Failure        400       {object}  utils.HTTPError
// @Failure        404       {object}  utils.HTTPError
// @Failure        500       {object}  utils.HTTPError
// @Router                   /api/v1/sessions [get]
func (c *SessionController) GetAllSessions(ctx *gin.Context) {
	//TODO deal with contexts correctly
	c2 := context.Background()
	sessions, err := c.SessionRepo.GetAllSessions(c2)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	var res []*entity.QuotationSessionShort
	for _, session := range sessions {
		res = append(res, &entity.QuotationSessionShort{
			SessionId: session.ID,
			Status:    session.Status,
		})
	}

	ctx.JSON(http.StatusOK, res)
}
