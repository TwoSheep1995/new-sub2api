<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-col justify-between gap-4 lg:flex-row lg:items-start">
          <div class="flex flex-1 flex-wrap items-center gap-3">
            <div class="relative w-full sm:w-80">
              <Icon
                name="search"
                size="md"
                class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-gray-500"
              />
              <input
                v-model.trim="searchQuery"
                type="text"
                class="input pl-10"
                :placeholder="t('admin.modelSquare.searchPlaceholder')"
              />
            </div>
            <select v-model="platformFilter" class="input w-full sm:w-44">
              <option value="all">{{ t('modelSquare.allPlatforms') }}</option>
              <option v-for="platform in platformOptions" :key="platform" :value="platform">
                {{ platformLabel(platform) }}
              </option>
            </select>
            <select v-model="visibilityFilter" class="input w-full sm:w-44">
              <option value="all">{{ t('admin.modelSquare.allRows') }}</option>
              <option value="enabled">{{ t('admin.modelSquare.enabledOnly') }}</option>
              <option value="disabled">{{ t('admin.modelSquare.disabledOnly') }}</option>
            </select>
          </div>

          <div class="flex w-full flex-shrink-0 flex-wrap items-center justify-end gap-3 lg:w-auto">
            <button type="button" class="btn btn-secondary" :disabled="loading" @click="loadRows">
              <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
            </button>
            <button type="button" class="btn btn-secondary" :disabled="saving || loading" @click="enableFiltered">
              {{ t('admin.modelSquare.enableFiltered') }}
            </button>
            <button type="button" class="btn btn-secondary" :disabled="saving || loading" @click="disableFiltered">
              {{ t('admin.modelSquare.disableFiltered') }}
            </button>
            <button type="button" class="btn btn-primary" :disabled="saving || loading" @click="saveRows">
              {{ saving ? t('common.saving') : t('common.save') }}
            </button>
          </div>
        </div>
      </template>

      <template #table>
        <div class="overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm dark:border-dark-800 dark:bg-dark-900">
          <div class="flex flex-wrap items-center justify-between gap-3 border-b border-gray-200 px-4 py-3 dark:border-dark-800">
            <div class="text-sm font-semibold text-gray-900 dark:text-white">
              {{ t('admin.modelSquare.title') }}
            </div>
            <div class="text-xs text-gray-500 dark:text-dark-400">
              {{ t('admin.modelSquare.summary', { enabled: enabledCount, total: rows.length }) }}
            </div>
          </div>

          <div v-if="loading" class="grid gap-3 p-4">
            <div v-for="idx in 8" :key="idx" class="h-14 animate-pulse rounded-lg bg-gray-100 dark:bg-dark-800"></div>
          </div>

          <div v-else-if="filteredRows.length === 0" class="px-4 py-14 text-center">
            <Icon name="search" size="xl" class="mx-auto text-gray-400" />
            <div class="mt-3 text-sm font-medium text-gray-900 dark:text-white">{{ t('admin.modelSquare.empty') }}</div>
          </div>

          <div v-else class="table-wrapper max-h-[calc(100vh-18rem)] overflow-auto">
            <table class="min-w-full text-left text-sm">
              <thead class="sticky top-0 z-10 bg-gray-50 text-xs uppercase text-gray-500 shadow-sm dark:bg-dark-950 dark:text-dark-400">
                <tr>
                  <th class="w-20 whitespace-nowrap px-4 py-3 font-medium">{{ t('admin.modelSquare.display') }}</th>
                  <th class="whitespace-nowrap px-4 py-3 font-medium">{{ t('modelSquare.model') }}</th>
                  <th class="whitespace-nowrap px-4 py-3 font-medium">{{ t('modelSquare.platform') }}</th>
                  <th class="whitespace-nowrap px-4 py-3 font-medium">{{ t('modelSquare.detail.channel') }}</th>
                  <th class="whitespace-nowrap px-4 py-3 font-medium">{{ t('modelSquare.group') }}</th>
                  <th class="whitespace-nowrap px-4 py-3 font-medium">{{ t('modelSquare.detail.rateMultiplier') }}</th>
                  <th class="whitespace-nowrap px-4 py-3 text-right font-medium">{{ t('admin.modelSquare.sortOrder') }}</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200 dark:divide-dark-800">
                <tr v-for="row in filteredRows" :key="rowKey(row)" class="hover:bg-gray-50/80 dark:hover:bg-dark-800/70">
                  <td class="whitespace-nowrap px-4 py-3">
                    <input
                      v-model="row.enabled"
                      type="checkbox"
                      class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                    />
                  </td>
                  <td class="min-w-64 px-4 py-3">
                    <div class="font-semibold text-gray-950 dark:text-white">{{ row.name }}</div>
                    <div class="mt-0.5 text-xs text-gray-500 dark:text-dark-400">{{ billingModeLabel(row) }}</div>
                  </td>
                  <td class="whitespace-nowrap px-4 py-3 text-gray-600 dark:text-dark-300">
                    {{ platformLabel(row.platform) }}
                  </td>
                  <td class="whitespace-nowrap px-4 py-3 text-gray-600 dark:text-dark-300">
                    {{ row.channel_name || '-' }}
                  </td>
                  <td class="whitespace-nowrap px-4 py-3 text-gray-600 dark:text-dark-300">
                    {{ row.group_name || '-' }}
                  </td>
                  <td class="whitespace-nowrap px-4 py-3 font-mono text-gray-600 dark:text-dark-300">
                    {{ formatMultiplier(row.rate_multiplier) }}x
                  </td>
                  <td class="whitespace-nowrap px-4 py-3 text-right">
                    <input
                      v-model.number="row.sort_order"
                      type="number"
                      min="0"
                      class="input ml-auto w-24 text-right"
                    />
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </template>
    </TablePageLayout>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import channelsAPI from '@/api/admin/channels'
