<template>
  <div>
    <div style="display:flex;justify-content:space-between;margin-bottom:16px">
      <el-input v-model="keyword" placeholder="搜索姓名/身份证号/电话" style="width:300px" clearable @clear="loadData" @keyup.enter="loadData" />
      <el-button type="primary" @click="openDialog()">新增船员</el-button>
    </div>

    <el-table :data="tableData" v-loading="loading" border stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="姓名" width="120" />
      <el-table-column label="性别" width="80">
        <template #default="{ row }">
          {{ row.gender === 1 ? '男' : row.gender === 2 ? '女' : '未知' }}
        </template>
      </el-table-column>
      <el-table-column prop="id_number" label="身份证号" width="200" />
      <el-table-column prop="phone" label="电话" width="140" />
      <el-table-column prop="rank" label="职务" width="120" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 0 ? 'info' : row.status === 1 ? 'success' : 'warning'">
            {{ row.status === 0 ? '待派' : row.status === 1 ? '在船' : '休假' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="openDialog(row)">编辑</el-button>
          <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div style="display:flex;justify-content:flex-end;margin-top:16px">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="loadData"
        @current-change="loadData"
      />
    </div>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑船员' : '新增船员'" width="520px" destroy-on-close>
      <el-form :model="form" label-width="80px">
        <el-form-item label="姓名">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="性别">
          <el-radio-group v-model="form.gender">
            <el-radio :value="1">男</el-radio>
            <el-radio :value="2">女</el-radio>
            <el-radio :value="0">未知</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="出生日期">
          <el-date-picker v-model="form.birthday" type="date" value-format="YYYY-MM-DD" style="width:100%" />
        </el-form-item>
        <el-form-item label="身份证号">
          <el-input v-model="form.id_number" />
        </el-form-item>
        <el-form-item label="电话">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="职务">
          <el-input v-model="form.rank" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status" style="width:100%">
            <el-option label="待派" :value="0" />
            <el-option label="在船" :value="1" />
            <el-option label="休假" :value="2" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getSeafarers, createSeafarer, updateSeafarer, deleteSeafarer } from '../api/modules.js'

const keyword = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const tableData = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const form = reactive({
  id: null,
  name: '',
  gender: 1,
  birthday: '',
  id_number: '',
  phone: '',
  email: '',
  rank: '',
  status: 0,
})

function resetForm() {
  form.id = null
  form.name = ''
  form.gender = 1
  form.birthday = ''
  form.id_number = ''
  form.phone = ''
  form.email = ''
  form.rank = ''
  form.status = 0
}

async function loadData() {
  loading.value = true
  try {
    const res = await getSeafarers({ page: page.value, page_size: pageSize.value, keyword: keyword.value })
    tableData.value = res.data.items || []
    total.value = res.data.total || 0
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

function openDialog(row) {
  resetForm()
  if (row) {
    isEdit.value = true
    Object.assign(form, {
      id: row.id,
      name: row.name,
      gender: row.gender,
      birthday: row.birthday ? row.birthday.slice(0, 10) : '',
      id_number: row.id_number,
      phone: row.phone,
      email: row.email,
      rank: row.rank,
      status: row.status,
    })
  } else {
    isEdit.value = false
  }
  dialogVisible.value = true
}

async function handleSubmit() {
  submitting.value = true
  try {
    const payload = { ...form }
    if (isEdit.value) {
      await updateSeafarer(form.id, payload)
    } else {
      await createSeafarer(payload)
    }
    ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
    dialogVisible.value = false
    loadData()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(`确定删除船员「${row.name}」？`, '提示', { type: 'warning' })
    await deleteSeafarer(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch {
    // cancelled
  }
}

onMounted(() => loadData())
</script>
