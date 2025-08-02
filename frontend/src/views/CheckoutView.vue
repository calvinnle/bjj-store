<template>
  <div class="checkout">
    <!-- Navigation breadcrumb -->
    <nav class="breadcrumb">
      <div class="page-container">
        <div class="breadcrumb-content">
          <router-link to="/" style="color: #3b82f6;">Home</router-link>
          <svg style="width: 16px; height: 16px;" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
          </svg>
          <span style="color: #111827;">Checkout</span>
        </div>
      </div>
    </nav>

    <div class="page-container" style="padding: 2rem 1rem;">
      <!-- Empty Cart State -->
      <div v-if="!cart_store.hasItems" class="text-center py-16">
        <svg class="w-24 h-24 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M3 3h2l.4 2M7 13h10l4-8H5.4m1.6 8L3 3H1m6 10v6a2 2 0 002 2h4a2 2 0 002-2v-6m-6 0a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2z"/>
        </svg>
        <h2 class="text-2xl font-bold text-gray-900 mb-4">Your cart is empty</h2>
        <p class="text-gray-600 mb-8">Add some items to your cart before checking out.</p>
        <router-link to="/" class="btn btn-primary">
          Continue Shopping
        </router-link>
      </div>

      <!-- Checkout Form -->
      <div v-else class="two-column">
        <!-- Checkout Form -->
        <div style="display: flex; flex-direction: column; gap: 2rem;">
          <div>
            <h1 class="text-2xl font-bold text-gray-900 mb-8">Checkout</h1>

            <!-- Customer Information -->
            <div class="card">
              <h2 class="card-title">Customer Information</h2>
              <div class="form-row">
                <div class="form-group">
                  <label class="form-label">First Name *</label>
                  <input
                    v-model="form.shipping_address.first_name"
                    type="text"
                    required
                    class="form-input"
                    :class="{ 'form-input-error': errors.first_name }"
                  />
                  <p v-if="errors.first_name" class="form-error-text">{{ errors.first_name }}</p>
                </div>
                <div class="form-group">
                  <label class="form-label">Last Name *</label>
                  <input
                    v-model="form.shipping_address.last_name"
                    type="text"
                    required
                    class="form-input"
                    :class="{ 'form-input-error': errors.last_name }"
                  />
                  <p v-if="errors.last_name" class="form-error-text">{{ errors.last_name }}</p>
                </div>
              </div>
              <div class="form-group">
                <label class="form-label">Email *</label>
                <input
                  v-model="form.guest_email"
                  type="email"
                  required
                  class="form-input"
                  :class="{ 'form-input-error': errors.guest_email }"
                />
                <p v-if="errors.guest_email" class="form-error-text">{{ errors.guest_email }}</p>
              </div>
            </div>

            <!-- Shipping Address -->
            <div class="bg-white p-6 rounded-lg shadow-sm border">
              <h2 class="text-lg font-semibold mb-4">Shipping Address</h2>
              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">Address Line 1 *</label>
                  <input
                    v-model="form.shipping_address.address1"
                    type="text"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    :class="{ 'border-red-500': errors.address1 }"
                  />
                  <p v-if="errors.address1" class="text-red-500 text-sm mt-1">{{ errors.address1 }}</p>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">Address Line 2</label>
                  <input
                    v-model="form.shipping_address.address2"
                    type="text"
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
                <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">City *</label>
                    <input
                      v-model="form.shipping_address.city"
                      type="text"
                      required
                      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      :class="{ 'border-red-500': errors.city }"
                    />
                    <p v-if="errors.city" class="text-red-500 text-sm mt-1">{{ errors.city }}</p>
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">State *</label>
                    <input
                      v-model="form.shipping_address.state"
                      type="text"
                      required
                      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      :class="{ 'border-red-500': errors.state }"
                    />
                    <p v-if="errors.state" class="text-red-500 text-sm mt-1">{{ errors.state }}</p>
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">ZIP Code *</label>
                    <input
                      v-model="form.shipping_address.zip_code"
                      type="text"
                      required
                      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      :class="{ 'border-red-500': errors.zip_code }"
                    />
                    <p v-if="errors.zip_code" class="text-red-500 text-sm mt-1">{{ errors.zip_code }}</p>
                  </div>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">Country *</label>
                  <select
                    v-model="form.shipping_address.country"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    :class="{ 'border-red-500': errors.country }"
                  >
                    <option value="">Select Country</option>
                    <option value="US">United States</option>
                    <option value="CA">Canada</option>
                    <option value="UK">United Kingdom</option>
                    <option value="AU">Australia</option>
                    <option value="VN">Vietnam</option>
                  </select>
                  <p v-if="errors.country" class="text-red-500 text-sm mt-1">{{ errors.country }}</p>
                </div>
              </div>
            </div>

            <!-- Payment Information -->
            <div class="bg-white p-6 rounded-lg shadow-sm border">
              <h2 class="text-lg font-semibold mb-4">Payment Information</h2>
              
              <!-- Test Card Info -->
              <div class="mb-4 p-3 bg-blue-50 border border-blue-200 rounded-lg">
                <p class="text-sm text-blue-800 font-medium mb-2">Test Payment - Use these cards:</p>
                <div class="text-xs text-blue-700 space-y-1">
                  <p>✅ Success: 4242 4242 4242 4242</p>
                  <p>❌ Declined: 4000 0000 0000 0002</p>
                  <p>⚠️ Error: 4000 0000 0000 0119</p>
                </div>
              </div>

              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">Card Number *</label>
                  <input
                    v-model="paymentForm.card_number"
                    type="text"
                    placeholder="1234 5678 9012 3456"
                    maxlength="19"
                    @input="formatCardNumber"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    :class="{ 'border-red-500': errors.card_number }"
                  />
                  <p v-if="errors.card_number" class="text-red-500 text-sm mt-1">{{ errors.card_number }}</p>
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">Name on Card *</label>
                  <input
                    v-model="paymentForm.name_on_card"
                    type="text"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    :class="{ 'border-red-500': errors.name_on_card }"
                  />
                  <p v-if="errors.name_on_card" class="text-red-500 text-sm mt-1">{{ errors.name_on_card }}</p>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">Expiry Date *</label>
                    <input
                      v-model="paymentForm.expiry_date"
                      type="text"
                      placeholder="MM/YY"
                      maxlength="5"
                      @input="formatExpiryDate"
                      required
                      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      :class="{ 'border-red-500': errors.expiry_date }"
                    />
                    <p v-if="errors.expiry_date" class="text-red-500 text-sm mt-1">{{ errors.expiry_date }}</p>
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">CVV *</label>
                    <input
                      v-model="paymentForm.cvv"
                      type="text"
                      placeholder="123"
                      maxlength="4"
                      required
                      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      :class="{ 'border-red-500': errors.cvv }"
                    />
                    <p v-if="errors.cvv" class="text-red-500 text-sm mt-1">{{ errors.cvv }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Order Summary -->
        <div class="space-y-6">
          <div class="bg-white p-6 rounded-lg shadow-sm border sticky top-4">
            <h2 class="text-lg font-semibold mb-4">Order Summary</h2>
            
            <!-- Order Items -->
            <div class="space-y-4 mb-6">
              <div 
                v-for="item in cart_store.itemsWithProducts" 
                :key="`${item.product_id}-${item.size}`"
                class="flex items-center space-x-4"
              >
                <div class="w-16 h-16 bg-gray-200 rounded-lg overflow-hidden flex-shrink-0">
                  <img 
                    v-if="item.product?.image_url" 
                    :src="item.product.image_url" 
                    :alt="item.product.name"
                    class="w-full h-full object-cover"
                  />
                  <div v-else class="w-full h-full flex items-center justify-center text-gray-400">
                    <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clip-rule="evenodd" />
                    </svg>
                  </div>
                </div>
                <div class="flex-1">
                  <h4 class="text-sm font-medium text-gray-900">{{ item.product?.name }}</h4>
                  <p class="text-sm text-gray-500">Size: {{ item.size }}</p>
                  <p class="text-sm text-gray-500">Qty: {{ item.quantity }}</p>
                </div>
                <div class="text-sm font-medium text-gray-900">
                  ${{ (item.price * item.quantity).toFixed(2) }}
                </div>
              </div>
            </div>

            <!-- Order Total -->
            <div class="border-t pt-4 space-y-2">
              <div class="flex justify-between text-sm">
                <span>Subtotal</span>
                <span>${{ cart_store.totalPrice.toFixed(2) }}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span>Shipping</span>
                <span>Free</span>
              </div>
              <div class="flex justify-between text-lg font-semibold border-t pt-2">
                <span>Total</span>
                <span>${{ cart_store.totalPrice.toFixed(2) }}</span>
              </div>
            </div>

            <!-- Place Order Button -->
            <button
              @click="submitOrder"
              :disabled="loading || !isFormValid"
              class="w-full bg-blue-600 text-white py-3 px-4 rounded-lg font-semibold hover:bg-blue-700 transition-colors mt-6 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="loading">Processing...</span>
              <span v-else>Place Order (${{ cart_store.totalPrice.toFixed(2) }})</span>
            </button>

            <!-- Error Message -->
            <div v-if="submitError" class="mt-4 bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
              {{ submitError }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '@/stores/cart_store'
import { orderService } from '@/services/orders'
import { paymentService } from '@/services/payment'
import type { Address } from '@/types'

const router = useRouter()
const cart_store = useCartStore()

// Form data
const form = ref({
  guest_email: '',
  shipping_address: {
    first_name: '',
    last_name: '',
    address1: '',
    address2: '',
    city: '',
    state: '',
    zip_code: '',
    country: ''
  } as Address
})

// Payment form data
const paymentForm = ref({
  card_number: '',
  expiry_date: '',
  cvv: '',
  name_on_card: ''
})

// Form validation
const errors = ref({} as Record<string, string>)
const loading = ref(false)
const submitError = ref('')

// Computed
const isFormValid = computed(() => {
  return form.value.guest_email &&
         form.value.shipping_address.first_name &&
         form.value.shipping_address.last_name &&
         form.value.shipping_address.address1 &&
         form.value.shipping_address.city &&
         form.value.shipping_address.state &&
         form.value.shipping_address.zip_code &&
         form.value.shipping_address.country &&
         paymentForm.value.card_number &&
         paymentForm.value.expiry_date &&
         paymentForm.value.cvv &&
         paymentForm.value.name_on_card
})

// Methods
const validateForm = () => {
  errors.value = {}

  if (!form.value.guest_email) {
    errors.value.guest_email = 'Email is required'
  } else if (!/\S+@\S+\.\S+/.test(form.value.guest_email)) {
    errors.value.guest_email = 'Please enter a valid email'
  }

  if (!form.value.shipping_address.first_name) {
    errors.value.first_name = 'First name is required'
  }

  if (!form.value.shipping_address.last_name) {
    errors.value.last_name = 'Last name is required'
  }

  if (!form.value.shipping_address.address1) {
    errors.value.address1 = 'Address is required'
  }

  if (!form.value.shipping_address.city) {
    errors.value.city = 'City is required'
  }

  if (!form.value.shipping_address.state) {
    errors.value.state = 'State is required'
  }

  if (!form.value.shipping_address.zip_code) {
    errors.value.zip_code = 'ZIP code is required'
  }

  if (!form.value.shipping_address.country) {
    errors.value.country = 'Country is required'
  }

  // Payment validation
  if (!paymentForm.value.card_number) {
    errors.value.card_number = 'Card number is required'
  } else if (paymentForm.value.card_number.replace(/\s/g, '').length < 13) {
    errors.value.card_number = 'Card number is too short'
  }

  if (!paymentForm.value.name_on_card) {
    errors.value.name_on_card = 'Name on card is required'
  }

  if (!paymentForm.value.expiry_date) {
    errors.value.expiry_date = 'Expiry date is required'
  } else if (!/^\d{2}\/\d{2}$/.test(paymentForm.value.expiry_date)) {
    errors.value.expiry_date = 'Invalid expiry date format (MM/YY)'
  }

  if (!paymentForm.value.cvv) {
    errors.value.cvv = 'CVV is required'
  } else if (paymentForm.value.cvv.length < 3) {
    errors.value.cvv = 'CVV must be at least 3 digits'
  }

  return Object.keys(errors.value).length === 0
}

// Payment form formatting
const formatCardNumber = (event: Event) => {
  const target = event.target as HTMLInputElement
  let value = target.value.replace(/\s/g, '').replace(/\D/g, '')
  value = value.replace(/(\d{4})(?=\d)/g, '$1 ')
  target.value = value
  paymentForm.value.card_number = value
}

const formatExpiryDate = (event: Event) => {
  const target = event.target as HTMLInputElement
  let value = target.value.replace(/\D/g, '')
  if (value.length >= 2) {
    value = value.substring(0, 2) + '/' + value.substring(2, 4)
  }
  target.value = value
  paymentForm.value.expiry_date = value
}

const submitOrder = async () => {
  if (!validateForm() || !cart_store.hasItems) return

  loading.value = true
  submitError.value = ''

  try {
    // Step 1: Create order
    const orderData = {
      guest_email: form.value.guest_email,
      shipping_address: form.value.shipping_address,
      items: cart_store.getCheckoutData()
    }

    const orderResponse = await orderService.createOrder(orderData)
    const order = orderResponse.order

    // Step 2: Process payment
    const paymentData = {
      order_id: order.id,
      amount: order.total_amount,
      card_number: paymentForm.value.card_number,
      expiry_date: paymentForm.value.expiry_date,
      cvv: paymentForm.value.cvv,
      name_on_card: paymentForm.value.name_on_card
    }

    await paymentService.processPayment(paymentData)
    
    // Step 3: Clear cart and redirect
    cart_store.clearCart()
    router.push(`/order-confirmation/${order.order_number}`)
    
  } catch (error: any) {
    console.error('Order/Payment error:', error)
    
    if (error.response?.data?.error) {
      submitError.value = error.response.data.error
    } else if (error.response?.data?.message) {
      submitError.value = error.response.data.message
    } else {
      submitError.value = 'Failed to process order. Please try again.'
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  // Load cart from storage if empty
  if (!cart_store.hasItems) {
    cart_store.loadFromStorage()
  }
})
</script>