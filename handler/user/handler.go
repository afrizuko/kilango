package user

import (
	"github.com/audit/model"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

type Handler struct {
	http.Handler
	service model.UserService
}

func NewHandler(service model.UserService) *Handler {

	handler := new(Handler)
	handler.service = service

	mux := chi.NewRouter()
	mux.Use(render.SetContentType(render.ContentTypeJSON))

	mux.Get("/", handler.GetUsers)
	mux.Get("/{id}", handler.GetUser)
	mux.Put("/{id}", handler.UpdateUser)
	mux.Post("/", handler.CreateUser)
	mux.Delete("/{id}", handler.DeleteUser)
	mux.Delete("/{id}/purge", handler.PurgeUser)

	handler.Handler = mux
	return handler
}

func DefaultHandler() *Handler {

	handler := new(Handler)
	handler.service = model.NewUserStub()

	mux := chi.NewRouter()
	mux.Use(render.SetContentType(render.ContentTypeJSON))

	mux.Get("/", handler.GetUsers)
	mux.Get("/{id}", handler.GetUser)
	mux.Put("/{id}", handler.UpdateUser)
	mux.Post("/", handler.CreateUser)
	mux.Delete("/{id}", handler.DeleteUser)
	mux.Delete("/{id}/purge", handler.PurgeUser)

	handler.Handler = mux
	return handler
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
