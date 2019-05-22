<template>
    <v-app id="large_text">
        <v-content>
            <v-container fluid fill-height>
                <v-layout align-center justify-center>
                    <v-flex xs12 sm8 md8>
                        <v-card class="elevation-12">
                            <v-toolbar dark color="teal">
                                <v-toolbar-title justify-center>Input large text</v-toolbar-title>
                            </v-toolbar>
                            <v-card-text>
                                <v-textarea
                                    name="text"
                                    box
                                    label="Input large text"
                                    auto-grow
                                    value="Lorem ipsum dolor sit amet"
                                    v-model="text"
                                ></v-textarea>
                                <v-btn @click.prevent="send" color="primary">Send</v-btn>
                            </v-card-text>
                        </v-card>
                    </v-flex>
                </v-layout>
            </v-container>
        </v-content>
    </v-app>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
        text: ""
    };
  },
  methods: {
    send () {
        let data = JSON.stringify({ 
            text: this.text, 
            number: 12 
        });
        axios.post('http://localhost:9090/api/send', data)
        .then(res => {
            console.log(res);
        });
    }
  },
}
</script>