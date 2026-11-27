package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/ElSennio/go-musthave-shortener-tpl/internal/config"
	"github.com/go-chi/chi/v5"
)

type URLService struct {
	urls        map[string]string
	idGenerator func() (string, error)
}

func generateID() (string, error) {
	b := make([]byte, 6)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (service *URLService) generateID() (string, error) {
	if service.idGenerator != nil {
		return service.idGenerator()
	}

	return generateID()
}

func (service *URLService) getMethod(res http.ResponseWriter, req *http.Request) {

	id := req.PathValue("id")
	originalURL, exists := service.urls[id]
	if !exists {
		http.NotFound(res, req)
		return
	}
	http.Redirect(res, req, originalURL, http.StatusTemporaryRedirect)

}

func (service *URLService) postMethod(res http.ResponseWriter, req *http.Request) {

	// 1. Читаем исходный URL из тела запроса
	data, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	originalURL := string(data)

	// 2. Генерируем случайную строку id, которую затем положим в мапу
	strID, err := service.generateID()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	service.urls[strID] = originalURL

	body := config.ConfigFlags.BaseURL + "/" + strID

	// 3. Выводим в печать для клиента
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	fmt.Fprint(res, body)
}

func main() {

	flag.Parse()

	service := &URLService{
		urls: make(map[string]string),
	}
	r := chi.NewRouter()

	r.Post("/", service.postMethod)
	r.Get("/{id}", service.getMethod)

	err := http.ListenAndServe(config.ConfigFlags.ServerAddress, r)
	if err != nil {
		panic(err)
	}
}
