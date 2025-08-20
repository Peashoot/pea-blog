<template>
  <div class="search-bar glass-effect">
    <div class="search-input-wrapper">
      <el-input
        v-model="searchQuery"
        :placeholder="$t('search_bar.placeholder')"
        size="large"
        clearable
        @keyup.enter="handleSearch"
        @clear="handleClear"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <button class="search-btn tech-button" @click="handleSearch">
        <el-icon><Search /></el-icon>
        <span v-if="!isMobile">{{ $t('search_bar.search_button') }}</span>
      </button>
    </div>
    
    <div class="filters" v-if="showFilters">
      <div class="filter-section">
        <label>{{ $t('search_bar.tags_filter') }}</label>
        <el-select
          v-model="selectedTags"
          multiple
          :placeholder="$t('search_bar.select_tags')"
          size="small"
          style="width: 200px"
        >
          <el-option
            v-for="tag in availableTags"
            :key="tag"
            :label="tag"
            :value="tag"
          />
        </el-select>
      </div>
      
      <div class="filter-section">
        <label>{{ $t('search_bar.sort_by') }}</label>
        <el-select v-model="sortBy" size="small" style="width: 120px">
          <el-option :label="$t('search_bar.latest')" value="created_at" />
          <el-option :label="$t('search_bar.most_viewed')" value="view_count" />
          <el-option :label="$t('search_bar.most_liked')" value="like_count" />
        </el-select>
      </div>
      
      <div class="filter-section">
        <label>{{ $t('search_bar.sort_order') }}</label>
        <el-select v-model="sortOrder" size="small" style="width: 100px">
          <el-option :label="$t('search_bar.descending')" value="desc" />
          <el-option :label="$t('search_bar.ascending')" value="asc" />
        </el-select>
      </div>
    </div>
    
    <div class="search-actions">
      <button 
        class="filter-toggle"
        @click="showFilters = !showFilters"
        :class="{ active: showFilters }"
      >
        <el-icon><Filter /></el-icon>
        <span v-if="!isMobile">{{ $t('search_bar.advanced_filter') }}</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useResponsive } from '@/composables'
import { debounce } from '@/utils'

const { isMobile } = useResponsive()

const searchQuery = ref('')
const selectedTags = ref<string[]>([])
const sortBy = ref('created_at')
const sortOrder = ref('desc')
const showFilters = ref(false)

const availableTags = ref([
  'Vue', 'TypeScript', 'JavaScript', 'React', 'Node.js', 
  'Python', 'Golang', '前端', '后端', '全栈',
  '算法', '数据结构', '系统设计', '数据库', 'Docker'
])

const emit = defineEmits<{
  search: [params: {
    keyword: string
    tags: string[]
    sortBy: string
    sortOrder: string
  }]
  clear: []
}>()

const handleSearch = () => {
  emit('search', {
    keyword: searchQuery.value,
    tags: selectedTags.value,
    sortBy: sortBy.value,
    sortOrder: sortOrder.value
  })
}

const handleClear = () => {
  searchQuery.value = ''
  selectedTags.value = []
  sortBy.value = 'created_at'
  sortOrder.value = 'desc'
  emit('clear')
}

const debouncedSearch = debounce(handleSearch, 500)

watch([selectedTags, sortBy, sortOrder], () => {
  if (searchQuery.value || selectedTags.value.length > 0) {
    debouncedSearch()
  }
})

watch(searchQuery, (newValue) => {
  if (newValue) {
    debouncedSearch()
  }
})
</script>

<style scoped>
.search-bar {
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.search-input-wrapper {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}

.search-input-wrapper .el-input {
  flex: 1;
}

.search-btn {
  padding: 0 1.5rem;
  white-space: nowrap;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.filters {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border-color);
}

.filter-section {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.filter-section label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  white-space: nowrap;
}

.search-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 1rem;
}

.filter-toggle {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid var(--border-color);
  color: var(--text-primary);
  padding: 0.5rem 1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.filter-toggle:hover {
  background: rgba(255, 255, 255, 0.2);
  border-color: var(--primary-color);
}

.filter-toggle.active {
  background: var(--tech-gradient);
  border-color: var(--primary-color);
}

@media (max-width: 768px) {
  .search-bar {
    padding: 1rem;
  }
  
  .search-input-wrapper {
    flex-direction: column;
    gap: 0.75rem;
  }
  
  .search-btn {
    justify-content: center;
  }
  
  .filters {
    flex-direction: column;
    gap: 0.75rem;
  }
  
  .filter-section {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.25rem;
  }
  
  .filter-section .el-select {
    width: 100% !important;
  }
}
</style>