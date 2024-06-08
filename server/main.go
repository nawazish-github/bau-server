package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nawazish-github/bau/server/handlers"
	"github.com/rs/cors"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	fmt.Println("Starting BAU server...")
	r := mux.NewRouter()
	r.HandleFunc("/login", login)
	r.HandleFunc("/", hello)
	r.HandleFunc("/bau/services", handlers.Services)
	r.HandleFunc("/bau/services/purchase", handlers.NewPurchase)
	r.HandleFunc("/bau/services/purchase/search", handlers.SearchPurchase)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "DELETE", "POST", "PUT", "HEAD", "OPTIONS"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)

	//time.Sleep(3 * time.Second)
	fmt.Println("BAU server listening on port 3333")
	//time.Sleep(1 * time.Second)
	fmt.Println("Access service on http://localhost:3333")
	err := http.ListenAndServe(":3333", handler)
	if err != nil {
		panic(err)
	}
}

var c = 0

func login(w http.ResponseWriter, r *http.Request) {
	c += 1
	fmt.Println("server counter value is: ", c)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only HTTP POST request is allowed"))
	}
	data, err := io.ReadAll(io.Reader(r.Body))
	if err != nil {
		w.Write([]byte("Unable to read request body: " + err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	l := LoginData{}
	err = json.Unmarshal(data, &l)

	if err != nil {
		w.Write([]byte("Unable to unmarshal request body: " + err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("Login request is: ", string(data))
	if (c % 2) == 0 {
		resp := &LoginResp{
			Status:  "successful",
			Message: "",
		}
		data, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Marshalling login successful response failed: " + err.Error()))
		}
		w.Write([]byte(data))
		w.WriteHeader(http.StatusOK)
	} else {
		resp := &LoginResp{
			Status:  "unsuccessful",
			Message: "Login failed. Username/password invalid. Try again...",
		}
		data, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Marshalling login unsuccessful response failed: " + err.Error()))
		}
		w.Write([]byte(data))
		w.WriteHeader(http.StatusUnauthorized)
	}
}
func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, BAU project!\n")
}
