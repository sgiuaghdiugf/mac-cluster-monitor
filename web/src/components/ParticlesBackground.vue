<template>
  <canvas ref="canvasRef" class="particles-canvas"></canvas>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const canvasRef = ref(null)
let animationId = null
let particles = []
let mouse = { x: null, y: null }

class Particle {
  constructor(canvas) {
    this.canvas = canvas
    this.x = Math.random() * canvas.width
    this.y = Math.random() * canvas.height
    this.vx = (Math.random() - 0.5) * 0.5
    this.vy = (Math.random() - 0.5) * 0.5
    this.radius = Math.random() * 2 + 1
    this.opacity = Math.random() * 0.5 + 0.2
    
    // 随机颜色（蓝色系）
    const colors = [
      '59, 130, 246',   // 蓝色
      '6, 182, 212',    // 青色
      '99, 102, 241',   // 靛蓝
      '139, 92, 246'    // 紫色
    ]
    this.color = colors[Math.floor(Math.random() * colors.length)]
  }

  update() {
    this.x += this.vx
    this.y += this.vy

    // 边界反弹
    if (this.x < 0 || this.x > this.canvas.width) {
      this.vx *= -1
    }
    if (this.y < 0 || this.y > this.canvas.height) {
      this.vy *= -1
    }

    // 鼠标交互：粒子被鼠标吸引/排斥
    if (mouse.x !== null && mouse.y !== null) {
      const dx = mouse.x - this.x
      const dy = mouse.y - this.y
      const distance = Math.sqrt(dx * dx + dy * dy)
      
      if (distance < 150) {
        const force = (150 - distance) / 150
        this.vx += dx / distance * force * 0.02
        this.vy += dy / distance * force * 0.02
      }
    }

    // 限制速度
    const speed = Math.sqrt(this.vx * this.vx + this.vy * this.vy)
    if (speed > 2) {
      this.vx = (this.vx / speed) * 2
      this.vy = (this.vy / speed) * 2
    }
  }

  draw(ctx) {
    ctx.beginPath()
    ctx.arc(this.x, this.y, this.radius, 0, Math.PI * 2)
    ctx.fillStyle = `rgba(${this.color}, ${this.opacity})`
    ctx.fill()
    
    // 发光效果
    ctx.shadowBlur = 10
    ctx.shadowColor = `rgba(${this.color}, ${this.opacity})`
  }
}

const initParticles = (canvas) => {
  particles = []
  const particleCount = Math.min(80, (canvas.width * canvas.height) / 10000)
  
  for (let i = 0; i < particleCount; i++) {
    particles.push(new Particle(canvas))
  }
}

const drawLines = (ctx, canvas) => {
  for (let i = 0; i < particles.length; i++) {
    for (let j = i + 1; j < particles.length; j++) {
      const dx = particles[i].x - particles[j].x
      const dy = particles[i].y - particles[j].y
      const distance = Math.sqrt(dx * dx + dy * dy)

      if (distance < 120) {
        const opacity = (120 - distance) / 120 * 0.3
        ctx.beginPath()
        ctx.moveTo(particles[i].x, particles[i].y)
        ctx.lineTo(particles[j].x, particles[j].y)
        ctx.strokeStyle = `rgba(59, 130, 246, ${opacity})`
        ctx.lineWidth = 1
        ctx.stroke()
      }
    }
  }
}

const animate = () => {
  const canvas = canvasRef.value
  if (!canvas) return
  
  const ctx = canvas.getContext('2d')
  
  // 清空画布，带拖尾效果
  ctx.fillStyle = 'rgba(15, 23, 42, 0.1)'
  ctx.fillRect(0, 0, canvas.width, canvas.height)

  // 更新和绘制粒子
  particles.forEach(particle => {
    particle.update()
    particle.draw(ctx)
  })

  // 绘制连线
  drawLines(ctx, canvas)

  animationId = requestAnimationFrame(animate)
}

const resizeCanvas = () => {
  const canvas = canvasRef.value
  if (!canvas) return
  
  canvas.width = window.innerWidth
  canvas.height = window.innerHeight
  
  initParticles(canvas)
}

onMounted(() => {
  resizeCanvas()
  animate()

  // 鼠标移动事件
  window.addEventListener('mousemove', (e) => {
    mouse.x = e.clientX
    mouse.y = e.clientY
  })

  window.addEventListener('mouseout', () => {
    mouse.x = null
    mouse.y = null
  })

  // 窗口大小改变事件
  window.addEventListener('resize', resizeCanvas)
})

onUnmounted(() => {
  if (animationId) {
    cancelAnimationFrame(animationId)
  }
  window.removeEventListener('resize', resizeCanvas)
})
</script>

<style scoped>
.particles-canvas {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  pointer-events: none;
}
</style>
