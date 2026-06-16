<template>
  <div>
    <el-tabs v-model="activeTab">
      <el-tab-pane label="航次合同" name="contracts">
        <div style="display:flex;gap:12px;margin-bottom:16px;align-items:center">
          <el-select v-model="contractFilter.seafarer_id" placeholder="选择船员" filterable clearable style="width:200px">
            <el-option v-for="s in seafarers" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
          <el-button type="primary" @click="loadContracts">搜索</el-button>
          <div style="flex:1" />
          <el-button type="primary" @click="openContractDialog()">新增合同</el-button>
        </div>
        <el-table :data="contracts" v-loading="contractLoading" border stripe>
          <el-table-column prop="contract_number" label="合同编号" width="160" />
          <el-table-column label="船员" width="100">
            <template #default="{ row }">{{ row.seafarer?.name || '' }}</template>
          </el-table-column>
          <el-table-column label="船舶" width="140">
            <template #default="{ row }">{{ row.ship?.name || '' }}</template>
          </el-table-column>
          <el-table-column label="开始日期" width="120">
            <template #default="{ row }">{{ row.start_date?.slice(0, 10) }}</template>
          </el-table-column>
          <el-table-column label="结束日期" width="120">
            <template #default="{ row }">{{ row.end_date?.slice(0, 10) }}</template>
          </el-table-column>
          <el-table-column label="实际结束日期" width="130">
            <template #default="{ row }">{{ row.actual_end_date?.slice(0, 10) || '-' }}</template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 0 ? 'danger' : row.status === 1 ? 'success' : 'info'">
                {{ row.status === 0 ? '终止' : row.status === 1 ? '执行中' : '已完成' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="上下船记录" name="embark">
        <div style="display:flex;gap:12px;margin-bottom:16px;align-items:center">
          <el-select v-model="embarkFilter.seafarer_id" placeholder="选择船员" filterable clearable style="width:200px">
            <el-option v-for="s in seafarers" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
          <el-button type="primary" @click="loadEmbarkRecords">搜索</el-button>
          <div style="flex:1" />
          <el-button type="primary" @click="openEmbarkDialog()">新增记录</el-button>
        </div>
        <el-table :data="embarkRecords" v-loading="embarkLoading" border stripe>
          <el-table-column label="船员" width="120">
            <template #default="{ row }">{{ row.seafarer?.name || '' }}</template>
          </el-table-column>
          <el-table-column label="船舶" width="140">
            <template #default="{ row }">{{ row.ship?.name || '' }}</template>
          </el-table-column>
          <el-table-column label="类型" width="100">
            <template #default="{ row }">
              <el-tag :type="row.record_type === 1 ? 'success' : 'warning'">
                {{ row.record_type === 1 ? '上船' : '下船' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="日期" width="120">
            <template #default="{ row }">{{ row.record_date?.slice(0, 10) }}</template>
          </el-table-column>
          <el-table-column prop="port" label="港口" width="140" />
          <el-table-column prop="reason" label="原因" />
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="休假记录" name="leave">
        <div style="display:flex;gap:12px;margin-bottom:16px;align-items:center">
          <el-select v-model="leaveFilter.seafarer_id" placeholder="选择船员" filterable clearable style="width:200px">
            <el-option v-for="s in seafarers" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
          <el-button type="primary" @click="loadLeaveRecords">搜索</el-button>
          <div style="flex:1" />
          <el-button type="primary" @click="openLeaveDialog()">新增休假</el-button>
        </div>
        <el-table :data="leaveRecords" v-loading="leaveLoading" border stripe>
          <el-table-column label="船员" width="120">
            <template #default="{ row }">{{ row.seafarer?.name || '' }}</template>
          </el-table-column>
          <el-table-column label="开始日期" width="120">
            <template #default="{ row }">{{ row.start_date?.slice(0, 10) }}</template>
          </el-table-column>
          <el-table-column label="结束日期" width="120">
            <template #default="{ row }">{{ row.end_date?.slice(0, 10) || '-' }}</template>
          </el-table-column>
          <el-table-column label="休假天数" width="100">
            <template #default="{ row }">{{ row.leave_days ?? '-' }}</template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 0 ? 'info' : 'success'">
                {{ row.status === 0 ? '已结束' : '休假中' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100">
            <template #default="{ row }">
              <el-button v-if="row.status === 1" link type="warning" @click="handleEndLeave(row)">结束休假</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="健康复检" name="health">
        <div style="display:flex;gap:12px;margin-bottom:16px;align-items:center">
          <el-select v-model="healthFilter.seafarer_id" placeholder="选择船员" filterable clearable style="width:200px">
            <el-option v-for="s in seafarers" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
          <el-button type="primary" @click="loadHealthRecords">搜索</el-button>
          <div style="flex:1" />
          <el-button type="primary" @click="openHealthDialog()">新增记录</el-button>
        </div>
        <el-table :data="healthRecords" v-loading="healthLoading" border stripe>
          <el-table-column label="船员" width="120">
            <template #default="{ row }">{{ row.seafarer?.name || '' }}</template>
          </el-table-column>
          <el-table-column label="检查日期" width="120">
            <template #default="{ row }">{{ row.exam_date?.slice(0, 10) }}</template>
          </el-table-column>
          <el-table-column label="下次检查日期" width="140">
            <template #default="{ row }">{{ row.next_exam_date?.slice(0, 10) || '-' }}</template>
          </el-table-column>
          <el-table-column label="结果" width="100">
            <template #default="{ row }">
              <el-tag :type="row.exam_result === 1 ? 'success' : row.exam_result === 2 ? 'danger' : 'warning'">
                {{ row.exam_result === 1 ? '合格' : row.exam_result === 2 ? '不合格' : '限制' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="exam_institution" label="检查机构" />
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="contractDialogVisible" title="新增合同" width="560px" destroy-on-close>
      <el-form :model="contractForm" label-width="100px">
        <el-form-item label="合同编号">
          <el-input v-model="contractForm.contract_number" />
        </el-form-item>
        <el-form-item label="船员">
          <el-select v-model="contractForm.seafarer_id" filterable placeholder="选择船员" style="width:100%">
            <el-option v-for="s in seafarers" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="船舶">
          <el-select v-model="contractForm.ship_id" filterable placeholder="选择船舶" style="width:100%">
            <el-option v-for="s in ships" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="开始日期">
          <el-date-picker v-model="contractForm.start_date" type="date" value-format="YYYY-MM-DD" style="width:100%" />
        </el-form-item>
        <el-form-item label="结束日期">
          <el-date-picker v-model="contractForm.end_date" type="date" value-format="YYYY-MM-DD" style="width:100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="contractDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="contractSubmitting" @click="handleContractSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="embarkDialogVisible" title="新增上下船记录" width="520px" destroy-on-close>
      <el-form :model="embarkForm" label-width="100px">
        <el-form-item label="船员">
          <el-select v-model="embarkForm.seafarer_id" filterable placeholder="选择船员" style="width:100%">
            <el-option v-for="s in seafarers" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="船舶">
          <el-select v-model="embarkForm.ship_id" filterable placeholder="选择船舶" style="width:100%">
            <el-option v-for="s in ships" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="embarkForm.record_type" style="width:100%">
            <el-option label="上船" :value="1" />
            <el-option label="下船" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期">
          <el-date-picker v-model="embarkForm.record_date" type="date" value-format="YYYY-MM-DD" style="width:100%" />
        </el-form-item>
        <el-form-item label="港口">
          <el-input v-model="embarkForm.port" />
        </el-form-item>
        <el-form-item label="原因">
          <el-input v-model="embarkForm.reason" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="embarkDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="embarkSubmitting" @click="handleEmbarkSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="leaveDialogVisible" title="新增休假记录" width="480px" destroy-on-close>
      <el-form :model="leaveForm" label-width="100px">
        <el-form-item label="船员">
          <el-select v-model="leaveForm.seafarer_id" filterable placeholder="选择船员" style="width:100%">
            <el-option v-for="s in seafarers" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="开始日期">
          <el-date-picker v-model="leaveForm.start_date" type="date" value-format="YYYY-MM-DD" style="width:100%" />
        </el-form-item>
        <el-form-item label="原因">
          <el-input v-model="leaveForm.reason" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="leaveDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="leaveSubmitting" @click="handleLeaveSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="healthDialogVisible" title="新增健康复检" width="520px" destroy-on-close>
      <el-form :model="healthForm" label-width="100px">
        <el-form-item label="船员">
          <el-select v-model="healthForm.seafarer_id" filterable placeholder="选择船员" style="width:100%">
            <el-option v-for="s in seafarers" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="检查日期">
          <el-date-picker v-model="healthForm.exam_date" type="date" value-format="YYYY-MM-DD" style="width:100%" />
        </el-form-item>
        <el-form-item label="下次检查日期">
          <el-date-picker v-model="healthForm.next_exam_date" type="date" value-format="YYYY-MM-DD" style="width:100%" />
        </el-form-item>
        <el-form-item label="结果">
          <el-select v-model="healthForm.exam_result" style="width:100%">
            <el-option label="合格" :value="1" />
            <el-option label="不合格" :value="2" />
            <el-option label="限制" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="检查机构">
          <el-input v-model="healthForm.exam_institution" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="healthDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="healthSubmitting" @click="handleHealthSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getSeafarers, getShips, getContracts, createContract, getEmbarkRecords, createEmbarkRecord, getLeaveRecords, createLeaveRecord, endLeave, getHealthRecords, createHealthRecord } from '../api/modules.js'

const activeTab = ref('contracts')

const seafarers = ref([])
const ships = ref([])

const contractFilter = reactive({ seafarer_id: null })
const contracts = ref([])
const contractLoading = ref(false)
const contractDialogVisible = ref(false)
const contractSubmitting = ref(false)
const contractForm = reactive({
  contract_number: '',
  seafarer_id: null,
  ship_id: null,
  start_date: '',
  end_date: '',
})

const embarkFilter = reactive({ seafarer_id: null })
const embarkRecords = ref([])
const embarkLoading = ref(false)
const embarkDialogVisible = ref(false)
const embarkSubmitting = ref(false)
const embarkForm = reactive({
  seafarer_id: null,
  ship_id: null,
  record_type: 1,
  record_date: '',
  port: '',
  reason: '',
})

const leaveFilter = reactive({ seafarer_id: null })
const leaveRecords = ref([])
const leaveLoading = ref(false)
const leaveDialogVisible = ref(false)
const leaveSubmitting = ref(false)
const leaveForm = reactive({
  seafarer_id: null,
  start_date: '',
  reason: '',
})

const healthFilter = reactive({ seafarer_id: null })
const healthRecords = ref([])
const healthLoading = ref(false)
const healthDialogVisible = ref(false)
const healthSubmitting = ref(false)
const healthForm = reactive({
  seafarer_id: null,
  exam_date: '',
  next_exam_date: '',
  exam_result: 1,
  exam_institution: '',
})

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

async function loadContracts() {
  contractLoading.value = true
  try {
    const params = {}
    if (contractFilter.seafarer_id) params.seafarer_id = contractFilter.seafarer_id
    const res = await getContracts(params)
    contracts.value = res.data?.items || res.data || []
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    contractLoading.value = false
  }
}

function openContractDialog() {
  contractForm.contract_number = ''
  contractForm.seafarer_id = null
  contractForm.ship_id = null
  contractForm.start_date = ''
  contractForm.end_date = ''
  contractDialogVisible.value = true
}

async function handleContractSubmit() {
  contractSubmitting.value = true
  try {
    await createContract({ ...contractForm })
    ElMessage.success('创建成功')
    contractDialogVisible.value = false
    loadContracts()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    contractSubmitting.value = false
  }
}

async function loadEmbarkRecords() {
  embarkLoading.value = true
  try {
    const params = {}
    if (embarkFilter.seafarer_id) params.seafarer_id = embarkFilter.seafarer_id
    const res = await getEmbarkRecords(params)
    embarkRecords.value = res.data?.items || res.data || []
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    embarkLoading.value = false
  }
}

function openEmbarkDialog() {
  embarkForm.seafarer_id = null
  embarkForm.ship_id = null
  embarkForm.record_type = 1
  embarkForm.record_date = ''
  embarkForm.port = ''
  embarkForm.reason = ''
  embarkDialogVisible.value = true
}

async function handleEmbarkSubmit() {
  embarkSubmitting.value = true
  try {
    await createEmbarkRecord({ ...embarkForm })
    ElMessage.success('创建成功')
    embarkDialogVisible.value = false
    loadEmbarkRecords()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    embarkSubmitting.value = false
  }
}

async function loadLeaveRecords() {
  leaveLoading.value = true
  try {
    const params = {}
    if (leaveFilter.seafarer_id) params.seafarer_id = leaveFilter.seafarer_id
    const res = await getLeaveRecords(params)
    leaveRecords.value = res.data?.items || res.data || []
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    leaveLoading.value = false
  }
}

function openLeaveDialog() {
  leaveForm.seafarer_id = null
  leaveForm.start_date = ''
  leaveForm.reason = ''
  leaveDialogVisible.value = true
}

async function handleLeaveSubmit() {
  leaveSubmitting.value = true
  try {
    await createLeaveRecord({ ...leaveForm })
    ElMessage.success('创建成功')
    leaveDialogVisible.value = false
    loadLeaveRecords()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    leaveSubmitting.value = false
  }
}

async function handleEndLeave(row) {
  try {
    await ElMessageBox.confirm('确定结束该休假？', '提示', { type: 'warning' })
    await endLeave(row.id)
    ElMessage.success('已结束休假')
    loadLeaveRecords()
  } catch {
    // cancelled
  }
}

async function loadHealthRecords() {
  healthLoading.value = true
  try {
    const params = {}
    if (healthFilter.seafarer_id) params.seafarer_id = healthFilter.seafarer_id
    const res = await getHealthRecords(params)
    healthRecords.value = res.data?.items || res.data || []
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    healthLoading.value = false
  }
}

function openHealthDialog() {
  healthForm.seafarer_id = null
  healthForm.exam_date = ''
  healthForm.next_exam_date = ''
  healthForm.exam_result = 1
  healthForm.exam_institution = ''
  healthDialogVisible.value = true
}

async function handleHealthSubmit() {
  healthSubmitting.value = true
  try {
    await createHealthRecord({ ...healthForm })
    ElMessage.success('创建成功')
    healthDialogVisible.value = false
    loadHealthRecords()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    healthSubmitting.value = false
  }
}

onMounted(() => {
  loadOptions()
  loadContracts()
  loadEmbarkRecords()
  loadLeaveRecords()
  loadHealthRecords()
})
</script>
