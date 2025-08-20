<template>
  <div class="article-editor">
    <div class="editor-header">
      <h1>{{ isEdit ? '编辑文章' : '新建文章' }}</h1>
      <div class="header-actions">
        <button class="secondary-btn" @click="saveDraft" :disabled="isLoading">
          <el-icon><DocumentCopy /></el-icon>
          保存草稿
        </button>
        <button class="tech-button" @click="publish" :disabled="isLoading">
          <el-icon><Promotion /></el-icon>
          {{ isEdit ? '更新文章' : '发布文章' }}
        </button>
      </div>
    </div>

    <div class="editor-content">
      <div class="editor-form glass-effect">
        <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
          <el-form-item label="标题" prop="title">
            <el-input
              v-model="form.title"
              placeholder="请输入文章标题"
              size="large"
            />
          </el-form-item>

          <el-form-item label="摘要" prop="summary">
            <el-input
              v-model="form.summary"
              type="textarea"
              :rows="3"
              placeholder="请输入文章摘要"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>

          <el-form-item label="标签" prop="tags">
            <el-select
              v-model="form.tags"
              multiple
              filterable
              allow-create
              default-first-option
              placeholder="请选择或输入标签"
              style="width: 100%"
            >
              <el-option
                v-for="tag in commonTags"
                :key="tag"
                :label="tag"
                :value="tag"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="封面图">
            <el-input
              v-model="form.coverImage"
              placeholder="请输入封面图URL（可选）"
            />
          </el-form-item>

          <el-form-item label="内容" prop="content">
            <div class="editor-container">
              <div class="editor-toolbar">
                <span class="toolbar-title">Markdown编辑器</span>
                <span class="toolbar-tip">支持Markdown语法</span>
              </div>
              <el-input
                v-model="form.content"
                type="textarea"
                :rows="20"
                placeholder="请输入文章内容（支持Markdown语法）"
                class="content-editor"
              />
            </div>
          </el-form-item>
        </el-form>
      </div>

      <div class="preview-panel glass-effect" v-if="form.content">
        <div class="preview-header">
          <h3>预览</h3>
        </div>
        <div class="preview-content" v-html="formattedContent"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useArticleStore } from '@/stores'
import { ElMessage, type FormInstance } from 'element-plus'

import { articleApi } from '@/api/articles';

const route = useRoute()
const router = useRouter()
const articleStore = useArticleStore()

const isEdit = computed(() => !!route.params.id)
const formRef = ref<FormInstance>()
const isLoading = ref(false)

const form = reactive({
  title: '',
  summary: '',
  content: '',
  tags: [] as string[],
  coverImage: '',
  status: 'draft' as 'draft' | 'published'
})

const rules = {
  title: [
    { required: true, message: '请输入文章标题', trigger: 'blur' },
    { min: 1, max: 200, message: '标题长度在1到200个字符', trigger: 'blur' }
  ],
  summary: [
    { required: true, message: '请输入文章摘要', trigger: 'blur' },
    { max: 500, message: '摘要长度不能超过500个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入文章内容', trigger: 'blur' }
  ]
}

const commonTags = [
  'Vue', 'TypeScript', 'JavaScript', 'React', 'Node.js',
  'Python', 'Golang', '前端', '后端', '全栈',
  '算法', '数据结构', '系统设计', '数据库', 'Docker'
]

const formattedContent = computed(() => {
  if (!form.content) return ''
  
  return form.content
    .replace(/\n/g, '<br>')
    .replace(/```(.*?)```/gs, '<pre><code>$1</code></pre>')
    .replace(/`([^`]+)`/g, '<code>$1</code>')
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.*?)\*/g, '<em>$1</em>')
    .replace(/#{3}\s*(.*)/g, '<h3>$1</h3>')
    .replace(/#{2}\s*(.*)/g, '<h2>$1</h2>')
    .replace(/#{1}\s*(.*)/g, '<h1>$1</h1>')
})

const saveDraft = async () => {
  if (!formRef.value) return
  
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  try {
    isLoading.value = true
    form.status = 'draft'
    
    if (isEdit.value) {
      await articleStore.updateArticle({
        id: Number(route.params.id),
        ...form
      })
      ElMessage.success('草稿保存成功')
    } else {
      await articleStore.createArticle(form)
      ElMessage.success('草稿创建成功')
      router.push('/admin/articles')
    }
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    isLoading.value = false
  }
}

const publish = async () => {
  if (!formRef.value) return
  
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  try {
    isLoading.value = true
    form.status = 'published'
    
    if (isEdit.value) {
      await articleStore.updateArticle({
        id: Number(route.params.id),
        ...form
      })
      ElMessage.success('文章更新成功')
    } else {
      await articleStore.createArticle(form)
      ElMessage.success('文章发布成功')
    }
    
    router.push('/admin/articles')
  } catch (error) {
    ElMessage.error('发布失败')
  } finally {
    isLoading.value = false
  }
}

onMounted(async () => {
  if (isEdit.value) {
    try {
      const article = await articleStore.fetchArticleById(Number(route.params.id))
      Object.assign(form, {
        title: article.title,
        summary: article.summary,
        content: article.content,
        tags: article.tags,
        coverImage: article.coverImage || '',
        status: article.status,
        publishedAt: article.publishedAt ? new Date(article.publishedAt) : null,
      })
    } catch (error) {
      ElMessage.error('加载文章失败')
      router.push('/admin/articles')
    }
  }
})
</script>

<style scoped>
.article-editor {
  max-width: 1400px;
  margin: 0 auto;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.editor-header h1 {
  color: var(--text-primary);
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 1rem;
}

.secondary-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.3s ease;
}

.secondary-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.editor-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
}

.editor-form {
  padding: 2rem;
}

.editor-container {
  width: 100%;
}

.editor-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1rem;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--border-color);
  border-bottom: none;
  border-radius: 8px 8px 0 0;
}

.toolbar-title {
  font-weight: 600;
  color: var(--text-primary);
}

.toolbar-tip {
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.content-editor {
  border-radius: 0 0 8px 8px !important;
}

.preview-panel {
  padding: 2rem;
  max-height: 80vh;
  overflow-y: auto;
}

.preview-header {
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--border-color);
}

.preview-header h3 {
  margin: 0;
  color: var(--text-primary);
}

.preview-content {
  color: var(--text-primary);
  line-height: 1.8;
}

.preview-content :deep(h1),
.preview-content :deep(h2),
.preview-content :deep(h3) {
  margin: 1.5rem 0 1rem 0;
  color: var(--text-primary);
}

.preview-content :deep(p) {
  margin: 1rem 0;
}

.preview-content :deep(pre) {
  background: rgba(0, 0, 0, 0.5);
  padding: 1rem;
  border-radius: 8px;
  overflow-x: auto;
  margin: 1rem 0;
}

.preview-content :deep(code) {
  background: rgba(255, 255, 255, 0.1);
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
}

.preview-content :deep(pre code) {
  background: none;
  padding: 0;
}

@media (max-width: 1200px) {
  .editor-content {
    grid-template-columns: 1fr;
  }
  
  .preview-panel {
    max-height: 60vh;
  }
}

@media (max-width: 768px) {
  .editor-header {
    flex-direction: column;
    gap: 1rem;
    align-items: flex-start;
  }
  
  .header-actions {
    width: 100%;
    justify-content: flex-end;
  }
  
  .editor-form,
  .preview-panel {
    padding: 1rem;
  }
}
</style>