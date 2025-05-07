package api

// Removed unused function performAuthenticatedRequest
/*
func performAuthenticatedRequest(router *gin.Engine, method, url string, body interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	var req *http.Request

	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	// Add authentication header
	req.Header.Set("Authorization", "Bearer test_token")
	req.Header.Set("Content-Type", "application/json")

	// Add user ID to context
	ctx := req.Context()
	ctx = context.WithValue(ctx, "user_id", "test_user_id")
	req = req.WithContext(ctx)

	router.ServeHTTP(w, req)
	return w
}
*/