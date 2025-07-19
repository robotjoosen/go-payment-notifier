package bunq

type NotificationFilters struct {
	NotificationFilter []NotificationFilter `json:"notification_filters"`
}

type NotificationFilter struct {
	Category           string `json:"category"`
	NotificationTarget string `json:"notification_target"`
}
