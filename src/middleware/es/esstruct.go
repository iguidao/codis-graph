package es

import "time"

type EsSource struct {
	Clientip  string    `json:"clientip"`
	Action    string    `json:"action"`
	Codisname string    `json:"codisname"`
	Codisip   string    `json:"codisip"`
	Timestamp time.Time `json:"@timestamp"`
}

type SouHits struct {
	Index    string   `json:"_index"`
	Type     string   `json:"_type"`
	Id       string   `json:"_id"`
	Score    float64  `json:"_score"`
	EsSource EsSource `json:"_source"`
}

type AllHits struct {
	Total    int       `json:"total"`
	MaxScore float64   `json:"max_score"`
	SouHits  []SouHits `json:"hits"`
}

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}
type EsData struct {
	ScrollId string  `json:"_scroll_id"`
	Took     int     `json:"took"`
	TimeOut  bool    `json:"timed_out"`
	Shards   Shards  `json:"_shards"`
	AllHits  AllHits `json:"hits"`
}
