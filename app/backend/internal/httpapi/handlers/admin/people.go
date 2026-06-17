// httpapi/handlers/admin/people.go
package admin

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ajn2004/daleego-hearth/backend/internal/db"
	"github.com/ajn2004/daleego-hearth/backend/internal/httpapi/response"
	httputil "github.com/ajn2004/daleego-hearth/backend/internal/httpapi/utils"
	"github.com/go-chi/chi/v5"
)

/*
API Structure

GET /admin/people
  Description:
    Lists all people.

POST /admin/people
  Payload:
    {
      "display_name": "Andrew",
      "role": "admin" // optional: "admin" | "member"
    }
  Description:
    Creates a new person.
  Returns:
    New person database object.

GET /admin/people/{person_id}
  Description:
    Gets one person by ID.
  Returns:
    Person database object.

PATCH /admin/people/{person_id}
  Payload:
    {
      "display_name": "Andrew Nelson"
    }
  Description:
    Updates a person's editable fields.
  Returns:
    Updated person database object.

DELETE /admin/people/{person_id}
  Description:
    Soft-deletes a person by setting deleted_at.
  Returns:
    Updated person database object.
*/

func (h *Handler) registerPeopleRoutes(r chi.Router) {
	r.Get("/people", h.ListPeople)
	r.Post("/people", h.CreatePerson)
	r.Get("/people/{person_id}", h.GetPerson)
	r.Patch("/people/{person_id}", h.UpdatePerson)
	r.Delete("/people/{person_id}", h.DeletePerson)
}

func (h *Handler) ListPeople(w http.ResponseWriter, r *http.Request) {
	people, err := h.queries.AdminGetAllPeople(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to list people")
	}
	response.WriteJSON(w, http.StatusOK, people)
}

type PersonRequest struct {
	DisplayName string `json:"display_name"`
}

func (h *Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var req PersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.DisplayName) == "" {
		response.WriteError(w, http.StatusBadRequest, "display_name is required")
		return
	}

	person, err := h.queries.CreatePerson(r.Context(), strings.TrimSpace(req.DisplayName))

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "could not create person")
		return
	}

	response.WriteJSON(w, http.StatusCreated, person)
}

func (h *Handler) GetPerson(w http.ResponseWriter, r *http.Request) {
	personID, err := httputil.ParseUUIDParam(chi.URLParam(r, "person_id"))
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "invalid person id")
		return
	}

	person, err := h.queries.GetPersonByID(r.Context(), personID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "invalid person id")
		return
	}
	response.WriteJSON(w, http.StatusOK, person)
}

func (h *Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	personID, err := httputil.ParseUUIDParam(chi.URLParam(r, "person_id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid person id")
		return
	}

	var req PersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.DisplayName) == "" {
		response.WriteError(w, http.StatusBadRequest, "display_name is required")
		return
	}

	person, err := h.queries.UpdatePersonName(r.Context(), db.UpdatePersonNameParams{
		PersonID:    personID,
		DisplayName: strings.TrimSpace(req.DisplayName),
	})

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "could not update person name")
		return
	}

	response.WriteJSON(w, http.StatusOK, person)
}

func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	personID, err := httputil.ParseUUIDParam(chi.URLParam(r, "person_id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid person id")
		return
	}

	person, err := h.queries.SetPersonToDeleted(r.Context(), personID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "could not delete person")
		return
	}

	response.WriteJSON(w, http.StatusOK, person)
}
