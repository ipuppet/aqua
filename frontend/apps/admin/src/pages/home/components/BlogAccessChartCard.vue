<template>
  <card-template>
    <el-row class="nav-box">
      <el-col :span="12">
        <el-select v-model="timeSpan" placeholder="Select" size="small" @change="update()">
          <el-option v-for="item in timeSpanSelection" :key="item" :label="`近 ${item} 天`" :value="item"> </el-option>
        </el-select>
      </el-col>
      <el-col :span="12">
        <div class="refresh" @click="update()">
          <Refresh />
        </div>
      </el-col>
    </el-row>
    <div class="charts" :id="elementId"></div>
  </card-template>
</template>

<script setup lang="ts">
import Refresh from "virtual:icons/mdi/refresh"
import { client, buildApiUrl } from "#/shared/lib/api"
import { sessionData } from "#/shared/lib/storage"
import * as echarts from "echarts"

const $echarts = echarts

const timeSpan = ref(7)
const timeSpanSelection = ref([7, 15, 30])
const elementId = ref("blogAccessChart")

let chart: any = null
let updateButtonLoading = false
let chartLoadingfinished = false
const blogAccess: { x: string[]; data: number[] } = {
  x: [],
  data: []
}
let blogArticleClick: any[] = []

function getDayArr(day: number, span = 1): string[] {
  const myDate = new Date()
  myDate.setDate(myDate.getDate() - (day - 1))
  const dateArray: string[] = []
  let dateTemp: string
  for (let i = 1; i < day; i++) {
    dateTemp = `${myDate.getMonth() + 1}月${myDate.getDate()}日`
    dateArray.push(dateTemp)
    myDate.setDate(myDate.getDate() + span)
  }
  dateArray.push(`${myDate.getMonth() + 1}月${myDate.getDate()}日`)
  return dateArray
}
function loadData(refresh = false, callback?: () => void) {
  const process = (data: any[]) => {
    // 清空数据，避免累积
    blogAccess.x = getDayArr(timeSpan.value)
    blogAccess.data = []
    const blogAccessMap: Record<string, any[]> = {}
    const blogArticleClickMap: Record<string, { name: string; value: number }> = {}

    // 确保 data 是数组
    if (Array.isArray(data) && data.length > 0) {
      data.forEach(item => {
        // 总访问量
        const date = new Date(item.timestamp)
        const accessKey = `${date.getMonth() + 1}月${date.getDate()}日`
        blogAccessMap[accessKey] ??= []
        blogAccessMap[accessKey].push(item)

        const aid = item.aid
        if (aid) {
          blogArticleClickMap[aid] ??= {
            name: item.title,
            value: 0
          }
          blogArticleClickMap[aid].value++
        }
      })
    }

    // 生成数据 - 确保每个日期都有对应的数据值
    blogAccess.x.forEach(key => {
      const len = blogAccessMap[key]?.length ?? 0
      blogAccess.data.push(len)
    })
    blogArticleClick = Object.values(blogArticleClickMap)

    // 确保数据数组长度匹配
    if (blogAccess.data.length !== blogAccess.x.length) {
      blogAccess.data = blogAccess.x.map(() => 0)
    }

    if (callback) callback()
    if (chart) chart.hideLoading()
  }
  const cachedData = sessionData.get("home.blogAccess")
  const cachedTimeSpan = sessionData.get("home.blogAccess.timeSpan")

  if (chart) chart.showLoading() // 加载动画
  if (refresh || cachedData === null || cachedTimeSpan !== timeSpan.value) {
    // 清空数据
    blogAccess.x = getDayArr(timeSpan.value)
    blogAccess.data = []
    blogArticleClick = []
    // 获取数据
    client
      .GET("/api/access/{day}", {
        baseUrl: buildApiUrl("blog"),
        params: { path: { day: timeSpan.value } }
      })
      .then(({ data: res, error }) => {
        if (error || !res) throw new Error(error?.message ?? "Error")
        sessionData.set("home.blogAccess", res.records)
        sessionData.set("home.blogAccess.timeSpan", timeSpan.value)
        process(res.records ?? [])
      })
      .catch(() => {
        // 处理错误，包括404
        // 设置空数据
        blogAccess.x = getDayArr(timeSpan.value)
        blogAccess.data = Array(blogAccess.x.length).fill(0)
        blogArticleClick = [{ name: "无数据", value: 0 }]
        if (callback) callback()
        if (chart) chart.hideLoading()
      })
  } else {
    // 使用缓存数据，但需要根据当前的 timeSpan 重新处理
    process(cachedData ?? [])
  }
}
function chartResize(): void {
  const elem = document.getElementById(elementId.value)
  if (elem && elem.parentElement) {
    // card 有 40px 的 padding，刷新按钮以及下拉菜单高度 25px，margin-bottom: 10px;
    elem.style.width = elem.parentElement.clientWidth - 40 + "px"
    elem.style.height = elem.parentElement.clientHeight - 75 + "px"
    if (chart) chart.resize()
  }
}

