<template>
  <div class="main">
    <div class="main-content">
      <el-card style="flex: 1;margin-right: 10px" shadow="">
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
            <el-input type="text" v-model="form.password" placeholder="请输入密码"/>
          </el-form-item>
          <el-form-item label="数据库" prop="database">
            <el-select v-model="form.database" @focus="getAllDatabases" @change="getAllTables"
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
          <el-form-item label="信息" prop="info">
            <el-input type="textarea" v-model="form.info" placeholder="请输入信息"/>
          </el-form-item>
          <el-form-item>
            <el-space>
              <el-button @click="test">测试</el-button>
              <el-button type="primary" @click="getAIAnswer" :loading>执行</el-button>
            </el-space>
          </el-form-item>
        </el-form>
      </el-card>
      <el-card style="flex: 2" shadow="">
        <div v-for="(item,index) in columnList" :key="index">
          <div style="margin-bottom: 7px">
            <el-text type="primary" tag="ins">{{ index }}</el-text>
          </div>
          <el-scrollbar>
            <el-space>
              <el-tag type="info" v-for="(c,i) in item" :key="i">{{ c.Field }}</el-tag>
            </el-space>
          </el-scrollbar>
        </div>
        <el-divider/>
        <el-scrollbar>
          <MdPreview editorId="MdPreview" :modelValue="aiInfo"/>
        </el-scrollbar>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import {onMounted, reactive, ref} from 'vue'
import {GetAIAnswer, GetAllDatabases, GetAllTables, GetTableStructure, TestDB} from '@api/internal/App'
import {ElMessage} from "element-plus";
import {MdPreview} from 'md-editor-v3';
import 'md-editor-v3/lib/preview.css';

const form = reactive({
  address: '',
  port: '',
  account: '',
  password: '',
  database: '',
  table: [],
  info: ''
})

const database = 'databaseInfo'

const loading = ref(false)
const aiInfo = ref<string>('')
const options = ref<any>([])
const columnList = ref<any>({})
const tableList = ref<any>([])

onMounted(() => {
  if (localStorage.getItem(database)) {
    const data = JSON.parse(localStorage.getItem(database) as string)
    Object.assign(form, data)
    getAllDatabases()
  }
})


function getAllDatabases() {
  GetAllDatabases(
      form.address,
      form.port,
      form.account,
      form.password
  ).then(res => {
    options.value = res
    getAllTables()
    localStorage.setItem(database, JSON.stringify(form))
  })
}

// 获取回答 问题+表结构
function getAIAnswer() {
  if (form.info) {
    loading.value = true
    GetAIAnswer(
        JSON.stringify(columnList.value),
        form.info
    ).then(res => {
      aiInfo.value = res
      ElMessage.success("获取成功")
      loading.value = false
    })
  } else {
    ElMessage.error('请输入信息')
  }
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
  columnList.value = {}
  if (form.table.length > 0) {
    form.table.forEach(item => {
      GetTableStructure(
          form.address,
          form.port,
          form.account,
          form.password,
          form.database,
          item
      ).then(res => {
        columnList.value[item] = res
      })
    })
  }
}

function handelClick() {
}

// 测试
function test() {
  TestDB(
      form.address,
      form.port,
      form.account,
      form.password,
      form.database
  ).then(res => {
    if (res) {
      ElMessage.success('测试成功')
    } else {
      ElMessage.error('测试失败')
    }
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