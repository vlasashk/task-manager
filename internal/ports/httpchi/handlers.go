package httpchi

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/vlasashk/todo-manager/internal/adapters/pgrepo"
	"github.com/vlasashk/todo-manager/internal/models/todo"
	"net/http"
	"time"
)

func (s Service) CreateTask(w http.ResponseWriter, r *http.Request) {
	taskRequest := todo.TaskReq{}
	log := *zerolog.Ctx(r.Context())
	if err := render.DecodeJSON(r.Body, &taskRequest); err != nil {
		log.Error().Err(err).Send()
		NewErr("", "", "bad JSON").Send(w, r, http.StatusBadRequest)
		return
	}
	if err := validateDate(taskRequest.DueDate); err != nil {
		log.Error().Err(err).Send()
		NewErr("date", taskRequest.DueDate, "bad date format").Send(w, r, http.StatusBadRequest)
		return
	}
	log.Info().Msg("request body decoded")
	if err := validator.New().Struct(taskRequest); err != nil {
		log.Error().Err(err).Send()
		NewErr("", "", "invalid JSON").Send(w, r, http.StatusUnprocessableEntity)
		return
	}
	newTask, err := s.DB.CreateTask(taskRequest)
	if err != nil {
		errorHandler(w, r, log, taskRequest.DueDate, "", err)
		return
	}
	log.Info().Msg("task created successfully")
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, newTask)
}

func (s Service) GetSingleTask(w http.ResponseWriter, r *http.Request) {
	log := *zerolog.Ctx(r.Context())
	taskID := chi.URLParam(r, "id")
	log.Info().Str("id", taskID).Msg("task id received")
	task, err := s.DB.GetTask(taskID)
	if err != nil {
		logID := log.With().Str("id", taskID).Logger()
		errorHandler(w, r, logID, "", taskID, err)
		return
	}
	log.Info().Str("id", taskID).Msg("received successfully")
	render.Status(r, http.StatusOK)
	render.JSON(w, r, task)
}

func (s Service) DeleteTask(w http.ResponseWriter, r *http.Request) {
	log := *zerolog.Ctx(r.Context())
	taskID := chi.URLParam(r, "id")
	log.Info().Str("id", taskID).Msg("task id received")
	err := s.DB.DeleteTask(taskID)
	if err != nil {
		logID := log.With().Str("id", taskID).Logger()
		errorHandler(w, r, logID, "", taskID, err)
		return
	}
	log.Info().Str("id", taskID).Msg("deleted successfully")
	NewMsg("success").Send(w, r, http.StatusOK)
}

func (s Service) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskUpd := todo.TaskReq{}
	log := *zerolog.Ctx(r.Context())
	taskID := chi.URLParam(r, "id")
	log.Info().Str("id", taskID).Msg("task id received")
	if err := render.DecodeJSON(r.Body, &taskUpd); err != nil {
		log.Error().Err(err).Send()
		NewErr("", "", "bad JSON").Send(w, r, http.StatusBadRequest)
		return
	}
	if err := validateDate(taskUpd.DueDate); err != nil {
		log.Error().Err(err).Send()
		NewErr("date", taskUpd.DueDate, "bad date format").Send(w, r, http.StatusBadRequest)
		return
	}
	log.Info().Msg("request body decoded")
	if err := validator.New().Struct(taskUpd); err != nil {
		log.Error().Err(err).Send()
		NewErr("", "", "invalid JSON").Send(w, r, http.StatusUnprocessableEntity)
		return
	}
	newTask, err := s.DB.UpdateTask(taskUpd, taskID)
	if err != nil {
		logID := log.With().Str("id", taskID).Logger()
		errorHandler(w, r, logID, taskUpd.DueDate, taskID, err)
		return
	}
	log.Info().Msg("task updated successfully")
	render.Status(r, http.StatusOK)
	render.JSON(w, r, newTask)
}

func validateDate(date string) error {
	layout := "2006-01-02"
	_, err := time.Parse(layout, date)
	return err
}

func errorHandler(w http.ResponseWriter, r *http.Request, log zerolog.Logger, date, taskID string, err error) {
	switch err.Error() {
	case pgrepo.InvalidIdErr:
		log.Warn().Err(err).Send()
		NewMsg(pgrepo.InvalidIdErr).Send(w, r, http.StatusBadRequest)
	case pgrepo.DateErr:
		log.Warn().Err(err).Send()
		NewErr("date", date, pgrepo.DateErr).Send(w, r, http.StatusBadRequest)
	default:
		log.Error().Err(err).Send()
		NewErr("id", taskID, "action fail").Send(w, r, http.StatusInternalServerError)
	}
}
