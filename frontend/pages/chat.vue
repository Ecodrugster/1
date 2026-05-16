<template>
  <div class="max-w-6xl mx-auto py-8 px-4 h-[calc(100vh-120px)]">
    <div class="bg-slate-900 border border-white/5 rounded-xl h-full shadow-xl flex overflow-hidden">
      <!-- Chat List -->
      <div class="w-80 border-r border-white/5 flex flex-col">
        <div class="p-4 border-b border-white/5">
          <h3 class="text-white font-semibold mb-2">Студенты</h3>
          <input 
            v-model="searchQuery"
            type="text" 
            placeholder="Поиск..." 
            class="w-full bg-slate-800 border-none rounded-lg px-4 py-2 text-sm text-white focus:ring-1 focus:ring-blue-500/50" 
          />
        </div>
        <div class="flex-grow overflow-y-auto">
          <div 
            v-for="user in filteredUsers" 
            :key="user.uid" 
            @click="selectUser(user)"
            class="p-4 hover:bg-white/5 cursor-pointer border-b border-white/5 last:border-0 transition-all flex items-center space-x-3"
            :class="{'bg-blue-600/10 border-l-2 border-l-blue-600': selectedUser?.uid === user.uid}"
          >
            <img v-if="user.photoURL" :src="user.photoURL" class="w-10 h-10 rounded-full" />
            <div v-else class="w-10 h-10 rounded-full bg-slate-800 flex items-center justify-center text-xs text-slate-500 font-bold">
              {{ (user.displayName || user.email || 'U')[0].toUpperCase() }}
            </div>
            <div class="flex-grow min-w-0">
              <h4 class="text-white text-sm font-medium truncate">{{ user.displayName || user.email }}</h4>
              <p class="text-slate-500 text-[10px] uppercase">Студент</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Chat Window -->
      <div class="flex-grow flex flex-col">
        <template v-if="selectedUser">
          <div class="p-4 border-b border-white/5 flex items-center justify-between bg-slate-900/50">
            <div class="flex items-center space-x-3">
              <img v-if="selectedUser.photoURL" :src="selectedUser.photoURL" class="w-10 h-10 rounded-full" />
              <div v-else class="w-10 h-10 rounded-full bg-slate-800 flex items-center justify-center text-xs font-bold text-white">
                {{ (selectedUser.displayName || selectedUser.email || 'U')[0].toUpperCase() }}
              </div>
              <div>
                <h4 class="text-white font-medium text-sm">{{ selectedUser.displayName || selectedUser.email }}</h4>
                <p class="text-green-500 text-[10px] uppercase font-bold">В сети</p>
              </div>
            </div>
          </div>

          <!-- Messages -->
          <div class="flex-grow overflow-y-auto p-6 space-y-4 bg-slate-950/30" ref="messageContainer">
            <div 
              v-for="msg in messages" 
              :key="msg.id"
              class="flex"
              :class="msg.senderId === userStore.user.uid ? 'justify-end' : 'justify-start'"
            >
              <div 
                class="max-w-[70%] rounded-2xl px-4 py-2 text-sm shadow-lg"
                :class="msg.senderId === userStore.user.uid ? 'bg-blue-600 text-white rounded-tr-none' : 'bg-slate-800 text-slate-200 rounded-tl-none'"
              >
                <p>{{ msg.text }}</p>
                <div class="text-[9px] mt-1 opacity-50 text-right">
                  {{ msg.createdAt ? new Date(msg.createdAt).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'}) : '...' }}
                </div>
              </div>
            </div>
          </div>

          <!-- Input -->
          <div class="p-4 border-t border-white/5">
            <form @submit.prevent="sendMessage" class="flex items-center space-x-4">
              <input 
                v-model="newMessage"
                type="text" 
                placeholder="Введите сообщение..." 
                class="flex-grow bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-1 focus:ring-blue-500/50" 
              />
              <button 
                type="submit"
                :disabled="!newMessage.trim()"
                class="p-3 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 rounded-lg text-white transition-all shadow-lg shadow-blue-600/20"
              >
                ✈️
              </button>
            </form>
          </div>
        </template>

        <div v-else class="flex-grow flex flex-col items-center justify-center text-slate-500 space-y-4">
          <div class="text-6xl text-slate-800">💬</div>
          <p class="italic">Выберите студента для начала общения</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: 'auth'
})

const userStore = useUserStore()
const { fetchApi: api } = useApi()
const route = useRoute()

const users = ref([])
const searchQuery = ref('')
const selectedUser = ref(null)
const messages = ref([])
const newMessage = ref('')
const messageContainer = ref(null)

let pollTimer = null

const filteredUsers = computed(() => {
  return users.value.filter(u => 
    u.uid !== userStore.user.uid && 
    (u.displayName?.toLowerCase().includes(searchQuery.value.toLowerCase()) || 
     u.email?.toLowerCase().includes(searchQuery.value.toLowerCase()))
  )
})

const stopPolling = () => {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messageContainer.value) {
      messageContainer.value.scrollTop = messageContainer.value.scrollHeight
    }
  })
}

const markAsRead = async (otherUserId) => {
  try {
    await api('/chat/messages/read', {
      method: 'POST',
      body: {
        user_id: otherUserId
      }
    })
  } catch (e) {
    console.error('[Chat Debug] Failed to mark messages as read:', e)
  }
}

const fetchMessages = async () => {
  if (!selectedUser.value) return

  try {
    const otherUserId = selectedUser.value.uid
    const data = await api(`/chat/messages?user_id=${encodeURIComponent(otherUserId)}&limit=200`)
    const nextMessages = data || []

    const prevLastId = messages.value.length ? messages.value[messages.value.length - 1].id : null
    const nextLastId = nextMessages.length ? nextMessages[nextMessages.length - 1].id : null

    messages.value = nextMessages

    await markAsRead(otherUserId)

    if (prevLastId !== nextLastId) {
      scrollToBottom()
    }
  } catch (e) {
    console.error('[Chat Debug] Failed to fetch messages:', e)
  }
}

const startPollingMessages = () => {
  stopPolling()
  fetchMessages()
  pollTimer = setInterval(fetchMessages, 2500)
}

const fetchUsers = async () => {
  try {
    const data = await api('/users')
    console.log('[Chat Debug] Fetched Users:', data)
    users.value = data || []

    const targetUID = typeof route.query.uid === 'string' ? route.query.uid.trim() : ''
    if (targetUID) {
      const userFromQuery = users.value.find((u) => u.uid === targetUID)
      if (userFromQuery) {
        selectUser(userFromQuery)
      }
    }
  } catch (e) {
    console.error('[Chat Debug] Failed to fetch users:', e)
  }
}

const selectUser = (user) => {
  messages.value = []
  selectedUser.value = user
  startPollingMessages()
}

const sendMessage = async () => {
  if (!newMessage.value.trim() || !selectedUser.value) return

  const otherUserId = selectedUser.value.uid

  const text = newMessage.value
  newMessage.value = ''

  try {
    await api('/chat/messages', {
      method: 'POST',
      body: {
        receiver_id: otherUserId,
        text
      }
    })

    await fetchMessages()
    scrollToBottom()
  } catch (e) {
    console.error('[Chat Debug] Error sending message:', e)
  }
}

onMounted(() => {
  fetchUsers()
})

onUnmounted(() => {
  stopPolling()
})
</script>

