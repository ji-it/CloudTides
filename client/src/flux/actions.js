import Constants from "./constants";
import AppDispatcher from "./dispatcher";
import ResourcesAPI from "../api/ResourcesAPI";
import TemplatesAPI from "../api/TemplatesAPI";
import PoliciesAPI from "../api/PoliciesAPI";

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

    addPolicy(data) {
        AppDispatcher.handleViewAction({
            actionType: Constants.ADD_POLICY,
            data: data
        });
    }

    toggleMenu() {
        AppDispatcher.handleViewAction({
            actionType: Constants.TOGGLE_SIDEBAR,
        });
    }

    switchResource(data) {
        AppDispatcher.handleViewAction({
            actionType: Constants.SWITCH_RESOURCE,
            data: data
        });
    }

    getResources(withPolling, interval = 15000) {
        AppDispatcher.handleViewAction({
            actionType: Constants.GET_RESOURCES,
        });
        if (withPolling)
            return ResourcesAPI.getListWithPolling(interval);
        ResourcesAPI.getList();
    }

    getVMS(withPolling, interval = 15000) {
        AppDispatcher.handleViewAction({
            actionType: Constants.GET_VMS,
        });
        if (withPolling)
            return ResourcesAPI.getVMListWithPolling(interval);
        ResourcesAPI.getVMList();
    }

    getDetailedResources(withPolling, interval = 15000) {
        AppDispatcher.handleViewAction({
            actionType: Constants.GET_DETAILED_RESOURCES,
        });
        if (withPolling)
            return ResourcesAPI.getDetailedListWithPolling(interval);
        ResourcesAPI.getDetailedList();
    }

    getOverview(withPolling, interval = 15000) {
        AppDispatcher.handleViewAction({
            actionType: Constants.GET_OVERVIEW,
        });
        if (withPolling)
            return ResourcesAPI.getOverviewWithPolling(interval);
        ResourcesAPI.getOverview();
    }

    getTemplates() {
        AppDispatcher.handleViewAction({
            actionType: Constants.GET_TEMPLATES,
        });
        TemplatesAPI.getList();
    }

    getPolicies() {
        AppDispatcher.handleViewAction({
            actionType: Constants.GET_POLICIES,
        });
        PoliciesAPI.getList();
    }
}


export default new Actions();