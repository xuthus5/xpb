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
                            <router-link :to='{"name": "Home"}' class="btn btn-success btn-shadow px-3 my-2 ml-0 text-left">
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
                <b-col xl="3" lg="3" md="2" sm="0"></b-col>
                <b-col xl="6" lg="6" md="8" sm="12">
                    <b-card :title="record.title" :sub-title="getSubInfo()">

                        <b-card-text>
                            <pre v-highlightjs="record.content">
                                <code class="record.lang mt-0"></code>
                            </pre>
                            <p class="text-left text-danger mt-0 pt-0">This paste expires on 2021-03-29.</p>
                        </b-card-text>

                        <a :href="getRawContent()" class="card-link">raw content</a>
                    </b-card>
                </b-col>
                <b-col xl="3" lg="3" md="2" sm="0"></b-col>
            </b-row>
        </b-container>

        <div class="container py-5">
            <p>Powered by <a href="https://xuthus.cc">xuthus</a>, it's open source :) - <a
                href="https://github.com/xuthus5/pastebin">pastebin</a></p>
        </div>
    </div>
</template>

<script>
export default {
    name: 'Show',
    data() {
        return {
            record: {},
        }
    },
    methods: {
        getSubInfo() {
            return "Paste from " + this.record.author + " at "+ this.set_time(this.record.created_at);
        },
        getRawContent() {
            let sk = this.$route.params.sk;
            if (sk === "") {

            }

            return 'http://192.168.3.3:8081/v1/get?format=raw&sk='+sk;
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
            return  year + '-' + month + '-' + day + ' ' + hours + ':' + minutes;
        }
    },
    created() {
        let sk = this.$route.params.sk;
        if (sk === "") {

        }

        let header = {
            'content-type': 'application/json',
        }

        this.$ajax.get("http://localhost:8081/v1/get?sk=" + sk, {headers: header}).then(response => {
            let data = response.data;
            console.log("get data: ", data);
            if (data.code === 200) {
                this.record = data.data;
            } else {
                console.log("get response err: ", data);
            }
        }).catch(error => {
            console.log("get err: ", error);
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