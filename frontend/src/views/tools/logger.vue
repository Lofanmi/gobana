<template>
  <div class="app-container logger-search">

    <el-form :inline="true" label-width="150px" size="mini">
      <el-form-item label="或者 [A || B || C]">
        <el-input v-model="form.or1" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.or2" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.or3" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.or4" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.or5" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.or6" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.or7" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.or8" class="search-input" @keyup.enter.native="getList(1)" />
      </el-form-item>
      <el-form-item label="并且 [A && B && C]">
        <el-input v-model="form.must1" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.must2" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.must3" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.must4" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.must5" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.must6" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.must7" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.must8" class="search-input" @keyup.enter.native="getList(1)" />
      </el-form-item>
      <el-form-item label="不包含 [!A && !B && !C]">
        <el-input v-model="form.must_not1" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.must_not2" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.must_not3" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.must_not4" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.must_not5" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.must_not6" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.must_not7" class="search-input" @keyup.enter.native="getList(1)" />
        <el-input v-model="form.must_not8" class="search-input" @keyup.enter.native="getList(1)" />
      </el-form-item>

      <el-form-item label="time">
        <el-date-picker v-model="timerange" type="datetimerange" start-placeholder="开始时间" end-placeholder="结束时间" :picker-options="picker_options" style="width: 324px" />
        <el-button-group>
          <el-button size="mini" @click="getList(1)">查询</el-button>
          <el-button size="mini" @click="getList(2)">最新</el-button>
        </el-button-group>
        <el-button-group>
          <el-button size="mini" @click="getList(3)">今日</el-button>
          <el-button size="mini" @click="getList(4)">昨日</el-button>
        </el-button-group>
        <el-button-group>
          <el-button v-permission="['role:admin','role:limited','permission:LoggerExport']" size="mini" @click="exportExcel">导出</el-button>
        </el-button-group>
        <el-button-group>
          <el-button size="mini" @click="changeChartsVisible">绘图</el-button>
        </el-button-group>

        <el-select v-model="form.search_options" multiple collapse-tags placeholder="日志分类" style="width: 160px">
          <el-option value="nginx" label="web日志" />
          <el-option value="json" label="程序日志" />
          <el-option value="string" label="字符串日志" />
        </el-select>

        <el-select v-model="form.without_options" multiple collapse-tags placeholder="过滤日志" style="width: 200px">
          <el-option value="total" label="不看总数" />
          <el-option value="track" label="过滤埋点" />
          <el-option value="orange" label="过滤网关本地转发" />
          <el-option value="ht" label="过滤发布机/后台" />
          <el-option value="tg" label="过滤广告推广" />
          <el-option value="source" label="过滤原始日志" />
        </el-select>

        <span class="replay-label">请求回放到机器</span>
        <el-select v-model="replay_host" allow-create filterable default-first-option placeholder="请选择机器" class="search-input">
          <el-option v-for="item in replay_servers" :key="item.value" :value="item.value" :label="item.label" />
        </el-select>

      </el-form-item>

    </el-form>

    <div v-show="charts_visible && charts.series && charts.series.data && charts.series.data.length > 0" v-loading="loading" class="chart" style="background-color: #ffffff;">
      <counter id="line" ref="counter" class="line" :height="'240px'" :width="'95%'" :chartdata="charts" @click-event="barClick" />
    </div>

    <el-table v-loading="loading" :data="list" size="mini" style="width: 100%" highlight-current-row>
      <el-table-column type="expand">
        <template slot-scope="props">
          <!-- index -->
          <div class="kv-detail">
            <span class="key">index:</span>
            <span class="value">{{ props.row.index | def }}</span>
          </div>
          <!-- url -->
          <div class="kv-detail">
            <span class="key">url:</span>
            <span v-clipboard:copy="props.row.level + (props.row.http_query ? '?' + props.row.http_query : '')" v-clipboard:success="clipboardSuccess" title="点击复制" class="value copy-line">
              {{ props.row.level + (props.row.http_query ? '?' + props.row.http_query : '') | def }}
            </span>
          </div>
          <!-- request_id -->
          <div class="kv-detail">
            <span class="key">request_id:</span>
            <span v-clipboard:copy="props.row.request_id" v-clipboard:success="clipboardSuccess" title="点击复制" class="value copy-line">
              {{ props.row.request_id | def }}
            </span>
          </div>
          <!-- log_path -->
          <div class="kv-detail">
            <span class="key">log_path:</span>
            <span class="value">{{ props.row.log_path | def }}</span>
          </div>
          <!-- nginx -->
          <template v-if="props.row.log_type === 'nginx'">
            <!-- http_version -->
            <div v-if="props.row.http_version" class="kv-detail">
              <span class="key">http_version:</span>
              <span class="value">
                {{ props.row.http_version | def }}
              </span>
            </div>
            <!-- http_user_agent -->
            <div v-if="props.row.http_user_agent" class="kv-detail">
              <span class="key">http_user_agent:</span>
              <span v-clipboard:copy="props.row.http_user_agent" v-clipboard:success="clipboardSuccess" title="点击复制" class="value copy-line">
                {{ props.row.http_user_agent | def }}
              </span>
            </div>
            <!-- http_referer -->
            <div v-if="props.row.http_referer" class="kv-detail">
              <span class="key">http_referer:</span>
              <span v-clipboard:copy="props.row.http_referer" v-clipboard:success="clipboardSuccess" title="点击复制" class="value copy-line">
                {{ props.row.http_referer | def }}
              </span>
            </div>
            <!-- http_cookie -->
            <div v-if="props.row.http_cookie" class="kv-detail">
              <span class="key">http_cookie:</span>
              <span v-clipboard:copy="props.row.http_cookie" v-clipboard:success="clipboardSuccess" title="点击复制" class="value copy-line">
                {{ props.row.http_cookie | def }}
              </span>
            </div>
            <!-- http_body -->
            <div class="kv-detail">
              <span class="key">http_body:</span>
              <span v-clipboard:copy="props.row.http_body" v-clipboard:success="clipboardSuccess" title="点击复制" class="value copy-line">
                {{ props.row.http_body | def }}
              </span>
            </div>
          </template>
          <!-- json -->
          <template v-if="props.row.log_type === 'json'">
            <!-- raw -->
            <div class="kv-detail">
              <span class="key">raw:</span>
              <span v-clipboard:copy="props.row.detail" v-clipboard:success="clipboardSuccess" title="点击复制" class="value copy-line">
                {{ props.row.detail | def }}
              </span>
            </div>
          </template>
          <!-- string -->
          <template v-if="props.row.log_type === 'string'">
            <!-- raw -->
            <div class="kv-detail">
              <span class="key">raw:</span>
              <span v-clipboard:copy="props.row.detail" v-clipboard:success="clipboardSuccess" title="点击复制" class="value copy-line">
                {{ props.row.detail | def }}
              </span>
            </div>
          </template>
          <!-- 复制请求/回放请求 -->
          <div v-if="props.row.log_type === 'nginx'" class="kv-detail">
            <br>
            <span title="点击复制" class="value copy-line" @click="copy($event, props.row, 'url')">[仅复制URL]</span>
            <span title="点击复制" class="value copy-line" @click="copy($event, props.row, 'query')">[仅复制query]</span>
            <span title="点击复制" class="value copy-line" @click="copy($event, props.row, 'body')">[仅复制body]</span>
            <br>
            <span title="点击复制" class="value copy-line font-bold" @click="copy($event, props.row, 'cURL-raw')">[复制线上cURL-原始协议]</span>
            <span title="点击复制" class="value copy-line" @click="copy($event, props.row, 'cURL-http')">[复制线上cURL-http]</span>
            <span title="点击复制" class="value copy-line" @click="copy($event, props.row, 'cURL-https')">[复制线上cURL-https]</span>
            <span title="点击复制" class="value copy-line font-bold" @click="copy($event, props.row, 'cURL2-raw')">[复制指定机器cURL-原始协议]</span>
            <span title="点击复制" class="value copy-line" @click="copy($event, props.row, 'cURL2-http')">[复制指定机器cURL-http]</span>
            <span title="点击复制" class="value copy-line" @click="copy($event, props.row, 'cURL2-https')">[复制指定机器cURL-https]</span>
            <br>
            <span title="点击回放" class="value copy-line font-bold" @click="replay(props.row, '', '')">[回放到线上-原始协议]</span>
            <span title="点击回放" class="value copy-line" @click="replay(props.row, 'http', '')">[回放到线上-http]</span>
            <span title="点击回放" class="value copy-line" @click="replay(props.row, 'https', '')">[回放到线上-https]</span>
            <span title="点击回放" class="value copy-line font-bold" @click="replay(props.row, '', replay_host)">[回放到指定机器-原始协议]</span>
            <span title="点击回放" class="value copy-line" @click="replay(props.row, 'http', replay_host)">[回放到指定机器-http]</span>
            <span title="点击回放" class="value copy-line" @click="replay(props.row, 'https', replay_host)">[回放到指定机器-https]</span>
          </div>

        </template>
      </el-table-column>

      <!-- time -->
      <el-table-column prop="time" label="time" width="150" :formatter="timestampFormatter" />
      <!-- tag -->
      <el-table-column label="tag" width="150">
        <template slot-scope="scope">
          <span v-if="scope.row.log_type !== 'json'" class="one-line text-center">
            {{ scope.row.tag | def }}
          </span>
          <el-tooltip v-else :content="scope.row.tag" placement="top">
            <span v-clipboard:copy="scope.row.tag" v-clipboard:success="clipboardSuccess" title="点击复制" class="one-line copy-line text-center">
              {{ scope.row.tag | def }}
            </span>
          </el-tooltip>
        </template>
      </el-table-column>
      <!-- server -->
      <el-table-column prop="server" label="server" width="180">
        <template slot-scope="scope">
          <el-tooltip :content="scope.row.server" placement="top">
            <span v-clipboard:copy="scope.row.server" v-clipboard:success="clipboardSuccess" title="点击复制" class="one-line copy-line text-center">
              {{ scope.row.server | def }}
            </span>
          </el-tooltip>
        </template>
      </el-table-column>
      <!-- level -->
      <el-table-column prop="level" label="level" width="240">
        <template slot-scope="scope">
          <span v-if="scope.row.level.indexOf('http') !== 0" class="one-line text-center">
            {{ scope.row.level | def }}
          </span>
          <el-tooltip v-else :content="scope.row.level" placement="top">
            <span v-clipboard:copy="scope.row.level" v-clipboard:success="clipboardSuccess" title="点击复制" class="one-line copy-line text-center">
              {{ scope.row.level | def }}
            </span>
          </el-tooltip>
        </template>
      </el-table-column>
      <!-- message -->
      <el-table-column prop="message" label="message">
        <template slot-scope="scope">
          <span class="one-line">
            {{ scope.row.message | def }}
          </span>
        </template>
      </el-table-column>
      <!-- request_id -->
      <el-table-column prop="request_id" label="request_id" width="100">
        <template slot-scope="scope">
          <el-tooltip v-if="scope.row.request_id" :content="scope.row.request_id" placement="top">
            <span v-clipboard:copy="scope.row.request_id" v-clipboard:success="clipboardSuccess" title="点击复制" class="one-line copy-line text-center">
              {{ scope.row.request_id | def }}
            </span>
          </el-tooltip>
          <span v-else> {{ scope.row.request_id | def }}</span>
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
import { date, deepClone } from '@/utils'
import { lasthour, today, yesterday, lastday } from '@/utils/timeshortcut'
import { loggerSearch, loggerReplay, loggerExport } from '@/api/logger'
import { cmdbServersRemote } from '@/api/cmdb'
import { jsonFormatter } from '@/filters'
import clip from '@/utils/clipboard'
import clipboard from '@/directive/clipboard/index.js'
import { applyVmByQuery, applyURLByForm } from '@/utils/form-binding'
import Counter from './Counter.vue'
import Vue from 'vue'

