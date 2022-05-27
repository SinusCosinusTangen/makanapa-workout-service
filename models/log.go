package models

import "time"

type Log struct {
	Index   string  `json:"index"`
	Type    string  `json:"type"`
	ID      string  `json:"id"`
	Version int     `json:"version"`
	Score   float64 `json:"score"`
	Source  Source  `json:"source"`
	Fields  Fields  `json:"fields"`
}

type Source struct {
	ProxyHost string    `json:"proxy_host"`
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Version   string    `json:"version"`
	ProxyPort string    `json:"proxy_port"`
	Method    string    `json:"method"`
	Type      string    `json:"type"`
	Host      string    `json:"host"`
	Port      string    `json:"port"`
}

type Fields struct {
	Timestamp []time.Time `json:"timestamp"`
}
