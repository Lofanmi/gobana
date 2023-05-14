<template>
  <div class="app-container logger-search">
    <query />
  </div>
</template>

<script>
import Query from './Query.vue'
import clipboard from '@/directive/clipboard/index.js'
// import { date, deepClone } from '@/utils'
// import { jsonFormatter } from '@/filters'
// import clip from '@/utils/clipboard'

export default {
  components: { Query },
  directives: { clipboard },
  data() {
    return {
      loading: false,
      list: [],
      total: 0,
      form: {
        page_no: 1,
        page_size: 20,
        count: 0
      },
      detail: null,
      dialog_visible: false
    }
  },
  computed: {
    // index() {
    //   if (this.form.search_options.length <= 0) {
    //     return `sy_*`
    //   }
    //   return this.form.search_options.map(i => `sy_*${i}*`).join(',')
    // }
  },
  watch: {
  },
  async created() {
  },
  methods: {
  }
}
</script>

<style lang="scss">
.logger-search {
  .el-table .cell {
    color: #333333 !important;
  }
  .el-table .is-leaf .cell {
    color: #555555 !important;
  }
  .el-table .one-line {
    white-space: nowrap;
    text-overflow: ellipsis;
    overflow: hidden;
    display: block;
  }
  .el-table__body-wrapper {
    overflow: auto;
  }
  .el-table .cell {
    padding: 0 0 0 5px;
  }
  .el-table--mini th, .el-table--mini td {
    padding: 0;
  }
  .el-table__expanded-cell[class*=cell] {
    padding: 10px 50px;
  }
}
</style>

<style lang="scss" scoped>
.replay-label {
  font-size: 12px;
  color: #606266;
  padding: 0 7px;
  font-weight: bold;
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
.pagination {
  margin: 0 10px;
}
</style>
