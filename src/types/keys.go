package types

type CtxKey string

const (
	CtxApiVersionKey CtxKey = "ApiVersionKey"
	CtxResourcesKey  CtxKey = "ResourcesKey"
)

// Values returns all known values for CtxKey.
func (CtxKey) Values() []CtxKey {
	return []CtxKey{
		"ApiVersionKey",
		"ResourcesKey",
	}
}
