<template>
  <div class="container mx-auto p-4">
    <h1 class="text-2xl font-bold mb-4">Contacts</h1>

    <!-- Add Contact Button -->
    <button @click="openAddModal" class="btn btn-primary mb-6">Add New Contact</button>

    <!-- Loading State -->
    <div v-if="loading" class="text-center p-8">
      <span class="loading loading-infinity loading-lg"></span>
      <p>Loading contacts...</p>
    </div>

    <!-- Contacts List -->
    <div class="grid gap-4">
      <div v-if="!loading && !contacts.length" class="text-center p-8 text-gray-500">
        No contacts found
      </div>

      <div v-for="contact in contacts" :key="contact.ID" class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">{{ contact.fullName }}</h2>
          <p>Phone: {{ contact.phoneNumber }}</p>
          <p>Note: {{ contact.note }}</p>
          <div class="card-actions justify-end mt-2">
            <button @click="openEditModal(contact)" class="btn btn-sm btn-ghost">Edit</button>
            <button @click="deleteContact(contact.ID)" class="btn btn-sm btn-error">Delete</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Pagination Controls -->
    <div class="flex justify-between mt-6">
      <button @click="prevPage" :disabled="currentPage === 0" class="btn btn-outline">Previous</button>
      <span>Page {{ currentPage + 1 }}</span>
      <button @click="nextPage" :disabled="!hasMore" class="btn btn-outline">Next</button>
    </div>

    <!-- Add Contact Modal -->
    <div class="modal" :class="{ 'modal-open': showAddModal }">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4">Add New Contact</h3>
        <div class="flex flex-col gap-3">
          <input v-model="newContact.fullName" placeholder="Full Name" class="input input-bordered">
          <input v-model="newContact.phoneNumber" placeholder="Phone Number" class="input input-bordered">
          <input v-model="newContact.note" placeholder="Note" class="input input-bordered">
        </div>
        <div class="modal-action">
          <button class="btn" @click="closeAddModal">Cancel</button>
          <button class="btn btn-primary" @click="createContact">Add</button>
        </div>
      </div>
    </div>

    <!-- Edit Contact Modal -->
    <div class="modal" :class="{ 'modal-open': showEditModal }">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4">Edit Contact</h3>
        <div class="flex flex-col gap-3">
          <input v-model="editingContact.fullName" class="input input-bordered">
          <input v-model="editingContact.phoneNumber" class="input input-bordered">
          <input v-model="editingContact.note" class="input input-bordered">
        </div>
        <div class="modal-action">
          <button class="btn" @click="closeEditModal">Cancel</button>
          <button class="btn btn-success" @click="updateContact">Save</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Contact {
  ID: number
  fullName: string
  phoneNumber: string
  note: string
}

interface ApiResponse {
  contacts?: Contact[]
  contact?: Contact
  error?: {
    msg: string
  }
}

const runtimeConfig = useRuntimeConfig()
const apiBase = runtimeConfig.public.apiBase

const contacts = ref<Contact[]>([])
const newContact = ref<Partial<Contact>>({})
const editingContact = ref<Partial<Contact>>({})
const showAddModal = ref(false)
const showEditModal = ref(false)
const loading = ref(false)

// Переменные пагинации
const currentPage = ref(0)
const pageSize = 10
const hasMore = ref(true)

const fetchContacts = async () => {
  loading.value = true
  try {
    const response = await $fetch<ApiResponse>(`${apiBase}/v1/contact`, {
      query: { page: currentPage.value, size: pageSize }
    })

    if (response.error) {
      throw new Error(response.error.msg)
    }

    // Записываем контакты, либо устанавливаем пустой массив
    contacts.value = response.contacts || []
    // Если вернулось меньше записей, чем ожидается, это последняя страница
    hasMore.value = (contacts.value.length === pageSize)
  } catch (error) {
    alert(error instanceof Error ? error.message : 'Error fetching contacts')
  } finally {
    loading.value = false
  }
}

const nextPage = async () => {
  if (!hasMore.value) return
  currentPage.value++
  await fetchContacts()
}

const prevPage = async () => {
  if (currentPage.value <= 0) return
  currentPage.value--
  await fetchContacts()
}

const createContact = async () => {
  try {
    const response = await $fetch<ApiResponse>(`${apiBase}/v1/contact`, {
      method: 'POST',
      body: newContact.value
    })

    if (response.error) {
      throw new Error(response.error.msg)
    }

    newContact.value = {}
    // Перезагружаем текущую страницу
    await fetchContacts()
    closeAddModal()
  } catch (error) {
    alert(error instanceof Error ? error.message : 'Error creating contact')
  }
}

const updateContact = async () => {
  try {
    if (!editingContact.value?.ID) return

    const response = await $fetch<ApiResponse>(`${apiBase}/v1/contact`, {
      method: 'PUT',
      body: editingContact.value
    })

    if (response.error) {
      throw new Error(response.error.msg)
    }

    await fetchContacts()
    closeEditModal()
  } catch (error) {
    alert(error instanceof Error ? error.message : 'Error updating contact')
  }
}

const deleteContact = async (id: number) => {
  try {
    const response = await $fetch<ApiResponse>(`${apiBase}/v1/contact`, {
      method: 'DELETE',
      body: { id },
      parseResponse: false
    })

    if (response?.error) {
      throw new Error(response.error.msg)
    }

    await fetchContacts()
  } catch (error) {
    alert(error instanceof Error ? error.message : 'Error deleting contact')
  }
}

const openAddModal = () => showAddModal.value = true
const closeAddModal = () => showAddModal.value = false

const openEditModal = (contact: Contact) => {
  editingContact.value = { ...contact }
  showEditModal.value = true
}

const closeEditModal = () => {
  editingContact.value = {}
  showEditModal.value = false
}

onMounted(() => {
  fetchContacts()
})
</script>
