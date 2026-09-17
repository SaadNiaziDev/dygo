<script setup lang="ts">
import { Check, ChevronDown } from '@lucide/vue'
import {
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuItemIndicator,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuRoot,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from 'reka-ui'

import Button from '../atoms/Button.vue'
import IconButton from '../atoms/IconButton.vue'
import type { DropdownMenuItemModel } from '../types'

withDefaults(defineProps<{
  label?: string
  items?: DropdownMenuItemModel[]
  align?: 'start' | 'center' | 'end'
  sideOffset?: number
  triggerType?: 'button' | 'icon' | 'slot'
  panelClass?: string
}>(), {
  align: 'end',
  sideOffset: 6,
  triggerType: 'button',
  panelClass: '',
})

const emit = defineEmits<{
  select: [key: string]
  'update:checked': [key: string, checked: boolean]
  'update:open': [value: boolean]
  'close-auto-focus': [event: Event]
}>()

function preventCheckboxClose(event: Event) {
  event.preventDefault()
}
</script>

<template>
  <DropdownMenuRoot @update:open="emit('update:open', $event)">
    <DropdownMenuTrigger as-child>
      <slot v-if="triggerType === 'slot'" name="trigger" />

      <IconButton
        v-else-if="triggerType === 'icon'"
        class="d-dropdown-menu__trigger"
        type="button"
        variant="secondary"
        :label="label ?? 'Menu'"
      >
        <slot name="trigger">
          <ChevronDown :size="13" :stroke-width="1.9" aria-hidden="true" />
        </slot>
      </IconButton>

      <Button
        v-else
        class="d-dropdown-menu__trigger"
        type="button"
        variant="secondary"
        :aria-label="label"
      >
        <slot name="trigger">
          {{ label }}
          <ChevronDown :size="13" :stroke-width="1.9" aria-hidden="true" />
        </slot>
      </Button>
    </DropdownMenuTrigger>

    <DropdownMenuPortal>
      <DropdownMenuContent
        :class="panelClass || 'd-dropdown-menu__content'"
        :align="align"
        :side-offset="sideOffset"
        @close-auto-focus="emit('close-auto-focus', $event)"
      >
        <slot v-if="$slots.default" />

        <template v-else>
          <template v-for="item in items ?? []" :key="item.key">
            <DropdownMenuLabel v-if="item.type === 'label'" class="d-dropdown-menu__label">
              {{ item.label }}
            </DropdownMenuLabel>

            <DropdownMenuSeparator v-else-if="item.type === 'separator'" class="d-dropdown-menu__separator" />

            <DropdownMenuCheckboxItem
              v-else-if="item.type === 'checkbox'"
              class="d-dropdown-menu__item d-dropdown-menu__item--checkbox"
              :model-value="item.checked"
              :disabled="item.disabled"
              @select="preventCheckboxClose"
              @update:model-value="emit('update:checked', item.key, Boolean($event))"
            >
              <DropdownMenuItemIndicator class="d-dropdown-menu__indicator">
                <Check :size="13" :stroke-width="2.2" aria-hidden="true" />
              </DropdownMenuItemIndicator>
              <span>{{ item.label }}</span>
            </DropdownMenuCheckboxItem>

            <DropdownMenuItem
              v-else
              class="d-dropdown-menu__item"
              :disabled="item.disabled"
              @select="emit('select', item.key)"
            >
              {{ item.label }}
            </DropdownMenuItem>
          </template>
        </template>
      </DropdownMenuContent>
    </DropdownMenuPortal>
  </DropdownMenuRoot>
</template>

<style scoped>
.d-dropdown-menu__trigger {
  flex: 0 0 auto;
}
</style>
