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
                                :sort-by="quantity"
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
        headers: [
          {
            text: 'Word',
            align: 'left',
            sortable: false,
            value: 'word'
          },
          { text: 'Quantity', value: 'quantity', mustSort: true },
        ],
        words: []
      }
    },
    methods: {
      results() {
        axios.get('http://localhost:9090/api/results')
        .then(res => {
          for (let word in res.data) {
            this.words.push({ word: word, quantity: res.data[word] });
          }
        })
        .catch(error =>  {
            console.log(error);
        });
      }
    }
  }
</script>
