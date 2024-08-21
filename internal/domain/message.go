package domain

import "time"

type Id string

type Message struct {
	Id                 Id        `json:"id"`
	Message            string    `json:"message"`
	Processed          bool      `json:"processed"`
	DateCreate         time.Time `json:"date_create"`
	DateProcessedStart time.Time `json:"date_processed_start"`
	DateProcessedEnd   time.Time `json:"date_processed_end"`
}
