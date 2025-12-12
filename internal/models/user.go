package models

import (
    "errors"
    "regexp"
    "strings"
)

type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func (u *User) ValidateForCreate() error {
    if strings.TrimSpace(u.Name) == "" {
        return errors.New("name is required")
    }
    if strings.TrimSpace(u.Email) == "" {
        return errors.New("email is required")
    }
    if !emailRegex.MatchString(u.Email) {
        return errors.New("invalid email format")
    }
    if len(strings.TrimSpace(u.Password)) < 6 {
        return errors.New("password must be at least 6 characters")
    }
    return nil
}

func (u *User) ValidateForUpdate() error {
    if u.Name != "" && strings.TrimSpace(u.Name) == "" {
        return errors.New("invalid name")
    }

    if u.Email != "" {
        if !emailRegex.MatchString(u.Email) {
            return errors.New("invalid email format")
        }
    }

    if u.Password != "" && len(u.Password) < 6 {
        return errors.New("password must be at least 6 characters")
    }

    return nil
}
