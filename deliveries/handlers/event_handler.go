package handlers

import (
	"fmt"
	"reflect"
	"strconv"
	"tupulung/config"
	"tupulung/deliveries/helpers"
	"tupulung/entities"
	"tupulung/entities/web"
	eventService "tupulung/services/event"

	"github.com/labstack/echo/v4"
)

type EventHandler struct {
	eventService *eventService.EventService
}

func NewEventHandler(service *eventService.EventService) *EventHandler {
	return &EventHandler{
		eventService: service,
	}
}

/*
 * -------------------------------------------
 * Show All events based on available queries
 * -------------------------------------------
 */
func (handler EventHandler) Index(c echo.Context) error {

	// Translate query param to map of filters
	filters := []map[string]string{}
	q := c.QueryParam("q")
	if q != "" {
		filters = append(filters, map[string]string{
			"field":    "title",
			"operator": "LIKE",
			"value":    "%" + q + "%",
		})
	}
	category_id := c.QueryParam("category_id")
	if category_id != "" {
		filters = append(filters, map[string]string{
			"field":    "category_id",
			"operator": "=",
			"value":    category_id,
		})
	}
	fmt.Println(category_id)
	// Sort parameter
	sorts := []map[string]interface{}{}
	sortLocation := c.QueryParam("sortLocation")
	if sortLocation != "" {
		switch sortLocation {
		case "1":
			sorts = append(sorts, map[string]interface{}{
				"field": "title",
				"desc":  true,
			})
		case "0":
			sorts = append(sorts, map[string]interface{}{
				"field": "title",
				"desc":  false,
			})
		}
	}
	links := map[string]string{"self": config.Get().App.BaseURL + "/api/events?limit=" + c.QueryParam("limit") + "&page=" + c.QueryParam("page")}

	// pagination param
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(400, helpers.MakeErrorResponse("ERROR", 400, "Limit Parameter format is invalid", links))
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		links := map[string]string{"self": config.Get().App.BaseURL}
		return c.JSON(400, helpers.MakeErrorResponse("ERROR", 400, "page Parameter format is invalid", links))
	}
	links["self"] = config.Get().App.BaseURL + "/api/events?limit=" + c.QueryParam("limit") + "&page=" + c.QueryParam("page")

	// Get all events
	eventsRes, err := handler.eventService.FindAll(limit, page, filters, sorts)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, helpers.MakeErrorResponse("ERROR", webErr.Code, webErr.Error(), links))
		}
		panic("not returning custom error")
	}

	// Get pagination data
	pagination, err := handler.eventService.GetPagination(limit, page, filters)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, helpers.MakeErrorResponse("ERROR", webErr.Code, webErr.Error(), links))
		}
		panic("not returning custom error")
	}

	links["first"] = config.Get().App.BaseURL + "/api/events?limit=" + c.QueryParam("limit") + "&page=1"
	links["last"] = config.Get().App.BaseURL + "/api/events?limit=" + c.QueryParam("limit") + "&page=" + strconv.Itoa(pagination.TotalPages)
	if pagination.Page > 1 {
		links["prev"] = config.Get().App.BaseURL + "/api/events?limit=" + c.QueryParam("limit") + "&page=" + strconv.Itoa(pagination.Page-1)
	}
	if pagination.Page < pagination.TotalPages {
		links["next"] = config.Get().App.BaseURL + "/api/events?limit=" + c.QueryParam("limit") + "&page=" + strconv.Itoa(pagination.Page+1)
	}

	// success response
	return c.JSON(200, web.SuccessListResponse{
		Status:     "OK",
		Code:       200,
		Error:      nil,
		Links:      links,
		Data:       eventsRes,
		Pagination: pagination,
	})
}

/*
 * -------------------------------------------
 * Show single event detail by ID
 * -------------------------------------------
 */
func (handler EventHandler) Show(c echo.Context) error {
	// Get param
	id, err := strconv.Atoi(c.Param("id"))
	links := map[string]string{"self": config.Get().App.BaseURL + "/events/" + c.Param("id")}
	if err != nil {
		return c.JSON(400, helpers.MakeErrorResponse("ERROR", 400, err.Error(), links))
	}
	// Get productdata
	event, err := handler.eventService.Find(id)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, helpers.MakeErrorResponse("ERROR", webErr.Code, webErr.Error(), links))
		}
	}
	// response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code:   200,
		Error:  nil,
		Links:  links,
		Data:   event,
	})
}

/*
 * -------------------------------------------
 * Get All user's events based on available queries
 * -------------------------------------------
 */
