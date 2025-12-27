func (h *Handler) CreateTwite(w http.ResponseWriter, r *http.Request) {
	var req struct {
			Content string `json:"content"`
	}

	json.NewDecoder(r.Body).Decode(&req)
	authorID := r.Context().Value("userId").(string)

	err := h.service.CreateTwite(r.Context(), authorID, req.Content)
	if err != nil {
			http.Error(w, err.Error(), 500)
			return
	}

	w.WriteHeader(http.StatusCreated)
}