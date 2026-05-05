<script setup>
const userStore = useUserStore()
const notificationStore = useNotificationStore()
const { logout } = useAuth()
const router = useRouter()

const searchQuery = ref('')
const handleSearch = () => {
  if (!searchQuery.value.trim()) return
  router.push({ path: '/search', query: { q: searchQuery.value } })
  searchQuery.value = ''
}

const handleLogout = async () => {
  await logout()
  router.push('/login')
}
</script>

<template>
  <div class="min-h-screen bg-slate-950 text-slate-200 selection:bg-blue-500/30">
    <header class="sticky top-0 z-50 bg-slate-950/80 backdrop-blur-md border-b border-white/5">
      <div class="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between">
        <NuxtLink to="/" class="flex items-center space-x-4 cursor-pointer">
          <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center font-bold text-white shadow-lg shadow-blue-600/30">
            S
          </div>
          <span class="text-lg font-bold bg-gradient-to-r from-white to-slate-400 bg-clip-text text-transparent">
            ITSTEP Social
          </span>
        </NuxtLink>
        
        <div class="hidden md:flex items-center bg-slate-900 border border-white/5 rounded-full px-4 py-1.5 w-96">
          <span class="text-slate-500 mr-2">🔍</span>
          <input 
            v-model="searchQuery"
            @keyup.enter="handleSearch"
            type="text" 
            placeholder="Поиск..." 
            class="bg-transparent border-none focus:outline-none text-sm w-full placeholder-slate-600 text-white"
          />
        </div>

        <div class="flex items-center space-x-6">
          <NuxtLink to="/" class="text-sm font-medium text-slate-400 hover:text-white transition-all">Лента</NuxtLink>
          <NuxtLink to="/clubs" class="text-sm font-medium text-slate-400 hover:text-white transition-all">Клубы</NuxtLink>
          <NuxtLink to="/guide" class="text-sm font-bold text-yellow-500 hover:text-yellow-400 transition-all flex items-center">
            <span class="mr-1">💡</span> Навигатор
          </NuxtLink>
          <NuxtLink to="/chat" class="text-sm font-medium text-slate-400 hover:text-white transition-all">Чат</NuxtLink>
          
          <!-- Role specific links -->
          <NuxtLink v-if="userStore.isTeacher" to="/teacher/grades" class="text-sm font-bold text-blue-500 hover:text-blue-400 transition-all flex items-center">
            <span class="mr-1">🎓</span> Журнал
          </NuxtLink>
          <NuxtLink v-if="!userStore.isAdmin && !userStore.isTeacher" to="/profile/grades" class="text-sm font-bold text-green-500 hover:text-green-400 transition-all flex items-center">
            <span class="mr-1">📝</span> Мои оценки
          </NuxtLink>
          <NuxtLink v-if="userStore.isAdmin" to="/admin" class="text-sm font-bold text-red-500 hover:text-red-400 transition-all flex items-center">
            <span class="mr-1">🛡️</span> Админ
          </NuxtLink>
        </div>
          <ClientOnly>
            <button class="p-2 hover:bg-white/5 rounded-full transition-all relative">
              <span class="text-xl">🔔</span>
              <span v-if="notificationStore.unreadCount > 0" class="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full border border-slate-950"></span>
            </button>
            
            <div v-if="userStore.user" class="flex items-center space-x-3">
              <div class="hidden md:block text-right">
                <div class="text-xs font-semibold text-white truncate w-24">
                  {{ userStore.user.displayName || userStore.user.email }}
                </div>
                <div class="text-[10px] text-slate-500 uppercase tracking-wider">Студент</div>
              </div>
              <NuxtLink to="/profile" class="relative group">
                <img 
                  v-if="userStore.user.photoURL"
                  :src="userStore.user.photoURL" 
                  class="w-9 h-9 rounded-full border border-white/10 group-hover:border-blue-500/50 transition-all shadow-lg"
                />
                <div v-else class="w-9 h-9 rounded-full bg-slate-800 border border-white/10 flex items-center justify-center text-xs group-hover:border-blue-500/50">
                  {{ (userStore.user.displayName || userStore.user.email || 'U')[0].toUpperCase() }}
                </div>
              </NuxtLink>
            </div>
          </ClientOnly>
      </div>
    </header>

    <main>
      <slot />
    </main>
  </div>
</template>
