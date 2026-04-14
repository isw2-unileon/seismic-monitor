<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { apiService } from '../services/api'

// Initialize router for programmatic navigation
const router = useRouter()

// Reactive state variables for the form
const email = ref('')
const password = ref('')
const errorMessage = ref('')
const isLoading = ref(false)

const handleLogin = async () => {
  // 1. Reset previous errors
  errorMessage.value = ''

  // 2. Client-side validation (Defensive programming to avoid unnecessary API calls)
  if (!email.value || !password.value) {
    errorMessage.value = 'Both email and password are required.'
    return
  }

  isLoading.value = true

  try {
    // 3. Execute the mock API call
    const response = await apiService.login({
      email: email.value,
      password: password.value
    })

    // 4. Validate the contract structure and persist the session
    if (response.status === 'success' && response.data.token) {
      // Store the JWT to satisfy the Navigation Guard
      localStorage.setItem('auth_token', response.data.token)
      
      // Store non-sensitive user data for UI purposes (like the alert radius)
      localStorage.setItem('user_data', JSON.stringify(response.data.user))

      // 5. Force redirect to the protected map view
      router.push({ name: 'map' })
    } else {
      throw new Error('Invalid authentication payload received')
    }
  } catch (error) {
    // Catch generic errors or specific HTTP mock rejections
    errorMessage.value = 'Authentication failed. Please verify your credentials.'
    console.error('Login process aborted:', error)
  } finally {
    // Always remove the loading state, even if the request fails
    isLoading.value = false
  }
}
</script>

<template>
  <div class="login-wrapper">
    <div class="login-card">
      <div class="brand-header">
        <h1>Seismic Monitor</h1>
        <p>Restricted Access</p>
      </div>

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
    </div>
  </div>
</template>

<style scoped>
.login-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #1a1a2e; /* Dark theme appropriate for monitoring tools */
  font-family: system-ui, -apple-system, sans-serif;
}

.login-card {
  background: #16213e;
  padding: 2.5rem;
  border-radius: 8px;
  box-shadow: 0 10px 25px rgba(0,0,0,0.5);
  width: 100%;
  max-width: 400px;
  border: 1px solid #2a3158;
}

.brand-header {
  text-align: center;
  margin-bottom: 2rem;
}

.brand-header h1 {
  color: #fff;
  margin: 0 0 0.5rem 0;
  font-size: 1.5rem;
}

.brand-header p {
  color: #e94560;
  margin: 0;
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 1px;
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