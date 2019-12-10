import React from "react";
import {Nav, Button, Container} from "shards-react";
import ChildSidebarNavItem from "./ChildSidebarNavItem";
import TreeChildbarNavItem from "./TreeChildbarNavItem";
import Store from "../../../flux/store";
import {Actions} from "../../../flux";
import ModalAddResource from "../../common/ModalAddResource";
import SuccessNotificationModal from "../../common/SuccessNotificationModal";

class RightSidebarNavItems extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            navItems: Store.getSidebarItems(),
            resourceDetails: Store.getDetailedResourceTableData(),
            policies: Store.getPoliciesTableData(),
        };

        this.onChange = this.onChange.bind(this);
    }

    componentWillMount() {
        Store.addChangeListener(this.onChange);
    }

    componentWillUnmount() {
        Store.removeChangeListener(this.onChange);
    }

    componentDidMount() {
        Actions.getPolicies();
    }

    onChange() {
        this.setState({
            ...this.state,
            navItems: Store.getSidebarItems(),
            resourceDetails: Store.getDetailedResourceTableData(),
            policies: Store.getPoliciesTableData(),
        });
    }

    toggleModal = state => {
        this.setState({
            [state]: !this.state[state]
        });
    };

    myCallBack = (data, showNotif) => {
        this.toggleModal(data);
        if (showNotif)
            this.successmodal.toggleModal("notificationModal");
    };

    render() {
        const {navItems: items, resourceDetails: resData, policies: policiesData} = this.state;
        const selectedItem = items.filter(item => item.to === this.props.location.pathname).shift();
        const isManageResources = (selectedItem.to === "/manage-resources");
        return (
            <div className="child-nav-wrapper">
                <h6 className="ml-10 mt-10 font-weight-bold text-black">{selectedItem.name}</h6>
                <hr className="mx-10"/>
                {/*<Button squared className="w-100 mb-4"*/}
                {/*        onClick={() => this.toggleModal("addModal")}*/}
                {/*>*/}
                {/*    <i className="fa fa-plus mr-1"></i>*/}
                {/*    Add Resources*/}
                {/*</Button>*/}
                {/*<ModalAddResource policiesData={policiesData} onExit={this.myCallBack}*/}
                {/*                  toggleState={this.state.addModal}/>*/}
                {/*<SuccessNotificationModal onRef={ref => (this.successmodal = ref)}/>*/}
                <Nav className="nav--no-borders flex-column">
                    {isManageResources ? <TreeChildbarNavItem data={resData}/> :
                        selectedItem.children && selectedItem.children.map((item, idx) => {
                            return <ChildSidebarNavItem key={idx} item={item}/>
                        })}
                </Nav>
            </div>
        )
    }
}

export default RightSidebarNavItems;
