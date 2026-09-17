<script setup lang="ts">
import { X } from '@lucide/vue'
import {
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogOverlay,
  DialogPortal,
  DialogRoot,
  DialogTitle,
  DialogTrigger,
} from 'reka-ui'

const props = withDefaults(defineProps<{
  open?: boolean
  title?: string
  description?: string
  titleSrOnly?: boolean
  showClose?: boolean
  dismissible?: boolean
  panelClass?: string
  panelAttrs?: Record<string, unknown>
  overlayClass?: string
  titleClass?: string
  descriptionClass?: string
  headerClass?: string
  closeLabel?: string
}>(), {
  titleSrOnly: false,
  showClose: false,
  dismissible: true,
  panelClass: '',
  panelAttrs: () => ({}),
  overlayClass: '',
  titleClass: '',
  descriptionClass: '',
  headerClass: '',
  closeLabel: 'Close',
})

const emit = defineEmits<{
  'update:open': [value: boolean]
  'close-auto-focus': [event: Event]
}>()

function preventDismiss(event: Event) {
  if (!props.dismissible) {
    event.preventDefault()
  }
}
</script>

<template>
  <DialogRoot :open="open" modal @update:open="emit('update:open', $event)">
    <DialogTrigger v-if="$slots.trigger" as-child>
      <slot name="trigger" />
    </DialogTrigger>

    <DialogPortal>
      <DialogOverlay class="d-dialog__overlay" :class="overlayClass" />
      <DialogContent
        class="d-dialog"
        :class="panelClass"
        v-bind="panelAttrs"
        @escape-key-down="preventDismiss"
        @pointer-down-outside="preventDismiss"
        @interact-outside="preventDismiss"
        @close-auto-focus="emit('close-auto-focus', $event)"
      >
        <div v-if="title && !titleSrOnly && showClose" class="d-dialog__header" :class="headerClass">
          <DialogTitle class="d-dialog__title" :class="titleClass">{{ title }}</DialogTitle>
          <DialogClose class="d-dialog__close" :aria-label="closeLabel">
            <X :size="16" :stroke-width="1.8" aria-hidden="true" />
          </DialogClose>
        </div>
        <DialogTitle v-else-if="title && !titleSrOnly" class="d-dialog__title" :class="titleClass">{{ title }}</DialogTitle>
        <DialogTitle v-else-if="title" class="sr-only">{{ title }}</DialogTitle>

        <DialogDescription v-if="description" class="d-dialog__description" :class="descriptionClass">
          {{ description }}
        </DialogDescription>

        <slot />
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>

<style scoped>
.d-dialog {
  background: var(--studio-surface);
  color: var(--studio-text);
  outline: none;
}

.d-dialog:focus-visible {
  outline: 2px solid var(--studio-focus);
  outline-offset: 2px;
}

.d-dialog__overlay {
  position: fixed;
  inset: 0;
  background: var(--studio-overlay);
}

.d-dialog__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  border-bottom: 1px solid var(--studio-border);
  padding: 10px 16px;
}

.d-dialog__title {
  color: var(--studio-text);
  font-size: 16px;
  font-weight: 700;
  line-height: 1.25;
  margin: 0;
}

.d-dialog__close {
  display: inline-flex;
  width: 28px;
  height: 28px;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  border: 0;
  border-radius: var(--studio-radius-control);
  background: transparent;
  color: var(--studio-text-muted);
  cursor: pointer;
}

.d-dialog__close:hover {
  background: var(--studio-surface-raised);
  color: var(--studio-text);
}

.d-dialog__close:focus-visible {
  outline: 2px solid var(--studio-focus);
  outline-offset: 1px;
}

.d-dialog__description {
  color: var(--studio-text-muted);
  font-size: 13px;
  line-height: 1.5;
  margin: 0;
}
</style>
