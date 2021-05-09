<template>
    <div class="archive">
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
                            <router-link :to='{"name": "Home"}'
                                         class="btn btn-success btn-shadow px-3 my-2 ml-0 text-left">
                                Home
                            </router-link>
                        </div>

                        <div class="text-darkgrey text-mono my-2">It is intended to be used directly by humans.</div>
                    </b-col>
                    <b-col xl="3" lg="3" md="2" sm="0"></b-col>
                </b-row>
            </b-container>
        </div>

        <b-container fluid="" class="content">
            <b-row>
                <b-col xl="3" lg="3" md="2" sm="0"></b-col>
                <b-col xl="6" lg="6" md="8" sm="12">
                    <b-table striped hover :items="items" :fields="fields">
                        <template #cell(title)="data">
                            <a class="text-info" :href="/s/+data.item.sk" target="_blank">{{ data.value }}</a>
                        </template>
                    </b-table>
                </b-col>
                <b-col xl="3" lg="3" md="2" sm="0">
                </b-col>
            </b-row>

        </b-container>

        <div class="container py-5">
            <p>Powered by <a href="https://xuthus.cc" target="_blank" class="text-danger">xuthus</a>, it's open source
                :) - <a
                    href="https://github.com/xuthus5/pbx" target="_blank" class="text-danger">pbx</a></p>
        </div>
    </div>
</template>

<script>
export default {
    name: 'Archive',
    data() {
        return {
            fields: [
                {key: 'title', label: 'TITLE'},
                {key: 'sk', label: 'POSTED'},
                {key: 'lang', label: 'SYNTAX'},
            ],
            items: []
        }
    },
    methods: {},
    created() {
        let header = {
            'content-type': 'application/json',
        }

        this.$ajax.get("/v1/archive", {headers: header}).then(response => {
            let data = response.data;
            if (data.code === 200) {
                this.items = data.data;
                return;
            }
            this.flashError(data.message);
            console.log("get err: ", response);
        }).catch(error => {
            this.flashError('get err: ', error);
            console.log("get err: ", error.response);
        })
    }
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
