package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	mv "URL-shortener/internal/middleware"
	service "URL-shortener/internal/service"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	LinkService *service.LinkService
	AuthService *service.AuthService
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ShortingRequest struct {
	URL   string `json:"url"`
	Alias string `json:"alias,omitempty"`
	Duration string `json:"duration,omitempty"`
}

func (h *Handler) RegisterRoutes() {
	http.Handle("/save", mv.AuthMiddleware(http.HandlerFunc(h.saveURL), h.AuthService.JWTSecret))
	http.HandleFunc("/", h.redirectURL)
	http.HandleFunc("/register", h.createUser)
	http.HandleFunc("/login", h.login)
	http.Handle("/metrics", promhttp.Handler())
}


func (h *Handler) saveURL(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost{
		var shReq ShortingRequest
		if err := json.NewDecoder(r.Body).Decode(&shReq); err != nil {
			http.Error(w, "Error decoding request", http.StatusBadRequest)
			return 
		}

		userIDctx := r.Context().Value("user_id")
		userIDctxFloat, ok := userIDctx.(float64)
		if !ok {
			http.Error(w, "Invalid user id", http.StatusBadRequest)
			return
		}
		userID := int(userIDctxFloat)

		shortURL, err := h.LinkService.CreateLink(shReq.URL, shReq.Alias, shReq.Duration, userID)
		if err != nil {
			switch{
			case errors.Is(err, service.ErrInvalidURL):
				http.Error(w, err.Error(), http.StatusBadRequest)
			case errors.Is(err, service.ErrInvalidTime):
				http.Error(w, err.Error(), http.StatusBadRequest)
			case errors.Is(err, service.ErrAliasAlreadyExists):
				http.Error(w, err.Error(), http.StatusConflict)
			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}
		w.Write([]byte(shortURL))

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) redirectURL(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		alias := r.URL.Path[1:]
		if len(alias) < 1 {
			http.Error(w, "Invalid key", http.StatusBadRequest)
			return
		}

		originalURL, err := h.LinkService.GetURL(alias)

		if err != nil {
			switch{
			case errors.Is(err, service.ErrLinkExpired):
				http.Error(w, "Link expired", http.StatusGone)
			case errors.Is(err, service.ErrLinkNotFound):
				http.Error(w, "Link not found", http.StatusNotFound)
			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		http.Redirect(w, r, originalURL, http.StatusFound)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var authReq AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
			http.Error(w, "Error decoding request", http.StatusBadRequest)
			return
		}

		if err := h.AuthService.CreateUser(authReq.Username, authReq.Password); err != nil {
			switch{
			case errors.Is(err, service.ErrInvalidPassword):
				http.Error(w, err.Error(), http.StatusUnauthorized)
			case errors.Is(err, service.ErrUserAlreadyExists):
				http.Error(w, "Username is already taken", http.StatusConflict)
			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var authReq AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
			http.Error(w, "Error decoding request", http.StatusBadRequest)
			return
		}

		tokenString, err := h.AuthService.LoginUser(authReq.Username, authReq.Password)
		if err != nil {
			switch{
			case errors.Is(err, service.ErrInvalidCredentials):
				http.Error(w, err.Error(), http.StatusUnauthorized)
			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
