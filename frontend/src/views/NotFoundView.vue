<template>
  <div class="not-found-page">
    <div class="not-found-container">
      <div class="not-found-content glass-effect">
        <div class="error-animation">
          <h1 class="error-code gradient-text">404</h1>
          <div class="error-particles">
            <div class="particle" v-for="n in 20" :key="n"></div>
          </div>
        </div>
        
        <h2>{{ $t('not_found_page.title') }}</h2>
        <p>{{ $t('not_found_page.message') }}</p>
        
        <div class="actions">
          <router-link to="/" class="tech-button">
            <el-icon><HomeFilled /></el-icon>
            {{ $t('not_found_page.back_to_home') }}
          </router-link>
          <button class="tech-button secondary" @click="goBack">
            <el-icon><Back /></el-icon>
            {{ $t('not_found_page.go_back') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goBack = () => {
  if (window.history.length > 1) {
    router.go(-1)
  } else {
    router.push('/')
  }
}
</script>

<style scoped>
.not-found-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  background: var(--dark-bg);
  background-image: 
    radial-gradient(circle at 20% 80%, rgba(120, 119, 198, 0.3) 0%, transparent 50%),
    radial-gradient(circle at 80% 20%, rgba(255, 119, 198, 0.3) 0%, transparent 50%),
    radial-gradient(circle at 40% 40%, rgba(120, 219, 255, 0.2) 0%, transparent 50%);
}

.not-found-container {
  width: 100%;
  max-width: 500px;
}

.not-found-content {
  padding: 3rem 2rem;
  text-align: center;
}

.error-animation {
  position: relative;
  margin-bottom: 2rem;
}

.error-code {
  font-size: 6rem;
  font-weight: 900;
  margin: 0;
  position: relative;
  z-index: 2;
}

.error-particles {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 200px;
  height: 200px;
  pointer-events: none;
}

.particle {
  position: absolute;
  width: 4px;
  height: 4px;
  background: var(--primary-color);
  border-radius: 50%;
  opacity: 0.7;
  animation: float 3s infinite ease-in-out;
}

.particle:nth-child(odd) {
  animation-delay: -1.5s;
}

.particle:nth-child(1) { top: 10%; left: 20%; animation-duration: 2s; }
.particle:nth-child(2) { top: 20%; left: 80%; animation-duration: 2.5s; }
.particle:nth-child(3) { top: 30%; left: 40%; animation-duration: 3s; }
.particle:nth-child(4) { top: 40%; left: 70%; animation-duration: 2.2s; }
.particle:nth-child(5) { top: 50%; left: 10%; animation-duration: 2.8s; }
.particle:nth-child(6) { top: 60%; left: 90%; animation-duration: 2.3s; }
.particle:nth-child(7) { top: 70%; left: 30%; animation-duration: 2.7s; }
.particle:nth-child(8) { top: 80%; left: 60%; animation-duration: 2.1s; }
.particle:nth-child(9) { top: 15%; left: 50%; animation-duration: 2.6s; }
.particle:nth-child(10) { top: 35%; left: 15%; animation-duration: 2.4s; }
.particle:nth-child(11) { top: 45%; left: 85%; animation-duration: 2.9s; }
.particle:nth-child(12) { top: 55%; left: 25%; animation-duration: 2.0s; }
.particle:nth-child(13) { top: 65%; left: 75%; animation-duration: 2.8s; }
.particle:nth-child(14) { top: 75%; left: 45%; animation-duration: 2.3s; }
.particle:nth-child(15) { top: 85%; left: 35%; animation-duration: 2.6s; }
.particle:nth-child(16) { top: 25%; left: 65%; animation-duration: 2.1s; }
.particle:nth-child(17) { top: 5%; left: 75%; animation-duration: 2.7s; }
.particle:nth-child(18) { top: 95%; left: 55%; animation-duration: 2.4s; }
.particle:nth-child(19) { top: 85%; left: 15%; animation-duration: 2.2s; }
.particle:nth-child(20) { top: 15%; left: 85%; animation-duration: 2.9s; }

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
    opacity: 0.7;
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
    opacity: 1;
  }
}

.not-found-content h2 {
  font-size: 1.5rem;
  margin: 0 0 1rem 0;
  color: var(--text-primary);
}

.not-found-content p {
  color: var(--text-secondary);
  margin: 0 0 2rem 0;
  line-height: 1.6;
}

.actions {
  display: flex;
  justify-content: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.tech-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  text-decoration: none;
}

.tech-button.secondary {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid var(--border-color);
}

.tech-button.secondary:hover {
  background: rgba(255, 255, 255, 0.2);
}

@media (max-width: 480px) {
  .not-found-content {
    padding: 2rem 1rem;
  }
  
  .error-code {
    font-size: 4rem;
  }
  
  .not-found-content h2 {
    font-size: 1.25rem;
  }
  
  .actions {
    flex-direction: column;
  }
}
</style>