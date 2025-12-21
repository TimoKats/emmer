package main

type RequestConfig struct {
	Method         string
	Endpoint       string
	Body           string
	ExpectedStatus int
}
