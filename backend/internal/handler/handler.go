package handler

import (
	"pea-blog-backend/internal/model"
	"pea-blog-backend/internal/service"
	"pea-blog-backend/pkg/logger"
	"pea-blog-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
	logger      *logger.Logger
}

func NewAuthHandler(authService *service.AuthService, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	loginResponse, err := h.authService.Login(req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, loginResponse)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	response.SuccessWithMessage(c, "Logged out successfully", nil)
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	user, err := h.authService.GetCurrentUser(userID.(int))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, user)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	response.InternalServerError(c, "Not implemented")
}

type ArticleHandler struct {
	articleService *service.ArticleService
	logger         *logger.Logger
}

func NewArticleHandler(articleService *service.ArticleService, logger *logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
		logger:         logger,
	}
}

func (h *ArticleHandler) GetArticles(c *gin.Context) {
	var params model.SearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, "Invalid query parameters")
		return
	}

	params.IncludeDrafts = true
	articles, err := h.articleService.GetArticles(params)
	if err != nil {
		h.logger.Error("Failed to get articles", "error", err)
		response.InternalServerError(c, "Failed to get articles: "+err.Error())
		return
	}

	response.Success(c, articles)
}

func (h *ArticleHandler) GetPublishedArticles(c *gin.Context) {
	var params model.SearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, "Invalid query parameters")
		return
	}

	articles, err := h.articleService.GetPublishedArticles(params)
	if err != nil {
		h.logger.Error("Failed to get published articles", "error", err)
		response.InternalServerError(c, "Failed to get published articles: "+err.Error())
		return
	}

	response.Success(c, articles)
}

func (h *ArticleHandler) GetArticleByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid article ID")
		return
	}

	article, err := h.articleService.GetArticleByID(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, article)
}

func (h *ArticleHandler) SearchArticles(c *gin.Context) {
	var params model.SearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, "Invalid query parameters")
		return
	}

	if params.Keyword == "" {
		response.BadRequest(c, "Keyword is required for search")
		return
	}

	articles, err := h.articleService.GetArticles(params)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, articles)
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req model.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	article, err := h.articleService.CreateArticle(req, userID.(int))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, article)
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid article ID")
		return
	}

	var req model.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	article, err := h.articleService.UpdateArticle(id, req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, article)
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid article ID")
		return
	}

	err = h.articleService.DeleteArticle(id)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Article deleted successfully", nil)
}

func (h *ArticleHandler) LikeArticle(c *gin.Context) {
	idStr := c.Param("id")
	articleID, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid article ID")
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		userID := 0
		err = h.articleService.LikeArticle(userID, articleID)
	} else {
		userID := userIDValue.(int)
		err = h.articleService.LikeArticle(userID, articleID)
	}

	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Article liked successfully", nil)
}

func (h *ArticleHandler) UnlikeArticle(c *gin.Context) {
	idStr := c.Param("id")
	articleID, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid article ID")
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		userID := 0
		err = h.articleService.UnlikeArticle(userID, articleID)
	} else {
		userID := userIDValue.(int)
		err = h.articleService.UnlikeArticle(userID, articleID)
	}

	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Article unliked successfully", nil)
}

func (h *ArticleHandler) UnpublishArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid article ID")
		return
	}

	err = h.articleService.UnpublishArticle(id)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Article unpublished successfully", nil)
}

type CommentHandler struct {
	commentService *service.CommentService
	logger         *logger.Logger
}

func NewCommentHandler(commentService *service.CommentService, logger *logger.Logger) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		logger:         logger,
	}
}

func (h *CommentHandler) GetCommentsByArticleID(c *gin.Context) {
	idStr := c.Param("id")
	articleID, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid article ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "15"))

	comments, err := h.commentService.GetCommentsByArticleID(articleID, page, pageSize)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, comments)
}

func (h *CommentHandler) GetRepliesByCommentID(c *gin.Context) {
	idStr := c.Param("id")
	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid comment ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "15"))

	replies, err := h.commentService.GetRepliesByCommentID(commentID, page, pageSize)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, replies)
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req model.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	var authorID int
	var fingerprint *string

	userIDValue, exists := c.Get("userID")
	if exists {
		authorID = userIDValue.(int)
	} else {
		fingerprint = req.Fingerprint
	}

	comment, err := h.commentService.CreateComment(req, authorID, fingerprint)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, comment)
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid comment ID")
		return
	}

	var authorID int
	var fingerprint *string

	userIDValue, exists := c.Get("userID")
	if exists {
		authorID = userIDValue.(int)
	} else {
		var req struct {
			Fingerprint string `json:"fingerprint"`
		}
		if err := c.ShouldBindJSON(&req); err == nil {
			fingerprint = &req.Fingerprint
		}
	}

	role, _ := c.Get("role")
	isAdmin := role == "admin"

	err = h.commentService.DeleteComment(commentID, authorID, isAdmin, fingerprint)
	if err != nil {
		if err.Error() == "user not authorized to delete this comment" {
			response.Forbidden(c, err.Error())
		} else {
			response.InternalServerError(c, err.Error())
		}
		return
	}

	response.SuccessWithMessage(c, "Comment deleted successfully", nil)
}

type Handler struct {
	Auth    *AuthHandler
	Article *ArticleHandler
	Comment *CommentHandler
	System  *SystemHandler
	Image   *ImageHandler
}

func New(services *service.Service, logger *logger.Logger) *Handler {
	return &Handler{
		Auth:    NewAuthHandler(services.Auth, logger),
		Article: NewArticleHandler(services.Article, logger),
		Comment: NewCommentHandler(services.Comment, logger),
		System:  nil, // 在main.go中单独设置
	}
}
