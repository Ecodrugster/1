<template>
  <div class="min-h-screen bg-slate-950 text-slate-200 flex">
    <AdminSidebar />

    <main class="flex-grow flex flex-col h-screen overflow-hidden">
      <header class="h-16 border-b border-white/5 bg-slate-950/60 backdrop-blur-md flex items-center justify-between px-6 md:px-8">
        <h3 class="text-sm font-medium text-slate-300">
          {{ currentPageTitle }}
        </h3>
        <div class="flex items-center gap-3">
          <NuxtLink
            to="/"
            class="hidden sm:inline-flex text-xs px-3 py-1.5 rounded-lg bg-white/5 text-slate-300 hover:bg-white/10 hover:text-white transition-all"
          >
            На главную
          </NuxtLink>
          <div class="text-right">
            <div class="text-xs font-bold text-white">{{ userDisplayName }}</div>
            <div class="text-[10px] text-blue-500 font-mono uppercase tracking-wider">Admin</div>
          </div>
          <img
            v-if="userStore.user?.photoURL"
            :src="userStore.user.photoURL"
            class="w-8 h-8 rounded-lg border border-white/10"
            alt="avatar"
          />
          <div v-else class="w-8 h-8 rounded-lg bg-slate-800 flex items-center justify-center text-xs font-bold">
            A
          </div>
        </div>
      </header>

      <div class="flex-grow overflow-y-auto p-6 md:p-8 custom-scrollbar">
        <slot />
      </div>
    </main>
  </div>
</template>

<script setup>
const userStore = useUserStore()
const route = useRoute()

const userDisplayName = computed(() => {
  return userStore.user?.displayName || userStore.profile?.display_name || userStore.profile?.displayName || userStore.profile?.email || 'Администратор'
})

const currentPageTitle = computed(() => {
  const titles = {
    '/admin': 'Обзор системы',
    '/admin/users': 'Управление пользователями',
    '/admin/schedule': 'Расписание пар',
    '/admin/news': 'Управление новостями',
    '/admin/clubs': 'Управление клубами',
    '/admin/moderation': 'Модерация контента'
  }
  return titles[route.path] || 'Панель управления'
})
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.08);
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.14);
}
</style>
