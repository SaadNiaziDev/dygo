<script setup lang="ts">
import { computed } from 'vue'

import { Button, Dialog } from '@/design'
import { useDialogStore, type StudioDialog } from './dialogs.store'

const dialogStore = useDialogStore()
const topDialog = computed(() => dialogStore.topDialog)

function onOpenChange(open: boolean) {
  if (!open) {
    dialogStore.dismissTop()
  }
}

function choose(dialog: StudioDialog, key: string) {
  dialogStore.selectAction(dialog.id, key)
}

</script>

<template>
  <Dialog
    :open="Boolean(topDialog)"
    :dismissible="topDialog?.dismissible ?? true"
    :title="topDialog?.title"
    :description="topDialog?.content"
    :panel-attrs="{ 'data-type': topDialog?.type }"
    panel-class="studio-dialog"
    overlay-class="studio-dialog__overlay"
    title-class="studio-dialog__title"
    description-class="studio-dialog__content"
    @update:open="onOpenChange"
  >
    <div v-if="topDialog" class="studio-dialog__actions">
      <Button
        v-for="action in topDialog.actions"
        :key="action.key"
        :variant="action.variant"
        size="sm"
        @click="choose(topDialog, action.key)"
      >
        {{ action.label }}
      </Button>
    </div>
  </Dialog>
</template>

<style>
.studio-dialog__overlay {
  position: fixed;
  inset: 0;
  z-index: 80;
  background: var(--studio-overlay);
}

.studio-dialog {
  position: fixed;
  z-index: 81;
  top: 50%;
  left: 50%;
  width: min(calc(100vw - 32px), 420px);
  transform: translate(-50%, -50%);
  border: 1px solid var(--studio-border);
  border-radius: var(--studio-radius-sheet);
  background: var(--studio-surface);
  box-shadow: var(--studio-shadow);
  padding: 18px;
  display: grid;
  gap: 12px;
}

.studio-dialog:focus-visible {
  outline: 2px solid var(--studio-focus);
  outline-offset: 2px;
}

.studio-dialog__title {
  color: var(--studio-text);
  font-size: 16px;
  font-weight: 700;
  line-height: 1.25;
  margin: 0;
}

.studio-dialog__content {
  color: var(--studio-text-muted);
  font-size: 13px;
  line-height: 1.5;
  margin: 0;
  white-space: pre-wrap;
}

.studio-dialog__actions {
  display: flex;
  justify-content: flex-end;
  flex-wrap: wrap;
  gap: 8px;
  padding-top: 4px;
}
</style>
