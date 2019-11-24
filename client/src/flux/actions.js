import Constants from "./constants";
import AppDispatcher from "./dispatcher";
import ResourcesAPI from "../api/ResourcesAPI";

class Actions {
    addResource(data) {
        AppDispatcher.handleViewAction({
            actionType: Constants.ADD_RESOURCE,
            data: data
        });
    }

    toggleMenu() {
        AppDispatcher.handleViewAction({
            actionType: Constants.TOGGLE_SIDEBAR,
        });
    }

    getResources() {
        AppDispatcher.handleViewAction({
            actionType: Constants.GET_RESOURCES,
        });
        ResourcesAPI.getList();
    }
}


export default new Actions();