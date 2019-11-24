import Constants from "./constants";
import getSidebarNavItems from "../data/sidebar-nav-items";
import AppDispatcher from "./dispatcher";
import {EventEmitter} from "events";

let _store = {
    menuVisible: false,
    navItems: getSidebarNavItems(),
    resources: [
        ["New York Datacenter", "Idle", "192.168.0.1", 0.40, "6/10GB", "10/25GB", "20", "SETI@home", true,],
        ["LA Datacenter", "Contributing", "192.168.0.1", 0.60, "6/10GB", "10/25GB", "20", "SETI@home", false,],
    ],
};

class Store extends EventEmitter {
    constructor() {
        super();

        this.addResource = this.addResource.bind(this);
        this.toggleSidebar = this.toggleSidebar.bind(this);
        this.updateResource = this.updateResource.bind(this);
        AppDispatcher.register(this.registerActions.bind(this));
    }

    registerActions({action}) {
        switch (action.actionType) {
            case Constants.TOGGLE_SIDEBAR:
                this.toggleSidebar();
                break;
            case Constants.ADD_RESOURCE:
                this.addResource(action.data);
                break;
            case Constants.GET_RESOURCES_RESPONSE:
                const results = action.response.results;
                this.updateResource(results);
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
        console.log(data)
        //Write to database and then use promise to push to 
        _store.resources.push(data);
        this.emit(Constants.CHANGE);
    }

    updateResource(data) {
        data.map((item, idx) => {
            _store.resources = [];
            _store.resources.push(item);
        });
        console.log(_store.resources)
        this.emit(Constants.CHANGE);
    }

    toggleSidebar() {
        _store.menuVisible = !_store.menuVisible;
        this.emit(Constants.CHANGE);
    }

    getMenuState() {
        return _store.menuVisible;
    }

    getSidebarItems() {
        return _store.navItems;
    }

    getResourceTableData() {
        return _store.resources;
    }
};

export default new Store();
