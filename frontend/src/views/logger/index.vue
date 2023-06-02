<template>
  <div class="app-container logger-search">
    <query :loading="loading" :charts="charts" style="width: 100%" @query="onQuery" />

    <el-table v-loading="loading" :data="list" :row-class-name="tableRowClassName" size="mini" style="width: 100%" highlight-current-row>
      <el-table-column type="expand">
        <template slot-scope="props">
          <template v-if="props.row.log_type === 'access-log'">
            <kv-detail label="storage" :value="props.row.storage" :copy="true" />
            <kv-detail label="log_type" :value="props.row.log_type" :copy="true" />
            <kv-detail label="request_id" :value="props.row.log.request_id" :copy="true" />
            <kv-detail label="uri" :value="props.row.log.uri" :copy="true" />
            <kv-detail label="body" :value="props.row.log.body" :copy="true" />
            <kv-detail label="http_version" :value="props.row.log.http_version" :copy="true" />
            <kv-detail label="user_agent" :value="props.row.log.user_agent" :copy="true" />
            <kv-detail label="referer" :value="props.row.log.referer" :copy="true" />
            <kv-detail label="x_forwarded_for" :value="props.row.log.x_forwarded_for" :copy="true" />
            <kv-detail label="cookie" :value="props.row.log.cookie" :copy="true" />
            <kv-detail label="remote_addr" :value="props.row.log.remote_addr" :copy="true" />
            <kv-detail label="ip_location" :value="props.row.log.ip_location" :copy="true" />
            <kv-detail label="message" :value="props.row.log.message" :copy="true" />
          </template>
          <template v-else-if="props.row.log_type === 'json-log'">
            <kv-detail label="storage" :value="props.row.storage" :copy="true" />
            <kv-detail label="log_type" :value="props.row.log_type" :copy="true" />
            <kv-detail label="request_id" :value="props.row.log.request_id" :copy="true" />
            <kv-detail label="path" :value="props.row.log.path" :copy="true" />
            <kv-detail label="message" :value="props.row.log.message" :copy="true" />
          </template>
          <template v-else-if="props.row.log_type === 'string-log'">
            <kv-detail label="storage" :value="props.row.storage" :copy="true" />
            <kv-detail label="log_type" :value="props.row.log_type" :copy="true" />
            <kv-detail label="path" :value="props.row.log.path" :copy="true" />
            <kv-detail label="message" :value="props.row.log.message" :copy="true" />
          </template>
          <template v-else>
            <span>未知的日志类型</span>
          </template>
        </template>
      </el-table-column>

      <!-- type -->
      <el-table-column label="type" width="48">
        <template slot-scope="scope">
          <div class="svg-icon">
            <template v-if="scope.row.log_type === 'access-log'">
              <svg-icon icon-class="access_log" />
            </template>
            <template v-else-if="scope.row.log_type === 'json-log'">
              <svg-icon icon-class="json_log" />
            </template>
            <template v-else-if="scope.row.log_type === 'string-log'">
              <svg-icon icon-class="string_log" />
            </template>
          </div>
        </template>
      </el-table-column>
      <!-- time -->
      <el-table-column prop="time" label="time" width="160" :formatter="timestampFormatter" />
      <!-- application/tag/path -->
      <el-table-column label="application/tag/path" width="150">
        <template slot-scope="scope">
          <template v-if="scope.row.log_type === 'access-log'">
            <column :content="scope.row.log.http_host" />
          </template>
          <template v-else-if="scope.row.log_type === 'json-log'">
            <column :content="scope.row.log.tag" />
          </template>
          <template v-else-if="scope.row.log_type === 'string-log'">
            <column :content="scope.row.log.path" />
          </template>
        </template>
      </el-table-column>
      <!-- hostname -->
      <el-table-column prop="hostname" label="hostname" width="180">
        <template slot-scope="scope">
          <column :content="scope.row.log.hostname" :center="false" />
        </template>
      </el-table-column>
      <!-- uri/level -->
      <el-table-column prop="level" label="uri/level" width="240">
        <template slot-scope="scope">
          <template v-if="scope.row.log_type === 'access-log'">
            <column :content="scope.row.log.uri" :copy="true" :tooltip="true" :center="false" />
          </template>
          <template v-else-if="scope.row.log_type === 'json-log'">
            <column :content="scope.row.log.level" :copy="false" :tooltip="false" :center="true" />
          </template>
          <template v-else-if="scope.row.log_type === 'string-log'">
            <span class="one-line" />
          </template>
        </template>
      </el-table-column>
      <!-- message -->
      <el-table-column prop="message" label="message">
        <template slot-scope="scope">
          <template v-if="scope.row.log_type === 'access-log'">
            <column :tooltip="false" :content="accessLogMessageFormatter(scope.row)" :center="false" />
          </template>
          <template v-else-if="scope.row.log_type === 'json-log'">
            <column :tooltip="false" :content="scope.row.log.message" :center="false" />
          </template>
          <template v-else-if="scope.row.log_type === 'string-log'">
            <column :tooltip="false" :content="scope.row.log.message" :center="false" />
          </template>
        </template>
      </el-table-column>
      <!-- request_id -->
      <el-table-column prop="request_id" label="request_id" width="100">
        <template slot-scope="scope">
          <span v-if="scope.row.log_type === 'string-log'" class="one-line" />
          <column v-else :content="scope.row.log.request_id" :color="hashStringColor(scope.row.log.request_id)" />
        </template>
      </el-table-column>

      <!-- 操作 -->
      <el-table-column label="操作" width="100">
        <template slot-scope="scope">
          <el-button type="text" size="mini" @click="handleDetail(logFormatter(scope.row))">格式化</el-button>
          <el-button type="text" size="mini" @click="handleDetail(jsonFormatter(scope.row))">原日志</el-button>
        </template>
      </el-table-column>

    </el-table>

    <div class="pagination">
      <el-pagination
        layout="sizes, prev, pager, next, total, jumper"
        :current-page="form.page_no"
        :page-sizes="[10, 20, 50, 100, 200]"
        :page-size="form.page_size"
        :total="form.count"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <el-dialog width="80%" :visible.sync="dialog_visible">
      <div v-if="detail" class="code-group">
        <div class="container">
          <code class="code">{{ detail }}</code>
        </div>
      </div>
      <div style="text-align:right;">
        <el-button type="danger" @click="dialog_visible=false">关闭</el-button>
      </div>
    </el-dialog>

  </div>
