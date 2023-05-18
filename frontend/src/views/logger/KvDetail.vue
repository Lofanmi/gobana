<template>
  <div class="kv-detail">
    <span class="label">{{ label }}:</span>
    <span
      v-if="copy && computedValue !== '-'"
      v-clipboard:copy="computedValue"
      v-clipboard:success="clipboardSuccess"
      title="点击复制"
      class="value copy-line"
    >
      {{ computedValue }}
    </span>
    <span v-else class="value">{{ computedValue }}</span>
  </div>
</template>

<script>
import clipboard from '@/directive/clipboard/index.js'

export default {
  directives: { clipboard },
  props: {
    label: {
      type: String,
      default: ''
    },
    value: {
      type: String,
      default: ''
    },
    copy: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    computedValue() {
      return this.value || '-'
    }
  },
  methods: {
    clipboardSuccess() { this.$message({ message: '复制成功', type: 'success', duration: 1500 }) }
  }
}
</script>

<style lang="scss" scoped>
.copy-line:hover {
  cursor: pointer;
  text-decoration: underline;
  color: #1890ff;
}
.kv-detail {
  line-height: 20px;
  .label {
    font-weight: bold;
  }
  .value {
    word-break: break-all;
    color: #333333 !important;
    &.copy-line:hover {
      color: #1890ff !important;
    }
    &.font-bold {
      font-weight: bold;
    }
  }
}
</style>
