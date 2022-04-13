package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"tupulung/config"
	"tupulung/entities/web"
	commentService "tupulung/services/comment"

	"github.com/labstack/echo/v4"
)

type CommentHandler struct {
	commentService *commentService.CommentService
}

func NewCommentHandler(commentService *commentService.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}


/*
 * Find All comment 
 * -------------------------------
 * Mengambil data comment berdasarkan filters dan sorts
 */
func (handler CommentHandler) Index(c echo.Context) error {
	
	// Url path parameter & link hateoas
	eventID, err := strconv.Atoi(c.Param("eventID"))
	links := map[string]string {}
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse {
			Status: "OK",
			Code: http.StatusBadRequest,
			Error: "Invalid eventID parameter format",
			Links: links,
		})
	}

	// pagination param
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 50
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	links["self"] = config.Get().App.BaseURL + "/api/events/" + c.Param("eventID") + "/comments?page=" + strconv.Itoa(page)


	// Service call
	filters := []map[string]string{}
	filters = append(filters, map[string]string{
		"field": "event_id",
		"operator": "=",
		"value": strconv.Itoa(eventID),
	})
	commentsRes, err := handler.commentService.FindAll(limit, page, filters, []map[string]interface{}{})
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, web.ErrorResponse{
				Status: "ERROR",
				Code: webErr.Code,
				Error: webErr.Error(),
				Links: links,
			})
		}
		return c.JSON(http.StatusInternalServerError, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusInternalServerError,
			Error: err.Error(),
			Links: links,
		})
	}

	// make pagination data & formatting pagination links
	paginationRes, err := handler.commentService.GetPagination(page, limit, filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusInternalServerError,
			Error: err.Error(),
			Links: links,
		})
	}
	pageUrl := fmt.Sprintf("%s/api/events/%s/comments?page=", config.Get().App.BaseURL, c.Param("eventID"))
	links["first"] = pageUrl + "1"
	links["last"] = pageUrl + strconv.Itoa(paginationRes.TotalPages)
	if paginationRes.Page > 1 {
		links["previous"] = pageUrl + strconv.Itoa(page - 1)
	}
	if paginationRes.Page < paginationRes.TotalPages {
		links["previous"] = pageUrl + strconv.Itoa(page + 1)
	}


	return c.JSON(http.StatusOK, web.SuccessListResponse{
		Status: "OK",
		Code: http.StatusOK,
		Error: nil,
		Links: links,
		Data: commentsRes,
		Pagination: paginationRes,
	}) 
}