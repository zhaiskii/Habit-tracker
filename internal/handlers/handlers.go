package handlers

import(
	"fmt"
	"net/http"
	"math/rand"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"habit21/internal/storage"
)

type HabitRequest struct {
	Name string `json:"name"`
	Progress int `json:"progress"`
	Completed bool `json:"completedToday"`
}

type HabitResponse struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Progress int `json:"progress"`
	Completed bool `json:"completedToday"`
}

type Handler struct {
	Storage *storage.Storage
}

func (conn Handler) Create(w http.ResponseWriter,r *http.Request) {
	var cur HabitRequest
	err := json.NewDecoder(r.Body).Decode(&cur)
	if err!=nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	id := rand.Int()
	id%=500000
	fmt.Println(cur)
	err = conn.Storage.Create(id, cur.Name, cur.Progress, cur.Completed)
	if err!=nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	cur2 := HabitResponse{ ID: id, Name: cur.Name, Progress: cur.Progress, Completed: cur.Completed }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cur2)
}

func (conn Handler) Delete(w http.ResponseWriter,r *http.Request) {
	id := chi.URLParam(r,"id")
	err := conn.Storage.Delete(id)
	if err!=nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (conn Handler) Show(w http.ResponseWriter,r *http.Request) {
	rows, err := conn.Storage.Show()
	w.Header().Set("Content-Type", "application/json")
	if err!=nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	res := []HabitResponse{}
	fmt.Println(res)
	for rows.Next() {
		var cur HabitResponse
		rows.Scan(&cur.ID, &cur.Name, &cur.Progress, &cur.Completed)
		res=append(res,cur)
	}
	fmt.Println(res)
	json.NewEncoder(w).Encode(res)
}

func (conn Handler) Update(w http.ResponseWriter, r *http.Request) {
	var nw HabitRequest
	err := json.NewDecoder(r.Body).Decode(&nw)
	id := chi.URLParam(r,"id")
	if err!=nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	err = conn.Storage.Update(id,nw.Progress,nw.Completed)
	if err!=nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	brbr:=0
	for i:=0; i<len(id); i++ {
		brbr*=10
		brbr+=int(id[i]-'0')
	}
	fmt.Println(id, brbr)
	fmt.Println(nw)
	nw2 := HabitResponse{ ID: brbr, Name: nw.Name, Progress: nw.Progress, Completed: nw.Completed }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nw2)
}
