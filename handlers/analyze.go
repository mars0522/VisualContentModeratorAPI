package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"VisulModerator/models"
	"VisulModerator/utils"
)

func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid form data")
		return
	}
	file, _, err := r.FormFile("image")
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Image file is required")
		return
	}
	defer file.Close()

	imageBytes, _ := io.ReadAll(file)
	credsFile := "credentials/your-vertex-creds.json"

	responseMap, err := utils.CallGemini(r.Context(), credsFile, imageBytes)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	flaggedItems, agentDecision := utils.AgentReason(responseMap)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.FinalResponse{
		ApiStatus:  "ok",
		Status:     "success",
		StatusCode: http.StatusOK,
		Response: models.AgenticResponse{
			Classification: responseMap,
			FlaggedItems:   flaggedItems,
			AgentDecision:  agentDecision,
		},
	})
}
