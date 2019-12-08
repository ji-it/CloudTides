import Constants from "./constants";
import AppDispatcher from "./dispatcher";

class ActionsServer {
    receiveResources(response) {
        AppDispatcher.handleServerAction({
            actionType: Constants.GET_RESOURCES_RESPONSE,
            response: response
        });
    }

    receiveTemplates(response) {
        AppDispatcher.handleServerAction({
            actionType: Constants.GET_TEMPLATES_RESPONSE,
            response: response
        });
    }

    receivePolicies(response) {
        AppDispatcher.handleServerAction({
            actionType: Constants.GET_POLICIES_RESPONSE,
            response: response
        });
    }
}


export default new ActionsServer();