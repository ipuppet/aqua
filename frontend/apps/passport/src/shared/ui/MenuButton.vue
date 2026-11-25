<template>
  <div>
    <!--size
        sm 小
        md 中
        lg 大-->
    <div :class="'container size-' + size" @click="click">
      <div :class="['item', 'item-first', { active: active }]"></div>
      <div :class="['item', 'item-second', { active: active }]"></div>
      <div :class="['item', 'item-third', { active: active }]"></div>
    </div>
    <div :class="['links', { active: active }]">
      <a :href="link.href" :key="i" target="_blank" v-for="(link, i) in links">{{ link.title }}</a>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"

interface Link {
  title: string
  href: string
}

interface Props {
  size?: string
  links?: Link[]
}

withDefaults(defineProps<Props>(), {
  size: "md",
  links: () => []
})

const active = ref(false)

function click() {
  active.value = !active.value
}
</script>

<style scoped>
.container {
  cursor: pointer;
  background-color: rgba(255, 255, 255, 0.15);
  padding: 5px;
  border-radius: 5%;
}

.size-sm {
  width: 40px;
  height: 27px;
}

.size-md {
  width: 50px;
  height: 33px;
}

.size-lg {
  width: 60px;
  height: 40px;
}

.item {
  width: 80%;
  height: 15%;
  margin-top: 8%;
  margin-left: 10%;
  border-radius: 5px;
  background-color: #d0d0d0;
  transition: 0.3s;
}

.item-first.active {
  margin-top: 24%;
  transform: rotate(40deg);
}

.item-second.active {
  display: none;
}

.item-third.active {
  margin-top: -8%;
  transform: rotate(-40deg);
}

.links {
  height: 0;
  overflow: hidden;
}

.links a {
  color: #2e2e2e;
  display: block;
  text-align: center;
  margin-bottom: 10px;
}

.links.active {
  margin-top: 20px;
  height: auto;
}
</style>
