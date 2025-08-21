<template>
  <div class="article-detail-page" ref="pageTop">
    <Navbar />
    
    <main class="main-content">
      <div class="container">
        <div v-if="isLoading" class="loading-container">
          <div class="loading-spinner"></div>
          <p>{{ $t('common.loading') }}</p>
        </div>

        <article v-else-if="article" class="article-detail glass-effect">
          <header class="article-header">
            <div class="article-meta">
              <div class="author-info">
                <el-avatar :size="40" :src="article.author?.avatar">
                  {{ article.author?.username?.charAt(0).toUpperCase() }}
                </el-avatar>
                <div class="author-details">
                  <span class="author-name">{{ article.author?.username }}</span>
                  <span class="publish-date">{{ $t('article_detail.publish_date') }}: {{ formatDate(article.created_at) }}</span>
                </div>
              </div>
              
              <div class="article-stats">
                <div class="stat">
                  <el-icon><View /></el-icon>
                  <span>{{ article.view_count }} {{ $t('article_detail.views') }}</span>
                </div>
                <div class="stat like-stat" @click="toggleLike">
                  <el-icon :class="{ liked: isLiked }">
                    <Star v-if="isLiked" />
                    <StarFilled v-else />
                  </el-icon>
                  <span>{{ article.like_count }} {{ $t('article_detail.likes') }}</span>
                </div>
              </div>
            </div>

            <h1 class="article-title">{{ article.title }}</h1>
            
            <div class="article-tags" v-if="article.tags.length">
              <span v-for="tag in article.tags" :key="tag" class="tag">
                #{{ tag }}
              </span>
            </div>
          </header>

          <div v-if="article.cover_image" class="article-cover">
            <img :src="article.cover_image" :alt="article.title" />
          </div>

          <div class="article-content" v-html="formattedContent"></div>
        </article>

        <div v-else class="error-state">
          <el-icon size="48" color="var(--text-secondary)"><DocumentDelete /></el-icon>
          <p>{{ $t('article_detail_page.not_found') }}</p>
          <router-link to="/" class="tech-button">{{ $t('article_detail_page.back_to_home') }}</router-link>
        </div>

        <!-- {{ $t('common.comments_section') }} -->
        <div v-if="article" class="comments-section">
          <div class="comments-header">
            <h3>{{ $t('article_detail.comments') }} ({{ article.comment_count }})</h3>
          </div>

          <div class="comment-form glass-effect">
            <div 
              class="comment-input"
              contenteditable="true"
              :placeholder="$t('article_detail.comment_placeholder')"
              @input="handleCommentInput"
              ref="commentInput"
            >
            </div>
            <div class="comment-actions">
              <button v-if="replyTo" class="tech-button secondary" @click="cancelReply">
                {{ $t('article_detail.cancel_reply') }}
              </button>
              <button 
                class="tech-button"
                @click="submitComment"
                :disabled="!newComment.trim() || isSubmittingComment"
              >
                {{ isSubmittingComment ? $t('article_detail_page.publishing') : $t('article_detail.submit_comment') }}
              </button>
            </div>
          </div>

          <div class="comments-list" ref="commentsList">
            <div v-if="isLoadingComments && comments.length === 0" class="loading-container">
              <div class="loading-spinner"></div>
              <p>{{ $t('article_detail_page.loading_comments') }}</p>
            </div>
            
            <div v-else-if="comments.length > 0">
              <div v-for="(page, pageIndex) in comments" :key="pageIndex" class="comment-page">
                <CommentItem 
                  v-for="comment in page" 
                  :key="comment.id" 
                  :comment="comment"
                  @reply="handleReply"
                  @delete="handleDeleteComment"
                />
              </div>
            </div>
            
            <div v-else class="empty-comments">
              <p>{{ $t('article_detail.no_comments') }}</p>
            </div>

            <div v-if="hasMoreComments" class="load-more-comments">
              <button class="tech-button" @click="loadComments(true)" :disabled="isLoadingComments">
                {{ isLoadingComments ? $t('article_detail_page.loading_comments') : $t('article_detail.load_more_comments') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useArticleStore, useAuthStore } from '@/stores'
import { commentApi } from '@/api'
import { formatDate } from '@/utils'
import Navbar from '@/components/Navbar.vue'
import CommentItem from '@/components/CommentItem.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Comment, User } from '@/types'
import { marked } from 'marked'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute()
const articleStore = useArticleStore()
const authStore = useAuthStore()

const articleTitle = computed(() => decodeURIComponent(route.params.title as string))
const article = computed(() => articleStore.currentArticle)
const isLoading = ref(false)
const isLoadingComments = ref(false)
const isSubmittingComment = ref(false)
const isLiked = ref(false)
const comments = ref<Comment[][]>([])
const newComment = ref('')
const currentPage = ref(1)
const totalComments = ref(0)
const hasMoreComments = computed(() => {
  const flattenedComments = comments.value.flat()
  return flattenedComments.length < totalComments.value
})
const replyTo = ref<Comment | null>(null)
const commentForm = ref<HTMLElement | null>(null)
const commentInput = ref<any>(null)

const formattedContent = computed(() => {
  if (!article.value?.content) return ''
  return marked(article.value.content)
})

const toggleLike = async () => {
  if (!article.value) return
  
  try {
    if (isLiked.value) {
      await articleStore.unlikeArticle(article.value.id)
      isLiked.value = false
    } else {
      await articleStore.likeArticle(article.value.id)
      isLiked.value = true
    }
  } catch (error) {
    ElMessage.error(t('common.operation_failed'))
  }
}

const loadComments = async (loadMore = false) => {
  if (!article.value?.id || isLoadingComments.value) return

  try {
    isLoadingComments.value = true
    if (!loadMore) {
      currentPage.value = 1
      comments.value = []
    }

    const response = await commentApi.getCommentsByArticleId(
      article.value.id,
      currentPage.value,
      15
    )

    if (response.comments) {
      comments.value.push(response.comments)
    }
    totalComments.value = response.total
    currentPage.value++
  } catch (error) {
    ElMessage.error(t('common.load_comments_failed'))
  } finally {
    isLoadingComments.value = false
  }
}

import { getFingerprint } from '@/utils/fingerprint'

const submitComment = async () => {
  if (!newComment.value.trim() || !article.value?.id) return

  try {
    isSubmittingComment.value = true

    let fingerprint: string | undefined
    if (!authStore.isLoggedIn) {
      fingerprint = await getFingerprint()
    }

    const comment = await commentApi.createComment({
      content: newComment.value.trim(),
      article_id: article.value.id,
      parent_id: replyTo.value?.id,
      fingerprint,
    })

    if (authStore.isLoggedIn) {
      comment.author = authStore.user as User
    }

    if (replyTo.value) {
      // Find the page and the comment to reply to
      for (const page of comments.value) {
        const parentComment = page.find(c => c.id === replyTo.value?.id)
        if (parentComment) {
          if (!parentComment.replies) {
            parentComment.replies = []
          }
          parentComment.replies.unshift(comment)
          break
        }
      }
    } else {
      if (comments.value.length === 0) {
        comments.value.push([])
      }
      comments.value[0].unshift(comment)
    }

    newComment.value = ''
    replyTo.value = null
    totalComments.value++
    if (article.value) {
      article.value.comment_count++
    }
    ElMessage.success(t('common.comment_submit_success'))
  } catch (error) {
    ElMessage.error(t('common.comment_submit_failed'))
  } finally {
    isSubmittingComment.value = false
  }
}

const handleReply = (parentComment: Comment) => {
  replyTo.value = parentComment
  nextTick(() => {
    if (commentInput.value) {
      commentInput.value.innerHTML = ''
      const tag = document.createElement('span')
      tag.className = 'el-tag el-tag--small reply-tag'
      tag.innerHTML = `@${parentComment.author?.username}`
      tag.contentEditable = 'false'
      commentInput.value.appendChild(tag)

      // Add a space after the tag
      const space = document.createTextNode(' ');
      commentInput.value.appendChild(space);

      commentInput.value.focus()

      // Move cursor to the end
      const range = document.createRange()
      const sel = window.getSelection()
      range.selectNodeContents(commentInput.value)
      range.collapse(false)
      sel?.removeAllRanges()
      sel?.addRange(range)

      commentForm.value?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    }
  })
}

const handleCommentInput = (event: Event) => {
  const target = event.target as HTMLElement
  newComment.value = target.innerText
}

const cancelReply = () => {
  replyTo.value = null
  newComment.value = ''
}

import { getStoredFingerprint } from '@/utils/fingerprint'

const handleDeleteComment = async (commentId: number) => {
  try {
    await ElMessageBox.confirm(
      t('article_detail.delete_comment_confirm_text'),
      t('article_detail.delete_comment_confirm_title'),
      {
        confirmButtonText: t('article_detail.confirm'),
        cancelButtonText: t('article_detail.cancel'),
        type: 'warning',
      }
    )

    let fingerprint: string | null = null
    if (!authStore.isLoggedIn) {
      fingerprint = getStoredFingerprint()
    }

    await commentApi.deleteComment(commentId, fingerprint ?? undefined)

    // Find and remove the comment from the nested arrays
    for (let i = 0; i < comments.value.length; i++) {
      const page = comments.value[i]
      const commentIndex = page.findIndex(c => c.id === commentId)
      if (commentIndex !== -1) {
        page.splice(commentIndex, 1)
        break
      }
    }

    totalComments.value--
    ElMessage.success(t('article_detail.delete_comment_success'))
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('article_detail.delete_comment_fail'))
    }
  }
}

