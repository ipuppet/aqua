<template>
  <div class="footer animate__animated animate__fadeIn">
    <el-row align="middle" type="flex">
      <el-col :span="12">
        <span>2018-{{ year }} @ {{ siteConfig.title }}</span>
      </el-col>
      <el-col :span="12">
        <div class="footer-links-desktop">
          <div class="links active">
            <a class="bottom-link" :href="link.href" :key="i" target="_blank" v-for="(link, i) in siteConfig.links">
              {{ link.title }}
            </a>
          </div>
        </div>
        <div class="footer-links-mobile">
          <div class="menu-button-box">
            <menu-button size="sm" :links="siteConfig.links" />
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue"
import MenuButton from "#/shared/ui/MenuButton.vue"
import { siteConfig } from "#/shared/config/site"

const year = computed(() => new Date().getFullYear())
</script>

<style scoped>
.footer {
  position: absolute;
  right: 0;
  left: 0;
  bottom: 0;
  padding: 19px 30px 20px;
  background-color: rgba(255, 255, 255, 0);
  transition: 0.3s;
}

.menu-button-box {
  float: right;
}

.footer-links-desktop {
  display: none;
}

.footer-links-mobile {
  display: block;
}

@media (min-width: 768px) {
  .footer-links-desktop {
    display: block;
  }

  .footer-links-mobile {
    display: none;
  }
}

.links {
  display: none;
  float: right;
}

.links.active {
  display: block;
}

.bottom-link {
  color: #98a6ad;
  text-align: center;
  text-decoration: underline;
  margin-right: 10px;
}

.icon-heart {
  margin: 0 5px 0 5px;
}
</style>
