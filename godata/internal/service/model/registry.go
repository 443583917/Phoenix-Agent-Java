package model

import (
	"context"
	"fmt"
	"sync"

	tmodel "trpc.group/trpc-go/trpc-agent-go/model"
)

type Registry struct {
	mu     sync.RWMutex
	models map[string]tmodel.Model
	active map[string]string
}

func NewRegistry() *Registry {
	return &Registry{
		models: make(map[string]tmodel.Model),
		active: make(map[string]string),
	}
}

func (r *Registry) Register(id string, m tmodel.Model) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.models[id] = m
}

func (r *Registry) GetActive(modelType string) (tmodel.Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.active[modelType]
	if !ok {
		return nil, fmt.Errorf("no active model for type %s", modelType)
	}
	m, ok := r.models[id]
	if !ok {
		return nil, fmt.Errorf("model %s not found", id)
	}
	return m, nil
}

func (r *Registry) SetActive(modelType, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.models[id]; !ok {
		return fmt.Errorf("model %s not registered", id)
	}
	r.active[modelType] = id
	return nil
}

func (r *Registry) List() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make(map[string]string)
	for k, v := range r.active {
		result[k] = v
	}
	return result
}

type Proxy struct {
	registry  *Registry
	modelType string
}

func NewProxy(registry *Registry, modelType string) *Proxy {
	return &Proxy{registry: registry, modelType: modelType}
}

func (p *Proxy) GenerateContent(ctx context.Context, req *tmodel.Request) (<-chan *tmodel.Response, error) {
	m, err := p.registry.GetActive(p.modelType)
	if err != nil {
		return nil, err
	}
	return m.GenerateContent(ctx, req)
}

func (p *Proxy) Info() tmodel.Info {
	m, err := p.registry.GetActive(p.modelType)
	if err != nil {
		return tmodel.Info{}
	}
	return m.Info()
}

var _ tmodel.Model = (*Proxy)(nil)
