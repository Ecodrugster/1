<template>
  <aside class="w-64 bg-slate-900 border-r border-white/5 flex flex-col h-full">
    <div class="p-6 border-b border-white/5">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center text-white font-bold">A</div>
        <h2 class="text-xl font-bold text-white">Админка</h2>
      </div>
    </div>
    
    <nav class="flex-grow p-4 space-y-2">
      <NuxtLink 
        v-for="item in menuItems" 
        :key="item.path"
        :to="item.path"
        class="flex items-center gap-3 px-4 py-3 rounded-xl transition-all text-sm font-medium"
        :class="$route.path === item.path ? 'bg-blue-600 text-white shadow-lg shadow-blue-600/20' : 'text-slate-400 hover:bg-white/5 hover:text-white'"
      >
        <span class="text-lg">{{ item.icon }}</span>
        {{ item.label }}
      </NuxtLink>
    </nav>
    
    <div class="p-4 border-t border-white/5">
      <NuxtLink to="/" class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-400 hover:bg-white/5 hover:text-white transition-all text-sm">
        <span>🏠</span> Вернуться на сайт
      </NuxtLink>
    </div>
  </aside>
</template>
  
<script setup>
const menuItems = [
  { label: 'Дашборд', path: '/admin', icon: '📊' },
  { label: 'Пользователи', path: '/admin/users', icon: '👥' },
  { label: 'Новости', path: '/admin/news', icon: '📰' },
  { label: 'Клубы', path: '/admin/clubs', icon: '🏢' },
  { label: 'Модерация', path: '/admin/moderation', icon: '🛡️' },
]
</script>
