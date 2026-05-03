export default defineNuxtRouteMiddleware((to, from) => {
  const userStore = useUserStore()
  
  // Если пользователь не авторизован или его роль не 'admin'
  const isAdmin = userStore.profile?.role === 'admin'
  
  if (!isAdmin) {
    console.warn('[Admin Middleware] Access denied. Not an admin.')
    return navigateTo('/')
  }
})
