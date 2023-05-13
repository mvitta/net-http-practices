package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type task struct {
	ID      int    `json:"ID"`
	Name    string `json:"Name"`
	Content string `json:"Content"`
}

var tasks = []task{
	{
		ID:      1,
		Name:    "Task One",
		Content: "Some Content",
	},
}

func Json(res http.ResponseWriter) []byte {
	tasks, err := json.Marshal(tasks)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)

	}
	return tasks
}

func MapToJson(m map[string]string, res http.ResponseWriter) []byte {
	b, err := json.Marshal(m)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	return b
}

func ToJson(res http.ResponseWriter, t task) []byte {
	json, error := json.Marshal(t)
	if error != nil {
		http.Error(res, error.Error(), http.StatusInternalServerError)
	}
	return json
}

func Parse(res http.ResponseWriter, b []byte) task {
	var newTask task
	errJson := json.Unmarshal(b, &newTask)
	if errJson != nil {
		http.Error(res, errJson.Error(), http.StatusInternalServerError)
	}
	return newTask
}

func SearchTask(id int) (task, int, error) {
	if id < 0 || id >= len(tasks) {
		return task{}, 0, errors.New("Index out of range")
	}
	i := 0
	for ; i < len(tasks); i++ {
		if tasks[i].ID == id {
			break
		}
	}
	return tasks[i], i, nil
}

func remove(index int, res http.ResponseWriter) {
	tasks = append(tasks[:index], tasks[index+1:]...)

}

func main() {

	//routes
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		url := req.URL.Path
		if url != "/" {
			http.Error(res, "404 not found", http.StatusNotFound)
			return
		}

		response := []byte(`
		{
			"status": "200",
			"message": "Hello world !"
		}
		`)

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(response)

	})

	http.HandleFunc("/addTask", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			b, errBody := io.ReadAll(req.Body)
			if errBody != nil {
				http.Error(res, errBody.Error(), http.StatusInternalServerError)

			} else {

				defer req.Body.Close()
				newTask := Parse(res, b)
				tasks = append(tasks, newTask)

				res.Header().Set("Content-Type", "application/json")
				res.WriteHeader(http.StatusOK)
				res.Write([]byte(`{"status":"add to list of task"}`))
			}

		} else {
			http.Error(res, "Invalid Request Method", http.StatusMethodNotAllowed)
		}

	})

	http.HandleFunc("/getTasks", func(res http.ResponseWriter, req *http.Request) {

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		tasks := Json(res)
		res.Write(tasks)

	})

	http.HandleFunc("/deleteTask", func(res http.ResponseWriter, req *http.Request) {

		id, error := strconv.Atoi(req.URL.Query().Get("id"))
		if error != nil {
			http.Error(res, error.Error(), http.StatusInternalServerError)
		} else {
			_, index, errorIndexSarch := SearchTask(id)
			if errorIndexSarch != nil {
				http.Error(res, errorIndexSarch.Error(), http.StatusInternalServerError)

			} else {
				remove(index, res)
				res.Header().Set("Content-Type", "application/json")
				res.Write([]byte(`{"Status":"Deleted" }`))

			}
		}

	})

	//listen server
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error en ListenAndServer")
	}
}
