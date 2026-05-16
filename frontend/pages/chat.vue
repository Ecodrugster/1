<template>
  <div class="mx-auto h-[calc(100vh-120px)] max-w-6xl px-4 py-8">
    <div class="flex h-full overflow-hidden rounded-xl border border-white/5 bg-slate-900 shadow-xl">
      <!-- Список пользователей -->
      <div class="flex w-80 flex-col border-r border-white/5">
        <div class="border-b border-white/5 p-4">
          <h3 class="mb-2 font-semibold text-white">Студенты</h3>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Поиск..."
            class="w-full rounded-lg border-none bg-slate-800 px-4 py-2 text-sm text-white focus:ring-1 focus:ring-blue-500/50"
          />
        </div>

        <div class="flex-grow overflow-y-auto">
          <div
            v-for="user in filteredUsers"
            :key="user.uid"
            class="flex cursor-pointer items-center space-x-3 border-b border-white/5 p-4 transition-all last:border-0 hover:bg-white/5"
            :class="{ 'border-l-2 border-l-blue-600 bg-blue-600/10': selectedUser?.uid === user.uid }"
            @click="selectUser(user)"
          >
            <img v-if="user.photoURL" :src="user.photoURL" class="h-10 w-10 rounded-full" />
            <div
              v-else
              class="flex h-10 w-10 items-center justify-center rounded-full bg-slate-800 text-xs font-bold text-slate-500"
            >
              {{ (user.displayName || user.email || 'U')[0].toUpperCase() }}
            </div>

            <div class="min-w-0 flex-grow">
              <h4 class="truncate text-sm font-medium text-white">{{ user.displayName || user.email }}</h4>
              <p class="text-[10px] uppercase text-slate-500">Студент</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Окно переписки -->
      <div class="flex flex-grow flex-col">
        <template v-if="selectedUser">
          <div class="flex items-center justify-between border-b border-white/5 bg-slate-900/50 p-4">
            <div class="flex items-center space-x-3">
              <img v-if="selectedUser.photoURL" :src="selectedUser.photoURL" class="h-10 w-10 rounded-full" />
              <div
                v-else
                class="flex h-10 w-10 items-center justify-center rounded-full bg-slate-800 text-xs font-bold text-white"
              >
                {{ (selectedUser.displayName || selectedUser.email || 'U')[0].toUpperCase() }}
              </div>

              <div>
                <h4 class="text-sm font-medium text-white">{{ selectedUser.displayName || selectedUser.email }}</h4>
                <p class="text-[10px] font-bold uppercase text-green-500">В сети</p>
              </div>
            </div>
          </div>

          <div ref="messageContainer" class="flex-grow space-y-4 overflow-y-auto bg-slate-950/30 p-6">
            <div
              v-for="msg in messages"
              :key="msg.id"
              class="flex"
              :class="msg.senderId === userStore.user.uid ? 'justify-end' : 'justify-start'"
            >
              <div
                class="max-w-[70%] rounded-2xl px-4 py-2 text-sm shadow-lg"
                :class="msg.senderId === userStore.user.uid ? 'rounded-tr-none bg-blue-600 text-white' : 'rounded-tl-none bg-slate-800 text-slate-200'"
              >
                <p>{{ msg.text }}</p>
                <div class="mt-1 text-right text-[9px] opacity-50">
                  {{ msg.createdAt ? new Date(msg.createdAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : '...' }}
                </div>
              </div>
            </div>
          </div>

          <div class="border-t border-white/5 p-4">
            <div
              v-if="sendError"
              class="mb-3 rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-300"
            >
              {{ sendError }}
            </div>

            <form class="flex items-center space-x-4" @submit.prevent="sendMessage">
              <input
                v-model="newMessage"
                type="text"
                placeholder="Введите сообщение..."
                class="flex-grow rounded-lg border-none bg-slate-800 px-4 py-3 text-white focus:ring-1 focus:ring-blue-500/50"
                @input="sendError = ''"
              />

              <button
                type="submit"
                :disabled="!newMessage.trim() || isSending"
                class="rounded-lg bg-blue-600 px-4 py-3 text-sm text-white shadow-lg shadow-blue-600/20 transition-all hover:bg-blue-500 disabled:opacity-50"
              >
                Отправить
              </button>
            </form>

            <div class="mt-2 text-right text-[11px]" :class="messageLength > MAX_MESSAGE_LENGTH ? 'text-red-400' : 'text-slate-500'">
              {{ messageLength }} / {{ MAX_MESSAGE_LENGTH }}
            </div>
          </div>
        </template>

        <div v-else class="flex flex-grow flex-col items-center justify-center space-y-4 text-slate-500">
          <div class="text-6xl text-slate-800">💬</div>
          <p class="italic">Выберите студента, чтобы начать диалог</p>
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

const MAX_MESSAGE_LENGTH = 10000

const users = ref([])
const searchQuery = ref('')
const selectedUser = ref(null)
const messages = ref([])
const newMessage = ref('')
const sendError = ref('')
const isSending = ref(false)
const messageContainer = ref(null)

let pollTimer = null

const messageLength = computed(() => Array.from(newMessage.value).length)

const filteredUsers = computed(() => {
  return users.value.filter((u) => {
    return (
      u.uid !== userStore.user.uid &&
      (u.displayName?.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        u.email?.toLowerCase().includes(searchQuery.value.toLowerCase()))
    )
  })
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
  sendError.value = ''
  selectedUser.value = user
  startPollingMessages()
}

const sendMessage = async () => {
  if (!newMessage.value.trim() || !selectedUser.value || isSending.value) return

  if (messageLength.value > MAX_MESSAGE_LENGTH) {
    sendError.value = `Сообщение слишком длинное. Максимум ${MAX_MESSAGE_LENGTH} символов.`
    return
  }

  const otherUserId = selectedUser.value.uid
  const text = newMessage.value
  isSending.value = true
  sendError.value = ''

  try {
    await api('/chat/messages', {
      method: 'POST',
      body: {
        receiver_id: otherUserId,
        text
      }
    })

    newMessage.value = ''
    await fetchMessages()
    scrollToBottom()
  } catch (e) {
    const backendMessage = String(e?.data?.error || e?.data?.message || e?.message || 'Неизвестная ошибка')
    const isTooLongError = /too long|maximum|длин/i.test(backendMessage.toLowerCase())

    if (isTooLongError) {
      const match = backendMessage.match(/(\d{2,6})/)
      const backendLimit = match ? Number(match[1]) : null
      const safeLimit = Number.isFinite(backendLimit) ? backendLimit : MAX_MESSAGE_LENGTH
      sendError.value = `Сообщение слишком длинное. Максимум ${safeLimit} символов.`
    } else {
      sendError.value = `Не удалось отправить сообщение: ${backendMessage}`
    }

    console.error('[Chat Debug] Error sending message:', e)
  } finally {
    isSending.value = false
  }
}

onMounted(() => {
  fetchUsers()
})

onUnmounted(() => {
  stopPolling()
})
</script>
