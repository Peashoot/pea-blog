package repository

import (
	"database/sql"
	"fmt"
	"pea-blog-backend/internal/model"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, username, email, password_hash, avatar, role, fingerprint, created_at, updated_at
		FROM users WHERE username = ?
	`
	row := r.db.QueryRow(query, username)
	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Avatar, &user.Role, &user.Fingerprint, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	if user.Role == "admin" {
		user.Fingerprint = nil
	}
	return user, nil
}

func (r *UserRepository) GetByID(id int) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, username, email, password_hash, avatar, role, fingerprint, created_at, updated_at
		FROM users WHERE id = ?
	`
	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Avatar, &user.Role, &user.Fingerprint, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	if user.Role == "admin" {
		user.Fingerprint = nil
	}
	return user, nil
}

func (r *UserRepository) GetByFingerprint(fingerprint string) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, username, email, password_hash, avatar, role, fingerprint, created_at, updated_at
		FROM users WHERE fingerprint = ?
	`
	row := r.db.QueryRow(query, fingerprint)
	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Avatar, &user.Role, &user.Fingerprint, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	if user.Role == "admin" {
		user.Fingerprint = nil
	}
	return user, nil
}

func (r *UserRepository) CreateGuestUser(fingerprint string) (*model.User, error) {
	// Need to add gofakeit to the project
	// go get github.com/brianvoe/gofakeit/v6
	username := gofakeit.Username()
	email := gofakeit.Email()
	user := &model.User{
		Username:    username,
		Email:       email,
		Password:    "", // No password for guest users
		Role:        "guest",
		Fingerprint: &fingerprint,
	}

	query := `
		INSERT INTO users (username, email, password_hash, role, fingerprint)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.Role, user.Fingerprint)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.ID = int(id)

	return r.GetByID(user.ID)
}

type ArticleRepository struct {
	db     *sql.DB
	dbType string
}

func NewArticleRepository(db *sql.DB) *ArticleRepository {
	// Detect database type
	dbType := "sqlite"
	if rows, err := db.Query("SELECT version()"); err == nil {
		rows.Close()
		dbType = "postgres"
	}

	return &ArticleRepository{db: db, dbType: dbType}
}

func (r *ArticleRepository) GetAll(params model.SearchParams) ([]model.Article, int, error) {
	var articles []model.Article
	var totalCount int

	// 安全检查：验证排序字段和顺序，防止SQL注入
	validSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"view_count": true,
		"like_count": true,
		"title":      true,
	}
	validSortOrders := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !validSortFields[params.SortBy] {
		params.SortBy = "created_at" // 默认安全值
	}
	if !validSortOrders[strings.ToLower(params.SortOrder)] {
		params.SortOrder = "desc" // 默认安全值
	}

	baseQuery := `
		SELECT a.id, a.title, a.content, a.summary, a.tags, a.author_id, a.status,
			   a.view_count, a.like_count, a.comment_count, a.cover_image,
			   a.created_at, a.updated_at,
			   u.id, u.username, u.email, u.avatar, u.role, u.created_at, u.updated_at
		FROM articles a
		JOIN users u ON a.author_id = u.id
	`

	countQuery := `SELECT COUNT(*) FROM articles a`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.Keyword != "" {
		if r.dbType == "postgres" {
			conditions = append(conditions, fmt.Sprintf("(a.title ILIKE $%d OR a.content ILIKE $%d OR a.summary ILIKE $%d)", argIndex, argIndex+1, argIndex+2))
		} else {
			conditions = append(conditions, "(a.title LIKE ? OR a.content LIKE ? OR a.summary LIKE ?)")
		}
		keyword := "%" + params.Keyword + "%"
		args = append(args, keyword, keyword, keyword)
		argIndex += 3
	}

	if params.Tags != "" {
		tags := strings.Split(params.Tags, ",")
		if r.dbType == "postgres" {
			conditions = append(conditions, fmt.Sprintf("a.tags && $%d", argIndex))
			args = append(args, pq.Array(tags))
		} else {
			// For SQLite, use simple string matching (simplified approach)
			for _, tag := range tags {
				if r.dbType == "postgres" {
					conditions = append(conditions, fmt.Sprintf("a.tags LIKE $%d", argIndex))
				} else {
					conditions = append(conditions, "a.tags LIKE ?")
				}
				args = append(args, "%"+tag+"%")
				argIndex++
			}
		}
		if r.dbType == "postgres" {
			argIndex++
		}
	}

	if !params.IncludeDrafts {
		conditions = append(conditions, "a.status = 'published'")
	}

	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	orderBy := fmt.Sprintf(" ORDER BY a.%s %s", params.SortBy, strings.ToUpper(params.SortOrder))
	baseQuery += orderBy

	offset := (params.Page - 1) * params.PageSize
	if r.dbType == "postgres" {
		baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	} else {
		baseQuery += " LIMIT ? OFFSET ?"
	}
	args = append(args, params.PageSize, offset)

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var article model.Article
		var author model.User
		var tagsStr string

		err := rows.Scan(
			&article.ID, &article.Title, &article.Content, &article.Summary,
			&tagsStr, &article.AuthorID, &article.Status,
			&article.ViewCount, &article.LikeCount, &article.CommentCount,
			&article.CoverImage, &article.CreatedAt, &article.UpdatedAt,
			&author.ID, &author.Username, &author.Email, &author.Avatar,
			&author.Role, &author.CreatedAt, &author.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Parse tags
		if r.dbType == "postgres" {
			if err := pq.Array(&article.Tags).Scan(tagsStr); err != nil {
				article.Tags = []string{}
			}
		} else {
			// For SQLite, parse comma-separated string
			if tagsStr != "" {
				article.Tags = strings.Split(tagsStr, ",")
			} else {
				article.Tags = []string{}
			}
		}

		article.Author = &author
		articles = append(articles, article)
	}

	return articles, totalCount, nil
}