func (handler EventHandler) GetUserEvent(c echo.Context) error {

	// Translate query param to map of filters
	filters := []map[string]string{}
	q := c.QueryParam("q")
	if q != "" {
		filters = append(filters, map[string]string{
			"field":    "title",
			"operator": "LIKE",
			"value":    "%" + q + "%",
		})
	}
	// Sort parameter
	sorts := []map[string]interface{}{}
	sortLocation := c.QueryParam("sortLocation")
	if sortLocation != "" {
		switch sortLocation {
		case "1":
			sorts = append(sorts, map[string]interface{}{
				"field": "location",
				"desc":  true,
			})
		case "0":
			sorts = append(sorts, map[string]interface{}{
				"field": "location",
				"desc":  false,
			})
		}
	}
	links := map[string]string{"self": config.Get().App.BaseURL + "/api/events"}

	// pagination param
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(400, helpers.MakeErrorResponse("ERROR", 400, "Limit Parameter format is invalid", links))
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		links := map[string]string{"self": config.Get().App.BaseURL}
		return c.JSON(400, helpers.MakeErrorResponse("ERROR", 400, "page Parameter format is invalid", links))
	}

	// get user param ID
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(400, helpers.MakeErrorResponse("ERROR", 400, "requested id is invalid", links))
	}
	filters = append(filters, map[string]string{
		"field":    "user_id",
		"operator": "=",
		"value":    strconv.Itoa(userID),
	})

	// Get all events
	eventsRes, err := handler.eventService.FindAll(limit, page, filters, sorts)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, helpers.MakeErrorResponse("ERROR", webErr.Code, webErr.Error(), links))
		}
		panic("not returning custom error")
	}

	// Get pagination data
	pagination, err := handler.eventService.GetPagination(limit, page, filters)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, helpers.MakeErrorResponse("ERROR", webErr.Code, webErr.Error(), links))
		}
		panic("not returning custom error")
	}

	// success response
	return c.JSON(200, web.SuccessListResponse{
		Status:     "OK",
		Code:       200,
		Error:      nil,
		Links:      links,
		Data:       eventsRes,
		Pagination: pagination,
	})
}

/*
 * -------------------------------------------
 * Create event resource
 * -------------------------------------------
 */
func (handler EventHandler) Create(c echo.Context) error {
	// Populate form
	eventReq := entities.EventRequest{}
	c.Bind(&eventReq)

	// Define hateoas links
	links := map[string]string{"self": config.Get().App.BaseURL + "/events"}

	token := c.Get("user")

	// Insert event
	eventRes, err := handler.eventService.Create(eventReq, token)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, helpers.MakeErrorResponse("ERROR", webErr.Code, webErr.Error(), links))
		}
	}

	// response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code:   200,
		Error:  nil,
		Links:  links,
		Data:   eventRes,
	})
}

/*
 * -------------------------------------------
 * Update event resource
 * -------------------------------------------
 */
func (handler EventHandler) Update(c echo.Context) error {
	// Populate form
	eventReq := entities.EventRequest{}
	c.Bind(&eventReq)

	id, err := strconv.Atoi(c.Param("eventID"))
	links := map[string]string{"self": config.Get().App.BaseURL + "/events/" + c.Param("id")}
	if err != nil {
		return c.JSON(400, helpers.MakeErrorResponse("ERROR", 400, err.Error(), links))
	}

	token := c.Get("user")

	// Product service call
	eventRes, err := handler.eventService.Update(eventReq, id, token)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, helpers.MakeErrorResponse("ERROR", webErr.Code, webErr.Error(), links))
		}
	}

	// response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code:   200,
		Error:  nil,
		Links:  links,
		Data:   eventRes,
	})
}

/*
 * -------------------------------------------
 * Delete event resource
 * -------------------------------------------
 */
func (handler EventHandler) Delete(c echo.Context) error {

	// Get params ID
	id, err := strconv.Atoi(c.Param("eventID"))
	links := map[string]string{"self": config.Get().App.BaseURL + "/events/" + c.Param("id")}
	if err != nil {
		return c.JSON(400, helpers.MakeErrorResponse("ERROR", 400, err.Error(), links))
	}

	token := c.Get("user")

	// call delete on event service
	err = handler.eventService.Delete(id, token)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, helpers.MakeErrorResponse("ERROR", webErr.Code, webErr.Error(), links))
		}
	}

	// response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code:   200,
		Error:  nil,
		Links:  links,
		Data: map[string]interface{}{
			"id": id,
		},
	})
}