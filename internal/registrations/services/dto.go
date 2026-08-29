package services

type RegisterMemberRequest struct {
	Name             string `json:"Name" binding:"required"`
	Email            string `json:"Email" binding:"required,email"`
	Phone            string `json:"Phone" binding:"required"`
	Type             string `json:"Type" binding:"required"` // e.g., member or guest
	EventName        string `json:"EventName" binding:"required"`
	Rating           string `json:"Rating"`
	LookingForwardTo string `json:"LookingForwardTo"`
	Suggestion       string `json:"Suggestion"`
}
