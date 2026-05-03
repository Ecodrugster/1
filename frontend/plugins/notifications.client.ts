import { 
  collectionGroup, 
  query, 
  where, 
  onSnapshot,
  orderBy
} from 'firebase/firestore'

export default defineNuxtPlugin((nuxtApp) => {
  const userStore = useUserStore()
  const notificationStore = useNotificationStore()
  const { $firestore } = nuxtApp

  if (process.client) {
    // Ждем, пока пользователь авторизуется
    watch(() => userStore.user, (user) => {
      if (user) {
        console.log('[Notifications] Starting global listener for UID:', user.uid)
        
        // Слушаем все сообщения во всех чатах, где мы получатели и сообщение не прочитано
        const q = query(
          collectionGroup($firestore, 'messages'),
          where('receiverId', '==', user.uid),
          where('read', '==', false)
        )

        onSnapshot(q, (snapshot) => {
          notificationStore.setUnreadCount(snapshot.size)
          
          snapshot.docChanges().forEach((change) => {
            if (change.type === 'added') {
              const msg = change.doc.data()
              console.log('[Notifications] New message received!', msg.text)
              // Здесь можно добавить Browser Notification API в будущем
            }
          })
        })
      } else {
        notificationStore.clear()
      }
    }, { immediate: true })
  }
})
