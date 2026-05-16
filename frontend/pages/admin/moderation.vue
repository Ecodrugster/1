<template>
  <div class="space-y-8">
    <section class="space-y-4">
      <div v-if="errorMessage" class="bg-red-500/10 border border-red-500/30 text-red-300 rounded-xl px-4 py-3 text-sm">
        {{ errorMessage }}
      </div>

      <div class="flex justify-between items-center">
        <h2 class="text-2xl font-bold text-white">Club Requests</h2>
        <div class="text-sm text-slate-500">{{ clubRequests.length }} in queue</div>
      </div>

      <div v-if="loadingClubs" class="space-y-3">
        <div v-for="i in 2" :key="i" class="h-24 bg-slate-900 border border-white/5 rounded-2xl animate-pulse"></div>
      </div>

      <div v-else-if="clubRequests.length === 0" class="bg-slate-900 border border-white/5 p-10 rounded-2xl text-center text-slate-500">
        No club requests to moderate.
      </div>

      <div v-else class="space-y-3">
        <div v-for="club in clubRequests" :key="club.id" class="bg-slate-900 border border-white/5 p-5 rounded-2xl flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
          <div>
            <div class="text-white font-bold">{{ club.name }}</div>
            <div class="text-sm text-slate-400 mt-1">{{ club.description }}</div>
            <div class="text-xs text-slate-500 mt-2">
              Creator: {{ club.created_by_name || club.created_by || 'unknown' }} • Members: {{ club.members?.length || 0 }}
            </div>
          </div>

          <div class="flex gap-2">
            <button
              @click="approveClub(club.id)"
              class="px-4 py-2 bg-green-500/10 hover:bg-green-500 text-green-400 hover:text-white rounded-lg text-sm transition-all border border-green-500/20"
            >
              Approve
            </button>
            <button
              @click="rejectClub(club.id)"
              class="px-4 py-2 bg-red-500/10 hover:bg-red-500 text-red-400 hover:text-white rounded-lg text-sm transition-all border border-red-500/20"
            >
              Reject
            </button>
          </div>
        </div>
      </div>
    </section>

    <section class="space-y-4">
      <div class="flex justify-between items-center">
        <h2 class="text-2xl font-bold text-white">Post Moderation</h2>
        <div class="text-sm text-slate-500">Latest 50 posts</div>
      </div>

      <div v-if="loadingPosts" class="space-y-4">
        <div v-for="i in 3" :key="i" class="bg-slate-900 border border-white/5 p-6 rounded-2xl animate-pulse h-32"></div>
      </div>

      <div v-else-if="posts.length === 0" class="bg-slate-900 border border-white/5 p-12 rounded-2xl text-center text-slate-500">
        No posts for moderation.
      </div>

      <div v-else class="space-y-4">
        <div v-for="post in posts" :key="post.id" class="bg-slate-900 border border-white/5 p-6 rounded-2xl flex justify-between items-start group hover:border-red-500/30 transition-all">
          <div class="flex gap-4">
            <div class="w-10 h-10 rounded-xl bg-slate-800 flex items-center justify-center text-xs text-slate-500">P</div>
            <div>
              <div class="flex items-center gap-2 mb-1">
                <span class="text-sm font-bold text-white">
                  Author: {{ post.author_name || post.author_id || 'unknown' }}
                </span>
                <span class="text-[10px] text-slate-500 font-mono">{{ formatDateTime(post.created_at) }}</span>
              </div>
              <p class="text-slate-300 text-sm leading-relaxed">{{ post.content }}</p>
            </div>
          </div>

          <button
            @click="deletePost(post.id)"
            class="opacity-0 group-hover:opacity-100 px-4 py-2 bg-red-500/10 hover:bg-red-500 text-red-500 hover:text-white rounded-xl text-xs font-bold transition-all border border-red-500/20"
          >
            Delete Post
          </button>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

const { fetchApi: api } = useApi()
const posts = ref([])
const clubRequests = ref([])
const loadingPosts = ref(true)
const loadingClubs = ref(true)
const errorMessage = ref('')

const formatDateTime = (value) => {
  if (!value) return 'time unknown'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'time unknown'
  return date.toLocaleString()
}

const fetchPosts = async () => {
  loadingPosts.value = true
  errorMessage.value = ''
  try {
    const data = await api('/admin/posts')
    posts.value = data || []
  } catch (e) {
    console.error('Failed to fetch posts:', e)
    const status = e?.status || e?.response?.status
    const message = e?.data?.error || e?.message || 'Failed to load posts'
    errorMessage.value = `Could not load posts (${status || 'no-status'}): ${message}`
  } finally {
    loadingPosts.value = false
  }
}

const fetchClubRequests = async () => {
  loadingClubs.value = true
  try {
    const data = await api('/admin/club-requests')
    clubRequests.value = data || []
  } catch (e) {
    console.error('Failed to fetch club requests:', e)
    const status = e?.status || e?.response?.status
    const message = e?.data?.error || e?.message || 'Failed to load club requests'
    errorMessage.value = `Could not load club requests (${status || 'no-status'}): ${message}`
  } finally {
    loadingClubs.value = false
  }
}

const deletePost = async (id) => {
  if (!confirm('Delete this post permanently?')) return

  try {
    await api(`/admin/posts/${id}`, { method: 'DELETE' })
    posts.value = posts.value.filter(p => p.id !== id)
    alert('Post deleted by moderator')
  } catch (e) {
    alert('Delete error: ' + (e?.data?.error || e.message))
  }
}

const approveClub = async (id) => {
  try {
    await api(`/admin/club-requests/${id}/approve`, { method: 'POST' })
    clubRequests.value = clubRequests.value.filter(c => c.id !== id)
  } catch (e) {
    alert('Approval error: ' + (e?.data?.error || e.message))
  }
}

const rejectClub = async (id) => {
  const reason = prompt('Reject reason (optional):') || ''
  try {
    await api(`/admin/club-requests/${id}/reject`, {
      method: 'POST',
      body: { comment: reason }
    })
    clubRequests.value = clubRequests.value.filter(c => c.id !== id)
  } catch (e) {
    alert('Reject error: ' + (e?.data?.error || e.message))
  }
}

onMounted(async () => {
  await Promise.all([fetchPosts(), fetchClubRequests()])
})
</script>
