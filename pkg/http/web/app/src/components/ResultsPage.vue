<template>
<v-app id="results">
        <v-content>
            <v-container fluid fill-height>
                <v-layout align-center justify-center>
                    <v-flex xs12 sm8 md4>
                        <v-card class="elevation-12">
                            <v-toolbar dark color="teal">
                                <v-toolbar-title justify-center>Results</v-toolbar-title>
                            </v-toolbar>
                            <v-data-table
                                :headers="headers"
                                :items="words"
                                :pagination.sync="pagination"
                            >
                                <v-toolbar dark color="teal">
                                                            <v-toolbar-title justify-center>Input large text</v-toolbar-title>
                                                        </v-toolbar>
                                <template v-slot:items="props">
                                <td class="text-xs-left">{{ props.item.word }}</td>
                                <td class="text-xs-left">{{ props.item.quantity }}</td>
                                </template>
                            </v-data-table>
                        <v-btn @click.prevent="results" color="primary">Update</v-btn>

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
    data () {
      return {
        pagination: {
          sortBy: 'quantity'
        },
        headers: [
          {
            text: 'Word',
            align: 'left',
            sortable: false,
            value: 'word'
          },
          { text: 'Quantity', value: 'quantity', mustSort: true },
        ],
        words: [
          {
            word: '',
            quantity: 0
          }
        ]
      }
    },
    methods: {
      results() {
        let that = this;
        
        axios.get('/api/results')
        .then(res => {
          that.words = [];
          for (let word in res.data) {
            that.words.push({ word: word, quantity: res.data[word] });
          }
        });
      }
    },
    created() {
      this.pagination.descending = true;
      this.results();
    }
  }
</script>
