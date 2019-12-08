import ActionsServer from '../flux/actions-server';
import axios from "axios";
import {devURL} from "../utils/urls";
import request from "../utils/request";

export default {

    getList() {
        const requestURL = devURL + "/api/resource/list/";
        request(requestURL, {method: 'GET'})
            .then((response) => {
                if (response.status === true) {
                    const {results} = response;
                    ActionsServer.receiveResources(results);
                }
            }).catch((err) => {
            console.log(err);
        });
    },

    getListWithPolling(interval) {
        return setInterval(() => this.getList(), interval);
    }

}