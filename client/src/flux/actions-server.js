import Constants from "./constants";
import AppDispatcher from "./dispatcher";

class ActionsServer {
    receiveResources(response) {
        AppDispatcher.handleServerAction({
            actionType: Constants.GET_RESOURCES_RESPONSE,
            response: response
        });
    }
}


export default new ActionsServer();