import api from './api'

export interface PaymentRequest {
  order_id: number
  amount: number
  card_number: string
  expiry_date: string
  cvv: string
  name_on_card: string
}

export interface PaymentResponse {
  success: boolean
  transaction_id: string
  message: string
  order?: any
}

export const paymentService = {
  async processPayment(paymentData: PaymentRequest): Promise<PaymentResponse> {
    const response = await api.post('/api/payment/process', paymentData)
    return response.data
  },

  // Mock validation for demo purposes
  validateCard(cardNumber: string): { valid: boolean; type: string } {
    const cleaned = cardNumber.replace(/\s/g, '')
    
    if (cleaned.startsWith('4') && cleaned.length === 16) {
      return { valid: true, type: 'visa' }
    }
    if (cleaned.startsWith('5') && cleaned.length === 16) {
      return { valid: true, type: 'mastercard' }
    }
    if (cleaned.startsWith('3') && cleaned.length === 15) {
      return { valid: true, type: 'amex' }
    }
    
    return { valid: false, type: 'unknown' }
  },

  // Test card numbers for different scenarios
  getTestCards() {
    return {
      success: '4242424242424242',
      declined: '4000000000000002',
      error: '4000000000000119',
      amex: '378282246310005'
    }
  }
}