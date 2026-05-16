export default defineNuxtPlugin(() => {
  const userStore = useUserStore()
  const notificationStore = useNotificationStore()
  const { fetchApi: api } = useApi()

  let pollTimer: ReturnType<typeof setInterval> | null = null

  const stopPolling = () => {
    if (pollTimer) {
      clearInterval(pollTimer)
      pollTimer = null
    }
  }

  const loadUnreadCount = async () => {
    if (!userStore.user) {
      notificationStore.clear()
      return
    }

    try {
      const data = await api<{ count?: number }>('/chat/unread-count')
      notificationStore.setUnreadCount(Number(data?.count || 0))
    } catch (e) {
      console.error('[Notifications] Failed to load unread count:', e)
    }
  }

  const startPolling = () => {
    stopPolling()
    loadUnreadCount()
    pollTimer = setInterval(loadUnreadCount, 5000)
  }

  watch(
    () => userStore.user,
    (user) => {
      if (user) {
        startPolling()
      } else {
        stopPolling()
        notificationStore.clear()
      }
    },
    { immediate: true }
  )
})
