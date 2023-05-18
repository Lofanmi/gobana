<template>
  <div>
    <el-tabs v-model="activeName">
      <el-tab-pane label="快捷查询" name="快捷查询">
        <el-form :inline="true" label-width="150px">
          <el-form-item label="或者 [A || B || C]">
            <el-input v-show="1 <= length" v-model="form.or1" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="2 <= length" v-model="form.or2" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="3 <= length" v-model="form.or3" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="4 <= length" v-model="form.or4" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="5 <= length" v-model="form.or5" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="6 <= length" v-model="form.or6" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="7 <= length" v-model="form.or7" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="8 <= length" v-model="form.or8" class="search-input" @keyup.enter.native="query(1)" />
          </el-form-item>
          <el-form-item label="并且 [A && B && C]">
            <el-input v-show="1 <= length" v-model="form.must1" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="2 <= length" v-model="form.must2" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="3 <= length" v-model="form.must3" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="4 <= length" v-model="form.must4" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="5 <= length" v-model="form.must5" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="6 <= length" v-model="form.must6" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="7 <= length" v-model="form.must7" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="8 <= length" v-model="form.must8" class="search-input" @keyup.enter.native="query(1)" />
          </el-form-item>
          <el-form-item label="不包含 [!A && !B && !C]">
            <el-input v-show="1 <= length" v-model="form.must_not1" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="2 <= length" v-model="form.must_not2" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="3 <= length" v-model="form.must_not3" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="4 <= length" v-model="form.must_not4" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="5 <= length" v-model="form.must_not5" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="6 <= length" v-model="form.must_not6" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="7 <= length" v-model="form.must_not7" class="search-input" @keyup.enter.native="query(1)" />
            <el-input v-show="8 <= length" v-model="form.must_not8" class="search-input" @keyup.enter.native="query(1)" />
          </el-form-item>
        </el-form>
      </el-tab-pane>
      <el-tab-pane label="原生查询" name="second">
        <el-button>原生查询（TODO）</el-button>
      </el-tab-pane>
    </el-tabs>

    <el-form :inline="true" label-width="150px">
      <el-form-item label="后端">
        <el-select v-model="backend_name" default-first-option placeholder="请选择后端" class="search-input">
          <el-option v-for="item in backend_list" :key="item.value" :value="item.value" :label="item.label" />
        </el-select>
        <el-select v-model="storage_name" default-first-option placeholder="请选择存储" class="search-input">
          <el-option v-for="item in storage_list" :key="item.value" :value="item.value" :label="item.label" />
        </el-select>

        <el-date-picker v-model="timerange" type="datetimerange" start-placeholder="开始时间" end-placeholder="结束时间" :picker-options="picker_options" style="width: 324px" />
        <el-button-group>
          <el-button @click="query(1)">查询</el-button>
          <el-button @click="query(2)">最新</el-button>
        </el-button-group>
        <el-button-group>
          <el-button @click="query(3)">今天</el-button>
          <el-button @click="query(4)">昨天</el-button>
          <el-button @click="query(5)">前天</el-button>
        </el-button-group>

      </el-form-item>
    </el-form>

  </div>
</template>

<script>
import Vue from 'vue'
import { lasthour, today, yesterday, lastday, dbyesterday } from '@/utils/timeshortcut'
import { configGetBackendList, configGetStorageList } from '@/api/config'

export default {
  data() {
    return {
      screenWidth: 0,
      backend_name: '',
      backend_list: [],
      storage_name: '',
      storage_list: [],
      activeName: '快捷查询',
      timerange: [],
      form: {
        or1: '', or2: '', or3: '', or4: '', or5: '', or6: '', or7: '', or8: '',
        must1: '', must2: '', must3: '', must4: '', must5: '', must7: '', must6: '', must8: '',
        must_not1: '', must_not2: '', must_not3: '', must_not4: '', must_not5: '', must_not6: '', must_not7: '', must_not8: '',
        time_a: 0, time_b: 0
      },
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
      }
    }
  },
  computed: {
    query_by() {
      return this.activeName === '快捷查询' ? 'query_by_human' : 'query_by_query_string'
    },
    or() {
      return [this.form.or1, this.form.or2, this.form.or3, this.form.or4, this.form.or5, this.form.or6, this.form.or7, this.form.or8].filter(i => !!i)
    },
    must() {
      return [this.form.must1, this.form.must2, this.form.must3, this.form.must4, this.form.must5, this.form.must6, this.form.must7, this.form.must8].filter(i => !!i)
    },
    must_not() {
      return [this.form.must_not1, this.form.must_not2, this.form.must_not3, this.form.must_not4, this.form.must_not5, this.form.must_not6, this.form.must_not7, this.form.must_not8].filter(i => !!i)
    },
    length() {
      if (this.screenWidth >= 1600) return 8
      if (this.screenWidth >= 1400) return 6
      return 5
    }
  },
  watch: {
    timerange: {
      handler(v) {
        this.form.time_a = +(v[0])
        this.form.time_b = +(v[1])
      },
      deep: true
    },
    form: {
      handler() {
      },
      deep: true
    },
    async backend_name() {
      if (!this.backend_name) return
      const res = await configGetStorageList({ backend_name: this.backend_name })
      if (res.code === 0) {
        this.storage_list = res.data.storage_list
        if (this.storage_list.length > 0) {
          this.storage_name = this.storage_list[0].value
        }
      }
    }
  },
  mounted() {
    this.screenWidth = document.body.clientWidth
    window.onresize = () => {
      return (() => {
        this.screenWidth = document.body.clientWidth
      })()
    }
  },
  async created() {
    const res = await configGetBackendList()
    if (res.code === 0) {
      this.backend_list = res.data.backend_list
      if (this.backend_list.length > 0) {
        this.backend_name = this.backend_list[0].value
      }
    }
    if (this.form.time_a <= 0 && this.form.time_b <= 0) {
      await this.setTimeRange(lasthour(0.25))
    } else {
      await this.setTimeRange([new Date(this.form.time_a), new Date(this.form.time_b)])
    }
  },
  methods: {
    async setTimeRange(v) {
      Vue.set(this.timerange, 0, v[0])
      Vue.set(this.timerange, 1, v[1])
    },
    async query(button = 1) {
      switch (button) {
        case 2:
          await this.setTimeRange([this.timerange[0], new Date()])
          break
        case 3:
          await this.setTimeRange(today())
          return
        case 4:
          await this.setTimeRange(yesterday())
          return
        case 5:
          await this.setTimeRange(dbyesterday())
          return
      }
      this.$emit('query', {
        button: button,
        params: {
          time_a: this.form.time_a,
          time_b: this.form.time_b,
          backend: this.backend_name,
          storage: this.storage_name,
          query_by: this.query_by,
          query: {
            or: this.or,
            must: this.must,
            must_not: this.must_not
          }
        }
      })
    }
  }
}
</script>

<style lang="scss">
.logger-search {
  .el-form-item__label {
    font-size: 12px !important;
  }
  .el-form-item--mini.el-form-item {
    margin-bottom: 4px;
  }
  .el-input--mini .el-input__inner {
    padding: 0 10px;
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
</style>
