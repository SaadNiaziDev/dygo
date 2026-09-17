<script setup lang="ts">
import { DropdownMenuItem } from 'reka-ui'

const emit = defineEmits<{
  select: [event: Event]
}>()
</script>

<template>
  <DropdownMenuItem class="d-menu__item" @select="emit('select', $event)">
    <slot />
  </DropdownMenuItem>
</template>
