package main

import "time"

type User struct {
	ID                 int
	Name               string
	Email              string
	Password           string
	ConfirmToken       string
	Confirmed          bool
	AttemptNumber      int64
	AttemptTime        time.Time
	Locked             time.Time
	RecoverToken       string
	RecoverTokenExpiry time.Time
}

type Token struct {
	Key   string
	Token string
}

type Digit struct {
	Position int64
	Digit    int64
}