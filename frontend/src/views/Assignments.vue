<template>
  <div>
    <div style="display:flex;gap:12px;margin-bottom:16px;flex-wrap:wrap;align-items:center">
      <el-select v-model="filters.ship_id" placeholder="选择船舶" filterable clearable style="width:180px">
        <el-option v-for="s in ships" :key="s.id" :label="s.name" :value="s.id" />
      </el-select>
      <el-select v-model="filters.seafarer_id" placeholder="选择船员" filterable clearable style="width:180px">
        <el-option v-for="s in seafarers" :key="s.id" :label="s.name" :value="s.id" />
      </el-select>
      <el-select v-model="filters.status" placeholder="状态" clearable style="width:140px">
        <el-option label="在船" :value="1" />
        <el-option label="已下船" :value="0" />
      </el-select>
      <el-button type="primary" @click="loadData">搜索</el-button>
      <div style="flex:1" />
      <el-button type="primary" @click="openAddDialog()">新增分配</el-button>
    </div>

    <el-table :data="tableData" v-loading="loading" border stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column label="船员名" width="120">
        <template #default="{ row }">{{ row.seafarer?.name || '' }}</template>
      </el-table-column>
      <el-table-column label="船名" width="140">
        <template #default="{ row }">{{ row.ship?.name || '' }}</template>
      </el-table-column>
      <el-table-column label="岗位" width="140">
        <template #default="{ row }">{{ row.ship_position?.position_name || '' }}</template>
      </el-table-column>
      <el-table-column label="上船日期" width="120">
        <template #default="{ row }">{{ row.embark_date?.slice(0, 10) }}</template>
      </el-table-column>
      <el-table-column label="预计下船日期" width="130">
        <template #default="{ row }">{{ row.expected_disembark_date?.slice(0, 10) || '-' }}</template>
      </el-table-column>
      <el-table-column label="实际下船日期" width="130">
        <template #default="{ row }">{{ row.actual_disembark_date?.slice(0, 10) || '-' }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? '在船' : '已下船' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.status === 1" link type="warning" @click="openDisembarkDialog(row)">下船</el-button>
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

    <el-dialog v-model="addDialogVisible" title="新增分配" width="520px" destroy-on-close>
      <el-form :model="addForm" label-width="100px">
        <el-form-item label="船员">
          <el-select v-model="addForm.seafarer_id" filterable placeholder="选择船员" style="width:100%">
            <el-option v-for="s in seafarers" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="船舶">
          <el-select v-model="addForm.ship_id" filterable placeholder="选择船舶" style="width:100%" @change="handleShipChange">
            <el-option v-for="s in ships" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="岗位">
          <el-select v-model="addForm.ship_position_id" filterable placeholder="选择岗位" style="width:100%">
            <el-option v-for="p in shipPositions" :key="p.id" :label="p.position_name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="上船日期">
          <el-date-picker v-model="addForm.embark_date" type="date" value-format="YYYY-MM-DD" style="width:100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="addSubmitting" @click="handleAddSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="disembarkDialogVisible" title="下船操作" width="420px" destroy-on-close>
      <el-form :model="disembarkForm" label-width="100px">
        <el-form-item label="实际下船日期">
          <el-date-picker v-model="disembarkForm.actual_date" type="date" value-format="YYYY-MM-DD" style="width:100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="disembarkDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="disembarkSubmitting" @click="handleDisembarkSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getSeafarers, getShips, getShipPositions, getAssignments, createAssignment, disembarkAssignment } from '../api/modules.js'

const seafarers = ref([])
const ships = ref([])
const shipPositions = ref([])

const filters = reactive({ ship_id: null, seafarer_id: null, status: null })
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const tableData = ref([])
const loading = ref(false)

const addDialogVisible = ref(false)
const addSubmitting = ref(false)
const addForm = reactive({
  seafarer_id: null,
  ship_id: null,
  ship_position_id: null,
  embark_date: '',
})

const disembarkDialogVisible = ref(false)
const disembarkSubmitting = ref(false)
const disembarkForm = reactive({ id: null, actual_date: '' })

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
    if (filters.ship_id) params.ship_id = filters.ship_id
    if (filters.seafarer_id) params.seafarer_id = filters.seafarer_id
    if (filters.status !== null && filters.status !== '') params.status = filters.status
    const res = await getAssignments(params)
    tableData.value = res.data.items || []
    total.value = res.data.total || 0
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function handleShipChange(shipId) {
  addForm.ship_position_id = null
  if (!shipId) {
    shipPositions.value = []
    return
  }
  try {
    const res = await getShipPositions({ ship_id: shipId })
    shipPositions.value = res.data || []
  } catch (e) {
    shipPositions.value = []
  }
}

function openAddDialog() {
  addForm.seafarer_id = null
  addForm.ship_id = null
  addForm.ship_position_id = null
  addForm.embark_date = ''
  shipPositions.value = []
  addDialogVisible.value = true
}

async function handleAddSubmit() {
  addSubmitting.value = true
  try {
    await createAssignment({ ...addForm })
    ElMessage.success('创建成功')
    addDialogVisible.value = false
    loadData()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    addSubmitting.value = false
  }
}

function openDisembarkDialog(row) {
  disembarkForm.id = row.id
  disembarkForm.actual_date = ''
  disembarkDialogVisible.value = true
}

async function handleDisembarkSubmit() {
  disembarkSubmitting.value = true
  try {
    await disembarkAssignment(disembarkForm.id, { actual_date: disembarkForm.actual_date })
    ElMessage.success('下船成功')
    disembarkDialogVisible.value = false
    loadData()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    disembarkSubmitting.value = false
  }
}

onMounted(() => {
  loadOptions()
  loadData()
})
</script>