const pageTop = ref<HTMLElement | null>(null)

onMounted(async () => {
  if (pageTop.value) {
    pageTop.value.scrollIntoView()
  }

  try {
    isLoading.value = true
    // First fetch the article
    await articleStore.fetchArticleByTitle(articleTitle.value)
    // Then load comments
    await loadComments()
  } catch (error) {
    ElMessage.error(t('common.load_articles_failed'))
  } finally {
    isLoading.value = false
  }
})
</script>

<style scoped>
.article-detail-page {
  min-height: 100vh;
}

.main-content {
  padding-top: 80px;
  padding-bottom: 2rem;
}

.loading-container, .error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 0;
  gap: 1rem;
  color: var(--text-secondary);
}

.article-detail {
  padding: 2rem;
  margin-bottom: 2rem;
}

.article-header {
  margin-bottom: 2rem;
}

.article-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.author-details {
  display: flex;
  flex-direction: column;
}

.author-name {
  font-weight: 600;
  color: var(--text-primary);
}

.publish-date {
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.article-stats {
  display: flex;
  gap: 1rem;
}

.stat {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.like-stat {
  cursor: pointer;
  transition: color 0.3s ease;
}

.like-stat:hover {
  color: var(--primary-color);
}

.like-stat .el-icon.liked {
  color: #f56565;
}

.article-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin: 0 0 1rem 0;
  color: var(--text-primary);
  line-height: 1.2;
}

.article-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}

