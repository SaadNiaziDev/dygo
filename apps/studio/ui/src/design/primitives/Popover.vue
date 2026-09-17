<script setup lang="ts">
import { PopoverContent, PopoverPortal, PopoverRoot, PopoverTrigger } from 'reka-ui'

withDefaults(defineProps<{
  open?: boolean
  align?: 'start' | 'center' | 'end'
  side?: 'top' | 'right' | 'bottom' | 'left'
  sideOffset?: number
  panelClass?: string
}>(), {
  align: 'start',
  side: 'bottom',
  sideOffset: 6,
  panelClass: '',
})

const emit = defineEmits<{
  'update:open': [value: boolean]
}>()
</script>

<template>
  <PopoverRoot :open="open" @update:open="emit('update:open', $event)">
    <PopoverTrigger as-child>
      <slot name="trigger" />
    </PopoverTrigger>

    <PopoverPortal>
      <PopoverContent
        class="d-popover"
        :class="panelClass"
        :align="align"
        :side="side"
        :side-offset="sideOffset"
      >
        <slot />
      </PopoverContent>
    </PopoverPortal>
  </PopoverRoot>
</template>

<style scoped>
.d-popover {
  border-radius: var(--studio-radius-control);
  background: var(--studio-surface);
  color: var(--studio-text);
  outline: none;
}
</style>
