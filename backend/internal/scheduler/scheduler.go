package scheduler

import (
	"pea-blog-backend/internal/service"
	"pea-blog-backend/pkg/logger"
	"time"
)

type Scheduler struct {
	articleService *service.ArticleService
	logger         *logger.Logger
}

func New(articleService *service.ArticleService, logger *logger.Logger) *Scheduler {
	return &Scheduler{
		articleService: articleService,
		logger:         logger,
	}
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.logger.Debug("Scheduler checking for articles to publish...")
		errors := s.articleService.PublishScheduledArticles()
		if len(errors) > 0 {
			for _, err := range errors {
				s.logger.Error("Failed to publish scheduled article", "error", err)
			}
		}
	}
}
