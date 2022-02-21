// +build !deps

package component

import "context"

type Manager interface {
	Install(ctx *context.Context) error
}