export default {
  components: { Counter },
  directives: { clipboard },
  data() {
    return {
      loading: false,
      timerange: [],
      list: [],
      charts: {},
      charts_visible: false,
      total: 0,
      form: {
        page_no: 1,
        page_size: 20,
        count: 0,
        elastic_search: '', index: '',
        or1: '', or2: '', or3: '', or4: '', or5: '', or6: '', or7: '', or8: '',
        must1: '', must2: '', must3: '', must4: '', must5: '', must7: '', must6: '', must8: '',
        must_not1: '', must_not2: '', must_not3: '', must_not4: '', must_not5: '', must_not6: '', must_not7: '', must_not8: '',
        time_a: 0, time_b: 0,
        search_options: [],
        without_options: []
      },
      detail: null,
      dialog_visible: false,
      picker_options: {
        shortcuts: [{
          text: '近15分钟',
          onClick(picker) {
            picker.$emit('pick', lasthour(0.25))
          }
        }, {
          text: '近30分钟',
          onClick(picker) {
            picker.$emit('pick', lasthour(0.5))
          }
        }, {
          text: '近1小时',
          onClick(picker) {
            picker.$emit('pick', lasthour(1))
          }
        }, {
          text: '近3小时',
          onClick(picker) {
            picker.$emit('pick', lasthour(3))
          }
        }, {
          text: '今天',
          onClick(picker) {
            picker.$emit('pick', today())
          }
        }, {
          text: '昨天',
          onClick(picker) {
            picker.$emit('pick', yesterday())
          }
        }, {
          text: '近1天',
          onClick(picker) {
            picker.$emit('pick', lastday(1))
          }
        }, {
          text: '近2天',
          onClick(picker) {
            picker.$emit('pick', lastday(2))
          }
        }, {
          text: '近1周',
          onClick(picker) {
            picker.$emit('pick', lastday(7))
          }
        }]
      },
      replay_servers: [],
      replay_host: ''
    }
  },
  computed: {
    or() {
      return [this.form.or1, this.form.or2, this.form.or3, this.form.or4, this.form.or5, this.form.or6, this.form.or7, this.form.or8].filter(i => !!i)
    },
    must() {
      return [this.form.must1, this.form.must2, this.form.must3, this.form.must4, this.form.must5, this.form.must6, this.form.must7, this.form.must8].filter(i => !!i)
    },
    must_not() {
      return [this.form.must_not1, this.form.must_not2, this.form.must_not3, this.form.must_not4, this.form.must_not5, this.form.must_not6, this.form.must_not7, this.form.must_not8].filter(i => !!i)
    },
    index() {
      if (this.form.search_options.length <= 0) {
        return `sy_*`
      }
      return this.form.search_options.map(i => `sy_*${i}*`).join(',')
    }
  },
  watch: {
    replay_host() {
      localStorage.setItem('user_host', this.replay_host)
    },
    timerange: {
      handler(v) {
        this.form.time_a = +(v[0])
        this.form.time_b = +(v[1])
      },
      deep: true
    },
    form: {
      handler() {
        applyURLByForm(this, this.form,
          [
            'elastic_search', 'index',
            'or1', 'or2', 'or3', 'or4', 'or5', 'or6', 'or7', 'or8',
            'must1', 'must2', 'must3', 'must4', 'must5', 'must7', 'must6', 'must8',
            'must_not1', 'must_not2', 'must_not3', 'must_not4', 'must_not5', 'must_not6', 'must_not7', 'must_not8',
            'time_a', 'time_b',
            'search_options', 'without_options'
          ])
      },
      deep: true
    },
    charts_visible() {
      this.$nextTick(() => {
        this.$refs.counter.resizeChart()
      })
    },
    charts() {
      this.$nextTick(() => {
        this.$refs.counter.resizeChart()
      })
    }
  },
  async created() {
    await this.getReplayServers()
    const userHost = localStorage.getItem('user_host')
    if (userHost) {
      this.replay_host = userHost
    }
    const query = deepClone(this.$route.query)
    applyVmByQuery(this.form, query, ['time_a', 'time_b'], ['search_options', 'without_options'])

    if (this.form.time_a < 1640966400000 && this.form.time_b < 1640966400000) {
      await this.setTimeRange(lasthour(0.25))
    } else {
      await this.setTimeRange([new Date(this.form.time_a), new Date(this.form.time_b)])
    }
    if (this.form.search_options.length <= 0) {
      this.form.search_options = ['nginx', 'json', 'string']
    }
    if (this.form.without_options.length <= 0) {
      this.form.without_options = ['total', 'track', 'orange', 'ht', 'tg', 'source']
    }
  },
  methods: {
    async setTimeRange(v) {
      Vue.set(this.timerange, 0, v[0])
      Vue.set(this.timerange, 1, v[1])
    },
    async barClick(form) {
      const a = new Date(form.name)
      const b = new Date((+a) + this.charts.interval * 1000)
      await this.setTimeRange([a, b])
      this.form.time_a = +a
      this.form.time_b = +b
      await this.getList()
    },
    clipboardSuccess() { this.$message({ message: '复制成功', type: 'success', duration: 1500 }) },
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
      if (row.log_type === 'nginx') { // 格式化请求信息
        let url = row.level
        if (row.http_query) {
          url += '?' + row.http_query
        }
        const query = JSON.stringify(this.parseQueryString('?' + row.http_query), null, 2)
        const http_body = row.http_body
        if (http_body.indexOf('<xml') === 0) {
          return `URL:\n${url}\n\nquery to JSON:\n${query}\n\nXML body:\n${http_body}`
        }
        let body
        if (http_body.indexOf('{') === 0) {
          body = JSON.stringify(JSON.parse(http_body), null, 2)
        } else {
          body = JSON.stringify(this.parseQueryString('?' + row.http_body), null, 2)
        }
        return `URL:\n${url}\n\nquery to JSON:\n${query}\n\nbody to JSON:\n${body}`
      }
      if (row.detail && row.detail.indexOf('app_message') >= 0) { // 有些应用的自定义数据可以提前解开
        const detail = JSON.parse(row.detail)
        if (detail.app_message) {
          detail.app_message = JSON.parse(detail.app_message)
          return JSON.stringify(detail, null, 2)
        }
      }
      const detail = JSON.parse(row.detail)
      const detail_friendly = JSON.stringify(detail, null, 2)
      if (detail && detail.info && detail.info.stack) { // Go 堆栈
        return `${detail.info.stack || ''}\n\n${detail_friendly}`
      }
      if (detail && detail.info && detail.info.debug_backtrace) { // PHP5堆栈格式化1
        return `${detail.info.debug_backtrace || ''}\n\n${detail_friendly}`
      }
      if (detail && detail.info && detail.info.trace) { // PHP5堆栈格式化2
        return `${detail.info.message || ''}\n\n${detail.info.trace || ''}\n\n${detail_friendly}`
      }
      if (detail && detail.message && detail.trace) { // PHP7(lumen)堆栈格式化
        return `${detail.message || ''}\n\n${detail.trace || ''}\n\n${detail_friendly}`
      }
      return detail_friendly
    },
    jsonFormatter(v) { return jsonFormatter(JSON.stringify(v)) },
    changeChartsVisible() {
      this.charts_visible = !this.charts_visible
      if (this.charts_visible) this.getList()
    },
    timestampFormatter(row) {
      const ts = (row.time * 0.001).toFixed()
      const ms = row.time % 1000
      return date('Y-m-d H:i:s', ts) + '.' + ms
    },
    copy(event, row, type = '') {
      let s = ''
      switch (type) {
        case 'url': // [仅复制URL]
          s = row.level
          break
        case 'query': // [仅复制query]
          s = row.http_query
          break
        case 'body': // [仅复制body]
          s = row.http_body
          break
        case 'cURL-raw': // [复制线上cURL-原始]
          s = row.http_curl.replaceAll('#SCHEME#', row.http_scheme).replaceAll('#HOST#', row.http_host)
          break
        case 'cURL-http': // [复制线上cURL-http]
          s = row.http_curl.replaceAll('#SCHEME#', 'http').replaceAll('#HOST#', row.http_host)
          break
        case 'cURL-https': // [复制线上cURL-https]
          s = row.http_curl.replaceAll('#SCHEME#', 'https').replaceAll('#HOST#', row.http_host)
          break
        case 'cURL2-raw': // [复制指定机器cURL-原始]
          s = row.http_curl.replaceAll('#SCHEME#', row.http_scheme).replaceAll('#HOST#', this.replay_host)
          break
        case 'cURL2-http': // [复制指定机器cURL-http]
          s = row.http_curl.replaceAll('#SCHEME#', 'http').replaceAll('#HOST#', this.replay_host)
          break
        case 'cURL2-https': // [复制指定机器cURL-https]
          s = row.http_curl.replaceAll('#SCHEME#', 'https').replaceAll('#HOST#', this.replay_host)
          break
      }
      if (!s) {
        this.$message({ message: '无可复制的数据', type: 'success', duration: 1500 })
        return
      }
      clip(s, event)
      this.$message({ message: '复制成功', type: 'success', duration: 1500 })
    },
    async getReplayServers() {
      this.replay_servers = [
        { label: '腾讯云 - QA自动化回归 - 119.29.46.21', value: '119.29.46.21' },
        { label: '腾讯云 - 开发专用1 - 81.71.12.52', value: '81.71.12.52' },
        { label: '腾讯云 - 开发专用2 - 42.194.153.5', value: '42.194.153.5' }
      ]
      const treeIDs = [22724] // 开发/测试
      this.loading = true
      try {
        for (let i = 0; i < treeIDs.length; i++) {
          const treeID = treeIDs[i]
          const res = await cmdbServersRemote(treeID)
          if (res.code !== 0) {
            this.$message({ type: 'error', message: res.message })
          } else {
            res.data = res.data || []
            this.replay_servers = this.replay_servers.concat(res.data.map(item => { return { label: `${item.configuration.zone} - ${item.ip}`, value: item.ip } }))
          }
        }
      } finally {
        this.loading = false
      }
    },
    async getList(button = 1) {
      switch (button) {
        case 2:
          await this.setTimeRange([this.timerange[0], new Date()])
          break
        case 3:
          await this.setTimeRange(today())
          break
        case 4:
          await this.setTimeRange(yesterday())
          break
      }
      this.loading = true
      try {
        const query = {
          page_no: this.form.page_no, page_size: this.form.page_size,
          elastic_search: 'http://sy-logsystem-es.39on.com', index: this.index,
          or: this.or, must: this.must, must_not: this.must_not,
          charts_visible: this.charts_visible,
          time_a: this.form.time_a, time_b: this.form.time_b,
          without_options: this.form.without_options
        }
        const res = await loggerSearch(query)
        if (res.code !== 0) {
          this.$message({ type: 'error', message: res.message })
        } else {
          this.form.page_no = res.data.page_no
          this.form.page_size = res.data.page_size
          this.form.count = res.data.count
          this.list = res.data ? res.data.list : []
          if (this.list && this.list.length > 0) {
            this.charts = res.data ? res.data.charts : []
          }
        }
      } finally {
        this.loading = false
      }
    },
    async replay(row, scheme, host) {
      if (!scheme) scheme = row.http_scheme
      if (!host) host = row.http_host
      this.loading = true
      try {
        const res = await loggerReplay({ scheme, host, log_item: row })
        if (res.code !== 0) {
          this.$message({ type: 'error', message: res.message })
          return
        }
        const end = '\r\n'
        const d = res.data
        const t = date()
        let urlQuery = row.http_uri
        if (row.http_query) urlQuery += '?' + row.http_query
        let detail = ''
        detail += '------------------------------' + end
        detail += `回放请求...${end}机器: ${host}${end}协议: ${scheme}${end}时间: ${t}` + end
        detail += '------------------------------' + end
        detail += end
        detail += '>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>' + end
        detail += `${row.http_method} ${urlQuery}` + end
        detail += d.request_headers.join(end) + end + end
        if (row.http_body) detail += row.http_body + end
        detail += '>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>' + end
        detail += end
        detail += '<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<' + end
        detail += d.response_headers.join(end) + end + end
        if (d.response_body) detail += d.response_body + end
        detail += '<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<' + end

        if (d.response_body_formatted) {
          detail += end
          detail += '------------------------------' + end
          detail += d.response_body_formatted + end
          detail += '------------------------------' + end
          detail += end
        }

        this.detail = detail
        this.dialog_visible = true
      } finally {
        this.loading = false
      }
    },
    handleSizeChange(val) {
      this.form.page_size = val
      this.getList()
    },
    handleCurrentChange(val) {
      this.form.page_no = val
      this.getList()
    },
    handleDetail(detail) {
      if (!detail) return
      this.detail = detail
      this.dialog_visible = true
    },
    async exportExcel() {
      this.$notify({
        title: '帮助',
        message: '由于浏览器限制, 系统按 400M 分割日志, 可以在本地使用命令 `cat 日志导出-* > merge.csv` 合并日志文件'
      })
      const size = 1000
      const form = {
        elastic_search: 'http://sy-logsystem-es.39on.com', index: this.index, size,
        time_a: this.form.time_a,
        time_b: this.form.time_b,
        or: this.or, must: this.must, must_not: this.must_not,
        without_options: this.form.without_options
      }
      const exportCallback = data => {
        import('@/vendor/Export2Csv').then(csv => {
          setTimeout(_ => {
            this.$message({ message: `正在生成 CSV 文件，成功后将自动下载。`, type: 'success', duration: 1500 })
            const t = date('Ymd-His')
            const t2 = +new Date()
            const filename = `日志导出-${t}-${t2}`
            csv.export_string_to_csv({ data, filename })
            console.log(`导出完成, 用时 ${+new Date() - t2}ms`)
          }, 1000)
        })
      }
      const maxPage = 20000 // 2000万
      let logs = ''
      let page = 0
      let finish = false
      // const res = await loggerExport(form)
      do {
        page++
        const message = `正在拉取第 ${page} 页日志，每页 ${size} 条，最大支持 ${maxPage} 页。`
        console.log(message)
        this.$message({ message, type: 'success', duration: 1500 })
        this.loading = true
        try {
          const res = await loggerExport(form)
          const l = res.data.logs
          if (l) {
            logs += l // 合并
            if (logs.length >= 400 * 1000 * 1000) { // 400 MB (按照磁盘格式)
              exportCallback(logs) // 分块导出
              logs = ''
            }
          }
          // 后置处理
          const id = res.data.id || false
          if (!id) {
            finish = true
          } else {
            form.id = id
          }
        } finally {
          this.loading = false
        }
      } while (page < maxPage && !finish)
      if (logs.length > 0) {
        exportCallback(logs)
      }
    }
  }
}
</script>

<style lang="scss">
.logger-search {
  .el-form-item__label {
    font-size: 12px !important;
  }
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
  .el-form-item--mini.el-form-item {
    margin-bottom: 4px;
  }
  .el-input--mini .el-input__inner {
    padding: 0 10px;
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
.logger-search {
  padding: 10px 0;
}
.search-input {
  width: 160px;
}
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
