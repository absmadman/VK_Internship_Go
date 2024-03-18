package db

type User struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type Quest struct {
	Id   int     `json:"id"`
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

type Event struct {
	UserId      int     `json:"user_id"`
	QuestId     int     `json:"quest_id"`
	UserBalance float64 `json:"user_balance"`
}

type EventResponse struct {
	UserId  int     `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	Quests  []int   `json:"quests"`
}
