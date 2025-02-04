package ml

import (
	"encoding/json"
	"net/http"
)

type PredictRequest struct {
	ModelType string                 `json:"modelType"`
	Features  map[string]float64     `json:"features"`
}

type PredictResponse struct {
	Predictions map[string]float64   `json:"predictions"`
	Metadata    ModelMetadata        `json:"metadata"`
}

func (s *ModelService) HandlePredict(w http.ResponseWriter, r *http.Request) {
	var req PredictRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	predictions, err := s.Predict(req.ModelType, req.Features)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := PredictResponse{
		Predictions: predictions,
		Metadata:    s.getModelMetadata(req.ModelType),
	}

	json.NewEncoder(w).Encode(resp)
} 