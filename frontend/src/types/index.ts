export interface User {
  id: number
  username: string
  email: string
  avatar?: string
  role: 'admin' | 'user' | 'guest'
  fingerprint?: string
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface Article {
  id: number
  title: string
  content: string
  summary: string
  tags: string[]
  author: User
  status: 'draft' | 'published'
  view_count: number
  like_count: number
  comment_count: number
  cover_image?: string
  created_at: string
  updated_at: string
  published_at?: string
  deleted_at?: string
}

export interface Comment {
  id: number
  content: string
  author: User
  article_id: number
  parent_id?: number
  replies?: Comment[]
  reply_count: number
  latest_reply?: Comment
  created_at: string
  updated_at: string
}

export interface CreateArticleRequest {
  title: string
  content: string
  summary: string
  tags: string[]
  status: 'draft' | 'published' | 'scheduled'
  cover_image?: string
  published_at?: string
}

export interface UpdateArticleRequest extends Partial<CreateArticleRequest> {
  id: number
}

export interface CreateCommentRequest {
  content: string
  article_id: number
  parent_id?: number
  fingerprint?: string
}

export interface ArticleListResponse {
  articles: Article[]
  total: number
  page: number
  page_size: number
}

export interface CommentListResponse {
  comments: Comment[]
  total: number
  page: number
  page_size: number
}

export interface SearchParams {
  keyword?: string
  tags?: string[]
  page?: number
  page_size?: number
  sort_by?: 'created_at' | 'view_count' | 'like_count'
  sort_order?: 'asc' | 'desc'
  include_drafts?: boolean
}