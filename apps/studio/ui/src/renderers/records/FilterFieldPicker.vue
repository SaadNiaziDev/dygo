<script setup lang="ts">
import { computed, ref } from 'vue'
import { FunnelPlus } from '@lucide/vue'
import { Combobox, ComboboxItem, IconButton, Popover } from '@/design'
import type { MetadataField } from '@/features/metadata/metadata.api'
const props = defineProps<{ fields: MetadataField[] }>()
const emit = defineEmits<{ select: [field: string] }>()
const open = ref(false)
const search = ref('')
const matches = computed(() => props.fields.filter((field) => `${field.label} ${field.name}`.toLowerCase().includes(search.value.toLowerCase())))
defineExpose({ open: () => { search.value = ''; open.value = true } })
function selectField(value: string) {
  emit('select', value)
  open.value = false
}
</script>

<template>
  <Popover v-model:open="open" panel-class="filter-field-picker" :side-offset="6" align="start" @update:open="search = ''">
    <template #trigger><IconButton label="Add filter"><FunnelPlus :size="14" /></IconButton></template>
    <Combobox
      :open="true"
      ignore-filter
      :portal="false"
      :search="search"
      label="Search filter fields"
      placeholder="Search fields"
      empty-text="No matching fields"
      @update:search="search = $event"
      @update:model-value="selectField"
    >
      <ComboboxItem v-for="field in matches" :key="field.name" :value="field.name">{{ field.label || field.name }}</ComboboxItem>
    </Combobox>
  </Popover>
</template>

<style>
.filter-field-picker { z-index: 60; width: 240px; max-height: 320px; overflow: auto; padding: 8px; border: 1px solid var(--studio-border); border-radius: var(--studio-radius-control); background: var(--studio-surface); box-shadow: var(--studio-shadow-sheet); color: var(--studio-text); }
.filter-field-picker input { width: 100%; padding: 6px; background: var(--studio-control-bg); border: 1px solid var(--studio-border); color: inherit; }
</style>