func (r *ArticleRepository) GetByID(id int) (*model.Article, error) {
	article := &model.Article{}
	author := &model.User{}
	var tagsStr string

	query := `
		SELECT a.id, a.title, a.content, a.summary, a.tags, a.author_id, a.status,
			   a.view_count, a.like_count, a.comment_count, a.cover_image,
			   a.created_at, a.updated_at,
			   u.id, u.username, u.email, u.avatar, u.role, u.created_at, u.updated_at
		FROM articles a
		JOIN users u ON a.author_id = u.id
		WHERE a.id = ?
	`

	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&article.ID, &article.Title, &article.Content, &article.Summary,
		&tagsStr, &article.AuthorID, &article.Status,
		&article.ViewCount, &article.LikeCount, &article.CommentCount,
		&article.CoverImage, &article.CreatedAt, &article.UpdatedAt,
		&author.ID, &author.Username, &author.Email, &author.Avatar,
		&author.Role, &author.CreatedAt, &author.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("article not found")
		}
		return nil, err
	}

	// Parse tags
	if r.dbType == "postgres" {
		if err := pq.Array(&article.Tags).Scan(tagsStr); err != nil {
			article.Tags = []string{}
		}
	} else {
		// For SQLite, parse comma-separated string
		if tagsStr != "" {
			article.Tags = strings.Split(tagsStr, ",")
		} else {
			article.Tags = []string{}
		}
	}

	article.Author = author

	_, err = r.db.Exec("UPDATE articles SET view_count = view_count + 1 WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	article.ViewCount++

	return article, nil
}

func (r *ArticleRepository) Create(article *model.Article) error {
	var tagsValue interface{}
	if r.dbType == "postgres" {
		tagsValue = pq.Array(article.Tags)
	} else {
		tagsValue = strings.Join(article.Tags, ",")
	}

	query := `
		INSERT INTO articles (title, content, summary, tags, author_id, status, cover_image, published_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		article.Title, article.Content, article.Summary,
		tagsValue, article.AuthorID, article.Status, article.CoverImage, article.PublishedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	article.ID = int(id)

	return nil
}

func (r *ArticleRepository) Update(article *model.Article) error {
	var tagsValue interface{}
	if r.dbType == "postgres" {
		tagsValue = pq.Array(article.Tags)
	} else {
		tagsValue = strings.Join(article.Tags, ",")
	}

	query := `
		UPDATE articles 
		SET title = ?, content = ?, summary = ?, tags = ?, status = ?, 
		    cover_image = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		article.Title, article.Content, article.Summary,
		tagsValue, article.Status, article.CoverImage, article.ID,
	)

	return err
}

func (r *ArticleRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM articles WHERE id = ?", id)
	return err
}

