import Constants from "./constants";
import AppDispatcher from "./dispatcher";
import ResourcesAPI from "../api/ResourcesAPI";
import TemplatesAPI from "../api/TemplatesAPI";
import PoliciesAPI from "../api/PoliciesAPI";
import UsageAPI from "../api/UsageAPI";
import UserAPI from "../api/UserAPI";

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

    getHostStats(withPolling, interval = 15000) {
        AppDispatcher.handleViewAction({
            actionType: Constants.GET_HOST_STATS,
        });
        if (withPolling)
            return UsageAPI.getHostStatsWithPolling(interval);
        UsageAPI.getHostStats();
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

    getUserDetails() {
        AppDispatcher.handleViewAction({
            actionType: Constants.GET_USER_DETAILS,
        });
        UserAPI.getDetails();
    }

    getPolicies() {
        AppDispatcher.handleViewAction({
            actionType: Constants.GET_POLICIES,
        });
        PoliciesAPI.getList();
    }
}


export default new Actions();