package sources

import "context"

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

func (s *Service) ByCode(ctx context.Context, code string) (*Source, error) {
	return s.repo.ByCode(ctx, code)
}

func (s *Service) Enabled(ctx context.Context) ([]Source, error) { return s.repo.Enabled(ctx) }
