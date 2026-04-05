package appmsg

import "time"

type GameListRequest struct {
}

type GameListResponse struct {
	Games []Game `json:"games"`
}

type Game struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Desc       string    `json:"desc"`
	Icon       string    `json:"icon"`
	Url        string    `json:"url"`
	Status     int       `json:"status"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}
