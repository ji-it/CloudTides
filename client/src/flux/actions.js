import Constants from "./constants";
import AppDispatcher from "./dispatcher";
import ResourcesAPI from "../api/ResourcesAPI";
import TemplatesAPI from "../api/TemplatesAPI";

class Actions {
    addResource(data) {
        AppDispatcher.handleViewAction({
            actionType: Constants.ADD_RESOURCE,
            data: data
        });
    }

    addTemplate(data) {
        AppDispatcher.handleViewAction({
            actionType: Constants.ADD_TEMPLATE,
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

    getTemplates() {
        AppDispatcher.handleViewAction({
            actionType: Constants.GET_TEMPLATES,
        });
        TemplatesAPI.getList();
    }
}


export default new Actions();