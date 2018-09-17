package main

import (
	"github.com/satori/go.uuid"
	"github.com/rs/xid"
)

func UUID() string {
	return uuid.NewV4().String()
}

// 只有20位
func UniqueId20() string {
	return xid.New().String()
}