package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Handler struct {
	Store Store
	Log   zerolog.Logger

	IsReadable bool
	IsWritable bool
}

func (s *Handler) handleGet(res http.ResponseWriter, req *http.Request, desc Description) error {
	if !s.IsReadable {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}

	err := s.Store.Get(req.Context(), desc, res)
	if err == nil {
		return nil
	}

	if errors.Is(err, ErrNotExist) {
		res.WriteHeader(http.StatusNotFound)
	}

	return err
}

func (s *Handler) handleHead(res http.ResponseWriter, req *http.Request, desc Description) error {
	size, err := s.Store.Head(req.Context(), desc)
	if err == nil {
		res.Header().Add("Content-Length", strconv.FormatInt(int64(size), 10))
		res.WriteHeader(http.StatusOK)
		return nil
	}

	if errors.Is(err, ErrNotExist) {
		res.WriteHeader(http.StatusNotFound)
		return nil
	}

	return err
}

func (s *Handler) handlePut(res http.ResponseWriter, req *http.Request, desc Description) error {
	if !s.IsWritable {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}

	err := s.Store.Put(req.Context(), desc, req.Body)
	if err == nil {
		res.WriteHeader(http.StatusOK)
		return nil
	}

	if errors.Is(err, ErrExist) {
		res.WriteHeader(http.StatusConflict)
		return nil
	}

	return err
}

type responseWriter struct {
	http.ResponseWriter
	status_code int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status_code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (s *Handler) parseDescription(res http.ResponseWriter, req *http.Request) (Description, error) {
	entries := strings.SplitN(req.URL.Path[1:], "/", 5)
	if len(entries) != 4 {
		res.WriteHeader(http.StatusNotFound)
		return Description{}, errors.New("invalid path")
	}

	switch req.Method {
	case http.MethodGet:
	case http.MethodHead:
	case http.MethodPut:
		break

	default:
		res.WriteHeader(http.StatusNotImplemented)
		return Description{}, errors.New("invalid method")
	}

	return Description{
        Triplet: entries[0],
		Name:    entries[1],
		Version: entries[2],
		Hash:    entries[3],
	}, nil
}

func (s *Handler) ServeHTTP(r http.ResponseWriter, req *http.Request) {
	if (req.URL.Path == "/") && (req.Method == http.MethodGet) {
		s.Log.Info().Msg("probe")
		r.WriteHeader(http.StatusOK)
		return
	}

	t0 := time.Now()
	res := &responseWriter{r, http.StatusOK}

	l := s.Log.With().Str("_", getTicket()).Logger()
	req = req.WithContext(l.WithContext(req.Context()))

	desc, err := s.parseDescription(res, req)
	{
		l := l.With().Str("url", req.URL.String()).Str("method", req.Method).Logger()
		if err != nil {
			l.Warn().Dur("dt", time.Since(t0)).Int("status", res.status_code).Msg("REQ " + err.Error())
			return
		}

		l.Info().Msg("")
	}

	l.Info().
		Str("triplet", desc.Triplet).
		Str("name", desc.Name).
		Str("version", desc.Version).
		Str("hash", desc.Hash).
		Msg("REQ " + req.Method)

	err = nil
	switch req.Method {
	case http.MethodGet:
		err = s.handleGet(res, req, desc)

	case http.MethodHead:
		err = s.handleHead(res, req, desc)

	case http.MethodPut:
		err = s.handlePut(res, req, desc)
	}

	l = l.With().Dur("dt", time.Since(t0)).Int("status", res.status_code).Logger()
	msg := "RES " + req.Method

	if err != nil {
		l.Error().Err(err).Msg(msg)
		return
	}
	if res.status_code >= 400 {
		l.Warn().Msg(msg)
		return
	}

	l.Info().Msg(msg)
}
