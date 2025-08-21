package service

import (
	"fmt"
	"pea-blog-backend/internal/model"
	"pea-blog-backend/internal/repository"
	"pea-blog-backend/internal/util"
	"pea-blog-backend/pkg/logger"
)

type AuthService struct {
	userRepo *repository.UserRepository
	logger   *logger.Logger
}

func NewAuthService(userRepo *repository.UserRepository, logger *logger.Logger) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *AuthService) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		s.logger.Error("Failed to get user by username", "username", req.Username, "error", err)
		return nil, fmt.Errorf("invalid credentials")
	}

	if !util.CheckPassword(req.Password, user.Password) {
		s.logger.Warn("Invalid password attempt", "username", req.Username)
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := util.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		s.logger.Error("Failed to generate JWT", "userID", user.ID, "error", err)
		return nil, fmt.Errorf("failed to generate token")
	}

	user.Password = ""

	return &model.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *AuthService) GetCurrentUser(userID int) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		s.logger.Error("Failed to get user by ID", "userID", userID, "error", err)
		return nil, fmt.Errorf("user not found")
	}

	user.Password = ""
	return user, nil
}

type ArticleService struct {
	articleRepo *repository.ArticleRepository
	userRepo    *repository.UserRepository
	logger      *logger.Logger
}

func NewArticleService(articleRepo *repository.ArticleRepository, userRepo *repository.UserRepository, logger *logger.Logger) *ArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
		userRepo:    userRepo,
		logger:      logger,
	}
}

