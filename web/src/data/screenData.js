/**
 * 数据大屏 - 静态数据配置文件
 * 后续可修改此文件中的数据来调整大屏展示内容
 */

// ============================================
// 规模指标分析 - 用户排行（展示前5名）
// ============================================
export const scaleRanking = [
  { rank: 1, name: '张先生', deviceCount: 12, profit: 856.50 },
  { rank: 2, name: '李先生', deviceCount: 10, profit: 723.80 },
  { rank: 3, name: '王先生', deviceCount: 8, profit: 589.20 },
  { rank: 4, name: '赵先生', deviceCount: 6, profit: 432.60 },
  { rank: 5, name: '陈先生', deviceCount: 5, profit: 368.90 }
]

// ============================================
// 设备总数量 - 在线/离线统计
// ============================================
export const deviceStats = {
  total: 128,          // 设备总数
  online: 96,          // 在线设备数
  offline: 32,         // 离线设备数
  onlineRate: 75.0     // 在线率（百分比）
}

// ============================================
// 设备上线情况 - 在线设备列表
// ============================================
export const onlineDevices = [
  { name: 'Mac Pro-001', mac: 'A1:B2:C3:D4:E5:F6', duration: '48小时32分', profit: 126.50 },
  { name: 'Mac Pro-002', mac: 'B2:C3:D4:E5:F6:A7', duration: '36小时15分', profit: 98.30 },
  { name: 'Mac Studio-003', mac: 'C3:D4:E5:F6:A7:B8', duration: '72小时08分', profit: 189.60 },
  { name: 'Mac Pro-004', mac: 'D4:E5:F6:A7:B8:C9', duration: '24小时45分', profit: 65.80 },
  { name: 'Mac Mini-005', mac: 'E5:F6:A7:B8:C9:D0', duration: '56小时22分', profit: 143.20 },
  { name: 'Mac Pro-006', mac: 'F6:A7:B8:C9:D0:E1', duration: '12小时30分', profit: 32.40 },
  { name: 'Mac Studio-007', mac: 'A7:B8:C9:D0:E1:F2', duration: '64小时18分', profit: 167.90 },
  { name: 'Mac Pro-008', mac: 'B8:C9:D0:E1:F2:A3', duration: '8小时52分', profit: 21.60 }
]

// ============================================
// 设备离线情况 - 离线设备列表
// ============================================
export const offlineDevices = [
  { name: 'Mac Pro-009', mac: 'C9:D0:E1:F2:A3:B4', errorCount: 3, status: '离线' },
  { name: 'Mac Mini-010', mac: 'D0:E1:F2:A3:B4:C5', errorCount: 1, status: '离线' },
  { name: 'Mac Studio-011', mac: 'E1:F2:A3:B4:C5:D6', errorCount: 5, status: '离线' },
  { name: 'Mac Pro-012', mac: 'F2:A3:B4:C5:D6:E7', errorCount: 2, status: '离线' },
  { name: 'Mac Mini-013', mac: 'A3:B4:C5:D6:E7:F8', errorCount: 0, status: '离线' },
  { name: 'Mac Pro-014', mac: 'B4:C5:D6:E7:F8:A9', errorCount: 4, status: '离线' }
]

// ============================================
// 收益情况 - 四季度饼图数据
// ============================================
export const profitByQuarter = [
  { name: '第一季度', value: 286.50 },
  { name: '第二季度', value: 352.80 },
  { name: '第三季度', value: 418.20 },
  { name: '第四季度', value: 325.60 }
]

// ============================================
// 总收益 - 动态增长配置
// ============================================
export const totalProfit = {
  initialValue: 1383.10,   // 总收益初始值
  increment: 0.05,         // 每5秒增加的金额
  interval: 5000           // 增加间隔（毫秒）
}

// ============================================
// 中国地图 - 设备分布数据（省份+设备数量）
// ============================================
export const mapDeviceData = [
  { name: '北京', value: 18 },
  { name: '上海', value: 15 },
  { name: '广东', value: 22 },
  { name: '浙江', value: 12 },
  { name: '江苏', value: 10 },
  { name: '四川', value: 8 },
  { name: '湖北', value: 6 },
  { name: '福建', value: 5 },
  { name: '山东', value: 7 },
  { name: '河南', value: 4 },
  { name: '湖南', value: 3 },
  { name: '安徽', value: 3 },
  { name: '辽宁', value: 4 },
  { name: '陕西', value: 3 },
  { name: '重庆', value: 2 },
  { name: '天津', value: 2 },
  { name: '云南', value: 1 },
  { name: '广西', value: 1 },
  { name: '山西', value: 2 },
  { name: '贵州', value: 1 }
]
