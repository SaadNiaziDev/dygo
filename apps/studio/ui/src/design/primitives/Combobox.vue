<script setup lang="ts">
import { ComboboxContent, ComboboxEmpty, ComboboxInput, ComboboxPortal, ComboboxRoot } from 'reka-ui'

withDefaults(defineProps<{
  modelValue?: string
  open?: boolean
  search?: string
  displayValue?: (value: string) => string
  placeholder?: string
  label?: string
  disabled?: boolean
  ignoreFilter?: boolean
  portal?: boolean
  panelClass?: string
  inputClass?: string
  emptyText?: string
  sideOffset?: number
}>(), {
  portal: true,
  ignoreFilter: false,
  disabled: false,
  panelClass: '',
  inputClass: '',
  emptyText: 'No matches',
  sideOffset: 0,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'update:open': [value: boolean]
  'update:search': [value: string]
}>()
</script>

<template>
  <ComboboxRoot
    :model-value="modelValue"
    :open="open"
    :disabled="disabled"
    :ignore-filter="ignoreFilter"
    @update:model-value="emit('update:modelValue', String($event))"
    @update:open="emit('update:open', $event)"
  >
    <ComboboxInput
      class="d-combobox__input"
      :class="inputClass"
      :model-value="search"
      :display-value="displayValue"
      :placeholder="placeholder"
      :aria-label="label"
      @update:model-value="emit('update:search', $event)"
    />

    <ComboboxPortal v-if="portal">
      <ComboboxContent class="d-combobox" :class="panelClass" position="popper" :side-offset="sideOffset">
        <ComboboxEmpty v-if="emptyText" class="d-combobox__empty">{{ emptyText }}</ComboboxEmpty>
        <slot />
      </ComboboxContent>
    </ComboboxPortal>

    <ComboboxContent v-else :class="panelClass">
      <ComboboxEmpty v-if="emptyText" class="d-combobox__empty">{{ emptyText }}</ComboboxEmpty>
      <slot />
    </ComboboxContent>
  </ComboboxRoot>
</template>

<style scoped>
.d-combobox {
  border: 1px solid var(--studio-border);
  border-radius: var(--studio-radius-control);
  background: var(--studio-surface);
  box-shadow: var(--studio-shadow-sheet);
  color: var(--studio-text);
  outline: none;
}

.d-combobox__empty {
  display: block;
  padding: 6px;
  color: var(--studio-text-muted);
  font-size: 13px;
}
</style>
