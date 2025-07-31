import api from './api'

export const uploadService = {
  // Upload image file
  async uploadImage(file: File): Promise<string> {
    const formData = new FormData()
    formData.append('image', file)

    try {
      const response = await api.post('/api/admin/upload/image', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        timeout: 30000, // 30 seconds timeout for file upload
      })

      return response.data.image_url
    } catch (error: any) {
      console.error('Upload service error:', error)
      
      if (error.response) {
        // Server responded with error
        const errorMessage = error.response.data?.error || `Upload failed with status ${error.response.status}`
        throw new Error(errorMessage)
      } else if (error.request) {
        // Network error
        throw new Error('Network error - please check your connection')
      } else {
        // Other error
        throw new Error(`Upload failed: ${error.message}`)
      }
    }
  },

  // Delete uploaded image
  async deleteImage(imageUrl: string): Promise<void> {
    try {
      await api.delete('/api/admin/upload/image', {
        data: { image_url: imageUrl },
      })
    } catch (error: any) {
      console.error('Delete image error:', error)
      throw new Error(error.response?.data?.error || 'Failed to delete image')
    }
  },
}