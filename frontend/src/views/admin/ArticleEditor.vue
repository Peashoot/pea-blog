<template>
  <div class="article-editor">
    <div class="editor-header">
      <h1>{{ isEdit ? $t('article_editor.edit_article') : $t('article_editor.create_article') }}</h1>
      <div class="header-actions">
        <button class="secondary-btn" @click="handleSaveDraft" :disabled="isLoading">
          <el-icon><DocumentCopy /></el-icon>
          {{ $t('article_editor.save_draft') }}
        </button>
        <button class="tech-button" @click="handlePublishNow" :disabled="isLoading">
          <el-icon><Promotion /></el-icon>
          {{ isEdit ? $t('article_editor.publish') : $t('article_editor.publish') }}
        </button>
        <button class="tech-button" @click="handleSchedulePublish" :disabled="isLoading">
          <el-icon><Timer /></el-icon>
          {{ $t('article_editor.schedule_publish') }}
        </button>

        <!-- Schedule Publish Dialog -->
        <el-dialog
          v-model="showScheduleDialog"
          :title="$t('article_editor.schedule_publish')"
          width="500px"
          :before-close="handleScheduleDialogClose"
        >
          <div class="schedule-dialog-content">
            <div class="dialog-section">
              <label class="dialog-label">{{ $t('article_editor.select_publish_time') }}</label>
              <el-date-picker
                v-model="scheduledTime"
                type="datetime"
                :placeholder="$t('article_editor.select_publish_time')"
                :disabled-date="disabledPastDate"
                value-format="YYYY-MM-DD HH:mm:ss"
                size="large"
                style="width: 100%"
                @change="onScheduledTimeChange"
              />
            </div>
            
                      </div>
          
          <template #footer>
            <div class="dialog-footer">
              <el-button @click="showScheduleDialog = false">
                {{ $t('article_editor_page.cancel') }}
              </el-button>
              <el-button 
                type="primary" 
                :disabled="!scheduledTime || isScheduling"
                @click="confirmSchedulePublish"
                :loading="isScheduling"
              >
                <el-icon><Check /></el-icon>
                {{ $t('article_editor.confirm_schedule') }}
              </el-button>
            </div>
          </template>
        </el-dialog>
      </div>
    </div>

    <div class="editor-layout">
      <div class="editor-form glass-effect">
        <el-form ref="formRef" :model="form" :rules="rules" label-width="110px" class="custom-form">
          <el-form-item :label="$t('article_management.article_title')" prop="title">
            <el-input
              v-model="form.title"
              :placeholder="$t('article_editor.title_placeholder')"
              size="large"
            />
          </el-form-item>

          <el-form-item :label="$t('article_editor_page.summary')" prop="summary">
            <el-input
              v-model="form.summary"
              type="textarea"
              :rows="3"
              :placeholder="$t('article_editor_page.summary_placeholder')"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>

          <el-form-item :label="$t('article_editor_page.tags')" prop="tags">
            <el-select
              v-model="form.tags"
              multiple
              filterable
              allow-create
              default-first-option
              :placeholder="$t('article_editor.tags_placeholder')"
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

          <el-form-item :label="$t('article_editor_page.cover_image')">
            <el-input
              v-model="form.coverImage"
              :placeholder="$t('article_editor.cover_image_placeholder')"
            />
          </el-form-item>

          <el-form-item v-if="form.status === 'scheduled' && form.publishedAt" :label="$t('article_editor.current_schedule_time')">
            <div class="current-schedule-info">
              <el-icon><Timer /></el-icon>
              <span>{{ formatDateTimeForDisplay(form.publishedAt) }}</span>
              <el-button type="primary" size="small" @click="handleReschedule">
                <el-icon><Edit /></el-icon>
                {{ $t('article_editor.reschedule') }}
              </el-button>
            </div>
          </el-form-item>

          
          <el-form-item :label="$t('article_editor_page.content')" prop="content">
            <el-input
              v-model="form.content"
              type="textarea"
              :rows="20"
              :placeholder="$t('article_editor_page.content_placeholder')"
              class="content-editor"
            />
          </el-form-item>
        </el-form>
      </div>

      <!-- Preview Panel -->
      <div class="preview-panel glass-effect">
        <div class="preview-header">
          <h3>{{ $t('article_management.preview') }}</h3>
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
import { Timer, Check, Edit } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { getDefaultScheduleTime, isValidScheduleTime, formatDateTimeForDisplay } from '@/utils'


