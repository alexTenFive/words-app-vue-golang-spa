import Vue from 'vue';
import VueRouter from 'vue-router';

import TextPage from './components/TextPage';
import ResultsPage from './components/ResultsPage';

Vue.use(VueRouter);

export default new VueRouter({
    mode: 'history',
    routes: [
        { path: '/', component: TextPage },
        { path: '/results', component: ResultsPage },
    ]
});