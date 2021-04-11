<template>
    <div class="home">
        <div class="jumbotron bg-transparent mb-0 radius-0">
            <b-container fluid="">
                <b-row>
                    <b-col xl="3" lg="3" md="2" sm="0"></b-col>
                    <b-col xl="6" lg="6" md="8" sm="12">
                        <h1 class="display-2">Paste Bin<span class="vim-caret">x</span></h1>
                        <div class="lead mb-3 text-mono text-success">Intended for use as a short-term exchange of
                            pasted information between parties.
                        </div>
                        <div class="text-mono">
                            <a href="#!"
                               title="Download Theme"
                               class="btn btn-success btn-shadow px-3 my-2 ml-0 text-left">
                                Start
                            </a>
                        </div>

                        <div class="text-darkgrey text-mono my-2">It is intended to be used directly by humans.</div>
                    </b-col>
                    <b-col xl="3" lg="3" md="2" sm="0"></b-col>
                </b-row>
            </b-container>
        </div>

        <b-container fluid="">
            <b-row>
                <b-col xl="3" lg="3" md="2" sm="0"></b-col>
                <b-col xl="6" lg="6" md="8" sm="12">
                    <codemirror
                        ref="cmEditor"
                        v-model="record.content"
                        :options="cmOptions"
                        class="input"
                        placeholder="Paste code here..."
                    />

                    <b-container fluid="" class="mt-2">
                        <b-row>
                            <b-col sm="12" md="12" lg="6" xl="6">
                                <b-form-group
                                    label="Title"
                                    label-for="input-title"
                                    label-align="left"
                                >
                                    <b-form-input id="input-title" v-model="record.title" trim
                                                  placeholder="Brief Title Description"></b-form-input>
                                </b-form-group>
                            </b-col>
                            <b-col sm="12" md="12" lg="6" xl="6">
                                <b-form-group
                                    label="Author / Editable"
                                    label-for="input-author"
                                    label-align="left"
                                >

                                    <b-input-group>
                                        <b-form-input id="input-title" v-model="record.author" trim
                                                      placeholder="Brief description"></b-form-input>
                                        <b-input-group-prepend is-text><b>Editable</b></b-input-group-prepend>
                                        <b-input-group-prepend is-text>
                                            <b-form-checkbox switch class="mr-n2" v-model="record.editable"
                                                             button-variant="primary">
                                                <span class="sr-only">Switch for following text input</span>
                                            </b-form-checkbox>
                                        </b-input-group-prepend>
                                    </b-input-group>
                                </b-form-group>
                            </b-col>
                        </b-row>
                        <b-row>
                            <b-col sm="12" md="12" lg="6" xl="6">
                                <b-form-group
                                    label="Language"
                                    label-for="input-lang"
                                    label-align="left"
                                >
                                    <b-form-select v-model="record.lang" :options="lang_opts"></b-form-select>
                                </b-form-group>
                            </b-col>
                            <b-col sm="12" md="12" lg="6" xl="6">
                                <b-form-group
                                    label="Tags"
                                    label-for="input-lang"
                                    label-align="left"
                                >
                                    <b-form-tags
                                        id="custom_tags"
                                        v-model="record.tags"
                                        tag-variant="primary"
                                        :limit="tags_limit"
                                        tag-pills
                                        separator=" "
                                        placeholder="Enter new tags separated by space"
                                    ></b-form-tags>
                                </b-form-group>
                            </b-col>
                        </b-row>
                        <b-row>
                            <b-col sm="12" md="12" lg="6" xl="6">
                                <b-form-group
                                    label="Expiration"
                                    label-for="input-exp"
                                    label-align="left"
                                    description="Submitted data is not guaranteed to be permanent."
                                >
                                    <b-form-select v-model="record.lifecycle" :options="lifecycle_opts"></b-form-select>
                                </b-form-group>
                            </b-col>
                            <b-col sm="12" md="12" lg="6" xl="6">
                                <b-form-group
                                    label="Password"
                                    label-for="input-password"
                                    label-align="left"
                                >
                                    <b-input-group>
                                        <b-input-group-prepend is-text>
                                            <b-form-checkbox switch class="mr-n2" v-model="need_password"
                                                             button-variant="primary">
                                                <span class="sr-only">Switch for following text input</span>
                                            </b-form-checkbox>
                                        </b-input-group-prepend>
                                        <b-form-input v-model="record.password" :disabled="!need_password"
                                                      placeholder="Enter a password"></b-form-input>
                                    </b-input-group>
                                </b-form-group>
                            </b-col>
                        </b-row>
                    </b-container>

                    <b-alert class="mt-2 mb-2 text-left"
                             :show="pasteErrorShow"
                             dismissible
                             variant="warning"
                             @dismissed="dismissCountDown=0"
                             @dismiss-count-down="countDownChanged"
                    >
                        Paste error: {{ pasteErrorMessage }}
                    </b-alert>

                </b-col>
                <b-col xl="3" lg="3" md="2" sm="0"></b-col>
            </b-row>

            <b-button variant="primary" class="mt-2" @click="paste">Paste!</b-button>
        </b-container>

        <div class="container py-5">
            <p>Powered by <a href="https://xuthus.cc" target="_blank" class="text-danger">xuthus</a>, it's open source
                :) - <a
                    href="https://github.com/xuthus5/pbx" target="_blank" class="text-danger">pbx</a></p>
        </div>
    </div>
