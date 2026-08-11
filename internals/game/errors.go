package game

import "errors"

var (
	AnimationIsOver     = errors.New("animation is over")
	AssetNotFound       = errors.New("asset not found!")
	TypeMismatchErr     = errors.New("expected type doesn't match the asset type ")
	ErrCharFilesMissing = errors.New("this charecter does not exist in this gameset")
)
