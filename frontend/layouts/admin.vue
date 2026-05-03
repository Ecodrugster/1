<template>
  <div class="min-h-screen bg-slate-950 text-slate-200 flex">
    <!-- Sidebar -->
    <AdminSidebar />
    
    <!-- Main Content -->
    <main class="flex-grow flex flex-col h-screen overflow-hidden">
      <!-- Header -->
      <header class="h-16 border-b border-white/5 bg-slate-950/50 backdrop-blur-md flex items-center justify-between px-8">
        <h3 class="text-sm font-medium text-slate-400">
          {{ currentPageTitle }}
        </h3>
        <div class="flex items-center gap-4">
          <div class="text-right">
            <div class="text-xs font-bold text-white">{{ userStore.user?.displayName }}</div>
            <div class="text-[10px] text-blue-500 font-mono uppercase tracking-wider">Administrator</div>
          </div>
          <img :src="userStore.user?.photoURL" v-if="userStore.user?.photoURL" class="w-8 h-8 rounded-lg border border-white/10" />
          <div v-else class="w-8 h-8 rounded-lg bg-slate-800 flex items-center justify-center text-xs font-bold">A</div>
        </div>
      </header>
      
      <!-- Page Body -->
      <div class="flex-grow overflow-y-auto p-8 custom-scrollbar">
        <slot />
      </div>
    </main>
  </div>
</template>

<script setup>
const userStore = useUserStore()
const route = useRoute()

const currentPageTitle = computed(() => {
  const titles = {
    '/admin': 'Обзор системы',
    '/admin/users': 'Управление пользователями',
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
  background: rgba(255, 255, 255, 0.05);
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.1);
}
</style>
