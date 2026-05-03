<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-2xl font-bold text-white">Управление клубами</h2>
      <button 
        @click="openModal()"
        class="bg-blue-600 hover:bg-blue-500 text-white px-6 py-2.5 rounded-xl font-semibold transition-all shadow-lg shadow-blue-600/20"
      >
        + Создать новый клуб
      </button>
    </div>

    <!-- Clubs Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div v-for="club in clubs" :key="club.id" class="bg-slate-900 border border-white/5 p-6 rounded-2xl shadow-xl flex gap-6 group">
        <div class="w-24 h-24 rounded-xl flex items-center justify-center text-4xl" :class="club.color || 'bg-blue-600/10'">
          {{ club.icon || '🏢' }}
        </div>
        <div class="flex-grow">
          <div class="flex justify-between items-start mb-2">
            <h3 class="text-white font-bold">{{ club.name }}</h3>
            <div class="flex gap-2">
              <button @click="openModal(club)" class="text-slate-500 hover:text-white transition-all text-sm">✏️</button>
              <button @click="deleteClub(club.id)" class="text-slate-500 hover:text-red-500 transition-all text-sm">🗑️</button>
            </div>
          </div>
          <p class="text-slate-400 text-xs mb-4 line-clamp-2">{{ club.description }}</p>
          <div class="text-[10px] text-slate-500 uppercase font-bold tracking-widest">
            {{ club.members?.length || 0 }} Участников
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm">
      <div class="bg-slate-900 border border-white/10 rounded-2xl p-8 max-w-md w-full shadow-2xl">
        <h2 class="text-2xl font-bold text-white mb-6">{{ isEditing ? 'Редактировать клуб' : 'Новый клуб' }}</h2>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-2">Название</label>
            <input v-model="form.name" type="text" class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-2">Описание</label>
            <textarea v-model="form.description" rows="3" class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50 resize-none"></textarea>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-slate-400 mb-2">Иконка</label>
              <input v-model="form.icon" type="text" placeholder="🎭" class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50 text-center" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-400 mb-2">Цвет (класс Tailwind)</label>
              <input v-model="form.color" type="text" placeholder="bg-blue-600/20" class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50 text-xs" />
            </div>
          </div>
        </div>
        <div class="flex space-x-4 mt-8">
          <button @click="showModal = false" class="flex-grow py-3 bg-white/5 hover:bg-white/10 text-white rounded-xl transition-all">Отмена</button>
          <button @click="saveClub" class="flex-grow py-3 bg-blue-600 hover:bg-blue-500 text-white rounded-xl font-semibold shadow-lg shadow-blue-600/20 transition-all">Сохранить</button>
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
const clubs = ref([])
const showModal = ref(false)
const isEditing = ref(false)
const currentId = ref(null)

const form = reactive({
  name: '',
  description: '',
  icon: '🎭',
  color: 'bg-blue-600/20'
})

const fetchClubs = async () => {
  const data = await api('/clubs')
  clubs.value = data || []
}

const openModal = (club = null) => {
  if (club) {
    isEditing.value = true
    currentId.value = club.id
    Object.assign(form, {
      name: club.name,
      description: club.description,
      icon: club.icon,
      color: club.color
    })
  } else {
    isEditing.value = false
    currentId.value = null
    Object.assign(form, {
      name: '',
      description: '',
      icon: '🎭',
      color: 'bg-blue-600/20'
    })
  }
  showModal.value = true
}

const saveClub = async () => {
  try {
    const method = isEditing.value ? 'PUT' : 'POST'
    const url = isEditing.value ? `/admin/clubs/${currentId.value}` : '/clubs'
    
    await api(url, {
      method,
      body: form
    })
    
    showModal.value = false
    fetchClubs()
  } catch (e) {
    alert('Ошибка: ' + e.message)
  }
}

const deleteClub = async (id) => {
  if (!confirm('Вы уверены? Это удалит клуб для всех участников.')) return
  try {
    await api(`/admin/clubs/${id}`, { method: 'DELETE' })
    fetchClubs()
  } catch (e) {
    alert('Ошибка удаления')
  }
}

onMounted(fetchClubs)
</script>
