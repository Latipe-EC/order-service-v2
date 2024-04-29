package msgDTO

var (
	NOTIFY_USER     = 1
	NOTIFY_STORE    = 2
	NOTIFY_CAMPAIGN = 3
)

type SendNotificationMessage struct {
	UserID       string `json:"user_id" validate:"required"`
	Type         int    `json:"type" validate:"required"`
	Title        string `json:"title" validate:"required"`
	Body         string `json:"body" validate:"required"`
	Image        string `json:"image"`
	PushToDevice bool   `json:"push_to_device"`
}

func NewNotificationMessage(_type int, userID string, title string, body string, image string) *SendNotificationMessage {
	return &SendNotificationMessage{
		UserID:       userID,
		Type:         _type,
		Title:        title,
		Body:         body,
		Image:        image,
		PushToDevice: true,
	}
}
