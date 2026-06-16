<template>
  <div>
    <div style="display:flex;justify-content:space-between;margin-bottom:16px">
      <el-input v-model="keyword" placeholder="搜索船名/IMO号" style="width:300px" clearable @clear="loadShips" @keyup.enter="loadShips" />
      <el-button type="primary" @click="openShipDialog()">新增船舶</el-button>
    </div>

    <el-table :data="shipData" v-loading="shipLoading" border stripe @row-click="handleRowClick" highlight-current-row>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="船名" width="140" />
      <el-table-column prop="imo_number" label="IMO号" width="100" />
      <el-table-column prop="mmsi" label="MMSI" width="120" />
      <el-table-column prop="ship_type" label="船舶类型" width="120" />
      <el-table-column label="总吨位" width="100">
        <template #default="{ row }">{{ row.gross_tonnage ?? '-' }}</template>
      </el-table-column>
      <el-table-column prop="flag_state" label="船旗国" width="100" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 0 ? 'danger' : row.status === 1 ? 'success' : 'warning'">
            {{ row.status === 0 ? '报废' : row.status === 1 ? '运营' : '维修' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click.stop="openShipDialog(row)">编辑</el-button>
          <el-button link type="danger" @click.stop="handleDeleteShip(row)">删除</el-button>
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
        @size-change="loadShips"
        @current-change="loadShips"
      />
    </div>

    <el-card v-if="selectedShip" style="margin-top:20px">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span>{{ selectedShip.name }} - 岗位配置</span>
          <el-button type="primary" size="small" @click="openPositionDialog()">添加岗位</el-button>
        </div>
      </template>
      <el-table :data="positions" v-loading="posLoading" border stripe>
        <el-table-column prop="position_name" label="岗位名" width="140" />
        <el-table-column prop="department" label="部门" width="120" />
        <el-table-column prop="required_count" label="编制人数" width="100" />
        <el-table-column label="在船人数" width="100">
          <template #default="{ row }">{{ row.assignments?.filter(a => a.status === 1).length || 0 }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button link type="danger" @click="handleDeletePosition(row)">删除岗位</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="shipDialogVisible" :title="isEditShip ? '编辑船舶' : '新增船舶'" width="560px" destroy-on-close>
      <el-form :model="shipForm" label-width="90px">
        <el-form-item label="船名">
          <el-input v-model="shipForm.name" />
        </el-form-item>
        <el-form-item label="IMO号">
          <el-input v-model="shipForm.imo_number" />
        </el-form-item>
        <el-form-item label="MMSI">
          <el-input v-model="shipForm.mmsi" />
        </el-form-item>
        <el-form-item label="船舶类型">
          <el-input v-model="shipForm.ship_type" />
        </el-form-item>
        <el-form-item label="总吨位">
          <el-input-number v-model="shipForm.gross_tonnage" :precision="2" :min="0" style="width:100%" />
        </el-form-item>
        <el-form-item label="船旗国">
          <el-input v-model="shipForm.flag_state" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="shipForm.status" style="width:100%">
            <el-option label="报废" :value="0" />
            <el-option label="运营" :value="1" />
            <el-option label="维修" :value="2" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="shipDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="shipSubmitting" @click="handleShipSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="posDialogVisible" title="添加岗位" width="480px" destroy-on-close>
      <el-form :model="posForm" label-width="90px">
        <el-form-item label="岗位名">
          <el-input v-model="posForm.position_name" />
        </el-form-item>
        <el-form-item label="部门">
          <el-input v-model="posForm.department" />
        </el-form-item>
        <el-form-item label="编制人数">
          <el-input-number v-model="posForm.required_count" :min="1" style="width:100%" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="posForm.sort_order" :min="0" style="width:100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="posDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="posSubmitting" @click="handlePosSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getShips, createShip, updateShip, deleteShip, getShipPositions, createShipPosition, deleteShipPosition } from '../api/modules.js'

const keyword = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const shipData = ref([])
const shipLoading = ref(false)

const selectedShip = ref(null)
const positions = ref([])
const posLoading = ref(false)

const shipDialogVisible = ref(false)
const isEditShip = ref(false)
const shipSubmitting = ref(false)
const shipForm = reactive({
  id: null,
  name: '',
  imo_number: '',
  mmsi: '',
  ship_type: '',
  gross_tonnage: null,
  flag_state: '',
  status: 1,
})

const posDialogVisible = ref(false)
const posSubmitting = ref(false)
const posForm = reactive({
  ship_id: null,
  position_name: '',
  department: '',
  required_count: 1,
  sort_order: 0,
})

function resetShipForm() {
  shipForm.id = null
  shipForm.name = ''
  shipForm.imo_number = ''
  shipForm.mmsi = ''
  shipForm.ship_type = ''
  shipForm.gross_tonnage = null
  shipForm.flag_state = ''
  shipForm.status = 1
}

function resetPosForm() {
  posForm.ship_id = null
  posForm.position_name = ''
  posForm.department = ''
  posForm.required_count = 1
  posForm.sort_order = 0
}

async function loadShips() {
  shipLoading.value = true
  try {
    const res = await getShips({ page: page.value, page_size: pageSize.value, keyword: keyword.value })
    shipData.value = res.data.items || []
    total.value = res.data.total || 0
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    shipLoading.value = false
  }
}

async function loadPositions() {
  if (!selectedShip.value) return
  posLoading.value = true
  try {
    const res = await getShipPositions({ ship_id: selectedShip.value.id })
    positions.value = res.data || []
  } catch (e) {
    ElMessage.error(e.message || '加载岗位失败')
  } finally {
    posLoading.value = false
  }
}

function handleRowClick(row) {
  selectedShip.value = row
  loadPositions()
}

function openShipDialog(row) {
  resetShipForm()
  if (row) {
    isEditShip.value = true
    Object.assign(shipForm, {
      id: row.id,
      name: row.name,
      imo_number: row.imo_number,
      mmsi: row.mmsi || '',
      ship_type: row.ship_type || '',
      gross_tonnage: row.gross_tonnage,
      flag_state: row.flag_state || '',
      status: row.status,
    })
  } else {
    isEditShip.value = false
  }
  shipDialogVisible.value = true
}

async function handleShipSubmit() {
  shipSubmitting.value = true
  try {
    const payload = { ...shipForm }
    if (isEditShip.value) {
      await updateShip(shipForm.id, payload)
    } else {
      await createShip(payload)
    }
    ElMessage.success(isEditShip.value ? '更新成功' : '创建成功')
    shipDialogVisible.value = false
    loadShips()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    shipSubmitting.value = false
  }
}

async function handleDeleteShip(row) {
  try {
    await ElMessageBox.confirm(`确定删除船舶「${row.name}」？`, '提示', { type: 'warning' })
    await deleteShip(row.id)
    ElMessage.success('删除成功')
    if (selectedShip.value?.id === row.id) {
      selectedShip.value = null
      positions.value = []
    }
    loadShips()
  } catch {
    // cancelled
  }
}

function openPositionDialog() {
  resetPosForm()
  posForm.ship_id = selectedShip.value.id
  posDialogVisible.value = true
}

async function handlePosSubmit() {
  posSubmitting.value = true
  try {
    await createShipPosition({ ...posForm })
    ElMessage.success('创建成功')
    posDialogVisible.value = false
    loadPositions()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    posSubmitting.value = false
  }
}

async function handleDeletePosition(row) {
  try {
    await ElMessageBox.confirm(`确定删除岗位「${row.position_name}」？`, '提示', { type: 'warning' })
    await deleteShipPosition(row.id)
    ElMessage.success('删除成功')
    loadPositions()
  } catch {
    // cancelled
  }
}

onMounted(() => loadShips())
</script>