</template>

<script>
import 'codemirror/mode/javascript/javascript.js'
// import 'codemirror/mode/go/go.js'

// closebrackets
import 'codemirror/addon/edit/closebrackets.js'

import 'codemirror/addon/display/placeholder.js'

// styleSelectedText
import 'codemirror/addon/selection/mark-selection.js'

// keyMap
import 'codemirror/mode/clike/clike.js'
import 'codemirror/addon/edit/matchbrackets.js'
import 'codemirror/addon/comment/comment.js'
import 'codemirror/keymap/sublime.js'

// foldGutter
import 'codemirror/addon/fold/foldgutter.css'
import 'codemirror/addon/fold/brace-fold.js'
import 'codemirror/addon/fold/comment-fold.js'
import 'codemirror/addon/fold/foldcode.js'
import 'codemirror/addon/fold/foldgutter.js'
import 'codemirror/addon/fold/indent-fold.js'
import 'codemirror/addon/fold/markdown-fold.js'
import 'codemirror/addon/fold/xml-fold.js'

import 'codemirror/theme/ayu-mirage.css'

export default {
    name: 'Home',
    components: {},
    data() {
        return {
            record: {
                title: "",
                content: "",
                author: "",
                lang: "plaintext",
                password: "",
                tags: [],
                lifecycle: 1,
                editable: false,
            },

            lang_opts: [
                {value: 'plaintext', text: 'Plain Text'},
                {value: 'c', text: 'C'},
                {value: 'cpp', text: 'C++'},
                {value: 'cs', text: 'C#'},
                {value: 'python', text: 'Python'},
                {value: 'go', text: 'Go'},
                {value: 'php', text: 'PHP'},
                {value: 'markdown', text: 'Markdown'},
                {value: 'java', text: 'Java'},
                {value: 'javascript', text: 'JavaScript'},
                {value: 'typescript', text: 'TypeScript'},
                {value: 'vbnet', text: 'VB.NET'},
                {value: 'xml', text: 'XML'},
                {value: 'html', text: 'HTML'},
                {value: 'perl', text: 'Perl'},
                {value: 'css', text: 'CSS'},
                {value: 'yaml', text: 'YAML'},
                {value: 'json', text: 'JSON / JSON with Comments'},
                {value: 'dart', text: 'Dart'},
                {value: 'bash', text: 'Bash'},
                {value: 'rust', text: 'Rust'},
                {value: 'cmake', text: 'CMake'},
                {value: 'coffeescript', text: 'CoffeeScript'},
                {value: 'css', text: 'CSS'},
                {value: 'd', text: 'D'},
                {value: 'delphi', text: 'Delphi'},
                {value: 'django', text: 'Django'},
                {value: 'dockerfile', text: 'Dockerfile'},
                {value: 'dos', text: 'DOS .bat'},
                {value: 'erlang', text: 'Erlang'},
                {value: 'fortran', text: 'Fortran'},
                {value: 'gradle', text: 'Gradle'},
                {value: 'groovy', text: 'Groovy'},
                {value: 'haskell', text: 'Haskell'},
                {value: 'julia', text: 'Julia'},
                {value: 'kotlin', text: 'Kotlin'},
                {value: 'less', text: 'Less'},
                {value: 'lisp', text: 'Lisp'},
                {value: 'llvm', text: 'LLVM IR'},
                {value: 'lua', text: 'Lua'},
                {value: 'makefile', text: 'Makefile'},
                {value: 'mathematica', text: 'Mathematica'},
                {value: 'matlab', text: 'Matlab'},
                {value: 'nginx', text: 'Nginx'},
                {value: 'objectivec', text: 'Objective-C'},
                {value: 'powershell', text: 'PowerShell'},
                {value: 'qml', text: 'QML'},
                {value: 'r', text: 'R'},
                {value: 'ruby', text: 'Ruby'},
                {value: 'scala', text: 'Scala'},
                {value: 'scss', text: 'SCSS'},
                {value: 'shell', text: 'Shell Session'},
                {value: 'sql', text: 'SQL (Structured Query Language)'},
                {value: 'stylus', text: 'Stylus'},
                {value: 'swift', text: 'Swift'},
                {value: 'tex', text: 'TeX'},
                {value: 'twig', text: 'Twig'},
                {value: 'vim', text: 'Vim Script'},
            ],
            need_password: false,
            tags_limit: 10,
            lifecycle_opts: [
                {value: 1, text: 'A day'},
                {value: 2, text: 'A week'},
                {value: 3, text: 'A month'},
                {value: 4, text: 'A year'},
            ],

            cmOptions: {
                tabSize: 4,
                autocorrect: true,//自动更正
                spellcheck: true,//拼写检查
                styleActiveLine: true, // 高亮选中行
                styleSelectedText: true,
                autoCloseBrackets: true, // 自动闭合
                matchTags: {bothTags: true},
                matchBrackets: true,
                lineWrapping: true,
                foldGutter: true,
                gutters: ['CodeMirror-lint-markers', 'CodeMirror-linenumbers', 'CodeMirror-foldgutter'],
                mode: {name: "javascript", json: true},
                theme: "ayu-mirage",
                lineNumbers: true,
            },

            pasteErrorDelay: 5,
            pasteErrorShow: 0,
            pasteErrorMessage: "",
        }
    },
    methods: {
        paste() {
            let header = {
                'content-type': 'application/json',
            }
            let payload = {
                "title": this.record.title,
                "content": this.record.content,
                "author": this.record.author,
                "lang": this.record.lang,
                "tags": this.record.tags,
                "lifecycle": this.record.lifecycle,
                "password": this.record.password,
                "editable": this.record.editable,
            }

            this.$ajax.post("/v1/add", payload, {headers: header}).then(response => {
                let data = response.data;
                if (data.code === 200) {
                    this.$router.push({name: 'Show', params: {sk: data.data.sk}})
                } else {
                    console.log("get response err: ", data);
                    this.showAlert(data.message);
                }
            }).catch(error => {
                console.log("get err: ", error.response);
                this.showAlert(error.response.data.message)
            })
        },

        countDownChanged(pasteErrorDelay) {
            this.pasteErrorShow = pasteErrorDelay;
        },
        showAlert(errMsg) {
            this.pasteErrorMessage = errMsg;
            this.pasteErrorShow = this.pasteErrorDelay;
        }
    },
    computed: {
        codemirror() {
            return this.$refs.cmEditor.codemirror;
        },
    },
}
</script>

<style>
textarea {
    /*tab-size: 4;*/
    text-align: left;
}

.input {
    text-align: left;
}

.CodeMirror, .vue-codemirror {
    height: 512px;
    text-align: left !important;
    font-family: "Fira Code", "Courier New", Courier, monospace !important;
}

.form-group small {
    text-align: left;
}

.b-form-tags.focus {
    color: #1ddd86 !important;
    background-color: #28293e;
}

.custom-switch .custom-control-label::before {
    background-color: #9ea9b6;
}

#custom_tags___input__ {
    color: wheat;
}

#custom_tags input::-webkit-input-placeholder {
    color: #A19FB9;
}
</style>
