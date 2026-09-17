<script setup lang="ts">
import { computed, nextTick, ref } from 'vue'
import { runStudioCommand } from '@/features/commands/context'
import { ariaShortcut, bindings, shortcutLabel } from '@/features/commands/shortcuts'
import { LogOut, Palette, RefreshCw } from '@lucide/vue'
import { useRouter } from 'vue-router'

import { queryClient } from '@/app/query'
import { reloadStudioApp } from '@/app/reload'
import Avatar from '@/design/atoms/Avatar.vue'
import {
  DropdownMenu,
  DropdownMenuItem,
  DropdownMenuRadioGroup,
  DropdownMenuSeparator,
  DropdownMenuSub,
} from '@/design'
import {
  getStudioThemePreference,
  isStudioThemePreference,
  setStudioThemePreference,
  studioThemeOptions,
  type StudioThemePreference,
} from '@/features/theme'
import { RouteName } from '@/router/routes'
import { useAuthStore } from '@/stores/auth.store'
import { usePreferencesStore } from '@/features/preferences/preferences.store'

withDefaults(defineProps<{
  userName?: string
  userAvatarUrl?: string
}>(), {
  userName: 'Studio user',
})

const router = useRouter()
const authStore = useAuthStore()
const reloading = ref(false)
const helpRequested = ref(false)
const trigger = ref<HTMLButtonElement | null>(null)
function menuClosed(event: Event) {
  if (!helpRequested.value) return
  event.preventDefault()
  helpRequested.value = false
  trigger.value?.focus()
  void nextTick(() => runStudioCommand('app:shortcuts'))
}
const preferences = usePreferencesStore()
const themePreference = computed<StudioThemePreference>(() => preferences.get('studio.theme', getStudioThemePreference()))

async function reloadApp() {
  if (reloading.value) {
    return
  }

  reloading.value = true
  try {
    await reloadStudioApp(router)
  } finally {
    reloading.value = false
  }
}

function onThemePreference(value: unknown) {
  if (!isStudioThemePreference(value)) {
    return
  }

  setStudioThemePreference(value)
}

async function logout() {
  await authStore.logout()
  queryClient.clear()
  await router.replace({ name: RouteName.Login })
}
</script>

<template>
  <DropdownMenu trigger-type="slot" panel-class="studio-user-menu__content" :side-offset="8" @close-auto-focus="menuClosed">
    <template #trigger>
      <button ref="trigger" class="studio-user-menu__trigger" type="button" :aria-label="`${userName} menu`">
        <Avatar :name="userName" :image-url="userAvatarUrl" />
      </button>
    </template>

    <DropdownMenuItem class="studio-user-menu__item" :aria-keyshortcuts="ariaShortcut(bindings['app:shortcuts']?.shortcut)" @select="helpRequested = true">
      <span>Keyboard shortcuts</span><kbd>{{ shortcutLabel(bindings['app:shortcuts']?.shortcut) }}</kbd>
    </DropdownMenuItem>
    <DropdownMenuItem class="studio-user-menu__item" :disabled="reloading" @select="reloadApp">
      <RefreshCw :size="14" :stroke-width="1.8" aria-hidden="true" />
      <span>Reload</span>
    </DropdownMenuItem>
    <DropdownMenuSub trigger-class="studio-user-menu__item" panel-class="studio-user-menu__content">
      <template #trigger>
        <Palette :size="14" :stroke-width="1.8" aria-hidden="true" />
        <span>Theme</span>
      </template>
      <DropdownMenuRadioGroup
        :model-value="themePreference"
        :options="studioThemeOptions"
        item-class="studio-user-menu__item studio-user-menu__item--radio"
        @update:model-value="onThemePreference"
      />
    </DropdownMenuSub>
    <DropdownMenuSeparator class="studio-user-menu__separator" />
    <DropdownMenuItem class="studio-user-menu__item" @select="logout">
      <LogOut :size="14" :stroke-width="1.8" aria-hidden="true" />
      <span>Logout</span>
    </DropdownMenuItem>
  </DropdownMenu>
</template>

<style scoped>
.studio-user-menu__trigger {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: inherit;
  padding: 0;
}

.studio-user-menu__trigger:focus-visible {
  outline: 2px solid var(--studio-focus);
  outline-offset: 2px;
}
</style>
