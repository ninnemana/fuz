package hosts

import "context"

type Service interface {
	GetLine(context.Context, GetParams) (*Record, error)
	Set(context.Context, Record) (*Record, error)
	List(context.Context, *ListParams) ([]Record, error)
	Delete(context.Context, GetParams) error
}
