<script setup lang="ts">
import { Check } from '@lucide/vue'
import {
  DropdownMenuItemIndicator,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
} from 'reka-ui'

import type { FieldOption } from '../types'

withDefaults(defineProps<{
  modelValue?: string
  options: readonly FieldOption[]
  itemClass?: string
}>(), {
  itemClass: '',
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  select: [value: string]
}>()
</script>

<template>
  <DropdownMenuRadioGroup
    :model-value="modelValue"
    @update:model-value="emit('update:modelValue', String($event))"
  >
    <DropdownMenuRadioItem
      v-for="option in options"
      :key="option.value"
      class="d-menu__item d-menu__radio-item"
      :class="itemClass"
      :value="option.value"
      :disabled="option.disabled"
      @select="emit('select', option.value)"
    >
      <DropdownMenuItemIndicator class="d-menu__indicator">
        <Check :size="13" :stroke-width="2.2" aria-hidden="true" />
      </DropdownMenuItemIndicator>
      <slot :option="option">{{ option.label }}</slot>
    </DropdownMenuRadioItem>
  </DropdownMenuRadioGroup>
</template>
