<template>
  <div class="space-y-8">
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div v-for="stat in stats" :key="stat.label" class="bg-slate-900 border border-white/5 p-6 rounded-2xl shadow-xl">
        <div class="flex items-center justify-between mb-4">
          <span class="text-2xl">{{ stat.icon }}</span>
          <span class="text-xs font-bold text-blue-500 bg-blue-500/10 px-2 py-1 rounded-full">+12%</span>
        </div>
        <div class="text-2xl font-bold text-white mb-1">{{ stat.value }}</div>
        <div class="text-xs text-slate-500 uppercase tracking-wider font-semibold">{{ stat.label }}</div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <div class="bg-slate-900 border border-white/5 rounded-2xl p-6">
        <h3 class="text-white font-bold mb-6 flex items-center gap-2">
          <span>🕒</span> Последние действия
        </h3>
        <div class="space-y-4">
          <div v-for="i in 5" :key="i" class="flex items-center gap-4 text-sm border-b border-white/5 pb-4 last:border-0">
            <div class="w-8 h-8 rounded-lg bg-slate-800 flex items-center justify-center text-xs">👤</div>
            <div>
              <span class="text-white font-medium">Модератор</span> удалил пост #{{ 2034 + i }}
              <div class="text-[10px] text-slate-500 mt-0.5">2 часа назад</div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="bg-slate-900 border border-white/5 rounded-2xl p-6 flex flex-col items-center justify-center text-center py-12">
        <div class="w-16 h-16 bg-blue-600/10 text-blue-500 rounded-full flex items-center justify-center text-2xl mb-4">🛡️</div>
        <h3 class="text-white font-bold mb-2">Система защищена</h3>
        <p class="text-slate-500 text-sm max-w-xs">Все модули работают в штатном режиме. Подозрительной активности не обнаружено.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

const stats = [
  { label: 'Всего студентов', value: '1,284', icon: '👥' },
  { label: 'Постов за неделю', value: '342', icon: '📝' },
  { label: 'Активных клубов', value: '18', icon: '🏢' }
]
</script>
