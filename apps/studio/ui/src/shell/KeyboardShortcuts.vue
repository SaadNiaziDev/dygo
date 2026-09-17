<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { Dialog } from '@/design'
import { studioCommands } from '@/features/commands/context'
import { shortcutLabel } from '@/features/commands/shortcuts'
import { useNavigationStore } from '@/stores/navigation.store'

const navigation = useNavigationStore()
const search = ref('')
let returnFocus: HTMLElement | null = null
function restoreFocus(event: Event) { event.preventDefault(); returnFocus?.focus() }
watch(() => navigation.shortcutsOpen, open => {
  if (open) { search.value = ''; returnFocus = document.activeElement instanceof HTMLElement ? document.activeElement : null }
}, { flush: 'sync' })
const groups = computed(() => {
  const result = new Map<string, typeof studioCommands.value>()
  for (const command of studioCommands.value) {
    if (!`${command.label} ${command.disabledReason ?? ''} ${shortcutLabel(command.shortcut)}`.toLowerCase().includes(search.value.trim().toLowerCase())) continue
    const group = command.group ?? 'This page'
    result.set(group, [...(result.get(group) ?? []), command])
  }
  return [...result]
})
</script>

<template>
  <Dialog
    v-model:open="navigation.shortcutsOpen"
    title="Keyboard shortcuts"
    description="Studio and current-page commands. Commands without a key are available from the command palette."
    description-class="sr-only"
    show-close
    close-label="Close keyboard shortcuts"
    panel-class="studio-command-menu__dialog"
    overlay-class="studio-command-menu__overlay"
    @close-auto-focus="restoreFocus"
  >
    <input v-model="search" class="studio-command-menu__input-wrap studio-command-menu__input" aria-label="Search keyboard shortcuts" placeholder="Search commands">
    <div class="studio-command-menu__list">
      <section v-for="[group, commands] in groups" :key="group">
        <h3 class="studio-command-menu__group-label">{{ group }}</h3>
        <div v-for="command in commands" :key="command.id" class="shortcut-row">
          <span>{{ command.label }}<small v-if="command.disabledReason">{{ command.disabledReason }}</small></span>
          <kbd>{{ shortcutLabel(command.shortcut) || 'Command palette' }}</kbd>
        </div>
      </section>
      <p v-if="!groups.length" role="status">No matching commands</p>
    </div>
  </Dialog>
</template>

<style scoped>
.shortcut-row { display: flex; align-items: center; justify-content: space-between; gap: 16px; padding: 10px 16px; font-size: 13px; }
.shortcut-row small { display: block; color: var(--studio-text-muted); }
.shortcut-row kbd { color: var(--studio-text-muted); white-space: nowrap; }
</style>
