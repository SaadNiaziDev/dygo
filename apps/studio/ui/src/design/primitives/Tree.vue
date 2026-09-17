<script setup lang="ts" generic="T extends { key: string; label: string; children?: T[] }">
import { ChevronRight } from '@lucide/vue'
import { TreeItem, TreeRoot, TreeVirtualizer } from 'reka-ui'

const props = withDefaults(defineProps<{
  items: T[]
  label: string
  expanded?: string[]
  estimateSize?: number
  indent?: number
  baseIndent?: number
}>(), {
  estimateSize: 34,
  indent: 20,
  baseIndent: 12,
})

const emit = defineEmits<{
  'update:expanded': [value: string[]]
  select: [item: T]
}>()

function toggleExpanded(item: T) {
  const keys = props.expanded ?? []
  emit('update:expanded', keys.includes(item.key) ? keys.filter(key => key !== item.key) : [...keys, item.key])
}

function preventClickToggle(event: Event) {
  const original = (event as CustomEvent<{ originalEvent?: Event }>).detail?.originalEvent
  if (original?.type === 'click') {
    event.preventDefault()
  }
}

function textContent(item: Record<string, unknown>): string {
  return String((item as T | undefined)?.label ?? '')
}
</script>

<template>
  <TreeRoot
    :expanded="expanded"
    :items="items"
    :get-key="(item: T) => item.key"
    :get-children="(item: T) => item.children"
    :aria-label="label"
    class="d-tree"
    @update:expanded="emit('update:expanded', $event)"
  >
    <TreeVirtualizer v-slot="{ item }" :estimate-size="estimateSize" :text-content="textContent">
      <TreeItem
        v-slot="{ isExpanded }"
        v-bind="item.bind"
        :value="item.value"
        class="d-tree__row"
        :style="{ paddingLeft: `${baseIndent + (item.level - 1) * indent}px` }"
        @select.prevent="emit('select', item.value as T)"
        @toggle="preventClickToggle"
      >
        <button
          v-if="item.value.children"
          type="button"
          tabindex="-1"
          class="d-tree__toggle"
          :aria-label="`${isExpanded ? 'Collapse' : 'Expand'} ${item.value.label}`"
          @click.stop="toggleExpanded(item.value as T)"
        >
          <slot name="toggle" :expanded="isExpanded">
            <ChevronRight
              class="d-tree__chevron"
              :class="{ 'd-tree__chevron--expanded': isExpanded }"
              :size="14"
              :stroke-width="1.8"
              aria-hidden="true"
            />
          </slot>
        </button>
        <span v-else class="d-tree__spacer" />

        <slot :item="item.value" :is-expanded="isExpanded" :level="item.level" />
      </TreeItem>
    </TreeVirtualizer>
  </TreeRoot>
</template>

<style scoped>
.d-tree {
  height: 100%;
  overflow: auto;
  padding: 4px 0;
  margin: 0;
  list-style: none;
}

.d-tree__row {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 34px;
  padding-right: 12px;
  color: var(--studio-text);
  cursor: pointer;
}

.d-tree__row:hover {
  background: var(--studio-surface-raised);
}

.d-tree__row:focus-visible {
  outline: 2px solid var(--studio-focus);
  outline-offset: -2px;
  background: var(--studio-surface-raised);
}

.d-tree__toggle {
  display: grid;
  place-items: center;
  flex: none;
  width: 20px;
  height: 26px;
  border: 0;
  background: transparent;
  color: var(--studio-text-muted);
  cursor: pointer;
}

.d-tree__chevron--expanded {
  transform: rotate(90deg);
}

.d-tree__spacer {
  width: 20px;
  flex: none;
}
</style>