const { t } = useI18n()
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
  status: 'draft' as 'draft' | 'published' | 'scheduled',
  publishedAt: null as string | null
})

const scheduledTime = ref<string | null>(null)
const isScheduling = ref(false)
const showScheduleDialog = ref(false)

const rules = {
  title: [
    { required: true, message: t('article_editor_page.title_required'), trigger: 'blur' },
    { min: 1, max: 200, message: t('article_editor_page.title_length_error'), trigger: 'blur' }
  ],
  summary: [
    { required: true, message: t('article_editor_page.summary_required'), trigger: 'blur' },
    { max: 500, message: t('article_editor_page.summary_length_error'), trigger: 'blur' }
  ],
  content: [
    { required: true, message: t('article_editor_page.content_required'), trigger: 'blur' },
    { 
      validator: (rule: any, value: string, callback: (error?: Error) => void) => {
        if (value && value.length > 100 * 1024 * 1024) {
          callback(new Error(t('article_editor_page.content_size_limit')))
        } else {
          callback()
        }
      }, 
      trigger: 'blur' 
    }
  ]
}

const commonTags = [
  'Vue', 'TypeScript', 'JavaScript', 'React', 'Node.js',
  'Python', 'Golang', t('common.frontend'), t('common.backend'), t('common.full_stack'),
  t('common.algorithm'), t('common.data_structure'), t('common.system_design'), t('common.database'), 'Docker'
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

const handlePublishNow = async () => {
  if (!formRef.value) return
  
  // First validate the form
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) {
    ElMessage.error(t('article_editor_page.form_validation_failed'))
    return
  }
  
  form.status = 'published'
  form.publishedAt = null
  publish()
}

const handleSaveDraft = async () => {
  if (!formRef.value) return
  
  // For draft, we only validate title (minimum requirement)
  const valid = await formRef.value.validateField('title').catch(() => false)
  if (!valid) {
    ElMessage.error(t('article_editor_page.title_required'))
    return
  }
  
  await saveDraft()
}

const handleReschedule = () => {
  // Set the scheduled time to the current time by default
  if (form.publishedAt) {
    scheduledTime.value = form.publishedAt
  } else {
    scheduledTime.value = getDefaultScheduleTime()
  }
  showScheduleDialog.value = true
}

const handleSchedulePublish = async () => {
  if (!formRef.value) return
  
  // First validate the form
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) {
    ElMessage.error(t('article_editor_page.form_validation_failed'))
    return
  }
  
  // Set default scheduled time to 1 hour from now if not set
  if (!scheduledTime.value) {
    scheduledTime.value = getDefaultScheduleTime()
  }
  
  showScheduleDialog.value = true
}

const handleScheduleDialogClose = (done: () => void) => {
  if (isScheduling.value) {
    return
  }
  done()
}


const onScheduledTimeChange = (value: string) => {
  scheduledTime.value = value
}

const confirmSchedulePublish = async () => {
  if (!scheduledTime.value) {
    ElMessage.error(t('article_editor_page.select_publish_time'))
    return
  }
  
  if (!isValidScheduleTime(scheduledTime.value)) {
    ElMessage.error(t('article_editor_page.select_publish_time'))
    return
  }
  
  form.status = 'scheduled'
  form.publishedAt = scheduledTime.value
  
  try {
    isScheduling.value = true
    await publish()
    ElMessage.success(t('article_editor_page.article_schedule_success'))
    showScheduleDialog.value = false
    scheduledTime.value = null
  } catch (error) {
    console.error('Schedule publish error:', error)
  } finally {
    isScheduling.value = false
  }
}


const disabledPastDate = (time: Date) => {
  return time.getTime() < Date.now() - 8.64e7 // Disable past dates
}

