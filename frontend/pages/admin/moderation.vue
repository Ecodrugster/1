<template>
  <div class="space-y-8">
    <section class="space-y-4">
      <div v-if="errorMessage" class="bg-red-500/10 border border-red-500/30 text-red-300 rounded-xl px-4 py-3 text-sm">
        {{ errorMessage }}
      </div>

      <div class="flex justify-between items-center">
        <h2 class="text-2xl font-bold text-white">Заявки на клубы</h2>
        <div class="text-sm text-slate-500">{{ clubRequests.length }} в очереди</div>
      </div>

      <div v-if="loadingClubs" class="space-y-3">
        <div v-for="i in 2" :key="i" class="h-24 bg-slate-900 border border-white/5 rounded-2xl animate-pulse"></div>
      </div>

      <div v-else-if="clubRequests.length === 0" class="bg-slate-900 border border-white/5 p-10 rounded-2xl text-center text-slate-500">
        Нет заявок на модерацию клубов.
      </div>

      <div v-else class="space-y-3">
        <div v-for="club in clubRequests" :key="club.id" class="bg-slate-900 border border-white/5 p-5 rounded-2xl flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
          <div>
            <div class="text-white font-bold">{{ club.name }}</div>
            <div class="text-sm text-slate-400 mt-1">{{ club.description }}</div>
            <div class="text-xs text-slate-500 mt-2">
              Создатель: {{ club.created_by_name || club.created_by || 'не указан' }} • Участников: {{ club.members?.length || 0 }}
            </div>
          </div>

          <div class="flex gap-2">
            <button
              @click="approveClub(club.id)"
              class="px-4 py-2 bg-green-500/10 hover:bg-green-500 text-green-400 hover:text-white rounded-lg text-sm transition-all border border-green-500/20"
            >
              Одобрить
            </button>
            <button
              @click="rejectClub(club.id)"
              class="px-4 py-2 bg-red-500/10 hover:bg-red-500 text-red-400 hover:text-white rounded-lg text-sm transition-all border border-red-500/20"
            >
              Отклонить
            </button>
          </div>
        </div>
      </div>
    </section>

    <section class="space-y-4">
      <div class="flex justify-between items-center">
        <h2 class="text-2xl font-bold text-white">Модерация постов</h2>
        <div class="text-sm text-slate-500">Последние 50 сообщений</div>
      </div>

      <div v-if="loadingPosts" class="space-y-4">
        <div v-for="i in 3" :key="i" class="bg-slate-900 border border-white/5 p-6 rounded-2xl animate-pulse h-32"></div>
      </div>

      <div v-else-if="posts.length === 0" class="bg-slate-900 border border-white/5 p-12 rounded-2xl text-center text-slate-500">
        Нет постов для модерации.
      </div>

      <div v-else class="space-y-4">
        <div v-for="post in posts" :key="post.id" class="bg-slate-900 border border-white/5 p-6 rounded-2xl flex justify-between items-start group hover:border-red-500/30 transition-all">
          <div class="flex gap-4">
            <div class="w-10 h-10 rounded-xl bg-slate-800 flex items-center justify-center text-xs text-slate-500">??</div>
            <div>
              <div class="flex items-center gap-2 mb-1">
                <span class="text-sm font-bold text-white">
                  Автор: {{ post.author_name || post.author_id || 'не указан' }}
                </span>
                <span class="text-[10px] text-slate-500 font-mono">{{ formatDateTime(post.created_at) }}</span>
              </div>
              <p class="text-slate-300 text-sm leading-relaxed">{{ post.content }}</p>
            </div>
          </div>

          <button
            @click="deletePost(post.id)"
            class="opacity-0 group-hover:opacity-100 px-4 py-2 bg-red-500/10 hover:bg-red-500 text-red-500 hover:text-white rounded-xl text-xs font-bold transition-all border border-red-500/20"
          >
            Удалить пост
          </button>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

const { fetchApi: api } = useApi()
const posts = ref([])
const clubRequests = ref([])
const loadingPosts = ref(true)
const loadingClubs = ref(true)
const errorMessage = ref('')

const formatDateTime = (value) => {
  if (!value) return 'время не указано'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'время не указано'
  return date.toLocaleString()
}

const fetchPosts = async () => {
  loadingPosts.value = true
  errorMessage.value = ''
  try {
    const data = await api('/admin/posts')
    posts.value = data || []
  } catch (e) {
    console.error('Failed to fetch posts:', e)
    const status = e?.status || e?.response?.status
    const message = e?.data?.error || e?.message || 'Ошибка загрузки постов'
    errorMessage.value = `Не удалось загрузить посты (${status || 'no-status'}): ${message}`
  } finally {
    loadingPosts.value = false
  }
}

const fetchClubRequests = async () => {
  loadingClubs.value = true
  try {
    const data = await api('/admin/club-requests')
    clubRequests.value = data || []
  } catch (e) {
    console.error('Failed to fetch club requests:', e)
    const status = e?.status || e?.response?.status
    const message = e?.data?.error || e?.message || 'Ошибка загрузки заявок клубов'
    errorMessage.value = `Не удалось загрузить заявки клубов (${status || 'no-status'}): ${message}`
  } finally {
    loadingClubs.value = false
  }
}

const deletePost = async (id) => {
  if (!confirm('Удалить этот пост безвозвратно?')) return

  try {
    await api(`/admin/posts/${id}`, { method: 'DELETE' })
    posts.value = posts.value.filter(p => p.id !== id)
    alert('Пост удален модератором')
  } catch (e) {
    alert('Ошибка при удалении: ' + (e?.data?.error || e.message))
  }
}

const approveClub = async (id) => {
  try {
    await api(`/admin/club-requests/${id}/approve`, { method: 'POST' })
    clubRequests.value = clubRequests.value.filter(c => c.id !== id)
  } catch (e) {
    alert('Ошибка одобрения: ' + (e?.data?.error || e.message))
  }
}

const rejectClub = async (id) => {
  const reason = prompt('Причина отклонения (необязательно):') || ''
  try {
    await api(`/admin/club-requests/${id}/reject`, {
      method: 'POST',
      body: { comment: reason }
    })
    clubRequests.value = clubRequests.value.filter(c => c.id !== id)
  } catch (e) {
    alert('Ошибка отклонения: ' + (e?.data?.error || e.message))
  }
}

onMounted(async () => {
  await Promise.all([fetchPosts(), fetchClubRequests()])
})
</script>

