import { 
  onIdTokenChanged, 
  signInWithEmailAndPassword, 
  createUserWithEmailAndPassword, 
  signOut,
  GoogleAuthProvider,
  signInWithPopup,
  sendEmailVerification
} from 'firebase/auth'

export const useAuth = () => {
  const { $auth } = useNuxtApp()
  const userStore = useUserStore()

  const initAuth = () => {
    if (!$auth) {
      console.warn('Auth not initialized yet')
      return
    }
    
    onIdTokenChanged($auth, async (user) => {
      if (user) {
        const token = await user.getIdToken()
        userStore.setUser(user)
        userStore.setToken(token)
        
        // Загружаем профиль с бэкенда (роль, доп. данные)
        try {
          const { fetchApi } = useApi()
          const profile = await fetchApi('/profile')
          userStore.setProfile(profile)
        } catch (e) {
          console.error('Failed to fetch user profile:', e)
        }
      } else {
        userStore.logout()
      }
    })
  }

  const login = async (email: string, pass: string) => {
    return signInWithEmailAndPassword($auth, email, pass)
  }

  const loginWithGoogle = async () => {
    const provider = new GoogleAuthProvider()
    const cred = await signInWithPopup($auth, provider)
    
    // Создаем/обновляем профиль на бэкенде
    const { fetchApi } = useApi()
    await fetchApi('/profile', {
      method: 'PUT',
      body: {
        email: cred.user.email,
        displayName: cred.user.displayName,
        photoURL: cred.user.photoURL,
        role: 'student'
      }
    })
    return cred
  }

  const register = async (email: string, pass: string) => {
    const cred = await createUserWithEmailAndPassword($auth, email, pass)
    
    // Создаем профиль на бэкенде
    const { fetchApi } = useApi()
    await fetchApi('/profile', {
      method: 'PUT',
      body: {
        email: email,
        displayName: email.split('@')[0],
        role: 'student'
      }
    })
    
    await sendEmailVerification(cred.user)
    return cred
  }

  const logout = async () => {
    return signOut($auth)
  }

  return {
    initAuth,
    login,
    loginWithGoogle,
    register,
    logout,
  }
}
