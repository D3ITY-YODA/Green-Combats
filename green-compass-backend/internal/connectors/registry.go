package connectors

import "fmt"

type Registry struct{ byCode map[string]Connector }

func NewRegistry(connectors ...Connector) (*Registry, error) {
	registry := &Registry{byCode: make(map[string]Connector, len(connectors))}
	for _, connector := range connectors {
		if connector == nil {
			return nil, fmt.Errorf("connector must not be nil")
		}
		code := connector.Code()
		if code == "" {
			return nil, fmt.Errorf("connector code must not be empty")
		}
		if _, exists := registry.byCode[code]; exists {
			return nil, fmt.Errorf("duplicate connector code %q", code)
		}
		registry.byCode[code] = connector
	}
	return registry, nil
}

func (r *Registry) ByCode(code string) (Connector, bool) {
	connector, ok := r.byCode[code]
	return connector, ok
}
