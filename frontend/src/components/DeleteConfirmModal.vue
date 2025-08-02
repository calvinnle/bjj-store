<template>
  <!-- Delete confirmation modal -->
  <div class="modal-backdrop">
    <div class="modal-content" style="max-width: 400px;">
        <!-- Warning icon -->
        <div style="display: flex; align-items: center; justify-content: center; width: 48px; height: 48px; margin: 0 auto 1rem; background-color: #fef2f2; border-radius: 50%;">
          <svg style="width: 24px; height: 24px; color: #dc2626;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"/>
          </svg>
        </div>

        <!-- Content -->
        <div class="text-center">
          <h3 style="font-size: 1.125rem; font-weight: 600; color: #111827; margin-bottom: 0.5rem;">Delete Item</h3>
          <p style="font-size: 0.875rem; color: #6b7280; margin-bottom: 1.5rem;">
            Are you sure you want to delete <strong>{{ itemName }}</strong>? This action cannot be undone.
          </p>
        </div>

        <!-- Action buttons -->
        <div style="display: flex; justify-content: center; gap: 1rem;">
          <button
            @click="$emit('cancel')"
            class="btn btn-secondary"
          >
            Cancel
          </button>
          <button
            @click="$emit('confirm')"
            class="btn btn-danger"
          >
            Delete
          </button>
        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  itemName: string
}

defineProps<Props>()
defineEmits<{
  confirm: []
  cancel: []
}>()
</script>