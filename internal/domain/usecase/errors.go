package usecase

import "errors"

// ErrInvalidOrderInput indicates ProcessOrderInput failed validation before any
// backend call. Errors wrapping this value may include detail after ": ".
var ErrInvalidOrderInput = errors.New("invalid process order input")