import type { ModelSquareEntryRequest } from '@/api/admin/channels'
import type { PublicModelPricing } from '@/api/channels'
import { BILLING_MODE_IMAGE, BILLING_MODE_PER_REQUEST, BILLING_MODE_TOKEN } from '@/constants/channel'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'

type VisibilityFilter = 'all' | 'enabled' | 'disabled'

const { t } = useI18n()
const appStore = useAppStore()

const rows = ref<PublicModelPricing[]>([])
const loading = ref(false)
const saving = ref(false)
const searchQuery = ref('')
const platformFilter = ref('all')
const visibilityFilter = ref<VisibilityFilter>('all')

const platformOptions = computed(() =>
  Array.from(new Set(rows.value.map((row) => row.platform))).sort((a, b) =>
    platformLabel(a).localeCompare(platformLabel(b)),
  ),
)

const filteredRows = computed(() => {
  const query = searchQuery.value.toLowerCase()
  return rows.value.filter((row) => {
    const matchesSearch =
      !query ||
      row.name.toLowerCase().includes(query) ||
      row.platform.toLowerCase().includes(query) ||
      (row.channel_name || '').toLowerCase().includes(query) ||
      (row.group_name || '').toLowerCase().includes(query)
    const matchesPlatform = platformFilter.value === 'all' || row.platform === platformFilter.value
    const matchesVisibility =
      visibilityFilter.value === 'all' ||
      (visibilityFilter.value === 'enabled' ? row.enabled !== false : row.enabled === false)
    return matchesSearch && matchesPlatform && matchesVisibility
  })
})

const enabledCount = computed(() => rows.value.filter((row) => row.enabled !== false).length)

async function loadRows() {
  loading.value = true
  try {
    rows.value = (await channelsAPI.listModelSquareCandidates()).map((row, index) => ({
      ...row,
      enabled: row.enabled !== false,
      sort_order: row.sort_order ?? index + 1,
    }))
  } catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('common.error')))
  } finally {
    loading.value = false
  }
}

async function saveRows() {
  saving.value = true
  try {
    const entries: ModelSquareEntryRequest[] = rows.value.map((row, index) => ({
      channel_id: row.channel_id || 0,
      group_id: row.group_id || 0,
      platform: row.platform,
      model_name: row.name,
      enabled: row.enabled !== false,
      sort_order: Number.isFinite(Number(row.sort_order)) ? Number(row.sort_order) : index + 1,
    }))
    rows.value = await channelsAPI.updateModelSquareEntries(entries)
    appStore.showSuccess(t('admin.modelSquare.saved'))
  } catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('common.error')))
  } finally {
    saving.value = false
  }
}

function enableFiltered() {
  for (const row of filteredRows.value) row.enabled = true
}

function disableFiltered() {
  for (const row of filteredRows.value) row.enabled = false
}

function rowKey(row: PublicModelPricing): string {
  return `${row.channel_id}:${row.group_id}:${row.platform}:${row.name}`
}

function platformLabel(platform: string): string {
  switch (platform) {
    case 'openai':
      return 'OpenAI'
    case 'anthropic':
      return 'Claude'
    case 'gemini':
      return 'Gemini'
    case 'antigravity':
      return 'Antigravity'
    default:
      return platform
  }
}

function billingModeLabel(row: PublicModelPricing): string {
  switch ((row.effective_pricing ?? row.pricing)?.billing_mode) {
    case BILLING_MODE_PER_REQUEST:
      return t('modelSquare.billingModePerRequest')
    case BILLING_MODE_IMAGE:
      return t('modelSquare.billingModeImage')
    case BILLING_MODE_TOKEN:
      return t('modelSquare.billingModeToken')
    default:
      return t('modelSquare.noPricing')
  }
}

function formatMultiplier(value: number): string {
  return Number.isFinite(value) ? value.toPrecision(8).replace(/\.?0+$/, '') : '-'
}

onMounted(loadRows)
</script>
