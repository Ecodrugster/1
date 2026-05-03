import { defineStore } from 'pinia'

export const useNotificationStore = defineStore('notifications', {
  state: () => ({
    unreadCount: 0,
    notifications: []
  }),
  actions: {
    increment() {
      this.unreadCount++
    },
    clear() {
      this.unreadCount = 0
    },
    setUnreadCount(count: number) {
      this.unreadCount = count
    }
  }
})
