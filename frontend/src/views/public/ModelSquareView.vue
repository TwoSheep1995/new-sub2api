<template>
  <div class="min-h-screen bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-white">
    <header class="sticky top-0 z-30 border-b border-gray-200/80 bg-white/90 backdrop-blur dark:border-dark-800 dark:bg-dark-950/90">
      <nav class="mx-auto flex max-w-7xl items-center justify-between px-4 py-3 sm:px-6">
        <router-link to="/home" class="flex min-w-0 items-center gap-3">
          <div class="h-9 w-9 overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <div class="hidden min-w-0 sm:block">
            <div class="truncate text-sm font-semibold text-gray-900 dark:text-white">{{ siteName }}</div>
            <div class="text-xs text-gray-500 dark:text-dark-400">{{ t('modelSquare.navSubtitle') }}</div>
          </div>
        </router-link>

        <div class="flex items-center gap-2">
          <LocaleSwitcher />
          <button
            type="button"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="toggleTheme"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="inline-flex items-center gap-1.5 rounded-lg bg-gray-900 px-3 py-2 text-xs font-medium text-white transition-colors hover:bg-gray-800 dark:bg-gray-800 dark:hover:bg-gray-700"
          >
            <Icon :name="isAuthenticated ? 'grid' : 'login'" size="sm" />
            <span>{{ isAuthenticated ? t('home.dashboard') : t('home.login') }}</span>
          </router-link>
        </div>
      </nav>
    </header>

    <main class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:py-10">
      <section class="mb-6 flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div class="max-w-3xl">
          <div class="mb-3 inline-flex items-center gap-2 rounded-full border border-primary-200 bg-primary-50 px-3 py-1 text-xs font-medium text-primary-700 dark:border-primary-900/70 dark:bg-primary-950/40 dark:text-primary-300">
            <Icon name="sparkles" size="xs" />
            <span>{{ t('modelSquare.badge') }}</span>
          </div>
          <h1 class="text-3xl font-bold tracking-normal text-gray-950 dark:text-white sm:text-4xl">
            {{ t('modelSquare.title') }}
          </h1>
          <p class="mt-3 text-sm leading-6 text-gray-600 dark:text-dark-300 sm:text-base">
            {{ t('modelSquare.subtitle') }}
          </p>
        </div>
        <div class="flex flex-wrap gap-2">
          <router-link
            to="/home"
            class="inline-flex items-center gap-1.5 rounded-lg border border-gray-200 bg-white px-3 py-2 text-xs font-medium text-gray-700 transition-colors hover:border-primary-300 hover:text-primary-600 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-200 dark:hover:border-primary-700 dark:hover:text-primary-400"
          >
            <Icon name="home" size="sm" />
            <span>{{ t('modelSquare.backHome') }}</span>
          </router-link>
          <button
            type="button"
            class="inline-flex items-center gap-1.5 rounded-lg bg-primary-600 px-3 py-2 text-xs font-medium text-white transition-colors hover:bg-primary-700 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="loading"
            @click="loadPublicPricing"
          >
            <Icon name="refresh" size="sm" :class="{ 'animate-spin': loading }" />
            <span>{{ t('modelSquare.refresh') }}</span>
          </button>
        </div>
      </section>

      <section class="mb-6 grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
        <div
          v-for="stat in stats"
          :key="stat.key"
          class="rounded-lg border border-gray-200 bg-white p-4 shadow-sm dark:border-dark-800 dark:bg-dark-900"
        >
          <div class="flex items-center justify-between gap-3">
            <span class="text-xs font-medium text-gray-500 dark:text-dark-400">{{ stat.label }}</span>
            <Icon :name="stat.icon" size="sm" class="text-primary-500" />
          </div>
          <div class="mt-2 text-2xl font-semibold text-gray-950 dark:text-white">{{ stat.value }}</div>
        </div>
      </section>

      <section class="mb-4 rounded-lg border border-gray-200 bg-white p-4 shadow-sm dark:border-dark-800 dark:bg-dark-900">
        <div class="grid gap-3 lg:grid-cols-[minmax(220px,1fr)_auto_auto] lg:items-center">
          <label class="relative block">
            <Icon name="search" size="sm" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            <input
              v-model.trim="searchTerm"
              type="search"
              class="w-full rounded-lg border border-gray-200 bg-white py-2 pl-9 pr-3 text-sm text-gray-900 outline-none transition-colors placeholder:text-gray-400 focus:border-primary-400 focus:ring-2 focus:ring-primary-100 dark:border-dark-700 dark:bg-dark-950 dark:text-white dark:focus:border-primary-600 dark:focus:ring-primary-950"
              :placeholder="t('modelSquare.searchPlaceholder')"
            />
          </label>

          <div class="flex flex-wrap items-center gap-2">
            <span class="inline-flex items-center gap-1 text-xs font-medium text-gray-500 dark:text-dark-400">
              <Icon name="filter" size="xs" />
              {{ t('modelSquare.platform') }}
            </span>
            <button
              v-for="option in platformOptions"
              :key="option.value"
              type="button"
              class="inline-flex items-center gap-1.5 rounded-lg border px-2.5 py-1.5 text-xs font-medium transition-colors"
              :class="selectedPlatform === option.value ? activeFilterClass : inactiveFilterClass"
              @click="selectedPlatform = option.value"
            >
              <PlatformIcon v-if="option.value !== 'all'" :platform="option.value as GroupPlatform" size="xs" />
              <Icon v-else name="globe" size="xs" />
              <span>{{ option.label }}</span>
            </button>
          </div>

          <div class="grid grid-cols-2 gap-2 sm:flex sm:items-center">
            <select
              v-model="selectedBillingMode"
              class="rounded-lg border border-gray-200 bg-white px-3 py-2 text-xs font-medium text-gray-700 outline-none focus:border-primary-400 focus:ring-2 focus:ring-primary-100 dark:border-dark-700 dark:bg-dark-950 dark:text-dark-200 dark:focus:border-primary-600 dark:focus:ring-primary-950"
            >
              <option value="all">{{ t('modelSquare.allBillingModes') }}</option>
              <option :value="BILLING_MODE_TOKEN">{{ t('modelSquare.billingModeToken') }}</option>
              <option :value="BILLING_MODE_PER_REQUEST">{{ t('modelSquare.billingModePerRequest') }}</option>
              <option :value="BILLING_MODE_IMAGE">{{ t('modelSquare.billingModeImage') }}</option>
            </select>
            <select
              v-model="sortBy"
              class="rounded-lg border border-gray-200 bg-white px-3 py-2 text-xs font-medium text-gray-700 outline-none focus:border-primary-400 focus:ring-2 focus:ring-primary-100 dark:border-dark-700 dark:bg-dark-950 dark:text-dark-200 dark:focus:border-primary-600 dark:focus:ring-primary-950"
            >
              <option value="platform">{{ t('modelSquare.sortByPlatform') }}</option>
              <option value="model">{{ t('modelSquare.sortByModel') }}</option>
              <option value="price">{{ t('modelSquare.sortByPrice') }}</option>
            </select>
          </div>
        </div>

        <div v-if="hasActiveFilters" class="mt-3 flex items-center justify-between gap-3 border-t border-gray-100 pt-3 dark:border-dark-800">
          <span class="text-xs text-gray-500 dark:text-dark-400">
            {{ t('modelSquare.filteredCount', { count: filteredRows.length, total: pricingRows.length }) }}
          </span>
          <button
            type="button"
            class="inline-flex items-center gap-1.5 rounded-lg px-2 py-1 text-xs font-medium text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            @click="clearFilters"
          >
            <Icon name="x" size="xs" />
            {{ t('modelSquare.clearFilters') }}
          </button>
        </div>
      </section>

      <section class="overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm dark:border-dark-800 dark:bg-dark-900">
        <div class="flex items-center justify-between gap-3 border-b border-gray-200 px-4 py-3 dark:border-dark-800">
          <div class="inline-flex items-center gap-2 text-sm font-semibold text-gray-900 dark:text-white">
            <Icon name="calculator" size="sm" class="text-primary-500" />
            <span>{{ t('modelSquare.tableTitle') }}</span>
          </div>
          <span class="text-xs text-gray-500 dark:text-dark-400">{{ t('modelSquare.cnyHint') }}</span>
        </div>

        <div v-if="loading" class="grid gap-3 p-4">
          <div v-for="idx in 8" :key="idx" class="h-14 animate-pulse rounded-lg bg-gray-100 dark:bg-dark-800"></div>
        </div>

        <div v-else-if="loadError" class="px-4 py-14 text-center">
          <Icon name="exclamationCircle" size="xl" class="mx-auto text-red-500" />
          <div class="mt-3 text-sm font-medium text-gray-900 dark:text-white">{{ t('modelSquare.loadFailed') }}</div>
          <button
            type="button"
            class="mt-4 inline-flex items-center gap-1.5 rounded-lg bg-primary-600 px-3 py-2 text-xs font-medium text-white transition-colors hover:bg-primary-700"
            @click="loadPublicPricing"
          >
            <Icon name="refresh" size="sm" />
            {{ t('modelSquare.retry') }}
          </button>
        </div>

        <div v-else-if="filteredRows.length === 0" class="px-4 py-14 text-center">
          <Icon name="search" size="xl" class="mx-auto text-gray-400" />
          <div class="mt-3 text-sm font-medium text-gray-900 dark:text-white">{{ t('modelSquare.empty') }}</div>
          <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ t('modelSquare.emptyHint') }}</p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="min-w-full text-left text-sm">
            <thead class="bg-gray-50 text-xs uppercase text-gray-500 dark:bg-dark-950/70 dark:text-dark-400">
              <tr>
                <th class="whitespace-nowrap px-4 py-3 font-medium">{{ t('modelSquare.model') }}</th>
                <th class="whitespace-nowrap px-4 py-3 font-medium">{{ t('modelSquare.platform') }}</th>
                <th class="hidden whitespace-nowrap px-4 py-3 font-medium md:table-cell">{{ t('modelSquare.group') }}</th>
                <th class="whitespace-nowrap px-4 py-3 font-medium">{{ t('modelSquare.billing') }}</th>
                <th class="whitespace-nowrap px-4 py-3 text-right font-medium">{{ t('modelSquare.input') }}</th>
                <th class="whitespace-nowrap px-4 py-3 text-right font-medium">{{ t('modelSquare.output') }}</th>
                <th class="hidden whitespace-nowrap px-4 py-3 text-right font-medium lg:table-cell">{{ t('modelSquare.cache') }}</th>
                <th class="whitespace-nowrap px-4 py-3 text-right font-medium">{{ t('modelSquare.actions') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-dark-800">
              <tr
                v-for="row in filteredRows"
                :key="`${row.platform}:${row.name}:${row.group_name}`"
                class="transition-colors hover:bg-gray-50/80 dark:hover:bg-dark-800/70"
              >
                <td class="min-w-64 px-4 py-3">
                  <div class="flex items-center gap-2.5">
                    <span class="inline-flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-gray-100 text-gray-600 dark:bg-dark-800 dark:text-dark-200">
                      <PlatformIcon :platform="row.platform as GroupPlatform" size="sm" />
                    </span>
                    <div class="min-w-0">
                      <div class="truncate font-semibold text-gray-950 dark:text-white">{{ row.name }}</div>
                      <div class="mt-0.5 text-xs text-gray-500 dark:text-dark-400 md:hidden">{{ row.group_name || '-' }}</div>
                    </div>
                  </div>
                </td>
                <td class="whitespace-nowrap px-4 py-3 text-gray-600 dark:text-dark-300">
                  {{ platformLabel(row.platform) }}
                </td>
                <td class="hidden whitespace-nowrap px-4 py-3 text-gray-600 dark:text-dark-300 md:table-cell">
                  {{ row.group_name || '-' }}
                </td>
                <td class="whitespace-nowrap px-4 py-3">
                  <span class="rounded-md bg-gray-100 px-2 py-1 text-xs font-medium text-gray-700 dark:bg-dark-800 dark:text-dark-200">
                    {{ billingModeLabel(row) }}
                  </span>
                </td>
                <td class="whitespace-nowrap px-4 py-3 text-right font-mono text-gray-900 dark:text-dark-100">
                  {{ formatInputPrice(row) }}
                </td>
                <td class="whitespace-nowrap px-4 py-3 text-right font-mono text-gray-900 dark:text-dark-100">
                  {{ formatOutputPrice(row) }}
                </td>
                <td class="hidden whitespace-nowrap px-4 py-3 text-right font-mono text-gray-600 dark:text-dark-300 lg:table-cell">
                  {{ formatCachePrice(row) }}
                </td>
                <td class="whitespace-nowrap px-4 py-3 text-right">
                  <button
                    type="button"
                    class="inline-flex items-center gap-1.5 rounded-lg border border-gray-200 bg-white px-2.5 py-1.5 text-xs font-medium text-gray-700 transition-colors hover:border-primary-300 hover:text-primary-600 dark:border-dark-700 dark:bg-dark-950 dark:text-dark-200 dark:hover:border-primary-700 dark:hover:text-primary-400"
                    @click="openDetails(row)"
                  >
                    <Icon name="eye" size="xs" />
                    <span>{{ t('modelSquare.viewDetails') }}</span>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </main>

    <div
      v-if="selectedModel"
      class="fixed inset-0 z-40 bg-gray-950/50 backdrop-blur-sm"
      @click.self="closeDetails"
    >
      <aside
        class="ml-auto flex h-full w-full max-w-2xl flex-col overflow-hidden bg-white shadow-2xl dark:bg-dark-900"
      >
        <div class="flex items-start justify-between gap-4 border-b border-gray-200 px-5 py-4 dark:border-dark-800">
          <div class="min-w-0">
            <div class="mb-2 inline-flex items-center gap-2 rounded-lg bg-gray-100 px-2.5 py-1 text-xs font-medium text-gray-700 dark:bg-dark-800 dark:text-dark-200">
              <PlatformIcon :platform="selectedModel.platform as GroupPlatform" size="xs" />
              <span>{{ platformLabel(selectedModel.platform) }}</span>
            </div>
            <h2 class="break-words text-xl font-semibold text-gray-950 dark:text-white">
              {{ selectedModel.name }}
            </h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">
              {{ t('modelSquare.detail.subtitle', { group: selectedModel.group_name || '-' }) }}
            </p>
          </div>
          <button
            type="button"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('common.close')"
            @click="closeDetails"
          >
            <Icon name="x" size="md" />
          </button>
        </div>

        <div class="flex-1 overflow-y-auto px-5 py-5">
          <div class="mb-5 grid gap-3 sm:grid-cols-2">
            <div
              v-for="item in selectedModelFacts"
              :key="item.label"
              class="rounded-lg border border-gray-200 bg-gray-50 p-3 dark:border-dark-800 dark:bg-dark-950"
            >
              <div class="text-xs font-medium text-gray-500 dark:text-dark-400">{{ item.label }}</div>
              <div class="mt-1 break-words text-sm font-semibold text-gray-950 dark:text-white">{{ item.value }}</div>
            </div>
          </div>

          <section class="mb-5 rounded-lg border border-gray-200 dark:border-dark-800">
            <div class="border-b border-gray-200 px-4 py-3 dark:border-dark-800">
              <h3 class="text-sm font-semibold text-gray-950 dark:text-white">{{ t('modelSquare.detail.pricingTitle') }}</h3>
            </div>
            <div class="grid gap-3 p-4 sm:grid-cols-2">
              <div
                v-for="item in selectedPricingFacts"
                :key="item.label"
                class="rounded-lg bg-gray-50 p-3 dark:bg-dark-950"
              >
                <div class="text-xs font-medium text-gray-500 dark:text-dark-400">{{ item.label }}</div>
                <div class="mt-1 font-mono text-sm text-gray-950 dark:text-white">{{ item.value }}</div>
              </div>
            </div>
            <div v-if="selectedModel.pricing?.intervals?.length" class="border-t border-gray-200 px-4 py-3 dark:border-dark-800">
              <div class="mb-2 text-xs font-medium text-gray-500 dark:text-dark-400">{{ t('modelSquare.detail.intervals') }}</div>
              <div class="space-y-2">
                <div
                  v-for="(interval, index) in selectedModel.pricing.intervals"
                  :key="`${interval.min_tokens}:${interval.max_tokens}:${index}`"
                  class="rounded-lg bg-gray-50 px-3 py-2 text-xs text-gray-700 dark:bg-dark-950 dark:text-dark-200"
                >
                  {{ formatInterval(interval) }}
                </div>
              </div>
            </div>
          </section>

          <section class="mb-5 rounded-lg border border-gray-200 dark:border-dark-800">
            <div class="border-b border-gray-200 px-4 py-3 dark:border-dark-800">
              <h3 class="text-sm font-semibold text-gray-950 dark:text-white">{{ t('modelSquare.detail.howToUseTitle') }}</h3>
            </div>
            <div class="space-y-4 p-4">
              <div class="grid gap-3 sm:grid-cols-2">
                <div class="rounded-lg bg-gray-50 p-3 dark:bg-dark-950">
                  <div class="text-xs font-medium text-gray-500 dark:text-dark-400">{{ t('modelSquare.detail.endpoint') }}</div>
                  <div class="mt-1 font-mono text-sm text-gray-950 dark:text-white">{{ endpointPath(selectedModel) }}</div>
                </div>
                <div class="rounded-lg bg-gray-50 p-3 dark:bg-dark-950">
                  <div class="text-xs font-medium text-gray-500 dark:text-dark-400">{{ t('modelSquare.detail.modelParam') }}</div>
                  <div class="mt-1 font-mono text-sm text-gray-950 dark:text-white">{{ selectedModel.name }}</div>
                </div>
              </div>

              <div>
                <div class="mb-2 text-xs font-medium text-gray-500 dark:text-dark-400">{{ t('modelSquare.detail.curlExample') }}</div>
                <pre class="overflow-x-auto rounded-lg bg-gray-950 p-4 text-xs leading-5 text-gray-100"><code>{{ curlExample(selectedModel) }}</code></pre>
              </div>

              <div>
                <div class="mb-2 text-xs font-medium text-gray-500 dark:text-dark-400">{{ t('modelSquare.detail.envExample') }}</div>
                <pre class="overflow-x-auto rounded-lg bg-gray-950 p-4 text-xs leading-5 text-gray-100"><code>{{ envExample(selectedModel) }}</code></pre>
              </div>
            </div>
          </section>

          <section class="rounded-lg border border-gray-200 dark:border-dark-800">
            <div class="border-b border-gray-200 px-4 py-3 dark:border-dark-800">
              <h3 class="text-sm font-semibold text-gray-950 dark:text-white">{{ t('modelSquare.detail.notesTitle') }}</h3>
            </div>
            <div class="space-y-3 p-4 text-sm leading-6 text-gray-600 dark:text-dark-300">
              <p>{{ usageDescription(selectedModel) }}</p>
              <p>{{ billingDescription(selectedModel) }}</p>
              <p>{{ t('modelSquare.detail.authNote') }}</p>
            </div>
          </section>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore, useAuthStore } from '@/stores'
import {
  getPublicModelPricing,
  type PublicModelPricing,
  type UserPricingInterval
} from '@/api/channels'
import {
  BILLING_MODE_IMAGE,
  BILLING_MODE_PER_REQUEST,
  BILLING_MODE_TOKEN,
  type BillingMode
} from '@/constants/channel'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import PlatformIcon from '@/components/common/PlatformIcon.vue'
import Icon from '@/components/icons/Icon.vue'
import { formatScaled } from '@/utils/pricing'
import type { GroupPlatform } from '@/types'

type SortBy = 'platform' | 'model' | 'price'
type FilterPlatform = 'all' | string
type FilterBillingMode = 'all' | BillingMode

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const serviceOrigin = computed(() => window.location.origin)
const isAuthenticated = computed(() => authStore.isAuthenticated)
const dashboardPath = computed(() => (authStore.isAdmin ? '/admin/dashboard' : '/dashboard'))
const isDark = ref(document.documentElement.classList.contains('dark'))

const pricingRows = ref<PublicModelPricing[]>([])
const loading = ref(false)
const loadError = ref(false)
const searchTerm = ref('')
const selectedPlatform = ref<FilterPlatform>('all')
const selectedBillingMode = ref<FilterBillingMode>('all')
const sortBy = ref<SortBy>('platform')
const selectedModel = ref<PublicModelPricing | null>(null)
const perMillionScale = 1_000_000

const activeFilterClass = 'border-primary-500 bg-primary-50 text-primary-700 dark:border-primary-700 dark:bg-primary-950/40 dark:text-primary-300'
const inactiveFilterClass = 'border-gray-200 bg-white text-gray-600 hover:border-primary-300 hover:text-primary-600 dark:border-dark-700 dark:bg-dark-950 dark:text-dark-300 dark:hover:border-primary-700 dark:hover:text-primary-400'

const platformOptions = computed(() => {
  const platforms = Array.from(new Set(pricingRows.value.map((row) => row.platform))).sort((a, b) =>
    platformLabel(a).localeCompare(platformLabel(b))
  )
  return [
    { value: 'all', label: t('modelSquare.allPlatforms') },
    ...platforms.map((platform) => ({ value: platform, label: platformLabel(platform) }))
  ]
})

const filteredRows = computed(() => {
  const query = searchTerm.value.toLowerCase()
  const rows = pricingRows.value.filter((row) => {
    const matchesSearch =
      !query ||
      row.name.toLowerCase().includes(query) ||
      row.platform.toLowerCase().includes(query) ||
      (row.group_name || '').toLowerCase().includes(query)
    const matchesPlatform = selectedPlatform.value === 'all' || row.platform === selectedPlatform.value
    const matchesBilling =
      selectedBillingMode.value === 'all' || row.pricing?.billing_mode === selectedBillingMode.value

    return matchesSearch && matchesPlatform && matchesBilling
  })

  return rows.sort((a, b) => {
    if (sortBy.value === 'model') return a.name.localeCompare(b.name)
    if (sortBy.value === 'price') return comparablePrice(a) - comparablePrice(b)

    const platformCompare = platformLabel(a.platform).localeCompare(platformLabel(b.platform))
    if (platformCompare !== 0) return platformCompare
    return a.name.localeCompare(b.name)
  })
})

const stats = computed(() => {
  const platforms = new Set(pricingRows.value.map((row) => row.platform))
  const tokenModels = pricingRows.value.filter((row) => row.pricing?.billing_mode === BILLING_MODE_TOKEN).length
  const requestModels = pricingRows.value.filter((row) =>
    row.pricing?.billing_mode === BILLING_MODE_PER_REQUEST || row.pricing?.billing_mode === BILLING_MODE_IMAGE
  ).length

  return [
    { key: 'total', label: t('modelSquare.totalModels'), value: pricingRows.value.length, icon: 'grid' as const },
    { key: 'platforms', label: t('modelSquare.providerCount'), value: platforms.size, icon: 'globe' as const },
    { key: 'token', label: t('modelSquare.tokenModels'), value: tokenModels, icon: 'calculator' as const },
    { key: 'request', label: t('modelSquare.requestModels'), value: requestModels, icon: 'bolt' as const }
  ]
})

const hasActiveFilters = computed(() =>
  Boolean(searchTerm.value) || selectedPlatform.value !== 'all' || selectedBillingMode.value !== 'all'
)

const selectedModelFacts = computed(() => {
  if (!selectedModel.value) return []
  return [
    { label: t('modelSquare.platform'), value: platformLabel(selectedModel.value.platform) },
    { label: t('modelSquare.group'), value: selectedModel.value.group_name || '-' },
    { label: t('modelSquare.detail.billingMode'), value: billingModeLabel(selectedModel.value) },
    { label: t('modelSquare.detail.rateMultiplier'), value: `${formatNumber(selectedModel.value.rate_multiplier)}x` }
  ]
})

const selectedPricingFacts = computed(() => {
  if (!selectedModel.value) return []
  return [
    { label: t('modelSquare.input'), value: formatInputPrice(selectedModel.value) },
    { label: t('modelSquare.output'), value: formatOutputPrice(selectedModel.value) },
    { label: t('modelSquare.cache'), value: formatCachePrice(selectedModel.value) },
    { label: t('modelSquare.detail.singleRequest'), value: formatPerRequestPrice(selectedModel.value) }
  ]
})

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  ) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

