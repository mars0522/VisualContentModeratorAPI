package utils

import (
	"VisulModerator/models"
)

func AgentReason(classification map[string]interface{}) ([]models.FlaggedItem, models.AgentDecision) {
	flagged := []models.FlaggedItem{}
	safe := true
	var objections []string
	for label, details := range classification {
		if dmap, ok := details.(map[string]interface{}); ok && dmap["flagged"] == true {
			reason, _ := dmap["reason"].(string)
			flagged = append(flagged, models.FlaggedItem{
				Label:  label,
				Reason: reason,
			})
			objections = append(objections, label)
			safe = false
		}
	}

	decision := models.AgentDecision{
		Safe:        safe,
		Explanation: "",
		Objections:  objections,
	}
	if safe {
		decision.Actions = []string{"store_safely"}
		decision.Explanation = "Image passed all moderation checks."
	} else {
		decision.Actions = []string{"notify_moderator", "log_violation"}
		decision.Explanation = "Unsafe content found in the image. Action required."
	}
	return flagged, decision
}
