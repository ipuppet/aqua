<template>
  <div class="container">
    <el-dropdown trigger="click" @command="handleCommand">
      <button class="pointer user">
        <img :src="info?.avatar ?? ''" alt="avatar" class="avatar" />
        <span class="nickname">Hi,&nbsp;{{ info?.nickname ?? "Guest" }}</span>
      </button>
      <template #dropdown>
        <el-dropdown-menu>
          <router-link to="/me/info/">
            <el-dropdown-item>
              <User class="icon" />
              <span>个人信息</span>
            </el-dropdown-item>
          </router-link>
          <router-link to="/me/security/">
            <el-dropdown-item>
              <Security class="icon" />
              <span>安全设置</span>
            </el-dropdown-item>
          </router-link>
          <el-dropdown-item command="logout">
            <Logout class="icon" />
            <span>Logout</span>
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup lang="ts">
import User from "virtual:icons/mdi/account"
import Security from "virtual:icons/mdi/shield-account"
import Logout from "virtual:icons/mdi/logout"
import { useUserStore } from "#/shared/stores/user"
import { storeToRefs } from "pinia"
import { useSession } from "#/features/session/model/useSession"

const userStore = useUserStore()
const { info } = storeToRefs(userStore)
const { logout } = useSession()

function handleCommand(command: string) {
  switch (command) {
    case "logout":
      logout()
      break
  }
}
</script>

<style scoped>
.container {
  float: right;
  margin: 10px 10px 0 0;
  padding: 0 10px 0 10px;
  height: 50px;
  border-radius: 50px;
  transition: background-color 0.3s;
}

.container:hover {
  background-color: #f0f3ff;
}

.icon {
  margin-right: 10px;
}

.user {
  background: rgba(255, 255, 255, 0);
  height: 100%;
  border: none;
  font-size: 12px;
  color: #000000;
}

.avatar {
  float: right;
  width: 38px;
  height: 38px;
  border-radius: 45px;
  margin-left: 10px;
  margin-top: 6px;
}

.nickname {
  float: left;
  font-size: 12px;
  line-height: 50px;
}

.user:focus {
  outline: none;
}
</style>
