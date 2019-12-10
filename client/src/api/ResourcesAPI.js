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
    },

    getVMList() {
        const requestURL = devURL + "/api/resource/get_vm_details/";
        request(requestURL, {method: 'GET'})
            .then((response) => {
                if (response.status === true) {
                    const {results} = response;
                    ActionsServer.receiveVMS(results);
                }
            }).catch((err) => {
            console.log(err);
        });
    },

    getVMListWithPolling(interval) {
        return setInterval(() => this.getVMList(), interval);
    },

    getDetailedList() {
        const requestURL = devURL + "/api/resource/get_details/";
        request(requestURL, {method: 'GET'})
            .then((response) => {
                if (response.status === true) {
                    const {results} = response;
                    ActionsServer.receiveDetailedResources(results);
                }
            }).catch((err) => {
            console.log(err);
        });
    },

    getDetailedListWithPolling(interval) {
        return setInterval(() => this.getDetailedList(), interval);
    },

    getOverview() {
        const requestURL = devURL + "/api/resource/overview/";
        request(requestURL, {method: 'GET'})
            .then((response) => {
                if (response.status === true) {
                    const {results} = response;
                    ActionsServer.receiveOverview(results);
                }
            }).catch((err) => {
            console.log(err);
        });
    },

    getOverviewWithPolling(interval) {
        return setInterval(() => this.getOverview(), interval);
    }
}