async function loadPublicPricing() {
  loading.value = true
  loadError.value = false
  try {
    pricingRows.value = await getPublicModelPricing()
  } catch {
    loadError.value = true
    pricingRows.value = []
  } finally {
    loading.value = false
  }
}

function clearFilters() {
  searchTerm.value = ''
  selectedPlatform.value = 'all'
  selectedBillingMode.value = 'all'
}

function openDetails(row: PublicModelPricing) {
  selectedModel.value = row
}

function closeDetails() {
  selectedModel.value = null
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
  switch (row.pricing?.billing_mode) {
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

function formatInputPrice(row: PublicModelPricing): string {
  const pricing = row.pricing
  if (!pricing || pricing.billing_mode !== BILLING_MODE_TOKEN) return '-'
  return formatScaled(pricing.input_price, perMillionScale)
}

function formatOutputPrice(row: PublicModelPricing): string {
  const pricing = row.pricing
  if (!pricing) return '-'
  if (pricing.billing_mode === BILLING_MODE_PER_REQUEST) {
    return formatScaled(pricing.per_request_price, 1)
  }
  if (pricing.billing_mode === BILLING_MODE_IMAGE) {
    return formatScaled(pricing.image_output_price ?? pricing.per_request_price, 1)
  }
  return formatScaled(pricing.output_price, perMillionScale)
}

function formatCachePrice(row: PublicModelPricing): string {
  const pricing = row.pricing
  if (!pricing || pricing.billing_mode !== BILLING_MODE_TOKEN) return '-'
  if (pricing.cache_write_price == null && pricing.cache_read_price == null) return '-'
  return `${formatScaled(pricing.cache_write_price, perMillionScale)} / ${formatScaled(pricing.cache_read_price, perMillionScale)}`
}

function formatPerRequestPrice(row: PublicModelPricing): string {
  const pricing = row.pricing
  if (!pricing) return '-'
  if (pricing.billing_mode === BILLING_MODE_PER_REQUEST) {
    return formatScaled(pricing.per_request_price, 1)
  }
  if (pricing.billing_mode === BILLING_MODE_IMAGE) {
    return formatScaled(pricing.image_output_price ?? pricing.per_request_price, 1)
  }
  return '-'
}

function formatNumber(value: number): string {
  return Number.isFinite(value) ? value.toPrecision(8).replace(/\.?0+$/, '') : '-'
}

function endpointPath(row: PublicModelPricing): string {
  if (row.platform === 'anthropic') return '/v1/messages'
  return '/v1/chat/completions'
}

function curlExample(row: PublicModelPricing): string {
  if (row.platform === 'anthropic') {
    return [
      `curl ${serviceOrigin.value}/v1/messages \\`,
      '  -H "x-api-key: sk-your-api-key" \\',
      '  -H "anthropic-version: 2023-06-01" \\',
      '  -H "content-type: application/json" \\',
      '  -d \'{',
      `    "model": "${row.name}",`,
      '    "max_tokens": 1024,',
      '    "messages": [',
      '      { "role": "user", "content": "你好，介绍一下你的能力" }',
      '    ]',
      '  }\''
    ].join('\n')
  }
  return [
    `curl ${serviceOrigin.value}/v1/chat/completions \\`,
    '  -H "Authorization: Bearer sk-your-api-key" \\',
    '  -H "content-type: application/json" \\',
    '  -d \'{',
    `    "model": "${row.name}",`,
    '    "messages": [',
    '      { "role": "user", "content": "你好，介绍一下你的能力" }',
    '    ]',
    '  }\''
  ].join('\n')
}

function envExample(row: PublicModelPricing): string {
  if (row.platform === 'anthropic') {
    return [
      `export ANTHROPIC_BASE_URL="${serviceOrigin.value}"`,
      'export ANTHROPIC_AUTH_TOKEN="sk-your-api-key"',
      `export ANTHROPIC_MODEL="${row.name}"`
    ].join('\n')
  }
  return [
    `export OPENAI_BASE_URL="${serviceOrigin.value}/v1"`,
    'export OPENAI_API_KEY="sk-your-api-key"',
    `export OPENAI_MODEL="${row.name}"`
  ].join('\n')
}

function usageDescription(row: PublicModelPricing): string {
  if (row.platform === 'anthropic' && row.group_name.includes('官key')) {
    return t('modelSquare.detail.officialKeyUsage')
  }
  if (row.platform === 'anthropic') {
    return t('modelSquare.detail.anthropicUsage')
  }
  return t('modelSquare.detail.openaiUsage')
}

function billingDescription(row: PublicModelPricing): string {
  const pricing = row.pricing
  if (!pricing) return t('modelSquare.detail.noPricingNote')
  if (pricing.billing_mode === BILLING_MODE_PER_REQUEST) {
    return t('modelSquare.detail.perRequestBilling')
  }
  if (pricing.billing_mode === BILLING_MODE_IMAGE) {
    return t('modelSquare.detail.imageBilling')
  }
  return t('modelSquare.detail.tokenBilling')
}

function formatInterval(interval: UserPricingInterval): string {
  const max = interval.max_tokens == null ? t('modelSquare.detail.noUpperLimit') : interval.max_tokens
  const label = interval.tier_label ? `${interval.tier_label}: ` : ''
  const input = formatScaled(interval.input_price, perMillionScale)
  const output = formatScaled(interval.output_price, perMillionScale)
  const request = formatScaled(interval.per_request_price, 1)
  return `${label}${interval.min_tokens} - ${max}: ${t('modelSquare.input')} ${input}, ${t('modelSquare.output')} ${output}, ${t('modelSquare.detail.singleRequest')} ${request}`
}

function comparablePrice(row: PublicModelPricing): number {
  const pricing = row.pricing
  if (!pricing) return Number.POSITIVE_INFINITY
  if (pricing.billing_mode === BILLING_MODE_TOKEN) return pricing.input_price ?? Number.POSITIVE_INFINITY
  if (pricing.billing_mode === BILLING_MODE_IMAGE) {
    return pricing.image_output_price ?? pricing.per_request_price ?? Number.POSITIVE_INFINITY
  }
  return pricing.per_request_price ?? Number.POSITIVE_INFINITY
}

onMounted(() => {
  initTheme()
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
  loadPublicPricing()
})
</script>