.tag {
  background: rgba(255, 255, 255, 0.1);
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-size: 0.875rem;
  color: var(--primary-color);
  border: 1px solid var(--border-color);
}

.article-cover {
  width: 100%;
  max-height: 400px;
  overflow: hidden;
  border-radius: 12px;
  margin-bottom: 2rem;
}

.article-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.article-content {
  font-size: 1.1rem;
  line-height: 1.8;
  color: var(--text-primary);
}

.article-content :deep(h1),
.article-content :deep(h2),
.article-content :deep(h3),
.article-content :deep(h4),
.article-content :deep(h5),
.article-content :deep(h6) {
  margin: 2rem 0 1rem 0;
  color: var(--text-primary);
}

.article-content :deep(p) {
  margin: 1rem 0;
}

.article-content :deep(pre) {
  background: rgba(0, 0, 0, 0.5);
  padding: 1rem;
  border-radius: 8px;
  overflow-x: auto;
  margin: 1rem 0;
}

.article-content :deep(code) {
  background: rgba(255, 255, 255, 0.1);
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
}

.article-content :deep(pre code) {
  background: none;
  padding: 0;
}

.comments-section {
  margin-top: 3rem;
}

.comments-header {
  margin-bottom: 1.5rem;
}

.comments-header h3 {
  font-size: 1.5rem;
  margin: 0;
  color: var(--text-primary);
}

.comment-form {
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.comment-input {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 0.5rem 1rem;
  min-height: 120px;
  background-color: white;
  color: black;
}

.comment-input :deep(.reply-tag) {
  background-color: #87CEEB; /* sky blue */
  color: white;
}

.comment-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 1rem;
  gap: 0.5rem;
}

.comments-list {
  space-y: 1rem;
}

.empty-comments {
  text-align: center;
  padding: 2rem;
  color: var(--text-secondary);
}

@media (max-width: 768px) {
  .article-detail {
    padding: 1.5rem;
  }
  
  .article-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .article-title {
    font-size: 1.75rem;
  }
  
  .article-content {
    font-size: 1rem;
  }
  
  .comment-form {
    padding: 1rem;
  }
}
</style>