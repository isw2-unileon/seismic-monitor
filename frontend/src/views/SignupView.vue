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
const successMessage = ref('')
const isLoading = ref(false)

const handleSignup = async () => {
  // 1. Reset previous messages
  errorMessage.value = ''
  successMessage.value = ''

  // 2. Client-side validation
  if (!email.value || !password.value) {
    errorMessage.value = 'Both email and password are required.'
    return
  }
  
  if (password.value.length < 6) {
    errorMessage.value = 'Password must be at least 6 characters long.'
    return
  }

  isLoading.value = true

  try {
    // 3. Execute the API call
    await apiService.register({
      email: email.value,
      password: password.value
    })

    successMessage.value = 'Registration successful! Redirecting to login...'
    
    // 4. Redirect to login after a short delay
    setTimeout(() => {
      router.push({ name: 'login' })
    }, 2000)

  } catch (error) {
    errorMessage.value = error.message || 'Registration failed. Please try again.'
    console.error('Registration process aborted:', error)
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="login-wrapper">
    <div class="login-card">
      <div class="brand-header">
        <h1>Seismic Monitor</h1>
        <p>Create Account</p>
      </div>

      <form @submit.prevent="handleSignup" class="auth-form">
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
            autocomplete="new-password"
          />
        </div>

        <div v-if="errorMessage" class="error-banner">
          {{ errorMessage }}
        </div>
        
        <div v-if="successMessage" class="success-banner">
          {{ successMessage }}
        </div>

        <button type="submit" :disabled="isLoading" class="submit-btn">
          {{ isLoading ? 'Creating Account...' : 'Sign Up' }}
        </button>
      </form>
      
      <div class="auth-links">
        <p>Already have an account? <router-link to="/login">Log in here</router-link></p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #1a1a2e;
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

.success-banner {
  background-color: rgba(76, 175, 80, 0.1);
  color: #4CAF50;
  padding: 0.75rem;
  border-radius: 4px;
  margin-bottom: 1.5rem;
  font-size: 0.875rem;
  text-align: center;
  border: 1px solid #4CAF50;
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
</style>