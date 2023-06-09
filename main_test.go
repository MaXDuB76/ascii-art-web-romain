package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEndpoints(t *testing.T) {
	// Création du routeur HTTP
	router := setupRouter()

	// Créer une requête GET pour l'endpoint "/"
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Vérification du code de statut
	if w.Code != http.StatusOK {
		t.Errorf("GET / returned wrong status code: got %v, want %v", w.Code, http.StatusOK)
	}

	// Créer une requête POST pour l'endpoint "/ascii-art"
	req = httptest.NewRequest("POST", "/ascii-art", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	/// Vérification du code de statut
	if w.Code != http.StatusOK {
		t.Errorf("POST /ascii-art returned wrong status code: got %v, want %v", w.Code, http.StatusOK)
	}
}

// Fonction pour configurer le routeur HTTP
func setupRouter() http.Handler {
	// Création du routeur HTTP
	router := http.NewServeMux()

	// Définir les handlers pour les endpoints
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Handler pour l'endpoint "/"
	})

	router.HandleFunc("/ascii-art", func(w http.ResponseWriter, r *http.Request) {
		// Handler pour l'endpoint "/ascii-art"
	})

	return router
}
