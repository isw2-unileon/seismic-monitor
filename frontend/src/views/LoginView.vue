<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { apiService } from '../services/api'

const router = useRouter()

const email = ref('')
const password = ref('')
const errorMessage = ref('')
const isLoading = ref(false)

const handleLogin = async () => {
  errorMessage.value = ''

  // Client-side validation (Defensive programming to avoid unnecessary API calls)
  if (!email.value || !password.value) {
    errorMessage.value = 'Both email and password are required.'
    return
  }

  isLoading.value = true

  try {
    const response = await apiService.login({
      email: email.value,
      password: password.value
    })

    if (response.token) {
      // Store the JWT to satisfy the Navigation Guard
      localStorage.setItem('auth_token', response.token)
      
      // Store non-sensitive user data for UI purposes (like the alert radius)
      localStorage.setItem('user_data', JSON.stringify(response.user))

      router.push({ name: 'map' })
    } else {
      throw new Error('Invalid authentication payload received')
    }
  } catch (error) {
    // Catch generic errors or specific HTTP mock rejections
    errorMessage.value = error.message || 'Authentication failed. Please verify your credentials.'
    console.error('Login process aborted:', error)
  } finally {
    // Always remove the loading state, even if the request fails
    isLoading.value = false
  }
}
</script>

<template>
  <div class="login-wrapper">
    <div class="brand-header">
      <h1>Seismic Monitor</h1>
      <p>Sign In to Your Account</p>
    </div>

    <div class="login-card">
      <form @submit.prevent="handleLogin" class="auth-form">
        <div class="form-group">
          <label for="email">Email Address</label>
          <input
            id="email"
            v-model="email"
            type="email"
            placeholder="operator@seismic-monitor.org"
            :disabled="isLoading"
            autocomplete="email"
          />
        </div>

        <div class="form-group">
          <label for="password">Password</label>
          <input
            id="password"
            v-model="password"
            type="password"
            placeholder="••••••••"
            :disabled="isLoading"
            autocomplete="current-password"
          />
        </div>

        <div v-if="errorMessage" class="error-banner">
          {{ errorMessage }}
        </div>

        <button type="submit" :disabled="isLoading" class="submit-btn">
          {{ isLoading ? 'Authenticating...' : 'Secure Login' }}
        </button>
      </form>

      <div class="auth-links">
        <p>Don't have an account? <router-link to="/signup">Sign up here</router-link></p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-links {
  margin-top: 1.5rem;
  text-align: center;
  color: #a0aab2;
  font-size: 0.875rem;
}

.auth-links a {
  color: #e94560;
  text-decoration: none;
}

.auth-links a:hover {
  text-decoration: underline;
}

.login-wrapper {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #1a1a2e; /* Dark theme appropriate for monitoring tools */
  font-family: system-ui, -apple-system, sans-serif;
  position: relative;
  overflow: hidden;
  gap: 2rem;
}

.login-wrapper::before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: url('/earthquake_photo.jpg');
  background-size: cover;
  background-position: center;
  filter: blur(2px);
  opacity: 0.5;
  z-index: 0;
}

.login-card {
  position: relative;
  z-index: 1;
  background: rgba(22, 33, 62, 0.95);
  padding: 2.5rem;
  border-radius: 8px;
  box-shadow: 0 10px 25px rgba(0,0,0,0.5);
  width: 100%;
  max-width: 400px;
  border: 1px solid #2a3158;
  animation: fadeInUp 0.6s ease-out forwards;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.brand-header {
  position: relative;
  z-index: 1;
  text-align: center;
}

.brand-header h1 {
  color: #fff;
  margin: 0 0 0.5rem 0;
  font-size: 2.5rem;
  text-shadow: 0 2px 10px rgba(0,0,0,0.5);
}

.brand-header p {
  color: #e94560;
  margin: 0;
  font-size: 1rem;
  text-transform: uppercase;
  letter-spacing: 2px;
  font-weight: bold;
}

.form-group {
  margin-bottom: 1.5rem;
  display: flex;
  flex-direction: column;
}

.form-group label {
  color: #a0aab2;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
}

.form-group input {
  padding: 0.75rem;
  background: #0f172a;
  border: 1px solid #2a3158;
  color: #fff;
  border-radius: 4px;
  outline: none;
  transition: border-color 0.2s;
}

.form-group input:focus {
  border-color: #e94560;
}

.form-group input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.error-banner {
  background-color: rgba(233, 69, 96, 0.1);
  color: #e94560;
  padding: 0.75rem;
  border-radius: 4px;
  margin-bottom: 1.5rem;
  font-size: 0.875rem;
  text-align: center;
  border: 1px solid #e94560;
}

.submit-btn {
  width: 100%;
  padding: 0.875rem;
  background: #e94560;
  color: white;
  border: none;
  border-radius: 4px;
  font-weight: bold;
  cursor: pointer;
  transition: background 0.2s;
}

.submit-btn:hover:not(:disabled) {
  background: #d63d56;
}

.submit-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>