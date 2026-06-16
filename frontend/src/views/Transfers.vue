<template>
  <div>
    <div style="display:flex;gap:12px;margin-bottom:16px;align-items:center">
      <el-select v-model="statusFilter" placeholder="状态筛选" clearable style="width:160px">
        <el-option label="待审批" :value="0" />
        <el-option label="已批准" :value="1" />
        <el-option label="已拒绝" :value="2" />
        <el-option label="已取消" :value="3" />
      </el-select>
      <el-button type="primary" @click="loadData">搜索</el-button>
      <div style="flex:1" />
      <el-button type="primary" @click="openCreateDialog()">新增调动</el-button>
    </div>

    <el-table :data="tableData" v-loading="loading" border stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column label="船员" width="100">
        <template #default="{ row }">{{ row.seafarer?.name || '' }}</template>
      </el-table-column>
      <el-table-column label="原船" width="120">
        <template #default="{ row }">{{ row.from_ship?.name || '' }}</template>
      </el-table-column>
      <el-table-column label="目标船" width="120">
        <template #default="{ row }">{{ row.to_ship?.name || '' }}</template>
      </el-table-column>
      <el-table-column label="原岗位" width="120">
        <template #default="{ row }">{{ row.from_position?.position_name || '' }}</template>
      </el-table-column>
      <el-table-column label="目标岗位" width="120">
        <template #default="{ row }">{{ row.to_position?.position_name || '' }}</template>
      </el-table-column>
      <el-table-column label="替换船员ID" width="110">
        <template #default="{ row }">{{ row.replacement_seafarer_id ?? '-' }}</template>
      </el-table-column>
      <el-table-column label="原船校验" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.from_ship_valid === 1" type="success">通过</el-tag>
          <el-tag v-else-if="row.from_ship_valid === 0" type="danger">未通过</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="目标船校验" width="110">
        <template #default="{ row }">
          <el-tag v-if="row.to_ship_valid === 1" type="success">通过</el-tag>
          <el-tag v-else-if="row.to_ship_valid === 0" type="danger">未通过</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.status === 0" link type="success" @click="openApproveDialog(row)">审批</el-button>
          <el-button v-if="row.status === 0" link type="danger" @click="openRejectDialog(row)">拒绝</el-button>
          <el-button v-if="row.status === 0" link type="info" @click="handleCancel(row)">取消</el-button>
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

    <el-dialog v-model="createDialogVisible" title="新增调动" width="600px" destroy-on-close>
      <el-form :model="createForm" label-width="110px">
        <el-form-item label="船员">
          <el-select v-model="createForm.seafarer_id" filterable placeholder="选择船员" style="width:100%">
            <el-option v-for="s in seafarers" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="原船">
          <el-select v-model="createForm.from_ship_id" filterable placeholder="选择原船" style="width:100%" @change="handleFromShipChange">
            <el-option v-for="s in ships" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="原岗位">
          <el-select v-model="createForm.from_position_id" filterable placeholder="选择原岗位" style="width:100%">
            <el-option v-for="p in fromPositions" :key="p.id" :label="p.position_name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标船">
          <el-select v-model="createForm.to_ship_id" filterable placeholder="选择目标船" style="width:100%" @change="handleToShipChange">
            <el-option v-for="s in ships" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标岗位">
          <el-select v-model="createForm.to_position_id" filterable placeholder="选择目标岗位" style="width:100%">
            <el-option v-for="p in toPositions" :key="p.id" :label="p.position_name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="替换船员ID">
          <el-select v-model="createForm.replacement_seafarer_id" filterable clearable placeholder="可选" style="width:100%">
            <el-option v-for="s in seafarers" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="原因">
          <el-input v-model="createForm.reason" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="createSubmitting" @click="handleCreateSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="approveDialogVisible" title="审批调动" width="440px" destroy-on-close>
      <el-form :model="approveForm" label-width="80px">
        <el-form-item label="审批人">
          <el-input v-model="approveForm.approver" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="approveForm.remark" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="approveDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="approveSubmitting" @click="handleApproveSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="rejectDialogVisible" title="拒绝调动" width="440px" destroy-on-close>
      <el-form :model="rejectForm" label-width="80px">
        <el-form-item label="审批人">
          <el-input v-model="rejectForm.approver" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="rejectForm.remark" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectDialogVisible = false">取消</el-button>
        <el-button type="danger" :loading="rejectSubmitting" @click="handleRejectSubmit">确定拒绝</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getSeafarers, getShips, getShipPositions, getTransfers, createTransfer, approveTransfer, rejectTransfer, cancelTransfer } from '../api/modules.js'

