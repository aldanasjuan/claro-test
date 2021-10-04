package main

import (
	"errors"
	"log"
	"runtime"
	"strings"
)

var ErrMissingRequired = errors.New("missing required information")
var ErrUsernameTaken = errors.New("username taken")
var ErrAlreadyLogged = errors.New("user already logged in")
var ErrWrongSignatureFormat = errors.New("wrong signature format")
var ErrInvalidSignature = errors.New("invalid signature")

func Log(err error) {
	_, fl, line, _ := runtime.Caller(1)

	file := strings.SplitAfterN(fl, "/", 2)
	var f *string
	if len(file) > 1 {
		f = &file[1]
	} else {
		f = &file[0]
	}
	log.Printf("error at file: %v line: %v => %v\n", *f, line, err)
}
