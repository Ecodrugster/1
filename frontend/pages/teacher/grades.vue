<template>
  <div class="max-w-6xl mx-auto py-8 px-4">
    <div class="flex justify-between items-center mb-8">
      <div>
        <h1 class="text-3xl font-bold text-white">Журнал оценок</h1>
        <p class="text-slate-400">Выставление оценок студентам ITSTEP</p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Student List -->
      <div class="lg:col-span-1 space-y-4">
        <div class="bg-slate-900 border border-white/5 rounded-2xl p-6 shadow-xl">
          <h2 class="text-xl font-bold text-white mb-4">Студенты</h2>
          <div class="space-y-2 max-h-[600px] overflow-y-auto pr-2">
            <button 
              v-for="student in students" 
              :key="student.uid"
              @click="selectedStudent = student"
              class="w-full flex items-center space-x-3 p-3 rounded-xl transition-all border border-transparent"
              :class="selectedStudent?.uid === student.uid ? 'bg-blue-600/20 border-blue-600/50' : 'hover:bg-white/5'"
            >
              <div class="w-10 h-10 rounded-full bg-slate-800 flex items-center justify-center text-xs font-bold text-white">
                {{ (student.displayName || student.email || 'U')[0].toUpperCase() }}
              </div>
              <div class="text-left overflow-hidden">
                <div class="text-sm font-medium text-white truncate">{{ student.displayName || student.email }}</div>
                <div class="text-[10px] text-slate-500 uppercase">Студент</div>
              </div>
            </button>
          </div>
        </div>
      </div>

      <!-- Grading Form -->
      <div class="lg:col-span-2">
        <div v-if="selectedStudent" class="bg-slate-900 border border-white/5 rounded-2xl p-8 shadow-xl">
          <div class="flex items-center space-x-4 mb-8 pb-8 border-b border-white/5">
            <div class="w-16 h-16 rounded-full bg-blue-600 flex items-center justify-center text-xl font-bold text-white shadow-lg shadow-blue-600/30">
              {{ (selectedStudent.displayName || selectedStudent.email || 'U')[0].toUpperCase() }}
            </div>
            <div>
              <h2 class="text-2xl font-bold text-white">{{ selectedStudent.displayName || selectedStudent.email }}</h2>
              <p class="text-slate-400">Выставление новой оценки</p>
            </div>
          </div>

          <form @submit.prevent="submitGrade" class="space-y-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-slate-400 mb-2">Предмет</label>
                <input 
                  v-model="form.subject" 
                  type="text" 
                  placeholder="Напр. Программирование на Go"
                  class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50" 
                  required
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-slate-400 mb-2">Оценка (1-12)</label>
                <input 
                  v-model.number="form.value" 
                  type="number" 
                  min="1" 
                  max="12"
                  class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50" 
                  required
                />
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-slate-400 mb-2">Комментарий</label>
              <textarea 
                v-model="form.comment" 
                rows="4" 
                placeholder="За что выставлена оценка..."
                class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50 resize-none"
              ></textarea>
            </div>

            <div class="flex justify-end">
              <button 
                type="submit" 
                :disabled="submitting"
                class="bg-blue-600 hover:bg-blue-500 text-white px-8 py-3 rounded-xl font-bold shadow-lg shadow-blue-600/20 transition-all disabled:opacity-50"
              >
                {{ submitting ? 'Отправка...' : 'Выставить оценку' }}
              </button>
            </div>
          </form>

          <!-- Recent Grades for this student -->
          <div class="mt-12 pt-8 border-t border-white/5">
            <h3 class="text-lg font-bold text-white mb-4">Недавние оценки</h3>
            <div class="space-y-3">
              <div v-for="g in studentGrades" :key="g.id" class="bg-white/5 p-4 rounded-xl flex justify-between items-center">
                <div>
                  <div class="text-sm font-bold text-white">{{ g.subject }}</div>
                  <div class="text-xs text-slate-500">{{ formatDate(g.created_at) }} • {{ g.comment }}</div>
                </div>
                <div class="text-2xl font-bold text-blue-500">{{ g.value }}</div>
              </div>
              <div v-if="studentGrades.length === 0" class="text-center text-slate-500 py-4 italic text-sm">
                Оценок пока нет
              </div>
            </div>
          </div>
        </div>

        <div v-else class="h-full bg-slate-900/50 border border-dashed border-white/10 rounded-2xl flex flex-col items-center justify-center text-slate-500 p-12">
          <div class="text-6xl mb-4">👨‍🎓</div>
          <p>Выберите студента из списка слева, чтобы выставить ему оценку</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const { fetchApi: api } = useApi()
const students = ref([])
const selectedStudent = ref(null)
const studentGrades = ref([])
const submitting = ref(false)

const form = reactive({
  subject: '',
  value: 12,
  comment: ''
})

const fetchStudents = async () => {
  const data = await api('/users')
  // Фильтруем только студентов (если в будущем будут другие роли)
  students.value = data.filter(u => u.role === 'student' || !u.role)
}

const fetchStudentGrades = async (uid) => {
  const data = await api(`/grades?student_id=${uid}`)
  studentGrades.value = data || []
}

watch(selectedStudent, (newVal) => {
  if (newVal) {
    fetchStudentGrades(newVal.uid)
  }
})

const submitGrade = async () => {
  if (!selectedStudent.value) return
  submitting.value = true
  try {
    await api('/teacher/grades', {
      method: 'POST',
      body: {
        student_id: selectedStudent.value.uid,
        ...form
      }
    })
    form.subject = ''
    form.comment = ''
    fetchStudentGrades(selectedStudent.value.uid)
    alert('Оценка успешно выставлена!')
  } catch (e) {
    alert('Ошибка: ' + e.message)
  } finally {
    submitting.value = false
  }
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return d.toLocaleDateString()
}

onMounted(fetchStudents)
</script>
