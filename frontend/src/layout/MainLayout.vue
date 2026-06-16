<template>
  <el-container style="min-height: 100vh">
    <el-aside width="220px" style="background: #304156">
      <div style="color:#fff;font-size:18px;padding:20px 16px;font-weight:bold;text-align:center">
        远洋船员证书排班
      </div>
      <el-menu
        :default-active="route.path"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <el-menu-item index="/seafarers">
          <el-icon><User /></el-icon>
          <span>船员管理</span>
        </el-menu-item>
        <el-menu-item index="/certificates">
          <el-icon><Document /></el-icon>
          <span>证书管理</span>
        </el-menu-item>
        <el-menu-item index="/ships">
          <el-icon><Ship /></el-icon>
          <span>船舶管理</span>
        </el-menu-item>
        <el-menu-item index="/assignments">
          <el-icon><Connection /></el-icon>
          <span>在船分配</span>
        </el-menu-item>
        <el-menu-item index="/transfers">
          <el-icon><Switch /></el-icon>
          <span>船员调动</span>
        </el-menu-item>
        <el-menu-item index="/alerts">
          <el-icon><Bell /></el-icon>
          <span>证书预警</span>
        </el-menu-item>
        <el-menu-item index="/records">
          <el-icon><Tickets /></el-icon>
          <span>合同与记录</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-main style="background:#f0f2f5;padding:20px">
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { User, Document, Connection, Switch, Bell, Tickets } from '@element-plus/icons-vue'

const Ship = Document
const route = useRoute()
</script>
