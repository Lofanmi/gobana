<template>
  <span>
    <template v-if="copy">
      <el-tooltip v-if="tooltip && computedValue !== '-'" :content="computedValue" placement="top">
        <span v-clipboard:copy="computedValue" v-clipboard:success="clipboardSuccess" :style="style" title="点击复制" class="one-line copy-line" :class="className">
          {{ computedValue }}
        </span>
      </el-tooltip>
      <span v-else v-clipboard:copy="computedValue" v-clipboard:success="clipboardSuccess" :style="style" title="点击复制" class="one-line copy-line" :class="className">
        {{ computedValue }}
      </span>
    </template>
    <template v-else>
      <el-tooltip v-if="tooltip && computedValue !== '-'" :content="computedValue" placement="top">
        <span :style="style" class="one-line" :class="className">
          {{ computedValue }}
        </span>
      </el-tooltip>
      <span v-else :style="style" class="one-line" :class="className">
        {{ computedValue }}
      </span>
    </template>
  </span>
</template>

<script>
import clipboard from '@/directive/clipboard/index.js'

export default {
  directives: { clipboard },
  props: {
    content: {
      type: String,
      default: ''
    },
    color: {
      type: String,
      default: ''
    },
    tooltip: {
      type: Boolean,
      default: true
    },
    copy: {
      type: Boolean,
      default: true
    },
    center: {
      type: Boolean,
      default: true
    }
  },
  computed: {
    computedValue() {
      return this.content || '-'
    },
    className() {
      return this.center ? 'text-center' : ''
    },
    style() {
      return this.color ? { color: this.color } : ''
    }
  },
  methods: {
    clipboardSuccess() { this.$message({ message: '复制成功', type: 'success', duration: 1500 }) }
  }
}
</script>

<style lang="scss" scoped>
.one-line {
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
  display: block;
}
.copy-line:hover {
  cursor: pointer;
  text-decoration: underline;
  color: #1890ff;
}
.kv-detail {
  line-height: 20px;
  .key {
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
