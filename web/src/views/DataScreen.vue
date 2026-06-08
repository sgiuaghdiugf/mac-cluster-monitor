<template>
  <div class="datascreen" ref="screenRef">
    <!-- 顶部标题 -->
    <header class="screen-header">
      <div class="header-decoration left-deco">
        <svg viewBox="0 0 400 60" preserveAspectRatio="none">
          <defs>
            <linearGradient id="headerLineL" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" stop-color="rgba(0,0,0,0)" />
              <stop offset="50%" stop-color="rgba(0,212,255,0.6)" />
              <stop offset="100%" stop-color="rgba(0,212,255,0.2)" />
            </linearGradient>
          </defs>
          <line x1="0" y1="30" x2="400" y2="30" stroke="url(#headerLineL)" stroke-width="1" />
          <circle cx="200" cy="30" r="3" fill="#00d4ff" opacity="0.8">
            <animate attributeName="opacity" values="0.4;1;0.4" dur="2s" repeatCount="indefinite" />
          </circle>
        </svg>
      </div>
      <div class="header-center">
        <h1 class="screen-title">Mac Cluster Monitor</h1>
        <p class="screen-subtitle">MAC CLUSTER INTELLIGENT MONITORING PLATFORM</p>
      </div>
      <div class="header-decoration right-deco">
        <svg viewBox="0 0 400 60" preserveAspectRatio="none">
          <defs>
            <linearGradient id="headerLineR" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" stop-color="rgba(0,212,255,0.2)" />
              <stop offset="50%" stop-color="rgba(0,212,255,0.6)" />
              <stop offset="100%" stop-color="rgba(0,0,0,0)" />
            </linearGradient>
          </defs>
          <line x1="0" y1="30" x2="400" y2="30" stroke="url(#headerLineR)" stroke-width="1" />
          <circle cx="200" cy="30" r="3" fill="#00d4ff" opacity="0.8">
            <animate attributeName="opacity" values="0.4;1;0.4" dur="2s" repeatCount="indefinite" />
          </circle>
        </svg>
      </div>
      <div class="header-time">{{ currentTime }}</div>
    </header>

    <!-- 主内容区 -->
    <div class="screen-body">
      <!-- 左侧面板 -->
      <div class="panel left-panel">
        <!-- 规模指标分析 -->
        <div class="panel-card">
          <div class="card-header">
            <span class="card-title">规模指标分析</span>
            <span class="card-subtitle">INDICATOR ANALYSIS</span>
          </div>
          <div class="card-body">
            <div class="ranking-list">
              <div class="ranking-header">
                <span class="col-rank">排名</span>
                <span class="col-name">用户</span>
                <span class="col-count">设备数</span>
                <span class="col-profit">累计收益</span>
              </div>
              <div class="ranking-body" ref="rankingScrollRef">
                <div class="ranking-scroll-inner" :style="{ transform: `translateY(-${rankingScrollOffset}px)`, transition: rankingScrolling ? 'transform 0.6s ease' : 'none' }">
                  <div
                    v-for="(item, idx) in rankingDisplayList"
                    :key="idx"
                    class="ranking-row"
                    :class="{ 'top-1': item.rank === 1, 'top-2': item.rank === 2, 'top-3': item.rank === 3 }"
                  >
                    <span class="col-rank">
                      <i v-if="item.rank <= 3" class="rank-badge" :class="`rank-${item.rank}`">{{ item.rank }}</i>
                      <span v-else class="rank-normal">{{ item.rank }}</span>
                    </span>
                    <span class="col-name">{{ item.name }}</span>
                    <span class="col-count">{{ item.deviceCount }}台</span>
                    <span class="col-profit">¥{{ item.profit.toLocaleString() }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 设备总数量 -->
        <div class="panel-card">
          <div class="card-header">
            <span class="card-title">设备总数量</span>
            <span class="card-subtitle">DEVICE STATISTICS</span>
          </div>
          <div class="card-body">
            <div class="device-stats-grid-2x2">
              <div class="stat-item-2 total">
                <div class="stat-value">{{ deviceStats.total }}</div>
                <div class="stat-label">设备总数</div>
              </div>
              <div class="stat-item-2 online">
                <div class="stat-value">{{ deviceStats.online }}</div>
                <div class="stat-label">在线设备</div>
              </div>
              <div class="stat-item-2 offline">
                <div class="stat-value">{{ deviceStats.offline }}</div>
                <div class="stat-label">离线设备</div>
              </div>
              <div class="stat-item-2 rate">
                <div class="stat-value">{{ deviceStats.onlineRate }}%</div>
                <div class="stat-label">在线率</div>
              </div>
            </div>
            <div class="device-bar">
              <div class="bar-online" :style="{ width: deviceStats.onlineRate + '%' }"></div>
            </div>
            <div class="device-total-profit">
              <span class="profit-label">累计总收益</span>
              <span class="profit-mini">¥{{ formatProfit(currentProfit) }}</span>
            </div>
          </div>
        </div>

        <!-- 设备上线情况 -->
        <div class="panel-card">
          <div class="card-header">
            <span class="card-title">设备上线情况</span>
            <span class="card-subtitle">ONLINE DEVICES</span>
          </div>
          <div class="card-body">
            <div class="device-table">
              <div class="table-header">
                <span>设备名称</span>
                <span>MAC地址</span>
                <span>在线时长</span>
                <span>累计收益</span>
              </div>
              <div class="table-scroll" ref="onlineScrollRef">
                <div class="table-scroll-inner" :style="{ transform: `translateY(-${onlineScrollOffset}px)`, transition: onlineScrolling ? 'transform 0.5s ease' : 'none' }">
                  <div v-for="(item, idx) in onlineDisplayList" :key="idx" class="table-row online-row">
                    <span class="cell-name">
                      <i class="status-dot online-dot"></i>
                      {{ item.name }}
                    </span>
                    <span class="cell-mac">{{ item.mac }}</span>
                    <span class="cell-duration">{{ item.duration }}</span>
                    <span class="cell-profit">¥{{ item.profit.toLocaleString() }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 中间地图区域 -->
      <div class="center-panel">
        <div class="map-container" ref="mapRef"></div>
        <div class="map-stats">
          <div class="map-stat-item" v-for="(item, idx) in mapStatItems" :key="idx">
            <span class="map-stat-value">{{ item.value }}</span>
            <span class="map-stat-label">{{ item.label }}</span>
          </div>
        </div>
      </div>

      <!-- 右侧面板 -->
      <div class="panel right-panel">
        <!-- 设备离线情况 -->
        <div class="panel-card">
          <div class="card-header">
            <span class="card-title">设备离线情况</span>
            <span class="card-subtitle">OFFLINE DEVICES</span>
          </div>
          <div class="card-body">
            <div class="device-table">
              <div class="table-header">
                <span>设备名称</span>
                <span>MAC地址</span>
                <span>异常次数</span>
                <span>状态</span>
              </div>
              <div class="table-scroll" ref="offlineScrollRef">
                <div class="table-scroll-inner" :style="{ transform: `translateY(-${offlineScrollOffset}px)`, transition: offlineScrolling ? 'transform 0.5s ease' : 'none' }">
                  <div v-for="(item, idx) in offlineDisplayList" :key="idx" class="table-row offline-row">
                    <span class="cell-name">
                      <i class="status-dot offline-dot"></i>
                      {{ item.name }}
                    </span>
                    <span class="cell-mac">{{ item.mac }}</span>
                    <span class="cell-error" :class="{ 'error-high': item.errorCount >= 3 }">{{ item.errorCount }}次</span>
                    <span class="cell-status offline-status">{{ item.status }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 收益情况 -->
        <div class="panel-card">
          <div class="card-header">
            <span class="card-title">收益情况</span>
            <span class="card-subtitle">PROFIT ANALYSIS</span>
          </div>
          <div class="card-body">
            <div class="profit-chart" ref="profitChartRef"></div>
          </div>
        </div>

        <!-- 总收益 -->
        <div class="panel-card total-profit-card">
          <div class="card-header">
            <span class="card-title">总收益</span>
            <span class="card-subtitle">TOTAL PROFIT</span>
          </div>
          <div class="card-body">
            <!-- 动态光环 -->
            <div class="profit-ring-wrapper">
              <svg class="profit-ring-svg" viewBox="0 0 120 120">
                <defs>
                  <linearGradient id="ringGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" stop-color="#fbbf24" />
                    <stop offset="50%" stop-color="#34d399" />
                    <stop offset="100%" stop-color="#00d4ff" />
                  </linearGradient>
                </defs>
                <circle cx="60" cy="60" r="54" fill="none" stroke="rgba(0,212,255,0.08)" stroke-width="3"/>
                <circle class="ring-anim ring-anim-1" cx="60" cy="60" r="48" fill="none" stroke="url(#ringGrad)" stroke-width="1.5"
                  stroke-dasharray="180 100" opacity="0.6"/>
                <circle class="ring-anim ring-anim-2" cx="60" cy="60" r="42" fill="none" stroke="#fbbf24" stroke-width="1"
                  stroke-dasharray="140 160" opacity="0.4"/>
              </svg>
              <div class="profit-center">
                <div class="total-profit-value">
                  <span class="profit-symbol">¥</span>
                  <span class="profit-number">{{ formatProfit(currentProfit) }}</span>
                </div>
              </div>
            </div>
            <!-- 子指标 -->
            <div class="profit-sub-metrics">
              <div class="sub-metric">
                <span class="sub-label">今日收益</span>
                <span class="sub-value today">¥{{ formatProfit(todayProfit) }}</span>
              </div>
              <div class="sub-metric-divider"></div>
              <div class="sub-metric">
                <span class="sub-label">昨日收益</span>
                <span class="sub-value yesterday">¥{{ formatProfit(yesterdayProfit) }}</span>
              </div>
            </div>
            <div class="profit-growth">
              <i class="growth-arrow"></i>
              <span>实时更新中</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 返回按钮 -->
    <div class="back-btn" @click="goBack">
      <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor">
        <path d="M20 11H7.83l5.59-5.59L12 4l-8 8 8 8 1.41-1.41L7.83 13H20v-2z"/>
      </svg>
      <span>返回</span>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import {
  scaleRanking,
  deviceStats,
  onlineDevices,
  offlineDevices,
  profitByQuarter,
  totalProfit,
  mapDeviceData
} from '../data/screenData'

const router = useRouter()
const screenRef = ref(null)
const mapRef = ref(null)
const profitChartRef = ref(null)

// 时间
const currentTime = ref('')
let timeTimer = null

// 总收益动态增长
const currentProfit = ref(totalProfit.initialValue)
let profitTimer = null

// 排行滚动
const rankingScrollOffset = ref(0)
const rankingScrolling = ref(false)
let rankingTimer = null

// 在线设备滚动
const onlineScrollOffset = ref(0)
const onlineScrolling = ref(false)
let onlineTimer = null

// 离线设备滚动
const offlineScrollOffset = ref(0)
const offlineScrolling = ref(false)
let offlineTimer = null

// 地图和饼图实例
let mapChart = null
let profitChart = null

// 排行展示列表（翻倍用于无缝滚动）
const rankingDisplayList = computed(() => [...scaleRanking, ...scaleRanking])
const onlineDisplayList = computed(() => [...onlineDevices, ...onlineDevices])
const offlineDisplayList = computed(() => [...offlineDevices, ...offlineDevices])

// 地图底部统计
const mapStatItems = computed(() => [
  { value: deviceStats.total, label: '设备总数' },
  { value: deviceStats.online, label: '在线设备' },
  { value: deviceStats.offline, label: '离线设备' },
  { value: deviceStats.onlineRate + '%', label: '在线率' }
])

// 今日收益 / 昨日收益（静态模拟数据）
const todayProfit = ref(86.50)
const yesterdayProfit = ref(81.20)

// 格式化收益数字
const formatProfit = (val) => {
  return val.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// 更新时间
const updateTime = () => {
  const now = new Date()
  const y = now.getFullYear()
  const m = String(now.getMonth() + 1).padStart(2, '0')
  const d = String(now.getDate()).padStart(2, '0')
  const h = String(now.getHours()).padStart(2, '0')
  const min = String(now.getMinutes()).padStart(2, '0')
  const s = String(now.getSeconds()).padStart(2, '0')
  currentTime.value = `${y}-${m}-${d} ${h}:${min}:${s}`
}

// 初始化地图
const initMap = async () => {
  if (!mapRef.value) return
  try {
    const response = await fetch('/china.json')
    const chinaJson = await response.json()
    echarts.registerMap('china', chinaJson)

    mapChart = echarts.init(mapRef.value)

    const mapData = mapDeviceData.map(item => ({
      name: item.name,
      value: item.value
    }))

    mapChart.setOption({
      tooltip: {
        trigger: 'item',
        backgroundColor: 'rgba(0, 20, 40, 0.9)',
        borderColor: '#00d4ff',
        borderWidth: 1,
        textStyle: { color: '#e0f4ff', fontSize: 14 },
        formatter: (params) => {
          if (params.value) {
            return `<b>${params.name}</b><br/>设备数量: ${params.value}台`
          }
          return `<b>${params.name}</b><br/>暂无设备`
        }
      },
      visualMap: {
        min: 0,
        max: 25,
        left: 20,
        bottom: 20,
        text: ['高', '低'],
        textStyle: { color: '#8ab4d6', fontSize: 12 },
        inRange: {
          color: ['#0a2e4a', '#0d4a6e', '#0e6ea8', '#1198d4', '#00d4ff']
        },
        calculable: true,
        itemWidth: 12,
        itemHeight: 80
      },
      geo: {
        map: 'china',
        roam: true,
        zoom: 1.2,
        center: [104, 36],
        label: { show: false },
        itemStyle: {
          areaColor: '#0a2e4a',
          borderColor: '#0e6ea8',
          borderWidth: 1,
          shadowColor: 'rgba(0, 212, 255, 0.3)',
          shadowBlur: 10
        },
        emphasis: {
          itemStyle: {
            areaColor: '#1198d4',
            borderColor: '#00d4ff',
            borderWidth: 2,
            shadowColor: 'rgba(0, 212, 255, 0.6)',
            shadowBlur: 20
          },
          label: {
            show: true,
            color: '#fff',
            fontSize: 14
          }
        }
      },
      series: [
        {
          name: '设备分布',
          type: 'map',
          map: 'china',
          geoIndex: 0,
          data: mapData
        },
        {
          name: '设备散点',
          type: 'effectScatter',
          coordinateSystem: 'geo',
          data: [
            { name: '北京', value: [116.46, 39.92, 18] },
            { name: '上海', value: [121.48, 31.22, 15] },
            { name: '广东', value: [113.23, 23.16, 22] },
            { name: '浙江', value: [120.19, 30.26, 12] },
            { name: '江苏', value: [118.78, 32.04, 10] },
            { name: '四川', value: [104.06, 30.67, 8] },
            { name: '湖北', value: [114.31, 30.52, 6] },
            { name: '福建', value: [119.30, 26.08, 5] }
          ],
          symbolSize: (val) => Math.max(val[2] * 0.8, 6),
          showEffectOn: 'render',
          rippleEffect: {
            brushType: 'stroke',
            scale: 4,
            period: 4
          },
          label: {
            show: true,
            formatter: '{b}',
            position: 'right',
            color: '#00d4ff',
            fontSize: 11
          },
          itemStyle: {
            color: '#00d4ff',
            shadowBlur: 10,
            shadowColor: 'rgba(0, 212, 255, 0.5)'
          }
        }
      ]
    })
  } catch (e) {
    console.error('加载地图数据失败', e)
  }
}

// 初始化收益饼图
const initProfitChart = () => {
  if (!profitChartRef.value) return
  profitChart = echarts.init(profitChartRef.value)

  const colors = ['#fbbf24', '#34d399', '#ef4444', '#f97316']

  profitChart.setOption({
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(0, 20, 40, 0.9)',
      borderColor: '#00d4ff',
      borderWidth: 1,
      textStyle: { color: '#e0f4ff' },
      formatter: '{b}: ¥{c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: '5%',
      top: 'middle',
      textStyle: { color: '#8ab4d6', fontSize: 12 },
      itemWidth: 10,
      itemHeight: 10,
      itemGap: 16
    },
    series: [
      {
        name: '收益分布',
        type: 'pie',
        radius: ['45%', '70%'],
        center: ['35%', '50%'],
        avoidLabelOverlap: false,
        label: { show: false },
        labelLine: { show: false },
        itemStyle: {
          borderRadius: 6,
          borderColor: '#0a1628',
          borderWidth: 2
        },
        data: profitByQuarter.map((item, idx) => ({
          value: item.value,
          name: item.name,
          itemStyle: { color: colors[idx] }
        }))
      }
    ]
  })
}

// 无缝滚动工具函数
const startScroll = (offsetRef, scrollingRef, itemCount, rowHeight, interval = 3000) => {
  let current = 0
  const totalHeight = itemCount * rowHeight

  return setInterval(() => {
    scrollingRef.value = true
    current++
    offsetRef.value = current * rowHeight

    setTimeout(() => {
      if (current >= itemCount) {
        scrollingRef.value = false
        current = 0
        offsetRef.value = 0
      }
    }, 600)
  }, interval)
}

// 返回
const goBack = () => {
  router.push('/dashboard')
}

// 窗口resize
const handleResize = () => {
  mapChart?.resize()
  profitChart?.resize()
}

onMounted(async () => {
  updateTime()
  timeTimer = setInterval(updateTime, 1000)

  // 总收益定时增长
  profitTimer = setInterval(() => {
    currentProfit.value += totalProfit.increment
  }, totalProfit.interval)

  await nextTick()

  // 初始化图表
  initMap()
  initProfitChart()

  // 启动滚动
  rankingTimer = startScroll(rankingScrollOffset, rankingScrolling, scaleRanking.length, 40, 3000)
  onlineTimer = startScroll(onlineScrollOffset, onlineScrolling, onlineDevices.length, 36, 3500)
  offlineTimer = startScroll(offlineScrollOffset, offlineScrolling, offlineDevices.length, 36, 4000)

  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  clearInterval(timeTimer)
  clearInterval(profitTimer)
  clearInterval(rankingTimer)
  clearInterval(onlineTimer)
  clearInterval(offlineTimer)
  mapChart?.dispose()
  profitChart?.dispose()
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
/* ============================================
   数据大屏 - 全局样式
   参考 sc-datav 四川省电力全景感知平台风格
   深蓝科技风 + 青色点缀
   ============================================ */
.datascreen {
  position: fixed;
  inset: 0;
  width: 100vw;
  height: 100vh;
  background: radial-gradient(ellipse at 50% 0%, #0a1e3d 0%, #060e1f 50%, #020810 100%);
  color: #e0f4ff;
  font-family: 'Rajdhani', 'Orbitron', 'Courier New', monospace;
  overflow: hidden;
  z-index: 9999;
}

/* 顶部标题栏 */
.screen-header {
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  background: linear-gradient(180deg, rgba(0, 40, 80, 0.6) 0%, transparent 100%);
  border-bottom: 1px solid rgba(0, 212, 255, 0.15);
}

.header-decoration {
  flex: 1;
  height: 60px;
  padding: 0 20px;
}

.header-decoration svg {
  width: 100%;
  height: 100%;
}

.header-center {
  text-align: center;
  padding: 0 40px;
}

.screen-title {
  font-size: 32px;
  font-weight: 700;
  letter-spacing: 8px;
  background: linear-gradient(180deg, #00d4ff 0%, #0088cc 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  text-shadow: 0 0 30px rgba(0, 212, 255, 0.3);
  margin: 0;
}

.screen-subtitle {
  font-size: 11px;
  letter-spacing: 6px;
  color: rgba(0, 212, 255, 0.5);
  margin: 2px 0 0;
}

.header-time {
  position: absolute;
  right: 30px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 16px;
  color: #8ab4d6;
  letter-spacing: 2px;
}

/* 主内容区 */
.screen-body {
  display: flex;
  height: calc(100vh - 80px);
  padding: 16px;
  gap: 16px;
}

/* 面板通用样式 */
.panel {
  width: 28%;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.center-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  position: relative;
}

/* 卡片样式 */
.panel-card {
  flex: 1;
  background: rgba(0, 30, 60, 0.5);
  border: 1px solid rgba(0, 212, 255, 0.15);
  border-radius: 4px;
  position: relative;
  overflow: hidden;
  backdrop-filter: blur(4px);
  display: flex;
  flex-direction: column;
}

/* 卡片四角装饰 */
.panel-card::before,
.panel-card::after {
  content: '';
  position: absolute;
  width: 12px;
  height: 12px;
  border-color: #00d4ff;
  border-style: solid;
  transition: all 0.3s ease;
  pointer-events: none;
  z-index: 2;
}

.panel-card::before {
  top: -1px;
  left: -1px;
  border-width: 2px 0 0 2px;
}

.panel-card::after {
  bottom: -1px;
  right: -1px;
  border-width: 0 2px 2px 0;
}

.panel-card:hover::before,
.panel-card:hover::after {
  width: 30px;
  height: 30px;
  opacity: 0.8;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 14px;
  border-bottom: 1px solid rgba(0, 212, 255, 0.1);
  background: rgba(0, 40, 80, 0.3);
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: #e0f4ff;
  padding-left: 10px;
  border-left: 3px solid #00d4ff;
}

.card-subtitle {
  font-size: 10px;
  color: rgba(0, 212, 255, 0.4);
  letter-spacing: 1px;
}

.card-body {
  flex: 1;
  padding: 10px 14px;
  overflow: hidden;
  min-height: 0;
}

/* ============================================
   规模指标排行
   ============================================ */
.ranking-list {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.ranking-header {
  display: flex;
  align-items: center;
  padding: 6px 0;
  border-bottom: 1px solid rgba(0, 212, 255, 0.1);
  font-size: 12px;
  color: #5a8aaa;
}

.ranking-body {
  flex: 1;
  overflow: hidden;
}

.ranking-row {
  display: flex;
  align-items: center;
  height: 40px;
  border-bottom: 1px solid rgba(0, 212, 255, 0.05);
  font-size: 13px;
  transition: background 0.3s;
}

.ranking-row:hover {
  background: rgba(0, 212, 255, 0.05);
}

.col-rank { width: 50px; text-align: center; }
.col-name { flex: 1; }
.col-count { width: 70px; text-align: center; }
.col-profit { width: 110px; text-align: right; color: #00d4ff; font-weight: 600; }

.rank-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border-radius: 4px;
  font-size: 12px;
  font-style: normal;
  font-weight: 700;
  color: #fff;
}

.rank-1 { background: linear-gradient(135deg, #ff6b35, #ff3d00); }
.rank-2 { background: linear-gradient(135deg, #00d4ff, #0088cc); }
.rank-3 { background: linear-gradient(135deg, #0e6ea8, #0a4a7a); }

.rank-normal {
  color: #5a8aaa;
  font-size: 12px;
}

/* ============================================
   设备统计 - 2x2布局 + 总收益
   ============================================ */
.device-stats-grid-2x2 {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.stat-item-2 {
  text-align: center;
  padding: 12px 8px;
  background: rgba(0, 30, 60, 0.5);
  border-radius: 4px;
  border: 1px solid rgba(0, 212, 255, 0.1);
  transition: all 0.3s ease;
}

.stat-item-2:hover {
  border-color: rgba(0, 212, 255, 0.3);
  background: rgba(0, 40, 80, 0.6);
}

.stat-item-2 .stat-value {
  font-size: 28px;
  font-weight: 700;
  line-height: 1.2;
}

.stat-item-2 .stat-label {
  font-size: 11px;
  color: #5a8aaa;
  margin-top: 4px;
}

.stat-item-2.total .stat-value { color: #e0f4ff; }
.stat-item-2.online .stat-value { color: #34d399; }
.stat-item-2.offline .stat-value { color: #ef4444; }
.stat-item-2.rate .stat-value { color: #00d4ff; }

.device-bar {
  height: 6px;
  background: rgba(0, 30, 60, 0.8);
  border-radius: 3px;
  overflow: hidden;
  margin-top: 12px;
}

.bar-online {
  height: 100%;
  background: linear-gradient(90deg, #00d4ff, #34d399);
  border-radius: 3px;
  transition: width 1s ease;
  box-shadow: 0 0 10px rgba(52, 211, 153, 0.4);
}

.device-total-profit {
  margin-top: 10px;
  padding-top: 8px;
  border-top: 1px solid rgba(0, 212, 255, 0.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.profit-label {
  font-size: 11px;
  color: #5a8aaa;
}

.profit-mini {
  font-size: 14px;
  font-weight: 600;
  color: #fbbf24;
  letter-spacing: 0.5px;
}

/* ============================================
   设备表格（上线/离线）
   ============================================ */
.device-table {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.table-header {
  display: flex;
  align-items: center;
  padding: 6px 0;
  border-bottom: 1px solid rgba(0, 212, 255, 0.1);
  font-size: 12px;
  color: #5a8aaa;
}

.table-header span {
  flex: 1;
}

.table-scroll {
  flex: 1;
  overflow: hidden;
}

.table-row {
  display: flex;
  align-items: center;
  height: 36px;
  border-bottom: 1px solid rgba(0, 212, 255, 0.03);
  font-size: 12px;
}

.table-row span {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cell-name {
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.online-dot {
  background: #00e676;
  box-shadow: 0 0 6px rgba(0, 230, 118, 0.6);
  animation: dotPulse 2s infinite;
}

.offline-dot {
  background: #ff5252;
  box-shadow: 0 0 6px rgba(255, 82, 82, 0.4);
}

@keyframes dotPulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

.cell-mac {
  font-family: 'Courier New', monospace;
  color: #5a8aaa;
  font-size: 11px;
}

.cell-profit {
  color: #00d4ff;
  font-weight: 600;
}

.cell-error {
  color: #ffab40;
}

.cell-error.error-high {
  color: #ff5252;
  font-weight: 600;
}

.cell-status {
  text-align: center;
}

.offline-status {
  color: #ff5252;
  background: rgba(255, 82, 82, 0.1);
  padding: 2px 8px;
  border-radius: 3px;
  font-size: 11px;
}

/* ============================================
   中间地图
   ============================================ */
.map-container {
  flex: 1;
  min-height: 0;
}

.map-stats {
  display: flex;
  justify-content: center;
  gap: 40px;
  padding: 12px 0;
  background: linear-gradient(0deg, rgba(0, 30, 60, 0.4) 0%, transparent 100%);
}

.map-stat-item {
  text-align: center;
}

.map-stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #00d4ff;
  display: block;
  text-shadow: 0 0 20px rgba(0, 212, 255, 0.3);
}

.map-stat-label {
  font-size: 12px;
  color: #5a8aaa;
  display: block;
  margin-top: 2px;
}

/* ============================================
   收益饼图
   ============================================ */
.profit-chart {
  width: 100%;
  height: 100%;
  min-height: 160px;
}

/* ============================================
   总收益 - 增强版（光环 + 子指标）
   ============================================ */
.total-profit-card .card-body {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
}

.profit-ring-wrapper {
  position: relative;
  width: 140px;
  height: 140px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.profit-ring-svg {
  position: absolute;
  width: 100%;
  height: 100%;
}

.ring-anim {
  transform-origin: center;
  animation: ringRotate 8s linear infinite;
}

.ring-anim-2 {
  animation-duration: 6s;
  animation-direction: reverse;
}

@keyframes ringRotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.profit-center {
  z-index: 1;
  text-align: center;
}

.total-profit-value {
  display: flex;
  align-items: baseline;
  gap: 2px;
  justify-content: center;
}

.profit-symbol {
  font-size: 16px;
  color: #fbbf24;
  font-weight: 300;
}

.profit-number {
  font-size: 22px;
  font-weight: 700;
  color: #e0f4ff;
  letter-spacing: 1px;
}

.profit-sub-metrics {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  width: 100%;
  margin-top: 8px;
}

.sub-metric {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.sub-label {
  font-size: 10px;
  color: #5a8aaa;
}

.sub-value {
  font-size: 14px;
  font-weight: 600;
}

.sub-value.today { color: #34d399; }
.sub-value.yesterday { color: #8ab4d6; }

.sub-metric-divider {
  width: 1px;
  height: 28px;
  background: rgba(0, 212, 255, 0.15);
}

.profit-growth {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 6px;
  font-size: 11px;
  color: #34d399;
}

.growth-arrow {
  display: inline-block;
  width: 0;
  height: 0;
  border-left: 5px solid transparent;
  border-right: 5px solid transparent;
  border-bottom: 8px solid #00e676;
  animation: arrowBounce 1.5s infinite;
}

@keyframes arrowBounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-3px); }
}

/* ============================================
   返回按钮
   ============================================ */
.back-btn {
  position: fixed;
  top: 20px;
  left: 20px;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: rgba(0, 30, 60, 0.6);
  border: 1px solid rgba(0, 212, 255, 0.2);
  border-radius: 6px;
  color: #8ab4d6;
  font-size: 13px;
  cursor: pointer;
  z-index: 10000;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.back-btn:hover {
  background: rgba(0, 60, 100, 0.8);
  border-color: #00d4ff;
  color: #00d4ff;
  box-shadow: 0 0 15px rgba(0, 212, 255, 0.2);
}

/* ============================================
   全局滚动条隐藏
   ============================================ */
.datascreen::-webkit-scrollbar {
  display: none;
}
</style>
