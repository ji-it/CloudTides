import Constants from "./constants";
import getSidebarNavItems from "../data/sidebar-nav-items";
import AppDispatcher from "./dispatcher";
import {EventEmitter} from "events";
import {devURL} from "../utils/urls";
import request from "../utils/request";
import auth from "../utils/auth";

let _store = {
    menuVisible: false,
    navItems: getSidebarNavItems(),
    resources: [],
    templates: [],
    policies: [],
    overview: {},
    resourceDetails: [],
    vms: [],
    resourceIndex: null
};

class Store extends EventEmitter {
    constructor() {
        super();

        this.addResource = this.addResource.bind(this);
        this.toggleSidebar = this.toggleSidebar.bind(this);
        this.updateResource = this.updateResource.bind(this);
        this.addTemplate = this.addTemplate.bind(this);
        this.updateTemplates = this.updateTemplates.bind(this);
        this.addPolicy = this.addPolicy.bind(this);
        this.updatePolicies = this.updatePolicies.bind(this);
        AppDispatcher.register(this.registerActions.bind(this));
    }

    registerActions({action}) {
        switch (action.actionType) {
            case Constants.TOGGLE_SIDEBAR:
                this.toggleSidebar();
                break;
            case Constants.SWITCH_RESOURCE:
                this.switchResourceIndex(action.data);
                break;
            case Constants.ADD_RESOURCE:
                this.addResource(action.data);
                break;
            case Constants.GET_RESOURCES_RESPONSE:
                this.updateResource(action.response);
                break;
            case Constants.ADD_TEMPLATE:
                this.addTemplate(action.data);
                break;
            case Constants.GET_VMS_RESPONSE:
                this.updateVMS(action.response);
                break;
            case Constants.GET_TEMPLATES_RESPONSE:
                this.updateTemplates(action.response);
                break;
            case Constants.ADD_POLICY:
                this.addPolicy(action.data);
                break;
            case Constants.GET_POLICIES_RESPONSE:
                this.updatePolicies(action.response);
                break;
            case Constants.GET_OVERVIEW_RESPONSE:
                this.updateOverview(action.response);
                break;
            case Constants.GET_DETAILED_RESOURCES_RESPONSE:
                this.updateDetailedResource(action.response);
                break;
            default:
                return true;
        }
    }

    addChangeListener(callback) {
        this.on(Constants.CHANGE, callback);
    }

    removeChangeListener(callback) {
        this.removeListener(Constants.CHANGE, callback);
    }

    addResource(data) {
        //Write to database and then use promise to push to
        const endpoint = '/api/resource/add/';
        const requestURL = devURL + endpoint;
        request(requestURL, {method: 'POST', body: data})
            .then((response) => {
                if (response.status === true) {
                    _store.resources = response.results;
                    this.emit(Constants.CHANGE);
                }
            }).catch((err) => {
            console.log(err);
        });
    }

    updateResource(data) {
        _store.resources = [];
        data.map((item, idx) => {
            _store.resources.push(item);
        });
        this.emit(Constants.CHANGE);
    }

    updateVMS(data) {
        _store.vms = [];
        data.map((item, idx) => {
            _store.vms.push(item);
        });
        this.emit(Constants.CHANGE);
    }

    updateDetailedResource(data) {
        _store.resourceDetails = [];
        data.map((item, idx) => {
            _store.resourceDetails.push(item);
        });
        this.emit(Constants.CHANGE);
    }

    updateOverview(data) {
        _store.overview = data;
        this.emit(Constants.CHANGE);
    }

    addTemplate(data) {
        const endpoint = '/api/template/add/';
        const requestURL = devURL + endpoint;
        request(requestURL, {method: 'POST', body: data}, false)
            .then((response) => {
                if (response.status === true) {
                    _store.templates = response.results;
                    this.emit(Constants.CHANGE);
                }
            }).catch((err) => {
            console.log(err);
        });
    }

    updateTemplates(data) {
        _store.templates = [];
        data.map((item, idx) => {
            _store.templates.push(item);
        });
        this.emit(Constants.CHANGE);
    }

    addPolicy(data) {
        const endpoint = '/api/policy/add/';
        const requestURL = devURL + endpoint;
        request(requestURL, {method: 'POST', body: data})
            .then((response) => {
                if (response.status === true) {
                    _store.policies = response.results;
                    this.emit(Constants.CHANGE);
                }
            }).catch((err) => {
            console.log(err);
        });
    }

    updatePolicies(data) {
        _store.policies = [];
        data.map((item, idx) => {
            _store.policies.push(item);
        });
        this.emit(Constants.CHANGE);
    }

    toggleSidebar() {
        _store.menuVisible = !_store.menuVisible;
        this.emit(Constants.CHANGE);
    }

    switchResourceIndex(data) {
        _store.resourceIndex = data;
        this.emit(Constants.CHANGE);
    }

    getMenuState() {
        return _store.menuVisible;
    }

    getResourceIndex() {
        return _store.resourceIndex
    }

    getSidebarItems() {
        return _store.navItems;
    }

    getResourceTableData() {
        return _store.resources;
    }

     getVMTableData() {
        return _store.vms;
    }

    getDetailedResourceTableData() {
        return _store.resourceDetails;
    }

    getTemplatesTableData() {
        return _store.templates;
    }

    getPoliciesTableData() {
        return _store.policies;
    }

    getOverviewData() {
        return _store.overview;
    }
}

export default new Store();
