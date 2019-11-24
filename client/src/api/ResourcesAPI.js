import ActionsServer from '../flux/actions-server';
import axios from "axios";

export default {

    getList() {
        var config = {
            headers: {'Access-Control-Allow-Origin': '*'},
        };
        axios.get("http://localhost:8080/resources", config).then(res => {
                if (res.status === 200) {
                    const {data} = res;
                    ActionsServer.receiveResources(data);
                }

            }
        ).catch(err => console.log(err))
    }

}