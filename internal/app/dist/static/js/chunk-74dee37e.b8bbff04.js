(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-74dee37e"],{"370b":function(e,t,n){},"5d95":function(e,t,n){"use strict";n("6222")},"616b":function(e,t,n){},6222:function(e,t,n){},6946:function(e,t,n){"use strict";n("936f")},"705a":function(e,t,n){"use strict";n.r(t);var r=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{staticClass:"app-container logger-search"},[n("query",{staticStyle:{width:"100%"},attrs:{loading:e.loading,charts:e.charts},on:{query:e.onQuery}}),e._v(" "),n("el-table",{directives:[{name:"loading",rawName:"v-loading",value:e.loading,expression:"loading"}],staticStyle:{width:"100%"},attrs:{data:e.list,"row-class-name":e.tableRowClassName,size:"mini","highlight-current-row":""}},[n("el-table-column",{attrs:{type:"expand"},scopedSlots:e._u([{key:"default",fn:function(t){return["access-log"===t.row.log_type?[n("kv-detail",{attrs:{label:"storage",value:t.row.storage,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"log_type",value:t.row.log_type,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"request_id",value:t.row.log.request_id,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"uri",value:t.row.log.uri,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"body",value:t.row.log.body,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"http_version",value:t.row.log.http_version,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"user_agent",value:t.row.log.user_agent,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"referer",value:t.row.log.referer,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"x_forwarded_for",value:t.row.log.x_forwarded_for,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"cookie",value:t.row.log.cookie,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"remote_addr",value:t.row.log.remote_addr,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"message",value:t.row.log.message,copy:!0}})]:"json-log"===t.row.log_type?[n("kv-detail",{attrs:{label:"storage",value:t.row.storage,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"log_type",value:t.row.log_type,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"request_id",value:t.row.log.request_id,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"path",value:t.row.log.path,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"message",value:t.row.log.message,copy:!0}})]:"string-log"===t.row.log_type?[n("kv-detail",{attrs:{label:"storage",value:t.row.storage,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"log_type",value:t.row.log_type,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"path",value:t.row.log.path,copy:!0}}),e._v(" "),n("kv-detail",{attrs:{label:"message",value:t.row.log.message,copy:!0}})]:[n("span",[e._v("未知的日志类型")])]]}}])}),e._v(" "),n("el-table-column",{attrs:{label:"type",width:"48"},scopedSlots:e._u([{key:"default",fn:function(t){return[n("div",{staticClass:"svg-icon"},["access-log"===t.row.log_type?[n("svg-icon",{attrs:{"icon-class":"access_log"}})]:"json-log"===t.row.log_type?[n("svg-icon",{attrs:{"icon-class":"json_log"}})]:"string-log"===t.row.log_type?[n("svg-icon",{attrs:{"icon-class":"string_log"}})]:e._e()],2)]}}])}),e._v(" "),n("el-table-column",{attrs:{prop:"time",label:"time",width:"160",formatter:e.timestampFormatter}}),e._v(" "),n("el-table-column",{attrs:{label:"application/tag/path",width:"150"},scopedSlots:e._u([{key:"default",fn:function(t){return["access-log"===t.row.log_type?[n("column",{attrs:{content:t.row.log.http_host}})]:"json-log"===t.row.log_type?[n("column",{attrs:{content:t.row.log.tag}})]:"string-log"===t.row.log_type?[n("column",{attrs:{content:t.row.log.path}})]:e._e()]}}])}),e._v(" "),n("el-table-column",{attrs:{prop:"hostname",label:"hostname",width:"180"},scopedSlots:e._u([{key:"default",fn:function(e){return[n("column",{attrs:{content:e.row.log.hostname,center:!1}})]}}])}),e._v(" "),n("el-table-column",{attrs:{prop:"level",label:"uri/level",width:"240"},scopedSlots:e._u([{key:"default",fn:function(t){return["access-log"===t.row.log_type?[n("column",{attrs:{content:t.row.log.uri,copy:!0,tooltip:!0,center:!1}})]:"json-log"===t.row.log_type?[n("column",{attrs:{content:t.row.log.level,copy:!1,tooltip:!1,center:!0}})]:"string-log"===t.row.log_type?[n("span",{staticClass:"one-line"})]:e._e()]}}])}),e._v(" "),n("el-table-column",{attrs:{prop:"message",label:"message"},scopedSlots:e._u([{key:"default",fn:function(t){return["access-log"===t.row.log_type?[n("column",{attrs:{tooltip:!1,content:e.accessLogMessageFormatter(t.row),center:!1}})]:"json-log"===t.row.log_type||"string-log"===t.row.log_type?[n("column",{attrs:{tooltip:!1,content:t.row.log.message,center:!1}})]:e._e()]}}])}),e._v(" "),n("el-table-column",{attrs:{prop:"request_id",label:"request_id",width:"100"},scopedSlots:e._u([{key:"default",fn:function(t){return["string-log"===t.row.log_type?n("span",{staticClass:"one-line"}):n("column",{attrs:{content:t.row.log.request_id,color:e.hashStringColor(t.row.log.request_id)}})]}}])}),e._v(" "),n("el-table-column",{attrs:{label:"操作",width:"100"},scopedSlots:e._u([{key:"default",fn:function(t){return[n("el-button",{attrs:{type:"text",size:"mini"},on:{click:function(n){e.handleDetail(e.logFormatter(t.row))}}},[e._v("格式化")]),e._v(" "),n("el-button",{attrs:{type:"text",size:"mini"},on:{click:function(n){e.handleDetail(e.jsonFormatter(t.row))}}},[e._v("原日志")])]}}])})],1),e._v(" "),n("div",{staticClass:"pagination"},[n("el-pagination",{attrs:{layout:"sizes, prev, pager, next, total, jumper","current-page":e.form.page_no,"page-sizes":[10,20,50,100,200],"page-size":e.form.page_size,total:e.form.count},on:{"size-change":e.handleSizeChange,"current-change":e.handleCurrentChange}})],1),e._v(" "),n("el-dialog",{attrs:{width:"80%",visible:e.dialog_visible},on:{"update:visible":function(t){e.dialog_visible=t}}},[e.detail?n("div",{staticClass:"code-group"},[n("div",{staticClass:"container"},[n("code",{staticClass:"code"},[e._v(e._s(e.detail))])])]):e._e(),e._v(" "),n("div",{staticStyle:{"text-align":"right"}},[n("el-button",{attrs:{type:"danger"},on:{click:function(t){e.dialog_visible=!1}}},[e._v("关闭")])],1)])],1)},a=[],s=(n("28a5"),n("96cf"),n("3b8d")),o=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("span",[e.copy?[e.tooltip&&"-"!==e.computedValue?n("el-tooltip",{attrs:{content:e.computedValue,placement:"top"}},[n("span",{directives:[{name:"clipboard",rawName:"v-clipboard:copy",value:e.computedValue,expression:"computedValue",arg:"copy"},{name:"clipboard",rawName:"v-clipboard:success",value:e.clipboardSuccess,expression:"clipboardSuccess",arg:"success"}],staticClass:"one-line copy-line",class:e.className,style:e.style,attrs:{title:"点击复制"}},[e._v("\n        "+e._s(e.computedValue)+"\n      ")])]):n("span",{directives:[{name:"clipboard",rawName:"v-clipboard:copy",value:e.computedValue,expression:"computedValue",arg:"copy"},{name:"clipboard",rawName:"v-clipboard:success",value:e.clipboardSuccess,expression:"clipboardSuccess",arg:"success"}],staticClass:"one-line copy-line",class:e.className,style:e.style,attrs:{title:"点击复制"}},[e._v("\n      "+e._s(e.computedValue)+"\n    ")])]:[e.tooltip&&"-"!==e.computedValue?n("el-tooltip",{attrs:{content:e.computedValue,placement:"top"}},[n("span",{staticClass:"one-line",class:e.className,style:e.style},[e._v("\n        "+e._s(e.computedValue)+"\n      ")])]):n("span",{staticClass:"one-line",class:e.className,style:e.style},[e._v("\n      "+e._s(e.computedValue)+"\n    ")])]],2)},i=[],l=n("b311");if(!l)throw new Error("you should npm install `clipboard` --save at first ");var u={bind:function(e,t){if("success"===t.arg)e._v_clipboard_success=t.value;else if("error"===t.arg)e._v_clipboard_error=t.value;else{var n=new l(e,{text:function(){return t.value},action:function(){return"cut"===t.arg?"cut":"copy"}});n.on("success",(function(t){var n=e._v_clipboard_success;n&&n(t)})),n.on("error",(function(t){var n=e._v_clipboard_error;n&&n(t)})),e._v_clipboard=n}},update:function(e,t){"success"===t.arg?e._v_clipboard_success=t.value:"error"===t.arg?e._v_clipboard_error=t.value:(e._v_clipboard.text=function(){return t.value},e._v_clipboard.action=function(){return"cut"===t.arg?"cut":"copy"})},unbind:function(e,t){"success"===t.arg?delete e._v_clipboard_success:"error"===t.arg?delete e._v_clipboard_error:(e._v_clipboard.destroy(),delete e._v_clipboard)}},c=function(e){e.directive("Clipboard",u)};window.Vue&&(window.clipboard=u,Vue.use(c)),u.install=c;var m=u,p={directives:{clipboard:m},props:{content:{type:String,default:""},color:{type:String,default:""},tooltip:{type:Boolean,default:!0},copy:{type:Boolean,default:!0},center:{type:Boolean,default:!0}},computed:{computedValue:function(){return this.content||"-"},className:function(){return this.center?"text-center":""},style:function(){return this.color?{color:this.color}:""}},methods:{clipboardSuccess:function(){this.$message({message:"复制成功",type:"success",duration:1500})}}},d=p,f=(n("5d95"),n("2877")),h=Object(f["a"])(d,o,i,!1,null,"5b9f3afd",null),v=h.exports,_=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{staticClass:"kv-detail"},[n("span",{staticClass:"label"},[e._v(e._s(e.label)+":")]),e._v(" "),e.copy&&"-"!==e.computedValue?n("span",{directives:[{name:"clipboard",rawName:"v-clipboard:copy",value:e.computedValue,expression:"computedValue",arg:"copy"},{name:"clipboard",rawName:"v-clipboard:success",value:e.clipboardSuccess,expression:"clipboardSuccess",arg:"success"}],staticClass:"value copy-line",attrs:{title:"点击复制"}},[e._v("\n    "+e._s(e.computedValue)+"\n  ")]):n("span",{staticClass:"value"},[e._v(e._s(e.computedValue))])])},g=[],y={directives:{clipboard:m},props:{label:{type:String,default:""},value:{type:String,default:""},copy:{type:Boolean,default:!1}},computed:{computedValue:function(){return this.value||"-"}},methods:{clipboardSuccess:function(){this.$message({message:"复制成功",type:"success",duration:1500})}}},b=y,k=(n("d09b"),Object(f["a"])(b,_,g,!1,null,"040dca40",null)),w=k.exports,x=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{staticStyle:{width:"100%"}},[n("el-tabs",{model:{value:e.query_by,callback:function(t){e.query_by=t},expression:"query_by"}},[n("el-tab-pane",{attrs:{label:"快捷查询",name:"query_by_human"}},[n("el-form",{attrs:{inline:!0,"label-width":"150px"},nativeOn:{submit:function(e){e.preventDefault()}}},[n("el-form-item",{attrs:{label:"或者 [A || B || C]"}},[n("el-input",{directives:[{name:"show",rawName:"v-show",value:1<=e.length,expression:"1 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.or1,callback:function(t){e.$set(e.form,"or1",t)},expression:"form.or1"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:2<=e.length,expression:"2 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.or2,callback:function(t){e.$set(e.form,"or2",t)},expression:"form.or2"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:3<=e.length,expression:"3 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.or3,callback:function(t){e.$set(e.form,"or3",t)},expression:"form.or3"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:4<=e.length,expression:"4 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.or4,callback:function(t){e.$set(e.form,"or4",t)},expression:"form.or4"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:5<=e.length,expression:"5 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.or5,callback:function(t){e.$set(e.form,"or5",t)},expression:"form.or5"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:6<=e.length,expression:"6 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.or6,callback:function(t){e.$set(e.form,"or6",t)},expression:"form.or6"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:7<=e.length,expression:"7 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.or7,callback:function(t){e.$set(e.form,"or7",t)},expression:"form.or7"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:8<=e.length,expression:"8 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.or8,callback:function(t){e.$set(e.form,"or8",t)},expression:"form.or8"}})],1),e._v(" "),n("el-form-item",{attrs:{label:"并且 [A && B && C]"}},[n("el-input",{directives:[{name:"show",rawName:"v-show",value:1<=e.length,expression:"1 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.must1,callback:function(t){e.$set(e.form,"must1",t)},expression:"form.must1"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:2<=e.length,expression:"2 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.must2,callback:function(t){e.$set(e.form,"must2",t)},expression:"form.must2"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:3<=e.length,expression:"3 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.must3,callback:function(t){e.$set(e.form,"must3",t)},expression:"form.must3"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:4<=e.length,expression:"4 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.must4,callback:function(t){e.$set(e.form,"must4",t)},expression:"form.must4"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:5<=e.length,expression:"5 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.must5,callback:function(t){e.$set(e.form,"must5",t)},expression:"form.must5"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:6<=e.length,expression:"6 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.must6,callback:function(t){e.$set(e.form,"must6",t)},expression:"form.must6"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:7<=e.length,expression:"7 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.must7,callback:function(t){e.$set(e.form,"must7",t)},expression:"form.must7"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:8<=e.length,expression:"8 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.must8,callback:function(t){e.$set(e.form,"must8",t)},expression:"form.must8"}})],1),e._v(" "),n("el-form-item",{attrs:{label:"不包含 [!A && !B && !C]"}},[n("el-input",{directives:[{name:"show",rawName:"v-show",value:1<=e.length,expression:"1 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.must_not1,callback:function(t){e.$set(e.form,"must_not1",t)},expression:"form.must_not1"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:2<=e.length,expression:"2 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.must_not2,callback:function(t){e.$set(e.form,"must_not2",t)},expression:"form.must_not2"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:3<=e.length,expression:"3 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.must_not3,callback:function(t){e.$set(e.form,"must_not3",t)},expression:"form.must_not3"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:4<=e.length,expression:"4 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.must_not4,callback:function(t){e.$set(e.form,"must_not4",t)},expression:"form.must_not4"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:5<=e.length,expression:"5 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.must_not5,callback:function(t){e.$set(e.form,"must_not5",t)},expression:"form.must_not5"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:6<=e.length,expression:"6 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.must_not6,callback:function(t){e.$set(e.form,"must_not6",t)},expression:"form.must_not6"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:7<=e.length,expression:"7 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.must_not7,callback:function(t){e.$set(e.form,"must_not7",t)},expression:"form.must_not7"}}),e._v(" "),n("el-input",{directives:[{name:"show",rawName:"v-show",value:8<=e.length,expression:"8 <= length"}],staticClass:"search-input",nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.must_not8,callback:function(t){e.$set(e.form,"must_not8",t)},expression:"form.must_not8"}})],1)],1)],1),e._v(" "),n("el-tab-pane",{attrs:{label:"Lucene",name:"query_by_lucene"}},[n("div",{staticClass:"lucene-help"},[n("el-alert",{attrs:{title:"Lucene 语法助手",type:"info",description:e.helper_lucene,"show-icon":""}})],1),e._v(" "),n("el-form",{attrs:{inline:!1,"label-width":"150px"},nativeOn:{submit:function(e){e.preventDefault()}}},[n("el-form-item",{attrs:{label:"Lucene"}},[n("el-input",{nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.query(1)}},model:{value:e.form.lucene,callback:function(t){e.$set(e.form,"lucene",t)},expression:"form.lucene"}})],1)],1)],1)],1),e._v(" "),n("el-form",{attrs:{inline:!1,"label-width":"150px"},nativeOn:{submit:function(e){e.preventDefault()}}},[n("el-form-item",{attrs:{label:"后端"}},[n("el-select",{staticClass:"search-input",attrs:{"default-first-option":"",placeholder:"请选择后端"},model:{value:e.backend_name,callback:function(t){e.backend_name=t},expression:"backend_name"}},e._l(e.backend_list,(function(e){return n("el-option",{key:e.value,attrs:{value:e.value,label:e.label}})})),1),e._v(" "),n("el-select",{staticClass:"search-input",attrs:{"default-first-option":"",placeholder:"请选择存储"},model:{value:e.storage_name,callback:function(t){e.storage_name=t},expression:"storage_name"}},e._l(e.storage_list,(function(e){return n("el-option",{key:e.value,attrs:{value:e.value,label:e.label}})})),1),e._v(" "),n("el-date-picker",{staticStyle:{width:"324px"},attrs:{type:"datetimerange","start-placeholder":"开始时间","end-placeholder":"结束时间","picker-options":e.picker_options},model:{value:e.timerange,callback:function(t){e.timerange=t},expression:"timerange"}}),e._v(" "),n("el-button-group",[n("el-button",{on:{click:function(t){return e.query(1)}}},[e._v("查询")]),e._v(" "),n("el-button",{on:{click:function(t){return e.query(2)}}},[e._v("最新")])],1),e._v(" "),n("el-button-group",[n("el-button",{on:{click:function(t){return e.query(3)}}},[e._v("今天")]),e._v(" "),n("el-button",{on:{click:function(t){return e.query(4)}}},[e._v("昨天")]),e._v(" "),n("el-button",{on:{click:function(t){return e.query(5)}}},[e._v("前天")])],1)],1),e._v(" "),n("el-form-item",{attrs:{label:"选项"}},[n("el-checkbox",{attrs:{size:"mini"},on:{change:e.onChartVisibleChange},model:{value:e.form.chart_visible,callback:function(t){e.$set(e.form,"chart_visible",t)},expression:"form.chart_visible"}},[e._v("日志柱状图")]),e._v(" "),n("el-checkbox",{attrs:{size:"mini"},model:{value:e.form.track_total_hits,callback:function(t){e.$set(e.form,"track_total_hits",t)},expression:"form.track_total_hits"}},[e._v("统计日志总数")])],1)],1),e._v(" "),n("div",{directives:[{name:"show",rawName:"v-show",value:e.showCharts,expression:"showCharts"}],staticClass:"chart",staticStyle:{"background-color":"#ffffff",width:"100%"}},[n("counter",{directives:[{name:"loading",rawName:"v-loading",value:e.loading,expression:"loading"}],ref:"counter",staticClass:"line",attrs:{id:"line",height:"160px",width:"100%",chartdata:e.charts},on:{"click-event":e.barClick}})],1)],1)},C=[],O=(n("7f7f"),n("2b0e"));n("a481"),n("456d"),n("ac6a"),n("7618");function q(e,t,n){var r,a,s,o,i,l=function l(){var u=+new Date-o;u<t&&u>0?r=setTimeout(l,t-u):(r=null,n||(i=e.apply(s,a),r||(s=a=null)))};return function(){for(var a=arguments.length,u=new Array(a),c=0;c<a;c++)u[c]=arguments[c];s=this,o=+new Date;var m=n&&!r;return r||(r=setTimeout(l,t)),m&&(i=e.apply(s,u),s=u=null),i}}function N(e,t){e=e||"Y-m-d H:i:s";var n=t||+new Date/1e3,r=["Mon","Tue","Wed","Thu","Fri","Sat","Sun"];n=new Date(1e3*(n>>>0));var a={year:n.getYear(),month:n.getMonth()+1,date:n.getDate(),day:r[n.getDay()],hours:n.getHours(),minutes:n.getMinutes(),seconds:n.getSeconds()};a.g=a.hours>12?Math.ceil(a.hours/2):a.hours;var s={Y:n.getFullYear(),y:a.year,m:a.month<10?"0"+a.month:a.month,n:a.month,d:a.date<10?"0"+a.date:a.date,j:a.date,D:a.day,H:a.hours<10?"0"+a.hours:a.hours,h:a.g<10?"0"+a.g:a.g,G:a.hours,g:a.g,i:a.minutes<10?"0"+a.minutes:a.minutes,s:a.seconds<10?"0"+a.seconds:a.seconds};for(var o in s)e=(""+e).replace(o,s[o]);return e}function $(e,t){var n=arguments.length>2&&void 0!==arguments[2]?arguments[2]:"0",r=String(e),a=(Array(t).join(n)+e).slice(-t);return r.length>a&&(a=r),a}function E(e){e=e||"";for(var t=0,n=e.length,r=5381;t<n;++t)r+=(r<<5)+e.charAt(t).charCodeAt();return 2147483647&r}function S(e,t){var n=60,r=parseInt(.001*(t-e)/n);return r=r<=1?1:r<=5?5:r<=10?10:r<=30?30:r<=60?60:r<=300?300:r<=900?900:r<=1800?1800:r<=3600?3600:r<=10800?10800:r<=32400?32400:r<=43200?43200:86400,r}var D=n("768b");function j(){var e=new Date,t=new Date;return e.setTime(1e3*Math.ceil(+new Date/1e3/3600/24)*3600*24-288e5),t.setTime(e.getTime()-864e5),e.setTime(e.getTime()-1e3),[t,e]}function z(){var e=j(),t=Object(D["a"])(e,2),n=t[0],r=t[1],a=new Date(n-864e5),s=new Date(r-864e5);return[a,s]}function R(){var e=z(),t=Object(D["a"])(e,2),n=t[0],r=t[1],a=new Date(n-864e5),s=new Date(r-864e5);return[a,s]}function V(e){var t=new Date,n=+t;return t.setTime(n-864e5*e),[t,new Date(n)]}function A(e){var t=new Date,n=+t;return t.setTime(n-36e5*e),[t,new Date(n)]}var T=n("bc3a"),F=n.n(T),B=n("5c96"),L=F.a.create({baseURL:Object({NODE_ENV:"production",BASE_URL:"/frontend/dist/"}).VUE_APP_BASE_API,withCredentials:!0,timeout:6e4});L.interceptors.request.use((function(e){return e}),(function(e){return console.log(e),Promise.reject(e)})),L.interceptors.response.use((function(e){var t=e.data;if(0===t.code)return t;Object(B["Message"])({message:t.message||"error",type:"error",duration:5e3})}),(function(e){return console.log(e),Object(B["Message"])({message:e.message,type:"error",duration:5e3}),Promise.reject(e)}));var Q=L,J="/api/gobana/v1";function H(e){return Q({url:J+"/config/backend_list",method:"get",params:e})}function M(e){return Q({url:J+"/config/storage_list",method:"get",params:e})}var W=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{class:e.className,style:{height:e.height,width:e.width},attrs:{id:e.id}})},P=[],I=n("3eba"),U=n.n(I),Y=(n("ef97"),n("94b1"),n("627c"),n("007d"),n("d28f"),{data:function(){return{$_sidebarElm:null}},mounted:function(){var e=this;this.__resizeHandler=q((function(){e.chart&&e.chart.resize()}),100),window.addEventListener("resize",this.__resizeHandler),this.$_sidebarElm=document.getElementsByClassName("sidebar-container")[0],this.$_sidebarElm&&this.$_sidebarElm.addEventListener("transitionend",this.$_sidebarResizeHandler)},beforeDestroy:function(){window.removeEventListener("resize",this.__resizeHandler),this.$_sidebarElm&&this.$_sidebarElm.removeEventListener("transitionend",this.$_sidebarResizeHandler)},methods:{$_sidebarResizeHandler:function(e){"width"===e.propertyName&&this.__resizeHandler()}}}),G={mixins:[Y],props:{className:{type:String,default:"chart"},id:{type:String,default:"chart"},width:{type:String,default:"90%"},height:{type:String,default:"400px"},chartdata:{type:Object,default:function(){}}},data:function(){return{chart:null,colors:{default:["#c23531","#2f4554","#61a0a8","#d48265","#91c7ae","#749f83","#ca8622","#bda29a","#6e7074","#546570","#c4ccd3"],light:["#37A2DA","#32C5E9","#67E0E3","#9FE6B8","#FFDB5C","#ff9f7f","#fb7293","#E062AE","#E690D1","#e7bcf3","#9d96f5","#8378EA","#96BFFF"],dark:["#dd6b66","#759aa0","#e69d87","#8dc1a9","#ea7e53","#eedd78","#73a373","#73b9bc","#7289ab","#91ca8c","#f49f42"]}}},computed:{xAxis:function(){return this.chartdata.xAxis},legend:function(){return this.chartdata.legend},series:function(){return this.chartdata.series}},watch:{chartdata:function(){this.setOption()}},mounted:function(){this.initChart()},beforeDestroy:function(){this.chart&&(this.chart.dispose(),this.chart=null)},methods:{initChart:function(){var e=this;this.chart=U.a.init(document.getElementById(this.id)),this.chart.on("click",(function(t){e.$emit("click-event",{seriesName:t.seriesName,name:t.name,value:t.value})}))},setOption:function(){this.chart.clear(),this.chart.setOption({title:{orient:"vertical",x:"center",text:""},tooltip:{trigger:"axis"},grid:{containLabel:!0,x:96,y:20,y2:20},color:"#9FE6B8",xAxis:{type:"category",boundaryGap:!0,data:this.xAxis},yAxis:{type:"value",min:0,logBase:10},series:this.series})},resizeChart:function(){this.chart&&this.chart.resize()}}},K=G,X=Object(f["a"])(K,W,P,!1,null,null,null),Z=X.exports,ee={components:{Counter:Z},props:{loading:{type:Boolean,default:!1},charts:{type:Object,default:null}},data:function(){return{screenWidth:0,backend_name:"",backend_list:[],storage_name:"",storage_list:[],query_by:"query_by_human",timerange:[],form:{or1:"",or2:"",or3:"",or4:"",or5:"",or6:"",or7:"",or8:"",must1:"",must2:"",must3:"",must4:"",must5:"",must7:"",must6:"",must8:"",must_not1:"",must_not2:"",must_not3:"",must_not4:"",must_not5:"",must_not6:"",must_not7:"",must_not8:"",time_a:0,time_b:0,lucene:"",chart_interval:0,chart_visible:!0,track_total_hits:!0},picker_options:{shortcuts:[{text:"近15分钟",onClick:function(e){e.$emit("pick",A(.25))}},{text:"近30分钟",onClick:function(e){e.$emit("pick",A(.5))}},{text:"近1小时",onClick:function(e){e.$emit("pick",A(1))}},{text:"近3小时",onClick:function(e){e.$emit("pick",A(3))}},{text:"今天",onClick:function(e){e.$emit("pick",j())}},{text:"昨天",onClick:function(e){e.$emit("pick",z())}},{text:"近1天",onClick:function(e){e.$emit("pick",V(1))}},{text:"近2天",onClick:function(e){e.$emit("pick",V(2))}},{text:"近1周",onClick:function(e){e.$emit("pick",V(7))}}]},helper_lucene:"message:ok AND grade:(60,80] AND NOT error"}},computed:{or:function(){return[this.form.or1,this.form.or2,this.form.or3,this.form.or4,this.form.or5,this.form.or6,this.form.or7,this.form.or8].filter((function(e){return!!e}))},must:function(){return[this.form.must1,this.form.must2,this.form.must3,this.form.must4,this.form.must5,this.form.must6,this.form.must7,this.form.must8].filter((function(e){return!!e}))},must_not:function(){return[this.form.must_not1,this.form.must_not2,this.form.must_not3,this.form.must_not4,this.form.must_not5,this.form.must_not6,this.form.must_not7,this.form.must_not8].filter((function(e){return!!e}))},queryParams:function(){return"query_by_human"===this.query_by?{or:this.or,must:this.must,must_not:this.must_not}:"query_by_lucene"===this.query_by?{lucene:this.form.lucene}:{}},length:function(){return this.screenWidth>=1600?8:this.screenWidth>=1400?6:5},showCharts:function(){return this.form.chart_visible&&this.charts.series&&this.charts.series.data&&this.charts.series.data.length>0}},watch:{timerange:{handler:function(e){this.form.time_a=+e[0],this.form.time_b=+e[1]},deep:!0},form:{handler:function(){},deep:!0},charts:function(){var e=this;this.$nextTick((function(){e.$refs.counter.resizeChart()}))},backend_name:function(){var e=Object(s["a"])(regeneratorRuntime.mark((function e(){var t;return regeneratorRuntime.wrap((function(e){while(1)switch(e.prev=e.next){case 0:if(this.backend_name){e.next=2;break}return e.abrupt("return");case 2:return e.next=4,M({backend_name:this.backend_name});case 4:t=e.sent,0===t.code&&(this.storage_list=t.data.storage_list,this.storage_list.length>0&&(this.storage_name=this.storage_list[0].value));case 6:case"end":return e.stop()}}),e,this)})));function t(){return e.apply(this,arguments)}return t}()},mounted:function(){var e=this;this.screenWidth=document.body.clientWidth,window.onresize=function(){return function(){e.screenWidth=document.body.clientWidth}()}},created:function(){var e=Object(s["a"])(regeneratorRuntime.mark((function e(){var t;return regeneratorRuntime.wrap((function(e){while(1)switch(e.prev=e.next){case 0:return e.next=2,H();case 2:if(t=e.sent,0===t.code&&(this.backend_list=t.data.backend_list,this.backend_list.length>0&&(this.backend_name=this.backend_list[0].value)),!(this.form.time_a<=0&&this.form.time_b<=0)){e.next=9;break}return e.next=7,this.setTimeRange(A(.25));case 7:e.next=11;break;case 9:return e.next=11,this.setTimeRange([new Date(this.form.time_a),new Date(this.form.time_b)]);case 11:case"end":return e.stop()}}),e,this)})));function t(){return e.apply(this,arguments)}return t}(),methods:{setTimeRange:function(){var e=Object(s["a"])(regeneratorRuntime.mark((function e(t){return regeneratorRuntime.wrap((function(e){while(1)switch(e.prev=e.next){case 0:O["default"].set(this.timerange,0,t[0]),O["default"].set(this.timerange,1,t[1]);case 2:case"end":return e.stop()}}),e,this)})));function t(t){return e.apply(this,arguments)}return t}(),query:function(){var e=Object(s["a"])(regeneratorRuntime.mark((function e(){var t,n=arguments;return regeneratorRuntime.wrap((function(e){while(1)switch(e.prev=e.next){case 0:t=n.length>0&&void 0!==n[0]?n[0]:1,e.t0=t,e.next=2===e.t0?4:3===e.t0?7:4===e.t0?10:5===e.t0?13:16;break;case 4:return e.next=6,this.setTimeRange([this.timerange[0],new Date]);case 6:return e.abrupt("break",16);case 7:return e.next=9,this.setTimeRange(j());case 9:return e.abrupt("return");case 10:return e.next=12,this.setTimeRange(z());case 12:return e.abrupt("return");case 13:return e.next=15,this.setTimeRange(R());case 15:return e.abrupt("return");case 16:this.$emit("query",{button:t,params:{time_a:this.form.time_a,time_b:this.form.time_b,backend:this.backend_name,storage:this.storage_name,query_by:this.query_by,query:this.queryParams,chart_interval:S(this.form.time_a,this.form.time_b),chart_visible:this.form.chart_visible,track_total_hits:this.form.track_total_hits}});case 17:case"end":return e.stop()}}),e,this)})));function t(){return e.apply(this,arguments)}return t}(),barClick:function(){var e=Object(s["a"])(regeneratorRuntime.mark((function e(t){var n,r;return regeneratorRuntime.wrap((function(e){while(1)switch(e.prev=e.next){case 0:return n=new Date(t.name),r=new Date(+n+1e3*this.charts.interval),e.next=4,this.setTimeRange([n,r]);case 4:return this.form.time_a=+n,this.form.time_b=+r,e.next=8,this.query();case 8:case"end":return e.stop()}}),e,this)})));function t(t){return e.apply(this,arguments)}return t}(),onChartVisibleChange:function(){var e=this;this.$nextTick((function(){e.$refs.counter.resizeChart()}))}}},te=ee,ne=(n("9e32"),n("a743"),Object(f["a"])(te,x,C,!1,null,"2cc34331",null)),re=ne.exports,ae=n("8975"),se="/api/gobana/v1";function oe(e){return Q({url:se+"/logger/search",method:"post",data:e})}var ie={components:{Column:v,KvDetail:w,Query:re},directives:{clipboard:m},data:function(){return{loading:!1,list:[],total:0,form:{page_no:1,page_size:10,count:0},lastQuery:null,detail:null,dialog_visible:!1,charts:{}}},watch:{},created:function(){var e=Object(s["a"])(regeneratorRuntime.mark((function e(){return regeneratorRuntime.wrap((function(e){while(1)switch(e.prev=e.next){case 0:case"end":return e.stop()}}),e)})));function t(){return e.apply(this,arguments)}return t}(),methods:{parseQueryString:function(e){var t={},n=e.substring(e.indexOf("?")+1,e.length).split("&");for(var r in n){var a=n[r].split("=");a[0]&&(t[a[0]]=decodeURIComponent(a[1]))}return t},logFormatter:function(e){if("string-log"===e.log_type)return e.log.message;if("access-log"===e.log_type){var t,n=e.log.uri,r=JSON.stringify(this.parseQueryString("?"+e.log.query),null,2),a=e.log.body;return 0===a.indexOf("<xml")?"URI:\n".concat(n,"\n\nquery to JSON:\n").concat(r,"\n\nXML body:\n").concat(a):(t=0===a.indexOf("{")?JSON.stringify(JSON.parse(a),null,2):JSON.stringify(this.parseQueryString("?"+e.log.body),null,2),"URI:\n".concat(n,"\n\nquery to JSON:\n").concat(r,"\n\nbody to JSON:\n").concat(t,"\n\n").concat(e.log.message))}var s=JSON.parse(e.log.message);return JSON.stringify(s,null,2)},jsonFormatter:function(e){return Object(ae["jsonFormatter"])(JSON.stringify(e))},timestampFormatter:function(e){var t=+new Date(e.log.time),n=(.001*t).toFixed(),r=$(t%1e3,3,"0");return N("Y-m-d H:i:s",n)+"."+r},accessLogMessageFormatter:function(e){var t=e.log;return"".concat(t.method," [").concat(t.remote_addr,"] status[").concat(t.status,"] duration[").concat(t.duration,"]")},hashStringColor:function(e){var t=["#2a9d2a","#645b93","#ff7f6a","#5fc0ea","#480048","#601848","#c04848","#f07241","#c71585","#008b8b","#7b68ee","#ff7f50"],n=E(e)%(t.length-1);return t[n]},tableRowClassName:function(e){var t=e.row;if(console.log(t.log_type),"access-log"===t.log_type)return t.log.status>=500?"error-row":t.log.status>=400?"warning-row":t.log.status>=300?"info-row":"success-row";if("json-log"===t.log_type)return"fatal"===t.log.level||"error"===t.log.level?"error-row":"warning"===t.log.level||"warn"===t.log.level?"warning-row":"info"===t.log.level?"info-row":"success-row";if("string-log"===t.log_type){var n=t.log.message.toLowerCase();return n.indexOf("fatal")>=0||n.indexOf("error")>=0?"error-row":n.indexOf("warning")>=0||n.indexOf("warn")>=0?"warning-row":n.indexOf("info")>=0?"info-row":"success-row"}return""},onQuery:function(){var e=Object(s["a"])(regeneratorRuntime.mark((function e(t){var n,r;return regeneratorRuntime.wrap((function(e){while(1)switch(e.prev=e.next){case 0:return this.lastQuery=t,this.loading=!0,e.prev=2,n=t.params,n.page_no=this.form.page_no,n.page_size=this.form.page_size,e.next=8,oe(n);case 8:r=e.sent,0!==r.code?this.$message({type:"error",message:r.message}):(this.form.page_no=r.data.page_no,this.form.page_size=r.data.page_size,this.form.count=r.data.count,this.list=r.data?r.data.list:[],this.charts=r.data?r.data.charts:[]);case 10:return e.prev=10,this.loading=!1,e.finish(10);case 13:case"end":return e.stop()}}),e,this,[[2,,10,13]])})));function t(t){return e.apply(this,arguments)}return t}(),handleSizeChange:function(e){this.form.page_size=e,this.lastQuery&&this.onQuery(this.lastQuery)},handleCurrentChange:function(e){this.form.page_no=e,this.lastQuery&&this.onQuery(this.lastQuery)},handleDetail:function(e){e&&(this.detail=e,this.dialog_visible=!0)}}},le=ie,ue=(n("6946"),n("b18c"),Object(f["a"])(le,r,a,!1,null,"a66ef04a",null));t["default"]=ue.exports},"7e38":function(e,t,n){},"861c":function(e,t,n){},"936f":function(e,t,n){},"9e32":function(e,t,n){"use strict";n("370b")},a743:function(e,t,n){"use strict";n("616b")},b18c:function(e,t,n){"use strict";n("7e38")},d09b:function(e,t,n){"use strict";n("861c")}}]);
//# sourceMappingURL=chunk-74dee37e.b8bbff04.js.map