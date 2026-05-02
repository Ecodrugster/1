import { defineStore } from 'pinia'

interface UserState {
  user: any | null
  token: string | null
  profile: any | null
  loading: boolean
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    user: null,
    token: null,
    profile: null,
    loading: true,
  }),
  actions: {
    setUser(user: any) {
      this.user = user
      this.loading = false
    },
    setToken(token: string | null) {
      this.token = token
    },
    setProfile(profile: any) {
      this.profile = profile
    },
    logout() {
      this.user = null
      this.token = null
      this.profile = null
      this.loading = false
    }
  },
  getters: {
    isLoggedIn: (state) => !!state.user,
  }
})
