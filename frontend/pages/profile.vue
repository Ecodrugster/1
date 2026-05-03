<template>
  <div class="max-w-4xl mx-auto py-8 px-4">
    <ClientOnly>
      <div v-if="userStore.user" class="bg-slate-900 border border-white/5 rounded-2xl shadow-xl overflow-hidden">
        <!-- Cover -->
        <div class="h-48 bg-gradient-to-r from-blue-600 to-purple-600"></div>
        
        <!-- Profile Info -->
        <div class="px-8 pb-8">
          <div class="relative -mt-16 mb-6">
            <img 
              v-if="userStore.user?.photoURL"
              :src="userStore.user.photoURL" 
              class="w-32 h-32 rounded-2xl bg-slate-800 border-4 border-slate-900 shadow-2xl object-cover"
            />
            <div v-else class="w-32 h-32 rounded-2xl bg-slate-800 border-4 border-slate-900 shadow-2xl flex items-center justify-center text-4xl font-bold text-slate-500">
              {{ (userStore.user?.displayName || userStore.user?.email || 'U')[0].toUpperCase() }}
            </div>
            <button class="absolute bottom-0 right-0 p-2 bg-blue-600 rounded-lg text-white shadow-lg hover:bg-blue-500 transition-all">📸</button>
          </div>
          
          <div class="flex justify-between items-start">
            <div>
              <h1 class="text-3xl font-bold text-white mb-1">
                {{ userStore.user?.displayName || 'Пользователь' }}
              </h1>
              <p class="text-slate-400">
                {{ userStore.user?.email }} • Студент ITSTEP
              </p>
            </div>
            <button 
              @click="openEditModal"
              class="px-6 py-2 bg-white/5 hover:bg-white/10 border border-white/10 rounded-lg text-sm font-medium transition-all"
            >
              Редактировать профиль
            </button>
          </div>

          <!-- Edit Profile Modal -->
          <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm">
            <div class="bg-slate-900 border border-white/10 rounded-2xl p-8 max-w-md w-full shadow-2xl">
              <h2 class="text-2xl font-bold text-white mb-6">Настройки профиля</h2>
              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-medium text-slate-400 mb-2">Имя</label>
                  <input v-model="editForm.displayName" type="text" class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50" />
                </div>
                <div>
                  <label class="block text-sm font-medium text-slate-400 mb-2">Ссылка на фото (URL)</label>
                  <input v-model="editForm.photoURL" type="text" class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50" />
                </div>
              </div>
              <div class="flex space-x-4 mt-8">
                <button @click="showEditModal = false" class="flex-grow py-3 bg-white/5 hover:bg-white/10 text-white rounded-xl transition-all">Отмена</button>
                <button @click="saveProfile" class="flex-grow py-3 bg-blue-600 hover:bg-blue-500 text-white rounded-xl font-semibold shadow-lg shadow-blue-600/20 transition-all">Сохранить</button>
              </div>
            </div>
          </div>

          <div class="mt-8 grid grid-cols-1 md:grid-cols-2 gap-8">
            <div class="space-y-6">
              <div>
                <h3 class="text-slate-300 font-semibold mb-3">О себе</h3>
                <p class="text-slate-400 text-sm leading-relaxed">
                  Это ваш личный профиль в соцсети ITSTEP. Здесь вы можете делиться своими интересами и следить за новостями колледжа.
                </p>
              </div>
              <div>
                <h3 class="text-slate-300 font-semibold mb-3">Интересы</h3>
                <div class="flex flex-wrap gap-2">
                  <span v-for="tag in ['Vue.js', 'Go', 'Nuxt', 'PostgreSQL', 'UI/UX']" :key="tag" class="px-3 py-1 bg-blue-600/10 text-blue-400 rounded-full text-xs">
                    {{ tag }}
                  </span>
                </div>
              </div>
            </div>

            <div class="bg-white/5 rounded-xl p-6 border border-white/5">
              <h3 class="text-slate-300 font-semibold mb-4">Статистика</h3>
              <div class="grid grid-cols-2 gap-4">
                <div class="text-center p-4 bg-slate-950/50 rounded-lg">
                  <div class="text-2xl font-bold text-white">{{ stats.posts }}</div>
                  <div class="text-xs text-slate-500 uppercase">Постов</div>
                </div>
                <div class="text-center p-4 bg-slate-950/50 rounded-lg">
                  <div class="text-2xl font-bold text-white">{{ stats.comments }}</div>
                  <div class="text-xs text-slate-500 uppercase">Комментариев</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="text-center py-20 text-slate-500">
        Загрузка профиля...
      </div>
    </ClientOnly>
  </div>
</template>

<script setup>
import { updateProfile } from 'firebase/auth'

definePageMeta({
  middleware: 'auth'
})

const { $auth } = useNuxtApp()
const userStore = useUserStore()
const { fetchApi: api } = useApi()

const showEditModal = ref(false)
const editForm = reactive({
  displayName: '',
  photoURL: ''
})

const stats = ref({ posts: 0, comments: 0 })

const fetchStats = async () => {
  try {
    const data = await api('/profile/stats')
    stats.value = data || { posts: 0, comments: 0 }
  } catch (e) {
    console.error('Failed to fetch stats:', e)
  }
}

const openEditModal = () => {
  editForm.displayName = userStore.user?.displayName || ''
  editForm.photoURL = userStore.user?.photoURL || ''
  showEditModal.value = true
}

const saveProfile = async () => {
  try {
    await updateProfile($auth.currentUser, {
      displayName: editForm.displayName,
      photoURL: editForm.photoURL
    })
    
    await api('/profile', {
      method: 'PUT',
      body: editForm
    })
    
    userStore.setUser({ ...$auth.currentUser })
    showEditModal.value = false
    alert('Профиль обновлен!')
  } catch (e) {
    alert('Ошибка при обновлении профиля: ' + e.message)
  }
}

onMounted(() => {
  fetchStats()
})
</script>
