<script setup lang="ts">
import { ChevronRight } from '@lucide/vue'
import {
  DropdownMenuPortal,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
} from 'reka-ui'

withDefaults(defineProps<{
  triggerClass?: string
  panelClass?: string
  sideOffset?: number
  align?: 'start' | 'center' | 'end'
}>(), {
  triggerClass: '',
  panelClass: '',
  sideOffset: 6,
  align: 'start',
})
</script>

<template>
  <DropdownMenuSub>
    <DropdownMenuSubTrigger class="d-menu__item d-menu__sub-trigger" :class="triggerClass">
      <slot name="trigger" />
      <ChevronRight class="d-menu__sub-chevron" :size="14" :stroke-width="1.8" aria-hidden="true" />
    </DropdownMenuSubTrigger>

    <DropdownMenuPortal>
      <DropdownMenuSubContent
        class="d-menu__content"
        :class="panelClass"
        :align="align"
        :side-offset="sideOffset"
      >
        <slot />
      </DropdownMenuSubContent>
    </DropdownMenuPortal>
  </DropdownMenuSub>
</template>
