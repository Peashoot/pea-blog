<template>
  <article class="article-card glass-effect" @click="$emit('click')">
    <div v-if="article.coverImage" class="article-cover">
      <img :src="article.coverImage" :alt="article.title" />
    </div>
    
    <div class="article-content">
      <div class="article-meta">
        <span class="author">{{ article.author?.username }}</span>
        <span class="separator">·</span>
        <span class="date">{{ formatDate(article.created_at) }}</span>
        <div class="tags" v-if="article.tags.length">
          <span 
            v-for="tag in article.tags.slice(0, 3)" 
            :key="tag" 
            class="tag"
          >
            #{{ tag }}
          </span>
        </div>
      </div>
      
      <h2 class="article-title">{{ article.title }}</h2>
      <p class="article-summary">{{ article.summary }}</p>
      
      <div class="article-stats">
        <div class="stat">
          <el-icon><View /></el-icon>
          <span>{{ article.view_count }}</span>
        </div>
        <div class="stat">
          <el-icon><Star /></el-icon>
          <span>{{ article.like_count }}</span>
        </div>
        <div class="stat">
          <el-icon><ChatDotRound /></el-icon>
           <span>{{ article.comment_count }}</span>
        </div>
      </div>
    </div>
  </article>
</template>

<script setup lang="ts">
import type { Article } from '@/types'
import { formatDate } from '@/utils'

defineProps<{
  article: Article
}>()

defineEmits<{
  click: []
}>()
</script>

<style scoped>
.article-card {
  cursor: pointer;
  transition: all 0.3s ease;
  overflow: hidden;
  margin-bottom: 1.5rem;
}

.article-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-heavy);
  border-color: var(--primary-color);
}

.article-cover {
  width: 100%;
  height: 200px;
  overflow: hidden;
  position: relative;
}

.article-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.article-card:hover .article-cover img {
  transform: scale(1.05);
}

.article-content {
  padding: 1.5rem;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.author {
  font-weight: 500;
  color: var(--primary-color);
}

.separator {
  opacity: 0.5;
}

.tags {
  display: flex;
  gap: 0.5rem;
  margin-left: auto;
}

.tag {
  background: rgba(255, 255, 255, 0.1);
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  color: var(--primary-color);
}

.article-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0 0 0.75rem 0;
  color: var(--text-primary);
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-summary {
  color: var(--text-secondary);
  line-height: 1.6;
  margin: 0 0 1rem 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-stats {
  display: flex;
  gap: 1rem;
  margin-top: auto;
  padding-top: 1rem;
  border-top: 1px solid var(--border-color);
}

.stat {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.stat .el-icon {
  font-size: 16px;
}

@media (max-width: 768px) {
  .article-content {
    padding: 1rem;
  }
  
  .article-title {
    font-size: 1.1rem;
  }
  
  .article-meta {
    flex-wrap: wrap;
  }
  
  .tags {
    margin-left: 0;
    margin-top: 0.5rem;
  }
}
</style>