import { useUserStore } from '@/stores/user'

/** 无权限时隐藏元素，勿 removeChild，否则会破坏 Vue 对表格等节点的 patch */
function apply(el, code) {
  if (!code) return
  const store = useUserStore()
  const allowed = store.hasPerm(code)
  if (allowed) {
    el.style.removeProperty('display')
    el.removeAttribute('data-auth-hidden')
  } else {
    el.style.display = 'none'
    el.setAttribute('data-auth-hidden', '1')
  }
}

export default {
  mounted(el, binding) {
    apply(el, binding.value)
  },
  updated(el, binding) {
    apply(el, binding.value)
  },
}
