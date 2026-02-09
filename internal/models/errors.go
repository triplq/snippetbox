package models

import "errors"

var ErrNoRecord = errors.New("models: no matching record found")

var ErrDuplicateEmail = errors.New("models: this email is used")

var ErrInvalidCredentials = errors.New("models: invalid credentials")
