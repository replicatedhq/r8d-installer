//go:build r8d
// +build r8d

package component

import "context"

type Manager interface {
	Install(ctx *context.Context) error
}