const statusFilter = ref(null)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const tableData = ref([])
const loading = ref(false)

const seafarers = ref([])
const ships = ref([])
const fromPositions = ref([])
const toPositions = ref([])

const createDialogVisible = ref(false)
const createSubmitting = ref(false)
const createForm = reactive({
  seafarer_id: null,
  from_ship_id: null,
  from_position_id: null,
  to_ship_id: null,
  to_position_id: null,
  replacement_seafarer_id: null,
  reason: '',
})

const approveDialogVisible = ref(false)
const approveSubmitting = ref(false)
const approveForm = reactive({ id: null, approver: '', remark: '' })

const rejectDialogVisible = ref(false)
const rejectSubmitting = ref(false)
const rejectForm = reactive({ id: null, approver: '', remark: '' })

function statusText(s) {
  return s === 0 ? '待审批' : s === 1 ? '已批准' : s === 2 ? '已拒绝' : '已取消'
}
function statusTagType(s) {
  return s === 0 ? 'info' : s === 1 ? 'success' : s === 2 ? 'danger' : 'info'
}

async function loadOptions() {
  try {
    const [sRes, shRes] = await Promise.all([
      getSeafarers({ page: 1, page_size: 1000 }),
      getShips({ page: 1, page_size: 1000 }),
    ])
    seafarers.value = sRes.data.items || []
    ships.value = shRes.data.items || []
  } catch (e) {
    ElMessage.error(e.message || '加载选项失败')
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (statusFilter.value !== null && statusFilter.value !== '') params.status = statusFilter.value
    const res = await getTransfers(params)
    tableData.value = res.data.items || []
    total.value = res.data.total || 0
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function loadPositionsForShip(shipId, target) {
  try {
    const res = await getShipPositions({ ship_id: shipId })
    if (target === 'from') fromPositions.value = res.data || []
    else toPositions.value = res.data || []
  } catch {
    if (target === 'from') fromPositions.value = []
    else toPositions.value = []
  }
}

function handleFromShipChange(shipId) {
  createForm.from_position_id = null
  if (!shipId) { fromPositions.value = []; return }
  loadPositionsForShip(shipId, 'from')
}

function handleToShipChange(shipId) {
  createForm.to_position_id = null
  if (!shipId) { toPositions.value = []; return }
  loadPositionsForShip(shipId, 'to')
}

function openCreateDialog() {
  createForm.seafarer_id = null
  createForm.from_ship_id = null
  createForm.from_position_id = null
  createForm.to_ship_id = null
  createForm.to_position_id = null
  createForm.replacement_seafarer_id = null
  createForm.reason = ''
  fromPositions.value = []
  toPositions.value = []
  createDialogVisible.value = true
}

async function handleCreateSubmit() {
  createSubmitting.value = true
  try {
    const payload = { ...createForm }
    if (!payload.replacement_seafarer_id) delete payload.replacement_seafarer_id
    await createTransfer(payload)
    ElMessage.success('创建成功')
    createDialogVisible.value = false
    loadData()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    createSubmitting.value = false
  }
}

function openApproveDialog(row) {
  approveForm.id = row.id
  approveForm.approver = ''
  approveForm.remark = ''
  approveDialogVisible.value = true
}

async function handleApproveSubmit() {
  approveSubmitting.value = true
  try {
    await approveTransfer(approveForm.id, { approver: approveForm.approver, remark: approveForm.remark })
    ElMessage.success('审批成功')
    approveDialogVisible.value = false
    loadData()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    approveSubmitting.value = false
  }
}

function openRejectDialog(row) {
  rejectForm.id = row.id
  rejectForm.approver = ''
  rejectForm.remark = ''
  rejectDialogVisible.value = true
}

async function handleRejectSubmit() {
  rejectSubmitting.value = true
  try {
    await rejectTransfer(rejectForm.id, { approver: rejectForm.approver, remark: rejectForm.remark })
    ElMessage.success('已拒绝')
    rejectDialogVisible.value = false
    loadData()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    rejectSubmitting.value = false
  }
}

async function handleCancel(row) {
  try {
    await ElMessageBox.confirm('确定取消该调动申请？', '提示', { type: 'warning' })
    await cancelTransfer(row.id)
    ElMessage.success('已取消')
    loadData()
  } catch {
    // cancelled
  }
}

onMounted(() => {
  loadOptions()
  loadData()
})
</script>
