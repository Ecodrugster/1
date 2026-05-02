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
    return signInWithPopup($auth, provider)
  }

  const register = async (email: string, pass: string) => {
    const cred = await createUserWithEmailAndPassword($auth, email, pass)
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
