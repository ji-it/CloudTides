import ActionsServer from '../flux/actions-server';
import axios from "axios";
import {devURL} from "../utils/urls";
import request from "../utils/request";

export default {

    getHostStats() {
        const requestURL = devURL + "/api/usage/getusage/";
        request(requestURL, {method: 'GET'})
            .then((response) => {
                if (response.status === true) {
                    const {results} = response;
                    ActionsServer.receiveHostStats(results);
                }
            }).catch((err) => {
            console.log(err);
        });
    },

    getHostStatsWithPolling(interval) {
        return setInterval(() => this.getHostStats(), interval);
    },
}