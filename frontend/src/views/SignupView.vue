<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { apiService } from '../services/api'

const router = useRouter()

const email = ref('')
const password = ref('')
const errorMessage = ref('')
const successMessage = ref('')
const isLoading = ref(false)

const handleSignup = async () => {
  errorMessage.value = ''
  successMessage.value = ''

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
    await apiService.register({
      email: email.value,
      password: password.value
    })

    successMessage.value = 'Registration successful! Redirecting to login...'
    
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
    <div class="brand-header">
      <h1>Seismic Monitor</h1>
      <p>Create Account</p>
    </div>

    <div class="login-card">
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
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #1a1a2e;
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
  background-image: url('/earthquake_photo_2.jpg');
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
  color: #ff9f43;
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
  border-color: #ff9f43;
}

.form-group input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.error-banner {
  background-color: rgba(255, 159, 67, 0.1);
  color: #ff9f43;
  padding: 0.75rem;
  border-radius: 4px;
  margin-bottom: 1.5rem;
  font-size: 0.875rem;
  text-align: center;
  border: 1px solid #ff9f43;
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
  background: #ff9f43;
  color: white;
  border: none;
  border-radius: 4px;
  font-weight: bold;
  cursor: pointer;
  transition: background 0.2s;
}

.submit-btn:hover:not(:disabled) {
  background: #e67e22;
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
  color: #ff9f43;
  text-decoration: none;
}

.auth-links a:hover {
  text-decoration: underline;
}
</style>