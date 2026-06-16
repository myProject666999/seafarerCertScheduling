<template>
  <div>
    <div style="display:flex;gap:20px">
      <div style="flex:1">
        <el-card>
          <template #header>
            <div style="display:flex;justify-content:space-between;align-items:center">
              <span>船员证书</span>
              <div style="display:flex;gap:8px">
                <el-select v-model="selectedSeafarerId" placeholder="选择船员" filterable style="width:200px" @change="loadCerts">
                  <el-option v-for="s in seafarers" :key="s.id" :label="s.name" :value="s.id" />
                </el-select>
                <el-button type="primary" :disabled="!selectedSeafarerId" @click="openCertDialog()">新增证书</el-button>
              </div>
            </div>
          </template>
          <el-table :data="certData" v-loading="certLoading" border stripe>
            <el-table-column prop="cert_number" label="证书编号" width="160" />
            <el-table-column label="证书类型名" width="140">
              <template #default="{ row }">
                {{ row.certificate_type?.name || certTypeMap[row.certificate_type_id] || '' }}
              </template>
            </el-table-column>
            <el-table-column label="签发日期" width="120">
              <template #default="{ row }">{{ row.issue_date?.slice(0, 10) }}</template>
            </el-table-column>
            <el-table-column label="到期日期" width="120">
              <template #default="{ row }">{{ row.expire_date?.slice(0, 10) }}</template>
            </el-table-column>
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 0 ? 'danger' : row.status === 1 ? 'success' : 'warning'">
                  {{ row.status === 0 ? '过期' : row.status === 1 ? '有效' : '即将过期' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="140">
              <template #default="{ row }">
                <el-button link type="primary" @click="openCertDialog(row)">编辑</el-button>
                <el-button link type="danger" @click="handleDeleteCert(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </div>

      <div style="flex:1">
        <el-card>
          <template #header>
            <div style="display:flex;justify-content:space-between;align-items:center">
              <span>证书类型管理</span>
              <el-button type="primary" @click="openCertTypeDialog()">新增类型</el-button>
            </div>
          </template>
          <el-table :data="certTypes" v-loading="certTypeLoading" border stripe>
            <el-table-column prop="name" label="名称" width="140" />
            <el-table-column prop="code" label="编码" width="120" />
            <el-table-column label="有效期(月)" width="120">
              <template #default="{ row }">{{ row.validity_months ?? '-' }}</template>
            </el-table-column>
            <el-table-column label="是否必须" width="100">
              <template #default="{ row }">
                <el-tag :type="row.is_required ? 'danger' : 'info'">{{ row.is_required ? '是' : '否' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="80">
              <template #default="{ row }">
                <el-button link type="danger" @click="handleDeleteCertType(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </div>
    </div>

    <el-dialog v-model="certDialogVisible" :title="isEditCert ? '编辑证书' : '新增证书'" width="520px" destroy-on-close>
      <el-form :model="certForm" label-width="100px">
        <el-form-item label="证书编号">
          <el-input v-model="certForm.cert_number" />
        </el-form-item>
        <el-form-item label="证书类型">
          <el-select v-model="certForm.certificate_type_id" style="width:100%">
            <el-option v-for="ct in certTypes" :key="ct.id" :label="ct.name" :value="ct.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="签发日期">
          <el-date-picker v-model="certForm.issue_date" type="date" value-format="YYYY-MM-DD" style="width:100%" />
        </el-form-item>
        <el-form-item label="到期日期">
          <el-date-picker v-model="certForm.expire_date" type="date" value-format="YYYY-MM-DD" style="width:100%" />
        </el-form-item>
        <el-form-item label="证书图片URL">
          <el-input v-model="certForm.cert_image_url" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="certDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="certSubmitting" @click="handleCertSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="certTypeDialogVisible" title="新增证书类型" width="480px" destroy-on-close>
      <el-form :model="certTypeForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="certTypeForm.name" />
        </el-form-item>
        <el-form-item label="编码">
          <el-input v-model="certTypeForm.code" />
        </el-form-item>
        <el-form-item label="有效期(月)">
          <el-input-number v-model="certTypeForm.validity_months" :min="0" />
        </el-form-item>
        <el-form-item label="是否必须">
          <el-select v-model="certTypeForm.is_required" style="width:100%">
            <el-option label="是" :value="1" />
            <el-option label="否" :value="0" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="certTypeDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="certTypeSubmitting" @click="handleCertTypeSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getSeafarers, getSeafarerCerts, createSeafarerCert, updateSeafarerCert, deleteSeafarerCert, getCertTypes, createCertType, deleteCertType } from '../api/modules.js'

const seafarers = ref([])
const selectedSeafarerId = ref(null)
const certData = ref([])
const certLoading = ref(false)
const certDialogVisible = ref(false)
const isEditCert = ref(false)
const certSubmitting = ref(false)
const certForm = reactive({
  id: null,
  seafarer_id: null,
  cert_number: '',
  certificate_type_id: null,
  issue_date: '',
  expire_date: '',
  cert_image_url: '',
})

const certTypes = ref([])
const certTypeMap = computed(() => {
  const m = {}
  certTypes.value.forEach(ct => { m[ct.id] = ct.name })
  return m
})
const certTypeLoading = ref(false)
const certTypeDialogVisible = ref(false)
const certTypeSubmitting = ref(false)
const certTypeForm = reactive({
  name: '',
  code: '',
  validity_months: null,
  is_required: 0,
})

function resetCertForm() {
  certForm.id = null
  certForm.seafarer_id = null
  certForm.cert_number = ''
  certForm.certificate_type_id = null
  certForm.issue_date = ''
  certForm.expire_date = ''
  certForm.cert_image_url = ''
}

function resetCertTypeForm() {
  certTypeForm.name = ''
  certTypeForm.code = ''
  certTypeForm.validity_months = null
  certTypeForm.is_required = 0
}

async function loadSeafarers() {
  try {
    const res = await getSeafarers({ page: 1, page_size: 1000 })
    seafarers.value = res.data.items || []
  } catch (e) {
    ElMessage.error(e.message || '加载船员列表失败')
  }
}

async function loadCertTypes() {
  certTypeLoading.value = true
  try {
    const res = await getCertTypes()
    certTypes.value = res.data || []
  } catch (e) {
    ElMessage.error(e.message || '加载证书类型失败')
  } finally {
    certTypeLoading.value = false
  }
}

async function loadCerts() {
  if (!selectedSeafarerId.value) {
    certData.value = []
    return
  }
  certLoading.value = true
  try {
    const res = await getSeafarerCerts({ seafarer_id: selectedSeafarerId.value })
    certData.value = res.data || []
  } catch (e) {
    ElMessage.error(e.message || '加载证书失败')
  } finally {
    certLoading.value = false
  }
}

function openCertDialog(row) {
  resetCertForm()
  if (row) {
    isEditCert.value = true
    Object.assign(certForm, {
      id: row.id,
      seafarer_id: row.seafarer_id,
      cert_number: row.cert_number,
      certificate_type_id: row.certificate_type_id,
      issue_date: row.issue_date?.slice(0, 10) || '',
      expire_date: row.expire_date?.slice(0, 10) || '',
      cert_image_url: row.cert_image_url || '',
    })
  } else {
    isEditCert.value = false
    certForm.seafarer_id = selectedSeafarerId.value
  }
  certDialogVisible.value = true
}

async function handleCertSubmit() {
  certSubmitting.value = true
  try {
    const payload = { ...certForm }
    if (isEditCert.value) {
      await updateSeafarerCert(certForm.id, payload)
    } else {
      await createSeafarerCert(payload)
    }
    ElMessage.success(isEditCert.value ? '更新成功' : '创建成功')
    certDialogVisible.value = false
    loadCerts()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    certSubmitting.value = false
  }
}

async function handleDeleteCert(row) {
  try {
    await ElMessageBox.confirm('确定删除该证书？', '提示', { type: 'warning' })
    await deleteSeafarerCert(row.id)
    ElMessage.success('删除成功')
    loadCerts()
  } catch {
    // cancelled
  }
}

function openCertTypeDialog() {
  resetCertTypeForm()
  certTypeDialogVisible.value = true
}

async function handleCertTypeSubmit() {
  certTypeSubmitting.value = true
  try {
    await createCertType({ ...certTypeForm })
    ElMessage.success('创建成功')
    certTypeDialogVisible.value = false
    loadCertTypes()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    certTypeSubmitting.value = false
  }
}

async function handleDeleteCertType(row) {
  try {
    await ElMessageBox.confirm(`确定删除证书类型「${row.name}」？`, '提示', { type: 'warning' })
    await deleteCertType(row.id)
    ElMessage.success('删除成功')
    loadCertTypes()
  } catch {
    // cancelled
  }
}

onMounted(() => {
  loadSeafarers()
  loadCertTypes()
})
</script>
