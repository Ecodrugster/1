export const useApi = () => {
  const userStore = useUserStore()
  const config = useRuntimeConfig()
  const baseUrl = config.public.apiBaseUrl || 'http://localhost:8080/api/v1'

  const fetchApi = async (url: string, options: any = {}) => {
    const headers = {
      ...options.headers,
    }

    if (userStore.token) {
      headers['Authorization'] = `Bearer ${userStore.token}`
    }

    return $fetch(`${baseUrl}${url}`, {
      ...options,
      headers,
    })
  }

  return {
    fetchApi,
  }
}