func (r *ArticleRepository) Like(userID, articleID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if r.dbType == "postgres" {
		_, err = tx.Exec("INSERT INTO likes (user_id, article_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, articleID)
	} else {
		_, err = tx.Exec("INSERT OR IGNORE INTO likes (user_id, article_id) VALUES (?, ?)", userID, articleID)
	}
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE articles SET like_count = (SELECT COUNT(*) FROM likes WHERE article_id = ?) WHERE id = ?", articleID, articleID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ArticleRepository) Unlike(userID, articleID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM likes WHERE user_id = ? AND article_id = ?", userID, articleID)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE articles SET like_count = (SELECT COUNT(*) FROM likes WHERE article_id = ?) WHERE id = ?", articleID, articleID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ArticleRepository) Unpublish(id int) error {
	_, err := r.db.Exec("UPDATE articles SET status = 'draft', published_at = NULL WHERE id = ?", id)
	return err
}

func (r *ArticleRepository) GetScheduledArticles() ([]model.Article, error) {
	var articles []model.Article

	query := `
		SELECT a.id, a.title, a.content, a.summary, a.tags, a.author_id, a.status,
			   a.view_count, a.like_count, a.comment_count, a.cover_image,
			   a.created_at, a.updated_at, a.published_at, a.deleted_at,
			   u.id, u.username, u.email, u.avatar, u.role, u.created_at, u.updated_at
		FROM articles a
		JOIN users u ON a.author_id = u.id
		WHERE a.status = 'scheduled' AND a.published_at <= CURRENT_TIMESTAMP AND a.deleted_at IS NULL
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var article model.Article
		var author model.User
		var tagsStr string

		err := rows.Scan(
			&article.ID, &article.Title, &article.Content, &article.Summary,
			&tagsStr, &article.AuthorID, &article.Status,
			&article.ViewCount, &article.LikeCount, &article.CommentCount,
			&article.CoverImage, &article.CreatedAt, &article.UpdatedAt, &article.PublishedAt, &article.DeletedAt,
			&author.ID, &author.Username, &author.Email, &author.Avatar,
			&author.Role, &author.CreatedAt, &author.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse tags
		if r.dbType == "postgres" {
			if err := pq.Array(&article.Tags).Scan(tagsStr); err != nil {
				article.Tags = []string{}
			}
		} else {
			// For SQLite, parse comma-separated string
			if tagsStr != "" {
				article.Tags = strings.Split(tagsStr, ",")
			} else {
				article.Tags = []string{}
			}
		}

		article.Author = &author
		articles = append(articles, article)
	}

	return articles, nil
}

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) GetByArticleID(articleID int, page int, pageSize int) ([]model.Comment, int, error) {
	var totalCount int
	countQuery := "SELECT COUNT(*) FROM comments WHERE article_id = ? AND parent_id IS NULL AND deleted_at IS NULL"
	err := r.db.QueryRow(countQuery, articleID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT c.id, c.content, c.author_id, c.article_id, c.parent_id, c.created_at, c.updated_at,
			   u.id, u.username, u.email, u.avatar, u.role, u.fingerprint, u.created_at, u.updated_at,
			   (SELECT COUNT(*) FROM comments r WHERE r.parent_id = c.id AND r.deleted_at IS NULL) as reply_count
		FROM comments c
		JOIN users u ON c.author_id = u.id
		WHERE c.article_id = ? AND c.parent_id IS NULL AND c.deleted_at IS NULL
		ORDER BY c.created_at DESC
		LIMIT ? OFFSET ?
	`
	offset := (page - 1) * pageSize
	rows, err := r.db.Query(query, articleID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		var author model.User

		err := rows.Scan(
			&comment.ID, &comment.Content, &comment.AuthorID, &comment.ArticleID,
			&comment.ParentID, &comment.CreatedAt, &comment.UpdatedAt,
			&author.ID, &author.Username, &author.Email, &author.Avatar,
			&author.Role, &author.Fingerprint, &author.CreatedAt, &author.UpdatedAt,
			&comment.ReplyCount,
		)
		if err != nil {
			return nil, 0, err
		}

		comment.Author = &author
		comments = append(comments, comment)
	}

	return comments, totalCount, nil
}

