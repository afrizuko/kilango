package user

import (
	"github.com/afrizuko/kilango/model"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"os"
	"strconv"
)

type Handler struct {
	http.Handler
	service model.UserService
}

func NewHandler(service model.UserService) *Handler {

	handler := new(Handler)
	handler.service = service

	handler.AddRoutes()
	return handler
}

func DefaultHandler() *Handler {

	handler := new(Handler)
	handler.service = model.NewUserServiceImpl()

	handler.AddRoutes()
	return handler
}

func (h *Handler) AddRoutes() {
	mux := chi.NewRouter()
	if os.Getenv("STATE") != "prod" {
		mux.Use(render.SetContentType(render.ContentTypeJSON))
	}

	mux.Get("/", h.GetUsers)
	mux.Get("/{id}", h.GetUser)
	mux.Put("/{id}", h.UpdateUser)
	mux.Post("/", h.CreateUser)
	mux.Delete("/{id}", h.DeleteUser)
	mux.Delete("/{id}/purge", h.PurgeUser)

	h.Handler = mux
}

func (h *Handler) GetUsers(writer http.ResponseWriter, request *http.Request) {
	page, limit := h.getPageQuery(request)
	users, _ := h.service.GetUsers(page, limit)
	render.JSON(writer, request, users)
}

func (h *Handler) GetUser(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.ParseUint(chi.URLParam(request, "id"), 10, 64)
	user, err := h.service.GetUser(uint(id))
	if err != nil {
		render.Status(request, http.StatusNotFound)
		render.JSON(writer, request, map[string]string{"error": err.Error()})
		return
	}
	render.JSON(writer, request, user)
}

func (h *Handler) CreateUser(writer http.ResponseWriter, request *http.Request) {

	user := new(model.User)
	_ = render.Bind(request, user)
	_ = h.service.CreateUser(user)

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, user)
}

func (h *Handler) UpdateUser(writer http.ResponseWriter, request *http.Request) {
	user := new(model.User)
	_ = render.Bind(request, user)

	id, _ := strconv.ParseUint(chi.URLParam(request, "id"), 10, 64)
	_ = h.service.ModifyUser(uint(id), user)

	render.Status(request, http.StatusAccepted)
	render.JSON(writer, request, user)
}

func (h *Handler) DeleteUser(writer http.ResponseWriter, request *http.Request) {

	id, _ := strconv.ParseUint(chi.URLParam(request, "id"), 10, 64)
	_ = h.service.DeleteUser(uint(id))

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, nil)
}

func (h *Handler) getPageQuery(request *http.Request) (int, int) {
	page, _ := strconv.Atoi(request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(request.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}
	return page, limit
}

func (h *Handler) PurgeUser(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.ParseUint(chi.URLParam(request, "id"), 10, 64)
	if err := h.service.PurgeUser(uint(id)); err != nil {
		render.Status(request, http.StatusNotFound)
		render.JSON(writer, request, err)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, nil)
}