// 设置图表选项的通用函数
function setupChartOptions(needsClear = false): void {
  if (!chart) return

  // 检查是否有数据
  const hasNoData =
    blogAccess.data.every(value => value === 0) &&
    (blogArticleClick.length === 0 || (blogArticleClick.length === 1 && blogArticleClick[0].name === "无数据"))

  // 如果需要清除之前的配置
  if (needsClear) chart.clear()

  if (hasNoData) {
    // 无数据时只显示提示文字
    chart.setOption({
      title: { text: "访问情况" },
      graphic: {
        elements: [
          {
            type: "text",
            left: "center",
            top: "middle",
            style: {
              text: "暂无数据",
              fontSize: 24,
              fill: "#999"
            }
          }
        ]
      }
    })
  } else {
    // 有数据时显示完整图表
    // 确保 x 轴数据不为空
    if (blogAccess.x.length === 0) {
      blogAccess.x = getDayArr(timeSpan.value)
      blogAccess.data = Array(blogAccess.x.length).fill(0)
    }

    // 确保数据长度匹配
    let barData = blogAccess.data.length === blogAccess.x.length ? blogAccess.data : blogAccess.x.map(() => 0)

    // 确保 barData 不为空且长度正确
    if (barData.length === 0 || barData.length !== blogAccess.x.length) {
      barData = Array(blogAccess.x.length).fill(0)
    }

    chart.setOption(
      {
        title: { text: "访问情况" },
        tooltip: {
          trigger: "axis",
          axisPointer: {
            type: "shadow"
          }
        },
        xAxis: {
          type: "category",
          data: blogAccess.x,
          axisLabel: {
            rotate: 45,
            formatter: (value: any) => String(value)
          },
          boundaryGap: true,
          show: true
        },
        yAxis: {
          type: "value",
          name: "访问量",
          minInterval: 1,
          show: true
        },
        grid: {
          top: "55%",
          bottom: "15%",
          left: "10%",
          right: "10%",
          containLabel: true
        },
        series: [
          {
            name: "访问量",
            type: "bar",
            emphasis: { focus: "series" },
            data: barData,
            itemStyle: {
              color: "#409EFF"
            },
            showBackground: true,
            backgroundStyle: {
              color: "rgba(180, 180, 180, 0.2)"
            }
          },
          {
            name: "点击量",
            type: "pie",
            radius: "35%",
            center: ["50%", "32%"],
            emphasis: { focus: "data" },
            data: blogArticleClick,
            tooltip: {
              trigger: "item",
              formatter: "{b}: {c} ({d}%)"
            }
          }
        ]
      },
      true
    ) // 使用 notMerge: true 确保完全替换配置
  }

  // 确保图表调整到正确大小
  chartResize()
}

function initChart(): void {
  // 先设置合适的容器大小
  chartResize()

  // 初始化图表实例
  const chartElement = document.getElementById(elementId.value)
  if (!chartElement) {
    // 如果元素不存在，延迟重试
    setTimeout(initChart, 100)
    return
  }

  chart = $echarts.init(chartElement)

  // 设置图表选项
  setupChartOptions()

  chart.on("finished", () => {
    if (chartLoadingfinished) return
    chartLoadingfinished = true
    chartResize() // 再次确保调整大小
  })
}

function update(): void {
  if (updateButtonLoading) return
  updateButtonLoading = true
  loadData(true, () => {
    updateButtonLoading = false
    // 使用共用函数设置图表选项，需要先清除现有配置
    setupChartOptions(true)
  })
}

onMounted(() => {
  // 加载数据
  loadData(false, () => {
    initChart()
    // 添加窗口调整大小的事件监听
    window.addEventListener("resize", chartResize)
  })
})

onBeforeUnmount(() => {
  // 组件卸载前移除事件监听器
  window.removeEventListener("resize", chartResize)
})
</script>

<style scoped>
.nav-box {
  height: 25px;
  overflow: hidden;
  margin-bottom: 10px;
}

.refresh {
  width: 25px;
  height: 25px;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: Center;
  position: relative;
  top: 0;
  margin-left: auto;
}

.charts {
  width: 400px;
  height: 400px;
}
</style>
