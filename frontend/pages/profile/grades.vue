<template>
  <div class="max-w-4xl mx-auto py-8 px-4">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-white">Мои оценки</h1>
      <p class="text-slate-400">Успеваемость в ITSTEP</p>
    </div>

    <div v-if="loading" class="space-y-4">
      <div v-for="i in 5" :key="i" class="h-20 bg-slate-900 animate-pulse rounded-2xl border border-white/5"></div>
    </div>

    <div v-else class="space-y-4">
      <div v-for="grade in grades" :key="grade.id" class="bg-slate-900 border border-white/5 p-6 rounded-2xl shadow-xl flex justify-between items-center group hover:border-blue-500/30 transition-all">
        <div class="flex items-center space-x-6">
          <div class="w-14 h-14 rounded-xl bg-blue-600/10 flex items-center justify-center text-2xl font-bold text-blue-500">
            {{ grade.value }}
          </div>
          <div>
            <h3 class="text-lg font-bold text-white">{{ grade.subject }}</h3>
            <div class="flex items-center space-x-2 mt-1">
              <span class="text-xs text-slate-500">{{ formatDate(grade.created_at) }}</span>
              <span class="text-xs text-slate-600">•</span>
              <span class="text-xs text-slate-400 italic">"{{ grade.comment || 'Без комментария' }}"</span>
            </div>
          </div>
        </div>
        
        <div class="text-right hidden md:block">
          <div class="text-[10px] text-slate-500 uppercase font-bold tracking-widest mb-1">Выставил(а)</div>
          <div class="text-xs text-slate-300 font-medium">Преподаватель</div>
        </div>
      </div>

      <div v-if="grades.length === 0" class="text-center py-20 bg-slate-900 rounded-2xl border border-dashed border-white/10">
        <div class="text-5xl mb-4">📝</div>
        <p class="text-slate-500">У вас пока нет выставленных оценок</p>
      </div>
    </div>
  </div>
</template>

<script setup>
const { fetchApi: api } = useApi()
const grades = ref([])
const loading = ref(true)

const fetchGrades = async () => {
  try {
    const data = await api('/grades')
    grades.value = data || []
  } catch (e) {
    console.error('Failed to fetch grades:', e)
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return d.toLocaleDateString('ru-RU', { 
    day: 'numeric', 
    month: 'long', 
    year: 'numeric' 
  })
}

onMounted(fetchGrades)
</script>
