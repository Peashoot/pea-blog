<template>
  <div class="comment-item glass-effect">
    <div class="comment-header">
      <div class="comment-author">
        <el-avatar :size="32" :src="comment.author?.avatar">
          {{ comment.author?.username?.charAt(0).toUpperCase() }}
        </el-avatar>
        <div class="author-info">
          <span class="author-name">{{ comment.author?.username }}</span>
          <span class="comment-date">{{ formatDate(comment.created_at) }}</span>
        </div>
      </div>
      
      <div class="comment-actions">
        <button class="action-btn" @click="$emit('reply', comment)">
          <el-icon><ChatRound /></el-icon>
          回复
        </button>
        <button 
          v-if="canDelete"
          class="action-btn delete-btn" 
          @click="$emit('delete', comment.id)"
        >
          <el-icon><Delete /></el-icon>
          删除
        </button>
      </div>
    </div>
    
    <div class="comment-content">
      {{ comment.content }}
    </div>

    <div v-if="comment.latest_reply" class="latest-reply">
      <CommentItem :comment="comment.latest_reply" :is-reply="true" />
    </div>

    <div class="comment-footer">
      <button v-if="comment.reply_count > 0 && !showReplies" class="action-btn" @click="toggleReplies">
        查看 {{ comment.reply_count }} 条回复
      </button>
      <button v-if="showReplies" class="action-btn" @click="toggleReplies">
        收起回复
      </button>
    </div>
    
    <div v-if="showReplies && comment.replies && comment.replies.length > 0" class="comment-replies">
      <CommentItem 
        v-for="reply in comment.replies" 
        :key="reply.id" 
        :comment="reply"
        :is-reply="true"
        @reply="$emit('reply', $event)"
        @delete="$emit('delete', $event)"
      />
      <button v-if="hasMoreReplies" class="action-btn load-more-replies" @click="loadReplies(true)">
        加载更多回复
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useAuthStore } from '@/stores'
import { formatDate } from '@/utils'
import type { Comment } from '@/types'
import { getStoredFingerprint } from '@/utils/fingerprint'
import { commentApi } from '@/api'

const props = defineProps<{
  comment: Comment
  isReply?: boolean
}>()

const emit = defineEmits<{
  reply: [comment: Comment]
  delete: [commentId: number]
}>()

const authStore = useAuthStore()
const showReplies = ref(false)
const repliesPage = ref(1)
const totalReplies = ref(0)
const hasMoreReplies = computed(() => (props.comment.replies?.length || 0) < totalReplies.value)

const canDelete = computed(() => {
  if (authStore.isLoggedIn) {
    return authStore.isAdmin || authStore.user?.id === props.comment.author?.id
  }

  const fingerprint = getStoredFingerprint()
  return fingerprint && fingerprint === props.comment.author?.fingerprint
})

const loadReplies = async (loadMore = false) => {
  if (!props.comment.id) return

  try {
    if (!loadMore) {
      repliesPage.value = 1
      props.comment.replies = []
    }

    const response = await commentApi.getRepliesByCommentId(
      props.comment.id,
      repliesPage.value,
      5 // Load 5 replies at a time
    )

    if (!props.comment.replies) {
      props.comment.replies = []
    }
    props.comment.replies.push(...(response.comments || []))
    totalReplies.value = response.total
    repliesPage.value++
  } catch (error) {
    console.error('Failed to fetch replies:', error)
  }
}

const toggleReplies = async () => {
  showReplies.value = !showReplies.value
  if (showReplies.value && (!props.comment.replies || props.comment.replies.length === 0)) {
    await loadReplies()
  }
}
</script>

<style scoped>
.comment-item {
  padding: 1rem;
  margin-bottom: 1rem;
}

.comment-item.is-reply {
  margin-left: 2rem;
  border-left: 2px solid var(--border-color);
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.comment-author {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.author-info {
  display: flex;
  flex-direction: column;
}

.author-name {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 0.875rem;
}

.comment-date {
  font-size: 0.75rem;
  color: var(--text-secondary);
}

.comment-actions {
  display: flex;
  gap: 0.5rem;
}

.action-btn {
  background: none;
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.75rem;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  transition: all 0.3s ease;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-primary);
}

.delete-btn:hover {
  border-color: #f56565;
  color: #f56565;
}

.comment-content {
  color: var(--text-primary);
  line-height: 1.6;
  white-space: pre-wrap;
}

.comment-replies {
  margin-top: 1rem;
  padding-left: 1rem;
  border-left: 2px solid var(--border-color);
}

@media (max-width: 768px) {
  .comment-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
  
  .comment-actions {
    align-self: flex-end;
  }
  
  .comment-item.is-reply {
    margin-left: 1rem;
  }
  
  .comment-replies {
    padding-left: 0.5rem;
  }
}
</style>