</template>

<script>
import Column from './Column.vue'
import KvDetail from './KvDetail.vue'
import Query from './Query.vue'
import clipboard from '@/directive/clipboard/index.js'
import { date, formatInteger, time33 } from '@/utils'
import { jsonFormatter } from '@/filters'
import { loggerSearch } from '@/api/logger'

export default {
  components: { Column, KvDetail, Query },
  directives: { clipboard },
  data() {
    return {
      loading: false,
      list: [],
      total: 0,
      form: {
        page_no: 1,
        page_size: 10,
        count: 0
      },
      lastQuery: null,
      detail: null,
      dialog_visible: false,
      charts: {}
    }
  },
  watch: {
  },
  async created() {
  },
  methods: {
    parseQueryString(url) {
      const obj = {}
      const form = url.substring(url.indexOf('?') + 1, url.length).split('&')
      for (const i in form) {
        const kv = form[i].split('=')
        if (kv[0]) {
          obj[kv[0]] = decodeURIComponent(kv[1])
        }
      }
      return obj
    },
    logFormatter(row) {
      if (row.log_type === 'string-log') {
        return row.log.message
      }
      if (row.log_type === 'access-log') {
        const uri = row.log.uri
        const query = JSON.stringify(this.parseQueryString('?' + row.log.query), null, 2)
        const http_body = row.log.body
        if (http_body.indexOf('<xml') === 0) {
          return `URI:\n${uri}\n\nquery to JSON:\n${query}\n\nXML body:\n${http_body}`
        }
        let body
        if (http_body.indexOf('{') === 0) {
          body = JSON.stringify(JSON.parse(http_body), null, 2)
        } else {
          body = JSON.stringify(this.parseQueryString('?' + row.log.body), null, 2)
        }
        return `URI:\n${uri}\n\nquery to JSON:\n${query}\n\nbody to JSON:\n${body}\n\n${row.log.message}`
      }
      const message = JSON.parse(row.log.message)
      return JSON.stringify(message, null, 2)
    },
    jsonFormatter(v) { return jsonFormatter(JSON.stringify(v)) },
    timestampFormatter(row) {
      const time = +(new Date(row.log.time))
      const ts = (time * 0.001).toFixed()
      const ms = formatInteger(time % 1000, 3, '0')
      return date('Y-m-d H:i:s', ts) + '.' + ms
    },
    accessLogMessageFormatter(row) {
      const s = row.log
      return `${s.method} [${s.remote_addr} ${s.ip_location}] [status:${s.status}] [${s.duration}] `
    },
    hashStringColor(s) {
      const colors = ['#2a9d2a', '#645b93', '#ff7f6a', '#5fc0ea', '#480048', '#601848', '#c04848', '#f07241', '#c71585', '#008b8b', '#7b68ee', '#ff7f50']
      const i = time33(s) % (colors.length - 1)
      return colors[i]
    },
    tableRowClassName({ row }) {
      console.log(row.log_type)
      if (row.log_type === 'access-log') {
        if (row.log.status >= 500) {
          return 'error-row'
        }
        if (row.log.status >= 400) {
          return 'warning-row'
        }
        if (row.log.status >= 300) {
          return 'info-row'
        }
        return 'success-row'
      } else if (row.log_type === 'json-log') {
        if (row.log.level === 'fatal' || row.log.level === 'error') {
          return 'error-row'
        }
        if (row.log.level === 'warning' || row.log.level === 'warn') {
          return 'warning-row'
        }
        if (row.log.level === 'info') {
          return 'info-row'
        }
        return 'success-row'
      } else if (row.log_type === 'string-log') {
        const message = row.log.message.toLowerCase()
        if (message.indexOf('fatal') >= 0 || message.indexOf('error') >= 0) {
          return 'error-row'
        }
        if (message.indexOf('warning') >= 0 || message.indexOf('warn') >= 0) {
          return 'warning-row'
        }
        if (message.indexOf('info') >= 0) {
          return 'info-row'
        }
        return 'success-row'
      } else {
        return ''
      }
    },
    async onQuery(data) {
      this.lastQuery = data
      this.loading = true
      try {
        const params = data.params
        params.page_no = this.form.page_no
        params.page_size = this.form.page_size
        const res = await loggerSearch(params)
        if (res.code !== 0) {
          this.$message({ type: 'error', message: res.message })
        } else {
          this.form.page_no = res.data.page_no
          this.form.page_size = res.data.page_size
          this.form.count = res.data.count
          this.list = res.data ? res.data.list : []
          this.charts = res.data ? res.data.charts : []
        }
      } finally {
        this.loading = false
      }
    },
    handleSizeChange(val) {
      this.form.page_size = val
      if (this.lastQuery) this.onQuery(this.lastQuery)
    },
    handleCurrentChange(val) {
      this.form.page_no = val
      if (this.lastQuery) this.onQuery(this.lastQuery)
    },
    handleDetail(detail) {
      if (!detail) return
      this.detail = detail
      this.dialog_visible = true
    }
  }
}
</script>

<style lang="scss">
.logger-search {
  width: 100%;

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
  .el-table .error-row {
    background: #ffe8e8;
  }
  .el-table .warning-row {
    background: oldlace;
  }
  .el-table .success-row {
    background: #f0f9ec;
  }
  .el-table .info-row {
    background: #f1faff;
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
.svg-icon {
  margin-left: 5px;
}
.pagination {
  margin: 0 10px;
}
</style>
