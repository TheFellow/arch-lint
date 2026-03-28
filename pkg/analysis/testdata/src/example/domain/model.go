package domain

import "example/usecase" // want `\[clean architecture\] forbidden import of "usecase"`

var _ = usecase.Service
