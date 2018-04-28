package main

import (
	"net/http"
	"encoding/json"
	"log"
	"io"
	"myworker"
)

var (
	MaxWorker = 10
	MaxQueue  = 10
	MaxLength = int64(1024)
)

func main() {
	service := myworker.NewService(MaxWorker, MaxQueue)

	service.Start()
	defer service.Stop()

	// 处理海量的任务
	http.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		log.Println("ok")

		// Job以JSON格式提交
		var jobs []myworker.Job
		err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&jobs)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// 处理任务
		for _, job := range jobs {
			service.AddJob(job)
		}

		// OK
		w.WriteHeader(http.StatusOK)
	})

	// 启动web服务
	log.Fatal(http.ListenAndServe(":8080", nil))
}
