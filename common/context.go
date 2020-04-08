package common

type ModuleContext struct {
	Asset string
}

func NewModuleContext(asset string) (*ModuleContext, error) {
	ctx := ModuleContext{
		Asset: asset,
	}

	return &ctx, nil
}
