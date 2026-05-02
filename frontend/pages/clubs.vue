<template>
  <div class="max-w-6xl mx-auto py-8 px-4">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold text-white">Студенческие клубы</h1>
      <button 
        @click="createSampleClub"
        class="bg-white/5 hover:bg-white/10 text-white px-4 py-2 rounded-lg font-medium transition-all text-sm border border-white/10"
      >
        Создать тестовый клуб
      </button>
    </div>

    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div v-for="i in 3" :key="i" class="bg-slate-900 border border-white/5 rounded-xl p-6 shadow-xl animate-pulse">
        <div class="w-full h-32 bg-slate-800 rounded-lg mb-4"></div>
        <div class="h-5 w-3/4 bg-slate-800 rounded mb-2"></div>
        <div class="h-4 w-full bg-slate-800 rounded mb-4"></div>
        <div class="w-full h-10 bg-slate-800 rounded-lg"></div>
      </div>
    </div>

    <div v-else-if="clubs.length === 0" class="text-center py-20 bg-slate-900/50 rounded-xl border border-white/5">
      <div class="text-6xl mb-4">🎭</div>
      <h3 class="text-xl text-white font-medium mb-2">Пока нет активных клубов</h3>
      <p class="text-slate-500">Нажмите "Создать тестовый клуб", чтобы добавить первый.</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div v-for="club in clubs" :key="club.id" class="bg-slate-900 border border-white/5 rounded-xl p-6 shadow-xl flex flex-col">
        <div class="w-full h-32 rounded-lg mb-4 flex items-center justify-center text-4xl" :class="club.color || 'bg-blue-600/20'">
          {{ club.icon || '🎭' }}
        </div>
        <h3 class="text-white font-semibold mb-2 text-lg">{{ club.name }}</h3>
        <p class="text-slate-400 text-sm mb-4 flex-grow">{{ club.description }}</p>
        
        <div class="flex items-center justify-between mb-4 text-xs text-slate-500 font-medium">
          <span>Участников: {{ club.members?.length || 0 }}</span>
        </div>

        <button 
          v-if="!isMember(club)"
          @click="joinClub(club)"
          :disabled="joining === club.id"
          class="w-full py-2.5 bg-blue-600 hover:bg-blue-500 text-white rounded-lg text-sm font-medium transition-all disabled:opacity-50"
        >
          {{ joining === club.id ? 'Вступаем...' : 'Вступить' }}
        </button>
        <button 
          v-else
          disabled
          class="w-full py-2.5 bg-green-500/10 border border-green-500/20 text-green-400 rounded-lg text-sm font-medium"
        >
          Вы участник ✓
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: 'auth'
})

const { fetchApi: api } = useApi()
const userStore = useUserStore()

const clubs = ref([])
const loading = ref(true)
const joining = ref(null)

const fetchClubs = async () => {
  loading.value = true
  try {
    const data = await api('/clubs')
    clubs.value = data || []
  } catch (e) {
    console.error('Failed to fetch clubs:', e)
  } finally {
    loading.value = false
  }
}

const isMember = (club) => {
  if (!club.members) return false
  return club.members.includes(userStore.user?.uid)
}

const joinClub = async (club) => {
  joining.value = club.id
  try {
    await api(`/clubs/${club.id}/join`, { method: 'POST' })
    // Optimistic update
    if (!club.members) club.members = []
    club.members.push(userStore.user.uid)
  } catch (e) {
    alert('Ошибка при вступлении: ' + e.message)
  } finally {
    joining.value = null
  }
}

const createSampleClub = async () => {
  const sampleClubs = [
    { name: 'Клуб Веб-разработки', description: 'Изучаем Vue, React и создание крутых интерфейсов.', icon: '💻', color: 'bg-purple-600/20' },
    { name: 'Go Backend Team', description: 'Пишем быстрые серверы на Go, изучаем микросервисы.', icon: '🐹', color: 'bg-cyan-600/20' },
    { name: 'Киберспорт ITSTEP', description: 'Тренируемся и участвуем в турнирах по CS:GO и Dota 2.', icon: '🎮', color: 'bg-red-600/20' }
  ]
  
  const randomClub = sampleClubs[Math.floor(Math.random() * sampleClubs.length)]
  
  try {
    await api('/clubs', {
      method: 'POST',
      body: randomClub
    })
    fetchClubs()
  } catch (e) {
    alert('Ошибка создания клуба')
  }
}

onMounted(() => {
  fetchClubs()
})
</script>
