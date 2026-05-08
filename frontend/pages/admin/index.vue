<template>
  <div class="space-y-8">
    <div v-if="errorMessage" class="bg-red-500/10 border border-red-500/30 text-red-300 rounded-xl px-4 py-3 text-sm">
      {{ errorMessage }}
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-5">
      <div
        v-for="card in statCards"
        :key="card.label"
        class="bg-slate-900 border border-white/5 rounded-2xl p-5"
      >
        <div class="text-xs uppercase tracking-wider text-slate-500 font-semibold">{{ card.label }}</div>
        <div class="text-3xl font-extrabold text-white mt-2">{{ card.value }}</div>
      </div>
    </div>

    <div class="grid grid-cols-1 xl:grid-cols-3 gap-6">
      <section class="xl:col-span-2 bg-slate-900 border border-white/5 rounded-2xl p-6">
        <div class="flex items-center justify-between mb-5">
          <h3 class="text-white font-bold text-lg">Recent Actions</h3>
          <div class="text-xs text-slate-500">Live data</div>
        </div>

        <div v-if="loading" class="space-y-3">
          <div v-for="i in 4" :key="i" class="h-14 rounded-xl bg-slate-800/60 animate-pulse"></div>
        </div>

        <div v-else-if="recentActions.length === 0" class="text-slate-500 text-sm">
          No actions in log yet.
        </div>

        <div v-else class="space-y-3">
          <div
            v-for="item in recentActions"
            :key="item.id"
            class="border border-white/5 rounded-xl px-4 py-3"
          >
            <div class="text-sm text-white">
              <span class="font-semibold">{{ item.actor_name || item.actor_uid || 'Admin' }}</span>
              {{ actionLabel(item.action) }}
              <span class="font-semibold">{{ actionTarget(item) }}</span>
            </div>
            <div class="text-xs text-slate-500 mt-1">{{ formatDateTime(item.created_at) }}</div>
          </div>
        </div>
      </section>

      <section class="bg-slate-900 border border-white/5 rounded-2xl p-6">
        <h3 class="text-white font-bold text-lg mb-5">System Summary</h3>
        <div class="space-y-3 text-sm">
          <div class="flex justify-between">
            <span class="text-slate-400">Students</span>
            <span class="text-white font-semibold">{{ stats.users.total_students }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-slate-400">Teachers</span>
            <span class="text-white font-semibold">{{ stats.users.total_teachers }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-slate-400">Admins</span>
            <span class="text-white font-semibold">{{ stats.users.total_admins }}</span>
          </div>
          <div class="h-px bg-white/5 my-2"></div>
          <div class="flex justify-between">
            <span class="text-slate-400">Posts in 7 days</span>
            <span class="text-white font-semibold">{{ stats.posts.posts_last_7_days }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-slate-400">Pending clubs</span>
            <span class="text-white font-semibold">{{ stats.clubs.pending_club_requests }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-slate-400">News in DB</span>
            <span class="text-white font-semibold">{{ stats.news.total_news }}</span>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

const { fetchApi: api } = useApi()

const loading = ref(true)
const recentActions = ref([])
const errorMessage = ref('')
const stats = ref({
  users: {
    total_users: 0,
    total_students: 0,
    total_teachers: 0,
    total_admins: 0
  },
  posts: {
    total_posts: 0,
    posts_last_7_days: 0
  },
  clubs: {
    active_clubs: 0,
    pending_club_requests: 0
  },
  news: {
    total_news: 0
  }
})

const statCards = computed(() => ([
  { label: 'Total users', value: stats.value.users.total_users },
  { label: 'Total posts', value: stats.value.posts.total_posts },
  { label: 'Active clubs', value: stats.value.clubs.active_clubs },
  { label: 'Club requests', value: stats.value.clubs.pending_club_requests }
]))

const actionLabel = (action) => {
  const dictionary = {
    'user.role.updated': 'changed user role',
    'user.group.updated': 'updated user group',
    'post.deleted': 'deleted a post',
    'news.deleted': 'deleted news',
    'news.updated': 'updated news',
    'club.created': 'created club',
    'club.updated': 'updated club',
    'club.deleted': 'deleted club',
    'club.request.approved': 'approved club request',
    'club.request.rejected': 'rejected club request',
    'schedule.created': 'created schedule pair',
    'schedule.updated': 'updated schedule pair',
    'schedule.deleted': 'deleted schedule pair'
  }
  return dictionary[action] || 'performed action'
}

const actionTarget = (item) => {
  return item.target_name || item.target_id || item.target_type || 'object'
}

const formatDateTime = (value) => {
  if (!value) return 'time unknown'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'time unknown'
  return date.toLocaleString()
}

const loadDashboard = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const data = await api('/admin/dashboard')
    stats.value = data?.stats || stats.value
    recentActions.value = data?.recent_actions || []
  } catch (e) {
    console.error('Failed to load admin dashboard:', e)
    const status = e?.status || e?.response?.status
    const message = e?.data?.error || e?.message || 'Failed to load admin dashboard'
    errorMessage.value = `Could not load dashboard (${status || 'no-status'}): ${message}`
  } finally {
    loading.value = false
  }
}

onMounted(loadDashboard)
</script>

