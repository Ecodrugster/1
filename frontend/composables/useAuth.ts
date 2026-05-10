import {
  type User,
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

  const ensureAuth = () => {
    if (!$auth) {
      throw new Error('Firebase Auth is not configured. Set FIREBASE_* variables in frontend/.env and restart Nuxt.')
    }
  }

  const syncSessionFromUser = async (user: User) => {
    const token = await user.getIdToken()
    userStore.setUser(user)
    userStore.setToken(token)
    return token
  }

  const initAuth = () => {
    if (!$auth) {
      console.warn('Auth not initialized yet')
      return
    }
    
    onIdTokenChanged($auth, async (user) => {
      if (user) {
        await syncSessionFromUser(user)
        
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
    ensureAuth()
    return signInWithEmailAndPassword($auth, email, pass)
  }

  const loginWithGoogle = async () => {
    ensureAuth()
    const provider = new GoogleAuthProvider()
    const cred = await signInWithPopup($auth, provider)
    await syncSessionFromUser(cred.user)
    
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
    ensureAuth()
    const cred = await createUserWithEmailAndPassword($auth, email, pass)
    await syncSessionFromUser(cred.user)
    
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
    ensureAuth()
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
