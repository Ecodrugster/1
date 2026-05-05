<template>
  <div class="max-w-6xl mx-auto py-8 px-4">
    <div class="mb-12">
      <div class="relative max-w-2xl mx-auto">
        <span class="absolute left-4 top-1/2 -translate-y-1/2 text-2xl">🔍</span>
        <input 
          v-model="query" 
          @input="handleSearch"
          type="text" 
          placeholder="Искать студентов, посты, клубы или новости..." 
          class="w-full bg-slate-900 border border-white/10 rounded-2xl pl-14 pr-6 py-4 text-xl text-white focus:ring-4 focus:ring-blue-500/20 transition-all outline-none"
        />
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-20">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
    </div>

    <div v-else-if="!query && !hasResults" class="text-center py-20 text-slate-500">
      <div class="text-6xl mb-4">🕵️‍♂️</div>
      <p>Введите поисковый запрос, чтобы начать поиск</p>
    </div>

    <div v-else class="space-y-12">
      <!-- Users -->
      <section v-if="results.users.length > 0">
        <h2 class="text-xl font-bold text-white mb-6 flex items-center">
          <span class="mr-2">👥</span> Студенты и сотрудники
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div v-for="user in results.users" :key="user.uid" class="bg-slate-900 border border-white/5 p-4 rounded-2xl flex items-center space-x-4">
            <div class="w-12 h-12 rounded-full bg-blue-600 flex items-center justify-center font-bold text-white">
              {{ (user.display_name || user.email)[0].toUpperCase() }}
            </div>
            <div>
              <div class="text-white font-bold">{{ user.display_name || 'Без имени' }}</div>
              <div class="text-xs text-slate-500 uppercase">{{ user.role || 'student' }}</div>
            </div>
            <NuxtLink :to="`/chat?uid=${user.uid}`" class="ml-auto text-blue-500 hover:text-blue-400">💬</NuxtLink>
          </div>
        </div>
      </section>

      <!-- Clubs -->
      <section v-if="results.clubs.length > 0">
        <h2 class="text-xl font-bold text-white mb-6 flex items-center">
          <span class="mr-2">🎭</span> Клубы
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div v-for="club in results.clubs" :key="club.id" class="bg-slate-900 border border-white/5 p-6 rounded-2xl flex items-center space-x-6">
            <div class="text-4xl">{{ club.icon || '🏢' }}</div>
            <div>
              <div class="text-white font-bold">{{ club.name }}</div>
              <p class="text-slate-500 text-sm line-clamp-1">{{ club.description }}</p>
            </div>
            <NuxtLink to="/clubs" class="ml-auto text-slate-400 hover:text-white">→</NuxtLink>
          </div>
        </div>
      </section>

      <!-- Posts -->
      <section v-if="results.posts.length > 0">
        <h2 class="text-xl font-bold text-white mb-6 flex items-center">
          <span class="mr-2">📝</span> Посты
        </h2>
        <div class="space-y-4">
          <div v-for="post in results.posts" :key="post.id" class="bg-slate-900 border border-white/5 p-6 rounded-2xl">
            <p class="text-slate-300 mb-2">{{ post.text }}</p>
            <div class="text-[10px] text-slate-600 uppercase font-bold">{{ formatDate(post.created_at) }}</div>
          </div>
        </div>
      </section>

      <!-- News -->
      <section v-if="results.news.length > 0">
        <h2 class="text-xl font-bold text-white mb-6 flex items-center">
          <span class="mr-2">📰</span> Новости
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div v-for="n in results.news" :key="n.id" class="bg-slate-900 border border-white/5 p-6 rounded-2xl">
            <div class="text-[10px] text-blue-500 font-bold uppercase mb-2">{{ n.category || 'Новость' }}</div>
            <h3 class="text-white font-bold mb-2">{{ n.title }}</h3>
            <p class="text-slate-500 text-sm line-clamp-2">{{ n.description }}</p>
          </div>
        </div>
      </section>

      <div v-if="!loading && query && !hasResults" class="text-center py-20 text-slate-500">
        <p>По вашему запросу «{{ query }}» ничего не найдено</p>
      </div>
    </div>
  </div>
</template>

<script setup>
const { fetchApi: api } = useApi()
const route = useRoute()
const query = ref(route.query.q || '')
const loading = ref(false)

const results = reactive({
  users: [],
  posts: [],
  clubs: [],
  news: []
})

const hasResults = computed(() => {
  return results.users.length > 0 || results.posts.length > 0 || results.clubs.length > 0 || results.news.length > 0
})

const handleSearch = async () => {
  if (!query.value || query.value.length < 2) {
    results.users = []
    results.posts = []
    results.clubs = []
    results.news = []
    return
  }

  loading.value = true
  try {
    // В идеале это должен быть один эндпоинт на бэкенде, но для гибкости сделаем параллельные запросы
    const q = query.value.toLowerCase()
    
    // Получаем все данные
    const [usersRes, posts, clubs, news] = await Promise.all([
      api('/users'),
      api('/posts'),
      api('/clubs'),
      api('/news')
    ])

    // /users возвращает { users: [], total: X }
    const usersList = usersRes.users || usersRes || []

    results.users = usersList.filter(u => 
      u.display_name?.toLowerCase().includes(q) || 
      u.email?.toLowerCase().includes(q)
    )

    results.posts = (posts || []).filter(p => 
      p.text?.toLowerCase().includes(q)
    )

    results.clubs = (clubs || []).filter(c => 
      c.name?.toLowerCase().includes(q) || 
      c.description?.toLowerCase().includes(q)
    )

    results.news = (news || []).filter(n => 
      n.title?.toLowerCase().includes(q) || 
      n.description?.toLowerCase().includes(q)
    )
  } catch (e) {
    console.error('Search failed:', e)
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleDateString()
}

// Первоначальный поиск если есть q в URL
onMounted(() => {
  if (query.value) handleSearch()
})
</script>
