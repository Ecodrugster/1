import { defineStore } from 'pinia'

export const useAdminStore = defineStore('admin', {
  state: () => ({
    users: [],
    posts: [],
    news: [],
    clubs: [],
    loading: false,
    stats: {
      totalUsers: 0,
      totalPosts: 0
    }
  }),
  actions: {
    setUsers(users: any[]) {
      this.users = users
    },
    setPosts(posts: any[]) {
      this.posts = posts
    },
    setLoading(status: boolean) {
      this.loading = status
    }
  }
})
