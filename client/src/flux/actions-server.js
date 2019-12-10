import Constants from "./constants";
import AppDispatcher from "./dispatcher";

class ActionsServer {
    receiveResources(response) {
        AppDispatcher.handleServerAction({
            actionType: Constants.GET_RESOURCES_RESPONSE,
            response: response
        });
    }

    receiveDetailedResources(response) {
        AppDispatcher.handleServerAction({
            actionType: Constants.GET_DETAILED_RESOURCES_RESPONSE,
            response: response
        });
    }

    receiveVMS(response) {
        AppDispatcher.handleServerAction({
            actionType: Constants.GET_VMS_RESPONSE,
            response: response
        });
    }


    receiveOverview(response) {
        AppDispatcher.handleServerAction({
            actionType: Constants.GET_OVERVIEW_RESPONSE,
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