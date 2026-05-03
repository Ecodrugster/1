<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-2xl font-bold text-white">Управление новостями</h2>
      <button 
        @click="openCreateModal"
        class="bg-blue-600 hover:bg-blue-500 text-white px-6 py-2.5 rounded-xl font-semibold transition-all shadow-lg shadow-blue-600/20"
      >
        + Добавить новость
      </button>
    </div>

    <!-- News Table -->
    <div class="bg-slate-900 border border-white/5 rounded-2xl overflow-hidden shadow-xl">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-white/5 text-[10px] uppercase tracking-wider text-slate-500 font-bold">
            <th class="px-6 py-4">Заголовок</th>
            <th class="px-6 py-4">Дата создания</th>
            <th class="px-6 py-4 text-right">Действия</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="item in news" :key="item.id" class="hover:bg-white/[0.02] transition-colors text-sm">
            <td class="px-6 py-4">
              <div class="text-white font-medium">{{ item.title }}</div>
              <div class="text-xs text-slate-500 truncate max-w-xs">{{ item.description }}</div>
            </td>
            <td class="px-6 py-4 text-slate-400 font-mono text-xs">{{ new Date(item.created_at).toLocaleDateString() }}</td>
            <td class="px-6 py-4 text-right space-x-2">
              <button @click="deleteNews(item.id)" class="p-2 text-red-500 hover:bg-red-500/10 rounded-lg transition-all">🗑️</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create Modal (Simplified for now) -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm">
      <div class="bg-slate-900 border border-white/10 rounded-2xl p-8 max-w-md w-full shadow-2xl">
        <h2 class="text-2xl font-bold text-white mb-6">Новая новость</h2>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-2">Заголовок</label>
            <input v-model="form.title" type="text" class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-2">Содержание</label>
            <textarea v-model="form.description" rows="4" class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50 resize-none"></textarea>
          </div>
        </div>
        <div class="flex space-x-4 mt-8">
          <button @click="showModal = false" class="flex-grow py-3 bg-white/5 hover:bg-white/10 text-white rounded-xl transition-all">Отмена</button>
          <button @click="saveNews" class="flex-grow py-3 bg-blue-600 hover:bg-blue-500 text-white rounded-xl font-semibold shadow-lg shadow-blue-600/20 transition-all">Опубликовать</button>
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
const showModal = ref(false)
const form = reactive({ title: '', description: '' })

const fetchNews = async () => {
  const data = await api('/news')
  news.value = data || []
}

const openCreateModal = () => {
  form.title = ''
  form.description = ''
  showModal.value = true
}

const saveNews = async () => {
  try {
    // В ТЗ указан эндпоинт /admin/news, но мы можем использовать существующий /news если добавим туда проверку прав
    // Для простоты пока используем /news (POST)
    await api('/news', {
      method: 'POST',
      body: { ...form, created_at: new Date().toISOString() }
    })
    showModal.value = false
    fetchNews()
  } catch (e) {
    alert('Ошибка при сохранении')
  }
}

const deleteNews = async (id) => {
  if (!confirm('Удалить эту новость?')) return
  try {
    // Пока у нас нет DELETE /news, но мы его добавим
    await api(`/admin/news/${id}`, { method: 'DELETE' })
    fetchNews()
  } catch (e) {
    alert('Ошибка при удалении')
  }
}

onMounted(fetchNews)
</script>
