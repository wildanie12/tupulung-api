package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"tupulung/config"
	"tupulung/entities"
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

	// response
	return c.JSON(http.StatusOK, web.SuccessListResponse{
		Status: "OK",
		Code: http.StatusOK,
		Error: nil,
		Links: links,
		Data: commentsRes,
		Pagination: paginationRes,
	}) 
}

/*
 * Create comment 
 * -------------------------------
 * Membuat commnet pada event untuk authenticated user
 */
func (handler CommentHandler) Create(c echo.Context) error {
	
	// Url path parameter & link hateoas
	eventID, err := strconv.Atoi(c.Param("eventID"))
	links := map[string]string {}
	links["self"] = config.Get().App.BaseURL + "/api/events/" + c.Param("eventID") + "/comments"
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse {
			Status: "OK",
			Code: http.StatusBadRequest,
			Error: "Invalid eventID parameter format",
			Links: links,
		})
	}

	// Populate form
	commentReq := entities.CommentRequest{}
	c.Bind(&commentReq)

	// token
	token := c.Get("user")

	// Insert comment
	commentRes, err := handler.commentService.Create(commentReq, eventID, token)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, web.ErrorResponse{
				Status: "ERROR",
				Code: webErr.Code,
				Error: webErr.Error(),
				Links: links,
			})
		} else if reflect.TypeOf(err).String() == "web.ValidationError" {
			valErr := err.(web.ValidationError)
			return c.JSON(valErr.Code, web.ValidationErrorResponse{
				Status: "ERROR",
				Code: valErr.Code,
				Error: valErr.Error(),
				Errors: valErr.Errors,
				Links: links,
			})
		}
	}

	// response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code: 201,
		Error: nil,
		Links: links,
		Data: map[string]int{
			"id": int(commentRes.ID),
		},
	})
}

/*
 * Update Comment
 * -------------------------------
 * Edit komentar user, hanya pemilik komentar yang dapat mengedit
 */
func (handler CommentHandler) Update(c echo.Context) error {
	
	// Url path parameter & link hateoas
	commentID, err := strconv.Atoi(c.Param("commentID"))
	links := map[string]string {}
	links["self"] = config.Get().App.BaseURL + "/api/events/comments/" + c.Param("commentID")
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse {
			Status: "OK",
			Code: http.StatusBadRequest,
			Error: "Invalid commentID parameter format",
			Links: links,
		})
	}

	// Populate form
	commentReq := entities.CommentRequest{}
	c.Bind(&commentReq)

	// token
	token := c.Get("user")

	// update comment
	commentRes, err := handler.commentService.Update(commentReq, commentID, token)
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
	}

	// response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code: 201,
		Error: nil,
		Links: links,
		Data: map[string]int{
			"id": int(commentRes.ID),
		},
	})
}

/*
 * Delete Comment
 * -------------------------------
 * Hapus komentar user, hanya pemilik komentar yang dapat mengedit
 */
func (handler CommentHandler) Delete(c echo.Context) error {
	
	// Url path parameter & link hateoas
	commentID, err := strconv.Atoi(c.Param("commentID"))
	links := map[string]string {}
	links["self"] = config.Get().App.BaseURL + "/api/events/comments/" + c.Param("commentID")
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse {
			Status: "OK",
			Code: http.StatusBadRequest,
			Error: "Invalid commentID parameter format",
			Links: links,
		})
	}

	// token
	token := c.Get("user")

	// update comment
	err = handler.commentService.Delete(commentID, token)
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
	}

	// response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code: 201,
		Error: nil,
		Links: links,
		Data: map[string]int{
			"id": commentID,
		},
	})
}