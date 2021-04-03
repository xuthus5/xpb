<template>
    <div class="show">
        <div class="jumbotron bg-transparent mb-0 radius-0 custom-header">
            <div class="container-fluid">
                <div class="row">
                    <div class="col-sm-0 col-md-2 col-lg-3 col-xl-3"></div>
                    <div class="col-sm-12 col-md-8 col-lg-6 col-xl-6"><h1 class="display-2">Paste Bin<span
                        class="vim-caret">x</span></h1>
                        <div class="lead mb-3 text-mono text-success">Intended for use as a short-term exchange of
                            pasted information between parties.
                        </div>
                        <div class="text-mono">
                            <router-link :to='{"name": "Home"}'
                                         class="btn btn-success btn-shadow px-3 my-2 ml-0 text-left">
                                Home
                            </router-link>
                        </div>
                        <div class="text-darkgrey text-mono my-2">It is intended to be used directly by humans.</div>
                    </div>
                    <div class="col-sm-0 col-md-2 col-lg-3 col-xl-3"></div>
                </div>
            </div>
        </div>

        <b-container fluid="" class="mt-2">
            <b-row>
                <b-col xl="2" lg="2" md="2" sm="0"></b-col>
                <b-col xl="8" lg="8" md="8" sm="12">
                    <b-card :title="record.title" :sub-title="getSubInfo()">

                        <b-card-text>
                            <pre v-highlightjs="record.content">
                                <code class="record.lang mt-0"></code>
                            </pre>
                            <div>
                                <span class="text-left text-danger mt-0 pt-0" v-if="infoNeedShow">{{ getExpiredInfo() }}</span>
                                <span class="text-right float-right text-danger mt-0 pt-0" v-if="infoNeedShow">
                                    <b-badge variant="success" v-for="(item, i) in record.tags" class="mr-1">{{item}}</b-badge>
                                </span>
                            </div>
                        </b-card-text>

                        <a :href="getEditLink()" class="card-link btn btn-secondary mr-2" target="_blank" v-if="record.editable">Public Edit</a>
                        <a :href="getRawLink()" class="card-link btn btn-primary" target="_blank" v-if="infoNeedShow">Raw Content</a>

                        <b-container fluid="">
                            <b-row>
                                <b-col sm="1" md="2" lg="4" xl="4"></b-col>
                                <b-col sm="10" md="8" lg="4" xl="4">
                                    <b-form inline class="text-center" v-if="needPassword">
                                        <b-form-input type="password" class="form-control" placeholder="Password"
                                                      v-model="record.password"></b-form-input>

                                        <button type="button" class="btn btn-primary ml-2" @click="showRecord">Go
                                        </button>
                                    </b-form>
                                </b-col>
                                <b-col sm="1" md="2" lg="4" xl="4"></b-col>
                            </b-row>
                        </b-container>

                    </b-card>
                </b-col>
                <b-col xl="2" lg="2" md="2" sm="0"></b-col>
            </b-row>
        </b-container>

        <div class="container py-5">
            <p>Powered by <a href="https://xuthus.cc" target="_blank" class="text-danger">xuthus</a>, it's open source
                :) - <a
                    href="https://github.com/xuthus5/xpb" target="_blank" class="text-danger">xpb</a></p>
        </div>
    </div>
</template>

<script>
export default {
    name: 'Show',
    data() {
        return {
            record: {},
            sk: '',
            infoNeedShow: false,
            needPassword: false
        }
    },
    methods: {
        getSubInfo() {
            if (this.record === null) return
            if (this.record === {}) return
            if (this.record.created_at === undefined) return

            if (this.record.author === undefined)
                this.record.author = 'Unknown'

            this.infoNeedShow = true;
            return "Paste from " + this.record.author + " at " + this.set_time(this.record.created_at);
        },
        getExpiredInfo() {
            return 'This paste expires on ' + this.set_time(this.record.expired_at) + ' .'
        },
        getRawLink() {
            return '/raw/' + this.sk;
        },
        getEditLink() {
            return '/edit/' + this.sk;
        },

        set_time(str) {
            let n = parseInt(str) * 1000;
            let D = new Date(n);
            let year = D.getFullYear();//四位数年份

            let month = D.getMonth() + 1;//月份(0-11),0为一月份
            month = month < 10 ? ('0' + month) : month;

            let day = D.getDate();//月的某一天(1-31)
            day = day < 10 ? ('0' + day) : day;

            let hours = D.getHours();//小时(0-23)
            hours = hours < 10 ? ('0' + hours) : hours;

            let minutes = D.getMinutes();//分钟(0-59)
            minutes = minutes < 10 ? ('0' + minutes) : minutes;

            // var seconds = D.getSeconds();//秒(0-59)
            // seconds = seconds<10?('0'+seconds):seconds;
            // var week = D.getDay();//周几(0-6),0为周日
            // var weekArr = ['周日','周一','周二','周三','周四','周五','周六'];
            return year + '-' + month + '-' + day + ' ' + hours + ':' + minutes;
        },

        showRecord() {
            let header = {
                'content-type': 'application/json',
            }

            this.$ajax.get("/v1/get?sk=" + this.sk + '&password=' + this.record.password, {headers: header}).then(response => {
                let data = response.data;
                if (data.code === 200) {
                    this.needPassword = false;
                    this.record = data.data;
                } else {
                    this.record = {content: data.message};
                    console.log("get response err: ", data);
                }
            }).catch(error => {
                var data = error.response.data;
                this.record = {content: data.message};
                if (data.code === 4002) {
                    this.needPassword = true;
                }
                console.log("get err: ", error.response);
            })
        }
    },
    created() {
        let sk = this.$route.params.sk;
        if (sk === "") {
            this.record = {content: "Record Not Found"};
            return
        }
        this.sk = sk;

        let header = {
            'content-type': 'application/json',
        }

        this.$ajax.get("/v1/get?sk=" + sk, {headers: header}).then(response => {
            let data = response.data;
            if (data.code === 200) {
                this.record = data.data;
            } else {
                this.record = {content: data.message};
                console.log("get response err: ", data);
            }
        }).catch(error => {
            var data = error.response.data;
            this.record = {content: data.message};
            if (data.code === 4002) {
                // need password to show
                this.needPassword = true;
            }
            console.log("get err: ", error.response);
        })
    }
}
</script>

<style scoped>
body {
    background-color: white;
}

.custom-header {
    background-color: #0c0d16 !important;
}

.card {
    background-color: white;
}

.card-title {
    color: black !important;
}

.card-subtitle {
    margin-top: 20px;
    color: #dd4814 !important;
}

.card-text {
    text-align: left;
}

pre {
    font-size: 1rem;
}

</style>