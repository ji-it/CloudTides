import ActionsServer from '../flux/actions-server';
import axios from "axios";
import {devURL} from "../utils/urls";
import request from "../utils/request";

export default {

    getDetails() {
        const requestURL = devURL + "/api/users/get_details/";
        request(requestURL, {method: 'GET'})
            .then((response) => {
                if (response.status === true) {
                    const {results} = response;
                    ActionsServer.receiveUserDetails(results);
                }
            }).catch((err) => {
            console.log(err);
        });
    },

}