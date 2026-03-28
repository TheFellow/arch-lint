package controllers

import "example/infrastructure" // want `\[controllers without infrastructure\] forbidden import of "infrastructure"`

var _ = infrastructure.DB
