package ml

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
)

func RegisterHandlers(r *mux.Router, service *ModelService) {
    r.HandleFunc("/api/ml/predict/{modelType}", service.HandlePredict).Methods("POST")
    r.HandleFunc("/api/ml/train/{modelType}", service.HandleTrain).Methods("POST")
    r.HandleFunc("/api/ml/models/{modelType}/status", service.HandleModelStatus).Methods("GET")
}

func (s *ModelService) HandleTrain(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    modelType := vars["modelType"]

    var req struct {
        Data []TrainingData `json:"data"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := s.Train(modelType, req.Data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "status": "success",
        "model":  modelType,
    })
}

func (s *ModelService) HandleModelStatus(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    modelType := vars["modelType"]

    model, exists := s.models[modelType]
    if !exists {
        http.Error(w, "model not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "version":      model.Version,
        "lastTrained": model.LastTrained,
        "accuracy":    model.Accuracy,
        "samples":     model.Samples,
    })
} 