const saveDraft = async () => {
  if (!formRef.value) return
  
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  try {
    isLoading.value = true
    form.status = 'draft'
    form.publishedAt = null
    
    if (isEdit.value) {
      await articleStore.updateArticle({
        id: Number(route.params.id),
        ...form
      })
      ElMessage.success(t('article_editor_page.draft_save_success'))
    } else {
      await articleStore.createArticle(form)
      ElMessage.success(t('article_editor_page.draft_create_success'))
      router.push('/admin/articles')
    }
  } catch {
    ElMessage.error(t('article_editor_page.save_fail'))
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
    
    if (isEdit.value) {
      await articleStore.updateArticle({
        id: Number(route.params.id),
        ...form
      })
      if (form.status === 'scheduled') {
        ElMessage.success(t('article_editor_page.article_schedule_success'))
      } else {
        ElMessage.success(t('article_editor_page.article_update_success'))
      }
    } else {
      await articleStore.createArticle(form)
      if (form.status === 'scheduled') {
        ElMessage.success(t('article_editor_page.article_schedule_success'))
      } else {
        ElMessage.success(t('article_editor_page.article_publish_success'))
      }
    }
    
    router.push('/admin/articles')
  } catch {
    ElMessage.error(t('article_editor_page.publish_fail'))
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
        coverImage: article.cover_image || '',
        status: article.status,
        publishedAt: article.published_at ? new Date(article.published_at).toISOString().slice(0, 19).replace('T', ' ') : null,
      })
    } catch {
      ElMessage.error(t('article_editor_page.load_fail'))
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

.editor-layout {
  display: flex;
  gap: 2rem;
  position: relative;
  height: calc(100vh - 145px);
}

.editor-form {
  flex: 1;
  min-width: 0;
  overflow-y: auto;
  padding: 2rem 1rem 1rem 0;
}

/* Editor form scrollbar */
.editor-form::-webkit-scrollbar {
  width: 6px;
}

.editor-form::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 3px;
}

.editor-form::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, rgba(0, 212, 255, 0.4) 0%, rgba(147, 51, 234, 0.4) 100%);
  border-radius: 3px;
  border: 1px solid rgba(255, 255, 255, 0.15);
}

.editor-form::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(180deg, rgba(0, 212, 255, 0.6) 0%, rgba(147, 51, 234, 0.6) 100%);
  border: 1px solid rgba(255, 255, 255, 0.25);
}

.preview-panel {
  width: 50%;
  transition: width 0.3s ease, opacity 0.3s ease;
}


.custom-form :deep(.el-form-item__label) {
  color: var(--text-label) !important;
  font-weight: 500;
}

.content-editor {
  border-radius: 8px !important;
}

.preview-panel {
  padding: 2rem;
  max-height: 80vh;
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
  overflow-y: auto;
}

/* Preview content scrollbar */
.preview-content::-webkit-scrollbar {
  width: 6px;
}

.preview-content::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.02);
  border-radius: 3px;
}

.preview-content::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, rgba(0, 212, 255, 0.3) 0%, rgba(147, 51, 234, 0.3) 100%);
  border-radius: 3px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.preview-content::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(180deg, rgba(0, 212, 255, 0.5) 0%, rgba(147, 51, 234, 0.5) 100%);
  border: 1px solid rgba(255, 255, 255, 0.2);
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


.schedule-dialog-content {
  padding: 1rem 0;
}

.dialog-section {
  margin-bottom: 1.5rem;
}

.dialog-section:last-child {
  margin-bottom: 0;
}

.dialog-label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: var(--text-primary);
  font-size: 0.9rem;
}

/* Fix date picker popup position */
:deep(.el-date-editor) {
  width: 100%;
}

:deep(.el-date-editor .el-input__wrapper) {
  width: 100%;
}

:deep(.el-picker-panel) {
  z-index: 3000 !important;
}

:deep(.el-date-picker__time-header) {
  margin-top: 8px;
}

:deep(.el-date-picker .el-input__prefix) {
  left: 12px;
}

:deep(.el-date-picker .el-input__suffix) {
  right: 12px;
}

.current-schedule-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 8px;
  color: #3b82f6;
  font-size: 0.875rem;
}

.current-schedule-info .el-button {
  margin-left: auto;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

@media (max-width: 1200px) {
  .editor-layout {
    flex-direction: column;
  }
  
  .preview-panel.show {
    width: 100%;
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
  
  .editor-form {
    padding: 1rem;
  }
  
  .preview-panel {
    padding: 1rem;
  }
}
</style>