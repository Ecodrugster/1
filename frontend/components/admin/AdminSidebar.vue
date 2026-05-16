<template>
  <aside class="w-72 bg-slate-900 border-r border-white/5 flex flex-col h-full">
    <div class="p-6 border-b border-white/5">
      <div class="flex items-center justify-between gap-3">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center text-white font-bold">A</div>
          <h2 class="text-lg font-bold text-white">Admin Panel</h2>
        </div>
        <NuxtLink
          to="/"
          class="text-xs px-3 py-1.5 rounded-lg bg-white/5 text-slate-300 hover:bg-white/10 hover:text-white transition-all"
        >
          Back To Site
        </NuxtLink>
      </div>
    </div>

    <nav class="flex-grow p-4 space-y-2">
      <NuxtLink
        v-for="item in menuItems"
        :key="item.path"
        :to="item.path"
        class="flex items-center gap-3 px-4 py-3 rounded-xl transition-all text-sm font-medium"
        :class="isActive(item.path) ? 'bg-blue-600 text-white shadow-lg shadow-blue-600/20' : 'text-slate-400 hover:bg-white/5 hover:text-white'"
      >
        <span class="text-xs w-5 text-center font-bold">{{ item.icon }}</span>
        {{ item.label }}
      </NuxtLink>
    </nav>

    <div class="p-4 border-t border-white/5">
      <NuxtLink to="/" class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-400 hover:bg-white/5 hover:text-white transition-all text-sm">
        <span class="w-5 text-center">&lt;</span>
        Back To Site
      </NuxtLink>
    </div>
  </aside>
</template>

<script setup>
const route = useRoute()

const menuItems = [
  { label: 'Dashboard', path: '/admin', icon: 'DB' },
  { label: 'Users', path: '/admin/users', icon: 'US' },
  { label: 'Schedule', path: '/admin/schedule', icon: 'SC' },
  { label: 'News', path: '/admin/news', icon: 'NW' },
  { label: 'Clubs', path: '/admin/clubs', icon: 'CL' },
  { label: 'Moderation', path: '/admin/moderation', icon: 'MD' }
]

const isActive = (path) => {
  if (path === '/admin') {
    return route.path === '/admin'
  }
  return route.path.startsWith(path)
}
</script>
