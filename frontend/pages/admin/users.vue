<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-2xl font-bold text-white">Пользователи</h2>
      <div class="flex gap-4">
        <input 
          v-model="search" 
          type="text" 
          placeholder="Поиск по имени или email..." 
          class="bg-slate-900 border border-white/10 rounded-xl px-4 py-2 text-sm text-white focus:ring-2 focus:ring-blue-500/50 w-64"
        />
      </div>
    </div>

    <div class="bg-slate-900 border border-white/5 rounded-2xl overflow-hidden shadow-xl">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-white/5 text-[10px] uppercase tracking-wider text-slate-500 font-bold">
            <th class="px-6 py-4">Пользователь</th>
            <th class="px-6 py-4">Email</th>
            <th class="px-6 py-4">Роль</th>
            <th class="px-6 py-4">Действия</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="user in filteredUsers" :key="user.uid" class="hover:bg-white/[0.02] transition-colors">
            <td class="px-6 py-4">
              <div class="flex items-center gap-3">
                <img v-if="user.photo_url" :src="user.photo_url" class="w-8 h-8 rounded-lg object-cover" />
                <div v-else class="w-8 h-8 rounded-lg bg-slate-800 flex items-center justify-center text-xs font-bold text-slate-500">
                  {{ (user.display_name || user.email)[0].toUpperCase() }}
                </div>
                <div class="text-sm font-medium text-white">{{ user.display_name || 'Без имени' }}</div>
              </div>
            </td>
            <td class="px-6 py-4 text-sm text-slate-400 font-mono">{{ user.email }}</td>
            <td class="px-6 py-4">
              <span 
                class="px-2 py-1 rounded-md text-[10px] font-bold uppercase tracking-tight"
                :class="user.role === 'admin' ? 'bg-blue-500/10 text-blue-500' : 'bg-slate-500/10 text-slate-500'"
              >
                {{ user.role || 'student' }}
              </span>
            </td>
            <td class="px-6 py-4">
              <button 
                @click="toggleRole(user)"
                :disabled="updating === user.uid"
                class="text-xs font-semibold px-3 py-1.5 rounded-lg transition-all"
                :class="user.role === 'admin' ? 'bg-red-500/10 text-red-500 hover:bg-red-500/20' : 'bg-blue-600 text-white hover:bg-blue-500'"
              >
                {{ user.role === 'admin' ? 'Снять права' : 'Сделать админом' }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      
      <div v-if="loading" class="p-12 text-center text-slate-500 animate-pulse">
        Загрузка списка пользователей...
      </div>
      <div v-else-if="filteredUsers.length === 0" class="p-12 text-center text-slate-500">
        Пользователи не найдены.
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
const search = ref('')
const users = ref([])
const loading = ref(true)
const updating = ref(null)

const fetchUsers = async () => {
  loading.value = true
  try {
    const data = await api('/admin/users')
    users.value = data.users || []
  } catch (e) {
    console.error('Failed to fetch users:', e)
  } finally {
    loading.value = false
  }
}

const filteredUsers = computed(() => {
  if (!search.value) return users.value
  const s = search.value.toLowerCase()
  return users.value.filter(u => 
    u.display_name?.toLowerCase().includes(s) || 
    u.email?.toLowerCase().includes(s)
  )
})

const toggleRole = async (user) => {
  const newRole = user.role === 'admin' ? 'student' : 'admin'
  if (!confirm(`Вы уверены, что хотите изменить роль для ${user.display_name || user.email} на ${newRole}?`)) return
  
  updating.value = user.uid
  try {
    await api(`/admin/users/${user.uid}/role`, {
      method: 'PUT',
      body: { role: newRole }
    })
    user.role = newRole // Local update
    alert('Роль успешно обновлена')
  } catch (e) {
    alert('Ошибка при смене роли: ' + e.message)
  } finally {
    updating.value = null
  }
}

onMounted(fetchUsers)
</script>
