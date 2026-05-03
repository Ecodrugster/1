<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-2xl font-bold text-white">Модерация постов</h2>
      <div class="text-sm text-slate-500">Последние 50 сообщений</div>
    </div>

    <div v-if="loading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="bg-slate-900 border border-white/5 p-6 rounded-2xl animate-pulse h-32"></div>
    </div>

    <div v-else-if="posts.length === 0" class="bg-slate-900 border border-white/5 p-12 rounded-2xl text-center text-slate-500">
      Нет постов для модерации.
    </div>

    <div v-else class="space-y-4">
      <div v-for="post in posts" :key="post.id" class="bg-slate-900 border border-white/5 p-6 rounded-2xl flex justify-between items-start group hover:border-red-500/30 transition-all">
        <div class="flex gap-4">
          <div class="w-10 h-10 rounded-xl bg-slate-800 flex items-center justify-center text-xs text-slate-500">👤</div>
          <div>
            <div class="flex items-center gap-2 mb-1">
              <span class="text-sm font-bold text-white">ID автора: {{ post.author_id.substring(0, 8) }}...</span>
              <span class="text-[10px] text-slate-500 font-mono">{{ new Date(post.created_at).toLocaleString() }}</span>
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
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

const { fetchApi: api } = useApi()
const posts = ref([])
const loading = ref(true)

const fetchPosts = async () => {
  loading.value = true
  try {
    const data = await api('/admin/posts')
    posts.value = data || []
  } catch (e) {
    console.error('Failed to fetch posts:', e)
  } finally {
    loading.value = false
  }
}

const deletePost = async (id) => {
  if (!confirm('Вы уверены, что хотите безвозвратно удалить этот пост?')) return
  
  try {
    await api(`/admin/posts/${id}`, { method: 'DELETE' })
    posts.value = posts.value.filter(p => p.id !== id)
    alert('Пост удален модератором')
  } catch (e) {
    alert('Ошибка при удалении: ' + e.message)
  }
}

onMounted(fetchPosts)
</script>
