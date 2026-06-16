<template>
  <div>
    <div style="display:flex;gap:16px;margin-bottom:20px">
      <el-card shadow="hover" style="flex:1;text-align:center">
        <div style="font-size:14px;color:#909399">90天预警数</div>
        <div style="font-size:28px;font-weight:bold;color:#409eff;margin-top:8px">{{ stats.level_1 || 0 }}</div>
      </el-card>
      <el-card shadow="hover" style="flex:1;text-align:center">
        <div style="font-size:14px;color:#909399">60天预警数</div>
        <div style="font-size:28px;font-weight:bold;color:#e6a23c;margin-top:8px">{{ stats.level_2 || 0 }}</div>
      </el-card>
      <el-card shadow="hover" style="flex:1;text-align:center">
        <div style="font-size:14px;color:#909399">30天预警数</div>
        <div style="font-size:28px;font-weight:bold;color:#f56c6c;margin-top:8px">{{ stats.level_3 || 0 }}</div>
      </el-card>
      <el-card shadow="hover" style="flex:1;text-align:center">
        <div style="font-size:14px;color:#909399">未处理总数</div>
        <div style="font-size:28px;font-weight:bold;color:#f56c6c;margin-top:8px">{{ stats.total_unhandled || 0 }}</div>
      </el-card>
    </div>

    <div style="display:flex;gap:12px;margin-bottom:16px;align-items:center">
      <el-select v-model="filters.level" placeholder="预警级别" clearable style="width:140px">
        <el-option label="90天预警" :value="1" />
        <el-option label="60天预警" :value="2" />
        <el-option label="30天预警" :value="3" />
      </el-select>
      <el-select v-model="filters.is_handled" placeholder="处理状态" clearable style="width:140px">
        <el-option label="未处理" :value="0" />
        <el-option label="已处理" :value="1" />
      </el-select>
      <el-button type="primary" @click="loadData">搜索</el-button>
      <div style="flex:1" />
      <el-button type="warning" :loading="scanning" @click="handleScan">手动扫描</el-button>
    </div>

    <el-table :data="tableData" v-loading="loading" border stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column label="船员" width="100">
        <template #default="{ row }">{{ row.seafarer?.name || '' }}</template>
      </el-table-column>
      <el-table-column label="证书" width="140">
        <template #default="{ row }">{{ row.certificate?.cert_number || '' }}</template>
      </el-table-column>
      <el-table-column label="预警级别" width="110">
        <template #default="{ row }">
          <el-tag :type="row.alert_level === 1 ? 'info' : row.alert_level === 2 ? 'warning' : 'danger'">
            {{ row.alert_level === 1 ? '90天' : row.alert_level === 2 ? '60天' : '30天' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="预警日期" width="120">
        <template #default="{ row }">{{ row.alert_date?.slice(0, 10) }}</template>
      </el-table-column>
      <el-table-column label="到期日期" width="120">
        <template #default="{ row }">{{ row.expire_date?.slice(0, 10) }}</template>
      </el-table-column>
      <el-table-column prop="days_remaining" label="剩余天数" width="100" />
      <el-table-column label="是否已处理" width="110">
        <template #default="{ row }">
          <el-tag :type="row.is_handled ? 'success' : 'danger'">{{ row.is_handled ? '已处理' : '未处理' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button v-if="!row.is_handled" link type="primary" @click="openHandleDialog(row)">处理</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="handleDialogVisible" title="处理预警" width="440px" destroy-on-close>
      <el-form :model="handleForm" label-width="80px">
        <el-form-item label="备注">
          <el-input v-model="handleForm.remark" type="textarea" :rows="4" placeholder="请输入处理备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="handleDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="handleSubmitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getAlerts, getAlertStats, handleAlert, runAlertScan } from '../api/modules.js'

const stats = reactive({ level_1: 0, level_2: 0, level_3: 0, total_unhandled: 0 })
const filters = reactive({ level: null, is_handled: null })
const tableData = ref([])
const loading = ref(false)
const scanning = ref(false)

const handleDialogVisible = ref(false)
const handleSubmitting = ref(false)
const handleForm = reactive({ id: null, remark: '' })

async function loadStats() {
  try {
    const res = await getAlertStats()
    Object.assign(stats, res.data || {})
  } catch (e) {
    ElMessage.error(e.message || '加载统计失败')
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = {}
    if (filters.level) params.level = filters.level
    if (filters.is_handled !== null && filters.is_handled !== '') params.is_handled = filters.is_handled
    const res = await getAlerts(params)
    tableData.value = res.data?.items || res.data || []
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function handleScan() {
  scanning.value = true
  try {
    const res = await runAlertScan()
    const count = res.data?.count || 0
    ElMessage.success(`扫描完成，新增预警${count}条`)
    loadStats()
    loadData()
  } catch (e) {
    ElMessage.error(e.message || '扫描失败')
  } finally {
    scanning.value = false
  }
}

function openHandleDialog(row) {
  handleForm.id = row.id
  handleForm.remark = ''
  handleDialogVisible.value = true
}

async function handleSubmit() {
  handleSubmitting.value = true
  try {
    await handleAlert(handleForm.id, { remark: handleForm.remark })
    ElMessage.success('处理成功')
    handleDialogVisible.value = false
    loadData()
    loadStats()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    handleSubmitting.value = false
  }
}

onMounted(() => {
  loadStats()
  loadData()
})
</script>
