<script setup lang="ts">
import { ComboboxItem } from 'reka-ui'
import type { AcceptableValue } from 'reka-ui'

defineProps<{
  value: AcceptableValue
  disabled?: boolean
  textValue?: string
}>()
</script>

<template>
  <ComboboxItem class="d-combobox__item" :value="value" :disabled="disabled" :text-value="textValue">
    <slot />
  </ComboboxItem>
</template>

<style scoped>
.d-combobox__item {
  display: flex;
  align-items: center;
  border-radius: 5px;
  color: var(--studio-text);
  cursor: pointer;
  font-size: 13px;
  outline: none;
  padding: 6px;
}

.d-combobox__item[data-highlighted] {
  background: var(--studio-surface-raised);
}
</style>
