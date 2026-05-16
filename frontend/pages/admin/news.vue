<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-2xl font-bold text-white">News Management</h2>
      <button
        @click="openCreateModal"
        class="bg-blue-600 hover:bg-blue-500 text-white px-6 py-2.5 rounded-xl font-semibold transition-all shadow-lg shadow-blue-600/20"
      >
        + Add News
      </button>
    </div>

    <div v-if="errorMessage" class="bg-red-500/10 border border-red-500/30 text-red-300 rounded-xl px-4 py-3 text-sm">
      {{ errorMessage }}
    </div>

    <div class="bg-slate-900 border border-white/5 rounded-2xl overflow-hidden shadow-xl">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-white/5 text-[10px] uppercase tracking-wider text-slate-500 font-bold">
            <th class="px-6 py-4">Title</th>
            <th class="px-6 py-4">Category</th>
            <th class="px-6 py-4">Created</th>
            <th class="px-6 py-4 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="item in news" :key="item.id" class="hover:bg-white/[0.02] transition-colors text-sm">
            <td class="px-6 py-4">
              <div class="text-white font-medium">{{ item.title }}</div>
              <div class="text-xs text-slate-500 truncate max-w-xs">{{ item.description }}</div>
            </td>
            <td class="px-6 py-4">
              <span class="px-2 py-1 bg-blue-500/10 text-blue-500 text-[10px] font-bold uppercase rounded-md">
                {{ item.category || 'news' }}
              </span>
            </td>
            <td class="px-6 py-4 text-slate-400 font-mono text-xs">
              {{ formatDate(item.created_at) }}
            </td>
            <td class="px-6 py-4 text-right space-x-2">
              <button
                @click="deleteNews(item.id)"
                class="px-3 py-2 text-red-500 hover:bg-red-500/10 rounded-lg transition-all text-xs font-semibold"
              >
                Delete
              </button>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-if="loading" class="p-12 text-center text-slate-500 animate-pulse">Loading news...</div>
      <div v-else-if="news.length === 0" class="p-12 text-center text-slate-500">No news yet.</div>
    </div>

    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm">
      <div class="bg-slate-900 border border-white/10 rounded-2xl p-8 max-w-md w-full shadow-2xl">
        <h2 class="text-2xl font-bold text-white mb-6">New News Item</h2>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-2">Title</label>
            <input v-model="form.title" type="text" class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-2">Category</label>
            <select v-model="form.category" class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50">
              <option value="news">News</option>
              <option value="announcement">Announcement</option>
              <option value="event">Event</option>
              <option value="deadline">Deadline</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-2">Description</label>
            <textarea v-model="form.description" rows="4" class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50 resize-none"></textarea>
          </div>
        </div>
        <div class="flex space-x-4 mt-8">
          <button @click="showModal = false" class="flex-grow py-3 bg-white/5 hover:bg-white/10 text-white rounded-xl transition-all">Cancel</button>
          <button @click="saveNews" class="flex-grow py-3 bg-blue-600 hover:bg-blue-500 text-white rounded-xl font-semibold shadow-lg shadow-blue-600/20 transition-all">Publish</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

const { fetchApi: api } = useApi()
const news = ref([])
const loading = ref(false)
const showModal = ref(false)
const errorMessage = ref('')
const form = reactive({ title: '', description: '', category: 'news' })

const formatDate = (value) => {
  if (!value) return 'date unknown'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'date unknown'
  return date.toLocaleDateString()
}

const fetchNews = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const data = await api('/news')
    news.value = data || []
  } catch (e) {
    const status = e?.status || e?.response?.status
    const message = e?.data?.error || e?.message || 'Failed to load news'
    errorMessage.value = `Could not load news (${status || 'no-status'}): ${message}`
    news.value = []
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  form.title = ''
  form.description = ''
  form.category = 'news'
  showModal.value = true
}

const saveNews = async () => {
  try {
    await api('/admin/news', {
      method: 'POST',
      body: { ...form }
    })
    showModal.value = false
    await fetchNews()
  } catch (e) {
    alert('Save error: ' + (e?.data?.error || e?.message || 'unknown error'))
  }
}

const deleteNews = async (id) => {
  if (!confirm('Delete this news item?')) return
  try {
    await api(`/admin/news/${id}`, { method: 'DELETE' })
    await fetchNews()
  } catch (e) {
    alert('Delete error: ' + (e?.data?.error || e?.message || 'unknown error'))
  }
}

onMounted(fetchNews)
</script>
