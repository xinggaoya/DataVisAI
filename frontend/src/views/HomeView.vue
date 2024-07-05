<template>
  <div class="main">
    <div class="main-content">
      <el-card style="flex: 1;margin-right: 10px" shadow="never">
        <el-form ref="formRef" :model="form" label-width="auto">
          <el-form-item label="地址" prop="address">
            <el-input type="text" v-model="form.address" placeholder="请输入地址"/>
          </el-form-item>
          <el-form-item label="端口" prop="port">
            <el-input type="number" v-model="form.port" placeholder="请输入端口"/>
          </el-form-item>
          <el-form-item label="账户" prop="account">
            <el-input type="text" v-model="form.account" placeholder="请输入账户"/>
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input type="password" v-model="form.password" placeholder="请输入密码"/>
          </el-form-item>
          <el-form-item>
            <div style="text-align: right;width: 100%">
              <el-button type="success" @click="test" :loading="testLoading">获取数据库</el-button>
            </div>
          </el-form-item>
          <el-form-item label="数据库" prop="database">
            <el-select v-model="form.database" @change="getAllTables"
                       placeholder="请选择数据库">
              <el-option
                  v-for="item in options"
                  :key="item"
                  :label="item"
                  :value="item"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="数据表" prop="table">
            <el-select v-model="form.table" multiple @change="getTableStructure" placeholder="请选择数据表">
              <el-option
                  v-for="item in tableList"
                  :key="item"
                  :label="item"
                  :value="item"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="问题" prop="info">
            <el-input type="textarea" v-model="form.info" placeholder="请输入信息"/>
          </el-form-item>
          <el-form-item>
            <el-alert type="info" title="仅支持通义千问APiKey"/>
          </el-form-item>
          <el-form-item label="apiKey" prop="apiKey">
            <el-input type="password" v-model="form.apiKey" placeholder="请输入通义千问apiKey"/>
          </el-form-item>
          <el-form-item>
            <div style="text-align: right;width: 100%">
              <el-button type="danger" @click="clearCache">清空缓存</el-button>
              <el-button type="primary" @click="getAIAnswer" :loading>执行</el-button>
            </div>
          </el-form-item>
        </el-form>
      </el-card>
      <el-card style="flex: 2" shadow="never">
        <template v-slot:default style="height: 100%">
          <div v-for="(item,index) in columnList" :key="index" style="margin-bottom: 7px;">
            <el-scrollbar>
              <el-tag type="primary" style="padding: 20px 10px">{{ item }}</el-tag>
            </el-scrollbar>
          </div>
          <el-divider/>
          <div style="height: 100%">
            <h3>AI生成结果</h3>
            <el-scrollbar>
              <el-text class="language-text">{{ aiInfo }}</el-text>
            </el-scrollbar>
          </div>
        </template>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import {onMounted, reactive, ref} from 'vue'
import {GetAIAnswer, GetAllDatabases, GetAllTables, GetTableStructure, TestDB} from '@api/internal/App'
import {ElMessage, ElMessageBox} from "element-plus";

const form = reactive({
  address: '',
  port: '',
  account: '',
  password: '',
  database: '',
  table: [],
  info: '',
  apiKey: '',
})

const database = 'databaseInfo'
const key = 'keyInfo'

const loading = ref(false)
const testLoading = ref(false)
const aiInfo = ref<string>('')
const options = ref<any>([])
const columnList = ref<any>({})
const tableList = ref<any>([])

onMounted(() => {
  if (localStorage.getItem(database)) {
    const data = JSON.parse(localStorage.getItem(database) as string)
    Object.assign(form, data)
  }
  if (localStorage.getItem(key)) {
    form.apiKey = localStorage.getItem(key) as string
  }
})


function getAllDatabases() {
  const params = {
    address: form.address,
    port: form.port,
    account: form.account,
    password: form.password
  }
  GetAllDatabases(
      params.address,
      params.port,
      params.account,
      params.password
  ).then(res => {
    options.value = res
    getAllTables()
    localStorage.setItem(database, JSON.stringify(params))
  }).catch(() => {
    ElMessage.error('数据库获取失败')
  })
}

// 获取回答 问题+表结构
function getAIAnswer() {
  if (form.info) {
    loading.value = true
    GetAIAnswer(
        form.apiKey,
        getQuestion(),
        form.info
    ).then(res => {
      aiInfo.value = res
      ElMessage.success("获取成功")
      localStorage.setItem(key, form.apiKey)
      loading.value = false
    })
  } else {
    ElMessage.error('请输入信息')
  }
}

function clearCache() {
  // 询问
  ElMessageBox.confirm('是否清空缓存', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(() => {
    localStorage.removeItem(database)
    localStorage.removeItem(key)
    ElMessage.success("清除成功,请重启应用！")
  }).catch(() => {
    ElMessage.info('已取消')
  })
}

// 生成结构问题
function getQuestion() {
  return JSON.stringify(columnList.value)
}

// 获取数据表
function getAllTables() {
  form.table = []
  GetAllTables(
      form.address,
      form.port,
      form.account,
      form.password,
      form.database
  ).then(res => {
    tableList.value = res
  })
}

function getTableStructure() {
  GetTableStructure(
      form.address,
      form.port,
      form.account,
      form.password,
      form.database,
      form.table
  ).then(res => {
    columnList.value = res
  })
}

// 测试
function test() {
  testLoading.value = true
  TestDB(
      form.address,
      form.port,
      form.account,
      form.password,
  ).then(res => {
    if (res) {
      ElMessage.success('测试成功')
      getAllDatabases()
    } else {
      ElMessage.error('测试失败')
    }
  }).finally(() => {
    testLoading.value = false
  })
}
</script>


<style lang="css" scoped>
.main {
  padding: 10px;
  height: 100%;

  .main-content {
    display: flex;
    height: 100%;
  }
}
</style>