func (s *ArticleService) GetArticles(params model.SearchParams) (*model.ArticleListResponse, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 || params.PageSize > 100 {
		params.PageSize = 10
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	params.IncludeDrafts = true
	articles, total, err := s.articleRepo.GetAll(params)
	if err != nil {
		s.logger.Error("Failed to get articles from repository", "error", err, "params", params)
		return nil, fmt.Errorf("failed to get articles from repository: %w", err)
	}

	return &model.ArticleListResponse{
		Articles: articles,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

func (s *ArticleService) GetPublishedArticles(params model.SearchParams) (*model.ArticleListResponse, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 || params.PageSize > 100 {
		params.PageSize = 10
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	articles, total, err := s.articleRepo.GetAll(params)
	if err != nil {
		s.logger.Error("Failed to get published articles from repository", "error", err, "params", params)
		return nil, fmt.Errorf("failed to get published articles from repository: %w", err)
	}

	return &model.ArticleListResponse{
		Articles: articles,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

func (s *ArticleService) GetArticleByID(id int) (*model.Article, error) {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get article by ID", "articleID", id, "error", err)
		return nil, fmt.Errorf("article not found")
	}

	return article, nil
}

func (s *ArticleService) GetArticleByTitle(title string) (*model.Article, error) {
	article, err := s.articleRepo.GetByTitle(title)
	if err != nil {
		s.logger.Error("Failed to get article by title", "title", title, "error", err)
		return nil, fmt.Errorf("article not found")
	}

	return article, nil
}

func (s *ArticleService) CreateArticle(req model.CreateArticleRequest, authorID int) (*model.Article, error) {
	article := &model.Article{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Tags:       req.Tags,
		AuthorID:   authorID,
		Status:     req.Status,
		CoverImage: req.CoverImage,
		PublishedAt: req.PublishedAt,
	}

	err := s.articleRepo.Create(article)
	if err != nil {
		s.logger.Error("Failed to create article", "authorID", authorID, "error", err)
		return nil, fmt.Errorf("failed to create article")
	}

	author, err := s.userRepo.GetByID(authorID)
	if err == nil {
		author.Password = ""
		article.Author = author
	}

	s.logger.Info("Article created", "articleID", article.ID, "authorID", authorID)
	return article, nil
}

func (s *ArticleService) UpdateArticle(id int, req model.UpdateArticleRequest) (*model.Article, error) {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("article not found")
	}

	if req.Title != nil {
		article.Title = *req.Title
	}
	if req.Content != nil {
		article.Content = *req.Content
	}
	if req.Summary != nil {
		article.Summary = *req.Summary
	}
	if req.Tags != nil {
		article.Tags = req.Tags
	}
	if req.Status != nil {
		article.Status = *req.Status
	}
	if req.CoverImage != nil {
		article.CoverImage = req.CoverImage
	}
	if req.PublishedAt != nil {
		article.PublishedAt = req.PublishedAt
	}

	err = s.articleRepo.Update(article)
	if err != nil {
		s.logger.Error("Failed to update article", "articleID", id, "error", err)
		return nil, fmt.Errorf("failed to update article")
	}

	s.logger.Info("Article updated", "articleID", id)
	return article, nil
}

func (s *ArticleService) DeleteArticle(id int) error {
	err := s.articleRepo.Delete(id)
	if err != nil {
		s.logger.Error("Failed to delete article", "articleID", id, "error", err)
		return fmt.Errorf("failed to delete article")
	}

	s.logger.Info("Article deleted", "articleID", id)
	return nil
}

func (s *ArticleService) UnpublishArticle(id int) error {
	err := s.articleRepo.Unpublish(id)
	if err != nil {
		s.logger.Error("Failed to unpublish article", "articleID", id, "error", err)
		return fmt.Errorf("failed to unpublish article")
	}

	s.logger.Info("Article unpublished", "articleID", id)
	return nil
}

func (s *ArticleService) LikeArticle(userID, articleID int) error {
	err := s.articleRepo.Like(userID, articleID)
	if err != nil {
		s.logger.Error("Failed to like article", "userID", userID, "articleID", articleID, "error", err)
		return fmt.Errorf("failed to like article")
	}

	return nil
}

func (s *ArticleService) UnlikeArticle(userID, articleID int) error {
	err := s.articleRepo.Unlike(userID, articleID)
	if err != nil {
		s.logger.Error("Failed to unlike article", "userID", userID, "articleID", articleID, "error", err)
		return fmt.Errorf("failed to unlike article")
	}

	return nil
}

func (s *ArticleService) PublishScheduledArticles() []error {
	articles, err := s.articleRepo.GetScheduledArticles()
	if err != nil {
		return []error{err}
	}

	now := time.Now()
	var errors []error
	for _, article := range articles {
		article.Status = "published"
		article.PublishedAt = &now
		err := s.articleRepo.Update(&article)
		if err != nil {
			errors = append(errors, err)
		}
		s.logger.Info("Article published from schedule", "articleID", article.ID, "scheduledTime", article.PublishedAt, "publishedAt", now)
	}

	return errors
}

func (s *ArticleService) GetAllArticlesForExport() ([]model.Article, error) {
	articles, _, err := s.articleRepo.GetAll(model.SearchParams{IncludeDrafts: true, PageSize: 9999})
	return articles, err
}

func (s *ArticleService) ImportArticles(articles []model.Article) error {
	for _, article := range articles {
		if err := s.articleRepo.Create(&article); err != nil {
			return err
		}
	}
	return nil
}

type CommentService struct {
	commentRepo *repository.CommentRepository
	userRepo    *repository.UserRepository
	logger      *logger.Logger
}

func NewCommentService(commentRepo *repository.CommentRepository, userRepo *repository.UserRepository, logger *logger.Logger) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		userRepo:    userRepo,
		logger:      logger,
	}
}

func (s *CommentService) GetCommentsByArticleID(articleID int, page int, pageSize int) (*model.CommentListResponse, error) {
	comments, total, err := s.commentRepo.GetByArticleID(articleID, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to get comments", "articleID", articleID, "error", err)
		return nil, fmt.Errorf("failed to get comments")
	}

	return &model.CommentListResponse{
		Comments: comments,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *CommentService) GetRepliesByCommentID(commentID int, page int, pageSize int) (*model.CommentListResponse, error) {
	replies, total, err := s.commentRepo.GetRepliesByCommentID(commentID, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to get replies", "commentID", commentID, "error", err)
		return nil, fmt.Errorf("failed to get replies")
	}

	return &model.CommentListResponse{
		Comments: replies,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *CommentService) CreateComment(req model.CreateCommentRequest, authorID int, fingerprint *string) (*model.Comment, error) {
	if authorID == 0 && fingerprint != nil {
		// Handle guest user
		user, err := s.userRepo.GetByFingerprint(*fingerprint)
		if err != nil {
			// Create new guest user if not found
			user, err = s.userRepo.CreateGuestUser(*fingerprint)
			if err != nil {
				s.logger.Error("Failed to create guest user", "fingerprint", *fingerprint, "error", err)
				return nil, fmt.Errorf("failed to create guest user")
			}
		}
		authorID = user.ID
	} else if authorID == 0 && fingerprint == nil {
		return nil, fmt.Errorf("user not authenticated and no fingerprint provided")
	}

	comment := &model.Comment{
		Content:   req.Content,
		AuthorID:  authorID,
		ArticleID: req.ArticleID,
		ParentID:  req.ParentID,
	}

	err := s.commentRepo.Create(comment)
	if err != nil {
		s.logger.Error("Failed to create comment", "authorID", authorID, "articleID", req.ArticleID, "error", err)
		return nil, fmt.Errorf("failed to create comment")
	}

	newComment, err := s.commentRepo.GetByID(comment.ID)
	if err != nil {
		s.logger.Error("Failed to get new comment by ID", "commentID", comment.ID, "error", err)
		return nil, fmt.Errorf("failed to get new comment by ID")
	}

	s.logger.Info("Comment created", "commentID", newComment.ID, "authorID", authorID, "articleID", req.ArticleID)
	return newComment, nil
}

func (s *CommentService) DeleteComment(id, userID int, isAdmin bool, fingerprint *string) error {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get comment by ID", "commentID", id, "error", err)
		return fmt.Errorf("comment not found")
	}

	if isAdmin {
		return s.commentRepo.Delete(id)
	}

	if userID != 0 {
		if comment.AuthorID == userID {
			return s.commentRepo.Delete(id)
		}
	} else if fingerprint != nil {
		author, err := s.userRepo.GetByID(comment.AuthorID)
		if err != nil {
			return fmt.Errorf("author not found")
		}
		if author.Fingerprint != nil && *author.Fingerprint == *fingerprint {
			return s.commentRepo.Delete(id)
		}
	}

	s.logger.Warn("User not authorized to delete comment", "commentID", id, "userID", userID)
	return fmt.Errorf("user not authorized to delete this comment")
}

type Service struct {
	Auth    *AuthService
	Article *ArticleService
	Comment *CommentService
}

func New(repos *repository.Repository, logger *logger.Logger) *Service {
	return &Service{
		Auth:    NewAuthService(repos.User, logger),
		Article: NewArticleService(repos.Article, repos.User, logger),
		Comment: NewCommentService(repos.Comment, repos.User, logger),
	}
}