func (r *CommentRepository) GetByID(id int) (*model.Comment, error) {
	comment := &model.Comment{}
	author := &model.User{}
	query := `
		SELECT c.id, c.content, c.author_id, c.article_id, c.parent_id, c.created_at, c.updated_at,
			   u.id, u.username, u.email, u.avatar, u.role, u.fingerprint, u.created_at, u.updated_at
		FROM comments c
		JOIN users u ON c.author_id = u.id
		WHERE c.id = ?
	`
	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&comment.ID, &comment.Content, &comment.AuthorID, &comment.ArticleID,
		&comment.ParentID, &comment.CreatedAt, &comment.UpdatedAt,
		&author.ID, &author.Username, &author.Email, &author.Avatar,
		&author.Role, &author.Fingerprint, &author.CreatedAt, &author.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, err
	}
	comment.Author = author
	return comment, nil
}

func (r *CommentRepository) GetRepliesByCommentID(commentID int, page int, pageSize int) ([]model.Comment, int, error) {
	var totalCount int
	countQuery := "SELECT COUNT(*) FROM comments WHERE parent_id = ?"
	err := r.db.QueryRow(countQuery, commentID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := `
		WITH RECURSIVE comment_tree AS (
			SELECT c.id, c.content, c.author_id, c.article_id, c.parent_id, c.created_at, c.updated_at
			FROM comments c
			WHERE c.parent_id = ?
			UNION ALL
			SELECT c.id, c.content, c.author_id, c.article_id, c.parent_id, c.created_at, c.updated_at
			FROM comments c
			JOIN comment_tree ct ON c.parent_id = ct.id
		)
		SELECT ct.id, ct.content, ct.author_id, ct.article_id, ct.parent_id, ct.created_at, ct.updated_at,
			   u.id, u.username, u.email, u.avatar, u.role, u.fingerprint, u.created_at, u.updated_at,
			   (SELECT COUNT(*) FROM comments r WHERE r.parent_id = ct.id AND r.deleted_at IS NULL) as reply_count
		FROM comment_tree ct
		JOIN users u ON ct.author_id = u.id
		ORDER BY ct.created_at DESC
	`
	rows, err := r.db.Query(query, commentID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	commentMap := make(map[int]*model.Comment)
	var allReplies []model.Comment

	for rows.Next() {
		var comment model.Comment
		var author model.User

		err := rows.Scan(
			&comment.ID, &comment.Content, &comment.AuthorID, &comment.ArticleID,
			&comment.ParentID, &comment.CreatedAt, &comment.UpdatedAt,
			&author.ID, &author.Username, &author.Email, &author.Avatar,
			&author.Role, &author.Fingerprint, &author.CreatedAt, &author.UpdatedAt,
			&comment.ReplyCount,
		)
		if err != nil {
			return nil, 0, err
		}

		comment.Author = &author
		allReplies = append(allReplies, comment)
		commentMap[comment.ID] = &allReplies[len(allReplies)-1]
	}

	var rootReplies []model.Comment
	for i := range allReplies {
		comment := &allReplies[i]
		if comment.ParentID != nil && *comment.ParentID == commentID {
			rootReplies = append(rootReplies, *comment)
		} else if comment.ParentID != nil {
			if parent, exists := commentMap[*comment.ParentID]; exists {
				parent.Replies = append(parent.Replies, *comment)
			}
		}
	}

	// Manual pagination for root replies
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(rootReplies) {
		return []model.Comment{}, totalCount, nil
	}
	if end > len(rootReplies) {
		end = len(rootReplies)
	}

	return rootReplies[start:end], totalCount, nil
}

func (r *CommentRepository) Create(comment *model.Comment) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO comments (content, author_id, article_id, parent_id)
		VALUES (?, ?, ?, ?)
	`

	result, err := tx.Exec(query, comment.Content, comment.AuthorID, comment.ArticleID, comment.ParentID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	comment.ID = int(id)

	_, err = tx.Exec("UPDATE articles SET comment_count = comment_count + 1 WHERE id = ?", comment.ArticleID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *CommentRepository) Delete(id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var articleID int
	err = tx.QueryRow("SELECT article_id FROM comments WHERE id = ?", id).Scan(&articleID)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE comments SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE articles SET comment_count = (SELECT COUNT(*) FROM comments WHERE article_id = ? AND deleted_at IS NULL) WHERE id = ?", articleID, articleID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

type Repository struct {
	User    *UserRepository
	Article *ArticleRepository
	Comment *CommentRepository
}

func New(db *sql.DB) *Repository {
	return &Repository{
		User:    NewUserRepository(db),
		Article: NewArticleRepository(db),
		Comment: NewCommentRepository(db),
	}
}
